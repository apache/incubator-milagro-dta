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
	"encoding/gob"
)

// GOBCodec implement Codec interface with GOB encoding
type GOBCodec struct{}

// NewGOBCodec creates a new GOB Codec
func NewGOBCodec() Codec {
	return &GOBCodec{}
}

// Marshal with GOB encoding
func (s *GOBCodec) Marshal(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unmarshal with GOB encoding
func (s *GOBCodec) Unmarshal(b []byte, v interface{}) error {
	buf := bytes.NewBuffer(b)
	if err := gob.NewDecoder(buf).Decode(v); err != nil {
		return err
	}

	return nil
}
