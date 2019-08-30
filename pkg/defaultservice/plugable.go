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

package defaultservice

import (
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/api"
)

// Plugable service methods
type Plugable interface {
	// service
	Name() string
	Vendor() string

	// order
	ValidateOrderRequest(req *api.OrderRequest) error
	ValidateOrderSecretRequest(req *api.OrderSecretRequest, order documents.OrderDoc) error
	PrepareOrderPart1(order *documents.OrderDoc, reqExtension map[string]string) (fulfillExtension map[string]string, err error)
	PrepareOrderResponse(orderPart2 *documents.OrderDoc, reqExtension, fulfillExtension map[string]string) (commitment string, extension map[string]string, err error)
	ProduceBeneficiaryEncryptedData(blsSK []byte, order *documents.OrderDoc, req *api.OrderSecretRequest) (encrypted []byte, extension map[string]string, err error)
	ProduceFinalSecret(seed, sikeSK []byte, order, orderPart4 *documents.OrderDoc, req *api.OrderSecretRequest, fulfillSecretRespomse *api.FulfillOrderSecretResponse) (secret, commitment string, extension map[string]string, err error)
}
