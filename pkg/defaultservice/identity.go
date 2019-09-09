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

	"github.com/apache/incubator-milagro-dta/libs/crypto"
	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/pkg/errors"
)

// CreateIdentity creates a new identity
func (s *Service) CreateIdentity(req *api.CreateIdentityRequest) (*api.CreateIdentityResponse, error) {
	name := req.Name

	//generate crypto random seed
	seed, err := cryptowallet.RandomBytes(48)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Generate random seed")
	}

        //Generate SIKE keys
	rc1, sikePublicKey, sikeSecretKey := crypto.SIKEKeys(seed)
	if rc1 != 0 {
		return nil, errors.New("Failed to generate SIKE keys")
	}

        //Generate BLS keys
	rc1, blsPublicKey, blsSecretKey := crypto.BLSKeys(seed, nil)
	if rc1 != 0 {
		return nil, errors.New("Failed to generate BLS keys")
	}

	ecPubKey, err := common.InitECKeys(seed)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Generate EC Pub Key from random seed")
	}
	//build ID Doc
	idDocument := documents.NewIDDoc()
	idDocument.AuthenticationReference = name
	idDocument.BeneficiaryECPublicKey = ecPubKey
	idDocument.SikePublicKey = sikePublicKey
	idDocument.BLSPublicKey = blsPublicKey
	idDocument.Timestamp = time.Now().Unix()
	//Encode the IDDoc to envelope byte stream
	rawDoc, err := documents.EncodeIDDocument(idDocument, blsSecretKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to encode IDDocument")
	}
	idDocumentCID, err := s.Ipfs.Add(rawDoc)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Save Raw Document into IPFS")
	}
	secrets := common.IdentitySecrets{
		Name:          name,
		Seed:          hex.EncodeToString(seed),
		BLSSecretKey:  hex.EncodeToString(blsSecretKey),
		SikeSecretKey: hex.EncodeToString(sikeSecretKey),
	}

	if err := s.Store.Set("id-doc", idDocumentCID, secrets, map[string]string{"time": time.Now().UTC().Format(time.RFC3339)}); err != nil {
		return nil, errors.Wrap(err, "Failed to Save ID Document - Write to Store")
	}

	return &api.CreateIdentityResponse{
		IDDocumentCID: idDocumentCID,
	}, nil
}

// GetIdentity retrieves an identity
func (s *Service) GetIdentity(req *api.GetIdentityRequest) (*api.GetIdentityResponse, error) {
	idDocumentCID := req.IDDocumentCID
	idDocument, err := common.RetrieveIDDocFromIPFS(s.Ipfs, idDocumentCID)
	if err != nil {
		return nil, err
	}
	return &api.GetIdentityResponse{
		IDDocumentCID:           idDocumentCID,
		AuthenticationReference: idDocument.AuthenticationReference,
		BeneficiaryECPublicKey:  hex.EncodeToString(idDocument.BeneficiaryECPublicKey),
		SikePublicKey:           hex.EncodeToString(idDocument.SikePublicKey),
		BLSPublicKey:            hex.EncodeToString(idDocument.BLSPublicKey),
		Timestamp:               idDocument.Timestamp,
	}, nil
}

// IdentityList reutrns the list of identities
func (s *Service) IdentityList(req *api.IdentityListRequest) (*api.IdentityListResponse, error) {
	page := req.Page
	perPage := req.PerPage
	sortBy := req.SortBy

	IDDocumentCIDes, err := s.Store.ListKeys("id-doc", "time", page*perPage, perPage, sortBy != "dateCreatedAsc")
	if err != nil {
		return nil, err
	}

	fullIDList := make([]api.GetIdentityResponse, len(IDDocumentCIDes))
	for i, idAddress := range IDDocumentCIDes {

		rawDocI, err := s.Ipfs.Get(idAddress)
		if err != nil {
			return nil, errors.Wrapf(err, "Read identity Doc")
		}

		idDocument := &documents.IDDoc{}
		if err = documents.DecodeIDDocument(rawDocI, idAddress, idDocument); err != nil {
			return nil, err
		}
		//Need to copy the whole object so I can append the idddocadderess
		idWithAddress := api.GetIdentityResponse{
			IDDocumentCID:           idAddress,
			AuthenticationReference: idDocument.AuthenticationReference,
			BeneficiaryECPublicKey:  hex.EncodeToString(idDocument.BeneficiaryECPublicKey),
			SikePublicKey:           hex.EncodeToString(idDocument.SikePublicKey),
			BLSPublicKey:            hex.EncodeToString(idDocument.BLSPublicKey),
			Timestamp:               idDocument.Timestamp,
		}

		fullIDList[i] = idWithAddress
	}

	return &api.IdentityListResponse{
		IDDocumentList: fullIDList,
	}, nil
}
