package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hibiken/asynq"
	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
)

type QueueService struct {
	client    *asynq.Client
	server    *asynq.Server
	scheduler *asynq.Scheduler
	inspector *asynq.Inspector
	handlers  map[string]asynq.HandlerFunc
}

type QueueConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	Concurrency   int
	Queues        map[string]int
}

func DefaultQueueConfig() *QueueConfig {
	host := app.GetEnv("REDIS_HOST", "localhost:6379")
	port := app.GetEnv("REDIS_PORT", "6379")
	password := app.GetEnv("REDIS_PASSWORD", "")
	db, _ := strconv.Atoi(app.GetEnv("REDIS_DB", "0"))

	return &QueueConfig{
		RedisHost:     host,
		RedisPort:     port,
		RedisPassword: password,
		RedisDB:       db,
		Concurrency:   10,
		Queues: map[string]int{
			"critical": 6,
			"high":     4,
			"default":  3,
			"low":      1,
		},
	}
}

func NewQueueService(config *QueueConfig) (*QueueService, error) {
	redisOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	}

	client := asynq.NewClient(redisOpt)

	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency:         config.Concurrency,
		Queues:              config.Queues,
		LogLevel:            asynq.InfoLevel,
		RetryDelayFunc:      asynq.DefaultRetryDelayFunc,
		HealthCheckInterval: 15 * time.Second,
	})

	scheduler := asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{
		LogLevel: asynq.InfoLevel,
	})

	inspector := asynq.NewInspector(redisOpt)

	return &QueueService{
		client:    client,
		server:    server,
		scheduler: scheduler,
		inspector: inspector,
		handlers:  make(map[string]asynq.HandlerFunc),
	}, nil
}

func (qs *QueueService) NewJobBuilder(taskName string) *JobBuilder {
	return NewJobBuilder(qs.client, taskName)
}

// RegisterHandler registers a handler function for a specific task name.
// The handler will be called when a task with the given taskName is processed.
//
// Example:
//
//	qs.RegisterHandler("email:send", func(ctx context.Context, task *asynq.Task) error {
//	    // Your email sending logic here
//	    return nil
//	})
func (qs *QueueService) RegisterHandler(taskName string, handler asynq.HandlerFunc) {
	qs.handlers[taskName] = handler
}

// RegisterHandlerFunc registers a handler function for a specific task name.
// This is a convenience method that wraps the provided function in asynq.HandlerFunc.
//
// Example:
//
//	qs.RegisterHandlerFunc("email:send", func(ctx context.Context, task *asynq.Task) error {
//	    var payload map[string]interface{}
//	    json.Unmarshal(task.Payload(), &payload)
//	    // Your logic here
//	    return nil
//	})
func (qs *QueueService) RegisterHandlerFunc(taskName string, handler func(context.Context, *asynq.Task) error) {
	qs.handlers[taskName] = asynq.HandlerFunc(handler)
}

// RegisterJSONHandler registers a handler that automatically unmarshals JSON payload.
// This is a convenience method for handlers that expect JSON payloads.
//
// Example:
//
//	qs.RegisterJSONHandler("email:send", func(ctx context.Context, payload map[string]interface{}) error {
//	    email := payload["email"].(string)
//	    // Your logic here
//	    return nil
//	})
func (qs *QueueService) RegisterJSONHandler(taskName string, handler func(context.Context, map[string]interface{}) error) {
	qs.handlers[taskName] = asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		var payload map[string]interface{}
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return fmt.Errorf("failed to unmarshal JSON payload for task %s: %w", taskName, err)
		}
		return handler(ctx, payload)
	})
}

// RegisterTypedHandler registers a handler with automatic JSON unmarshaling into a specific type.
// The payloadType should be a pointer to the struct you want to unmarshal into.
//
// Example:
//
//	type EmailPayload struct {
//	    Email string `json:"email"`
//	    Name  string `json:"name"`
//	}
//
//	qs.RegisterTypedHandler("email:send", &EmailPayload{}, func(ctx context.Context, payload interface{}) error {
//	    emailData := payload.(*EmailPayload)
//	    // Your logic here
//	    return nil
//	})
func (qs *QueueService) RegisterTypedHandler(taskName string, payloadType interface{}, handler func(context.Context, interface{}) error) {
	qs.handlers[taskName] = asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		// Create a new instance of the payload type
		payload := payloadType
		if err := json.Unmarshal(task.Payload(), payload); err != nil {
			return fmt.Errorf("failed to unmarshal typed payload for task %s: %w", taskName, err)
		}
		return handler(ctx, payload)
	})
}

func (qs *QueueService) Start() error {
	mux := asynq.NewServeMux()
	for taskName, handler := range qs.handlers {
		mux.HandleFunc(taskName, handler)
	}

	go func() {
		if err := qs.scheduler.Run(); err != nil {
			log.Printf("Failed to start scheduler: %v", err)
		}
	}()

	return qs.server.Run(mux)
}

func (qs *QueueService) Stop() {
	qs.scheduler.Shutdown()
	qs.server.Shutdown()
	if err := qs.client.Close(); err != nil {
		log.Printf("Failed to close queue client: %v", err)
	}
}

func (qs *QueueService) GetQueueNames() ([]string, error) {
	return qs.inspector.Queues()
}

func (qs *QueueService) GetQueueInfo(qname string) (*asynq.QueueInfo, error) {
	return qs.inspector.GetQueueInfo(qname)
}

func (qs *QueueService) GetTaskInfo(queue, taskID string) (*asynq.TaskInfo, error) {
	return qs.inspector.GetTaskInfo(queue, taskID)
}

func (qs *QueueService) SchedulePeriodicTask(cronspec, taskName string, payload map[string]interface{}, opts ...asynq.Option) (string, error) {
	task := qs.NewJobBuilder(taskName).WithPayload(payload)
	payloadBytes, err := task.marshalPayload()
	if err != nil {
		return "", err
	}

	asynqTask := asynq.NewTask(taskName, payloadBytes, opts...)

	entryID, err := qs.scheduler.Register(cronspec, asynqTask)
	if err != nil {
		return "", err
	}

	return entryID, nil
}

func (qs *QueueService) DeletePeriodicTask(entryID string) error {
	return qs.scheduler.Unregister(entryID)
}

// UnmarshalTaskPayload is a utility function to unmarshal task payload into a struct.
// This can be used within handler functions for custom payload parsing.
//
// Example:
//
//	type EmailPayload struct {
//	    Email string `json:"email"`
//	    Name  string `json:"name"`
//	}
//
//	var payload EmailPayload
//	if err := qs.UnmarshalTaskPayload(task, &payload); err != nil {
//	    return err
//	}
func (qs *QueueService) UnmarshalTaskPayload(task *asynq.Task, v interface{}) error {
	return json.Unmarshal(task.Payload(), v)
}

// ValidateRequiredFields is a utility function to validate required fields in a payload.
// This can be used within handler functions for payload validation.
//
// Example:
//
//	payload := map[string]interface{}{"email": "test@example.com"}
//	if err := qs.ValidateRequiredFields(payload, "email", "name"); err != nil {
//	    return err
//	}
func (qs *QueueService) ValidateRequiredFields(payload map[string]interface{}, fields ...string) error {
	for _, field := range fields {
		if _, exists := payload[field]; !exists {
			return fmt.Errorf("required field '%s' is missing from payload", field)
		}
		if payload[field] == nil {
			return fmt.Errorf("required field '%s' cannot be null", field)
		}
	}
	return nil
}

func (jb *JobBuilder) marshalPayload() ([]byte, error) {
	return jb.marshalPayloadInternal()
}

func (jb *JobBuilder) marshalPayloadInternal() ([]byte, error) {
	if len(jb.payload) == 0 {
		return []byte("{}"), nil
	}

	payloadBytes, err := json.Marshal(jb.payload)
	if err != nil {
		return nil, err
	}
	return payloadBytes, nil
}
