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
Package identity - manage Identity document and keys
*/
package identity

import (
	"bytes"
	"time"

	"github.com/apache/incubator-milagro-dta/pkg/api"

	"github.com/apache/incubator-milagro-dta/pkg/tendermint"

	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/libs/keystore"
	"github.com/pkg/errors"
)

// CreateIdentity creates a new identity
// returns Identity document and secret
func CreateIdentity(name string) (idDocument *documents.IDDoc, rawIDDoc, seed []byte, err error) {
	//generate crypto random seed
	seed, err = cryptowallet.RandomBytes(48)
	if err != nil {
		err = errors.Wrap(err, "Failed to generate random seed")
		return
	}

	sikePublicKey, _, err := GenerateSIKEKeys(seed)
	if err != nil {
		return
	}

	blsPublicKey, blsSecretKey, err := GenerateBLSKeys(seed)
	if err != nil {
		return
	}

	ecPublicKey, err := GenerateECPublicKey(seed)
	if err != nil {
		return
	}

	// build ID Doc
	idDocument = documents.NewIDDoc()
	idDocument.AuthenticationReference = name
	idDocument.BeneficiaryECPublicKey = ecPublicKey
	idDocument.SikePublicKey = sikePublicKey
	idDocument.BLSPublicKey = blsPublicKey
	idDocument.Timestamp = time.Now().Unix()

	// encode ID Doc
	rawIDDoc, err = documents.EncodeIDDocument(idDocument, blsSecretKey)
	if err != nil {
		err = errors.Wrap(err, "Failed to encode IDDocument")
		return
	}

	return
}

// StoreIdentity writes IDDocument to blockchain and secret to keystore
func StoreIdentity(rawIDDoc, secret []byte, tmConn *tendermint.NodeConnector, store keystore.Store) (idDocumentID string, err error) {

	chainTX := &api.BlockChainTX{
		Processor:              api.TXFulfillOrderSecretRequest,
		SenderID:               tmConn.NodeID(),
		RecipientID:            "",
		AdditionalRecipientIDs: []string{},
		Payload:                rawIDDoc,
		Tags:                   map[string]string{},
	}

	idDocumentID = chainTX.CalcHash()

	if _, err := tmConn.PostTx(chainTX, "IDDocument"); err != nil {
		return "", err
	}

	// store the seed
	return idDocumentID, store.Set("seed", secret)
}

// CheckIdentity verifies the IDDocument
func CheckIdentity(id, name string, tmConn *tendermint.NodeConnector, store keystore.Store) error {

	rawIDDoc, err := tmConn.GetTx(id)
	if err != nil {
		return errors.Wrap(err, "ID Document not found")
	}

	idDoc := &documents.IDDoc{}
	if err := documents.DecodeIDDocument(rawIDDoc.Payload, id, idDoc); err != nil {
		return errors.Wrap(err, "Decode ID document")
	}

	if idDoc.AuthenticationReference != name {
		return errors.New("Name doesn't match the authentication reference")
	}

	seed, err := store.Get("seed")
	if err != nil {
		return errors.Wrap(err, "Seed not found")
	}

	sikePublic, _, err := GenerateSIKEKeys(seed)
	if !bytes.Equal(idDoc.SikePublicKey, sikePublic) {
		return errors.New("SIKE keys are different")
	}
	blsPublic, _, err := GenerateBLSKeys(seed)
	if !bytes.Equal(idDoc.BLSPublicKey, blsPublic) {
		return errors.New("BLS keys are different")
	}

	return nil
}
