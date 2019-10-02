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
Package plugins - registers plugins with the service
*/
package plugins

import (
	"io"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/ipfs"
	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/pkg/config"
	"github.com/apache/incubator-milagro-dta/pkg/defaultservice"
	"github.com/apache/incubator-milagro-dta/pkg/service"
)

var plugins []Plugin

// Plugin is the main plugin interface
type Plugin interface {
	Name() string
	Vendor() string
}

// ServicePlugin interface
type ServicePlugin interface {
	defaultservice.Plugable
	service.Service

	Init(plugin defaultservice.Plugable, logger *logger.Logger, rng io.Reader, store *datastore.Store, ipfsConnector ipfs.Connector, cfg *config.Config) error
}

func registerPlugin(p Plugin) {
	plugins = append(plugins, p)
}

// FindServicePlugin returns a registered ServicePlugin by name
// Returns nil if the plugin is not loaded
func FindServicePlugin(name string) ServicePlugin {

	for _, p := range plugins {
		sp, ok := p.(ServicePlugin)
		if !ok {
			continue
		}

		if sp.Name() == name {
			return sp
		}
	}

	return nil
}
