package internal

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
)

// Process struct represents a system process
type Process struct {
	Name     string
	ID       int
	APICalls []string // Populated with actual API/system call data from strace
}

// MonitorKernelAPICalls monitors API calls for processes not in SecureProcesses using strace.
func MonitorKernelAPICalls(config Config, alertChan chan Alert) {
	var wg sync.WaitGroup

	for {
		// Get running processes
		processes := getRunningProcesses()

		for _, proc := range processes {
			// Skip secure processes
			if isSecureProcess(proc, config.SecureProcesses) {
				continue
			}

			wg.Add(1)
			// Run strace in a separate goroutine for each unsecure process
			go func(p Process) {
				defer wg.Done()
				straceAPICalls := captureSyscalls(p.ID)

				for _, apiCall := range config.UnsecureAPIChainCalls {
					if hasUnsecureAPICall(straceAPICalls, apiCall) {
						// Send alert if unsecure API call found
						alert := Alert{
							Description: apiCall,
							ProcessName: p.Name,
							ProcessID:   p.ID,
						}
						alertChan <- alert
					}
				}
			}(proc)
		}

		wg.Wait()

		// Monitor at intervals (e.g., every second)
		time.Sleep(1 * time.Second)
	}
}

// getRunningProcesses fetches the list of running processes using gopsutil
func getRunningProcesses() []Process {
	var processes []Process

	// Get a list of all running processes
	procs, err := process.Processes()
	if err != nil {
		fmt.Println("Error fetching processes:", err)
		return processes
	}

	// Iterate over each process
	for _, p := range procs {
		// Get process name and ID
		name, err := p.Name()
		if err != nil {
			continue
		}
		pid := int(p.Pid)

		// Initialize empty API calls slice
		apiCalls := []string{}

		// Add the process to the list
		processes = append(processes, Process{
			Name:     name,
			ID:       pid,
			APICalls: apiCalls, // This will be populated by strace
		})

		// Debug print
		fmt.Printf("Process: %s, PID: %d\n", name, pid)
	}

	return processes
}

// captureSyscalls uses strace to capture system calls for a given process by PID
func captureSyscalls(pid int) []string {
	var syscalls []string

	// Command to run strace on the given PID
	cmd := exec.Command("strace", "-p", fmt.Sprintf("%d", pid), "-e", "trace=all")

	// Get stdout pipe to read strace output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error getting stdout:", err)
		return nil
	}

	// Start strace command
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting strace:", err)
		return nil
	}

	// Read strace output line by line
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if syscall := extractSyscall(line); syscall != "" {
			syscalls = append(syscalls, syscall)
		}
	}

	// Wait for strace command to finish
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for strace:", err)
	}

	return syscalls
}

// extractSyscall parses strace output and extracts system calls
func extractSyscall(straceOutput string) string {
	// Simple parsing of strace output to extract the syscall name
	// Example: openat(AT_FDCWD, "/etc/ld.so.cache", O_RDONLY|O_CLOEXEC) = 3
	parts := strings.Split(straceOutput, "(")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0]) // The syscall name is before the '('
	}
	return ""
}

// isSecureProcess checks if a process is in the SecureProcesses list
func isSecureProcess(proc Process, secureProcesses []string) bool {
	for _, sp := range secureProcesses {
		if strings.Contains(proc.Name, sp) {
			return true
		}
	}
	return false
}

// hasUnsecureAPICall checks if a process contains an unsecure API call
func hasUnsecureAPICall(apiCalls []string, apiCall string) bool {
	for _, call := range apiCalls {
		if strings.Contains(call, apiCall) {
			return true
		}
	}
	return false
}

// Example of running the monitoring function
/*func main() {
	config := Config{
		SecureProcesses:       []string{"systemd", "bash"},
		UnsecureAPIChainCalls: []string{"open", "exec", "socket"},
	}

	alertChan := make(chan Alert)

	// Start monitoring in a separate goroutine
	go MonitorKernelAPICalls(config, alertChan)

	// Handle alerts in a separate goroutine
	go func() {
		for alert := range alertChan {
			fmt.Printf("ALERT: Process %s (PID: %d) used unsecure API call: %s\n", alert.ProcessName, alert.ProcessID, alert.Description)
		}
	}()

	// Keep the main program running
	select {}
}*/
