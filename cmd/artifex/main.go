package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mborders/artifex"
)

var success = 0
var fails = 0

func main() {
	// 700 working to 100k
	const maxWorkers = 600
	const maxQueue = 100000
	d := artifex.NewDispatcher(maxWorkers, maxQueue)
	d.Start()

	for i := 1; i <= maxQueue; i++ {
		log.Print(i)
		d.Dispatch(executePayment)
	}

	for {
		log.Print(success, fails)
		time.Sleep(1 * time.Second)

		if success == maxQueue {
			break
		}
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
			fails++
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("Request did not return OK %v", resp)
			time.Sleep(1 * time.Second)
			fails++
			continue
		}

		break
	}

	success++
	delay := rand.Intn(1000) + 1000 // a random delay (currently randing to 1501 + 1000 * ms works really well)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
