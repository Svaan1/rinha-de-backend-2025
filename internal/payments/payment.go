package payments

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
)

type Payment struct {
	CorrelationID string    `json:"correlationId"`
	Amount        float64   `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
}

func usePaymentProcessor(pp *PaymentProcessor, reqBody string) error {
	// Create the request object
	req, err := http.NewRequest("POST", pp.Endpoint, strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the post request
	resp, err := globals.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle non 200
	if resp.StatusCode != http.StatusOK {
		pp.Status.Failing = true
		return fmt.Errorf("failed to post request, status %d", resp.StatusCode)
	}

	return nil
}

func ExecutePayment(payment Payment) {
	// Create the payload
	reqBody := fmt.Sprintf(`{"correlationId":"%s","amount":%.2f,"requestedAt":"%s"}`, payment.CorrelationID, payment.Amount, payment.RequestedAt.Format("2006-01-02T15:04:05.000Z"))

	log.Print(reqBody)

	var pp *PaymentProcessor
	var err error
	for {
		// Choose the appropriate payment processor else we wait
		pp, err = choosePaymentProcessor()
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		// Apply the payment using the chosen processor
		err = usePaymentProcessor(pp, reqBody)
		if err != nil {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		break
	}

	// Store the successful payment at redis
	paymentKey := pp.Name + ":" + payment.CorrelationID
	fixedAmount := int64(math.Round(payment.Amount * 100))
	timestamp := float64(payment.RequestedAt.UnixMicro())

	err = globals.RedisClient.CreatePayment(paymentKey, fixedAmount, timestamp)
	if err != nil {
		return
	}
}
