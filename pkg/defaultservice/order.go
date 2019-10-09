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
	"time"

	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/apache/incubator-milagro-dta/pkg/identity"
	"github.com/pkg/errors"
)

// GetOrder retreives an order
func (s *Service) GetOrder(req *api.GetOrderRequest) (*api.GetOrderResponse, error) {
	// TODO: Get Order
	// orderReference := req.OrderReference

	// var cid string
	// if err := s.Store.Get("order", orderReference, &cid); err != nil {
	// 	return nil, err
	// }

	// localIDDoc, err := common.RetrieveIDDoc(s.Tendermint, s.NodeID())
	// if err != nil {
	// 	return nil, err
	// }

	// // SIKE key
	// keyseed, err := s.KeyStore.Get("seed")
	// if err != nil {
	// 	return nil, err
	// }
	// _, sikeSK, err := identity.GenerateSIKEKeys(keyseed)
	// if err != nil {
	// 	return nil, err
	// }

	// order, err := common.RetrieveOrderFromIPFS(s.Ipfs, cid, sikeSK, s.NodeID(), localIDDoc.BLSPublicKey)
	// if err != nil {
	// 	return nil, err
	// }

	// orderByte, err := json.Marshal(order)
	// if err != nil {
	// 	return nil, err
	// }
	// return &api.GetOrderResponse{
	// 	OrderCID: cid,
	// 	Order:    string(orderByte),
	// }, nil

	return nil, errors.New("Not implemented")
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
func (s *Service) PrepareOrderResponse(orderPart2 *documents.OrderDoc) (commitment string, extension map[string]string, err error) {
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

// Order creates a new deposit order
func (s *Service) Order(req *api.OrderRequest) (string, error) {
	if err := s.Plugin.ValidateOrderRequest(req); err != nil {
		return "", err
	}

	//Initialise values from Request object
	beneficiaryIDDocumentCID := req.BeneficiaryIDDocumentCID
	nodeID := s.NodeID()
	recipientList, err := common.BuildRecipientList(s.Tendermint, nodeID, s.MasterFiduciaryNodeID())
	if err != nil {
		return "", err
	}

	//Create Order
	order, err := common.CreateNewDepositOrder(beneficiaryIDDocumentCID, nodeID, s.Plugin.Name())
	if err != nil {
		return "", err
	}

	//Populate extension fields
	if order.OrderReqExtension == nil {
		order.OrderReqExtension = make(map[string]string)
	}
	for key, value := range req.Extension {
		order.OrderReqExtension[key] = value
	}

	// BLS key
	keyseed, err := s.KeyStore.Get("seed")
	if err != nil {
		return "", err
	}
	_, blsSK, err := identity.GenerateBLSKeys(keyseed)
	if err != nil {
		return "", err
	}

	//Create a new Transaction payload and TX
	tx, txID, err := common.CreateOrderTX(nodeID, api.TXFulfillRequest, s.Store, blsSK, order, recipientList, s.MasterFiduciaryNodeID())
	s.Logger.Debug("Created new %v Tx: %v", api.TXFulfillRequest, txID)

	if _, err := s.Tendermint.PostTx(tx); err != nil {
		return "", err
	}

	return order.Reference, nil
}

// OrderSecret initiates the process for retreiving the order secret
func (s *Service) OrderSecret(req *api.OrderSecretRequest) (string, error) {
	orderReference := req.OrderReference
	var previousOrderHash string
	if err := s.Store.Get("order", orderReference, &previousOrderHash); err != nil {
		return "", err
	}

	nodeID := s.NodeID()
	recipientList, err := common.BuildRecipientList(s.Tendermint, nodeID, s.MasterFiduciaryNodeID())
	if err != nil {
		return "", err
	}
	remoteIDDoc, err := common.RetrieveIDDoc(s.Tendermint, s.MasterFiduciaryNodeID())
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

	localIDDoc, err := common.RetrieveIDDoc(s.Tendermint, s.NodeID())
	if err != nil {
		return "", err
	}

	tx, err := s.Tendermint.GetTx(previousOrderHash)
	if err != nil {
		return "", err
	}

	order := &documents.OrderDoc{}
	err = documents.DecodeOrderDocument(tx.Payload, previousOrderHash, order, sikeSK, nodeID, remoteIDDoc.BLSPublicKey)
	if err != nil {
		err = documents.DecodeOrderDocument(tx.Payload, previousOrderHash, order, sikeSK, nodeID, localIDDoc.BLSPublicKey)
		if err != nil {
			return "", errors.Wrap(err, "Fail to retrieve existing order")
		}
	}

	//Populate extension fields
	if order.OrderSecretReqExtension == nil {
		order.OrderSecretReqExtension = make(map[string]string)
	}
	for key, value := range req.Extension {
		order.OrderSecretReqExtension[key] = value
	}

	if err := s.Plugin.ValidateOrderSecretRequest(req, *order); err != nil {
		return "", err
	}

	if req.BeneficiaryIDDocumentCID != "" {
		order.BeneficiaryCID = req.BeneficiaryIDDocumentCID
	}

	//Create a piece of data that is destined for the beneficiary, passed via the Master Fiduciary
	beneficiaryEncryptedData, _, err := s.Plugin.ProduceBeneficiaryEncryptedData(blsSK, order, req)
	if err != nil {
		return "", err
	}

	//Create a request Object
	order.OrderPart3 = &documents.OrderPart3{
		Redemption:               "SignedReferenceNumber",
		PreviousOrderCID:         previousOrderHash,
		BeneficiaryEncryptedData: beneficiaryEncryptedData,
		Timestamp:                time.Now().Unix(),
	}

	tx, txID, err := common.CreateOrderTX(nodeID, api.TXFulfillOrderSecretRequest, s.Store, blsSK, order, recipientList, s.MasterFiduciaryNodeID())
	s.Logger.Debug("Created new %v Tx: %v", api.TXFulfillOrderSecretRequest, txID)

	if _, err := s.Tendermint.PostTx(tx); err != nil {
		return "", err
	}

	return order.Reference, nil
}
