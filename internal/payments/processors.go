package payments

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
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
	resp, err := http.Get(globals.HealthCheckEndpoint)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health service returned non ok")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var paymentProcessorHealthCheck PaymentProcessorHealthCheck
	if err = sonic.Unmarshal(body, &paymentProcessorHealthCheck); err != nil {
		return err
	}

	DefaultPaymentProcessor.Status = paymentProcessorHealthCheck.Default
	FallbackPaymentProcessor.Status = paymentProcessorHealthCheck.Fallback

	return nil
}

func choosePaymentProcessor() (*PaymentProcessor, error) {
	if !DefaultPaymentProcessor.Status.Failing {
		return &DefaultPaymentProcessor, nil
	}

	if !FallbackPaymentProcessor.Status.Failing {
		return &FallbackPaymentProcessor, nil
	}

	return nil, fmt.Errorf("no payment processor available")
}
