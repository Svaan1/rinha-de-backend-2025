package payments

import (
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
	"github.com/valyala/fasthttp"
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
	// Execute the post request
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(globals.HealthCheckEndpoint)
	req.Header.SetMethod("GET")

	err := globals.HTTPClient.Do(req, resp)
	if err != nil {
		return err
	}

	body := resp.Body()

	// Handle non 200
	if resp.StatusCode() != 200 {
		return fmt.Errorf("health service returned non ok")
	}

	// Parse the response body
	var paymentProcessorHealthCheck PaymentProcessorHealthCheck
	if err = sonic.Unmarshal(body, &paymentProcessorHealthCheck); err != nil {
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

	if !FallbackPaymentProcessor.Status.Failing {
		return &FallbackPaymentProcessor, nil
	}

	return nil, fmt.Errorf("no payment processor available")
}
