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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Mine2(t *testing.T) {

	// Generate keys

	sikePK, _ := hex.DecodeString("f013ed197f8d2756eb7fb237bfe6714ec0d3f12ddc7431b8881f4bece686cafa67e9cfcfc75d5117114e857b379b1fe6b0e67f11dc7936d4ed4e8c3488b2300a80b12e3bee303b2d1e5d24c2c9b3e5855aeddfa45508e587f71a95c23d643523bf8385a7d830254b6ac7f231cf945795311ad8915da9605ff5e6543a01b923e924b37d7d2e464d7a9d78cc5ad4b456fc735ff7d60bd9225558bbad053a3d724a82d3a77348d4aac4296c6f79e0db21a49695e158abf16724e635b930e2704620eccdc1aab49d0655f7ca1eb36ba8c99129b2de8bee6e843a849b802b90775090a7617dc6e4a9ef289c0eac292cc00a80f9a76df2471695e9c433337c649acacc4ba0a2a5f80ccb594553861161cb4771bceabc4cb430bf233307e54f3ebbc0291e75db9cac6d9760e39cd39f373d8e852cac56585bbdf0167c62f728e2889d8387abc0b25f40f990dc7efd93d17839f04ad1e8f3dbe6b518f4f629d69ac7ce04ea614b977a60deb2148fd8e32dfe690da6bfd2aa614dd60688a86989e80cb718f69fa0d02d1ee9796510a9cb523e03fa3ad2a083fb42d11df0e5584873a80b86725d806fbaee92bd35ba2548a79a1a477ea1b1c5d4e33434db2b25ffd6713c887d25e774a0b7c6ca039fcf9e527739449a918d9c884dadf6ac1f11d9cc235a44ff9f69d568c7e5e1999d0353f37710ce71d5f59f756c9a89cd7fe318ebecb1ba7e408171da3e516a1e0e6bef2379e6b519986d6e75e720be3a8892ade537926760a60b011ef627db482a7edfce13fcc767f1eb1c")
	sikeSK, _ := hex.DecodeString("5a8ac3a8704ac0b906f905cbfbf62acf81046b9f7f24e7f95b43a62c37d483c7110fd01e31fe6961421b50a0672a0450153894a9a221bb62d059d7e5589e4765e93c6d5bc713fc822f44276b2eab4700f013ed197f8d2756eb7fb237bfe6714ec0d3f12ddc7431b8881f4bece686cafa67e9cfcfc75d5117114e857b379b1fe6b0e67f11dc7936d4ed4e8c3488b2300a80b12e3bee303b2d1e5d24c2c9b3e5855aeddfa45508e587f71a95c23d643523bf8385a7d830254b6ac7f231cf945795311ad8915da9605ff5e6543a01b923e924b37d7d2e464d7a9d78cc5ad4b456fc735ff7d60bd9225558bbad053a3d724a82d3a77348d4aac4296c6f79e0db21a49695e158abf16724e635b930e2704620eccdc1aab49d0655f7ca1eb36ba8c99129b2de8bee6e843a849b802b90775090a7617dc6e4a9ef289c0eac292cc00a80f9a76df2471695e9c433337c649acacc4ba0a2a5f80ccb594553861161cb4771bceabc4cb430bf233307e54f3ebbc0291e75db9cac6d9760e39cd39f373d8e852cac56585bbdf0167c62f728e2889d8387abc0b25f40f990dc7efd93d17839f04ad1e8f3dbe6b518f4f629d69ac7ce04ea614b977a60deb2148fd8e32dfe690da6bfd2aa614dd60688a86989e80cb718f69fa0d02d1ee9796510a9cb523e03fa3ad2a083fb42d11df0e5584873a80b86725d806fbaee92bd35ba2548a79a1a477ea1b1c5d4e33434db2b25ffd6713c887d25e774a0b7c6ca039fcf9e527739449a918d9c884dadf6ac1f11d9cc235a44ff9f69d568c7e5e1999d0353f37710ce71d5f59f756c9a89cd7fe318ebecb1ba7e408171da3e516a1e0e6bef2379e6b519986d6e75e720be3a8892ade537926760a60b011ef627db482a7edfce13fcc767f1eb1c")
	secret, _ := hex.DecodeString("5a8ac3a8704ac0b906f905cbfbf62acf81046b9f7f24e7f95b43a62c37d483c7110fd01e31fe6961421")
	final, _ := hex.DecodeString("5a8ac3a8704ac0b906f905cbfbf62acf81046b9f7f24e7f95b43a62c37d483c7110fd01e31fe6961421")
	iv, _ := hex.DecodeString("9640061f9e3c29fdd52945feb678de83")

	_, cipherText, encapsulatedKey := EncapsulateEncrypt(secret, iv, sikePK)

	println(encapsulatedKey)
	println(cipherText)
	encapsulatedKeyHex := hex.EncodeToString(encapsulatedKey)
	cipherTextHex := hex.EncodeToString(cipherText)
	fmt.Printf("encapsulatedKey : %s \n", encapsulatedKeyHex)
	fmt.Printf("BLSscipherTextk : %s \n", cipherTextHex)

	_, aes := DecapsulateDecrypt(cipherText, iv, sikeSK, encapsulatedKey)

	print(hex.EncodeToString(aes))
	assert.Equal(t, aes, final, "Secret doesnt match")
}

func Test_Mine(t *testing.T) {
	SEEDHex := "2c82c5a6b14f6ce0fca9b83e929f6ca091fb25b6648676b3c387e8a13b0f4cee92a54d42e388db3fbb0e906b32e880f4"
	SEED, _ := hex.DecodeString(SEEDHex)

	// Generate keys
	RC1, SIKEpk, SIKEsk, BLSpk, BLSsk := Keys(SEED)
	if RC1 != 0 {
		fmt.Println("Panicking!")
		panic("Failed to create keys")
	}

	SIKEpkHex := hex.EncodeToString(SIKEpk)
	SIKEskHex := hex.EncodeToString(SIKEsk)
	fmt.Printf("SIKEpk : %s \n", SIKEpkHex)
	fmt.Printf("SIKEsk : %s \n", SIKEskHex)

	BLSpkHex := hex.EncodeToString(BLSpk)
	BLSskHex := hex.EncodeToString(BLSsk)
	fmt.Printf("BLSpk : %s \n", BLSpkHex)
	fmt.Printf("BLSsk : %s \n", BLSskHex)

	secret, _ := hex.DecodeString("79e54957d823668872b41f6bd6394a2132935902bc9c9192474562d58225e129")
	iv, _ := hex.DecodeString("f7427ca6749f696c0dab97582b96222f")

	_, cipherText, encapsulatedKey := EncapsulateEncrypt(secret, iv, SIKEpk)

	println(encapsulatedKey)
	println(cipherText)
	encapsulatedKeyHex := hex.EncodeToString(encapsulatedKey)
	cipherTextHex := hex.EncodeToString(cipherText)
	fmt.Printf("encapsulatedKey : %s \n", encapsulatedKeyHex)
	fmt.Printf("BLSscipherTextk : %s \n", cipherTextHex)

	_, aes := DecapsulateDecrypt(cipherText, iv, SIKEsk, encapsulatedKey)

	print(hex.EncodeToString(aes))
	assert.Equal(t, aes, secret, "Secret doesnt match")
}

func Test_Smoke_Test(t *testing.T) {
	SEEDHex := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f30"
	SEED, _ := hex.DecodeString(SEEDHex)

	// AES-256 Key
	KHex := "af5f8452d644d131c35164fee8c8300fb29725b03b00eaef411c823293c469d8"

	// AES IVs
	IV1Hex := "de4724534df5f50160a28cf3b3caec80"
	IV2Hex := "a08576336500e79b2593dc4c10b0f36c"

	// Messsage to encrypt and sign. Note it is zero padded.
	PHex := "48656c6c6f20426f622120546869732069732061206d6573736167652066726f6d20416c696365000000000000000000"

	// Generate keys
	RC1, SIKEpk, SIKEsk, BLSpk, BLSsk := Keys(SEED)
	if RC1 != 0 {
		fmt.Println("Panicking!")
		panic("Failed to create keys")
	}

	SIKEpkHex := hex.EncodeToString(SIKEpk)
	SIKEskHex := hex.EncodeToString(SIKEsk)
	fmt.Printf("SIKEpk : %s \n", SIKEpkHex)
	fmt.Printf("SIKEsk : %s \n", SIKEskHex)

	BLSpkHex := hex.EncodeToString(BLSpk)
	BLSskHex := hex.EncodeToString(BLSsk)
	fmt.Printf("BLSpk : %s \n", BLSpkHex)
	fmt.Printf("BLSsk : %s \n", BLSskHex)

	// Encrypt message
	K, _ := hex.DecodeString(KHex)
	IV1, _ := hex.DecodeString(IV1Hex)
	P1, _ := hex.DecodeString(PHex)
	fmt.Printf("P1 : %s \n", P1)
	C1 := AESCBCEncrypt(K, IV1, P1)
	C1Hex := hex.EncodeToString(C1)
	fmt.Printf("C1Hex : %s \n", C1Hex)

	// Encrypt AES Key, K, and returned encapsulated key used for
	// encryption
	IV2, _ := hex.DecodeString(IV2Hex)
	P2 := K
	P2Hex := hex.EncodeToString(P2)
	fmt.Printf("P2Hex : %s \n", P2Hex)
	RC2, C2, EK := EncapsulateEncrypt(P2, IV2, SIKEpk)
	if RC2 != 0 {
		fmt.Println("Panicking!")
		panic("Failed to encrypt and encapsulate key")
	}
	C2Hex := hex.EncodeToString(C2)
	fmt.Printf("C2Hex : %s \n", C2Hex)
	EKHex := hex.EncodeToString(EK)
	fmt.Printf("EKHex : %s \n", EKHex)

	// Decapsulate the AES Key and use it to decrypt the ciphertext.
	// P2 and P3 should be the same. This value is the AES-256 key
	// used to encrypt the plaintext P1
	RC3, P3 := DecapsulateDecrypt(C2, IV2, SIKEsk, EK)
	if RC3 != 0 {
		fmt.Println("Panicking!")
		panic("Failed to decapsulate key and decrypt ciphertext")
	}
	P3Hex := hex.EncodeToString(P3)
	fmt.Printf("P3Hex : %s \n", P3Hex)

	// Decrypt the ciphertext to recover the orignal plaintext
	// contained in P1
	K2 := P3
	P4 := AESCBCDecrypt(K2, IV1, C1)
	fmt.Printf("P4 : %s \n", P4)

	// Sign a message
	RC5, S := Sign(P1, BLSsk)
	if RC5 != 0 {
		fmt.Println("Panicking!")
		panic("Failed to sign message")
	}
	SHex := hex.EncodeToString(S)
	fmt.Printf("S : %s \n", SHex)

	// Verify signature
	RC6 := Verify(P1, BLSpk, S)
	if RC6 != 0 {
		fmt.Println("Panicking!")
		panic("Failed to verify signature")
	} else {
		fmt.Println("Signature Verified")
	}
}

func TestPQNIST_AES_CBC_ENCRYPT(t *testing.T) {
	KHex := "6ed76d2d97c69fd1339589523931f2a6cff554b15f738f21ec72dd97a7330907"
	IVHex := "851e8764776e6796aab722dbb644ace8"
	PHex := "6282b8c05c5c1530b97d4816ca434762"
	want := "6acc04142e100a65f51b97adf5172c41"

	K, _ := hex.DecodeString(KHex)
	IV, _ := hex.DecodeString(IVHex)
	P, _ := hex.DecodeString(PHex)

	C := AESCBCEncrypt(K, IV, P)
	got := hex.EncodeToString(C)
	fmt.Printf("C1 : %s \n", want)
	fmt.Printf("C2 : %s \n", got)

	// verify
	assert.Equal(t, want, got, "Should be equal")
}

func TestPQNIST_AES_CBC_DECRYPT(t *testing.T) {
	KHex := "43e953b2aea08a3ad52d182f58c72b9c60fbe4a9ca46a3cb89e3863845e22c9e"
	IVHex := "ddbbb0173f1e2deb2394a62aa2a0240e"
	CHex := "d51d19ded5ca4ae14b2b20b027ffb020"
	want := "07270d0e63aa36daed8c6ade13ac1af1"

	K, _ := hex.DecodeString(KHex)
	IV, _ := hex.DecodeString(IVHex)
	C, _ := hex.DecodeString(CHex)

	P := AESCBCDecrypt(K, IV, C)
	got := hex.EncodeToString(P)
	fmt.Printf("P1 : %s \n", want)
	fmt.Printf("P2 : %s \n", got)

	// verify
	assert.Equal(t, want, got, "Should be equal")
}
func TestENCAP_DECAP(t *testing.T) {
	SEEDHex := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f30"
	SEED, _ := hex.DecodeString(SEEDHex)

	RC1, SIKEpk, SIKEsk, _, _ := Keys(SEED)
	assert.Equal(t, 0, RC1, "Should be equal")

	SIKEpkHex := hex.EncodeToString(SIKEpk)
	SIKEskHex := hex.EncodeToString(SIKEsk)
	fmt.Printf("BLSpk : %s \n", SIKEpkHex)
	fmt.Printf("BLSsk : %s \n", SIKEskHex)

	IVHex := "851e8764776e6796aab722dbb644ace8"
	want := "6282b8c05c5c1530b97d4816ca434762"
	IV, _ := hex.DecodeString(IVHex)
	P, _ := hex.DecodeString(want)

	RC2, C, EK := EncapsulateEncrypt(P, IV, SIKEpk)
	assert.Equal(t, 0, RC2, "Should be equal")
	CHex := hex.EncodeToString(C)
	fmt.Printf("C : %s \n", CHex)
	EKHex := hex.EncodeToString(EK)
	fmt.Printf("EK : %s \n", EKHex)

	RC3, P2 := DecapsulateDecrypt(C, IV, SIKEsk, EK)
	assert.Equal(t, 0, RC3, "Should be equal")
	got := hex.EncodeToString(P2)
	fmt.Printf("want : %s \n", want)
	fmt.Printf("got : %s \n", got)

	assert.Equal(t, want, got, "Should be equal")
}

func TestSIGN_VERIFY(t *testing.T) {
	SEEDHex := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	SEED, _ := hex.DecodeString(SEEDHex)

	RC1, _, _, BLSpk, BLSsk := Keys(SEED)
	assert.Equal(t, 0, RC1, "Should be equal")

	BLSpkHex := hex.EncodeToString(BLSpk)
	BLSskHex := hex.EncodeToString(BLSsk)
	fmt.Printf("BLSpk : %s \n", BLSpkHex)
	fmt.Printf("BLSsk : %s \n", BLSskHex)

	PHex := "48656c6c6f20426f622120546869732069732061206d6573736167652066726f6d20416c696365"
	P, _ := hex.DecodeString(PHex)

	RC2, S := Sign(P, BLSsk)
	assert.Equal(t, 0, RC2, "Should be equal")
	// SHex := hex.EncodeToString(S)
	// fmt.Printf("S : %s \n", SHex)

	want := 0
	got := Verify(P, BLSpk, S)

	assert.Equal(t, want, got, "Should be equal")
}
