package internal

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"ransomguard/utils"
)

const configFile = "/usr/local/ransomguard/config.txt" // Path to config file
const encryptionKey = "hardcoded-encryption-key"       // Simple hardcoded key

// Config holds the configuration details for ransomguard
type Config struct {
	HoneypotFiles        []HoneypotFile `json:"honeypot_files"`
	CustomerEmail        string         `json:"customer_email"`
	CompanyEmail         string         `json:"company_email"`
	StaticFileKeywords   []string       `json:"static_file_keywords"`
	SecureProcesses      []string       `json:"secure_processes"`
	UnsecureAPIChainCalls []string      `json:"unsecure_api_chain_calls"`
	RollbackWant         bool           `json:"rollback_want"`         // New item 1
	RollbackMethod       string         `json:"rollback_method"`       // New item 2
	RollbackParams       []string       `json:"rollback_params"`       // New item 3
}

// LoadConfig reads the config file, decrypts it, and unmarshals the JSON data
func LoadConfig() (*Config, error) {
	encryptedData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println("Config file not found.")
		return nil, errors.New("config file not found")
	}

	// Decrypt data
	decryptedData, err := utils.DecryptData(encryptedData, encryptionKey)
	if err != nil {
		log.Println("Failed to decrypt config file.")
		return nil, errors.New("failed to decrypt config file")
	}

	// Unmarshal JSON
	var config Config
	err = json.Unmarshal(decryptedData, &config)
	if err != nil {
		log.Println("Failed to unmarshal config file.")
		return nil, errors.New("failed to parse config file")
	}

	log.Println("Configuration loaded successfully.")
	return &config, nil
}

// SaveConfig marshals and encrypts the config struct, and writes it to a file
func SaveConfig(config *Config) error {
	// Marshal to JSON
	data, err := json.Marshal(config)
	if err != nil {
		log.Println("Failed to marshal config.")
		return err
	}

	// Encrypt data
	encryptedData, err := utils.EncryptData(data, encryptionKey)
	if err != nil {
		log.Println("Failed to encrypt config.")
		return err
	}

	// Write to file
	err = ioutil.WriteFile(configFile, encryptedData, 0644)
	if err != nil {
		log.Println("Failed to write config to file.")
		return err
	}

	log.Println("Configuration saved successfully.")
	return nil
}
