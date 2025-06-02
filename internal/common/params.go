// Package common provides common constants and utility functions for the MPC project.
package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"

	"github.com/bnb-chain/tss-lib/v2/tss"
)

// Define a struct to hold configuration parameters
type Config struct {
	Threshold  int `json:"Threshold"`
	PartyCount int `json:"PartyCount"`
}

// Global variables to store configuration
var (
	Threshold  int
	PartyCount int
)

func init() {
	config, err := LoadConfig()
	if err != nil {
		// Handle error loading config, perhaps panic or log and exit
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}
	Threshold = config.Threshold
	PartyCount = config.PartyCount
}

// LoadConfig reads configuration from data/config.json. If the file does not exist,
// it creates a default config file and returns the default values.
func LoadConfig() (*Config, error) {
	configPath := filepath.Join("data", "config.json")

	// Try to open the config file
	configFile, err := os.Open(configPath)
	if err != nil {
		// If the file does not exist, create a default one
		if os.IsNotExist(err) {
			fmt.Printf("Config file not found. Creating default config file at %s\n", configPath)
			defaultConfig := Config{Threshold: 2, PartyCount: 5}

			// Ensure data directory exists
			if err := os.MkdirAll("data", 0755); err != nil {
				return nil, fmt.Errorf("failed to create data directory for default config: %v", err)
			}

			bytes, err := json.MarshalIndent(defaultConfig, "", "  ")
			if err != nil {
				return nil, fmt.Errorf("failed to marshal default config to JSON: %v", err)
			}

			if err := os.WriteFile(configPath, bytes, 0644); err != nil {
				return nil, fmt.Errorf("failed to write default config file: %v", err)
			}

			return &defaultConfig, nil
		} else {
			// For other errors, return the error
			return nil, fmt.Errorf("failed to open config file: %v", err)
		}
	}
	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config JSON: %v", err)
	}

	return &config, nil
}

// GenerateTestPartyIDs creates a slice of PartyID for testing purposes.
// The PartyID.Index is set starting from 1 to avoid index 0 issues with tss-lib.
func GenerateTestPartyIDs(count int) tss.UnSortedPartyIDs {
	pIDs := make(tss.UnSortedPartyIDs, 0, count)
	for i := 0; i < count; i++ {
		// Change index from i -> i+1 to avoid index = 0
		pIDs = append(pIDs, tss.NewPartyID(string(rune(i)), "", big.NewInt(int64(i+1))))
	}
	return pIDs
}

// Error messages
const (
	ErrInvalidThreshold = "invalid threshold: expected at least %d parties to sign"
	ErrPartyNotFound    = "party data for index %d not found"
	ErrSigningFailed    = "signing process failed: %s"
	ErrVerifyFailed     = "signature verification failed: %v"
)
