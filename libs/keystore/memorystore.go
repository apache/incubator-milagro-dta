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

package keystore

import (
	"sync"
)

// MemoryStore is the in-memory implementation of key store
type MemoryStore struct {
	sync.RWMutex
	keys map[string][]byte
}

// NewMemoryStore creates a new MemoryStore
func NewMemoryStore() (Store, error) {
	return &MemoryStore{
		keys: map[string][]byte{},
	}, nil
}

// Set stores multiple keys at once
func (f *MemoryStore) Set(name string, key []byte) error {
	f.Lock()
	defer f.Unlock()

	f.keys[name] = make([]byte, len(key))
	copy(f.keys[name], key)

	return nil
}

// Get retrieves multiple keys
func (f *MemoryStore) Get(name string) ([]byte, error) {
	f.RLock()
	defer f.RUnlock()

	key, ok := f.keys[name]
	if !ok {
		return nil, ErrKeyNotFound
	}

	return key, nil
}
