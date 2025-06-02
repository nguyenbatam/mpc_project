// Package main provides the entry point for the signing CLI.
package main

import (
	"flag"
	"fmt"
	"github.com/ipfs/go-log"
	log2 "github.com/ipfs/go-log/v2"
	mpcsigning "github.com/nguyenbatam/mpc_project/internal/signing"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.SetAllLoggers(log2.LevelError)
	// Define command line parameters
	partyIndexesStr := flag.String("p", "0,1,2,3", "Comma-separated list of party indexes (e.g., 0,1)")
	message := flag.String("m", "Hello", "Sample message to sign")
	flag.Parse()

	// Convert the index string to an integer array
	partyIndexesStrArr := strings.Split(*partyIndexesStr, ",")
	partyIndexes := make([]int, 0, len(partyIndexesStrArr))

	for _, indexStr := range partyIndexesStrArr {
		index, err := strconv.Atoi(strings.TrimSpace(indexStr))
		if err != nil {
			fmt.Printf("Error: Invalid party index '%s': %v\n\n", indexStr, err)
			os.Exit(1)
		}
		partyIndexes = append(partyIndexes, index)
	}

	fmt.Printf("Starting signing process for message: '%s'\n", *message)
	fmt.Printf("Using participating party indexes: %v\n", partyIndexes)

	signature, err := mpcsigning.SignMessage(*message, partyIndexes)
	if err != nil {
		fmt.Printf("Signing failed: %v", err)
	}

	fmt.Printf("Signature generated successfully for message: '%s' \n", *message)
	fmt.Println(signature)
}
