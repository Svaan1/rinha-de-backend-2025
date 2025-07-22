package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/mborders/artifex"
)

var success, fails int64

const maxWorkers = 700
const maxQueue = 100000

const endpoint = "http://payment-processor-default:8080/payments"

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

	paymentTask := func() {
		executePayment(client)
	}

	for i := 1; i <= maxQueue; i++ {
		d.Dispatch(paymentTask)
	}

	for (atomic.LoadInt64(&success) + atomic.LoadInt64(&fails)) < maxQueue {
		log.Printf("Progress -> Success: %d, Fails: %d", atomic.LoadInt64(&success), atomic.LoadInt64(&fails))
		time.Sleep(1 * time.Second)
	}

	log.Printf("Finished! Final Score -> Success: %d, Fails: %d", success, fails)
}

func executePayment(client *http.Client) {
	correlationID := uuid.New().String()
	requestedAt := time.Now().UTC().Format("2006-01-02T15:04:05")
	reqBody := fmt.Sprintf(`{"correlationId":"%s","amount":19,"requestedAt":"%s"}`, correlationID, requestedAt)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(reqBody))
	if err != nil {
		atomic.AddInt64(&fails, 1)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		atomic.AddInt64(&fails, 1)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		atomic.AddInt64(&fails, 1)
		return
	}

	atomic.AddInt64(&success, 1)
}
