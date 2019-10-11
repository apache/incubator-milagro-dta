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
Package safeguardsecret - is an example of a D-TA plugin
*/
package safeguardsecret

import (
	"errors"
	"fmt"
	"strings"

	"github.com/apache/incubator-milagro-dta/libs/crypto"
	"github.com/apache/incubator-milagro-dta/libs/cryptowallet"
	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/pkg/defaultservice"
)

//Constants describe plugin and its creator
var (
	extensionVendor = "Milagro"
	pluginName      = "safeguardsecret"
)

// Service implements Safeguard secret plugin service
type Service struct {
	defaultservice.Service
}

//NewService returns a new Safeguard secret implementation of Service
func NewService() *Service {
	return &Service{}
}

// Name of the plugin
func (s *Service) Name() string {
	return pluginName
}

// Vendor of the plugin
func (s *Service) Vendor() string {
	return extensionVendor
}

// PrepareOrderResponse gets the updated order and returns the commitment and extension
func (s *Service) PrepareOrderResponse(order *documents.OrderDoc) (commitment string, extension map[string]string, err error) {
	finalPublicKey := order.OrderPart2.CommitmentPublicKey
	plainText := order.OrderReqExtension["plainText"]
	delete(order.OrderReqExtension, "plainText") //we don't want to store this
	c, v, t, err := crypto.Secp256k1Encrypt(plainText, finalPublicKey)
	CiphertextWithParams := fmt.Sprintf("%s:%s:%s", t, v, c)
	return finalPublicKey, map[string]string{"CipherText": CiphertextWithParams}, nil
	//	return finalPublicKey, map[string]string{"Ciphertext": c, "v": v, "t": t}, nil
}

// ProduceFinalSecret -
func (s *Service) ProduceFinalSecret(seed, sikeSK []byte, order, orderPart4 *documents.OrderDoc, beneficiaryCID string) (secret, commitment string, extension map[string]string, err error) {
	finalPrivateKey := orderPart4.OrderDocument.OrderPart4.Secret
	//Derive the Public key from the supplied Private Key
	finalPublicKey, _, err := cryptowallet.PublicKeyFromPrivate(finalPrivateKey)
	if err != nil {
		return "", "", nil, err
	}
	cypherTextWithParams := order.OrderPart2.Extension["CipherText"]
	cypherparts := strings.SplitN(cypherTextWithParams, ":", 3) //maps parts to t:v:c
	if len(cypherparts) != 3 {
		return "", "", nil, errors.New("Invalid Ciphertext")
	}
	plainText, err := crypto.Secp256k1Decrypt(cypherparts[2], cypherparts[1], cypherparts[0], finalPrivateKey)
	if err != nil {
		return "", "", nil, err
	}
	return finalPrivateKey, finalPublicKey, map[string]string{"PlainText": plainText}, nil
}
