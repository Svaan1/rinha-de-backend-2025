package redis

type Payments struct {
	TotalRequests int64   `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type PaymentSummary struct {
	Default  Payments `json:"default"`
	Fallback Payments `json:"fallback"`
}

type ProcessorStatus struct {
	Failing         bool  `json:"failing"`
	MinResponseTime int64 `json:"minResponseTime"`
}

type ProcessorStatusSummary struct {
	Default  ProcessorStatus `json:"default"`
	Fallback ProcessorStatus `json:"fallback"`
}
