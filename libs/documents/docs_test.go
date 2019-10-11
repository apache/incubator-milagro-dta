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
	"crypto/sha256"
	"encoding/hex"
	mrand "math/rand"
	"strings"
	"testing"
	"time"

	"github.com/apache/incubator-milagro-dta/libs/crypto"
	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/go-test/deep"
	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

//Test to check that data is being encrypted in Order & not encrypted in ID
func Test_EnvelopeEncryption(t *testing.T) {
	s1, id1, _, _, _, blsSK := BuildTestIDDoc()

	order, _ := BuildTestOrderDoc()
	order.IPFSID = "NEW IPFS ID"
	recipients := map[string]*IDDoc{
		id1: s1,
	}
	testText := "SEARCH_FOR_THIS123"
	testTextBytes := []byte(testText)
	order.Reference = testText
	raw, _ := EncodeOrderDocument(id1, order, blsSK, recipients)
	contains := bytes.Contains(raw, testTextBytes)
	assert.False(t, contains, "Testtext should not be found inside the Envelope - its inside the cipherText")

	iddoc, _, _, _, _, blsSK := BuildTestIDDoc()
	iddoc.AuthenticationReference = testText

	raw2, _ := EncodeIDDocument(iddoc, blsSK)
	contains2 := bytes.Contains(raw2, testTextBytes)

	assert.True(t, contains2, "Testtext should  be found inside the Envelope -its inside the plaintext")

}

func Test_EncodeDecodeOrderDoc(t *testing.T) {
	s1, id1, _, sikeSK, blsPK, blsSK := BuildTestIDDoc()
	order, _ := BuildTestOrderDoc()
	order.IPFSID = "NEW IPFS ID"
	recipients := map[string]*IDDoc{
		id1: s1,
	}
	raw, _ := EncodeOrderDocument(id1, order, blsSK, recipients)
	reconstitutedOrder := OrderDoc{}
	_ = DecodeOrderDocument(raw, "NEW IPFS ID", &reconstitutedOrder, sikeSK, id1, blsPK)
	order.Header.Recipients[0].CipherText = reconstitutedOrder.Header.Recipients[0].CipherText
	differences := deep.Equal(reconstitutedOrder, order)
	var failed = false
	for _, diff := range differences {
		//ignore differences with XXX_ in the names, as these are protobuffer additional fields
		if strings.Contains(diff, "XXX_") == false {
			failed = true
		}
	}
	assert.False(t, failed, "Reconstituted Fields don't match")
	assert.NotNil(t, reconstitutedOrder.OrderPart4.Secret, "Reconstituted Fields dont match")
}

func Test_EncodeDecodeID(t *testing.T) {
	iddoc, tag, _, _, _, blsSK := BuildTestIDDoc()
	raw, _ := EncodeIDDocument(iddoc, blsSK)
	reconstitutedIDDoc := NewIDDoc()
	_ = DecodeIDDocument(raw, tag, reconstitutedIDDoc)
	differences := deep.Equal(reconstitutedIDDoc, iddoc)
	var failed = false
	for _, diff := range differences {
		if strings.Contains(diff, "XXX_") == false {
			failed = true
		}
	}
	assert.False(t, failed, "Reconstituted Fields dont match")
	assert.NotNil(t, reconstitutedIDDoc.DateTime, "Reconstituted Fields dont match")
}

func Test_AESPadding(t *testing.T) {
	for i := 0; i < 1000; i++ {
		randCount := mrand.Intn(100)
		plainText, _ := cryptowallet.RandomBytes(randCount)
		cipherText, aesKey, iv, _ := aesEncrypt(plainText)
		decryptedPlaintext, _ := aesDecrypt(cipherText, iv, aesKey)
		assert.Equal(t, decryptedPlaintext, plainText, "AES round trip fails")
	}
}

func Test_EncodeDecode(t *testing.T) {
	//These are some DocID for local user
	s1, id1, _, sikeSK1, _, _ := BuildTestIDDoc()
	recipients := map[string]*IDDoc{
		id1: s1,
	}
	seed, _ := cryptowallet.RandomBytes(16)
	//Now generate a test Document & use some temp keys, as once made, we needs the TestID above to decode
	_, blsPK, blsSK := crypto.BLSKeys(seed, nil)
	secretBody := &SimpleString{Content: "B"}
	plainText := &SimpleString{Content: "A"}
	header := &Header{}
	rawDoc, err := Encode(id1, plainText, secretBody, header, blsSK, recipients)
	assert.Nil(t, err, "Failed to Encode")

	//Now test Decode,
	sigEnv := &SignedEnvelope{}
	_ = proto.Unmarshal(rawDoc, sigEnv)
	err = Verify(*sigEnv, blsPK)
	assert.Nil(t, err, "Verify fails")

	reconPlainText := &SimpleString{}
	reconSecretBody := &SimpleString{}
	tag := "this is the ipfs id tag"

	reconHeader, err := Decode(rawDoc, tag, sikeSK1, id1, reconPlainText, reconSecretBody, blsPK)

	assert.Nil(t, err, "Verify fails")
	assert.Equal(t, plainText.Content, reconPlainText.Content, "Verify fails")
	assert.Equal(t, secretBody.Content, reconSecretBody.Content, "Verify fails")
	assert.NotNil(t, reconHeader.DateTime, "Header not populated")
	assert.Equal(t, reconHeader.IPFSID, tag, "tag not loaded into header")
}

func BuildTestOrderDoc() (OrderDoc, error) {
	reference, err := uuid.NewUUID()
	if err != nil {
		return OrderDoc{}, err
	}
	order := NewOrderDoc()
	//oder.Type will be used to extend the things that an order can do.
	order.Type = "Safeguard_Secret"
	order.Reference = reference.String()
	order.Coin = 0
	order.BeneficiaryCID = "TESTBeneficiaryIDDocumentCID"
	order.Timestamp = time.Now().Unix()

	o2 := &OrderPart2{}
	o2.CommitmentPublicKey = "2CommitmentPublicKey"
	o2.PreviousOrderCID = "2PreviousOrderCID"
	o2.Timestamp = time.Now().Unix()

	o3 := &OrderPart3{}
	o3.Redemption = "3Redemption"
	o3.PreviousOrderCID = "3PreviousOrderCID"
	o3.BeneficiaryEncryptedData = []byte("3BeneficiaryEncryptedData")
	o3.Timestamp = time.Now().Unix()

	o4 := &OrderPart4{}
	o4.Secret = "4Secret"
	o4.PreviousOrderCID = "4PreviousOrderCID"
	o4.Timestamp = time.Now().Unix()

	order.OrderPart2 = o2
	order.OrderPart3 = o3
	order.OrderPart4 = o4
	return order, nil
}

func BuildTestIDDoc() (*IDDoc, string, []byte, []byte, []byte, []byte) {
	//make some test ID docs
	seed, _ := cryptowallet.RandomBytes(16)

	_, sikePK, sikeSK := crypto.SIKEKeys(seed)

	_, blsPK, blsSK := crypto.BLSKeys(seed, nil)

	//id := []byte("TestID1")
	envelope := Envelope{}
	header := Header{}
	idDocument := IDDocument{}
	idDocument.SikePublicKey = sikePK
	idDocument.BLSPublicKey = blsPK

	//assemble
	envelope.Header = &header
	envelope.Body, _ = proto.Marshal(&idDocument)
	envelope.EncryptedBody = nil
	header.EncryptedBodyIV = nil
	signedEnvelope, _ := sign(envelope, blsSK, "TESTID")
	ipfsID, _ := createIDForSignedEnvelope(signedEnvelope)
	iddoc := &IDDoc{}
	rawDocI, _ := proto.Marshal(&signedEnvelope)
	_ = DecodeIDDocument(rawDocI, ipfsID, iddoc)
	return iddoc, ipfsID, sikePK, sikeSK, blsPK, blsSK
}

//createIDForSignedEnvelope - create a hash for the document to be used as an ID
func createIDForSignedEnvelope(signedEnvelope SignedEnvelope) (string, []byte) {
	dat, _ := proto.Marshal(&signedEnvelope)
	return sha256Hash(string(dat)), dat
}

func sha256Hash(data string) string {
	hasher := sha256.New()
	_, _ = hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

//Use - helper to remove warnings
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
