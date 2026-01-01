package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

type JobPriority int

const (
	PriorityLow JobPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

type JobOptions struct {
	Priority  JobPriority
	Delay     time.Duration
	MaxRetry  int
	Timeout   time.Duration
	Queue     string
	ProcessAt *time.Time
	ProcessIn *time.Duration
	Retention time.Duration
	UniqueKey string
}

type JobBuilder struct {
	client   *asynq.Client
	taskName string
	payload  map[string]interface{}
	options  JobOptions
}

func NewJobBuilder(client *asynq.Client, taskName string) *JobBuilder {
	return &JobBuilder{
		client:   client,
		taskName: taskName,
		payload:  make(map[string]interface{}),
		options: JobOptions{
			Priority:  PriorityNormal,
			MaxRetry:  3,
			Timeout:   30 * time.Second,
			Queue:     "default",
			Retention: 24 * time.Hour,
		},
	}
}

func (jb *JobBuilder) WithPayload(payload map[string]interface{}) *JobBuilder {
	jb.payload = payload
	return jb
}

func (jb *JobBuilder) WithData(key string, value interface{}) *JobBuilder {
	jb.payload[key] = value
	return jb
}

func (jb *JobBuilder) WithPriority(priority JobPriority) *JobBuilder {
	jb.options.Priority = priority
	return jb
}

func (jb *JobBuilder) WithDelay(delay time.Duration) *JobBuilder {
	jb.options.Delay = delay
	return jb
}

func (jb *JobBuilder) WithMaxRetry(maxRetry int) *JobBuilder {
	jb.options.MaxRetry = maxRetry
	return jb
}

func (jb *JobBuilder) WithTimeout(timeout time.Duration) *JobBuilder {
	jb.options.Timeout = timeout
	return jb
}

func (jb *JobBuilder) WithQueue(queue string) *JobBuilder {
	jb.options.Queue = queue
	return jb
}

func (jb *JobBuilder) WithProcessAt(processAt time.Time) *JobBuilder {
	jb.options.ProcessAt = &processAt
	return jb
}

func (jb *JobBuilder) WithProcessIn(duration time.Duration) *JobBuilder {
	jb.options.ProcessIn = &duration
	return jb
}

func (jb *JobBuilder) WithRetention(retention time.Duration) *JobBuilder {
	jb.options.Retention = retention
	return jb
}

func (jb *JobBuilder) WithUniqueKey(key string) *JobBuilder {
	jb.options.UniqueKey = key
	return jb
}

func (jb *JobBuilder) Enqueue(ctx context.Context) (*asynq.TaskInfo, error) {
	payloadBytes, err := json.Marshal(jb.payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal job payload: %w", err)
	}

	task := asynq.NewTask(jb.taskName, payloadBytes)

	opts := jb.buildAsynqOptions()

	info, err := jb.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue job %s: %w", jb.taskName, err)
	}

	return info, nil
}

func (jb *JobBuilder) buildAsynqOptions() []asynq.Option {
	var opts []asynq.Option

	opts = append(opts, asynq.MaxRetry(jb.options.MaxRetry))

	opts = append(opts, asynq.Timeout(jb.options.Timeout))

	queueName := jb.getQueueNameForPriority()
	opts = append(opts, asynq.Queue(queueName))

	if jb.options.ProcessAt != nil {
		opts = append(opts, asynq.ProcessAt(*jb.options.ProcessAt))
	} else if jb.options.ProcessIn != nil {
		opts = append(opts, asynq.ProcessIn(*jb.options.ProcessIn))
	} else if jb.options.Delay > 0 {
		opts = append(opts, asynq.ProcessIn(jb.options.Delay))
	}

	opts = append(opts, asynq.Retention(jb.options.Retention))

	if jb.options.UniqueKey != "" {
		opts = append(opts, asynq.Unique(time.Hour))
	}

	return opts
}

func (jb *JobBuilder) getQueueNameForPriority() string {
	if jb.options.Queue != "default" {
		return jb.options.Queue
	}

	switch jb.options.Priority {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "default"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return "default"
	}
}
