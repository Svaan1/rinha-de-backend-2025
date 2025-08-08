package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb redis.Client
}

func New(addr string, password string, DB int) RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return RedisClient{
		rdb: *rdb,
	}
}

func (r *RedisClient) PersistPayment(ctx context.Context, key string, amount int64, timestamp int64) error {
	keys := []string{key, "payments:by_time"}
	args := []interface{}{timestamp, amount}
	if err := createPaymentScript.Run(ctx, r.rdb, keys, args).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) GetPaymentSummary(ctx context.Context, fromTimestamp int64, toTimestamp int64) (PaymentSummary, error) {
	keys := []string{"payments:by_time"}
	args := []interface{}{fromTimestamp, toTimestamp}
	result, err := getPaymentSummaryScript.Run(ctx, r.rdb, keys, args).Result()
	if err != nil {
		return PaymentSummary{}, err
	}

	resultSlice := result.([]interface{})

	totalDefault := resultSlice[1].(int64)
	totalFallback := resultSlice[3].(int64)

	return PaymentSummary{
		Default: Payments{
			TotalRequests: resultSlice[0].(int64),
			TotalAmount:   float64(totalDefault) / 100,
		},
		Fallback: Payments{
			TotalRequests: resultSlice[2].(int64),
			TotalAmount:   float64(totalFallback) / 100,
		},
	}, nil
}

func (r *RedisClient) Purge(ctx context.Context) {
	r.rdb.FlushAll(ctx)
}
