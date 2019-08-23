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
Package config - create default config and read config values
*/
package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

const (
	configFileName = "config.yaml"
)

var (
	// ErrConfigExists is returned on init if config file exists
	ErrConfigExists = errors.New("Config file already exists")
	// ErrConfigNotFound is returned when configuration not found
	ErrConfigNotFound = errors.New("Configuration not found")
)

// HTTPConfig -
type HTTPConfig struct {
	ListenAddr    string `yaml:"listenAddr"`
	MetricsAddr   string `yaml:"metricsAddr"`
	OIDCProvider  string `yaml:"oidcProvider"`
	OIDCClientID  string `yaml:"oidcClientID"`
	OIDCClientKey string `yaml:"oidcClientKey"`
	CorsAllow     string `yaml:"corsAllow"`
}

// LogConfig -
type LogConfig struct {
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
}

// IPFSConfig -
type IPFSConfig struct {
	Connector     string   `yaml:"connector"`
	Bootstrap     []string `yaml:"bootstrap"`
	ListenAddress string   `yaml:"listenAddress"`
	APIAddress    string   `yaml:"apiAddress"`
}

// NodeConfig -
type NodeConfig struct {
	NodeType              string `yaml:"nodeType"`
	MasterFiduciaryServer string `yaml:"masterFiduciaryServer"`
	MasterFiduciaryNodeID string `yaml:"masterFiduciaryNodeID"`
	NodeID                string `yaml:"nodeID"`
	NodeName              string `yaml:"nodeName"`
	Datastore             string `yaml:"dataStore"`
}

// PluginsConfig -
type PluginsConfig struct {
	Service string `yaml:"service"`
}

// Config -
type Config struct {
	HTTP    HTTPConfig    `yaml:"http"`
	Node    NodeConfig    `yaml:"node"`
	Log     LogConfig     `yaml:"log"`
	IPFS    IPFSConfig    `yaml:"ipfs"`
	Plugins PluginsConfig `yaml:"plugins"`
}

// Init initialise config folder with default options
func Init(folder string, config *Config) error {
	configFilePath := filepath.Join(folder, configFileName)

	_, err := os.Stat(configFilePath)
	if err == nil {
		return ErrConfigExists
	}
	if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(folder, 0700); err != nil {
		return err
	}

	return SaveConfig(folder, config)
}

// ParseConfig parses configuration file
func ParseConfig(folder string) (*Config, error) {
	configFilePath := filepath.Join(folder, configFileName)
	if _, err := os.Stat(configFilePath); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigNotFound
		}
		return nil, err
	}

	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	return cfg, err
}

// SaveConfig stores configuration
func SaveConfig(folder string, cfg *Config) error {
	configFilePath := filepath.Join(folder, configFileName)

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFilePath, b, 0600)
}
