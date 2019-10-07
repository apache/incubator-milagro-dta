package common

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/pkg/errors"
)

// CreateTX creates the transaction ready for write to the chain
func CreateTX(nodeID string, store *datastore.Store, blsSecretKey []byte, id string, order *documents.OrderDoc, recipients map[string]*documents.IDDoc) ([]byte, []byte, error) {
	rawDoc, err := documents.EncodeOrderDocument(nodeID, *order, blsSecretKey, recipients)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to encode IDDocument")
	}
	TXID := sha256.Sum256(rawDoc)
	TXIDhex := hex.EncodeToString(TXID[:])
	//Write order to store
	if err := WriteOrderToStore(store, order.Reference, TXIDhex); err != nil {
		return nil, nil, errors.New("Save Order to store")
	}
	return TXID[:], rawDoc, nil
}

//PeekTX Decode a transaction for header data but don't decrypt it
func PeekTX(tx []byte) (string, error) {
	signerCID, err := documents.OrderDocumentSigner(tx)
	if err != nil {
		return "", errors.Wrap(err, "Error peeking signer")
	}
	return signerCID, nil

}
