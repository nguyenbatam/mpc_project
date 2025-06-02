// Package verification implements the signature verification process.
package verification

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"fmt"

	"github.com/nguyenbatam/mpc_project/internal/common"
)

// VerifySignature verifies an ECDSA signature for a given message.
func VerifySignature(message, signatureHex string) (bool, error) {
	// Load data of the first party to get the public key
	data, err := common.LoadPartyData(0)
	if err != nil {
		return false, fmt.Errorf("failed to load party data: %v", err)
	}

	// Get the public key
	pubKey := common.GetECDSAPublicKey(data.SavedData.ECDSAPub)

	// Convert signature from hex
	r, s, err := common.HexToSignature(signatureHex)
	if err != nil {
		return false, fmt.Errorf(common.ErrVerifyFailed, err)
	}

	// Hash the message
	msgHash := sha256.Sum256([]byte(message))

	// Verify the signature using standard ecdsa package
	valid := ecdsa.Verify(pubKey, msgHash[:], r, s)
	return valid, nil
}
