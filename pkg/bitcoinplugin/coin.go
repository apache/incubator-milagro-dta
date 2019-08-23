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

/*
The coin constans are defined here
https://github.com/satoshilabs/slips/blob/master/slip-0044.md
https://www.thepolyglotdeveloper.com/2018/02/generate-cryptocurrency-private-keys-public-addresses-golang/
https://godoc.org/bitbucket.org/dchapes/ripple/crypto/rkey#example-package--Address

*/

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

var (
	errUnsupportedCoin = errors.New("unsupported coin")
)

//Network list of Cryptocurrency networks
type Network struct {
	name        string
	symbol      string
	xpubkey     byte
	xprivatekey byte
}

var network = map[string]Network{
	"btc":     {name: "bitcoin", symbol: "btc", xpubkey: 0x00, xprivatekey: 0x80},
	"testbtc": {name: "bitcoin testnet", symbol: "btc", xpubkey: 0x6f, xprivatekey: 0xef},
}

func addressForPublicKey(publicKey string, coinType int) (string, error) {
	pubK, _ := hex.DecodeString(publicKey)
	pubKeyBytes := []byte(pubK)

	switch coinType {
	case 0: //Bitcoin
		return network["btc"].pubkeyToAddress(pubKeyBytes, false)
	case 1: //Bitcoin Testnet
		return network["testbtc"].pubkeyToAddress(pubKeyBytes, false)
	}
	return "", errUnsupportedCoin
}

func addressForPrivateKey(privateKey []byte, coinType int) (string, error) {
	switch coinType {
	case 0: //Bitcoin
		return network["btc"].privkeyToAddress(privateKey, false)
	case 1: //Bitcoin Testnet
		return network["testbtc"].privkeyToAddress(privateKey, false)
	}
	return "", errUnsupportedCoin
}

func (network Network) privkeyToAddress(privateKey []byte, compressed bool) (string, error) {
	wif, _ := network.createPrivateKeyFromBytes(privateKey)
	address, _ := network.getAddress(wif)
	if compressed == true {
		address.SetFormat(btcutil.PKFCompressed)
	} else {
		address.SetFormat(btcutil.PKFUncompressed)
	}
	return address.EncodeAddress(), nil
}

func (network Network) pubkeyToAddress(publicKey []byte, compressed bool) (string, error) {
	mainNetAddr, err := btcutil.NewAddressPubKey(publicKey, network.getNetworkParams())
	if err != nil {
		return "", errors.Wrap(err, "Failed to decode Public Key")
	}
	if compressed == true {
		mainNetAddr.SetFormat(btcutil.PKFCompressed)
	} else {
		mainNetAddr.SetFormat(btcutil.PKFUncompressed)
	}
	return mainNetAddr.EncodeAddress(), nil
}

func (network Network) getNetworkParams() *chaincfg.Params {
	networkParams := &chaincfg.MainNetParams
	networkParams.PubKeyHashAddrID = network.xpubkey
	networkParams.PrivateKeyID = network.xprivatekey
	return networkParams
}

func (network Network) createPrivateKeyFromBytes(privateKey []byte) (*btcutil.WIF, error) {
	secret, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKey)
	return btcutil.NewWIF(secret, network.getNetworkParams(), true)
}

func (network Network) createPrivateKey() (*btcutil.WIF, error) {
	secret, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}
	return btcutil.NewWIF(secret, network.getNetworkParams(), true)
}

func (network Network) importWIF(wifStr string) (*btcutil.WIF, error) {
	wif, err := btcutil.DecodeWIF(wifStr)
	if err != nil {
		return nil, err
	}
	if !wif.IsForNet(network.getNetworkParams()) {
		return nil, errors.New("The WIF string is not valid for the `" + network.name + "` network")
	}
	return wif, nil
}

func (network Network) getAddress(wif *btcutil.WIF) (*btcutil.AddressPubKey, error) {
	return btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), network.getNetworkParams())
}
