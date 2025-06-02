// Package signing implements the signing process for the MPC system.
package signing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	tsscommon "github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"
	mpccommon "github.com/nguyenbatam/mpc_project/internal/common"
)

// SignMessage signs a message using a subset of participating parties.
func SignMessage(message string, partyIndexes []int) (string, error) {
	// Check the number of participating parties
	if len(partyIndexes) < mpccommon.Threshold+1 {
		return "", fmt.Errorf(mpccommon.ErrInvalidThreshold, mpccommon.Threshold)
	}

	// Hash the message
	msgHashBytes := sha256.Sum256([]byte(message))
	msgHashInt := new(big.Int).SetBytes(msgHashBytes[:])

	// Create a list of participating parties
	parties := make(map[string]*signing.LocalParty)
	partyData := make([]*mpccommon.PartyData, 0, len(partyIndexes))

	// Read data for each party
	for _, index := range partyIndexes {
		data, err := mpccommon.LoadPartyData(index)
		if err != nil {
			return "", fmt.Errorf(mpccommon.ErrPartyNotFound, index)
		}
		partyData = append(partyData, data)
	}

	// Initialize PartyID objects from loaded data
	unsortedPIDs := make(tss.UnSortedPartyIDs, len(partyData))
	for i, data := range partyData {
		unsortedPIDs[i] = data.PartyID
	}

	sortedPIDs := tss.SortPartyIDs(unsortedPIDs)

	// Create PeerContext from sortedPIDs
	pCtx := tss.NewPeerContext(sortedPIDs)

	// Channel for party events
	outCh := make(chan tss.Message, len(partyIndexes)*len(partyIndexes))
	endCh := make(chan *tsscommon.SignatureData, len(partyIndexes))
	errCh := make(chan *tss.Error, len(partyIndexes))

	// Create parties for the signing process
	for _, data := range partyData {
		params := tss.NewParameters(tss.S256(), pCtx, data.PartyID, len(sortedPIDs), mpccommon.Threshold)
		P := signing.NewLocalParty(msgHashInt, params, *data.SavedData, outCh, endCh).(*signing.LocalParty)
		parties[data.PartyID.Id] = P
		go func(p *signing.LocalParty) {
			if err := p.Start(); err != nil {
				errCh <- err
			}
		}(P)
	}

	// Coordinate message exchange between parties
	var signature *tsscommon.SignatureData
signing:
	for {
		select {
		case msg := <-outCh:
			dest := msg.GetTo()
			paredMsg := msg.(tss.ParsedMessage)
			if dest == nil { // Broadcast message
				for _, P := range parties {
					if P.PartyID().Index != msg.GetFrom().Index {
						go func(P *signing.LocalParty, parsedMsg tss.ParsedMessage) {
							if _, err := P.Update(parsedMsg); err != nil {
								errCh <- err
							}
						}(P, paredMsg)
					}
				}
			} else { // Unicast message
				if dest[0].Id != msg.GetFrom().Id {
					P := parties[dest[0].Id]
					if P == nil {
						fmt.Printf("Party %d not found\n", dest[0].Id)
						os.Exit(1)
					}
					go func(P *signing.LocalParty, parsedMsg tss.ParsedMessage) {
						if _, err := P.Update(parsedMsg); err != nil {
							errCh <- err
						}
					}(P, paredMsg)
				}
			}
		case err := <-errCh:
			fmt.Printf("Error: %s\n", err)
			break signing
		case sigData := <-endCh:
			signature = sigData
			break signing
		}
	}
	if signature == nil {
		return "", fmt.Errorf(mpccommon.ErrSigningFailed, "no signature received")
	}
	return hex.EncodeToString(signature.GetSignature()), nil
}
