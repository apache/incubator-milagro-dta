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
Package cryptowallet - generates and manipulates SECP256 keys
*/
package cryptowallet

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/pkg/errors"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	// CoinTypeBitcoinMain is Bitcoin main network coin
	CoinTypeBitcoinMain = 0
	// CoinTypeBitcoinTestNet is Bitcoin test network coin
	CoinTypeBitcoinTestNet = 1
)

var (
	// ErrUnsupportedCoin is returned when the coinType is not supported
	errUnsupportedCoin = errors.New("unsupported coin")
)

//BIP 32 - xPub/xPriv from seed
//BIP 39 - Mnemonic Wordlist
//BIP 44 - m / purpose' / coin_type' / account' / change / address_index

/*
   Some Useful tools for testsing
   Mnemonic Code Converter tool
   https://iancoleman.io/bip39/

   Key Convertor
   https://www.bitaddress.org
*/

// userWallet holds the persistent user data

//chainParams chain config based on coin type
func chainParams(coinType int) (*chaincfg.Params, error) {
	switch coinType {
	case CoinTypeBitcoinMain:
		return &chaincfg.MainNetParams, nil
	case CoinTypeBitcoinTestNet:
		return &chaincfg.TestNet3Params, nil
	}
	return nil, errUnsupportedCoin
}

// Bip44Address -  generates a bitcoin address & private key for a given BIP44 path  -return btc address, private key, error
func Bip44Address(seed []byte, coinType int, account int, change int, addressIndex int64) (string, *btcec.PublicKey, *btcec.PrivateKey, error) {
	chain, err := chainParams(coinType)
	if err != nil {
		return "", nil, nil, err
	}
	bip32extended, err := bip32Extended(seed, coinType, account, change)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "Failed to BIP32 Ext from seed,coin,acc,change")
	}
	priveAddKey, err := bip32extended.Child(uint32(addressIndex))
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "Failed to derive child from Bip32 ext")
	}
	priv, err := priveAddKey.ECPrivKey()
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "Failed to extract Private key from ECPriv")
	}
	pubKey, err := priveAddKey.ECPubKey()
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "Failed to extract Public key from ECPriv")
	}
	pub, err := priveAddKey.Address(chain)
	if err != nil {
		return "", nil, nil, errors.Wrap(err, "Failed to extract public key from priv key")
	}
	pubAdd := pub.EncodeAddress()
	return pubAdd, pubKey, priv, nil
}

// MasterKeyFromSeed Generates a MasterKey or XPriv from a seed
func masterKeyFromSeed(seed []byte, coinType int) (*hdkeychain.ExtendedKey, error) {
	chain, err := chainParams(coinType)
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, chain)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive masterkey from seed")
	}
	return masterKey, nil
}

//Bip32Extended get Bip32 extended Keys for path
func bip32Extended(seed []byte, coinType int, account int, change int) (*hdkeychain.ExtendedKey, error) {
	masterKey, err := masterKeyFromSeed(seed, coinType)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive masterkey from seed")
	}
	hkStart := uint32(0x80000000)
	child1, err := masterKey.Child(44 + hkStart)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive child1 from seed")
	}
	child2, err := child1.Child(uint32(coinType) + hkStart)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive child2 from seed")
	}
	child3, err := child2.Child(uint32(account) + hkStart)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive child3 from seed")
	}
	bip32extended, err := child3.Child(uint32(change))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive BIP32ext from child3")
	}
	return bip32extended, nil
}

//SeedFromEntropy generate a Seed from supplied entropy string (from HSM)
func seedFromEntropy(entropy string) ([]byte, error) {
	entropyDecoded, err := hex.DecodeString(entropy)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to decode entropy as hex")
	}
	mnemonic, err := entropy2Mnemonic(entropyDecoded)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive mnemonic from Entropy")
	}
	seed, err := mnemonic2Seed(mnemonic)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive seed from mnemonic")
	}
	return seed, nil
}

//Entropy2Mnemonic convert a seed (entropy) to a Mmemonic String
func entropy2Mnemonic(entropy []byte) (string, error) {
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", errors.Wrap(err, "Failed to derive mnemonic from Entropy")
	}
	return mnemonic, nil
}

//Mnemonic2Seed convert BIP39 mnemonic (recovery phrase) to a seed byte array
func mnemonic2Seed(mnemonic string) ([]byte, error) {
	password := "" //we arent using passwords
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to derive seed from Mnemonic")
	}
	return seed, nil
}
