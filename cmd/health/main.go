package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/svaan1/rinha-de-backend-2025/internal/payments"
)

const (
	defaultPaymentProcessorEndpoint  = "http://payment-processor-default:8080/payments/service-health"
	fallbackPaymentProcessorEndpoint = "http://payment-processor-fallback:8080/payments/service-health"
)

var (
	healthStatus   payments.PaymentProcessorHealthCheck
	healthStatusMu sync.RWMutex
	httpClient     = &http.Client{Timeout: 10 * time.Second}
)

func getHealth(endpoint string) *payments.PaymentProcessorHealth {
	resp, err := httpClient.Get(endpoint)
	if err != nil {
		log.Printf("Error fetching health from %s: %v", endpoint, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: Received non-200 status code %d from %s", resp.StatusCode, endpoint)
		return nil
	}

	var health payments.PaymentProcessorHealth
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		log.Printf("Error decoding response from %s: %v", endpoint, err)
		return nil
	}

	return &health
}

func healthTicker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	performCheck := func() {
		var wg sync.WaitGroup
		wg.Add(2)

		var defaultHealth, fallbackHealth *payments.PaymentProcessorHealth

		go func() {
			defer wg.Done()
			defaultHealth = getHealth(defaultPaymentProcessorEndpoint)
		}()

		go func() {
			defer wg.Done()
			fallbackHealth = getHealth(fallbackPaymentProcessorEndpoint)
		}()

		wg.Wait()

		healthStatusMu.Lock()
		healthStatus.Default = *defaultHealth
		healthStatus.Fallback = *fallbackHealth
		healthStatusMu.Unlock()
	}

	performCheck()

	for range ticker.C {
		performCheck()
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	healthStatusMu.RLock()
	defer healthStatusMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(healthStatus); err != nil {
		log.Printf("Error writing health response: %v", err)
		http.Error(w, "Failed to encode health status", http.StatusInternalServerError)
	}
}

func main() {
	go healthTicker()

	http.HandleFunc("/", healthCheckHandler)

	log.Print("Starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Failed to start server %v", err)
	}
}
