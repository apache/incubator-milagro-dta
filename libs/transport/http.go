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
Package transport - HTTP request and response methods
*/
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// TODO: Add request rate limiter
// TODO: Add circuit breaker

type contextKey int

const (
	contextKeyCorsOrigin     contextKey = 10000
	contextQueryParams       contextKey = 10001
	contextURLParams         contextKey = 10002
	contextHTTPHeaders       contextKey = 10003
	contextProtectedEndpoint contextKey = 10004
	contextAuthorized        contextKey = 10005
	contextAuthorizeError    contextKey = 10006
	contextUserInfo          contextKey = 10007
)

var (
	// ErrMethodNotAllowed when HTTP request method is not handled
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrInvalidRequest when HTTP request is invalid
	ErrInvalidRequest = errors.New("invalid request")
	// ErrInvalidJSON when the request is invalid json
	ErrInvalidJSON = errors.New("invalid json")
	// ErrUnauthorized when a protected endpoint is not authorized
	ErrUnauthorized = errors.New("unauthorized")
	// ErrHTTPClientError when HTTP request status is not 200
	ErrHTTPClientError = errors.New("request error")
)

// HTTPEndpoints is a map of endpoints
// Key is the name of the endpoint
type HTTPEndpoints map[string]HTTPEndpoint

// HTTPEndpoint defines a single endpoint
type HTTPEndpoint struct {
	Path        string
	Method      string
	Endpoint    endpoint.Endpoint
	NewRequest  func() interface{}
	NewResponse func() interface{}
	ErrStatus   ErrorStatus
	Options     []httptransport.ServerOption
}

// ClientEndpoints is a map of all exported client endpoints
type ClientEndpoints map[string]endpoint.Endpoint

// ErrorStatus is a map of errors to http response status codes
type ErrorStatus map[error]int

// ResponseStatus returns the status code and status text based on the error type
func (e ErrorStatus) ResponseStatus(err error) (statusCode int, statusText string) {
	statusCode = http.StatusOK
	statusText = ""

	if err != nil {
		statusCode = http.StatusInternalServerError

		if sc, ok := e[errors.Cause(err)]; ok {
			statusCode = sc
			statusText = decodeError(err.Error())
		}

		// Add transport errors
		switch errors.Cause(err) {
		case ErrInvalidJSON, ErrInvalidRequest:
			statusCode = http.StatusBadRequest
		case ErrMethodNotAllowed:
			statusCode = http.StatusMethodNotAllowed
		case ErrUnauthorized:
			statusCode = http.StatusUnauthorized
			statusText = decodeError(err.Error())
		case ErrHTTPClientError:
			statusCode = http.StatusBadRequest
			statusText = decodeError(err.Error())
		}
	}

	if statusText == "" {
		statusText = http.StatusText(statusCode)
	}

	return statusCode, statusText
}

func decodeError(strerr string) string {
	ew := &errorWrapper{}
	if er := json.Unmarshal([]byte(strerr), ew); er != nil {
		return strerr
	}
	return ew.Error
}

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
func NewHTTPHandler(endpoints HTTPEndpoints, logger *logger.Logger, duration metrics.Histogram) *mux.Router {
	m := mux.NewRouter()
	for eName, e := range endpoints {
		endpoint := loggingMiddleware(
			logger,
			e.ErrStatus,
		)(e.Endpoint)
		endpoint = metricsDurationMiddleware(duration.With("method", eName))(endpoint)

		options := []httptransport.ServerOption{
			httptransport.ServerErrorEncoder(errorEncoder(e.ErrStatus, logger)),
			parseQueryParams,
			parseURLParams,
			parseHTTPHeaders,
			httptransport.ServerBefore(httptransport.PopulateRequestContext),
		}

		if e.Options != nil {
			options = append(options, e.Options...)
		}

		m.Path(e.Path).Methods(http.MethodOptions).Handler(
			httptransport.NewServer(
				func(ctx context.Context, request interface{}) (response interface{}, err error) {
					return nil, nil
				},
				decodeOptionsRequest,
				func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
					return nil
				},
				options...,
			),
		)

		m.Path(e.Path).Methods(e.Method).Handler(
			httptransport.NewServer(
				endpoint,
				decodeJSONRequest(e),
				encodeJSONResponse,
				options...,
			),
		)
	}

	return m
}

// GetParams returns Query params of a request
func GetParams(ctx context.Context) url.Values {
	v := ctx.Value(contextQueryParams)
	if v == nil {
		return url.Values{}
	}

	values, ok := v.(url.Values)
	if !ok {
		return url.Values{}
	}

	return values
}

// GetURLParams returns URL params of the request path
func GetURLParams(ctx context.Context) url.Values {
	v := ctx.Value(contextURLParams)
	if v == nil {
		return url.Values{}
	}

	values, ok := v.(url.Values)
	if !ok {
		return url.Values{}
	}

	return values
}

// GetHeaders returns the HTTP Request headers
func GetHeaders(ctx context.Context) http.Header {
	v := ctx.Value(contextHTTPHeaders)
	if v == nil {
		return http.Header{}
	}

	headers, ok := v.(http.Header)
	if !ok {
		return http.Header{}
	}

	return headers
}

// GetUserInfo returns the authorized user info
func GetUserInfo(ctx context.Context) *UserClaims {
	v := ctx.Value(contextUserInfo)
	if v == nil {
		return &UserClaims{}
	}

	uc, ok := v.(*UserClaims)
	if !ok {
		return &UserClaims{}
	}

	return uc
}

// Options helper handlers

// ServerOptions returns an array of httptransport.ServerOption
func ServerOptions(opt ...httptransport.ServerOption) []httptransport.ServerOption {
	return opt
}

// SetCors sets CORS HTTP response headers for specific domains
func SetCors(domain string) httptransport.ServerOption {
	return func(s *httptransport.Server) {
		httptransport.ServerBefore(
			func(ctx context.Context, r *http.Request) context.Context {
				origin := r.Header.Get("Origin")
				if domain != "*" && origin != domain {
					return ctx
				}
				return context.WithValue(ctx, contextKeyCorsOrigin, origin)
			},
		)(s)

		httptransport.ServerAfter(
			func(ctx context.Context, w http.ResponseWriter) context.Context {
				setCORSHeaders(ctx, w)
				return ctx
			},
		)(s)
	}
}

func setCORSHeaders(ctx context.Context, w http.ResponseWriter) {
	cv := ctx.Value(contextKeyCorsOrigin)
	if cv == nil {
		return
	}
	origin, ok := cv.(string)
	if !ok {
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// AuthorizeOIDC checks the Authorization header OIDC token, verifies it and gets the user info
func AuthorizeOIDC(authorizer Authorizer, allowUnauthorized bool) httptransport.ServerOption {
	return func(s *httptransport.Server) {
		httptransport.ServerBefore(
			func(ctx context.Context, r *http.Request) context.Context {
				ctx = context.WithValue(ctx, contextProtectedEndpoint, !allowUnauthorized)

				userInfo, err := authorizer.Authorize(r.Header)
				if err != nil {
					ctx = context.WithValue(ctx, contextAuthorized, false)
					ctx = context.WithValue(ctx, contextAuthorizeError, err.Error())
					return ctx
				}
				ctx = context.WithValue(ctx, contextAuthorized, true)
				ctx = context.WithValue(ctx, contextUserInfo, userInfo)
				return ctx
			},
		)(s)
	}
}

// parseQueryParams and store them in context
func parseQueryParams(s *httptransport.Server) {
	httptransport.ServerBefore(
		func(ctx context.Context, r *http.Request) context.Context {
			return context.WithValue(ctx, contextQueryParams, r.URL.Query())
		},
	)(s)
}

// parseURLParams and store them in context
func parseURLParams(s *httptransport.Server) {
	httptransport.ServerBefore(
		func(ctx context.Context, r *http.Request) context.Context {
			muxV := mux.Vars(r)
			qp := url.Values{}
			for k, v := range muxV {
				qp.Set(k, v)
			}
			return context.WithValue(ctx, contextURLParams, qp)
		},
	)(s)
}

// parseHTTPHeaders and store them in context
func parseHTTPHeaders(s *httptransport.Server) {
	httptransport.ServerBefore(
		func(ctx context.Context, r *http.Request) context.Context {
			return context.WithValue(ctx, contextHTTPHeaders, r.Header)
		},
	)(s)
}

func decodeJSONRequest(e HTTPEndpoint) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		if e.Method != r.Method {
			return nil, ErrMethodNotAllowed
		}

		// Check authorization
		if ctx.Value(contextProtectedEndpoint) == true {
			if ctx.Value(contextAuthorized) == false {
				errString, ok := ctx.Value(contextAuthorizeError).(string)
				if ok {
					return nil, errors.Wrap(ErrUnauthorized, errString)
				}
				return nil, ErrUnauthorized
			}
		}

		// Decode request only if it's expected
		if e.NewRequest == nil {
			return nil, nil
		}
		req := e.NewRequest()
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			return nil, errors.Wrap(ErrInvalidJSON, err.Error())
		}

		return req, nil
	}
}

func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if response == nil {
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeOptionsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	corsEnabled := ctx.Value(contextKeyCorsOrigin)
	if corsEnabled == nil {
		return nil, ErrMethodNotAllowed
	}

	return nil, nil
}

func errorEncoder(errorStatus ErrorStatus, logger *logger.Logger) httptransport.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		statusCode, statusText := errorStatus.ResponseStatus(err)

		// _ = logger.Log("status", fmt.Sprintf("%d %s", statusCode, statusText))
		setCORSHeaders(ctx, w)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(errorWrapper{Error: statusText})
	}
}

func loggingMiddleware(logger *logger.Logger, errorStatus ErrorStatus) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			method := getContextStr(ctx, httptransport.ContextKeyRequestMethod)
			path := getContextStr(ctx, httptransport.ContextKeyRequestPath)
			reqID := getReqID(ctx)

			defer func(begin time.Time) {
				statusCode, statusText := errorStatus.ResponseStatus(err)
				if err != nil {
					if statusCode == http.StatusInternalServerError {
						logger.Error("reqID: %v, err: %v", reqID, err.Error())
					}
				}
				logger.Response(reqID, method, path, statusCode, statusText, time.Since(begin))
			}(time.Now())

			logger.Request(reqID, method, path)

			response, err = next(ctx, request)
			return
		}
	}
}

func getContextStr(ctx context.Context, key interface{}) string {
	if v := ctx.Value(key); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getReqID(ctx context.Context) string {
	// check if request ID is passed as HTTP header
	reqID := getContextStr(ctx, httptransport.ContextKeyRequestXRequestID)
	if reqID == "" {
		// generate request ID
		id := uuid.New()
		reqID = id.String()
	}
	return reqID
}

func metricsDurationMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				if duration != nil {
					duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
				}
			}(time.Now())
			return next(ctx, request)
		}
	}
}

type errorWrapper struct {
	Error string `json:"error"`
}

// HTTP Client

// addHeader adds the content of http.Header to an existing http.Header
// to ensure it keeps the headers that are already set
func addHeader(dst, src http.Header) {
	for k, v := range src {
		for _, hv := range v {
			dst.Add(k, hv)
		}
	}
}

func encodeJSONRequest(ctx context.Context, r *http.Request, request interface{}) error {
	headers, ok := ctx.Value(contextHTTPHeaders).(http.Header)
	if ok {
		addHeader(r.Header, headers)
	}

	if request == nil {
		return nil
	}

	if urlV, ok := request.(url.Values); ok {
		r.URL.RawQuery = urlV.Encode()
		return nil
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeJSONResponse(e HTTPEndpoint) httptransport.DecodeResponseFunc {
	return func(_ context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode != http.StatusOK {
			strBody, _ := ioutil.ReadAll(r.Body)
			defer r.Body.Close()

			return nil, errors.Wrapf(ErrHTTPClientError, "status: %v (%s)", r.Status, decodeError(string(strBody)))
		}

		if e.NewResponse == nil {
			return nil, nil
		}
		resp := e.NewResponse()
		err := json.NewDecoder(r.Body).Decode(resp)
		return resp, err
	}
}

// NewHTTPClient returns an HTTP handler that makes a set of endpoints
func NewHTTPClient(instance string, endpoints HTTPEndpoints, logger *logger.Logger) (ClientEndpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	ce := ClientEndpoints{}

	for ename, e := range endpoints {
		ce[ename] = httptransport.NewClient(
			e.Method,
			copyURL(u, e.Path),
			encodeJSONRequest,
			decodeJSONResponse(e),
		).Endpoint()
	}

	return ce, nil
}

// AddClientHeader adds header to the http client context
func AddClientHeader(ctx context.Context, header http.Header) context.Context {
	currentHeader, ok := ctx.Value(contextHTTPHeaders).(http.Header)
	if !ok {
		currentHeader = http.Header{}
	}

	addHeader(currentHeader, header)

	return context.WithValue(ctx, contextHTTPHeaders, currentHeader)
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
