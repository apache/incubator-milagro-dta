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
Package crypto - wrapper for encryption libraries required by service
*/
package crypto

/*
#cgo CFLAGS:  -O2 -I/amcl -I/usr/local/include/amcl
#cgo LDFLAGS: -L. -lpqnist -lamcl_bls_BLS381 -lamcl_pairing_BLS381 -lamcl_curve_BLS381 -lamcl_core -loqs
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <amcl/bls_BLS381.h>
#include <oqs/oqs.h>
#include <amcl/pqnist.h>
*/
import "C"
import "unsafe"

// AES
const pqnistAesKeyLength int = int(C.PQNIST_AES_KEY_LENGTH)
const pqnistAesIvLength int = int(C.PQNIST_AES_IV_LENGTH)

// CSPRNG
const pqnistSeedLength int = int(C.PQNIST_SEED_LENGTH)

// SIKE
const oqsKemSikeP751LengthPublicKey int = int(C.OQS_KEM_sike_p751_length_public_key)
const oqsKemSikeP751LengthSecretKey int = int(C.OQS_KEM_sike_p751_length_secret_key)
const oqsKemSikeP751LengthCiphertext int = int(C.OQS_KEM_sike_p751_length_ciphertext)

// BFSBLS381 Field size
const BFSBLS381 int = int(C.BFS_BLS381)

// BGSBLS381 Group size
const BGSBLS381 int = int(C.BGS_BLS381)

// G2Len G2 point size
const G2Len int = 4 * BFSBLS381

// SIGLen Signature length
const SIGLen int = BFSBLS381 + 1

/*BLSKeys Generate BLS keys

Generate public and private key pair. If the seed value is nil then
generate the public key using the input secret key.

@param seed             seed value for CSPRNG.
@param ski              input secret key
@param pk               public key
@param sko              output secret key
@param rc               Return code. Zero for success or else an error code
*/
func BLSKeys(seed []byte, ski []byte) (rc int, pk []byte, sko []byte) {

	// Allocate memory
	ppk := C.malloc(C.size_t(G2Len))
	defer C.free(ppk)
	var sk []byte

	if seed == nil {
		rtn := C.pqnist_bls_keys(nil, (*C.char)(ppk), (*C.char)(unsafe.Pointer(&ski[0])))
		rc = int(rtn)

		pk = C.GoBytes(ppk, C.int(G2Len))
		sk = ski

	} else {
		psk := C.malloc(C.size_t(BGSBLS381))
		defer C.free(psk)

		rtn := C.pqnist_bls_keys((*C.char)(unsafe.Pointer(&seed[0])), (*C.char)(ppk), (*C.char)(psk))
		rc = int(rtn)

		pk = C.GoBytes(ppk, C.int(G2Len))
		sko = C.GoBytes(psk, C.int(BGSBLS381))
		sk = sko
	}
	return rc, pk, sk
}

/*BLSSign Sign a message

  The message is signed using the BLS algorithm

  @param M            Message to be signed
  @param sk           secret key
  @param S            Signature
  @param rc           Return code. Zero for success or else an error code
*/
func BLSSign(m []byte, sk []byte) (rc int, s []byte) {
	// Allocate memory
	pS := C.malloc(C.size_t(SIGLen))
	defer C.free(pS)

	rtn := C.pqnist_bls_sign(
		(*C.char)(unsafe.Pointer(&m[0])),
		(*C.char)(unsafe.Pointer(&sk[0])),
		(*C.char)(pS))

	rc = int(rtn)

	s = C.GoBytes(pS, C.int(SIGLen))

	return rc, s
}

/*BLSVerify Verify a signature

  Verify a signature using the BLS algorithm

  @param M            Message that was signed
  @param pk           public key
  @param S            Signature
  @param rc           Return code. Zero for success or else an error code
*/
func BLSVerify(m []byte, pk []byte, s []byte) (rc int) {

	rtn := C.pqnist_bls_verify(
		(*C.char)(unsafe.Pointer(&m[0])),
		(*C.char)(unsafe.Pointer(&pk[0])),
		(*C.char)(unsafe.Pointer(&s[0])))

	rc = int(rtn)

	return rc
}

/*BLSAddG1 Add two members from the group G1

  Add two members from the group G1

  @param R1           member of G1
  @param R2           member of G1
  @param R            member of G1. r = r1+r2
  @param rc           Return code. Zero for success or else an error code
*/
func BLSAddG1(R1 []byte, R2 []byte) (rc int, R []byte) {

	// Allocate memory
	pR := C.malloc(C.size_t(SIGLen))
	defer C.free(pR)

	rtn := C.pqnist_bls_addg1(
		(*C.char)(unsafe.Pointer(&R1[0])),
		(*C.char)(unsafe.Pointer(&R2[0])),
		(*C.char)(pR))

	rc = int(rtn)

	R = C.GoBytes(pR, C.int(SIGLen))

	return rc, R
}

/*BLSAddG2 Add two members from the group G2

  Add two members from the group G2

  @param R1           member of G2
  @param R2           member of G2
  @param R            member of G2. r = r1+r2
  @param rc           Return code. Zero for success or else an error code
*/
func BLSAddG2(R1 []byte, R2 []byte) (rc int, R []byte) {

	// Allocate memory
	pR := C.malloc(C.size_t(G2Len))
	defer C.free(pR)

	rtn := C.pqnist_bls_addg2(
		(*C.char)(unsafe.Pointer(&R1[0])),
		(*C.char)(unsafe.Pointer(&R2[0])),
		(*C.char)(pR))

	rc = int(rtn)

	R = C.GoBytes(pR, C.int(G2Len))

	return rc, R
}

/*BLSMakeShares Use Shamir's secret sharing to distribute BLS secret keys

Use Shamir's secret sharing to distribute BLS secret keys

@param  k       Threshold
@param  n       Number of shares
@param  seed    seed value for CSPRNG
@param  ski     Secret key to be shared.
@param  x       x values
@param  y       y values. Valid BLS secret keys
@param  rc      Zero for success or else an error code
*/
func BLSMakeShares(k int, n int, seed []byte, ski []byte) (rc int, x []byte, y []byte, sko []byte) {

	// Allocate memory
	pX := C.malloc(C.size_t(BGSBLS381 * n))
	defer C.free(pX)
	pY := C.malloc(C.size_t(BGSBLS381 * n))
	defer C.free(pY)
	pSKO := C.malloc(C.size_t(BGSBLS381))
	defer C.free(pSKO)

	rtn := C.pqnist_bls_make_shares(C.int(k), C.int(n), (*C.char)(unsafe.Pointer(&seed[0])), (*C.char)(pX), (*C.char)(pY), (*C.char)(unsafe.Pointer(&ski[0])), (*C.char)(pSKO))
	rc = int(rtn)

	sko = C.GoBytes(pSKO, C.int(BGSBLS381))
	x = C.GoBytes(pX, C.int(BGSBLS381*n))
	y = C.GoBytes(pY, C.int(BGSBLS381*n))
	return rc, x, y, sko
}

/*BLSRecoverSecret Use Shamir's secret sharing to recover a BLS secret key

Use Shamir's secret sharing to recover a BLS secret key

@param  k       Threshold
@param  x       x values
@param  y       y values. Valid BLS secret keys
@param  sk      Secret key that is recovered
@param  rc      Zero for success or else an error code
*/
func BLSRecoverSecret(k int, x []byte, y []byte) (rc int, sk []byte) {

	// Allocate memory
	pSK := C.malloc(C.size_t(BGSBLS381))
	defer C.free(pSK)

	rtn := C.pqnist_bls_recover_secret(C.int(k), (*C.char)(unsafe.Pointer(&x[0])), (*C.char)(unsafe.Pointer(&y[0])), (*C.char)(pSK))
	rc = int(rtn)

	sk = C.GoBytes(pSK, C.int(BGSBLS381))

	return rc, sk
}

/*BLSRecoverSignature Use Shamir's secret sharing to recover a BLS signature

Use Shamir's secret sharing to recover a BLS signature

@param  k       Threshold
@param  x       x values
@param  y       y values. Valid BLS signatures
@param  sig     Signature that is recovered
@param  rc      Zero for success or else an error code
*/
func BLSRecoverSignature(k int, x []byte, y []byte) (rc int, sig []byte) {

	// Allocate memory
	pSIG := C.malloc(C.size_t(SIGLen))
	defer C.free(pSIG)

	rtn := C.pqnist_bls_recover_signature(C.int(k), (*C.char)(unsafe.Pointer(&x[0])), (*C.char)(unsafe.Pointer(&y[0])), (*C.char)(pSIG))
	rc = int(rtn)

	sig = C.GoBytes(pSIG, C.int(SIGLen))

	return rc, sig
}

/*AESCBCEncrypt AES-CBC Encryption

  AES encryption using CBC mode

  @param K            Key
  @param IV           Initialization vector IV (16 bytes)
  @param P            Plaintext
  @param C            Ciphertext
*/
func AESCBCEncrypt(k []byte, iv []byte, p []byte) (c []byte) {
	C.pqnist_aes_cbc_encrypt(
		(*C.char)(unsafe.Pointer(&k[0])),
		C.int(len(k)),
		(*C.char)(unsafe.Pointer(&iv[0])),
		(*C.char)(unsafe.Pointer(&p[0])),
		C.int(len(p)))

	return p
}

/*AESCBCDecrypt AES-CBC Decryption

  AES decryption using CBC mode

  @param K            Key
  @param IV           Initialization vector IV (16 bytes)
  @param C            Ciphertext
  @param P            Plaintext
*/
func AESCBCDecrypt(k []byte, iv []byte, c []byte) (p []byte) {

	C.pqnist_aes_cbc_decrypt(
		(*C.char)(unsafe.Pointer(&k[0])),
		C.int(len(k)),
		(*C.char)(unsafe.Pointer(&iv[0])),
		(*C.char)(unsafe.Pointer(&c[0])),
		C.int(len(c)))

	return c
}

/*SIKEKeys Generate SIKE keys

Generate SIKE public and private key pair

@param seed             seed value for CSPRNG
@param sikePK           SIKE public key
@param sikeSK           SIKE secret key
@param rc               Return code. Zero for success or else an error code
*/
func SIKEKeys(seed []byte) (rc int, sikePK []byte, sikeSK []byte) {

	// Allocate memory
	psikePK := C.malloc(C.size_t(oqsKemSikeP751LengthPublicKey))
	defer C.free(psikePK)
	psikeSK := C.malloc(C.size_t(oqsKemSikeP751LengthSecretKey))
	defer C.free(psikeSK)

	rtn := C.pqnist_sike_keys((*C.char)(unsafe.Pointer(&seed[0])), (*C.char)(psikePK), (*C.char)(psikeSK))
	rc = int(rtn)

	sikePK = C.GoBytes(psikePK, C.int(oqsKemSikeP751LengthPublicKey))
	sikeSK = C.GoBytes(psikeSK, C.int(oqsKemSikeP751LengthSecretKey))
	return rc, sikePK, sikeSK
}

/*EncapsulateEncrypt Encrypt a message and encapsulate the AES Key for a recipient.

  The  message is encrypted using AES-256. The key
  is generated inside this function as an output
  from the encapsulation function. The ciphertext
  is returned using the P paramter.

  @param P            Plaintext to be encrypted
  @param IV           Initialization vector IV (16 bytes)
  @param sikePK       SIKE public key
  @param C            Ciphertext
  @param EK           Encapsulated key
  @param rc           Return code. Zero for success or else an error code
*/
func EncapsulateEncrypt(p []byte, iv []byte, sikePK []byte) (rc int, c []byte, ek []byte) {

	// Allocate memory
	pEK := C.malloc(C.size_t(oqsKemSikeP751LengthCiphertext))
	defer C.free(pEK)

	rtn := C.pqnist_encapsulate_encrypt(
		(*C.char)(unsafe.Pointer(&p[0])),
		C.int(len(p)),
		(*C.char)(unsafe.Pointer(&iv[0])),
		(*C.char)(unsafe.Pointer(&sikePK[0])),
		(*C.char)(pEK))
	rc = int(rtn)

	ek = C.GoBytes(pEK, C.int(oqsKemSikeP751LengthCiphertext))

	return rc, p, ek
}

/*DecapsulateDecrypt Decapsulate the AES Key and decrypt the message

Decapsulate the AES key and use it to decrypt the
ciphertext. The plaintext is returned using the C
parameter.

@param C            Ciphertext to be decrypted
@param IV           Initialization vector IV
@param sikeSK       SIKE secret key
@param EK           Encapsulated key
@param P            Plaintext
@param rc           Return code. Zero for success or else an error code
*/
func DecapsulateDecrypt(c []byte, iv []byte, sikeSK []byte, ek []byte) (rc int, p []byte) {

	rtn := C.pqnist_decapsulate_decrypt(
		(*C.char)(unsafe.Pointer(&c[0])),
		C.int(len(c)),
		(*C.char)(unsafe.Pointer(&iv[0])),
		(*C.char)(unsafe.Pointer(&sikeSK[0])),
		(*C.char)(unsafe.Pointer(&ek[0])))
	rc = int(rtn)

	return rc, c
}
