package queue

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type PaymentProcessor struct {
	MinResponseTime int64
	Failing         bool
	Endpoint        string
	RedisKey        string
}

type Payment struct {
	CorrelationID string    `json:"correlationId"`
	Amount        int64     `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
}

func choosePaymentProcessor() (PaymentProcessor, error) {
	return PaymentProcessor{
		MinResponseTime: 0,
		Failing:         false,
		Endpoint:        "http://payment-processor-default:8080/payments",
		RedisKey:        "default",
	}, nil
}

func usePaymentProcessor(endpoint, reqBody string, client http.Client) error {
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func ExecutePayment(payment Payment, client http.Client) {
	reqBody := fmt.Sprintf(`{"correlationId":"%s","amount":%d,"requestedAt":"%s"}`, payment.CorrelationID, payment.Amount, payment.RequestedAt.Format(time.RFC3339))

	pp, err := choosePaymentProcessor()
	if err != nil {
		return
	}

	err = usePaymentProcessor(pp.Endpoint, reqBody, client)
	if err != nil {
		return
	}

}
