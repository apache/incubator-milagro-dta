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

package cryptowallet

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/pkg/errors"
)

var (
	errEntropyError = errors.New("Failed to supply entropy")
)

//RedeemSecret - using supplied seed return the 1st entry in the HD walet
func RedeemSecret(entropy string) (secret string, err error) {
	if entropy == "" {
		return "", errEntropyError
	}
	seed, err := seedFromEntropy(entropy)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create seed from entropy")
	}
	_, _, privKey, err := Bip44Address(seed, 0, 0, 0, 0)
	if err != nil {
		return "", errors.Wrap(err, "Failed to derive Wallet Key 0")
	}
	privateKeyStr := hex.EncodeToString(privKey.Serialize())
	return privateKeyStr, nil
}

//RedeemPublicKey - using supplied seed return the 1st entry in the HD walet
func RedeemPublicKey(entropy string) (secret string, err error) {
	if entropy == "" {
		return "", errEntropyError
	}
	seed, err := seedFromEntropy(entropy)
	if err != nil {
		return "", errors.Wrap(err, "Failed to create seed from entropy")
	}
	_, pubKey, _, err := Bip44Address(seed, 0, 0, 0, 0)
	if err != nil {
		return "", errors.Wrap(err, "Failed to derive Wallet Key 0")
	}
	pubKeyStr := hex.EncodeToString(pubKey.SerializeUncompressed())
	return pubKeyStr, nil
}

//PublicKeyFromPrivate Derive EC Public key from a private key
func PublicKeyFromPrivate(priv string) (string, string, error) {
	remotePrivKeyBytes, err := hex.DecodeString(priv)
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to hex decode PrivateKey")
	}

	_, cpub1 := btcec.PrivKeyFromBytes(btcec.S256(), remotePrivKeyBytes)
	remotePubKeyStr := hex.EncodeToString(cpub1.SerializeUncompressed())
	remotePubKeyCompressedStr := hex.EncodeToString(cpub1.SerializeCompressed())
	return remotePubKeyStr, remotePubKeyCompressedStr, nil
}

//RandomBytes - generate n random bytes
func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}
