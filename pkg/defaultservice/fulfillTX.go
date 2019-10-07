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

package defaultservice

import (
	"encoding/hex"
	"time"

	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/apache/incubator-milagro-dta/pkg/identity"
)

// FulfillOrder TX
func (s *Service) FulfillOrder(tx *api.BlockChainTX) (string, error) {
	nodeID := s.NodeID()
	reqPayload := tx.Payload
	txHashString := hex.EncodeToString(tx.TXhash)

	//Get signer by peeking inside the document
	signerID, err := documents.OrderDocumentSigner(reqPayload)
	if err != nil {
		return "", err
	}
	remoteIDDocCID := signerID

	// SIKE and BLS keys
	keyseed, err := s.KeyStore.Get("seed")
	if err != nil {
		return "", err
	}
	_, sikeSK, err := identity.GenerateSIKEKeys(keyseed)
	if err != nil {
		return "", err
	}
	_, blsSK, err := identity.GenerateBLSKeys(keyseed)
	if err != nil {
		return "", err
	}

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, remoteIDDocCID)
	if err != nil {
		return "", err
	}

	//Decode the Order from the supplied TX
	order := &documents.OrderDoc{}
	err = documents.DecodeOrderDocument(reqPayload, txHashString, order, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		return "", err
	}

	//Recipient list is principal and self
	recipientList, err := common.BuildRecipientList(s.Ipfs, order.PrincipalCID, nodeID)
	if err != nil {
		return "", err
	}

	//Generate the secret and store for later redemption
	seed, err := common.MakeRandomSeedAndStore(s.Store, s.Rng, order.Reference)
	if err != nil {
		return "", err
	}

	//Generate the Public Key (Commitment) from the Seed/Secret
	commitmentPublicKey, err := cryptowallet.RedeemPublicKey(seed)
	if err != nil {
		return "", err
	}

	//Populate Order part 2
	order.OrderPart2 = &documents.OrderPart2{
		CommitmentPublicKey: commitmentPublicKey,
		PreviousOrderCID:    txHashString,
		Timestamp:           time.Now().Unix(),
	}

	//Create a new Transaction payload and TX
	txHash, payload, err := common.CreateTX(nodeID, s.Store, blsSK, nodeID, order, recipientList)

	//Write the requests to the chain
	chainTX := &api.BlockChainTX{
		Processor:              api.TXFulfillResponse,
		SenderID:               nodeID,
		RecipientID:            order.PrincipalCID,
		AdditionalRecipientIDs: []string{},
		Payload:                payload,
		TXhash:                 txHash,
		Tags:                   map[string]string{"reference": order.Reference, "txhash": hex.EncodeToString(txHash)},
	}

	return s.Tendermint.PostTx(chainTX, "FulfillOrder")
}

// FulfillOrderSecret -
func (s *Service) FulfillOrderSecret(tx *api.BlockChainTX) (string, error) {
	nodeID := s.NodeID()
	reqPayload := tx.Payload
	txHashString := hex.EncodeToString(tx.TXhash)

	//Get signer by peeking inside the document
	signerID, err := documents.OrderDocumentSigner(reqPayload)
	if err != nil {
		return "", err
	}
	remoteIDDocCID := signerID

	// SIKE and BLS keys
	keyseed, err := s.KeyStore.Get("seed")
	if err != nil {
		return "", err
	}
	_, sikeSK, err := identity.GenerateSIKEKeys(keyseed)
	if err != nil {
		return "", err
	}
	_, blsSK, err := identity.GenerateBLSKeys(keyseed)
	if err != nil {
		return "", err
	}

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, remoteIDDocCID)
	if err != nil {
		return "", err
	}

	//Decode the Order from the supplied TX
	order := &documents.OrderDoc{}
	err = documents.DecodeOrderDocument(reqPayload, txHashString, order, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		return "", err
	}

	//Recipient list is beneficiary and self
	recipientList, err := common.BuildRecipientList(s.Ipfs, nodeID, order.BeneficiaryCID)
	if err != nil {
		return "", err
	}

	//Retrieve the Seed
	seed, err := common.RetrieveSeed(s.Store, order.Reference)
	if err != nil {
		return "", err
	}

	//Generate the Secert from the Seed
	commitmentPrivateKey, err := cryptowallet.RedeemSecret(seed)
	if err != nil {
		return "", err
	}

	//Populate Order part 4
	order.OrderPart4 = &documents.OrderPart4{
		Secret:           commitmentPrivateKey,
		PreviousOrderCID: txHashString,
		Timestamp:        time.Now().Unix(),
	}

	//Create a new Transaction payload and TX
	txHash, payload, err := common.CreateTX(nodeID, s.Store, blsSK, nodeID, order, recipientList)

	//Write the requests to the chain
	chainTX := &api.BlockChainTX{
		Processor:              api.TXFulfillOrderSecretResponse,
		SenderID:               nodeID,
		RecipientID:            order.BeneficiaryCID,
		AdditionalRecipientIDs: []string{},

		Payload: payload,
		Tags:    map[string]string{"reference": order.Reference, "txhash": hex.EncodeToString(txHash)},
	}

	return s.Tendermint.PostTx(chainTX, "FulfillOrderSecret")
}
