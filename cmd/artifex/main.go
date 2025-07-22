package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mborders/artifex"
)

var success = 0
var failed = 0

func main() {
	const maxWorkers = 500
	const maxQueue = 100000
	d := artifex.NewDispatcher(maxWorkers, maxQueue)
	d.Start()

	for i := 1; i <= maxQueue; i++ {
		log.Print(i)
		d.Dispatch(executePayment)
	}

	for {
		time.Sleep(10 * time.Millisecond)
	}
}

func executePayment() {
	client := &http.Client{Timeout: 1 * time.Second}

	correlationID := uuid.New().String()
	amount := 19

	reqBody := fmt.Sprintf(`{"correlationId":"%s","amount":%d}`, correlationID, amount)

	for {
		resp, err := client.Post("http://localhost:8001/payments", "application/json", strings.NewReader(reqBody))
		if err != nil {
			log.Printf("Error on request %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Request did not return OK %v", resp)
			time.Sleep(1 * time.Second)
			continue
		}

		break
	}

	success++
	time.Sleep(1 * time.Second)
}
