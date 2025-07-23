package api

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
	"github.com/svaan1/rinha-de-backend-2025/internal/payments"
)

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the body
	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to parse body %v", err)
		return
	}

	// Dispatch the task
	go func() {
		task := func() {
			payments.ExecutePayment(payments.Payment{
				CorrelationID: req.CorrelationID,
				Amount:        req.Amount,
				RequestedAt:   time.Now().UTC(),
			})
		}

		globals.QueueDispatcher.Dispatch(task)
	}()

	// Return 201
	w.WriteHeader(http.StatusCreated)
}

func PaymentSummaryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the query params
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	// Parse the timestamps
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

	// Query from redis
	result, err := globals.RedisClient.GetPaymentSummary(fromTimestamp, toTimestamp)
	if err != nil {
		log.Printf("Failed to get payment summary: %v", err)
	}

	// Return the encoded struct
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func PurgePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	globals.RedisClient.Purge()
}
