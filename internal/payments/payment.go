package payments

import (
	"fmt"
	"math"
	"time"

	"github.com/bytedance/sonic"
	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
	"github.com/valyala/fasthttp"
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

func ExecutePayment(payment PaymentRequest, pp *PaymentProcessor) error {
	requestedAt := time.Now().UTC()

	payload := PaymentPayload{
		CorrelationID: payment.CorrelationID,
		Amount:        fmt.Sprintf("%.2f", payment.Amount),
		RequestedAt:   requestedAt.Format("2006-01-02T15:04:05.000Z"),
	}

	paymentKey := pp.Name + ":" + payment.CorrelationID
	fixedAmount := int64(math.Round(payment.Amount * 100))
	timestamp := requestedAt.UnixMicro()

	body, err := sonic.Marshal(payload)
	if err != nil {
		return err
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(pp.Endpoint)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")

	req.SetBody(body)

	err = globals.HTTPClient.Do(req, resp)
	if err != nil {
		return err
	}

	_ = resp.Body()

	// Handle non 200
	if resp.StatusCode() != 200 {
		pp.Status.Failing = true
		return fmt.Errorf("failed to post request, status, %d", resp.StatusCode())
	}

	err = globals.RedisClient.CreatePayment(paymentKey, fixedAmount, timestamp)
	if err != nil {
		return err
	}

	return nil
}

func PaymentTask(payment PaymentRequest) {
	for {
		// Choose the appropriate payment processor else we wait
		pp, err := choosePaymentProcessor()
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// Apply the payment using the chosen processor
		err = ExecutePayment(payment, pp)
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		return
	}
}
