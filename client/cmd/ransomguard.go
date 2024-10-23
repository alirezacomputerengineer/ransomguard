package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"ransomguard/internal"
	"syscall"
	"time"
)

func main() {
	// Load configuration from the encrypted config file
	config, err := internal.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize honeypot files if they don't already exist
	err = internal.CreateHoneypotFiles(config)
	if err != nil {
		log.Fatalf("Error initializing honeypot files: %v", err)
	}

	// Channel for OS signals to allow graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal an issue, such as ransomware detection, across all goroutines
	alertChan := make(chan internal.Alert)

	// Start monitoring honeypot files for any modifications
	go func() {
		err := internal.MonitorHoneypotFiles(config, alertChan)
		if err != nil {
			log.Fatalf("Error in honeypot file watching: %v", err)
		}
	}()

	// Start kernel API watch for unsecure API call chains
	go func() {
		err := internal.MonitorKernelAPI(config, alertChan)
		if err != nil {
			log.Fatalf("Error in kernel API monitoring: %v", err)
		}
	}()

	// Start static file watch for suspicious new files
	go func() {
		err := internal.MonitorStaticFiles(config, alertChan)
		if err != nil {
			log.Fatalf("Error in static file watching: %v", err)
		}
	}()

	// Monitor for alerts from all goroutines
	go func() {
		for alert := range alertChan {
			handleAlert(alert, config)
		}
	}()

	// Wait for termination signals to gracefully exit
	<-sigChan
	log.Println("Ransomguard is shutting down.")
}

// handleAlert handles alerts when a potential ransomware attack is detected
func handleAlert(alert internal.Alert, config *internal.Config) {
	log.Printf("ALERT: %s detected on process: %s", alert.Description, alert.ProcessName)

	// Step 2.1: Stop the process that triggered the alert
	err := internal.TerminateProcess(alert.ProcessID)
	if err != nil {
		log.Printf("Failed to terminate process %d: %v", alert.ProcessID, err)
	} else {
		log.Printf("Process %d terminated successfully.", alert.ProcessID)
	}

	// Step 2.2: Send the process file to quarantine
	err = internal.QuarantineProcess(alert.ProcessName)
	if err != nil {
		log.Printf("Failed to quarantine process: %v", err)
	} else {
		log.Printf("Process %s moved to quarantine.", alert.ProcessName)
	}

	// Step 2.3: Send an alert email to both the customer and ransomguard company
	err = internal.SendAlertEmail(config.CustomerEmail, config.CompanyEmail, alert)
	if err != nil {
		log.Printf("Failed to send alert email: %v", err)
	} else {
		log.Println("Alert email sent successfully.")
	}

	// Step 2.4: If rollback is enabled, perform the rollback
	if config.RollbackWant {
		err = internal.PerformRollback(config)
		if err != nil {
			log.Printf("Rollback failed: %v", err)
		} else {
			log.Println("Rollback completed successfully.")
		}
	}
}
