package queue

import "log"

func NewClient() *QueueService {
	config := DefaultQueueConfig()
	client, err := NewQueueService(config)
	if err != nil {
		log.Fatalf("could not create queue client: %v", err)
	}

	return client
}
