package internal

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// CreateHoneypotFiles generates and monitors honeypot files based on the configuration
func CreateHoneypotFiles(config *Config) error {
	for _, honeypot := range config.HoneypotFiles {
		filePath := filepath.Join(honeypot.Path, honeypot.Name+honeypot.Extension)

		// Check if file already exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Create a new honeypot file
			err := ioutil.WriteFile(filePath, []byte("This is a honeypot file."), 0644)
			if err != nil {
				log.Printf("Failed to create honeypot file %s: %v", filePath, err)
				return err
			}
			log.Printf("Honeypot file created: %s", filePath)
		} else {
			log.Printf("Honeypot file already exists: %s", filePath)
		}

		// Start monitoring the honeypot file for changes
		err := MonitorHoneypotFile(filePath, config)
		if err != nil {
			log.Printf("Failed to monitor honeypot file %s: %v", filePath, err)
			return err
		}
	}
	return nil
}

// MonitorHoneypotFile watches a single honeypot file for modifications
func MonitorHoneypotFile(filePath string, config *Config) error {
	log.Printf("Monitoring honeypot file: %s", filePath)

	// TODO: Implement file monitoring mechanism to detect file modifications
	// You can use inotify (Linux file system watcher) or polling for changes

	// If any modification is detected:
	// 1. Trigger an alert to the system admin and ransomguard company
	// 2. Send email alerts (refer to email sending logic)

	// Example alert when modification is detected
	modificationDetected := false // Placeholder for actual logic

	if modificationDetected {
		SendAlert(config.CustomerEmail, config.CompanyEmail, filePath)
	}

	return nil
}

// SendAlert triggers an alert by sending emails to the system admin and ransomguard company
func SendAlert(customerEmail, companyEmail, filePath string) {
	// TODO: Implement email sending logic (SMTP or other email service)

	// Example log for alerting
	log.Printf("Alert! Honeypot file modified: %s", filePath)
	log.Printf("Sending email alert to customer: %s and company: %s", customerEmail, companyEmail)
}
