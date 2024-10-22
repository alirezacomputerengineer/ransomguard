package internal

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// CreateHoneypotFiles generates and monitors honeypot files based on the configuration
func CreateHoneypotFiles(config *Config) error {
	for _, honeypot := range config.HoneypotFiles {
		filePath := filepath.Join(honeypot.Route, honeypot.Name+honeypot.Extension)

		// Check if file already exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Create a new honeypot file
			err := createFileWithVolume(filePath, honeypot.Volume)
			if err != nil {
				log.Printf("Failed to create honeypot file %s: %v", filePath, err)
				return err
			}
			log.Printf("Honeypot file created: %s", filePath)
		} else {
			log.Printf("Honeypot file already exists: %s", filePath)
		}
	}
	return nil
}

// createFileWithVolume creates a file and fills it with dummy data to reach the specified volume in KB
func createFileWithVolume(filePath string, volumeKB int) error {
	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Fill the file with dummy data to match the volume in kilobytes
	data := make([]byte, 1024) // 1 KB of dummy data
	for i := 0; i < volumeKB; i++ {
		_, err := file.Write(data)
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	}

	return nil
}

// MonitorHoneypotFile monitors honeypot files for modifications and sends alerts
func MonitorHoneypotFile(config *Config, alertChan chan Alert) {
	for _, honeypot := range config.HoneypotFiles {
		filePath := filepath.Join(honeypot.Route, honeypot.Name+honeypot.Extension)
		go monitorFile(filePath, alertChan)
	}
}

func monitorFile(h HoneypotFile, alertChan chan Alert) {
	// Create a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating file watcher: %v", err)
	}
	defer watcher.Close()

	// Add the honeypot file to the watcher
	err = watcher.Add(h.Route)
	if err != nil {
		log.Fatalf("Error adding file to watcher: %v", err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Check if the event is a write (modification)
			if event.Op&fsnotify.Write == fsnotify.Write {
				processName, processID := getProcessDetails(h.Route)
				alert := Alert{
					Description: fmt.Sprintf("Modified: %s - %s", event.Name, h.Name),
					ProcessName: processName,
					ProcessID:   processID,
				}
				alertChan <- alert
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error watching file: %v", err)
		}
	}
}
