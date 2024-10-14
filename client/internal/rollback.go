package internal

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// PerformRollback checks if the rollback is enabled and executes the appropriate rollback method
func PerformRollback(config *Config) error {
	// Check if rollback is enabled by customer
	if !config.RollbackWant {
		log.Println("Rollback functionality is disabled by the customer.")
		return nil
	}

	// Switch based on the rollback method provided in the config
	switch config.RollbackMethod {
	case "btrfs":
		return rollbackBtrfs(config.RollbackParams)
	// Add more rollback methods as needed
	default:
		log.Printf("Unsupported rollback method: %s", config.RollbackMethod)
		return fmt.Errorf("unsupported rollback method")
	}
}

// rollbackBtrfs handles rollback using the BTRFS file system
func rollbackBtrfs(params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("no rollback parameters provided for BTRFS")
	}

	subvolume := params[0]           // e.g., "subvolume" or path to snapshot
	restorePath := params[1]         // e.g., "/restored-directory"
	snapshotPath := params[2]        // e.g., "/backup/snapshot"

	log.Printf("Starting BTRFS rollback: subvolume: %s, snapshot: %s, restorePath: %s", subvolume, snapshotPath, restorePath)

	// Example command (customize based on your rollback logic)
	cmd := exec.Command("btrfs", "subvolume", "snapshot", snapshotPath, restorePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to rollback using BTRFS: %v\nCommand output: %s", err, string(output))
		return fmt.Errorf("rollback failed: %v", err)
	}

	log.Printf("BTRFS rollback successful. Restored snapshot to: %s", restorePath)
	return nil
}
