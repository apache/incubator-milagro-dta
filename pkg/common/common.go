// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

/*
Package common - helper functions that enable service to get and set encrypted envelopes
*/
package common

import (
	"encoding/hex"
	"io"
	"time"

	"github.com/apache/incubator-milagro-dta/pkg/tendermint"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// CreateNewDepositOrder - Generate an empty new Deposit Order with random reference
func CreateNewDepositOrder(BeneficiaryIDDocumentCID string, nodeID string, orderType string) (*documents.OrderDoc, error) {
	//Create a reference for this order
	reference, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	order := documents.NewOrderDoc()

	if BeneficiaryIDDocumentCID == "" {
		order.BeneficiaryType = documents.OrderDocument_Beneficiary_Unknown_at_Start
	} else {
		order.BeneficiaryType = documents.OrderDocument_Beneficiary_Known_at_start
	}

	//oder.Type will be used to extend the things that an order can do.
	order.Type = orderType
	order.PrincipalCID = nodeID
	order.Reference = reference.String()
	order.BeneficiaryCID = BeneficiaryIDDocumentCID
	order.Timestamp = time.Now().Unix()
	return &order, nil
}

// RetrieveIDDoc finds and parses the IDDocument
func RetrieveIDDoc(tmConn *tendermint.NodeConnector, id string) (*documents.IDDoc, error) {
	rawDocI, err := tmConn.GetTx(id)
	if err != nil {
		return nil, err
	}
	iddoc := &documents.IDDoc{}
	err = documents.DecodeIDDocument(rawDocI.Payload, id, iddoc)
	return iddoc, err
}

// MakeRandomSeedAndStore genefates and stores a random seed
func MakeRandomSeedAndStore(store *datastore.Store, rng io.Reader, reference string) (seedHex string, err error) {
	seed := make([]byte, 32)
	if _, err := io.ReadFull(rng, seed); err != nil {
		return "", err
	}
	i := len(seed)
	if i > 32 {
		i = 32
	}
	var byte32 [32]byte
	copy(byte32[:], seed[:i])
	seedHex = hex.EncodeToString(seed)
	if err := store.Set("keySeed", reference, seedHex, nil); err != nil {
		return "", errors.Wrap(err, "store seed")
	}
	return seedHex, nil
}

// RetrieveSeed gets the seed from the key store
func RetrieveSeed(store *datastore.Store, reference string) (seedHex string, err error) {
	if err := store.Get("keySeed", reference, &seedHex); err != nil {
		return "", nil
	}
	return seedHex, nil
}

//WriteOrderToStore stores an order
func WriteOrderToStore(store *datastore.Store, orderReference string, address string) error {
	if err := store.Set("order", orderReference, address, map[string]string{"time": time.Now().UTC().Format(time.RFC3339)}); err != nil {
		return errors.New("Save Order to store")
	}
	return nil
}

// BuildRecipientList builds a list of recipients who are able to decrypt the encrypted envelope
func BuildRecipientList(tmConn *tendermint.NodeConnector, ids ...string) (map[string]*documents.IDDoc, error) {
	recipients := make(map[string]*documents.IDDoc)
	for _, v := range ids {
		iddoc, err := RetrieveIDDoc(tmConn, v)
		if err != nil {
			return nil, err
		}
		recipients[v] = iddoc
	}
	return recipients, nil
}
