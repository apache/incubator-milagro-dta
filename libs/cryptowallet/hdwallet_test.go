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
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/stretchr/testify/assert"
)

func Test_48ByteSeedGeneration(t *testing.T) {
	seed, err := RandomBytes(65)
	assert.Nil(t, err, "Error should be nil")
	_, _, _, err = Bip44Address(seed, 0, 0, 0, 0)
	assert.NotNil(t, err, "Error should be thrown Seed should be <=64")

	seed, err = RandomBytes(64)
	assert.Nil(t, err, "Error should be nil")
	_, _, _, err = Bip44Address(seed, 0, 0, 0, 0)
	assert.Nil(t, err, "Error should be nil")
}

func Test_OnlyGO(t *testing.T) {
	hkStart := uint32(0x80000000)
	masterSeed, err := hex.DecodeString("3779b041fab425e9c0fd55846b2a03e9a388fb12784067bd8ebdb464c2574a05bcc7a8eb54d7b2a2c8420ff60f630722ea5132d28605dbc996c8ca7d7a8311c0")
	assert.Nil(t, err, "Error should be nil")
	master, err := hdkeychain.NewMaster(masterSeed, &chaincfg.MainNetParams)
	assert.Nil(t, err, "Error should be nil")
	println(master.String())
	ek1, err := master.Child(44 + hkStart)
	assert.Nil(t, err, "Error should be nil")
	ek2, err := ek1.Child(0 + hkStart)
	assert.Nil(t, err, "Error should be nil")
	ek3, err := ek2.Child(0 + hkStart)
	assert.Nil(t, err, "Error should be nil")
	ek4, err := ek3.Child(0)
	assert.Nil(t, err, "Error should be nil")

	//XPUB
	pubKey, err := ek4.Neuter()
	assert.Nil(t, err, "Error should be nil")
	println(pubKey.String())           //XPUB
	pubAddKey0, err := pubKey.Child(0) //0 address from XPUB
	assert.Nil(t, err, "Error should be nil")
	addr1, err := pubAddKey0.Address(&chaincfg.MainNetParams)
	assert.Nil(t, err, "Error should be nil")
	println(addr1.EncodeAddress())

	//XPRIV
	println(ek4.String())             //XPRIV
	priveAddKey0, err := ek4.Child(0) //0 address from Priv
	assert.Nil(t, err, "Error should be nil")
	addr2, err := priveAddKey0.Address(&chaincfg.MainNetParams)
	assert.Nil(t, err, "Error should be nil")
	println(addr2.EncodeAddress())
	priv, err := priveAddKey0.ECPrivKey()
	assert.Nil(t, err, "Error should be nil")
	wif, err := btcutil.NewWIF(priv, &chaincfg.MainNetParams, true)
	assert.Nil(t, err, "Error should be nil")
	println(wif.String())
}

//Test_BTCVectors - Cycle through the list of test vectors below taken from https://www.coinomi.com/recovery-phrase-tool.html
//Need to add significantly more vectors for more edge cases, coins, accounts, change etc.
func Test_BTCVectors(t *testing.T) {

	for _, testVector := range vectors {

		Chain := &testVector.net
		startingEntropy, err := hex.DecodeString(testVector.entropy)
		assert.Nil(t, err, "Error should be nil")

		mnemonic, err := entropy2Mnemonic(startingEntropy)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, mnemonic, testVector.mnemonic, "Mnemonic is incorrect")

		seed, err := mnemonic2Seed(mnemonic)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, hex.EncodeToString(seed), testVector.seed, "Seed from Mnemonic is incorrect")

		seed2, err := seedFromEntropy(testVector.entropy)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, hex.EncodeToString(seed2), testVector.seed, "Seed from Mnemonic is incorrect")

		xpriv, err := masterKeyFromSeed(seed, testVector.coin)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, xpriv.String(), testVector.bip32Root, "Invalid xPriv")
		print("\n" + xpriv.String())

		bip32Extended, err := bip32Extended(seed, testVector.coin, testVector.account, testVector.change)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, bip32Extended.String(), testVector.bip32ExPriv, "Invalid BIP32 Priv")
		print("\n" + bip32Extended.String())

		xPubKey, err := bip32Extended.Neuter()
		assert.Nil(t, err, "Error should be nil")
		xPub := xPubKey.String()
		assert.Equal(t, xPub, testVector.bip32ExPub, "Invalid BIP32 Public (1)")

		btcAdd, _, btcPrivKey, err := Bip44Address(seed, testVector.coin, testVector.account, testVector.change, testVector.addressIndex)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, btcAdd, testVector.address, "Invalid BTC Address")

		wifComp, err := btcutil.NewWIF(btcPrivKey, Chain, true)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, wifComp.String(), testVector.privKey, "Invalid BTC Address")
	}
}

var vectors = []struct {
	net          chaincfg.Params
	coin         int
	account      int
	change       int
	addressIndex int64
	entropy      string
	mnemonic     string
	seed         string
	bip32Root    string
	bip32ExPriv  string
	bip32ExPub   string
	address      string
	privKey      string
}{
	{
		chaincfg.MainNetParams,
		0, 0, 0, 0,
		"000102030405060708090a0b0c0d0e0f",
		"abandon amount liar amount expire adjust cage candy arch gather drum buyer",
		"3779b041fab425e9c0fd55846b2a03e9a388fb12784067bd8ebdb464c2574a05bcc7a8eb54d7b2a2c8420ff60f630722ea5132d28605dbc996c8ca7d7a8311c0",
		"xprv9s21ZrQH143K2XojduRLQnU8D8K59KSBoMuQKGx8dW3NBitFDMkYGiJPwZdanjZonM7eXvcEbxwuGf3RdkCyyXjsbHSkwtLnJcsZ9US42Gd",
		"xprvA2QWrMvVn11Cnc8Wv5XH22Phaz1eLLYUtUVCJxjRu3eSbPZk3WphdkqGBnAKiKtg3bxkL48zbf9C8jJKtbDhB4kTJuNfv3KZVRjxseHNNWk",
		"xpub6FPsFsTPcNZW16Cz274HPALS91r8joGLFhQo7M93TPBRUBttb48xBZ9k34oiG29Bvqfry9QyXPsGXSRE1kjut92Dgik1w6Whm1GU4F122n8",
		"128BCBZndgrPXzEgF4QbVR3jnQGwzRtEz5",
		"L35qaFLpbCc9yCzeTuWJg4qWnTs9BaLr5CDYcnJ5UnGmgLo8JBgk",
	},
	{
		chaincfg.MainNetParams,
		0, 1, 1, 19,
		"000102030405060708090a0b0c0d0e0f",
		"abandon amount liar amount expire adjust cage candy arch gather drum buyer",
		"3779b041fab425e9c0fd55846b2a03e9a388fb12784067bd8ebdb464c2574a05bcc7a8eb54d7b2a2c8420ff60f630722ea5132d28605dbc996c8ca7d7a8311c0",
		"xprv9s21ZrQH143K2XojduRLQnU8D8K59KSBoMuQKGx8dW3NBitFDMkYGiJPwZdanjZonM7eXvcEbxwuGf3RdkCyyXjsbHSkwtLnJcsZ9US42Gd",
		"xprvA1WdMXWAt7uoSc1STdHwXLLgWaiJovFwYAyHtowP666fCAZNA5T3msYkKqfiYFYwRVQqr8SbuYgtf2tZ1PJwWhWNHMjdknEnGkDmZrkpFn4",
		"xpub6EVym334iVU6f65uZepwtUHR4cYoDNynuPtthCLzeRde4xtWhcmJKfsEB6vXKpYyh8FwwXoWE8NUDNvvsNAxDdXr3bwWhLfekyd3qcbbtuj",
		"1N2pq8QzgYRmSyszkaGgNngpcUTqejsmKn",
		"L3U3y6CFHwq9dS15VvtgHGorZFicanWBkfbee5gDp3m7uFHVoqAy",
	},
	{
		chaincfg.TestNet3Params,
		1, 0, 0, 0,
		"000102030405060708090a0b0c0d0e0f",
		"abandon amount liar amount expire adjust cage candy arch gather drum buyer",
		"3779b041fab425e9c0fd55846b2a03e9a388fb12784067bd8ebdb464c2574a05bcc7a8eb54d7b2a2c8420ff60f630722ea5132d28605dbc996c8ca7d7a8311c0",
		"tprv8ZgxMBicQKsPdM3GJUGqaS67XFjHNqUC8upXBhNb7UXqyKdLCj6HnTfqrjoEo6x89neRY2DzmKXhjWbAkxYvnb1U7vf4cF4qDicyb7Y2mNa",
		"tprv8hJrzKEmbFfBx44tsRe1wHh25i5QGztsawJGmxeqryPwdXdKrgxMgJUWn35dY2nrYmomRWWL7Y9wJrA6EvKJ27BfQTX1tWzZVxAXrR2pLLn",
		"tpubDDzu8jH1jdLrqX6gm5JcLhM8ejbLSL5nAEu44Uh9HFCLU1t6V5mwro6NxAXCfR2jUJ9vkYkUazKXQSU7WAaA9cbEkxdWmbLxHQnWqLyQ6uR",
		"mq1VMMXiZKLdY2WLeaqocJxXijhEFoQu3X",
		"cMwbkii126fSsPtWBUuUPrKZS5KK3qCjSNuRhcuw6sJ8HmVsrmHq",
	},
	{
		chaincfg.TestNet3Params,
		1, 1, 1, 19,
		"000102030405060708090a0b0c0d0e0f",
		"abandon amount liar amount expire adjust cage candy arch gather drum buyer",
		"3779b041fab425e9c0fd55846b2a03e9a388fb12784067bd8ebdb464c2574a05bcc7a8eb54d7b2a2c8420ff60f630722ea5132d28605dbc996c8ca7d7a8311c0",
		"tprv8ZgxMBicQKsPdM3GJUGqaS67XFjHNqUC8upXBhNb7UXqyKdLCj6HnTfqrjoEo6x89neRY2DzmKXhjWbAkxYvnb1U7vf4cF4qDicyb7Y2mNa",
		"tprv8hfdp3VTvAW219XTuczuzC97BuNCpWZfRjPLEXgvrxfv8YW4yXFKAcic5S1AtPNWe36sXkUNWdAh6PSuQsYk3nabNKQnaPhGTx5etiDFSsW",
		"tpubDEMfxTXi4YBgtcZFoGfWPboDkvt8yqka12z7X3jEHEUJy2kqbv4uM7LUFYcaMEApY7TbKj9FVqAhwUcvXVbLHRyyNRsUEFJy7x46dXMUHb2",
		"mwXHRw1hMmWQhn54r3mXXZ62kaozQTkeTa",
		"cPHvnD9R6ELjEWKamse2szPva2r8rJsxUHcKngfZghSTzaeaqa2g",
	},
}
