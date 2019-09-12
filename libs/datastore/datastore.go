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

/*
Package datastore - enables data to be persisted in built in datebase (Bolt)
*/
package datastore

import (
	"errors"
)

var (
	// ErrBackendNotInitialized is returned when the backend is not set
	ErrBackendNotInitialized = errors.New("backend not initialized")
	// ErrCodecNotInitialized is returned when the codec is not set
	ErrCodecNotInitialized = errors.New("codec not initialized")
	// ErrKeyNotFound is returned when an attempt to load a value of a missing key is made
	ErrKeyNotFound = errors.New("key not found")
)

// Store provides key-value data storage with
// specific backend and codec
type Store struct {
	backend Backend
	codec   Codec
}

// NewStore constructs a new store
func NewStore(options ...StoreOption) (*Store, error) {
	s := &Store{}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Set stores the value for a key of datatype
func (s *Store) Set(datatype, key string, v interface{}, indexData map[string]string) error {
	if err := s.checkInit(); err != nil {
		return err
	}

	vbytes, err := s.codec.Marshal(v)
	if err != nil {
		return err
	}

	if err := s.backend.Set(datatype, key, vbytes, indexData); err != nil {
		return err
	}

	return nil
}

// Get retreives the value for specified key and datatyoe
// Returns ErrKeyNotFound if the key has no value set
func (s *Store) Get(datatype, key string, v interface{}) error {
	if err := s.checkInit(); err != nil {
		return err
	}

	vbytes, err := s.backend.Get(datatype, key)
	if err != nil {
		return err
	}

	if err := s.codec.Unmarshal(vbytes, v); err != nil {
		return err
	}

	return nil
}

// Del deletes a key and all the indexes
func (s *Store) Del(datatype, key string) error {
	if err := s.checkInit(); err != nil {
		return err
	}

	return s.backend.Del(datatype, key)
}

// Close closes the database
func (s *Store) Close() error {
	return s.backend.Close()
}

// ListKeys lists keys by index
func (s *Store) ListKeys(datatype, index string, skip, limit int, reverse bool) (keys []string, err error) {
	return s.backend.ListKeys(datatype, index, skip, limit, reverse)
}

func (s *Store) checkInit() error {
	if s.backend == nil {
		return ErrBackendNotInitialized
	}
	if s.codec == nil {
		return ErrBackendNotInitialized
	}
	return nil
}

// StoreOption sets additional parameters to the Store
type StoreOption func(s *Store) error

// WithBackend sets the store backend
func WithBackend(b Backend) StoreOption {
	return func(s *Store) error {
		s.backend = b
		return nil
	}
}

// WithCodec sets the store codec
func WithCodec(enc Codec) StoreOption {
	return func(s *Store) error {
		s.codec = enc
		return nil
	}
}

// Backend provides data storage interface
type Backend interface {
	Set(datatype, key string, value []byte, indexData map[string]string) error
	Get(datatype, key string) ([]byte, error)
	Del(datatype, key string) error
	ListKeys(datatype, index string, skip, limit int, reverse bool) (keys []string, err error)
	Close() error
}

// Codec probides data serialization interface
type Codec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}
