package globals

import (
	"time"

	"github.com/mborders/artifex"
	"github.com/svaan1/rinha-de-backend-2025/internal/redis"
	"github.com/valyala/fasthttp"
)

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

var HTTPClient = &fasthttp.Client{
	MaxResponseBodySize: 4 * 1024 * 1024,
	ReadTimeout:         30 * time.Second,
	WriteTimeout:        30 * time.Second,

	Dial: (&fasthttp.TCPDialer{
		Concurrency: MaxWorkers * 2,
	}).Dial,

	MaxConnsPerHost:               MaxWorkers * 2,
	MaxIdleConnDuration:           90 * time.Second,
	MaxConnDuration:               0,
	MaxConnWaitTimeout:            5 * time.Second,
	DisablePathNormalizing:        true,
	DisableHeaderNamesNormalizing: true,
}

var RedisClient = redis.New(redisAddr, redisPasswd, redisDB)

var QueueDispatcher = artifex.NewDispatcher(MaxWorkers, MaxQueue)
