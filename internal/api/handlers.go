package api

import (
	"math"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
	"github.com/svaan1/rinha-de-backend-2025/internal/payments"
)

func PaymentHandler(c *fiber.Ctx) error {
	body := c.Body()

	go func() {
		var data PaymentRequest
		if err := sonic.Unmarshal(body, &data); err != nil {
			return
		}

		task := func() {
			payments.ExecutePayment(payments.Payment{
				CorrelationID: data.CorrelationID,
				Amount:        data.Amount,
				RequestedAt:   time.Now().UTC(),
			})
		}

		globals.QueueDispatcher.Dispatch(task)
	}()

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

	result, err := globals.RedisClient.GetPaymentSummary(fromTimestamp, toTimestamp)
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
	globals.RedisClient.Purge()
	return c.SendStatus(fiber.StatusOK)
}
