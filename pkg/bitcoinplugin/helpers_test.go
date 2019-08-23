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
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/assert"
)

func Test_AddKeys(t *testing.T) {
	privKey1, _ := btcec.NewPrivateKey(btcec.S256())
	privKey2, _ := btcec.NewPrivateKey(btcec.S256())
	privKey1Str := hex.EncodeToString(privKey1.Serialize())
	privKey2Str := hex.EncodeToString(privKey2.Serialize())

	pubKey1 := btcec.PublicKey(privKey1.ToECDSA().PublicKey)
	pubKey2 := btcec.PublicKey(privKey2.ToECDSA().PublicKey)
	pubKey1Str := hex.EncodeToString(pubKey1.SerializeUncompressed())
	pubKey2Str := hex.EncodeToString(pubKey2.SerializeUncompressed())

	//run Additions
	privAdd, _ := addPrivateKeys(privKey1Str, privKey2Str)
	pubAdd, _ := addPublicKeys(pubKey1Str, pubKey2Str)

	privAddBytes, _ := hex.DecodeString(privAdd)
	_, pubKeyFromPrivate := btcec.PrivKeyFromBytes(btcec.S256(), privAddBytes)
	pubKeyFromPrivateString := hex.EncodeToString(pubKeyFromPrivate.SerializeCompressed())
	assert.Equal(t, pubKeyFromPrivateString, pubAdd, "Addition failed")
}
