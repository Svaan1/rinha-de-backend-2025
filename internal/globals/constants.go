package globals

const (
	DefaultPaymentProcessorEndpoint  = "http://payment-processor-default:8080/payments"
	FallbackPaymentProcessorEndpoint = "http://payment-processor-fallback:8080/payments"
	HealthCheckEndpoint              = "http://health:8080/"

	MaxWorkers  = 300
	MaxQueue    = 5000
	redisAddr   = "redis:6379"
	redisPasswd = ""
	redisDB     = 0
)
