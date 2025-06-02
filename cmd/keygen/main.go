// Package main provides the entry point for the key generation CLI.
package main

import (
	"fmt"
	"github.com/ipfs/go-log"
	log2 "github.com/ipfs/go-log/v2"
	mpckeygen "github.com/nguyenbatam/mpc_project/internal/keygen"
	"os"
)

func main() {
	fmt.Println("Starting TSS 2-of-3 distributed key generation process...")
	log.SetAllLoggers(log2.LevelInfo)
	err := mpckeygen.GenerateDistributedKey()
	if err != nil {
		fmt.Printf("Key generation failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Key generation process completed. Data saved in the data/ directory.")
}
