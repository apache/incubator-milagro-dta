package common

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/pkg/errors"
)

// CreateTX creates the transaction ready for write to the chain
func CreateTX(nodeID string, store *datastore.Store, id string, order *documents.OrderDoc, recipients map[string]documents.IDDoc) ([]byte, []byte, error) {
	secrets := &IdentitySecrets{}
	if err := store.Get("id-doc", nodeID, secrets); err != nil {
		return nil, nil, errors.New("load secrets from store")
	}
	blsSecretKey, err := hex.DecodeString(secrets.BLSSecretKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Decode identity secrets")
	}
	rawDoc, err := documents.EncodeOrderDocument(nodeID, *order, blsSecretKey, "previousID", recipients)
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

//Decode a transaction for header data but don't decrypt it
func PeekTX(tx []byte) (string, error) {
	signerCID, err := documents.OrderDocumentSigner(tx)
	print(signerCID)
	if err != nil {
		return "", errors.Wrap(err, "Error peeking signer")
	}
	return signerCID, nil

}
