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
Package api - service integration and contract types
*/
package api

import (
	"context"
	"net/http"

	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/libs/transport"
)

var (
	apiVersion = "v1"
)

// ClientService interface
type ClientService interface {
	Order(token string, req *OrderRequest) (*OrderResponse, error)
	OrderSecret(token string, req *OrderSecretRequest) (*OrderSecretResponse, error)
	Status(token string) (*StatusResponse, error)
}

// MilagroClientService - implements Service Interface
type MilagroClientService struct {
	endpoints transport.ClientEndpoints
}

// ClientEndpoints return only the exported endpoints
func ClientEndpoints() transport.HTTPEndpoints {
	return transport.HTTPEndpoints{
		"Order": {
			Path:        "/" + apiVersion + "/order",
			Method:      http.MethodPost,
			NewRequest:  func() interface{} { return &OrderRequest{} },
			NewResponse: func() interface{} { return &OrderResponse{} },
		},
		"OrderSecret": {
			Path:        "/" + apiVersion + "/order/secret",
			Method:      http.MethodPost,
			NewRequest:  func() interface{} { return &OrderSecretRequest{} },
			NewResponse: func() interface{} { return &OrderSecretResponse{} },
		},
		"Status": {
			Path:        "/" + apiVersion + "/status",
			Method:      http.MethodGet,
			NewResponse: func() interface{} { return &StatusResponse{} },
		},
	}
}

// NewHTTPClient returns Service backed by an HTTP server living at the remote instance
func NewHTTPClient(instance string, logger *logger.Logger) (ClientService, error) {
	clientEndpoints, err := transport.NewHTTPClient(instance, ClientEndpoints(), logger)
	return MilagroClientService{clientEndpoints}, err

}

// Order makes a request for a new order
func (c MilagroClientService) Order(token string, req *OrderRequest) (*OrderResponse, error) {
	endpoint := c.endpoints["Order"]
	ctx := context.Background()
	ctx = transport.SetJWTAuthHeader(ctx, token)

	resp, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*OrderResponse), nil
}

// OrderSecret makes a request for initiate the order secret
func (c MilagroClientService) OrderSecret(token string, req *OrderSecretRequest) (*OrderSecretResponse, error) {
	endpoint := c.endpoints["OrderSecret"]
	ctx := context.Background()
	ctx = transport.SetJWTAuthHeader(ctx, token)

	resp, err := endpoint(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.(*OrderSecretResponse), nil
}

//Status - Allows a client to see the status of the server that it is connecting too
func (c MilagroClientService) Status(token string) (*StatusResponse, error) {
	endpoint := c.endpoints["Status"]
	ctx := context.Background()
	ctx = transport.SetJWTAuthHeader(ctx, token)

	resp, err := endpoint(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp.(*StatusResponse), nil
}
