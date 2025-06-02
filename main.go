package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nguyenbatam/mpc_project/internal/common"
)

func main() {
	fmt.Println("Threshold Signature Scheme based Multi-Party Computation (MPC) System")
	fmt.Println("=================================================================")

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. Key Generation")
		fmt.Println("2. Signing")
		fmt.Println("3. Verification")
		fmt.Println("4. Exit")

		var choice int
		fmt.Print("Enter your choice (1-4): ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			runKeygen()
		case 2:
			runSign()
		case 3:
			runVerify()
		case 4:
			fmt.Println("Exiting.")
			return // Exit the program
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func runKeygen() {
	// Display total number of parties and threshold for key generation
	fmt.Printf("Total parties for key generation: %d\n", common.PartyCount)
	fmt.Printf("Threshold for key generation: %d\n", common.Threshold)

	cmd := exec.Command("go", "run", filepath.Join("cmd", "keygen", "main.go"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running key generation process: %v\n", err)
	}
}

func runSign() {
	var partyIndexes, message string

	// Display available parties and the required threshold
	fmt.Printf("Available parties for signing: 0 to %d\n", common.PartyCount-1)   // PartyCount = 5 (0-indexed)
	fmt.Printf("Minimum required parties for signing: %d \n", common.Threshold+1) // Threshold = 2, minimum parties = Threshold + 1

	fmt.Print("Enter participating party indexes (comma-separated, e.g., 0,1,2): ")
	fmt.Scanln(&partyIndexes)

	fmt.Print("Enter message to sign: ")
	fmt.Scanln(&message)

	cmd := exec.Command("go", "run", filepath.Join("cmd", "sign", "main.go"), "-p", partyIndexes, "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running signing process: %v\n", err)
	}
}

func runVerify() {
	var signature, message string

	fmt.Print("Enter signature (hex string): ")
	fmt.Scanln(&signature)

	fmt.Print("Enter original message: ")
	fmt.Scanln(&message)

	cmd := exec.Command("go", "run", filepath.Join("cmd", "verify", "main.go"), "-s", signature, "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running verification process: %v\n", err)
	}
}
