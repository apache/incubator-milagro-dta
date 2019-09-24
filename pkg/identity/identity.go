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
Package identity - manage Identity document and keys
*/
package identity

import (
	"encoding/hex"
	"time"

	"github.com/apache/incubator-milagro-dta/libs/crypto"
	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/libs/ipfs"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/pkg/errors"
)

// CreateIdentity creates a new identity
// returns Identity secrets and Identity document
func CreateIdentity(name string, ipfsConn ipfs.Connector, store *datastore.Store) (idDocumentCID string, err error) {
	//generate crypto random seed
	seed, err := cryptowallet.RandomBytes(48)
	if err != nil {
		err = errors.Wrap(err, "Failed to generate random seed")
		return
	}

	//Generate SIKE keys
	rc1, sikePublicKey, sikeSecretKey := crypto.SIKEKeys(seed)
	if rc1 != 0 {
		err = errors.New("Failed to generate SIKE keys")
		return
	}

	//Generate BLS keys
	rc1, blsPublicKey, blsSecretKey := crypto.BLSKeys(seed, nil)
	if rc1 != 0 {
		err = errors.New("Failed to generate BLS keys")
		return
	}

	ecPubKey, err := common.InitECKeys(seed)
	if err != nil {
		err = errors.Wrap(err, "Failed to generate EC Public Key")
		return
	}

	//build ID Doc
	idDocument := documents.NewIDDoc()
	idDocument.AuthenticationReference = name
	idDocument.BeneficiaryECPublicKey = ecPubKey
	idDocument.SikePublicKey = sikePublicKey
	idDocument.BLSPublicKey = blsPublicKey
	idDocument.Timestamp = time.Now().Unix()

	rawIDDoc, err := documents.EncodeIDDocument(idDocument, blsSecretKey)
	if err != nil {
		err = errors.Wrap(err, "Failed to encode IDDocument")
		return
	}

	idDocumentCID, err = ipfsConn.Add(rawIDDoc)

	secrets := common.IdentitySecrets{
		Name:          name,
		Seed:          hex.EncodeToString(seed),
		BLSSecretKey:  hex.EncodeToString(blsSecretKey),
		SikeSecretKey: hex.EncodeToString(sikeSecretKey),
	}

	if store != nil {
		err = store.Set("id-doc", idDocumentCID, secrets, map[string]string{"time": time.Now().UTC().Format(time.RFC3339)})
		if err != nil {
			err = errors.Wrap(err, "Failed to store ID Document")
			return
		}
	}

	return idDocumentCID, nil
}
