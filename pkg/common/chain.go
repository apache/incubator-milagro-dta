package common

import (
	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/pkg/errors"
)

// CreateOrderTX creates the transaction ready for write to the chain
func CreateOrderTX(
	nodeID string,
	processor string,
	store *datastore.Store,
	blsSecretKey []byte,
	order *documents.OrderDoc,
	recipientDocs map[string]*documents.IDDoc,
	recipientID string,
) (*api.BlockChainTX, string, error) {
	rawDoc, err := documents.EncodeOrderDocument(nodeID, *order, blsSecretKey, recipientDocs)
	if err != nil {
		return nil, "", errors.Wrap(err, "Failed to encode IDDocument")
	}

	tx, txID := api.NewBlockChainTX(processor, rawDoc, order.Reference, nodeID, recipientID)

	//Write order to store
	if err := WriteOrderToStore(store, order.Reference, txID); err != nil {
		return nil, "", errors.New("Save Order to store")
	}

	return tx, txID, nil
}

//PeekTX Decode a transaction for header data but don't decrypt it
func PeekTX(tx []byte) (string, error) {
	signerCID, err := documents.OrderDocumentSigner(tx)
	if err != nil {
		return "", errors.Wrap(err, "Error peeking signer")
	}
	return signerCID, nil

}
