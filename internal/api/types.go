package api

import "time"

// route: /payments
type PaymentRequest struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

// route: /payments-summary
type PaymentSummaryRequest struct {
	From time.Time
	To   time.Time
}

type PaymentSummary struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type PaymentSummaryResponse struct {
	Default  PaymentSummary `json:"default"`
	Fallback PaymentSummary `json:"fallback"`
}
