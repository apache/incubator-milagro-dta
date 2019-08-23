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

package crypto

/*
#include <amcl/randapi.h>
*/
import "C"
import (
	"encoding/hex"
	"unsafe"
)

// Octet adds functionality around C octet
type Octet = C.octet

// NewOctet creates an empty Octet with a given size
func NewOctet(maxSize int) *Octet {
	return &Octet{
		len: C.int(maxSize),
		max: C.int(maxSize),
		val: (*C.char)(C.calloc(1, C.size_t(maxSize))),
	}
}

// CreateOctet creates new Octet with a value
func CreateOctet(val []byte) *Octet {
	if val == nil {
		return nil
	}

	return &Octet{
		len: C.int(len(val)),
		max: C.int(len(val)),
		val: C.CString(string(val)),
	}
}

// Free frees the allocated memory
func (o *Octet) Free() {
	if o == nil {
		return
	}
	C.free(unsafe.Pointer(o.val))
}

// ToBytes returns the bytes representation of the Octet
func (o *Octet) ToBytes() []byte {
	return C.GoBytes(unsafe.Pointer(o.val), o.len)
}

// ToString returns the hex encoded representation of the Octet
func (o *Octet) ToString() string {
	return hex.EncodeToString(o.ToBytes())
}

// Rand is a cryptographically secure random number generator
type Rand C.csprng

// NewRand create new seeded Rand
func NewRand(seed []byte) *Rand {
	sOct := CreateOctet(seed)
	defer sOct.Free()

	var rand C.csprng
	C.CREATE_CSPRNG(&rand, sOct)
	return (*Rand)(&rand)
}
