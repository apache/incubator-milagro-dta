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
Package ipfs - connect, get and set operations for embedded or external IPFS node
*/
package ipfs

import (
	"github.com/pkg/errors"
)

var (
	// ErrNodeConnection when the ipfs node is not accessible
	ErrNodeConnection = errors.New("ipfs node connection problem")
	// ErrDocumentNotValid when the ipfs document is not well formatted
	ErrDocumentNotValid = errors.New("document not valid")
	// ErrDocumentNotFound when the ipfs document is not dound
	ErrDocumentNotFound = errors.New("document not found")
)

// Connector is the IPFS connector interface
type Connector interface {
	GetID() string
	Add(data []byte) (string, error)
	Get(path string) ([]byte, error)
	AddJSON(data interface{}) (string, error)
	GetJSON(path string, data interface{}) error
}
