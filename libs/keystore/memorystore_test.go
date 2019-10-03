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
	"testing"
)

func TestMemoryStore(t *testing.T) {
	keys := map[string][]byte{"key1": {1}, "key2": {1, 2}, "key3": {1, 2, 3}}

	ms, err := NewMemoryStore()
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range keys {
		ms.Set(k, v)
	}

	for name, v := range keys {
		key, err := ms.Get(name)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(v, key) {
			t.Errorf("Key not match: %v. Expected: %v, Found: %v", name, v, key)
		}

	}
}
