package globals

import (
	"github.com/mborders/artifex"
	"github.com/svaan1/rinha-de-backend-2025/internal/redis"
)

const (
	DefaultPaymentProcessorEndpoint  = "http://payment-processor-default:8080/payments"
	FallbackPaymentProcessorEndpoint = "http://payment-processor-fallback:8080/payments"
	HealthCheckEndpoint              = "http://health:8080/"

	MaxWorkers  = 30
	MaxQueue    = 100_000
	redisAddr   = "redis:6379"
	redisPasswd = ""
	redisDB     = 0
)

var RedisClient = redis.New(redisAddr, redisPasswd, redisDB)

var QueueDispatcher = artifex.NewDispatcher(MaxWorkers, MaxQueue)
