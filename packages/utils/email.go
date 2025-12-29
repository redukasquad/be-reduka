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
)

// BrevoSender represents the sender in Brevo API
type BrevoSender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// BrevoRecipient represents a recipient in Brevo API
type BrevoRecipient struct {
	Email string `json:"email"`
}

// BrevoRequest represents the request body for Brevo API
type BrevoRequest struct {
	Sender      BrevoSender      `json:"sender"`
	To          []BrevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HTMLContent string           `json:"htmlContent"`
}

// BrevoResponse represents the response from Brevo API
type BrevoResponse struct {
	MessageID string `json:"messageId"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

func SendEmail(to string, subject string, body string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	senderEmail := os.Getenv("BREVO_SENDER_EMAIL")
	senderName := os.Getenv("BREVO_SENDER_NAME")

	// Default sender if not set
	if senderEmail == "" {
		senderEmail = "noreply@reduka.com"
	}
	if senderName == "" {
		senderName = "Reduka"
	}

	log.Printf("[EMAIL] Attempting to send email to: %s via Brevo", to)
	log.Printf("[EMAIL] BREVO_API_KEY present: %v, length: %d", apiKey != "", len(apiKey))

	// Validate configuration
	if apiKey == "" {
		log.Printf("[EMAIL] ERROR: BREVO_API_KEY is not set")
		return fmt.Errorf("BREVO_API_KEY is not set")
	}

	// Convert plain text body to HTML
	htmlBody := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; max-width: 600px; margin: 0 auto;">
			<h2 style="color: #333;">Reduka</h2>
			<p style="font-size: 16px; color: #555;">%s</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="font-size: 12px; color: #999;">This is an automated message, please do not reply.</p>
		</div>
	`, body)

	// Build request body
	reqBody := BrevoRequest{
		Sender: BrevoSender{
			Email: senderEmail,
			Name:  senderName,
		},
		To: []BrevoRecipient{
			{Email: to},
		},
		Subject:     subject,
		HTMLContent: htmlBody,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("[EMAIL] ERROR marshaling request: %v", err)
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("[EMAIL] Request body: %s", string(jsonBody))

	// Create HTTP request to Brevo API
	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("[EMAIL] ERROR creating request: %v", err)
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers - Brevo uses api-key header
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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

	log.Printf("[EMAIL] Response status: %d, body: %s", resp.StatusCode, string(respBody))

	// Check if successful (2xx status code)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("[EMAIL] SUCCESS: Email sent to %s via Brevo", to)
		return nil
	}

	// Parse error response
	var brevoResp BrevoResponse
	if err := json.Unmarshal(respBody, &brevoResp); err != nil {
		log.Printf("[EMAIL] ERROR parsing response: %v", err)
		return fmt.Errorf("Brevo error: status %d", resp.StatusCode)
	}

	log.Printf("[EMAIL] ERROR from Brevo: %s - %s", brevoResp.Code, brevoResp.Message)
	return fmt.Errorf("Brevo error: %s - %s", brevoResp.Code, brevoResp.Message)
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(code)
}
