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
Package documents - data is signed and nested in "encrypted envelope"
*/
package documents

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

//DocType - defines a document that is parseable
//It is necessary to build this list because there is no inheritance in

var (
	//EnvelopeVersion the versioning of the entire Envelope, (not individual documents/contents)
	EnvelopeVersion float32 = 1.0
)

//IDDoc wrapper to encapsulate Header & IDDocument into one object
type IDDoc struct {
	*Header
	*IDDocument
}

//OrderDoc puts and order in encrypted wrapper
type OrderDoc struct {
	*Header
	*OrderDocument
}

//NewIDDoc generate a new empty IDDoc
func NewIDDoc() IDDoc {
	ret := IDDoc{}
	ret.Header = &Header{}
	ret.IDDocument = &IDDocument{}
	return ret
}

//NewOrderDoc generate a new order
func NewOrderDoc() OrderDoc {
	ret := OrderDoc{}
	ret.Header = &Header{}
	ret.OrderDocument = &OrderDocument{}
	return ret
}

//EncodeIDDocument encode an IDDoc into a raw bytes stream for the wire
func EncodeIDDocument(idDocument IDDoc, blsSK []byte) ([]byte, error) {
	header := idDocument.Header
	plaintext := idDocument.IDDocument
	rawDoc, err := Encode("", plaintext, nil, header, blsSK, nil)
	return rawDoc, err
}

//EncodeOrderDocument encode an OrderDoc into a raw bytes stream for the wire
func EncodeOrderDocument(nodeID string, orderDoc OrderDoc, blsSK []byte, previousCID string, recipients map[string]IDDoc) ([]byte, error) {
	header := orderDoc.Header
	header.PreviousCID = previousCID
	//	rawDoc, err := Encode(orderDoc.OrderDocument, nil, header, blsSK, nil)
	rawDoc, err := Encode(nodeID, nil, orderDoc.OrderDocument, header, blsSK, recipients)
	return rawDoc, err
}

//DecodeIDDocument - decode a raw byte stream into an IDDocument
func DecodeIDDocument(rawdoc []byte, tag string, idDocument *IDDoc) error {
	plainText := IDDocument{}
	header, err := Decode(rawdoc, tag, nil, "", &plainText, nil, nil)
	if err != nil {
		return errors.Wrap(err, "DecodeIDDocument Failed to Decode")
	}
	idDocument.Header = header
	idDocument.IDDocument = &plainText

	//validate the order document
	err = idDocument.IDDocument.Validate()
	if err != nil {
		return err
	}
	return nil
}

//PeekOrderDocument - look at the header inside an order document before decryption
func OrderDocumentSigner(rawDoc []byte) (string, error) {
	signedEnvelope := SignedEnvelope{}
	err := proto.Unmarshal(rawDoc, &signedEnvelope)
	if err != nil {
		return "", errors.New("Protobuf - Failed to unmarshal Signed Envelope")
	}
	return signedEnvelope.SignerCID, nil

}

//DecodeOrderDocument -
func DecodeOrderDocument(rawdoc []byte, tag string, orderdoc *OrderDoc, sikeSK []byte, recipientCID string, sendersBlsPK []byte) error {
	cipherText := OrderDocument{}
	header, err := Decode(rawdoc, tag, sikeSK, recipientCID, nil, &cipherText, sendersBlsPK)
	if err != nil {
		return errors.Wrap(err, "DecodeIDDocument Failed to Decode")
	}
	orderdoc.Header = header
	orderdoc.OrderDocument = &cipherText

	//validate the order document
	err = orderdoc.OrderDocument.Validate()
	if err != nil {
		return err
	}
	return nil
}

//Decode - Given a raw envelope, Sike Secret Key & ID - decode into plaintext, ciphertext(decrypted) and header
func Decode(rawDoc []byte, tag string, sikeSK []byte, recipientID string, plainText proto.Message, encryptedText proto.Message, sendersBlsPK []byte) (header *Header, err error) {
	signedEnvelope := SignedEnvelope{}
	err = proto.Unmarshal(rawDoc, &signedEnvelope)
	if err != nil {
		return &Header{}, errors.New("Protobuf - Failed to unmarshal Signed Envelope")
	}

	//check the message verification if we have a key for it
	if sendersBlsPK != nil {
		//if the document is locally signed check with our signature
		err = Verify(signedEnvelope, sendersBlsPK)
		if err != nil {
			return nil, err
		}
	}

	//Decode the  envelope
	message := signedEnvelope.Message
	envelope := Envelope{}
	err = proto.Unmarshal(message, &envelope)
	if err != nil {
		return &Header{}, errors.New("Protobuf - Failed to unmarshal Envelope")
	}
	header = envelope.Header
	header.IPFSID = tag
	//Decode the plaintext
	if plainText != nil {
		err = proto.Unmarshal(envelope.Body, plainText)
		if err != nil {
			return &Header{}, errors.New("Protobuf - Failed to unmarshall plaintext")
		}
	}
	//Decrypt the cipherText & decode into the correct object
	if encryptedText != nil {
		recipientList := header.Recipients
		aesKey, err := decapsulate(recipientID, recipientList, sikeSK)

		if err != nil {
			return &Header{}, errors.Wrap(err, "Failed to Decapsulate Encrypted Text in Envelope Decode")
		}

		decryptedCipherText, err := aesDecrypt(envelope.EncryptedBody, header.EncryptedBodyIV, aesKey)
		if err != nil {
			return &Header{}, errors.Wrap(err, "Failed to AES Decrypt Envelope cipherText")
		}
		err = proto.Unmarshal(decryptedCipherText, encryptedText)
		if err != nil {
			return &Header{}, errors.New("Failed to unmarshall ciphertext")
		}
	}
	return header, nil
}

//Encode - convert the header, secret and plaintext into a message for the wire
//The Header can be pre-populated with any nece
func Encode(nodeID string, plainText proto.Message, secretText proto.Message, header *Header, blsSK []byte, recipients map[string]IDDoc) (rawDoc []byte, err error) {
	plainTextDocType, err := detectDocType(plainText)
	if err != nil {
		return nil, errors.New("Plaintext Document - Unknown Type")
	}
	encryptedTextDocType, err := detectDocType(secretText)
	if err != nil {
		return nil, errors.New("Encrypted Document - Unknown Type")
	}
	//build Header
	header.Version = EnvelopeVersion
	//TODO:  change datetime to real value - its 0 to create fixed IPFS refs
	header.DateTime = 0 //time.Now().Unix()
	header.BodyTypeCode = plainTextDocType.TypeCode
	header.BodyVersion = plainTextDocType.Version
	header.EncryptedBodyTypeCode = encryptedTextDocType.TypeCode
	header.EncryptedBodyVersion = encryptedTextDocType.Version
	//Plaintext to bytes
	var body []byte
	if plainText != nil {
		body, err = proto.Marshal(plainText)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to marshall Plaintext to Protobuf")
		}
	}
	//Ciphertext
	var cipherText []byte
	if secretText != nil {
		secretBody, err := proto.Marshal(secretText)
		if err != nil {
			return nil, err
		}
		var aesKey, iv []byte
		cipherText, aesKey, iv, err = aesEncrypt(secretBody)

		if err != nil {
			return nil, err
		}
		header.EncryptedBodyIV = iv
		recipientList, err := encapsulateKeyForRecipient(recipients, aesKey)
		if err != nil {
			return nil, err
		}
		header.Recipients = recipientList
	}
	//assemble
	envelope := Envelope{}
	envelope.Header = header
	envelope.Body = body
	envelope.EncryptedBody = cipherText
	//sign
	signedEnvelope, err := sign(envelope, blsSK, nodeID)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to Sign Envelope")
	}

	//SignedEnvelope to bytes
	rawDoc, err = proto.Marshal(&signedEnvelope)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshall Signed Envelope to Protobuf")
	}
	return rawDoc, nil
}
