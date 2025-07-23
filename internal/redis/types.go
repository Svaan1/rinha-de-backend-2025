package redis

type Payments struct {
	TotalRequests int64   `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type PaymentSummary struct {
	Default  Payments `json:"default"`
	Fallback Payments `json:"fallback"`
}
