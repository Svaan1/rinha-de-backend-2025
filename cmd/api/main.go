package main

import (
	"log"
	"net/http"

	"github.com/svaan1/rinha-de-backend-2025/internal/api"
	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
	"github.com/svaan1/rinha-de-backend-2025/internal/payments"
)

func main() {
	go payments.StartHealthCheckTicker()
	globals.QueueDispatcher.Start()

	http.HandleFunc("/payments", api.PaymentHandler)
	http.HandleFunc("/payments-summary", api.PaymentSummaryHandler)
	http.HandleFunc("/purge-payments", api.PurgePaymentsHandler)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Printf("Failed to start server %v", err)
	}
}
