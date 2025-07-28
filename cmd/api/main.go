package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/svaan1/rinha-de-backend-2025/internal/api"
	"github.com/svaan1/rinha-de-backend-2025/internal/globals"
	"github.com/svaan1/rinha-de-backend-2025/internal/payments"
)

func main() {
	go payments.StartHealthCheckTicker()
	globals.QueueDispatcher.Start()

	app := fiber.New()

	app.Post("/payments", api.PaymentHandler)
	app.Get("/payments-summary", api.PaymentSummaryHandler)
	app.Delete("/payments", api.PurgePaymentsHandler)

	if err := app.Listen(":80"); err != nil {
		log.Fatalf("Listen failed: %v", err)
	}
}
