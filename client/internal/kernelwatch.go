package internal

import (
	"log"
	"strings"
)

// MonitorKernelAPICalls monitors the kernel API calls of processes in the system
func MonitorKernelAPICalls(config *Config) error {
	log.Println("Starting Kernel API monitoring...")

	// Placeholder for real-time monitoring logic (using ptrace, seccomp, or other)
	// This could be continuously checking system processes

	// Example simulated process monitoring (replace with actual process monitoring logic)
	activeProcesses := []string{"process1", "process2", "malicious_process"}

	for _, process := range activeProcesses {
		if isSecureProcess(process, config.SecureProcesses) {
			log.Printf("Process %s is secure, skipping monitoring.", process)
			continue
		}

		// Simulate detecting API calls (in real implementation, hook actual API calls)
		apiCalls := []string{"rename", "delete_shadow_copy", "encrypt"}

		for _, apiCall := range apiCalls {
			if isUnsecureAPICall(apiCall, config.UnsecureAPIChainCalls) {
				log.Printf("Unsecure API call detected in process %s: %s", process, apiCall)
				SendAlert(config.CustomerEmail, config.CompanyEmail, process)
			}
		}
	}
	return nil
}

// isSecureProcess checks if the given process is part of the secure processes list
func isSecureProcess(process string, secureProcesses []string) bool {
	for _, secure := range secureProcesses {
		if strings.Contains(process, secure) {
			return true
		}
	}
	return false
}

// isUnsecureAPICall checks if the given API call is part of the unsecure API calls list
func isUnsecureAPICall(apiCall string, unsecureAPICalls []string) bool {
	for _, unsecure := range unsecureAPICalls {
		if strings.Contains(apiCall, unsecure) {
			return true
		}
	}
	return false
}
