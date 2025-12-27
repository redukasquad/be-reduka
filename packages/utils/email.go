package utils

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_AUTH_EMAIL")
	smtpPassword := os.Getenv("SMTP_AUTH_PASSWORD")
	smtpSenderName := os.Getenv("SMTP_SENDER_NAME")

	// Debug: Log SMTP configuration (without password)
	log.Printf("[EMAIL] Attempting to send email to: %s", to)
	log.Printf("[EMAIL] SMTP Config - Host: %s, Port: %s, Email: %s, SenderName: %s", smtpHost, smtpPortStr, smtpEmail, smtpSenderName)

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Printf("[EMAIL] ERROR: Invalid SMTP_PORT: %s", smtpPortStr)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpSenderName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpEmail, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("[EMAIL] ERROR sending email to %s: %v", to, err)
		return err
	}

	log.Printf("[EMAIL] SUCCESS: Email sent to %s", to)
	return nil
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999-100000) + 100000
	return strconv.Itoa(code)
}
