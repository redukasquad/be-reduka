package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// PromailerRequest represents the request body for Promailer API
type PromailerRequest struct {
	MessageID string `json:"messageId"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	HTML      string `json:"html"`
}

// PromailerResponse represents the response from Promailer API
type PromailerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		MessageID string `json:"messageId"`
	} `json:"data"`
}

func SendEmail(to string, subject string, body string) error {
	apiKey := os.Getenv("API_MAIL_KEY")

	log.Printf("[EMAIL] Attempting to send email to: %s via Promailer API", to)
	log.Printf("[EMAIL] API_MAIL_KEY present: %v, length: %d", apiKey != "", len(apiKey))

	// Validate configuration
	if apiKey == "" {
		log.Printf("[EMAIL] ERROR: API_MAIL_KEY is not set")
		return fmt.Errorf("API_MAIL_KEY is not set")
	}

	// Convert plain text body to HTML
	htmlBody := fmt.Sprintf("<div style=\"font-family: Arial, sans-serif; padding: 20px;\">%s</div>", body)

	// Build request body
	reqBody := PromailerRequest{
		MessageID: uuid.New().String(),
		To:        to,
		Subject:   subject,
		HTML:      htmlBody,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[EMAIL] ERROR marshaling request: %v", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("[EMAIL] Request body: %s", string(jsonBody))

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://mailserver.automationlounge.com/api/v1/messages/send", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("[EMAIL] ERROR creating request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request with timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[EMAIL] ERROR sending request: %v", err)
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[EMAIL] ERROR reading response: %v", err)
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var promailerResp PromailerResponse
	if err := json.Unmarshal(respBody, &promailerResp); err != nil {
		log.Printf("[EMAIL] ERROR parsing response: %v, body: %s", err, string(respBody))
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check if successful
	if !promailerResp.Success {
		log.Printf("[EMAIL] ERROR from Promailer: %s", promailerResp.Message)
		return fmt.Errorf("Promailer error: %s", promailerResp.Message)
	}

	log.Printf("[EMAIL] SUCCESS: Email sent to %s via Promailer, messageId: %s", to, promailerResp.Data.MessageID)
	return nil
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(code)
}
