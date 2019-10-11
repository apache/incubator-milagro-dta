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
	"encoding/hex"
	"log"
	"testing"
	"time"
)

//Results on I9 laptop
// BLSKeys took 1.367643ms
// AESCBCEncrypt took 1.831µs
// EncapsulateEncrypt took 163.514139ms
// DecapsulateDecrypt took 165.586105ms
// AESCBCDecrypt took 1.791µs
// BLSSign took 652.917µs
// BLSVerify took 3.441072ms

//Test_Bench1 Some simple benchmarks to show speeds of the different crypto functions
func Test_Bench1(t *testing.T) {
	start := time.Now()
	elapsed := time.Since(start)

	var SIKEpk []byte
	var SIKEsk []byte
	var BLSsk []byte
	var BLSpk []byte
	var C2 []byte
	var C1 []byte
	var EK []byte
	var P3 []byte
	var S []byte

	SEEDHex := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f30"
	SEED, _ := hex.DecodeString(SEEDHex)
	// AES-256 Key
	KHex := "af5f8452d644d131c35164fee8c8300fb29725b03b00eaef411c823293c469d8"

	// AES IVs
	IV1Hex := "de4724534df5f50160a28cf3b3caec80"
	IV2Hex := "a08576336500e79b2593dc4c10b0f36c"

	// Messsage to encrypt and sign. Note it is zero padded.
	PHex := "48656c6c6f20426f622120546869732069732061206d6573736167652066726f6d20416c696365000000000000000000"

	// Generate SIKE keys

	start = time.Now()
	for i := 0; i < 100; i++ {
		_, SIKEpk, SIKEsk = SIKEKeys(SEED)
	}
	elapsed = time.Since(start)
	log.Printf("SIKEKeys took %s", elapsed/100)

	// Generate BLS keys

	start = time.Now()
	for i := 0; i < 1000; i++ {
		_, BLSpk, BLSsk = BLSKeys(SEED, nil)
	}
	elapsed = time.Since(start)
	log.Printf("BLSKeys took %s", elapsed/1000)

	// Encrypt message
	K, _ := hex.DecodeString(KHex)
	IV1, _ := hex.DecodeString(IV1Hex)
	P1, _ := hex.DecodeString(PHex)

	start = time.Now()
	for i := 0; i < 1000; i++ {
		C1 = AESCBCEncrypt(K, IV1, P1)
	}
	elapsed = time.Since(start)
	log.Printf("AESCBCEncrypt took %s", elapsed/1000)

	// Encrypt AES Key, K, and returned encapsulated key used for
	// encryption
	IV2, _ := hex.DecodeString(IV2Hex)
	P2 := K

	start = time.Now()
	for i := 0; i < 100; i++ {
		_, C2, EK = EncapsulateEncrypt(P2, IV2, SIKEpk)
	}
	elapsed = time.Since(start)
	log.Printf("EncapsulateEncrypt took %s", elapsed/100)

	// Decapsulate the AES Key and use it to decrypt the cipherText.
	// P2 and P3 should be the same. This value is the AES-256 key
	// used to encrypt the plaintext P1

	start = time.Now()
	for i := 0; i < 100; i++ {
		_, P3 = DecapsulateDecrypt(C2, IV2, SIKEsk, EK)
	}
	elapsed = time.Since(start)
	log.Printf("DecapsulateDecrypt took %s", elapsed/100)

	K2 := P3

	start = time.Now()
	for i := 0; i < 1000; i++ {
		_ = AESCBCDecrypt(K2, IV1, C1)
	}
	elapsed = time.Since(start)
	log.Printf("AESCBCDecrypt took %s", elapsed/1000)

	// BLS Sign a message
	start = time.Now()
	for i := 0; i < 1000; i++ {
		_, S = BLSSign(P1, BLSsk)
	}
	elapsed = time.Since(start)
	log.Printf("BLSSign took %s", elapsed/1000)

	// BLS Verify signature
	start = time.Now()
	for i := 0; i < 1000; i++ {
		_ = BLSVerify(P1, BLSpk, S)
	}
	elapsed = time.Since(start)
	log.Printf("BLSVerify took %s", elapsed/1000)

}
