package payments

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
)

type PaymentRequest struct {
	CorrelationID string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}

type PaymentPayload struct {
	CorrelationID string `json:"correlationId"`
	Amount        string `json:"amount"`
	RequestedAt   string `json:"requestedAt"`
}

func ExecutePayment(ctx context.Context, payment PaymentRequest, pp *PaymentProcessor) error {
	requestedAt := time.Now().UTC()

	payload := PaymentPayload{
		CorrelationID: payment.CorrelationID,
		Amount:        fmt.Sprintf("%.2f", payment.Amount),
		RequestedAt:   requestedAt.Format(time.RFC3339Nano),
	}

	paymentKey := pp.Name + ":" + payment.CorrelationID
	fixedAmount := int64(math.Round(payment.Amount * 100))
	timestamp := requestedAt.UnixMicro()

	body, err := sonic.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(pp.Endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		pp.Status.Failing = true
		return fmt.Errorf("failed to post request, status, %d", resp.StatusCode)
	}

	err = globals.RedisClient.PersistPayment(ctx, paymentKey, fixedAmount, timestamp)
	if err != nil {
		return err
	}

	return nil
}

func PaymentTask(ctx context.Context, payment PaymentRequest) {
	for {
		pp, err := choosePaymentProcessor()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		err = ExecutePayment(ctx, payment, pp)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		return
	}
}
