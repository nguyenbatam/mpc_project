// Package common provides common constants and utility functions for the MPC project.
package common

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/bnb-chain/tss-lib/v2/crypto"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
)

// PartyData structure to store party-specific data.
type PartyData struct {
	Index     int                        `json:"index"`
	PartyID   *tss.PartyID               `json:"party_id"`
	SavedData *keygen.LocalPartySaveData `json:"saved_data"`
}

// DataDir is the directory path to store party data.
const DataDir = "./data"

// SavePartyData saves the data of a party to a JSON file.
func SavePartyData(index int, partyID *tss.PartyID, data *keygen.LocalPartySaveData) error {
	if err := os.MkdirAll(DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal party data to JSON: %v", err)
	}

	filePath := filepath.Join(DataDir, fmt.Sprintf("party_%d.json", index))
	if err := os.WriteFile(filePath, bytes, 0644); err != nil {
		return fmt.Errorf("failed to write party data to file: %v", err)
	}

	return nil
}

// LoadPartyData loads the data of a party from a JSON file.
func LoadPartyData(index int) (*PartyData, error) {
	filePath := filepath.Join(DataDir, fmt.Sprintf("party_%d.json", index))
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read party data from file: %v", err)
	}

	var partyData PartyData
	partyData.Index = index
	partyData.SavedData = &keygen.LocalPartySaveData{}
	if err := json.Unmarshal(bytes, partyData.SavedData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal party data from JSON: %v", err)
	}
	for _, kbxj := range partyData.SavedData.BigXj {
		kbxj.SetCurve(tss.S256())
	}
	partyData.SavedData.ECDSAPub.SetCurve(tss.S256())

	pMoniker := fmt.Sprintf("%d", partyData.SavedData.ShareID.Int64())
	partyData.PartyID = tss.NewPartyID(pMoniker, pMoniker, partyData.SavedData.ShareID)
	return &partyData, nil
}

// HexToSignature converts a heHexToSignaturex string to an ECDSA signature (R, S).
func HexToSignature(sigHex string) (*big.Int, *big.Int, error) {
	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode hex string: %v", err)
	}

	// Split the bytes to get r and s
	halfway := len(sigBytes) / 2
	r := new(big.Int).SetBytes(sigBytes[:halfway])
	s := new(big.Int).SetBytes(sigBytes[halfway:])

	return r, s, nil
}

// GetECDSAPublicKey converts a tss-lib ECPoint to a standard ECDSA PublicKey.
func GetECDSAPublicKey(pk *crypto.ECPoint) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: tss.EC(), // Assuming P256 curve, adjust if necessary
		X:     pk.X(),
		Y:     pk.Y(),
	}
}
