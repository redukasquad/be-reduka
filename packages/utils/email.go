package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ResendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

func SendEmail(to string, subject string, body string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	senderEmail := os.Getenv("RESEND_FROM_EMAIL")

	log.Printf("[EMAIL] Attempting to send email to: %s", to)

	if apiKey == "" {
		log.Printf("[EMAIL] ERROR: RESEND_API_KEY is not set")
		return fmt.Errorf("RESEND_API_KEY is not set")
	}

	if senderEmail == "" {
		senderEmail = "onboarding@resend.dev" // Default Resend sender for testing
	}

	emailReq := ResendEmailRequest{
		From:    senderEmail,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	jsonData, err := json.Marshal(emailReq)
	if err != nil {
		log.Printf("[EMAIL] ERROR marshaling request: %v", err)
		return err
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[EMAIL] ERROR creating request: %v", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[EMAIL] ERROR sending request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errorBody bytes.Buffer
		errorBody.ReadFrom(resp.Body)
		log.Printf("[EMAIL] ERROR response from Resend (status %d): %s", resp.StatusCode, errorBody.String())
		return fmt.Errorf("resend API error: %s", errorBody.String())
	}

	log.Printf("[EMAIL] SUCCESS: Email sent to %s", to)
	return nil
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(code)
}
