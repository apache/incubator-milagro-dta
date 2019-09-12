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

/*
#cgo CFLAGS: -O2
#cgo LDFLAGS: -lamcl_curve_SECP256K1 -lamcl_core
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <amcl/amcl.h>
#include <amcl/randapi.h>
#include <amcl/ecdh_SECP256K1.h>
*/
import "C"
import (
	"crypto/rand"
	"encoding/hex"

	"github.com/pkg/errors"
)

const (
	hashFunc = C.int(C.SHA256)
)

var (
	paramP1 = CreateOctet([]byte{0, 1, 2})
	paramP2 = CreateOctet([]byte{0, 1, 2, 3})
)

// Secp256k1Encrypt encrypts a message using ECP_SECP256K1_ECIES
func Secp256k1Encrypt(message, publicKey string) (C, V, T string, err error) {
	dec := &hexBatchDecoder{}
	wOctet := dec.decodeOctet(publicKey)
	if dec.err != nil {
		err = dec.err
		return
	}

	seed := make([]byte, 32)
	rand.Read(seed)
	rng := NewRand(seed)

	mOctet := CreateOctet([]byte(message))

	//Results
	vPtr := NewOctet(65)
	hmacPtr := NewOctet(32)
	cypherPtr := NewOctet(len(message) + 16 - (len(message) % 16))

	C.ECP_SECP256K1_ECIES_ENCRYPT(hashFunc, paramP1, paramP2, (*C.csprng)(rng), (*C.octet)(wOctet), (*C.octet)(mOctet), C.int(12), vPtr, cypherPtr, hmacPtr)

	return hex.EncodeToString(cypherPtr.ToBytes()), hex.EncodeToString(vPtr.ToBytes()), hex.EncodeToString(hmacPtr.ToBytes()), nil
}

// Secp256k1Decrypt decrypts an encrypoted message using ECP_SECP256K1_ECIES
func Secp256k1Decrypt(C, V, T, sK string) (message string, err error) {
	dec := &hexBatchDecoder{}
	cOct := dec.decodeOctet(C)
	vOct := dec.decodeOctet(V)
	tOct := dec.decodeOctet(T)
	uOct := dec.decodeOctet(sK)
	if dec.err != nil {
		err = dec.err
		return
	}

	//Cast the cypherText back to Octets
	mOct := NewOctet(len(C) + 16 - (len(C) % 16))

	if C.ECP_SECP256K1_ECIES_DECRYPT(hashFunc, paramP1, paramP2, vOct, cOct, tOct, uOct, mOct) != 1 {
		return "", errors.New("Cannot decrypt cyphertext")
	}

	b := mOct.ToBytes()
	return string(b), nil
}

type hexBatchDecoder struct {
	err error
}

func (d *hexBatchDecoder) decodeOctet(s string) *Octet {
	if d.err != nil {
		return nil
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		d.err = err
		return nil
	}
	return CreateOctet(b)
}
