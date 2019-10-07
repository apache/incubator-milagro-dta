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

package api

/*
 For validation see:
 https://github.com/go-playground/validator

*/

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const (
	// OrderIn Types for the Blockchain

	// TXOrderRequest -
	TXOrderRequest = "v1/order1"
	// TXFulfillRequest -
	TXFulfillRequest = "v1/fulfill/order"
	// TXFulfillResponse -
	TXFulfillResponse = "v1/order2"
	// TXOrderResponse -
	TXOrderResponse = "dump"
	// TXOrderSecretRequest -
	TXOrderSecretRequest = "v1/order/secret1"
	// TXFulfillOrderSecretRequest -
	TXFulfillOrderSecretRequest = "v1/fulfill/order/secret"
	// TXFulfillOrderSecretResponse -
	TXFulfillOrderSecretResponse = "v1/order/secret2"
	// TXOrderSecretResponse -
	TXOrderSecretResponse = "dump"
)

//BlockChainTX - struct for on chain req/resp
type BlockChainTX struct {
	Processor              string
	SenderID               string
	RecipientID            string
	AdditionalRecipientIDs []string
	Payload                []byte
	TXhash                 []byte
	Tags                   map[string]string
}

// CalcHash calculates, sets the TXhash and returns the string representation
func (tx *BlockChainTX) CalcHash() string {
	txSha := sha256.Sum256(tx.Payload)
	txHash := hex.EncodeToString(txSha[:])

	tx.TXhash = txSha[:]
	tx.Tags["txhash"] = txHash

	return txHash

}

//CreateIdentityRequest -
type CreateIdentityRequest struct {
	Name      string            `json:"name,omitempty" validate:"required,alphanum"`
	Extension map[string]string `json:"extension,omitempty"`
}

//CreateIdentityResponse -
type CreateIdentityResponse struct {
	IDDocumentCID string            `json:"idDocumentCID,omitempty"`
	Extension     map[string]string `json:"extension,omitempty"`
}

//GetIdentityRequest -
type GetIdentityRequest struct {
	IDDocumentCID string `json:"idDocumentCID"  validate:"hashID"`
}

//GetIdentityResponse -
type GetIdentityResponse struct {
	IDDocumentCID           string            `json:"idDocumentCID,omitempty"`
	AuthenticationReference string            `json:"authenticationReference,omitempty"`
	BeneficiaryECPublicKey  string            `json:"beneficiaryECPublicKey,omitempty"`
	SikePublicKey           string            `json:"sikePublicKey,omitempty"`
	BLSPublicKey            string            `json:"blsPublicKey,omitempty"`
	Handle                  string            `json:"handle,omitempty"`
	Email                   string            `json:"email,omitempty"`
	Username                string            `json:"string,omitempty"`
	Timestamp               int64             `json:"timestamp,omitempty"`
	Extension               map[string]string `json:"extension,omitempty"`
}

//IdentityListRequest -
type IdentityListRequest struct {
	Page      int               `json:"page,omitempty"`
	PerPage   int               `json:"perPage,omitempty"`
	SortBy    string            `json:"sortBy,omitempty"`
	Extension map[string]string `json:"extension,omitempty"`
}

//IdentityListResponse -
type IdentityListResponse struct {
	IDDocumentList []GetIdentityResponse `json:"idDocumentList,omitempty"`
	Extension      map[string]string     `json:"extension,omitempty"`
}

//OrderRequest -
type OrderRequest struct {
	// BeneficiaryIDDocumentCID string            `json:"BeneficiaryIDDocumentCID,omitempty" validate:"omitempty,hashID"`
	BeneficiaryIDDocumentCID string            `json:"beneficiaryIDDocumentCID,omitempty"`
	Extension                map[string]string `json:"extension,omitempty"`
}

//OrderResponse -
type OrderResponse struct {
	// OrderPart1CID  string            `json:"orderPart1CID,omitempty" validate:"omitempty,hashID"`
	// OrderPart2CID  string            `json:"orderPart2CID,omitempty" validate:"omitempty,hashID"`
	OrderReference string            `json:"orderReference,omitempty" validate:"omitempty"`
	Commitment     string            `json:"commitment,omitempty"`
	CreatedAt      int64             `json:"createdAt,omitempty"`
	Extension      map[string]string `json:"extension,omitempty"`
}

//OrderListRequest -
type OrderListRequest struct {
	Page      int               `json:"page,omitempty"`
	PerPage   int               `json:"perPage,omitempty"`
	SortBy    string            `json:"sortBy,omitempty"`
	Extension map[string]string `json:"extension,omitempty"`
}

//OrderListResponse -
type OrderListResponse struct {
	OrderReference []string          `json:"orderReference,omitempty"`
	Extension      map[string]string `json:"extension,omitempty"`
}

//GetOrderRequest -
type GetOrderRequest struct {
	OrderReference string            `json:"orderReference,omitempty"`
	Extension      map[string]string `json:"extension,omitempty"`
}

//GetOrderResponse -
type GetOrderResponse struct {
	OrderCID  string            `json:"orderCID,omitempty"`
	Order     string            `json:"order,omitempty"`
	Timestamp int64             `json:"timestamp,omitempty"`
	Extension map[string]string `json:"extension,omitempty"`
}

//OrderSecretRequest -
type OrderSecretRequest struct {
	OrderReference           string            `json:"orderReference,omitempty" validate:"omitempty"`
	BeneficiaryIDDocumentCID string            `json:"beneficiaryIDDocumentCID,omitempty" validate:"omitempty,hashID"`
	Extension                map[string]string `json:"extension,omitempty"`
}

//OrderSecretResponse -
type OrderSecretResponse struct {
	Secret         string            `json:"secret,omitempty"`
	Commitment     string            `json:"commitment,omitempty"`
	OrderReference string            `json:"orderReference,omitempty" validate:"omitempty"`
	Extension      map[string]string `json:"extension,omitempty"`
}

//FulfillOrderSecretRequest -
type FulfillOrderSecretRequest struct {
	OrderPart3CID     string            `json:"orderPart3CID,omitempty" validate:"hashID"`
	SenderDocumentCID string            `json:"documentCID,omitempty" validate:"hashID"`
	Extension         map[string]string `json:"extension,omitempty"`
}

//FulfillOrderSecretResponse -
type FulfillOrderSecretResponse struct {
	OrderPart4CID string            `json:"orderPart4CID,omitempty"`
	Extension     map[string]string `json:"extension,omitempty"`
}

//FulfillOrderRequest -
type FulfillOrderRequest struct {
	OrderPart1CID string            `json:"orderPart1CID,omitempty" validate:"hashID"`
	DocumentCID   string            `json:"documentCID,omitempty" validate:"hashID"`
	Extension     map[string]string `json:"extension,omitempty"`
}

//FulfillOrderResponse -
type FulfillOrderResponse struct {
	OrderPart2CID string            `json:"orderPart2CID,omitempty"`
	Extension     map[string]string `json:"extension,omitempty"`
}

//StatusResponse -
type StatusResponse struct {
	Application     string            `json:"application,omitempty"`
	TimeStamp       time.Time         `json:"timeStamp,omitempty"`
	APIVersion      string            `json:"apiVersion,omitempty"`
	NodeCID         string            `json:"nodeCID,omitempty"`
	ExtensionVendor string            `json:"extensionVendor,omitempty"`
	Extension       map[string]string `json:"extension,omitempty"`
	Plugin          string            `json:"plugin,omitempty"`
	NodeType        string            `json:"nodeType,omitempty"`
}
