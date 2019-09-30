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
Package service - defines core Milagro D-TA interface
*/
package service

import "github.com/apache/incubator-milagro-dta/pkg/api"

// Service is the CustodyService interface
type Service interface {
	//Identity
	CreateIdentity(req *api.CreateIdentityRequest) (*api.CreateIdentityResponse, error)
	GetIdentity(req *api.GetIdentityRequest) (*api.GetIdentityResponse, error)
	IdentityList(req *api.IdentityListRequest) (*api.IdentityListResponse, error)

	//Order
	GetOrder(req *api.GetOrderRequest) (*api.GetOrderResponse, error)
	OrderList(req *api.OrderListRequest) (*api.OrderListResponse, error)

	//Order processing
	OrderSecret1(req *api.OrderSecretRequest) (string, error)
	OrderSecret2(req *api.FulfillOrderSecretResponse) (string, error)

	Order1(req *api.OrderRequest) (string, error)
	Order2(tx *api.BlockChainTX) (string, error)

	//Fullfill processing
	BCFulfillOrder(tx *api.BlockChainTX) (string, error)
	FulfillOrder(req *api.FulfillOrderRequest) (string, error)
	FulfillOrderSecret(req *api.FulfillOrderSecretRequest) (string, error)

	NodeID() string
	MasterFiduciaryNodeID() string
	SetNodeID(nodeID string)
	SetMasterFiduciaryNodeID(masterFiduciaryNodeID string)

	//System
	Dump(tx *api.BlockChainTX) error //Decrypt and dump the order
	Status(apiVersion, nopdeType string) (*api.StatusResponse, error)
}
