package client

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Rizz404/inventory-api/internal/client/smtp"
	"github.com/wneessen/go-mail"
)

// InitSMTP initializes SMTP mail client
func InitSMTP() *smtp.Client {
	enableSMTP := os.Getenv("ENABLE_SMTP") == "true"
	if !enableSMTP {
		log.Printf("SMTP services disabled via ENABLE_SMTP environment variable")
		return nil
	}

	host := os.Getenv("SMTP_HOST")
	if host == "" {
		log.Printf("Warning: SMTP_HOST not set. SMTP services will be disabled.")
		return nil
	}

	portStr := os.Getenv("SMTP_PORT")
	port := 587 // default port
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("SMTP_FROM_EMAIL")
	fromName := os.Getenv("SMTP_FROM_NAME")

	if username == "" || password == "" {
		log.Printf("Warning: SMTP credentials not set. SMTP services will be disabled.")
		return nil
	}

	if fromEmail == "" {
		fromEmail = username
	}
	if fromName == "" {
		fromName = "Inventory API"
	}

	// Create go-mail client
	mailClient, err := mail.NewClient(
		host,
		mail.WithPort(port),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(30*time.Second),
	)
	if err != nil {
		log.Printf("Warning: Failed to create SMTP client: %v. SMTP services will be disabled.", err)
		return nil
	}

	log.Printf("SMTP client initialized successfully for %s", host)

	return &smtp.Client{
		MailClient: mailClient,
		FromEmail:  fromEmail,
		FromName:   fromName,
	}
}
