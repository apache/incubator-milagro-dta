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
	"time"
)

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
	IDDocumentCID string `json:"idDocumentCID"  validate:"IPFS"`
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
	// BeneficiaryIDDocumentCID string            `json:"BeneficiaryIDDocumentCID,omitempty" validate:"omitempty,IPFS"`
	BeneficiaryIDDocumentCID string            `json:"beneficiaryIDDocumentCID,omitempty"`
	Policy                   Policy            `json:"Policy,omitempty"`
	Extension                map[string]string `json:"extension,omitempty"`
}

//OrderResponse -
type OrderResponse struct {
	// OrderPart1CID  string            `json:"orderPart1CID,omitempty" validate:"omitempty,IPFS"`
	// OrderPart2CID  string            `json:"orderPart2CID,omitempty" validate:"omitempty,IPFS"`
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
	BeneficiaryIDDocumentCID string            `json:"beneficiaryIDDocumentCID,omitempty" validate:"omitempty,IPFS"`
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
	OrderPart3CID     string            `json:"orderPart3CID,omitempty" validate:"IPFS"`
	SenderDocumentCID string            `json:"documentCID,omitempty" validate:"IPFS"`
	Extension         map[string]string `json:"extension,omitempty"`
}

//FulfillOrderSecretResponse -
type FulfillOrderSecretResponse struct {
	OrderPart4CID string            `json:"orderPart4CID,omitempty"`
	Extension     map[string]string `json:"extension,omitempty"`
}

//FulfillOrderRequest -
type FulfillOrderRequest struct {
	OrderPart1CID string            `json:"orderPart1CID,omitempty" validate:"IPFS"`
	DocumentCID   string            `json:"documentCID,omitempty" validate:"IPFS"`
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
}

//Policy - JSON polciy definition
type Policy struct {
	// WalletRef is the customer ref ID for the wallet
	WalletRef string `json:"walletRef"`
	// NodeID as provided by the customer registration
	NodeID string `json:"nodeId"`
	// coinType is the coin type
	// 0 - Bitcoin
	// 1 - Bitcoin test
	CoinType int `json:"coin"`
	// Sharing groups. Each group defines a list of identities.
	// The leader and the node are added automatically.
	SharingGroups []SharingGroup `json:"sharingGroups"`
	// The number of participants
	ParticipantCount uint `json:"participantCount"`
	// The required number of signatures
	Threshold uint `json:"threshold"`
	// Slice that points to signer's keys
	Signers []uint `json:"signers"`
	//Public Address for deposit
	PublicAddress string `json:"publicaddress"`
	//The beneficiary who is able to spend
	BeneficiaryDocID string `json:"beneficiarydocid"`
	//TimeStamps
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SharingGroup defines a list of identities that will hold
// the secret after split
type SharingGroup struct {
	GroupID   int        `json:"groupId"`
	GroupRef  string     `json:"groupref"`
	IDs       []Identity `json:"ids"`
	Threshold int        `json:"threshold"`
	//Temporary store of pSig during withdrawal
	Signature []byte    `json:"signature,omitempty"`
	TimeStamp time.Time `json:"timeStamp"`
	Status    string    `json:"status"`
}

// Identity of a sharing group
type Identity struct {
	// ID is a verifiable identity of IDType
	ID string `json:"id"`
	// IDRef is the identity reference
	IDRef string `json:"idRef"`
	// IDType is the identity type, e.g. "phone"
	IDType string `json:"idType"`
	//Temporary store of vss shares during withdrawal
	Share     []byte    `json:"share,omitempty"`
	Status    string    `json:"status"`
	TimeStamp time.Time `json:"timeStamp"`
}
