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

package documents

import (
	"bytes"
	"fmt"

	"github.com/apache/incubator-milagro-dta/libs/crypto"
	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	proto "github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

var (
	errRecipientNotFound      = errors.New("Recipient not found")
	errFailedDecapsulation    = errors.New("Failed to decapsulate AES key")
	errFailedToGenerateAESKey = errors.New("Failed to generate Random aesKey")
)

//decapsulate - decapsulate the aes for Recipient ID in the list
func decapsulate(recipientCID string, recipients []*Recipient, sikeSK []byte) ([]byte, error) {
	for _, recipient := range recipients {
		if recipient.CID == recipientCID {
			return decapsulateWithRecipient(*recipient, sikeSK)
		}
	}
	return nil, errRecipientNotFound
}

func decapsulateWithRecipient(recipient Recipient, sikeSK []byte) ([]byte, error) {
	cipherText := recipient.CipherText
	encapsulatedKey := recipient.EncapsulatedKey
	encapIV := recipient.IV

	cipherTextTemp := make([]byte, len(cipherText))
	encapIVTemp := make([]byte, len(encapIV))
	sikeSKTemp := make([]byte, len(sikeSK))
	encapsulatedKeyTemp := make([]byte, len(encapsulatedKey))

	copy(cipherTextTemp, cipherText)
	copy(encapIVTemp, encapIV)
	copy(sikeSKTemp, sikeSK)
	copy(encapsulatedKeyTemp, encapsulatedKey)

	rc, recreatedAesKey := crypto.DecapsulateDecrypt(cipherTextTemp, encapIVTemp, sikeSKTemp, encapsulatedKeyTemp)

	if rc != 0 {
		return nil, errFailedDecapsulation
	}
	return recreatedAesKey, nil
}

func encapsulateKeyForRecipient(recipientsIDDocs map[string]IDDoc, secret []byte) (recipientList []*Recipient, err error) {
	for id, idDocument := range recipientsIDDocs {
		r := &Recipient{}
		iv, err := cryptowallet.RandomBytes(16)
		if err != nil {
			return nil, errFailedToGenerateAESKey
		}
		r.CID = id
		r.IV = iv
		sikePK := idDocument.SikePublicKey

		//Make Copies of EncapsulateEncrypt's input params as its destructive
		secretTemp := make([]byte, len(secret))
		ivTemp := make([]byte, len(iv))
		sikePKtemp := make([]byte, len(sikePK))
		copy(secretTemp, secret)
		copy(ivTemp, iv)
		copy(sikePKtemp, sikePK)

		rc, cipherText, encapsulatedKey := crypto.EncapsulateEncrypt(secretTemp, ivTemp, sikePKtemp)

		if rc != 0 {
			return nil, errFailedToGenerateAESKey
		}
		r.EncapsulatedKey = encapsulatedKey
		r.CipherText = cipherText
		recipientList = append(recipientList, r)
	}

	return recipientList, nil
}

func aesEncrypt(plainText []byte) (cipherText []byte, aesKey []byte, iv []byte, err error) {
	aesKey, err = cryptowallet.RandomBytes(32)
	if err != nil {
		return nil, nil, nil, errFailedToGenerateAESKey
	}
	iv, err = cryptowallet.RandomBytes(16)
	if err != nil {
		return nil, nil, nil, errFailedToGenerateAESKey
	}
	paddedText, err := pkcs7Pad(plainText, 32)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "Failed to pad Document secret")
	}
	cipherText = crypto.AESCBCEncrypt(aesKey, iv, paddedText)
	return cipherText, aesKey, iv, nil
}

func aesDecrypt(cipherText []byte, iv []byte, aesKey []byte) (plainText []byte, err error) {
	pt := crypto.AESCBCDecrypt(aesKey, iv, cipherText)
	plainText, err = pkcs7Unpad(pt, 32)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

//Sign - generate a Signed envelope from the envelope
func sign(envelope Envelope, blsSK []byte, signerNodeID string) (SignedEnvelope, error) {
	envelopeBytes, err := proto.Marshal(&envelope)
	if err != nil {
		return SignedEnvelope{}, errors.Wrap(err, "Failed to serialize envelope in SignBLS")
	}
	rc, signature := crypto.BLSSign(envelopeBytes, blsSK)
	if rc != 0 {
		return SignedEnvelope{}, errors.Wrap(err, "Failed to sign envelope in in SignBLS")
	}
	signedEnvelope := SignedEnvelope{}
	signedEnvelope.SignerCID = signerNodeID
	signedEnvelope.Message = envelopeBytes
	signedEnvelope.Signature = signature
	return signedEnvelope, nil
}

//Verify verify the envelopes BLS signature
func Verify(signedEnvelope SignedEnvelope, blsPK []byte) error {
	message := signedEnvelope.Message
	signature := signedEnvelope.Signature

	rc := crypto.BLSVerify(message, blsPK, signature)
	if rc == 0 {
		return nil
	}
	return errors.New("invalid signature")
}

// Appends padding.
func pkcs7Pad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	padlen := 1
	for ((len(data) + padlen) % blocklen) != 0 {
		padlen = padlen + 1
	}
	pad := bytes.Repeat([]byte{byte(padlen)}, padlen)
	return append(data, pad...), nil
}

// Returns slice of the original data without padding.
func pkcs7Unpad(data []byte, blocklen int) ([]byte, error) {
	if blocklen <= 0 {
		return nil, fmt.Errorf("invalid blocklen %d", blocklen)
	}
	if len(data)%blocklen != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen > blocklen || padlen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	// check padding
	pad := data[len(data)-padlen:]
	for i := 0; i < padlen; i++ {
		if pad[i] != byte(padlen) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:len(data)-padlen], nil
}
