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

package datastore

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	opSet = iota
	opGet
	opDel
	opListKeys
)

func TestBoltBackend(t *testing.T) {
	type tcOp struct {
		optype      int
		doctype     string
		key         string
		err         error
		result      interface{}
		listLimit   int
		listSkip    int
		listReverse bool
	}

	seedValues := []struct {
		doctype string
		key     string
		value   []byte
		index   map[string]string
	}{
		{"test", "3", []byte{3, 4, 5}, map[string]string{"index": "3"}},
		{"test", "1", []byte{1, 2, 3}, map[string]string{"index": "1"}},
		{"test", "4", []byte{4, 5}, map[string]string{"index": "4"}},
		{"test", "2", []byte{2, 3, 4}, map[string]string{"index": "2"}},
		{"test", "5", []byte{5, 6}, map[string]string{"index": "5"}},
	}

	testCases := []tcOp{
		{opGet, "test", "", ErrKeyNotFound, nil, 0, 0, false},
		{opGet, "", "1", ErrKeyNotFound, nil, 0, 0, false},
		{opGet, "test", "1", nil, []byte{1, 2, 3}, 0, 0, false},
		{opListKeys, "test", "index", nil, []string{"1", "2", "3", "4", "5"}, 0, 0, false},
		{opDel, "test", "", nil, nil, 0, 0, false},
		{opDel, "", "1", nil, nil, 0, 0, false},
		{opDel, "test", "1", nil, nil, 0, 0, false},
		{opDel, "test", "1", nil, nil, 0, 0, false},
		{opListKeys, "test", "index", nil, []string{"2", "3", "4", "5"}, 0, 0, false},
		{opListKeys, "test", "index", nil, []string{"2"}, 0, 1, false},
		{opListKeys, "test", "index", nil, []string{"3", "4"}, 1, 2, false},
		{opListKeys, "test", "index", nil, []string{"4", "5"}, 2, 2, false},
		{opListKeys, "test", "index", nil, []string{"5"}, 3, 2, false},
		{opListKeys, "test", "index", nil, []string{}, 4, 2, false},
		{opListKeys, "test", "index", nil, []string{}, 5, 0, false},
		{opListKeys, "test", "index", nil, []string{"2", "3", "4", "5"}, 0, 10, false},
		{opListKeys, "test", "index", nil, []string{"5", "4", "3", "2"}, 0, 0, true},
		{opListKeys, "test", "index", nil, []string{"4", "3", "2"}, 1, 0, true},
		{opListKeys, "test", "index", nil, []string{"5"}, 0, 1, true},
		{opListKeys, "test", "index", nil, []string{"4", "3"}, 1, 2, true},
		{opListKeys, "test", "index", nil, []string{"3", "2"}, 2, 2, true},
		{opListKeys, "test", "index", nil, []string{"2"}, 3, 2, true},
		{opListKeys, "test", "index", nil, []string{}, 4, 2, true},
		{opListKeys, "test", "index", nil, []string{}, 5, 0, true},
		{opListKeys, "test", "index", nil, []string{"5", "4", "3", "2"}, 0, 10, true},

		{opListKeys, "test-invalid", "index", nil, []string{}, 0, 5, false},
		{opListKeys, "test", "index-invalid", nil, []string{}, 0, 5, false},
	}

	dbName := genTempFilename()
	defer os.Remove(dbName)

	b, err := NewBoltBackend(dbName)
	if err != nil {
		t.Fatal(err)
	}

	for _, sv := range seedValues {
		if err := b.Set(sv.doctype, sv.key, sv.value, sv.index); err != nil {
			t.Fatal(err)
		}
	}

	for itc, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", itc), func(t *testing.T) {
			switch tc.optype {
			case opGet:
				result, err := b.Get(tc.doctype, tc.key)
				if err != tc.err {
					t.Fatalf("invalid Get error response. Expected: %v, found: %v", tc.err, err)
				}
				if err != nil {
					break
				}
				if !bytes.Equal(result, tc.result.([]byte)) {
					t.Fatalf("invalid Get result. Expected: %v, found: %v", tc.result, result)
				}
			case opDel:
				if err := b.Del(tc.doctype, tc.key); err != tc.err {
					t.Fatalf("invalid Del error response. Expected: %v, found: %v", tc.err, err)
				}
			case opListKeys:
				result, err := b.ListKeys(tc.doctype, tc.key, tc.listLimit, tc.listSkip, tc.listReverse)
				if err != tc.err {
					t.Fatalf("invalid ListKeys error response. Expected: %v, found: %v", tc.err, err)
				}
				if err != nil {
					break
				}

				if strings.Join(result, "") != strings.Join(tc.result.([]string), "") {
					t.Fatalf("invalid ListKeys result. Expected: %v, found: %v", tc.result, result)
				}
			}
		})
	}

	if err := b.Close(); err != nil {
		t.Fatalf("Close database: %v", err)
	}

}

func genTempFilename() string {
	tempDir := os.TempDir()
	filename := fmt.Sprintf("milagro-test-bolt-%v.db", time.Now().UnixNano())
	return filepath.Join(tempDir, filename)
}
