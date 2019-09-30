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

	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/apache/incubator-milagro-dta/pkg/tendermint"
	"github.com/pkg/errors"
)

// GetOrder retreives an order
func (s *Service) GetOrder(req *api.GetOrderRequest) (*api.GetOrderResponse, error) {
	orderReference := req.OrderReference

	var cid string
	if err := s.Store.Get("order", orderReference, &cid); err != nil {
		return nil, err
	}

	localIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, s.NodeID())
	if err != nil {
		return nil, err
	}

	_, _, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, s.NodeID())
	if err != nil {
		return nil, err
	}

	order, err := common.RetrieveOrderFromIPFS(s.Ipfs, cid, sikeSK, s.NodeID(), localIDDoc.BLSPublicKey)
	if err != nil {
		return nil, err
	}

	orderByte, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	return &api.GetOrderResponse{
		OrderCID: cid,
		Order:    string(orderByte),
	}, nil
}

// OrderList retrieves the list of orders
func (s *Service) OrderList(req *api.OrderListRequest) (*api.OrderListResponse, error) {
	page := req.Page
	perPage := req.PerPage
	sortBy := req.SortBy

	orderref, err := s.Store.ListKeys("order", "time", page*perPage, perPage, sortBy != "dateCreatedAsc")
	if err != nil {
		return nil, err
	}

	//Pagnination - Show everything by default
	start := 0
	stop := len(orderref)

	if perPage != 0 && page < len(orderref)/perPage && page*perPage < len(orderref) {
		start = page
		stop = perPage
	}

	return &api.OrderListResponse{
		OrderReference: orderref[start:stop],
	}, nil
}

// ValidateOrderRequest returns error if the request values are invalid
func (s *Service) ValidateOrderRequest(req *api.OrderRequest) error {
	return nil
}

//ValidateOrderSecretRequest - Validate fields in the Order Secret
func (s *Service) ValidateOrderSecretRequest(req *api.OrderSecretRequest, order documents.OrderDoc) error {
	return nil
}

// PrepareOrderPart1 is called before the order is send
func (s *Service) PrepareOrderPart1(order *documents.OrderDoc, reqExtension map[string]string) (fulfillExtension map[string]string, err error) {
	return nil, nil
}

// PrepareOrderResponse gets the updated order and returns the commitment and extension
func (s *Service) PrepareOrderResponse(orderPart2 *documents.OrderDoc, reqExtension, fulfillExtension map[string]string) (commitment string, extension map[string]string, err error) {
	return orderPart2.OrderPart2.CommitmentPublicKey, nil, nil
}

// ProduceBeneficiaryEncryptedData -
func (s *Service) ProduceBeneficiaryEncryptedData(blsSK []byte, order *documents.OrderDoc, req *api.OrderSecretRequest) (encrypted []byte, extension map[string]string, err error) {
	return nil, nil, nil
}

// ProduceFinalSecret -
func (s *Service) ProduceFinalSecret(seed, sikeSK []byte, order, orderPart4 *documents.OrderDoc, beneficiaryIDDocumentCID string) (secret, commitment string, extension map[string]string, err error) {
	finalPrivateKey := orderPart4.OrderDocument.OrderPart4.Secret
	//Derive the Public key from the supplied Private Key
	finalPublicKey, _, err := cryptowallet.PublicKeyFromPrivate(finalPrivateKey)
	return finalPrivateKey, finalPublicKey, nil, err
}

// Order1 -
func (s *Service) Order1(req *api.OrderRequest) (string, error) {
	if err := s.Plugin.ValidateOrderRequest(req); err != nil {
		return "", err
	}

	//Initialise values from Request object
	beneficiaryIDDocumentCID := req.BeneficiaryIDDocumentCID
	nodeID := s.NodeID()
	recipientList, err := common.BuildRecipientList(s.Ipfs, nodeID, s.MasterFiduciaryNodeID())
	if err != nil {
		return "", err
	}

	//Create Order
	order, err := common.CreateNewDepositOrder(beneficiaryIDDocumentCID, nodeID)
	if err != nil {
		return "", err
	}

	fulfillExtension, err := s.Plugin.PrepareOrderPart1(order, req.Extension)
	if err != nil {
		return "", err
	}

	//Write Order to IPFS
	orderPart1CID, err := common.WriteOrderToIPFS(nodeID, s.Ipfs, s.Store, nodeID, order, recipientList)
	if err != nil {
		return "", err
	}

	request := &api.FulfillOrderRequest{
		DocumentCID:   nodeID,
		OrderPart1CID: orderPart1CID,
		Extension:     fulfillExtension,
	}

	//This is serialized and output to the chain
	txHash, payload, err := common.CreateTX(nodeID, s.Store, nodeID, order, recipientList)

	marshaledRequest, _ := json.Marshal(request)
	_ = marshaledRequest

	//Write the requests to the chain
	chainTX := &api.BlockChainTX{
		Processor:   api.TXFulfillRequest,
		SenderID:    nodeID,
		RecipientID: []string{s.MasterFiduciaryNodeID(), nodeID},
		Payload:     payload, //marshaledRequest,
		TXhash:      txHash,
		Tags:        map[string]string{"reference": order.Reference},
	}
	tendermint.PostToChain(chainTX, "Order1")
	return order.Reference, nil
}

// Order2 -
func (s *Service) Order2(tx *api.BlockChainTX) (string, error) {
	nodeID := s.NodeID()
	reqPayload := tx.Payload
	txHashString := hex.EncodeToString(tx.TXhash)

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, s.MasterFiduciaryNodeID())
	if err != nil {
		return "", err
	}

	//Get the updated order out of IPFS
	_, _, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return "", err
	}

	order := &documents.OrderDoc{}
	err = documents.DecodeOrderDocument(reqPayload, txHashString, order, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)

	// updatedOrder, err := common.RetrieveOrderFromIPFS(s.Ipfs, req.OrderPart2CID, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	// if err != nil {
	// 	return "", errors.Wrap(err, "Fail to retrieve Order from IPFS")
	// }

	//update OrderPartCID for order id

	commitment, extension, err := s.Plugin.PrepareOrderResponse(order)
	if err != nil {
		return "", errors.Wrap(err, "Generating Final Public Key")
	}

	order.OrderPart2.CommitmentPublicKey = commitment

	//Populate Extension

	if order.OrderPart2.Extension == nil {
		order.OrderPart2.Extension = make(map[string]string)
	}
	for key, value := range extension {

		order.OrderPart2.Extension[key] = value
	}

	// err = common.WriteOrderToStore(s.Store, order.Reference, req.OrderPart2CID)
	// if err != nil {
	// 	return "", errors.Wrap(err, "Saving new CID to Order reference")
	// }

	// response := &api.OrderResponse{
	// 	OrderReference: order.Reference,
	// 	Commitment:     commitment,
	// 	CreatedAt:      time.Now().Unix(),
	// 	Extension:      extension,
	// }
	recipientList, err := common.BuildRecipientList(s.Ipfs, nodeID, nodeID)
	if err != nil {
		return "", err
	}
	txHash, payload, err := common.CreateTX(nodeID, s.Store, nodeID, order, recipientList)

	//Write the Order2 results to the chain
	chainTX := &api.BlockChainTX{
		Processor:   api.TXOrderResponse,
		SenderID:    "", //use no Sender so we can read our own Result for testing
		RecipientID: []string{nodeID},
		Payload:     payload,
		TXhash:      txHash,
		Tags:        map[string]string{"reference": order.Reference},
	}
	return tendermint.PostToChain(chainTX, "Order2")

}

// OrderSecret1 -
func (s *Service) OrderSecret1(req *api.OrderSecretRequest) (string, error) {
	orderReference := req.OrderReference
	var orderPart2CID string
	if err := s.Store.Get("order", orderReference, &orderPart2CID); err != nil {
		return "", err
	}

	nodeID := s.NodeID()
	recipientList, err := common.BuildRecipientList(s.Ipfs, nodeID, s.MasterFiduciaryNodeID())
	if err != nil {
		return "", err
	}
	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, s.MasterFiduciaryNodeID())
	if err != nil {
		return "", err
	}

	_, _, blsSK, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return "", err
	}

	localIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, s.NodeID())
	if err != nil {
		return "", err
	}

	//If we already did a transfer the Order doc is self signed so, check with own Key so we can re-process the transfer
	order, err := common.RetrieveOrderFromIPFS(s.Ipfs, orderPart2CID, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		//check if we are re-trying the call, so the OrderDoc is locally signed
		order, err = common.RetrieveOrderFromIPFS(s.Ipfs, orderPart2CID, sikeSK, nodeID, localIDDoc.BLSPublicKey)
		if err != nil {
			return "", errors.Wrap(err, "Fail to retrieve Order from IPFS")
		}
	}

	if err := s.Plugin.ValidateOrderSecretRequest(req, *order); err != nil {
		return "", err
	}

	if req.BeneficiaryIDDocumentCID != "" {
		order.BeneficiaryCID = req.BeneficiaryIDDocumentCID
	}

	//Create a piece of data that is destined for the beneficiary, passed via the Master Fiduciary
	beneficiaryEncryptedData, extension, err := s.Plugin.ProduceBeneficiaryEncryptedData(blsSK, order, req)
	if err != nil {
		return "", err
	}

	//Create a request Object in IPFS
	orderPart3CID, err := common.CreateAndStorePart3(s.Ipfs, s.Store, order, orderPart2CID, nodeID, beneficiaryEncryptedData, recipientList)
	if err != nil {
		return "", err
	}

	//Post the address of the updated doc to the custody node
	request := &api.FulfillOrderSecretRequest{
		SenderDocumentCID: nodeID,
		OrderPart3CID:     orderPart3CID,
		Extension:         extension,
	}

	marshaledRequest, _ := json.Marshal(request)

	//Write the requests to the chain
	chainTX := &api.BlockChainTX{
		Processor:   api.TXFulfillOrderSecretRequest,
		SenderID:    nodeID,
		RecipientID: []string{s.MasterFiduciaryNodeID()},
		Payload:     marshaledRequest,
		Tags:        map[string]string{"reference": order.Reference},
	}
	return tendermint.PostToChain(chainTX, "OrderSecret1")
}

// OrderSecret2 -
func (s *Service) OrderSecret2(req *api.FulfillOrderSecretResponse) (string, error) {
	nodeID := s.NodeID()
	_, _, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return "", err
	}

	remoteIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, s.MasterFiduciaryNodeID())
	if err != nil {
		return "", err
	}

	//Retrieve the response Order from IPFS
	orderPart4, err := common.RetrieveOrderFromIPFS(s.Ipfs, req.OrderPart4CID, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)

	if orderPart4.BeneficiaryCID != nodeID {
		return "", errors.New("Invalid Processor")
	}

	_, seed, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return "", err
	}

	finalPrivateKey, finalPublicKey, ext, err := s.Plugin.ProduceFinalSecret(seed, sikeSK, orderPart4, orderPart4, nodeID)
	if err != nil {
		return "", err
	}

	request := &api.OrderSecretResponse{
		Secret:         finalPrivateKey,
		Commitment:     finalPublicKey,
		OrderReference: orderPart4.Reference,
		Extension:      ext,
	}

	marshaledRequest, _ := json.Marshal(request)

	//Write the requests to the chain
	chainTX := &api.BlockChainTX{
		Processor:   api.TXOrderSecretResponse, //NONE
		SenderID:    "",                        // so we can view it
		RecipientID: []string{nodeID},          //don't send this to chain, seed compromise becomes fatal, sent just debugging
		Payload:     marshaledRequest,
		Tags:        map[string]string{"reference": orderPart4.Reference, "part": "4"},
	}
	return tendermint.PostToChain(chainTX, "OrderSecret2")
}
