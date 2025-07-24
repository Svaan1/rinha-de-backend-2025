package payments

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
)

type PaymentProcessorHealth struct {
	Failing         bool
	MinResponseTime int64
}

type PaymentProcessorHealthCheck struct {
	Default  PaymentProcessorHealth `json:"default"`
	Fallback PaymentProcessorHealth `json:"fallback"`
}

type PaymentProcessor struct {
	Name     string
	Endpoint string
	Status   PaymentProcessorHealth
}

var DefaultPaymentProcessor = PaymentProcessor{
	Name:     "default",
	Endpoint: globals.DefaultPaymentProcessorEndpoint,
	Status: PaymentProcessorHealth{
		Failing:         false,
		MinResponseTime: 0,
	},
}

var FallbackPaymentProcessor = PaymentProcessor{
	Name:     "fallback",
	Endpoint: globals.FallbackPaymentProcessorEndpoint,
	Status: PaymentProcessorHealth{
		Failing:         false,
		MinResponseTime: 0,
	},
}

func StartHealthCheckTicker() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := UpdateProcessorHealth(); err != nil {
			log.Printf("Health check failed: %v", err)
		}
	}
}

func UpdateProcessorHealth() error {
	// Fetch the health service
	resp, err := globals.HttpClient.Get(globals.HealthCheckEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// If non 200, return an error
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health service returned non ok")
	}

	// Parse the response body
	var paymentProcessorHealthCheck PaymentProcessorHealthCheck
	if err := json.NewDecoder(resp.Body).Decode(&paymentProcessorHealthCheck); err != nil {
		return err
	}

	// Update the object
	DefaultPaymentProcessor.Status = paymentProcessorHealthCheck.Default
	FallbackPaymentProcessor.Status = paymentProcessorHealthCheck.Fallback

	return nil
}

func choosePaymentProcessor() (*PaymentProcessor, error) {
	if !DefaultPaymentProcessor.Status.Failing {
		return &DefaultPaymentProcessor, nil
	}

	// if !FallbackPaymentProcessor.Status.Failing {
	// 	return &FallbackPaymentProcessor, nil
	// }

	return nil, fmt.Errorf("no payment processor available")
}
