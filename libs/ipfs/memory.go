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
	"encoding/json"
	"sync"

	multihash "github.com/multiformats/go-multihash"
	"github.com/pkg/errors"
)

// MemoryConnector ipmlements IPFS Connector interface with Memory storage
type MemoryConnector struct {
	mutex sync.RWMutex
	store map[string][]byte
}

// NewMemoryConnector creates a new MemoryConnector struct
func NewMemoryConnector() (Connector, error) {
	return &MemoryConnector{
		mutex: sync.RWMutex{},
		store: map[string][]byte{},
	}, nil
}

// Add adds data to Memory IPFS
func (m *MemoryConnector) Add(data []byte) (string, error) {
	id, err := genIPFSID(data)
	if err != nil {
		return "", err
	}

	m.mutex.Lock()
	m.store[id] = data
	m.mutex.Unlock()

	return id, nil
}

// Get gets data from Memory IPFS
func (m *MemoryConnector) Get(path string) ([]byte, error) {
	m.mutex.RLock()
	data, ok := m.store[path]
	m.mutex.RUnlock()
	if !ok {
		return nil, errors.New("Key not found")
	}

	return data, nil
}

// AddJSON encodes data to JSON and stores it in Memory IPFS
func (m *MemoryConnector) AddJSON(data interface{}) (string, error) {
	jd, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return m.Add(jd)
}

// GetJSON gets data from IPFS and decodes it from JSON
func (m *MemoryConnector) GetJSON(path string, data interface{}) error {
	jd, err := m.Get(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jd, data); err != nil {
		return errors.Wrap(ErrDocumentNotValid, err.Error())
	}

	return nil
}

// GetID returns the local id
func (m *MemoryConnector) GetID() string {
	return ""
}

func genIPFSID(data []byte) (string, error) {
	mh, err := multihash.Sum(data, multihash.SHA2_256, 32)
	if err != nil {
		return "", err
	}

	return mh.B58String(), nil
}
