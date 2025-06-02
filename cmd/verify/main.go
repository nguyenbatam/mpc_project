// Package main provides the entry point for the verification CLI.
package main

import (
	"flag"
	"fmt"
	"log"

	mpcverification "github.com/nguyenbatam/mpc_project/internal/verification"
)

func main() {
	message := flag.String("m", "Hello", "Sample message to sign")
	signatureHex := flag.String("s", "c70b63e7fb0286c2432d02704885372ee08cb04f88b73024987c8264d2fd52915b61f2574445ac8874326cc3ba66760c7e683c5243cf773c056372db67901fe8", "Signature Hex")

	flag.Parse()

	fmt.Printf("Verifying signature for message: '%s'\n", *message)
	fmt.Printf("Signature (Hex): %s\n", *signatureHex)

	valid, err := mpcverification.VerifySignature(*message, *signatureHex)
	if err != nil {
		log.Fatalf("Verification failed: %v", err)
	}

	if valid {
		fmt.Println("Signature is valid.")
	} else {
		fmt.Println("Signature is invalid.")
	}
}
