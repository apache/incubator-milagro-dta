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
	"fmt"
	"testing"
	"time"

	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	multihash "github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
)

var (
	pass = true
	fail = false
)

func Test_1(t *testing.T) {
	now := time.Now().Unix()

	msg := &OrderDocument{
		Coin:      10,
		Reference: "ABC",
		Timestamp: now,
	}
	err := msg.Validate()
	assert.Nil(t, err, "Validation Failed")
}

func Test_SignedEnvelope(t *testing.T) {
	se := &SignedEnvelope{
		Signature: []byte("012345678901234567890123456789"),
		SignerCID: "",
		Message:   nil,
	}
	err := se.Validate()
	assert.Nil(t, err, "Validation Failed1")

}

func Test_IPFSRegex(t *testing.T) {

	rec := &Recipient{
		Version: 1,
		CID:     "QmSkKsExuUZJ9ETraPB9ZA9KySYGynvx1jK7LmWgRxinhx",
	}
	err := rec.Validate()
	assert.Nil(t, err, "Validation Failed1")

	rec = &Recipient{
		Version: 1,
		CID:     "QmSkKsExuUZJ9ETraPB9ZA9KySYGynvx1jK7LmWgRxinhx1", //extra char
	}
	err = rec.Validate()
	assert.NotNil(t, err, "Validation Failed2")

	rec = &Recipient{
		Version: 1,
		CID:     "QmSkKsExuUZJ9ETraPB9ZA9KySYGynvx1jK7LmWgRxinh", //less 1 char
	}
	err = rec.Validate()
	assert.NotNil(t, err, "Validation Failed3")

	rec = &Recipient{
		Version: 1,
		CID:     "QmSkKsExuUZJ9ETraPB9;A9KySYGynvx1jK7LmWgRxinhx", //puncuation
	}
	err = rec.Validate()
	assert.NotNil(t, err, "Validation Failed4")

	rec = &Recipient{
		Version: 1,
		CID:     "",
	}
	err = rec.Validate()
	assert.Nil(t, err, "Validation Failed5")

	rec = &Recipient{}
	err = rec.Validate()
	assert.Nil(t, err, "Validation Failed6")

	for i := 0; i < 1000; i++ {
		rnddata, _ := cryptowallet.RandomBytes(32)
		mh, _ := multihash.Sum(rnddata, multihash.SHA2_256, 32)
		cid := mh.B58String()
		cid = fmt.Sprintf("%s", cid)
		rec := &Recipient{
			Version: 1,
			CID:     cid,
		}
		err := rec.Validate()
		assert.Nil(t, err, "Validation Failed1")
	}
}
