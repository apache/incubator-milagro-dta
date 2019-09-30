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
	"encoding/json"
	"time"

	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/apache/incubator-milagro-dta/pkg/tendermint"
)

func (s *Service) BCFulfillOrder(tx *api.BlockChainTX) (string, error) {

	reqPayload := tx.Payload
	txHashString := hex.EncodeToString(tx.TXhash)

	//Decode the incoming TX
	//Peek inside the TX
	//Pull out the header - to get PrincipalID
	signerID, err := documents.OrderDocumentSigner(reqPayload)
	if err != nil {
		return "", err
	}

	//orderPart1CID := req.OrderPart1CID
	nodeID := s.NodeID()
	remoteIDDocCID := signerID

	_, _, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return "", err
	}

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, remoteIDDocCID)
	if err != nil {
		return "", err
	}

	//Retrieve the order from IPFS
	order := &documents.OrderDoc{}
	err = documents.DecodeOrderDocument(reqPayload, txHashString, order, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)

	//	order, err := common.RetrieveOrderFromIPFS(s.Ipfs, orderPart1CID, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		return "", err
	}

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

	//Create an order part 2
	Part2 := documents.OrderPart2{
		CommitmentPublicKey: commitmentPublicKey,
		PreviousOrderCID:    txHashString,
		Timestamp:           time.Now().Unix(),
	}
	order.OrderPart2 = &Part2

	txHash, payload, err := common.CreateTX(nodeID, s.Store, nodeID, order, recipientList)
	//_ = txHashID

	// orderPart2CID, err := common.CreateAndStoreOrderPart2(s.Ipfs, s.Store, order, orderPart1CID, commitmentPublicKey, nodeID, recipientList)
	// if err != nil {
	// 	return "", err
	// }

	//marshaledRequest, _ := json.Marshal(response)

	//Write the requests to the chain
	chainTX := &api.BlockChainTX{
		Processor:   api.TXFulfillResponse,
		SenderID:    nodeID,
		RecipientID: []string{order.PrincipalCID, nodeID},
		Payload:     payload,
		TXhash:      txHash,
		Tags:        map[string]string{"reference": order.Reference},
	}
	return tendermint.PostToChain(chainTX, "FulfillOrder")

}

// FulfillOrder -
func (s *Service) FulfillOrder(req *api.FulfillOrderRequest) (string, error) {
	orderPart1CID := req.OrderPart1CID
	nodeID := s.NodeID()
	remoteIDDocCID := req.DocumentCID
	_, _, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return "", err
	}

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, remoteIDDocCID)
	if err != nil {
		return "", err
	}

	//Retrieve the order from IPFS
	order, err := common.RetrieveOrderFromIPFS(s.Ipfs, orderPart1CID, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		return "", err
	}

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

	//Create an order response in IPFS
	orderPart2CID, err := common.CreateAndStoreOrderPart2(s.Ipfs, s.Store, order, orderPart1CID, commitmentPublicKey, nodeID, recipientList)
	if err != nil {
		return "", err
	}

	response := &api.FulfillOrderResponse{
		OrderPart2CID: orderPart2CID,
	}

	marshaledRequest, _ := json.Marshal(response)

	//Write the requests to the chain
	chainTX := &api.BlockChainTX{
		Processor:   api.TXFulfillResponse,
		SenderID:    nodeID,
		RecipientID: []string{order.PrincipalCID, nodeID},
		Payload:     marshaledRequest,
		Tags:        map[string]string{"reference": order.Reference},
	}
	return tendermint.PostToChain(chainTX, "FulfillOrder")

}

// FulfillOrderSecret -
func (s *Service) FulfillOrderSecret(req *api.FulfillOrderSecretRequest) (string, error) {
	//Initialise values from Request object
	orderPart3CID := req.OrderPart3CID
	nodeID := s.NodeID()
	remoteIDDocCID := req.SenderDocumentCID
	_, _, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return "", err
	}

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, remoteIDDocCID)
	if err != nil {
		return "", err
	}

	//Retrieve the order from IPFS
	order, err := common.RetrieveOrderFromIPFS(s.Ipfs, orderPart3CID, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		return "", err
	}

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

	//Create an order response in IPFS
	orderPart4ID, err := common.CreateAndStoreOrderPart4(s.Ipfs, s.Store, order, commitmentPrivateKey, orderPart3CID, nodeID, recipientList)
	if err != nil {
		return "", err
	}

	response := &api.FulfillOrderSecretResponse{
		OrderPart4CID: orderPart4ID,
	}

	marshaledRequest, _ := json.Marshal(response)

	//Write the requests to the chain
	chainTX := &api.BlockChainTX{
		Processor:   api.TXFulfillOrderSecretResponse,
		SenderID:    nodeID,
		RecipientID: []string{s.MasterFiduciaryNodeID(), order.BeneficiaryCID},
		Payload:     marshaledRequest,
		Tags:        map[string]string{"reference": order.Reference},
	}
	return tendermint.PostToChain(chainTX, "FulfillOrderSecret")
}
