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
	"encoding/json"
	"io/ioutil"
	"net"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/pkg/errors"
)

// APIConnector is IPFS Shell API Connector
type APIConnector struct {
	shell *shell.Shell
	id    string
}

// NewAPIConnector inisialises new IPFS API connector
func NewAPIConnector(options ...APIConnectorOption) (Connector, error) {

	cb := &APIConnectorBuilder{
		nodeAddr: "localhost:5001",
	}

	for _, option := range options {
		if err := option(cb); err != nil {
			return nil, err
		}
	}

	sh := shell.NewShell(cb.nodeAddr)

	outID, err := sh.ID()
	if err != nil {
		return nil, errors.Wrap(ErrNodeConnection, err.Error())
	}

	if !sh.IsUp() {
		return nil, ErrNodeConnection
	}

	if cb.swarmPeerAddr != "" {
		_ = sh.SwarmConnect(context.Background(), cb.swarmPeerAddr)
	}

	return &APIConnector{
		shell: sh,
		id:    outID.ID,
	}, nil
}

// Add adds a data to the IPFS network and returns the ipfs path
func (c *APIConnector) Add(data []byte) (string, error) {
	return c.shell.Add(bytes.NewReader(data), shell.Pin(true), shell.Progress(false))
}

// Get gets a data from ipfs path
func (c *APIConnector) Get(path string) ([]byte, error) {
	r, err := c.shell.Cat(path)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Pin the document
	if err := c.shell.Pin(path); err != nil {
		return nil, err
	}

	return b, nil
}

// AddJSON encodes the data in JSON and adds it to the IPFS
func (c *APIConnector) AddJSON(data interface{}) (string, error) {
	jd, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return c.Add(jd)
}

// GetJSON gets data from IPFS and decodes it from JSON
func (c *APIConnector) GetJSON(path string, data interface{}) error {
	jd, err := c.Get(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jd, data); err != nil {
		return errors.Wrap(ErrDocumentNotValid, err.Error())
	}

	return nil
}

// GetID returns the local id
func (c *APIConnector) GetID() string {
	return c.id
}

// APIConnectorBuilder for building the IPFS API connector
type APIConnectorBuilder struct {
	nodeAddr      string
	swarmPeerAddr string
}

// APIConnectorOption function
type APIConnectorOption func(*APIConnectorBuilder) error

// NodeAddr specifies the IPFS local node address
func NodeAddr(addr string) APIConnectorOption {
	return func(cb *APIConnectorBuilder) error {
		cb.nodeAddr = addr
		return nil
	}
}

// PeerAddr specifies the IPFS node peer to connect
func PeerAddr(peerAddr string) APIConnectorOption {
	return func(cb *APIConnectorBuilder) error {
		cb.swarmPeerAddr = peerAddr
		return nil
	}
}

// PeerDomain returns the IPFS peer addr from domain TXT record
func PeerDomain(domain string) APIConnectorOption {
	return func(cb *APIConnectorBuilder) error {
		var addr string
		records, err := net.LookupTXT(domain)
		if err != nil {
			return err
		}

		for _, r := range records {
			a := strings.TrimPrefix(strings.TrimSuffix(r, `"`), `"`)
			if len(a) > 0 && a[0] == '/' {
				addr = a
			}
		}

		if addr == "" {
			return errors.Errorf("unable to find the IPFS address from TXT record of %v", domain)
		}

		cb.swarmPeerAddr = addr
		return nil
	}
}
