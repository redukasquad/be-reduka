package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// GetSMTPConfig returns SMTP configuration from environment variables
func GetSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM_EMAIL"),
	}
}

func SendEmail(to string, subject string, body string) error {
	config := GetSMTPConfig()

	log.Printf("[EMAIL] Attempting to send email to: %s", to)

	// Validate configuration
	if config.Host == "" {
		log.Printf("[EMAIL] ERROR: SMTP_HOST is not set")
		return fmt.Errorf("SMTP_HOST is not set")
	}
	if config.Port == "" {
		config.Port = "587" // Default SMTP port
	}
	if config.Username == "" {
		log.Printf("[EMAIL] ERROR: SMTP_USERNAME is not set")
		return fmt.Errorf("SMTP_USERNAME is not set")
	}
	if config.Password == "" {
		log.Printf("[EMAIL] ERROR: SMTP_PASSWORD is not set")
		return fmt.Errorf("SMTP_PASSWORD is not set")
	}
	if config.From == "" {
		config.From = config.Username // Default to username if from is not set
	}

	// Setup authentication
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// Build the email message with proper headers
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n%s\r\n%s",
		config.From,
		to,
		subject,
		mime,
		body,
	))

	// Send email
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	err := smtp.SendMail(addr, auth, config.From, []string{to}, msg)
	if err != nil {
		log.Printf("[EMAIL] ERROR sending email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("[EMAIL] SUCCESS: Email sent to %s", to)
	return nil
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(code)
}
