package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/svaan1/rinha-de-backend-2025/internal/queue"
)

func main() {
	const maxJobs = 1000
	const maxWorkers = 100

	q := queue.New(maxJobs, maxWorkers)
	q.Run()

	for i := 1; i <= maxJobs; i++ {
		job := queue.Job{
			CorrelationID: uuid.New().String(),
			Amount:        19,
		}
		q.Enqueue(job)
	}

	time.Sleep(5 * time.Second)
}
