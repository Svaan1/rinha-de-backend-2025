package api

import (
	"io"
	"math"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
	"github.com/svaan1/rinha-de-backend-2025/internal/payments"
)

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	go func() {
		// Unmarshal de JSON
		var data PaymentRequest
		if err := sonic.Unmarshal(body, &data); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		// Apply the task to the queue
		task := func() {
			payments.ExecutePayment(payments.Payment{
				CorrelationID: data.CorrelationID,
				Amount:        data.Amount,
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
		return
	}

	// Marshal
	jsonBytes, err := sonic.Marshal(result)
	if err != nil {
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func PurgePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	globals.RedisClient.Purge()
}
