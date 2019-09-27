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

package identity

import (
	"fmt"

	"github.com/apache/incubator-milagro-dta/libs/crypto"
	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/pkg/errors"
)

// Secrets - keys required for decryption and signing
type Secrets struct {
	Seed          []byte
	SikeSecretKey []byte
	BLSSecretKey  []byte
}

// GenerateBLSKeys generate BLS keys from seed
func GenerateBLSKeys(seed []byte) (blsPublic, blsSecret []byte, err error) {
	rc1, blsPublic, blsSecret := crypto.BLSKeys(seed, nil)
	if rc1 != 0 {
		err = fmt.Errorf("Failed to generate BLS keys: %v", rc1)
	}
	return
}

// GenerateSIKEKeys generate SIKE keys from seed
func GenerateSIKEKeys(seed []byte) (sikePublic, sikeSecret []byte, err error) {
	rc1, sikePublic, sikeSecret := crypto.SIKEKeys(seed)
	if rc1 != 0 {
		err = fmt.Errorf("Failed to generate SIKE keys: %v", rc1)
	}
	return
}

// GenerateECPublicKey - generate EC keys using BIP44 HD Wallets (as bitcoin) from seed
func GenerateECPublicKey(seed []byte) (ecPublic []byte, err error) {
	//EC ADD Keypair Protocol
	_, pubKeyECADD, _, err := cryptowallet.Bip44Address(seed, cryptowallet.CoinTypeBitcoinMain, 0, 0, 0)
	if err != nil {
		err = errors.Wrap(err, "Failed to derive EC HD Wallet Key")
		return
	}

	return pubKeyECADD.SerializeCompressed(), nil
}
