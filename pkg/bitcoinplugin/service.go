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
Package bitcoinplugin - Milagro D-TA plugin that generates bitcoin addresses
*/
package bitcoinplugin

import (
	"strconv"

	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/apache/incubator-milagro-dta/pkg/defaultservice"

	"github.com/pkg/errors"
)

var (
	extensionVendor = "Milagro"
	pluginName      = "bitcoinwallet"
)

// Service is the Milagro bitcoin service
type Service struct {
	defaultservice.Service
}

// NewService returns a Milagro implementation of Service
func NewService() *Service {
	return &Service{}
}

// Name of the plugin
func (s *Service) Name() string {
	return pluginName
}

// Vendor of the plugin
func (s *Service) Vendor() string {
	return extensionVendor
}

// ValidateOrderRequest checks if the Coin type is valid
func (s *Service) ValidateOrderRequest(req *api.OrderRequest) error {
	if _, err := strconv.ParseInt(req.Extension["coin"], 10, 64); err != nil {
		return errors.Wrap(err, "Failed to Parse Coin Type")
	}

	return nil
}

//ValidateOrderSecretRequest - checks incoming OrderSecret fields for Error, comparing to the Original Order
func (s *Service) ValidateOrderSecretRequest(req *api.OrderSecretRequest, order documents.OrderDoc) error {
	//These are deliberately overly long winded, but it makes the case I'm trapping more obvious to the reader

	//There is no beneficiary supplided in either the Deposit or Redemption
	if order.BeneficiaryType == documents.OrderDocument_Unspecified && req.BeneficiaryIDDocumentCID == "" {
		return errors.New("Beneficiary must be supplied")
	}

	//A beneficiary is specified in both, but they aren't the same
	if order.BeneficiaryCID != "" && req.BeneficiaryIDDocumentCID != "" && order.BeneficiaryCID != req.BeneficiaryIDDocumentCID {
		return errors.New("Beneficiaries in order & order/secret don't match")
	}

	//order & order/secret beneficiary are the same order/secret is not required - discard
	if order.BeneficiaryCID != "" && req.BeneficiaryIDDocumentCID != "" && order.BeneficiaryCID == req.BeneficiaryIDDocumentCID {
		req.BeneficiaryIDDocumentCID = ""
	}
	return nil
}

// PrepareOrderPart1 adds the coin type to the order
func (s *Service) PrepareOrderPart1(order *documents.OrderDoc, reqExtension map[string]string) (fulfillExtension map[string]string, err error) {
	coin, err := strconv.ParseInt(reqExtension["coin"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Parse Coin Type")
	}

	order.Coin = coin
	return nil, nil
}

// PrepareOrderResponse gets the updated order and returns the commitment and extension
func (s *Service) PrepareOrderResponse(orderPart2 *documents.OrderDoc) (commitment string, extension map[string]string, err error) {
	pubKeyPart2of2 := orderPart2.OrderPart2.CommitmentPublicKey
	finalPublicKey, cryptoAddress, err := generateFinalPubKey(s, pubKeyPart2of2, *orderPart2)

	return finalPublicKey, map[string]string{"address": cryptoAddress}, nil
}

// ProduceBeneficiaryEncryptedData -
func (s *Service) ProduceBeneficiaryEncryptedData(blsSK []byte, order *documents.OrderDoc, req *api.OrderSecretRequest) (encrypted []byte, extension map[string]string, err error) {

	enc, err := adhocEncryptedEnvelopeEncode(s, s.NodeID(), *order, blsSK)
	return enc, nil, err
}

// ProduceFinalSecret -
func (s *Service) ProduceFinalSecret(seed, sikeSK []byte, order, orderPart4 *documents.OrderDoc, beneficiaryIDDocumentCID string) (secret, commitment string, extension map[string]string, err error) {
	//retrieve principal IDDoc
	principalDocID, err := common.RetrieveIDDocFromIPFS(s.Ipfs, order.PrincipalCID)
	if err != nil {
		return "", "", nil, err
	}

	finalPrivateKey, err := deriveFinalPrivateKey(s, *orderPart4, sikeSK, seed, beneficiaryIDDocumentCID, s.NodeID(), principalDocID.BLSPublicKey)
	if err != nil {
		return "", "", nil, err
	}
	//Generate Public key & derive crypto address
	finalPublicKey, finalPublicKeyCompressed, err := cryptowallet.PublicKeyFromPrivate(finalPrivateKey)
	if err != nil {
		return "", "", nil, err
	}

	addressForPublicKey, err := addressForPublicKey(finalPublicKey, int(order.Coin))
	if err != nil {
		return "", "", nil, err
	}

	return finalPrivateKey, finalPublicKeyCompressed, map[string]string{"address": addressForPublicKey}, nil
}
