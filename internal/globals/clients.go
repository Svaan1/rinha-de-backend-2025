package globals

import (
	"net/http"
	"time"

	"github.com/mborders/artifex"
	"github.com/svaan1/rinha-de-backend-2025/internal/redis"
)

const MaxWorkers = 600
const MaxQueue = 100000
const redisAddr = "redis:6379"

var HttpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        MaxWorkers * 2,
		MaxIdleConnsPerHost: MaxWorkers * 2,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	},
}

var RedisClient = redis.New(redisAddr, "", 0)

var QueueDispatcher = artifex.NewDispatcher(MaxWorkers, MaxQueue)
