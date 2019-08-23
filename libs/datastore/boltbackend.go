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
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

// BoltBackend implements Backend interface with embedded Bolt storage
type BoltBackend struct {
	db *bolt.DB
}

// NewBoltBackend creates a new Bolt backend using CoreOS bbolt implementation
func NewBoltBackend(filename string) (Backend, error) {
	db, err := bolt.Open(
		filename,
		0600,
		&bolt.Options{
			Timeout: 1 * time.Second,
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "initialize bolt datastore backend")
	}

	return &BoltBackend{
		db: db,
	}, nil
}

// Set stores the value for a key of datatype using bolt datastore
func (bb *BoltBackend) Set(datatype, key string, value []byte, indexData map[string]string) error {
	return bb.db.Update(func(tx *bolt.Tx) error {
		// Get or create the root bucket for the datatype
		bk, err := tx.CreateBucketIfNotExists([]byte(datatype))
		if err != nil {
			return errors.Wrap(err, "failed to create root bucket")
		}

		// Get or crerate the data bucket for storing the key/values
		dataBk, err := bk.CreateBucketIfNotExists([]byte("data"))
		if err != nil {
			return errors.Wrap(err, "failed to create data bucket")
		}
		// Store data in the data bucket
		if err := dataBk.Put([]byte(key), value); err != nil {
			return errors.Wrap(err, "failed to put data")
		}

		// Delete the previous indexes if they exists
		if err := deleteIndexes(key, bk); err != nil {
			return errors.Wrap(err, "delete old indexes")
		}

		// Perform indexing
		if err := createIndexes(key, indexData, bk); err != nil {
			return errors.Wrap(err, "create indexes")
		}

		return nil
	})
}

// Get retreives the value for specified key and datatyoe
// Returns ErrKeyNotFound if the key has no value set
func (bb *BoltBackend) Get(datatype, key string) (data []byte, err error) {
	err = bb.db.View(func(tx *bolt.Tx) error {
		// Get the root bucket
		bk := tx.Bucket([]byte(datatype))
		if bk == nil {
			return ErrKeyNotFound
		}
		dataBk := bk.Bucket([]byte("data"))
		if dataBk == nil {
			return ErrKeyNotFound
		}

		data = dataBk.Get([]byte(key))
		if data == nil {
			return ErrKeyNotFound
		}

		return nil
	})

	return
}

// Del deletes a key and all the indexes
func (bb *BoltBackend) Del(datatype, key string) error {
	return bb.db.Update(func(tx *bolt.Tx) error {
		// Get the root bucket
		bk := tx.Bucket([]byte(datatype))
		if bk == nil {
			return nil
		}
		dataBk := bk.Bucket([]byte("data"))
		if dataBk == nil {
			return nil
		}

		if err := deleteIndexes(key, bk); err != nil {
			return errors.Wrap(err, "delete indexes")
		}

		return dataBk.Delete([]byte(key))
	})
}

type iterFunc func() ([]byte, []byte)

// ListKeys lists all keys for specified datatype
func (bb *BoltBackend) ListKeys(datatype, index string, skip, limit int, reverse bool) (keys []string, err error) {
	err = bb.db.View(func(tx *bolt.Tx) error {
		// Get the root bucket
		bk := tx.Bucket([]byte(datatype))
		if bk == nil {
			return nil
		}
		indexBk := bk.Bucket([]byte(fmt.Sprintf("index-%s", index)))
		if indexBk == nil {
			return nil
		}
		c := indexBk.Cursor()

		var first, next iterFunc
		switch reverse {
		default:
			first = c.First
			next = c.Next
		case true:
			first = c.Last
			next = c.Prev
		}

		i := 0
		for k, v := first(); k != nil; k, v = next() {
			i++
			if i <= skip {
				continue
			}
			keys = append(keys, string(v))
			if limit > 0 && len(keys) >= limit {
				break
			}
		}

		return nil
	})

	return
}

// Close closes the database
func (bb *BoltBackend) Close() error {
	return errors.Wrap(bb.db.Close(), "close bolt datastore backend database")
}

func createIndexes(key string, indexData map[string]string, rootBucket *bolt.Bucket) error {
	// Iterate over index values
	for indexName, v := range indexData {
		indexBk, err := rootBucket.CreateBucketIfNotExists([]byte(fmt.Sprintf("index-%s", indexName)))
		if err != nil {
			return errors.Wrapf(err, "failed to create index bucket for index %s", indexName)
		}

		// generate unique and sortable key
		valueKey := createIndexValueKey(v)
		if err := indexBk.Put([]byte(valueKey), []byte(key)); err != nil {
			return errors.Wrap(err, "failed to put index data")
		}

		// update the index data with the actual key
		indexData[indexName] = valueKey
	}

	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(&indexData); err != nil {
		return errors.Wrap(err, "failed to encode index data")
	}

	return rootBucket.Put([]byte(fmt.Sprintf("indexes-%s", key)), b.Bytes())
}

func deleteIndexes(key string, rootBucket *bolt.Bucket) error {
	kiData := rootBucket.Get([]byte(fmt.Sprintf("indexes-%s", key)))
	if kiData == nil {
		return nil
	}

	indexData := map[string]string{}
	b := bytes.NewBuffer(kiData)
	if err := gob.NewDecoder(b).Decode(&indexData); err != nil {
		return errors.Wrap(err, "invalid index data")
	}

	for indexName, v := range indexData {
		indexBk := rootBucket.Bucket([]byte(fmt.Sprintf("index-%s", indexName)))
		if indexBk == nil {
			return errors.Errorf("index not found: %s", indexName)
		}

		if err := indexBk.Delete([]byte(v)); err != nil {
			return nil
		}
	}

	return rootBucket.Delete([]byte(fmt.Sprintf("indexes-%s", key)))
}

func createIndexValueKey(v string) string {
	randBytes := make([]byte, 8)
	rand.Read(randBytes)
	return fmt.Sprintf("%s\x00%s%x", v, time.Now().UTC().Format(time.RFC3339), randBytes)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
