package utils

import (
	"crypto/tls"
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
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	if fromEmail == "" {
		fromEmail = os.Getenv("SMTP_USERNAME") // Fall back to username
	}

	return SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     fromEmail,
	}
}

func SendEmail(to string, subject string, body string) error {
	config := GetSMTPConfig()

	log.Printf("[EMAIL] Attempting to send email to: %s", to)
	log.Printf("[EMAIL] Using SMTP Host: %s:%s", config.Host, config.Port)
	log.Printf("[EMAIL] From: %s", config.From)

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
		config.From = config.Username
	}

	// Build the email message
	msg := []byte(
		"From: " + config.From + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n\r\n" +
			body,
	)

	// Connect to the SMTP server with TLS
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	// TLS configuration
	tlsConfig := &tls.Config{
		ServerName: config.Host,
	}

	// Connect to the server
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		// If direct TLS fails, try STARTTLS approach
		log.Printf("[EMAIL] Direct TLS failed, trying STARTTLS: %v", err)
		return sendWithSTARTTLS(config, to, msg)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		log.Printf("[EMAIL] ERROR creating SMTP client: %v", err)
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	if err = client.Auth(auth); err != nil {
		log.Printf("[EMAIL] ERROR authenticating: %v", err)
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender and recipient
	if err = client.Mail(config.From); err != nil {
		log.Printf("[EMAIL] ERROR setting sender: %v", err)
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err = client.Rcpt(to); err != nil {
		log.Printf("[EMAIL] ERROR setting recipient: %v", err)
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		log.Printf("[EMAIL] ERROR getting data writer: %v", err)
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		log.Printf("[EMAIL] ERROR writing message: %v", err)
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		log.Printf("[EMAIL] ERROR closing writer: %v", err)
		return fmt.Errorf("failed to close writer: %w", err)
	}

	client.Quit()

	log.Printf("[EMAIL] SUCCESS: Email sent to %s", to)
	return nil
}

func sendWithSTARTTLS(config SMTPConfig, to string, msg []byte) error {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	// Connect to server without TLS first
	conn, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("[EMAIL] ERROR dialing SMTP: %v", err)
		return fmt.Errorf("failed to dial SMTP: %w", err)
	}
	defer conn.Close()

	// Say hello
	if err = conn.Hello("localhost"); err != nil {
		log.Printf("[EMAIL] ERROR HELO: %v", err)
		return fmt.Errorf("failed HELO: %w", err)
	}

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName: config.Host,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		log.Printf("[EMAIL] ERROR STARTTLS: %v", err)
		return fmt.Errorf("failed STARTTLS: %w", err)
	}

	// Authenticate
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	if err = conn.Auth(auth); err != nil {
		log.Printf("[EMAIL] ERROR authenticating: %v", err)
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err = conn.Mail(config.From); err != nil {
		log.Printf("[EMAIL] ERROR setting sender: %v", err)
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err = conn.Rcpt(to); err != nil {
		log.Printf("[EMAIL] ERROR setting recipient: %v", err)
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message
	w, err := conn.Data()
	if err != nil {
		log.Printf("[EMAIL] ERROR getting data writer: %v", err)
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		log.Printf("[EMAIL] ERROR writing message: %v", err)
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = w.Close()
	if err != nil {
		log.Printf("[EMAIL] ERROR closing writer: %v", err)
		return fmt.Errorf("failed to close writer: %w", err)
	}

	conn.Quit()

	log.Printf("[EMAIL] SUCCESS: Email sent to %s (via STARTTLS)", to)
	return nil
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(code)
}
