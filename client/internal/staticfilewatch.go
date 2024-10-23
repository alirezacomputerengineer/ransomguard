package internal

import (
	"io/ioutil"
	"log"
	"strings"
)

// MonitorStaticFiles monitors newly created static files for suspicious content
func MonitorStaticFiles(config Config, alertChan chan Alert) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error creating file watcher:", err)
		return
	}
	defer watcher.Close()

	// Watch directories for new or modified files (you can customize these paths)
	err = watcher.Add("/path/to/monitor") // You may want to monitor a directory like "/usr/bin" or "/opt"
	if err != nil {
		fmt.Println("Error adding path to watcher:", err)
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			// Check for new or modified executable files
			if event.Op&(fsnotify.Create|fsnotify.Write) != 0 {
				go scanFileForKeywords(event.Name, config.StaticFileKeywords, alertChan)
			}
		case err := <-watcher.Errors:
			fmt.Println("Error:", err)
		}
	}
}

func scanFileForKeywords(filePath string, keywords []string, alertChan chan Alert) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, keyword := range keywords {
			if strings.Contains(line, keyword) {
				// Send an alert if keyword is found
				alert := Alert{
					Description: fmt.Sprintf("Keyword '%s' found", keyword),
					ProcessName: "Static File",
					ProcessID:   filePath,
				}
				alertChan <- alert
				// You can stop after the first match or keep searching
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", filePath, err)
	}
}
