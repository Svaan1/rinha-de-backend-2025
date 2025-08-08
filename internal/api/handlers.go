package api

import (
	"context"
	"math"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
	"github.com/svaan1/rinha-de-backend-2025/internal/payments"
)

func PaymentHandler(c *fiber.Ctx) error {
	body := c.Body()

	go func(body []byte) {
		var data PaymentRequest
		if err := sonic.Unmarshal(body, &data); err != nil {
			return
		}

		task := func() {
			payments.PaymentTask(context.Background(), payments.PaymentRequest{
				CorrelationID: data.CorrelationID,
				Amount:        data.Amount,
			})
		}

		globals.QueueDispatcher.Dispatch(task)
	}(body)

	return c.SendStatus(fiber.StatusCreated)
}

func PaymentSummaryHandler(c *fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")

	var fromTimestamp, toTimestamp int64 = 0, math.MaxInt64
	if from != "" {
		if t, err := time.Parse(time.RFC3339Nano, from); err == nil {
			fromTimestamp = t.UnixMicro()
		}
	}
	if to != "" {
		if t, err := time.Parse(time.RFC3339Nano, to); err == nil {
			toTimestamp = t.UnixMicro()
		}
	}

	result, err := globals.RedisClient.GetPaymentSummary(c.Context(), fromTimestamp, toTimestamp)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	jsonBytes, err := sonic.Marshal(result)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set("Content-Type", "application/json")
	return c.Send(jsonBytes)
}

func PurgePaymentsHandler(c *fiber.Ctx) error {
	globals.RedisClient.Purge(c.Context())
	return c.SendStatus(fiber.StatusOK)
}
