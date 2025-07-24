package globals

import (
	"net/http"
	"time"

	"github.com/mborders/artifex"
	"github.com/svaan1/rinha-de-backend-2025/internal/redis"
)

var HttpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        MaxWorkers * 2,
		MaxIdleConnsPerHost: MaxWorkers * 2,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	},
}

var RedisClient = redis.New(redisAddr, redisPasswd, redisDB)

var QueueDispatcher = artifex.NewDispatcher(MaxWorkers, MaxQueue)
