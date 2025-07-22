package redis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb redis.Client
	ctx context.Context
}

func New(addr string, password string, DB int) RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	ctx := context.Background()

	return RedisClient{
		rdb: *rdb,
		ctx: ctx,
	}
}

func (r *RedisClient) CreatePayment(key string, amount int64, timestamp float64) error {
	keys := []string{key, "payments:by_time"}
	args := []interface{}{timestamp, amount}
	if err := createPaymentScript.Run(r.ctx, r.rdb, keys, args).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisClient) GetProcessorStatusSummary() (ProcessorStatusSummary, error) {
	result, err := fetchPaymentProcessorsScript.Run(r.ctx, r.rdb, nil).Result()
	if err != nil {
		return ProcessorStatusSummary{}, err
	}

	resultSlice := result.([]interface{})

	defaultFailing := resultSlice[0].(string) == "1"
	defaultMinResponseTime, err := strconv.ParseInt(resultSlice[1].(string), 0, 64)
	if err != nil {
		return ProcessorStatusSummary{}, err
	}

	fallbackFailing := resultSlice[2].(string) == "1"
	fallbackMinResponseTime, err := strconv.ParseInt(resultSlice[3].(string), 0, 64)
	if err != nil {
		return ProcessorStatusSummary{}, err
	}

	return ProcessorStatusSummary{
		Default: ProcessorStatus{
			Failing:         defaultFailing,
			MinResponseTime: defaultMinResponseTime,
		},
		Fallback: ProcessorStatus{
			Failing:         fallbackFailing,
			MinResponseTime: fallbackMinResponseTime,
		},
	}, nil
}

func (r *RedisClient) GetPaymentSummary(fromTimestamp int64, toTimestamp int64) (PaymentSummary, error) {
	keys := []string{"payments:by_time"}
	args := []interface{}{fromTimestamp, toTimestamp}
	result, err := getPaymentSummaryScript.Run(r.ctx, r.rdb, keys, args).Result()
	if err != nil {
		return PaymentSummary{}, err
	}

	resultSlice := result.([]interface{})

	return PaymentSummary{
		Default: Payments{
			TotalRequests: resultSlice[0].(int64),
			TotalAmount:   float64(resultSlice[1].(int64) / 100),
		},
		Fallback: Payments{
			TotalRequests: resultSlice[2].(int64),
			TotalAmount:   float64(resultSlice[3].(int64) / 100),
		},
	}, nil
}

func (r *RedisClient) Purge() {
	purgeScript.Run(r.ctx, r.rdb, nil)
}

func (r *RedisClient) UpdateDefaultFailure() {
	r.rdb.HSet(r.ctx, "payment-processor", "failing", "1")
}

func (r *RedisClient) UpdateFallbackFailure() {
	r.rdb.HSet(r.ctx, "payment-processor-fallback", "failing", "1")
}
