# Threshold Signature Scheme (TSS) using tss-lib

This project implements a Threshold Signature Scheme (TSS) for ECDSA using the `github.com/bnb-chain/tss-lib/v2` library.

It allows multiple parties to collaboratively generate a shared private key and sign messages without any single party holding the full key.

## Features:

-   **Key Generation:** Generate shared private keys for multiple parties.
-   **Signing:** Collaboratively sign a message using the shared key shares.
-   **Verification:** Verify the validity of a generated signature.
-   **Error Handling:** Basic error handling for key generation, signing, and verification.

## Requirements:

-   Go 1.18 or later.
-   `github.com/bnb-chain/tss-lib/v2` (Managed by Go Modules).

## Setup:

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/nguyenbatam/mpc_project
    cd mpc_project
    ```

2.  **Download dependencies:**

    ```bash
    go mod tidy
    ```

## Usage:

To run the TSS application, navigate to the project's root directory and execute the main program:

```bash
go run main.go
```

This will start an interactive menu in your terminal. Follow the prompts to perform key generation, signing, or verification.

### Interactive Menu Options:

1.  **Key Generation:**
    - Select option `1` from the menu.
    - The program will display the total number of parties and the threshold used for key generation.
    - Upon successful completion, each party's share data will be saved in the `./data` directory (`party_*.json`). The generated public key will also be printed to the console.
    - **Important:** Running key generation again will overwrite the existing data in the `./data` directory.

2.  **Signing:**
    - Select option `2` from the menu.
    - The program will display the available party indexes and the minimum required parties for signing (Threshold + 1).
    - You will be prompted to enter the participating party indexes (comma-separated, e.g., `0,1,2`).
    - You will then be prompted to enter the message to sign.
    - If successful, the R and S components of the signature will be printed to the console as a hex string.

3.  **Verification:**
    - Select option `3` from the menu.
    - You will be prompted to enter the signature (as a hex string).
    - You will then be prompted to enter the original message.
    - The command will print whether the signature is valid or invalid.

4.  **Exit:**
    - Select option `4` to exit the application.

## Error Handling:

Basic error handling is included to catch issues during the processes. Errors will be printed to the console.

## Additional Challenges (Optional):

-   Modify the `Threshold` and `PartyCount` constants in `internal/common/params.go` to change the MPC configuration (e.g., 3-of-5).
-   Implement a network layer to simulate communication between parties running on different machines.
-   Add persistence for the signed message and signature.

## Configuration:

To change the MPC configuration, such as the number of parties (`PartyCount`) and the signing threshold (`Threshold`), you need to modify the `config.json` file located in the `data/` directory at the root of the project.

If the `data/config.json` file does not exist, the program will create a default one with `PartyCount: 5` and `Threshold: 2` when it runs for the first time. You can then modify this file.

The `config.json` file contains two main fields:

-   `PartyCount`: The total number of parties involved in the key generation process.
-   `Threshold`: The minimum number of parties required to sign a message.

Update the values for `PartyCount` and `Threshold` in this file to your desired configuration.

**Note:** Ensure that `Threshold` is always less than `PartyCount` and that the signing process requires `Threshold + 1` parties.

After saving the changes to `data/config.json`, you will need to re-run the key generation process (`Option 1` in the interactive menu) to generate new key shares based on the updated configuration.
