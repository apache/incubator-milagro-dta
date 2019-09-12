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

	"github.com/stretchr/testify/assert"
)

func Test_AddressForPublicKey(t *testing.T) {
	pubKey := "0487DBF8D88A860270AB7D689EB44C2DFFF768D2F7851A753FACF356978B82CE4ACB5C9B061FC884668D9BB46B83D6BF180A2099F397142785D2E03DACCEF03D01"
	btcAddress, err := addressForPublicKey(pubKey, 0)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "1MwiNcg3v19BLeawNJKL8L18m4Tzmtua5T", btcAddress)

	btcTestNetAddress, err := addressForPublicKey(pubKey, 1)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "n2Tfffm2j2aS7m4Z5sHhxFDTd44hfnDhPz", btcTestNetAddress)

	_, err = addressForPublicKey(pubKey, 9999999999)
	assert.EqualError(t, err, "unsupported coin")
}
func Test_AddressForPrivateKey(t *testing.T) {
	privKey, _ := hex.DecodeString("EB354D4B18E0B4AC6E63369F33D0CFFE7F3C09101D29678877A6CE8879D7E152")
	btcAddress, err := addressForPrivateKey(privKey, 0)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "1MwiNcg3v19BLeawNJKL8L18m4Tzmtua5T", btcAddress)

	btcTestNetAddress, err := addressForPrivateKey(privKey, 1)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "n2Tfffm2j2aS7m4Z5sHhxFDTd44hfnDhPz", btcTestNetAddress)

	_, err = addressForPrivateKey(privKey, 9999999999)
	assert.EqualError(t, err, "unsupported coin")
}
