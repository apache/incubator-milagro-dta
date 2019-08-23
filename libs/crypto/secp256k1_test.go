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

package crypto

import (
	"fmt"
	"testing"
)

const (
	//For a longer testString see below
	testString     = "The Qu1ck Br0wn F0x Jumps 0v3r Th3 L@zy D0g"
	secpPublicKey  = "041adb36d23c6d01a77d0c2064d491a948922718b848a5f422d17f64b19d65c195986e0ee28049d34e912b3f8022eeb5ec60bcd6562d0c1ee507427e183bdbdd66"
	secpPrivateKey = "8ba005a4ab3b655205435cd61da3630655ccba3c5365e32207eb9bdad561b38f"

	testCypherText = "16751918cf55801daf36e6e6e595a41f3e31d7ba2db55790693e90dfff61ba617fa5ad63fb5fd0c52ccf4b2a85c1f527"
	testT          = "68d7e2c2a4dbceb4e7b885d3"
	testV          = "04a68004e3de100c2b76537e0b3d6eb95ce4f03e4dfac2e01527f73f723b4387c956d1120c6b64a812ccde3658ceeed80cf062e6ea6bf6a95395315e0ef6f2140f"
)

//TestEncrypt - C.ECP_SECP256K1_ECIES_ENCRYPT & DECRYPT
func TestEncrypt(t *testing.T) {
	C, V, T, err := Secp256k1Encrypt(testString, secpPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	plainText, err := Secp256k1Decrypt(C, V, T, secpPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	if plainText != testString {
		t.Fatal(fmt.Sprintf("Failed to decrypt string, expects %s got %s", testString, plainText))
	}
}

func TestDecrypt(t *testing.T) {

	plainText, err := Secp256k1Decrypt(testCypherText, testV, testT, secpPrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	if plainText != testString {
		t.Fatal(fmt.Sprintf("Failed to decrypt string, expects %s got %s", testString, plainText))
	}

}
