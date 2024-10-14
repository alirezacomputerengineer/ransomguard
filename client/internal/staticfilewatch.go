package internal

import (
	"io/ioutil"
	"log"
	"strings"
)

// MonitorStaticFiles monitors newly created static files for suspicious content
func MonitorStaticFiles(config *Config) error {
	log.Println("Starting static file monitoring...")

	// Simulate scanning files in a directory
	filesToCheck := []string{"example.exe", "malware.db", "document.pdf"}

	for _, file := range filesToCheck {
		content, err := ioutil.ReadFile(file) // Read the file content
		if err != nil {
			log.Printf("Failed to read file %s: %v", file, err)
			continue
		}

		if containsSuspiciousContent(string(content), config.StaticFileKeywords) {
			log.Printf("Suspicious content found in file %s", file)
			SendAlert(config.CustomerEmail, config.CompanyEmail, file)
		}
	}
	return nil
}

// containsSuspiciousContent checks if the file content contains any suspicious keywords
func containsSuspiciousContent(fileContent string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(fileContent, keyword) {
			return true
		}
	}
	return false
}
