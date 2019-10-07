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

package bitcoinplugin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/common"
	"github.com/btcsuite/btcd/btcec"
	"github.com/pkg/errors"
)

func deriveFinalPrivateKey(s *Service, order documents.OrderDoc, beneficiariesSikeSK []byte, beneficiariesSeed []byte, beneficiaryIDDocumentCID string, nodeID string, signingBlsPK []byte) (string, error) {

	switch order.BeneficiaryType {
	case documents.OrderDocument_Beneficiary_Unknown_at_Start:
		//we are using the beneficiary specified in order Part 3
		beneficiaryBlob := order.OrderPart3.BeneficiaryEncryptedData

		//Decrypt the Envelope intented for the Beneficiary
		privateKeyPart1of1, err := adhocEncryptedEnvelopeDecode(s, beneficiariesSikeSK, beneficiaryBlob, beneficiaryIDDocumentCID, signingBlsPK)
		if err != nil {
			return "", err
		}
		//Calculate the final private key by Eliptical Key addition of both parts
		privateKeyPart2of2 := order.OrderDocument.OrderPart4.Secret

		finalPrivateKey, err := addPrivateKeys(privateKeyPart1of1, privateKeyPart2of2)

		if err != nil {
			return "", err
		}
		return finalPrivateKey, err

	case documents.OrderDocument_Beneficiary_Known_at_start:
		//we are using the beneficiary specified in the order part 1
		privateKeyPart2of2 := order.OrderDocument.OrderPart4.Secret
		// if order.OrderDocument.BeneficiaryCID != nodeID {
		// 	//need to forward this data to the beneficiary to complete redemption
		// 	return "", errors.New("Currently beneficiary must be the same as the Principal")
		// }
		//restore the Seed
		_, _, ecAddPrivateKey, err := cryptowallet.Bip44Address(beneficiariesSeed, cryptowallet.CoinTypeBitcoinMain, 0, 0, 0)
		if err != nil {
			return "", err
		}
		privateKeyPart1of1 := hex.EncodeToString(ecAddPrivateKey.Serialize())
		finalPrivateKey, err := addPrivateKeys(privateKeyPart1of1, privateKeyPart2of2)
		if err != nil {
			return "", err
		}
		return finalPrivateKey, err
	default:
		return "", errors.New("Critical Error Unknown Beneficiary Type")
	}
}

func adhocEncryptedEnvelopeEncode(s *Service, nodeID string, order documents.OrderDoc, blsSK []byte) ([]byte, error) {
	//Regenerate the original Princaipal Priv Key based on Order
	beneficiaryIDDocumentCID := order.BeneficiaryCID
	seedHex, err := common.RetrieveSeed(s.Store, order.Reference)
	if err != nil {
		return nil, err
	}
	seedOrderModifier := order.OrderDocument.OrderPart2.PreviousOrderCID
	seed, err := hex.DecodeString(seedHex)
	if err != nil {
		return nil, err
	}
	concatenatedSeeds := append(seed, seedOrderModifier...)
	finalSeed := sha256.Sum256(concatenatedSeeds)
	finalSeedHex := hex.EncodeToString(finalSeed[:])
	privateKeyPart1of2, err := cryptowallet.RedeemSecret(finalSeedHex)
	if err != nil {
		return nil, err
	}
	beneficiaryIDDocument, err := common.RetrieveIDDoc(s.Tendermint, beneficiaryIDDocumentCID)
	if err != nil {
		return nil, err
	}
	secretBody := &documents.SimpleString{Content: privateKeyPart1of2}
	header := &documents.Header{}
	recipients := map[string]*documents.IDDoc{
		beneficiaryIDDocumentCID: beneficiaryIDDocument,
	}
	docEnv, err := documents.Encode(nodeID, nil, secretBody, header, blsSK, recipients)
	if err != nil {
		return nil, err
	}
	return docEnv, err
}

func adhocEncryptedEnvelopeDecode(s *Service, sikeSK []byte, beneficiaryBlob []byte, beneficiaryIDDocumentCID string, signingBlsPK []byte) (string, error) {
	//Regenerate the original Principal Priv Key based on Order
	secretBody := &documents.SimpleString{}
	_, err := documents.Decode(beneficiaryBlob, "INTERNAL", sikeSK, beneficiaryIDDocumentCID, nil, secretBody, signingBlsPK)
	if err != nil {
		return "", err
	}
	return secretBody.Content, err
}

func generateFinalPubKey(s *Service, pubKeyPart2of2 string, order documents.OrderDoc) (string, string, error) {
	beneficiaryIDDocumentCID := order.OrderDocument.BeneficiaryCID
	coinType := order.Coin
	var pubKeyPart1of2 string

	if order.BeneficiaryType == documents.OrderDocument_Beneficiary_Unknown_at_Start {
		//There is no beneficiary ID so we do it all locally based on
		//Retrieve the Local Seed
		seedHex, err := common.RetrieveSeed(s.Store, order.Reference)
		if err != nil {
			return "", "", err
		}
		seedOrderModifier := order.OrderDocument.OrderPart2.PreviousOrderCID
		seed, err := hex.DecodeString(seedHex)
		if err != nil {
			return "", "", err
		}
		concatenatedSeeds := append(seed, seedOrderModifier...)
		finalSeed := sha256.Sum256(concatenatedSeeds)
		finalSeedHex := hex.EncodeToString(finalSeed[:])
		//Use HD Wallet to obtain the local Public Key
		pubKeyPart1of2, err = cryptowallet.RedeemPublicKey(finalSeedHex)

		if err != nil {
			return "", "", err
		}
	}

	if order.BeneficiaryType == documents.OrderDocument_Beneficiary_Known_at_start {
		//There is a BeneficiaryID use it to generate the key
		benIDDoc, err := common.RetrieveIDDoc(s.Tendermint, beneficiaryIDDocumentCID)
		if err != nil {
			return "", "", err
		}

		pubKeyPart1of2 = hex.EncodeToString(benIDDoc.BeneficiaryECPublicKey)
	}

	finalPublicKey, err := addPublicKeys(pubKeyPart2of2, pubKeyPart1of2)

	if err != nil {
		return "", "", err
	}
	addressForPublicKey, err := addressForPublicKey(finalPublicKey, int(coinType))
	if err != nil {
		return "", "", err
	}
	return finalPublicKey, addressForPublicKey, nil
}

//AddPrivateKeys Perform eliptical key additon on 2 privates keys
func addPrivateKeys(key1 string, key2 string) (string, error) {
	curveOrder := new(big.Int)
	curveOrder.SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	priv1 := new(big.Int)
	priv2 := new(big.Int)
	priv3 := new(big.Int)
	priv1.SetString(key1, 16)
	priv2.SetString(key2, 16)
	priv3.Add(priv1, priv2)
	if curveOrder.Cmp(priv3) == -1 {
		priv3.Sub(priv3, curveOrder)
	}
	priv4 := fmt.Sprintf("%064x", priv3)
	return priv4, nil
}

//AddPublicKeys Perform eliptical key addition on 2 public keys
func addPublicKeys(key1 string, key2 string) (string, error) {
	pub1Hex, err := hex.DecodeString(key1)
	if err != nil {
		return "", errors.Wrap(err, "Failed to hex decode String")
	}
	dpub1, err := btcec.ParsePubKey(pub1Hex, btcec.S256())
	if err != nil {
		return "", errors.Wrap(err, "Failed to Parse Public Key")
	}
	pub2Hex, err := hex.DecodeString(key2)
	if err != nil {
		return "", errors.Wrap(err, "Failed to hex decode String")
	}
	dpub2, err := btcec.ParsePubKey(pub2Hex, btcec.S256())
	if err != nil {
		return "", errors.Wrap(err, "Failed to Parse Public Key")
	}
	x, y := btcec.S256().Add(dpub1.X, dpub1.Y, dpub2.X, dpub2.Y)
	comp := fmt.Sprintf("04%064X%064X", x, y)
	compByte, err := hex.DecodeString(comp)
	if err != nil {
		return "", err
	}
	pubKey, err := btcec.ParsePubKey([]byte(compByte), btcec.S256())
	if err != nil {
		return "", err
	}
	pubKeyString := hex.EncodeToString(pubKey.SerializeCompressed())
	return pubKeyString, nil
}
