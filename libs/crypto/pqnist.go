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
Package crypto - wrapper for encryption libraries required by encrypted envelope
*/
package crypto

/*
#cgo CFLAGS: -O2
#cgo LDFLAGS: -lpqnist -lamcl_bls_BLS381 -lamcl_pairing_BLS381 -lamcl_curve_BLS381 -lamcl_core -loqs
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/utils.h>
#include <amcl/randapi.h>
#include <amcl/bls_BLS381.h>
#include <oqs/oqs.h>
#include <pqnist.h>
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

// BLS
const bfsBls381 int = int(C.BFS_BLS381)
const bgsBls381 int = int(C.BGS_BLS381)
const g2Len int = 4 * bfsBls381
const siglen int = bfsBls381 + 1

/*Keys Generate SIKE and BLS keys

Generate SIKE and BLS key public and private key pairs

@param seed             seed value for CSPRNG
@param sikePK           SIKE public key
@param sikeSK           SIKE secret key
@param blsPK            BLS public key
@param blsSK            BLS secret key
@param rc               Return code. Zero for success or else an error code
*/
func Keys(seed []byte) (rc int, sikePK []byte, sikeSK []byte, blsPK []byte, blsSK []byte) {

	// Allocate memory
	psikePK := C.malloc(C.size_t(oqsKemSikeP751LengthPublicKey))
	defer C.free(psikePK)
	psikeSK := C.malloc(C.size_t(oqsKemSikeP751LengthSecretKey))
	defer C.free(psikeSK)
	pblsPK := C.malloc(C.size_t(g2Len))
	defer C.free(pblsPK)
	pblsSK := C.malloc(C.size_t(bgsBls381))
	defer C.free(pblsSK)

	rtn := C.pqnist_keys((*C.char)(unsafe.Pointer(&seed[0])), (*C.char)(psikePK), (*C.char)(psikeSK), (*C.char)(pblsPK), (*C.char)(pblsSK))
	rc = int(rtn)

	sikePK = C.GoBytes(psikePK, C.int(oqsKemSikeP751LengthPublicKey))
	sikeSK = C.GoBytes(psikeSK, C.int(oqsKemSikeP751LengthSecretKey))
	blsPK = C.GoBytes(pblsPK, C.int(g2Len))
	blsSK = C.GoBytes(pblsSK, C.int(bgsBls381))
	return rc, sikePK, sikeSK, blsPK, blsSK
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

/*EncapsulateEncrypt Encrypt a message and encapsulate the AES Key for
  a recipient.

  The  message is encrypted using AES-256. The key
  is generated inside this function as an output
  from the encapsulation function.

  @param P            Plaintext to be encrypted
  @param IV           Initialization vector IV (16 bytes)
  @param sikePK       SIKE public key
  @param C            Ciphertext
  @param EK           Encapsulated key
  @param rc           Return code. Zero for success or else an error code
*/
func EncapsulateEncrypt(pOrig []byte, iv []byte, sikePK []byte) (rc int, c []byte, ek []byte) {
	p := make([]byte, len(pOrig))
	_ = copy(p, pOrig)

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

/*DecapsulateDecrypt  @brief Decapsulate the AES Key and decrypt the message

Decapsulate the AES key and use it to decrypt the
ciphertext.

@param C            Ciphertext to be decrypted
@param IV           Initialization vector IV
@param sikeSK       SIKE secret key
@param EK           Encapsulated key
@param P            Plaintext
@param rc           Return code. Zero for success or else an error code
*/
func DecapsulateDecrypt(cOrig []byte, iv []byte, sikeSK []byte, ek []byte) (rc int, p []byte) {

	c := make([]byte, len(cOrig))
	_ = copy(c, cOrig)

	rtn := C.pqnist_decapsulate_decrypt(
		(*C.char)(unsafe.Pointer(&c[0])),
		C.int(len(c)),
		(*C.char)(unsafe.Pointer(&iv[0])),
		(*C.char)(unsafe.Pointer(&sikeSK[0])),
		(*C.char)(unsafe.Pointer(&ek[0])))
	rc = int(rtn)

	return rc, c
}

/*Sign a message

  The message is signed using the BLS algorithm

  @param M            Message to be signed
  @param blsSK        BLS secret key
  @param S            Signature
  @param rc           Return code. Zero for success or else an error code
*/
func Sign(m []byte, blsSK []byte) (rc int, s []byte) {
	// Allocate memory
	pS := C.malloc(C.size_t(siglen))
	defer C.free(pS)

	rtn := C.pqnist_sign(
		(*C.char)(unsafe.Pointer(&m[0])),
		(*C.char)(unsafe.Pointer(&blsSK[0])),
		(*C.char)(pS))

	rc = int(rtn)

	s = C.GoBytes(pS, C.int(siglen))

	return rc, s
}

/*Verify a signature

  Verify a signature using the BLS algorithm

  @param M            Message that was signed
  @param blsPK        BLS public key
  @param S            Signature
  @param rc           Return code. Zero for success or else an error code
*/
func Verify(m []byte, blsPK []byte, s []byte) (rc int) {

	rtn := C.pqnist_verify(
		(*C.char)(unsafe.Pointer(&m[0])),
		(*C.char)(unsafe.Pointer(&blsPK[0])),
		(*C.char)(unsafe.Pointer(&s[0])))

	rc = int(rtn)

	return rc
}
