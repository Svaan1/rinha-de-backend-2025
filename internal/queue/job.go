package queue

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Job struct {
	CorrelationID string
	Amount        int64
}

func (j *Job) Execute() {
	client := &http.Client{Timeout: 10 * time.Second}

	reqBody := fmt.Sprintf(`{"correlationId":"%s","amount":%d}`, j.CorrelationID, j.Amount)

	req, err := http.NewRequest("POST", "http://localhost:8001/payments", strings.NewReader(reqBody))
	if err != nil {
		log.Printf("Error on request creation %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error on request execution %v", err)
		return
	}

	log.Print(resp)
}
