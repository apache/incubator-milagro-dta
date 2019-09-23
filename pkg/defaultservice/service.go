// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownershis.  The ASF licenses this file
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
Package defaultservice - Service that Milagro D-TA provides with no-plugins
*/
package defaultservice

import (
	"io"
	"time"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/ipfs"
	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/libs/transport"
	"github.com/apache/incubator-milagro-dta/pkg/api"
)

var (
	extensionVendor = "Milagro"
	pluginName      = "milagro"
)

// Service implements the default functionality
// It also implements the ServicePlugin interface
type Service struct {
	Plugin                Plugable
	Logger                *logger.Logger
	Rng                   io.Reader
	Store                 *datastore.Store
	Ipfs                  ipfs.Connector
	MasterFiduciaryServer api.ClientService
	nodeID                string
	masterFiduciaryNodeID string
}

//NewService returns a default implementation of Service
func NewService() *Service {
	s := &Service{}
	return s
}

// Name of the plugin
func (s *Service) Name() string {
	return pluginName
}

// Vendor of the plugin
func (s *Service) Vendor() string {
	return "Milagro"
}

// NodeID returns the node CID
func (s *Service) NodeID() string {
	return s.nodeID
}

// SetNodeID sets the Node CID
func (s *Service) SetNodeID(nodeID string) {
	s.nodeID = nodeID
}

// MasterFiduciaryNodeID returns the Master Fiduciary NodeID
func (s *Service) MasterFiduciaryNodeID() string {
	return s.masterFiduciaryNodeID
}

// SetMasterFiduciaryNodeID sets the Master Fiduciary NodeID
func (s *Service) SetMasterFiduciaryNodeID(masterFiduciaryNodeID string) {
	s.masterFiduciaryNodeID = masterFiduciaryNodeID
}

// Status of the server
func (s *Service) Status(apiVersion, nodeType string) (*api.StatusResponse, error) {
	return &api.StatusResponse{
		Application:     "Milagro Distributed Trust",
		APIVersion:      apiVersion,
		ExtensionVendor: s.Vendor(),
		NodeType:        nodeType,
		NodeCID:         s.nodeID,
		TimeStamp:       time.Now(),
		Plugin:          s.Plugin.Name(),
	}, nil
}

// Endpoints for extending the plugin endpoints
func (s *Service) Endpoints() (namespace string, endpoints transport.HTTPEndpoints) {
	return s.Name(), nil
}
