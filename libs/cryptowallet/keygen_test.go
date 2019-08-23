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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RedeemSecret(t *testing.T) {
	secret, err := RedeemSecret("000102030405060708090a0b0c0d0e0f")
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "aef9737d5472303c32318fd6e33aade1fc3ac66159febaa02bd0f1a06834ba61", secret, "HD Wallet Address is incorrect")
}

func Test_RedeemPublicKey(t *testing.T) {
	pubkey, err := RedeemPublicKey("000102030405060708090a0b0c0d0e0f")
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "04ce9b978595558053580d557ff40f9f99a4f1a7609c25268863ee64de7e4abbdad899f36b68872d13348a8098ba6132fb0a0ec4150058a237c7bcaf66e4e16ca7", pubkey, "HD Wallet Address is incorrect")
}

func Test_PublicKeyFromPrivate(t *testing.T) {
	pub, _, err := PublicKeyFromPrivate("C70C9D95F4C1612C53886D2E07A2BAE5AA931F36C65E6AF13BFBA410A0CA1BD0")
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "04db1f3069f0e75feb1bdeeac9e29c8b8eb2ad1bbc9869f72481def876c1f50645b10ee6f9a640bd6dfd8d34286358b493f7b54c37e80ac9d97a4c01a41c0cc8eb", pub, "Pub key not derived from priv key")
}
