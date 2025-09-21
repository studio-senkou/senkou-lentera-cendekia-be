package queue

import (
	"os"
	"strconv"
	"time"
)

func LoadConfigFromEnv() *QueueConfig {
	config := DefaultQueueConfig()

	if host := os.Getenv("REDIS_HOST"); host != "" {
		config.RedisHost = host
	}

	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		config.RedisPassword = password
	}

	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if db, err := strconv.Atoi(dbStr); err == nil {
			config.RedisDB = db
		}
	}

	if concurrencyStr := os.Getenv("QUEUE_CONCURRENCY"); concurrencyStr != "" {
		if concurrency, err := strconv.Atoi(concurrencyStr); err == nil {
			config.Concurrency = concurrency
		}
	}

	return config
}

type QueueManagerConfig struct {
	Queue  *QueueConfig
	Worker WorkerConfig
}

type WorkerConfig struct {
	Enabled           bool
	ShutdownTimeout   time.Duration
	HealthCheckPeriod time.Duration
}

func DefaultWorkerConfig() WorkerConfig {
	return WorkerConfig{
		Enabled:           true,
		ShutdownTimeout:   30 * time.Second,
		HealthCheckPeriod: 15 * time.Second,
	}
}

func LoadWorkerConfigFromEnv() WorkerConfig {
	config := DefaultWorkerConfig()

	if enabledStr := os.Getenv("QUEUE_WORKER_ENABLED"); enabledStr != "" {
		if enabled, err := strconv.ParseBool(enabledStr); err == nil {
			config.Enabled = enabled
		}
	}

	if timeoutStr := os.Getenv("QUEUE_SHUTDOWN_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			config.ShutdownTimeout = timeout
		}
	}

	if periodStr := os.Getenv("QUEUE_HEALTH_CHECK_PERIOD"); periodStr != "" {
		if period, err := time.ParseDuration(periodStr); err == nil {
			config.HealthCheckPeriod = period
		}
	}

	return config
}
