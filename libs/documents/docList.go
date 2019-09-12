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

// DocList is an array of all available document types currently supported.
// It is necessary because the envelope body can store many different types of proto structs
// You can't make an instance of a struct by the name of its type in golang, so this list
// facilitates the creation of empty objects of the correct type.

// In SmartDecodeEnvelope
// 1)	DocList is used to determine a textual description of a message at runtime
// 2)	When the message body types are unknown (eg. in the case of the command line tool)
// 	The payload types are taken from the header, these are then used to obtain the DocType
// 	and its subsequent (empty) message (proto.Message).
// 	The objects can then be correctly populated using the type specific DecodeEnvelope

// 	DecodeEnvelope has supplied types which are populated by the envelope and returned

import (
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

//DocType - defines a document which is parseable
//It is necessary to build this list because there is no inheritance in
type DocType struct {
	Name     string //English name of the document
	TypeCode float32
	Version  float32
	Message  proto.Message //An empty instance of the document for comparison
}

//DocList This is a master list of documents supported
var DocList = []DocType{
	//      Name. TypeCode, Version, Message
	{"empty", 0, 1.0, nil},
	{"Simple", 1, 1.0, &SimpleString{}},
	{"PlainTestMessage1", 2, 1.0, &PlainTestMessage1{}},
	{"EncryptTestMessage1", 3, 1.0, &EncryptTestMessage1{}},
	{"IDDocument", 100, 1.0, &IDDocument{}},
	{"OrderDocument", 101, 1.0, &OrderDocument{}},
	{"PolicyDocument", 102, 1.0, &Policy{}},
}

//GetDocTypeForType return the DocType for a given Doc Code
func GetDocTypeForType(docType float32) DocType {
	for _, doc := range DocList {
		if doc.TypeCode == docType {
			return doc
		}
	}
	return DocList[0]
}

//detectDocType Detect the DocType for the supplied protobuf Message, using the DocList
func detectDocType(message proto.Message) (DocType, error) {
	if message == nil {
		return DocList[0], nil
	}
	for _, doc := range DocList {
		if reflect.TypeOf(message) == reflect.TypeOf(doc.Message) {
			return doc, nil
		}
	}
	return DocType{}, errors.New("Document not found")
}
