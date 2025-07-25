package globals

import (
	"time"

	"github.com/mborders/artifex"
	"github.com/svaan1/rinha-de-backend-2025/internal/redis"
	"github.com/valyala/fasthttp"
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
