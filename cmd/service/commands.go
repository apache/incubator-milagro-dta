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

package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/apache/incubator-milagro-dta/pkg/config"
)

const (
	envMilagroHome      = "MILAGRO_HOME"
	milagroConfigFolder = ".milagro"

	cmdInit   = "init"
	cmdDaemon = "daemon"
)

func configFolder() string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		userHome = ""
	}

	return getEnv(envMilagroHome, filepath.Join(userHome, milagroConfigFolder))
}

func getEnv(name, defaultValue string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return v
}

func parseCommand() (cmd string, args []string) {
	if len(os.Args) < 2 {
		return
	}
	return os.Args[1], os.Args[2:]
}

func printHelp() string {
	return `Milagro DTA
USAGE
	milagro <command> [options]
	
COMMANDS
	init	Initialize configuration
	daemon	Starts the milagro daemon
	`
}

func parseConfig(args []string) (*config.Config, error) {
	cfg, err := config.ParseConfig(configFolder())
	if err != nil {
		return nil, err
	}

	fs := flag.NewFlagSet("daemon", flag.ExitOnError)
	fs.StringVar(&(cfg.Plugins.Service), "service", cfg.Plugins.Service, "Service plugin")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return cfg, nil
}
