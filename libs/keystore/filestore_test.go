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
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileStore(t *testing.T) {
	keys := map[string][]byte{"key1": {1}, "key2": {1, 2}, "key3": {1, 2, 3}}

	fn := tmpFileName()
	defer func() {
		if err := os.Remove(fn); err != nil {
			t.Logf("Warning! Temp file could not be deleted (%v): %v", err, fn)
		}
	}()

	fs, err := NewFileStore(fn)
	if err != nil {
		t.Fatal(err)
	}

	// Set keys
	for k, v := range keys {
		fs.Set(k, v)
	}

	// Get Keys
	for name, v := range keys {
		key, err := fs.Get(name)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(v, key) {
			t.Errorf("Key not match: %v. Expected: %v, Found: %v", name, v, key)
		}

	}

	fs1, err := NewFileStore(fn)
	if err != nil {
		t.Fatal(err)
	}
	// Get Keys
	for name, v := range keys {
		key, err := fs1.Get(name)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(v, key) {
			t.Errorf("Key not match: %v. Expected: %v, Found: %v", name, v, key)
		}

	}
}

func tmpFileName() string {
	rnd := make([]byte, 8)
	rand.Read(rnd)
	ts := time.Now().UnixNano()

	return filepath.Join(os.TempDir(), fmt.Sprintf("keystore-%v-%x.tmp", ts, rnd))
}
