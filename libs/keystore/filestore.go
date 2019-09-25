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
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/pkg/errors"
)

var (
	// ErrKeyNotFound is returned when a key is not found in the store
	ErrKeyNotFound = errors.New("Key not found")
)

// FileStore is the key Store implementation storing the keys in a file
type FileStore struct {
	sync.RWMutex
	filePath string
	keys     map[string][]byte
}

// NewFileStore creates a new FileStore
func NewFileStore(filePath string) (Store, error) {
	fs := &FileStore{
		filePath: filePath,
		keys:     map[string][]byte{},
	}

	if err := fs.loadKeys(); err != nil {
		return nil, err
	}

	return fs, nil
}

// Set stores multiple keys at once
func (f *FileStore) Set(keys ...Key) error {
	for _, key := range keys {
		f.keys[key.Name] = make([]byte, len(key.Key))
		copy(f.keys[key.Name], key.Key)
	}

	return f.storeKeys()
}

// Get retrieves multiple keys
func (f *FileStore) Get(names ...string) (map[string][]byte, error) {
	keys := map[string][]byte{}
	for _, name := range names {
		k, ok := f.keys[name]
		if !ok {
			return nil, ErrKeyNotFound
		}
		keys[name] = make([]byte, len(k))
		copy(keys[name], k)
	}
	return keys, nil
}

// TODO: Lock the file

func (f *FileStore) loadKeys() error {
	f.RLock()
	defer f.RUnlock()
	rawKeys, err := ioutil.ReadFile(f.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.Wrap(err, "Load keys")
	}

	return json.Unmarshal(rawKeys, &(f).keys)
}

func (f *FileStore) storeKeys() error {
	f.Lock()
	defer f.Unlock()
	rawKeys, err := json.Marshal(f.keys)
	if err != nil {
		return err
	}

	// Get the file permissions
	var perm os.FileMode
	fi, err := os.Stat(f.filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "Get key file permissions")
		}
		perm = 0600
	} else {
		perm = fi.Mode().Perm()
	}

	if err := ioutil.WriteFile(f.filePath, rawKeys, perm); err != nil {
		return errors.Wrap(err, "Store keys")
	}

	return nil
}
