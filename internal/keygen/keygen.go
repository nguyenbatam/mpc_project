// Package keygen implements the key generation process for the MPC system.
package keygen

import (
	"fmt"
	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/test"
	"sync/atomic"

	mpccommon "github.com/nguyenbatam/mpc_project/internal/common"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
)

// GenerateDistributedKey generates a distributed key among n parties with a threshold t.
func GenerateDistributedKey() error {
	threshold := mpccommon.Threshold
	pCount := mpccommon.PartyCount

	// Validate the threshold
	if threshold > pCount {
		return fmt.Errorf(mpccommon.ErrInvalidThreshold, pCount)
	}

	pIDs := tss.GenerateTestPartyIDs(pCount)

	p2pCtx := tss.NewPeerContext(pIDs)
	parties := make([]*keygen.LocalParty, 0, len(pIDs))

	errCh := make(chan *tss.Error, len(pIDs))
	outCh := make(chan tss.Message, len(pIDs))
	endCh := make(chan *keygen.LocalPartySaveData, len(pIDs))

	updater := test.SharedPartyUpdater

	// init the parties
	for i := 0; i < len(pIDs); i++ {
		var P *keygen.LocalParty
		params := tss.NewParameters(tss.S256(), p2pCtx, pIDs[i], len(pIDs), threshold)
		P = keygen.NewLocalParty(params, outCh, endCh).(*keygen.LocalParty)
		parties = append(parties, P)
		go func(P *keygen.LocalParty) {
			if err := P.Start(); err != nil {
				errCh <- err
			}
		}(P)
	}

	// PHASE: keygen
	var ended int32
keygen:
	for {
		select {
		case err := <-errCh:
			common.Logger.Errorf("Error: %s", err)
			break keygen

		case msg := <-outCh:
			dest := msg.GetTo()
			if dest == nil { // broadcast!
				for _, P := range parties {
					if P.PartyID().Index == msg.GetFrom().Index {
						continue
					}
					go updater(P, msg, errCh)
				}
			} else { // point-to-point!
				if dest[0].Index != msg.GetFrom().Index {

					go test.SharedPartyUpdater(parties[dest[0].Index], msg, errCh)
				}
			}

		case save := <-endCh:
			index, _ := save.OriginalIndex()
			err := mpccommon.SavePartyData(index, nil, save)
			if err == nil {
				fmt.Printf("Saved data for party %d\n", index)
			}
			atomic.AddInt32(&ended, 1)
			if atomic.LoadInt32(&ended) == int32(len(pIDs)) {
				break keygen
			}
		}
	}
	/*// Coordinate message exchange between parties and handle results
	go func() {
		for {
			select {
			case msg := <-outCh:
				dest := msg.GetTo()
				// Type assertion to tss.ParsedMessage
				//paredMsg := msg.(tss.ParsedMessage)
				if dest == nil { // Broadcast message
					for _, P := range parties {
						// Do not send message back to the sender
						if P.PartyID().Index != msg.GetFrom().Index {
							go func() {
								test.SharedPartyUpdater(P, msg, errCh)
							}()
							/*go func(P *keygen.LocalParty, parsedMsg tss.ParsedMessage) {
								if _, err := P.Update(parsedMsg); err != nil {
									errCh <- err
								}
							}(P, paredMsg)*/
	/*}
		}
	} else { // Unicast message
		if dest[0].Index != msg.GetFrom().Index {
			P := parties[dest[0].Index]
			go func() {
				test.SharedPartyUpdater(P, msg, errCh)
			}()
			/*go func(P *keygen.LocalParty, parsedMsg tss.ParsedMessage) {
				if _, err := P.Update(parsedMsg); err != nil {
					errCh <- err
				}
			}(P, paredMsg)*/
	/*}
				}
			case err := <-errCh:
				fmt.Printf("Error from TSS party: %s\n", err) // Log the error from a party
			case save := <-endCh:
				fmt.Println("Receive save info id :", save.ShareID.Int64())
				// Save party data
				// Convert ShareID (big.Int) to index
				index := int(save.ShareID.Int64()) % pCount
				err := mpccommon.SavePartyData(index, partyIDs[index], save)
				if err != nil {
					fmt.Printf("Failed to save data for party %d: %v\n", index, err)
				} else {
					fmt.Printf("Successfully saved data for party %d\n", index)
				}
				savedCount.Done() // Decrement the counter when data is successfully saved
			}
		}
	}()

	// Wait for all parties to finish their Start() method AND all party data to be saved
	partyWg.Wait()
	savedCount.Wait()*/

	fmt.Println("Key generation successful!") // Now this is printed after data is likely saved

	// Verify public key
	// Read data of the first party to get the public key
	data0, err := mpccommon.LoadPartyData(0)
	if err != nil {
		// If this fails after waiting for saves, it indicates a different issue.
		return fmt.Errorf("failed to read data for party 0 after saving: %v", err)
	}

	pubKey := data0.SavedData.ECDSAPub
	fmt.Printf("Generated public key: (%s, %s)\n", pubKey.X().String(), pubKey.Y().String())

	return nil
}
