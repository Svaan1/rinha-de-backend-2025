package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mborders/artifex"
	"github.com/svaan1/rinha-de-backend-2025/internal/queue"
)

const maxWorkers = 600
const maxQueue = 100000

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        maxWorkers * 2,
			MaxIdleConnsPerHost: maxWorkers * 2,
			IdleConnTimeout:     90 * time.Second,
			DisableKeepAlives:   false,
		},
	}

	d := artifex.NewDispatcher(maxWorkers, maxQueue)
	d.Start()

	task := func() {
		queue.ExecutePayment(queue.Payment{
			CorrelationID: uuid.NewString(),
			Amount:        19,
			RequestedAt:   time.Now(),
		}, *client)
	}

	for i := 1; i <= maxQueue; i++ {
		d.Dispatch(task)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
