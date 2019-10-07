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

//Package classification Milagro Custody Node API
//
//This application creates a distributed network of nodes that collaborate to keep secrets safe
//swagger:meta

/*
Package endpoints - HTTP API mapping
*/
package endpoints

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/libs/transport"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	apiVersion = "v1"
)

// Endpoints returns all the exported endpoints
func Endpoints(svc service.Service, corsAllow string, authorizer transport.Authorizer, logger *logger.Logger, nodeType string, pluginEndpoints service.Endpoints) transport.HTTPEndpoints {
	principalEndpoints := transport.HTTPEndpoints{

		"Order1": {
			Path:        "/" + apiVersion + "/order1",
			Method:      http.MethodPost,
			Endpoint:    MakeOrder1Endpoint(svc, logger),
			NewRequest:  func() interface{} { return &api.OrderRequest{} },
			NewResponse: func() interface{} { return &api.OrderResponse{} },
			Options: transport.ServerOptions(
				transport.SetCors(corsAllow),
				transport.AuthorizeOIDC(authorizer, false),
			),
			ErrStatus: transport.ErrorStatus{
				transport.ErrInvalidRequest: http.StatusUnprocessableEntity,
			},
		},

		"GetOrder": {
			Path:        "/" + apiVersion + "/order/{OrderReference}",
			Method:      http.MethodGet,
			Endpoint:    MakeGetOrderEndpoint(svc, logger),
			NewResponse: func() interface{} { return &api.GetOrderResponse{} },
			Options: transport.ServerOptions(
				transport.SetCors(corsAllow),
				transport.AuthorizeOIDC(authorizer, false),
			),
			ErrStatus: transport.ErrorStatus{
				transport.ErrInvalidRequest: http.StatusUnprocessableEntity,
			},
		},
		"OrderList": {
			Path:        "/" + apiVersion + "/order",
			Method:      http.MethodGet,
			Endpoint:    MakeOrderListEndpoint(svc, logger),
			NewResponse: func() interface{} { return &api.OrderListResponse{} },
			Options: transport.ServerOptions(
				transport.SetCors(corsAllow),
				transport.AuthorizeOIDC(authorizer, false),
			),
			ErrStatus: transport.ErrorStatus{
				transport.ErrInvalidRequest: http.StatusUnprocessableEntity,
			},
		},
		"OrderSecret1": {
			Path:        "/" + apiVersion + "/order/secret1",
			Method:      http.MethodPost,
			Endpoint:    MakeOrderSecret1Endpoint(svc, logger),
			NewRequest:  func() interface{} { return &api.OrderSecretRequest{} },
			NewResponse: func() interface{} { return &api.OrderSecretResponse{} },
			Options: transport.ServerOptions(
				transport.SetCors(corsAllow),
				transport.AuthorizeOIDC(authorizer, false),
			),
			ErrStatus: transport.ErrorStatus{
				transport.ErrInvalidRequest: http.StatusUnprocessableEntity,
			},
		},
	}
	masterFiduciaryEndpoints := transport.HTTPEndpoints{}

	statusEndPoints := transport.HTTPEndpoints{
		"Status": {
			Path:        "/" + apiVersion + "/status",
			Method:      http.MethodGet,
			Endpoint:    MakeStatusEndpoint(svc, logger, nodeType),
			NewResponse: func() interface{} { return &api.StatusResponse{} },
			Options: transport.ServerOptions(
				transport.SetCors(corsAllow),
				transport.AuthorizeOIDC(authorizer, false),
			),
			ErrStatus: transport.ErrorStatus{
				transport.ErrInvalidRequest: http.StatusUnprocessableEntity,
			},
		},
	}

	endpoints := transport.HTTPEndpoints{}
	switch strings.ToLower(nodeType) {
	case "multi":
		endpoints = concatEndpoints(masterFiduciaryEndpoints, principalEndpoints, statusEndPoints)
	case "principal":
		endpoints = concatEndpoints(principalEndpoints, statusEndPoints)
	case "fiduciary", "masterfiduciary":
		endpoints = concatEndpoints(masterFiduciaryEndpoints, statusEndPoints)
	}

	plugNamespace, plugEndpoints := pluginEndpoints.Endpoints()
	endpoints = concatPluginEndpoints(logger, endpoints, plugNamespace, plugEndpoints)

	return endpoints
}

func concatEndpoints(endpoints ...transport.HTTPEndpoints) transport.HTTPEndpoints {
	var res = make(transport.HTTPEndpoints)
	for _, endpoint := range endpoints {
		for k, v := range endpoint {
			res[k] = v
		}
	}
	return res
}

func concatPluginEndpoints(logger *logger.Logger, dst transport.HTTPEndpoints, namespace string, endpoints ...transport.HTTPEndpoints) transport.HTTPEndpoints {
	for _, endpoint := range endpoints {
		for k, v := range endpoint {
			v.Path = "/" + apiVersion + "/ext/" + namespace + v.Path
			logger.Info("Registering plugin endpoint %v", v.Path)
			dst["namespace."+k] = v
		}
	}
	return dst
}

//MakeOrderListEndpoint -
func MakeOrderListEndpoint(m service.Service, log *logger.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		params := transport.GetParams(ctx)
		sortBy := params.Get("sortBy")
		perPage, err := strconv.Atoi(params.Get("perPage"))
		if err != nil {
			return nil, err
		}
		page, err := strconv.Atoi(params.Get("page"))
		if err != nil {
			return nil, err
		}
		req := &api.OrderListRequest{
			Page:    page,
			PerPage: perPage,
			SortBy:  sortBy,
		}
		if err := validateRequest(req); err != nil {
			return "", err
		}
		debugRequest(ctx, req, log)
		return m.OrderList(req)
	}
}

//MakeGetOrderEndpoint -
func MakeGetOrderEndpoint(m service.Service, log *logger.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		params := transport.GetURLParams(ctx)
		orderReference := params.Get("OrderReference")

		req := &api.GetOrderRequest{
			OrderReference: orderReference,
		}
		debugRequest(ctx, req, log)
		return m.GetOrder(req)
	}
}

//MakeOrder1Endpoint -
func MakeOrder1Endpoint(m service.Service, log *logger.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*api.OrderRequest)
		if !ok {
			return nil, transport.ErrInvalidRequest
		}
		if err := validateRequest(req); err != nil {
			return "", err
		}
		debugRequest(ctx, req, log)
		return m.Order1(req)
	}
}

//MakeOrderSecret1Endpoint -
func MakeOrderSecret1Endpoint(m service.Service, log *logger.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*api.OrderSecretRequest)
		if !ok {
			return nil, transport.ErrInvalidRequest
		}
		if err := validateRequest(req); err != nil {
			return "", err
		}
		debugRequest(ctx, req, log)
		return m.OrderSecret1(req)
	}
}

//MakeStatusEndpoint -
func MakeStatusEndpoint(m service.Service, log *logger.Logger, nodeType string) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return m.Status(apiVersion, nodeType)
	}
}

func validateRequest(req interface{}) error {
	validate := validator.New()
	validate.RegisterAlias("hashID", "min=64,max=64")
	if err := validate.Struct(req); err != nil {
		return errors.Wrap(transport.ErrInvalidRequest, err.Error())
	}
	return nil
}

func debugRequest(ctx context.Context, req interface{}, log *logger.Logger) {
	log.Debug("%s : %+v", transport.ReqID(ctx), req)
}
