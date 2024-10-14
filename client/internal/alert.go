package internal

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// Email credentials and SMTP server config (update as needed)
const (
	smtpServer     = "smtp.example.com"
	smtpPort       = "587"
	smtpUsername   = "noreply@example.com" // Your email address
	smtpPassword   = "yourpassword"        // Your email password
	senderEmail    = "noreply@example.com" // Sender email
)

// SendAlert sends an email alert to both the system admin and ransomguard company
func SendAlert(customerEmail, companyEmail, modifiedFile string) {
	subject := "RansomGuard Alert: Honeypot File Modified!"
	body := fmt.Sprintf("Alert: The following honeypot file has been modified: %s", modifiedFile)

	// Prepare email
	message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)
	recipients := []string{customerEmail, companyEmail}
	err := sendEmail(recipients, message)
	if err != nil {
		log.Printf("Failed to send alert email: %v", err)
		return
	}
	log.Printf("Alert email sent to: %s, %s", customerEmail, companyEmail)
}

// sendEmail sends an email to the recipients using the provided message body
func sendEmail(recipients []string, message string) error {
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	// Join recipients into a single string
	to := strings.Join(recipients, ", ")

	// Send email
	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpServer, smtpPort), auth, senderEmail, recipients, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
