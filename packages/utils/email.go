package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
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

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}

	return SMTPConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     fromEmail,
	}
}

func SendEmail(to string, subject string, body string) error {
	config := GetSMTPConfig()

	log.Printf("[EMAIL] Attempting to send email to: %s", to)
	log.Printf("[EMAIL] Using SMTP: %s:%s from %s", config.Host, config.Port, config.From)

	// Validate configuration
	if config.Host == "" {
		log.Printf("[EMAIL] ERROR: SMTP_HOST is not set")
		return fmt.Errorf("SMTP_HOST is not set")
	}
	if config.Username == "" {
		log.Printf("[EMAIL] ERROR: SMTP_USERNAME is not set")
		return fmt.Errorf("SMTP_USERNAME is not set")
	}
	if config.Password == "" {
		log.Printf("[EMAIL] ERROR: SMTP_PASSWORD is not set")
		return fmt.Errorf("SMTP_PASSWORD is not set")
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

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	// Try based on port
	if config.Port == "465" {
		// Port 465 uses implicit TLS (SSL)
		return sendWithSSL(config, to, msg, addr)
	}

	// Port 587 or 25 uses STARTTLS
	return sendWithSTARTTLS(config, to, msg, addr)
}

func sendWithSSL(config SMTPConfig, to string, msg []byte, addr string) error {
	log.Printf("[EMAIL] Using SSL/TLS connection")

	// Connect with TLS directly
	tlsConfig := &tls.Config{
		ServerName: config.Host,
	}

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, tlsConfig)
	if err != nil {
		log.Printf("[EMAIL] ERROR connecting with TLS: %v", err)
		return fmt.Errorf("TLS dial failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		log.Printf("[EMAIL] ERROR creating SMTP client: %v", err)
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer client.Close()

	// Authenticate
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	if err = client.Auth(auth); err != nil {
		log.Printf("[EMAIL] ERROR authenticating: %v", err)
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Send email
	if err = client.Mail(config.From); err != nil {
		log.Printf("[EMAIL] ERROR setting sender: %v", err)
		return fmt.Errorf("set sender failed: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		log.Printf("[EMAIL] ERROR setting recipient: %v", err)
		return fmt.Errorf("set recipient failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		log.Printf("[EMAIL] ERROR getting data writer: %v", err)
		return fmt.Errorf("data command failed: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		log.Printf("[EMAIL] ERROR writing message: %v", err)
		return fmt.Errorf("write message failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		log.Printf("[EMAIL] ERROR closing writer: %v", err)
		return fmt.Errorf("close writer failed: %w", err)
	}

	client.Quit()
	log.Printf("[EMAIL] SUCCESS: Email sent to %s via SSL", to)
	return nil
}

func sendWithSTARTTLS(config SMTPConfig, to string, msg []byte, addr string) error {
	log.Printf("[EMAIL] Using STARTTLS connection")

	// Connect with timeout
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		log.Printf("[EMAIL] ERROR dialing: %v", err)
		return fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		log.Printf("[EMAIL] ERROR creating client: %v", err)
		return fmt.Errorf("client creation failed: %w", err)
	}
	defer client.Close()

	// STARTTLS
	tlsConfig := &tls.Config{
		ServerName: config.Host,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		log.Printf("[EMAIL] ERROR STARTTLS: %v", err)
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	// Authenticate
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	if err = client.Auth(auth); err != nil {
		log.Printf("[EMAIL] ERROR authenticating: %v", err)
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Send email
	if err = client.Mail(config.From); err != nil {
		log.Printf("[EMAIL] ERROR setting sender: %v", err)
		return fmt.Errorf("set sender failed: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		log.Printf("[EMAIL] ERROR setting recipient: %v", err)
		return fmt.Errorf("set recipient failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		log.Printf("[EMAIL] ERROR getting data writer: %v", err)
		return fmt.Errorf("data command failed: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		log.Printf("[EMAIL] ERROR writing message: %v", err)
		return fmt.Errorf("write message failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		log.Printf("[EMAIL] ERROR closing writer: %v", err)
		return fmt.Errorf("close writer failed: %w", err)
	}

	client.Quit()
	log.Printf("[EMAIL] SUCCESS: Email sent to %s via STARTTLS", to)
	return nil
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(code)
}
