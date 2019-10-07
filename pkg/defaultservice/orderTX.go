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
	"fmt"

	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/apache/incubator-milagro-dta/pkg/identity"
	"github.com/pkg/errors"
)

// Order2 - Process an incoming Blockchain Order transaction from a MasterFiduciary, to generate the final public key/address
func (s *Service) Order2(tx *api.BlockChainTX) (string, error) {
	nodeID := s.NodeID()
	reqPayload := tx.Payload
	txHashString := hex.EncodeToString(tx.TXhash)

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, s.MasterFiduciaryNodeID())
	if err != nil {
		return "", err
	}

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

	//Decode the Order from the supplied TX
	order := &documents.OrderDoc{}
	err = documents.DecodeOrderDocument(reqPayload, txHashString, order, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)

	//Generate commitment
	commitment, extension, err := s.Plugin.PrepareOrderResponse(order)
	if err != nil {
		return "", errors.Wrap(err, "Generating Final Public Key")
	}

	//TODO: Do something with the Commitment, which should only be visible to the Principal
	//For now, we will put it in a TX and broadcast with only the Principal as Recipients
	//The Processor for the TX is 'dump' - So the principal will pick up the TX and display
	//its contents.

	recipientList, err := common.BuildRecipientList(s.Ipfs, nodeID)
	if err != nil {
		return "", err
	}

	//Populate extension fields
	order.OrderPart2.CommitmentPublicKey = commitment
	if order.OrderPart2.Extension == nil {
		order.OrderPart2.Extension = make(map[string]string)
	}
	for key, value := range extension {
		order.OrderPart2.Extension[key] = value
	}

	//Generate a transaction
	txHash, payload, err := common.CreateTX(nodeID, s.Store, blsSK, nodeID, order, recipientList)

	//Write the Order2 results to the chain
	chainTX := &api.BlockChainTX{
		Processor:              api.TXOrderResponse,
		SenderID:               nodeID,
		RecipientID:            nodeID,
		AdditionalRecipientIDs: []string{},
		Payload:                payload,
		Tags:                   map[string]string{"reference": order.Reference, "txhash": hex.EncodeToString(txHash)},
	}

	return s.Tendermint.PostTx(chainTX, "Order2")
}

// OrderSecret2 - Process an incoming Blockchain Order/Secret transaction from a MasterFiduciary, to generate the final secret
func (s *Service) OrderSecret2(tx *api.BlockChainTX) (string, error) {
	nodeID := s.NodeID()
	reqPayload := tx.Payload
	txHashString := hex.EncodeToString(tx.TXhash)

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

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, s.MasterFiduciaryNodeID())
	if err != nil {
		return "", err
	}

	//Decode the Order from the supplied TX
	order := &documents.OrderDoc{}
	err = documents.DecodeOrderDocument(reqPayload, txHashString, order, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		fmt.Println("ERROR DEcode Order:", err)
		return "", err
	}

	if order.BeneficiaryCID != nodeID {
		return "", errors.New("Invalid Processor")
	}

	finalPrivateKey, _, extension, err := s.Plugin.ProduceFinalSecret(keyseed, sikeSK, order, order, nodeID)
	if err != nil {
		return "", err
	}

	//TODO: Do something with the Final Private, which should only be visible to the Beneficiary
	//For now, we will put it in a TX and broadcast with only the Beneficiary as Recipient
	//The Processor for the TX is 'dump' - So the Beneficiary will pick up the TX and display
	//its contents.

	if order.OrderPart4.Extension == nil {
		order.OrderPart4.Extension = make(map[string]string)
	}
	for key, value := range extension {
		order.OrderPart4.Extension[key] = value
	}
	order.OrderPart4.Extension["FinalPrivateKey"] = finalPrivateKey

	//Output Only to self for autoviewing
	recipientList, err := common.BuildRecipientList(s.Ipfs, nodeID)
	if err != nil {
		return "", err
	}
	txHash, payload, err := common.CreateTX(nodeID, s.Store, blsSK, nodeID, order, recipientList)

	//Write the requests to the chain
	chainTX := &api.BlockChainTX{
		Processor:              api.TXOrderSecretResponse, //NONE
		SenderID:               nodeID,
		RecipientID:            nodeID,
		AdditionalRecipientIDs: []string{},
		Payload:                payload,
		Tags:                   map[string]string{"reference": order.Reference, "txhash": hex.EncodeToString(txHash)},
	}

	return s.Tendermint.PostTx(chainTX, "OrderSecret2")
}
