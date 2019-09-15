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

package transport

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	oidc "github.com/coreos/go-oidc"
	"github.com/pkg/errors"
)

const (
//	jwtRequestTimeout = time.Second * 2
)

// Authorizer interface for implementing API Authorization
type Authorizer interface {
	Authorize(header http.Header) (*UserClaims, error)
}

// OIDCAuthorizer implements OIDC/OAuth2 JWT token authorization token authorizer
type OIDCAuthorizer struct {
	ClientID        string
	RequestUserInfo bool
	ClaimsSupported []string
	Provider        *oidc.Provider
}

// NewOIDCAuthorizer creates a new instance of OIDCAuthorizer
func NewOIDCAuthorizer(clientID, oidcProvider string, claims ...string) (Authorizer, error) {
	provider, err := oidc.NewProvider(context.Background(), oidcProvider)
	if err != nil {
		return nil, errors.Wrap(err, "init oidc provider")
	}

	return &OIDCAuthorizer{
		ClientID:        clientID,
		RequestUserInfo: true,
		ClaimsSupported: claims[:],
		Provider:        provider,
	}, nil
}

// Authorize checks the IDToken
func (a *OIDCAuthorizer) Authorize(header http.Header) (*UserClaims, error) {
	rawIDToken, err := extractJWTToken(header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	verifier := a.Provider.Verifier(&oidc.Config{ClientID: a.ClientID})
	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return nil, err
	}
	uc := &UserClaims{}
	if err := idToken.Claims(uc); err != nil {
		return nil, errors.Wrap(err, "parse claims")
	}

	return uc, nil
}

// UserClaims holds information about the claims from the UserInfo endpoint
type UserClaims map[string]interface{}

// GetString returns a value of a key if exists
func (uc *UserClaims) GetString(key string) string {
	if uc == nil {
		return ""
	}
	vi, ok := (*uc)[key]
	if ok {
		v, ok := vi.(string)
		if ok {
			return v
		}
	}
	return ""
}

// LocalAuthorizer implements JWT Authorizer without verifying the signature
// For testing purposes only
type LocalAuthorizer struct{}

// Authorize checks the IDToken
func (a *LocalAuthorizer) Authorize(header http.Header) (*UserClaims, error) {
	rawIDToken, err := extractJWTToken(header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	uc := &UserClaims{}
	if err := parseJWTToken(rawIDToken, uc); err != nil {
		return nil, err
	}

	return uc, nil
}

func extractJWTToken(authHeader string) (string, error) {
	split := strings.Split(authHeader, " ")
	if len(split) < 2 || split[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}
	return split[1], nil
}

func parseJWTToken(p string, v interface{}) error {
	parts := strings.Split(p, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return errors.Wrap(err, "invalid token")
	}

	if err := json.Unmarshal(payload, v); err != nil {
		return errors.Wrap(err, "invalid token payload")
	}

	return nil
}

// SetJWTAuthHeader sets the Bearer Authorization header to the http client context
func SetJWTAuthHeader(ctx context.Context, jwtToken string) context.Context {
	h := http.Header{}
	h.Set("Authorization", "Bearer "+jwtToken)
	return AddClientHeader(ctx, h)
}

// EmptyAuthorizer implements empty Authorizer
type EmptyAuthorizer struct{}

// Authorize checks the IDToken
func (a *EmptyAuthorizer) Authorize(header http.Header) (*UserClaims, error) {
	return &UserClaims{}, nil
}
