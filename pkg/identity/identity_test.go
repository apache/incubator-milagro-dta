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

package identity

import (
	"testing"

	"github.com/apache/incubator-milagro-dta/libs/ipfs"
	"github.com/apache/incubator-milagro-dta/libs/keystore"
)

func TestCreateIdentity(t *testing.T) {

	ipfsNode, err := ipfs.NewMemoryConnector()
	if err != nil {
		t.Fatal(err)
	}

	store, _ := keystore.NewMemoryStore()

	_, rawIDDoc, secret, err := CreateIdentity("test")
	if err != nil {
		t.Fatal(err)
	}

	idDocID, err := StoreIdentity(rawIDDoc, secret, ipfsNode, store)

	if err := CheckIdentity(idDocID, "test", ipfsNode, store); err != nil {
		t.Fatal(err)
	}

}
