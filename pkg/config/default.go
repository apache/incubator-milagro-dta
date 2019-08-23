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

package config

import (
	ipfsCfg "github.com/ipfs/go-ipfs-config"
)

// DefaultConfig -
func DefaultConfig() *Config {
	return &Config{
		HTTP:    defaultHTTPConfig(),
		Node:    defaultNodeConfig(),
		Log:     defaultLogConfig(),
		IPFS:    defaultIPFSConfig(),
		Plugins: defaultPluginsConfig(),
	}
}

func defaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		ListenAddr:    ":5556",
		MetricsAddr:   ":5557",
		OIDCProvider:  "",
		OIDCClientID:  "",
		OIDCClientKey: "",
		CorsAllow:     "http://localhost:5556",
	}
}

// LogConfig -
func defaultLogConfig() LogConfig {
	return LogConfig{
		Format: "text",
		Level:  "info",
	}
}

// IPFSConfig -
func defaultIPFSConfig() IPFSConfig {
	return IPFSConfig{
		Bootstrap:     ipfsCfg.DefaultBootstrapAddresses,
		Connector:     "embedded",
		ListenAddress: "/ip4/0.0.0.0/tcp/4001",
		APIAddress:    "http://localhost:5001",
	}
}

// NodeConfig -
func defaultNodeConfig() NodeConfig {
	return NodeConfig{
		NodeType:              "multi",
		MasterFiduciaryServer: "http://localhost:5556",
		MasterFiduciaryNodeID: "",
		NodeID:                "",
		Datastore:             "embedded",
	}
}

// PluginsConfig -
func defaultPluginsConfig() PluginsConfig {
	return PluginsConfig{
		Service: "milagro",
	}
}
