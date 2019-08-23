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

package ipfs

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestNodeConnector(t *testing.T) {
	timeOut := time.Second * 5
	ctx, cancel := context.WithCancel(context.Background())

	errChan := make(chan error)
	go func(ctx context.Context, errChan chan error) {
		addr1 := "/ip4/127.0.0.1/tcp/53221"
		addr2 := "/ip4/127.0.0.1/tcp/53222"

		ipfs1, err := NewNodeConnector(
			WithContext(ctx),
			AddLocalAddress(addr1),
			WithMemoryDatastore(),
		)
		if err != nil {
			errChan <- err
			return
		}

		ipfs2, err := NewNodeConnector(
			WithContext(ctx),
			AddLocalAddress(addr2),
			AddBootstrapPeer(fmt.Sprintf("%s/ipfs/%s", addr1, ipfs1.GetID())),
			WithMemoryDatastore(),
		)
		if err != nil {
			t.Fatal(err)
		}

		bsend := []byte{1, 2, 3}
		cid, err := ipfs1.Add(bsend)
		if err != nil {
			errChan <- err
			return
		}

		b, err := ipfs2.Get(cid)
		if err != nil {
			errChan <- err
			return
		}

		if !bytes.Equal(bsend, b) {
			errChan <- errors.Errorf("Expected: %v, found: %v", bsend, b)
			return
		}

		errChan <- nil
	}(ctx, errChan)

	time.AfterFunc(timeOut, func() {
		t.Error("Timeout. Test took too long")
		cancel()
	})

	if err := <-errChan; err != nil {
		t.Fatal(err)
	}

}
