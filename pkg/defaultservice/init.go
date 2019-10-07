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
	"io"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/ipfs"
	"github.com/apache/incubator-milagro-dta/libs/keystore"
	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/pkg/config"
	"github.com/apache/incubator-milagro-dta/pkg/tendermint"
)

// ServiceOption function to set Service properties
type ServiceOption func(s *Service) error

// Init sets-up the service options. It's called when the plugin gets registered
func (s *Service) Init(plugin Plugable, options ...ServiceOption) error {
	s.Plugin = plugin

	for _, opt := range options {
		if err := opt(s); err != nil {
			return err
		}
	}

	return nil
}

// WithLogger adds logger to the Service
func WithLogger(logger *logger.Logger) ServiceOption {
	return func(s *Service) error {
		s.Logger = logger
		return nil
	}
}

// WithRng adds rng to the Service
func WithRng(rng io.Reader) ServiceOption {
	return func(s *Service) error {
		s.Rng = rng
		return nil
	}
}

// WithDataStore adds store to the Service
func WithDataStore(store *datastore.Store) ServiceOption {
	return func(s *Service) error {
		s.Store = store
		return nil
	}
}

// WithKeyStore adds store to the Service
func WithKeyStore(store keystore.Store) ServiceOption {
	return func(s *Service) error {
		s.KeyStore = store
		return nil
	}
}

// WithIPFS adds ipfs connector to the Service
func WithIPFS(ipfsConnector ipfs.Connector) ServiceOption {
	return func(s *Service) error {
		s.Ipfs = ipfsConnector
		return nil
	}
}

// WithTendermint adds tendermint node connector to the Service
func WithTendermint(tmConnector *tendermint.NodeConnector) ServiceOption {
	return func(s *Service) error {
		s.Tendermint = tmConnector
		return nil
	}
}

// WithMasterFiduciary adds master fiduciary connector to the Service
func WithMasterFiduciary(masterFiduciaryNodeID string) ServiceOption {
	return func(s *Service) error {
		s.SetMasterFiduciaryNodeID(masterFiduciaryNodeID)
		return nil
	}
}

// WithConfig adds config settings to the Service
func WithConfig(cfg *config.Config) ServiceOption {
	return func(s *Service) error {
		s.SetNodeID(cfg.Node.NodeID)
		s.SetMasterFiduciaryNodeID(cfg.Node.MasterFiduciaryNodeID)
		return nil
	}
}
