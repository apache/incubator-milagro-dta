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

	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/apache/incubator-milagro-dta/pkg/identity"
	"github.com/pkg/errors"
)

// CreateIdentity creates a new identity
func (s *Service) CreateIdentity(req *api.CreateIdentityRequest) (*api.CreateIdentityResponse, error) {
	idDocumentCID, err := identity.CreateIdentity(req.Name, s.Ipfs, s.Store)
	if err != nil {
		return nil, err
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
