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
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type initOptions struct {
	NodeName             string
	MasterFidNodeID      string
	MasterFidNodeAddress string
	ServicePlugin        string
	Interactive          bool
}

func parseInitOptions(args []string) (*initOptions, error) {
	i := &initOptions{}

	var masterFidNode string

	fs := flag.NewFlagSet("init", flag.ExitOnError)
	fs.StringVar(&i.NodeName, "nodename", "", "Node name")
	fs.StringVar(&masterFidNode, "masterfiduciarynode", "", "Master fiduciary node")
	fs.StringVar(&i.ServicePlugin, "service", "milagro", "Service plugin")
	fs.BoolVar(&i.Interactive, "interactive", false, "Interactive setup")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if masterFidNode != "" {
		spl := strings.Split(masterFidNode, ",")
		if len(spl) != 2 {
			return nil, errors.New("Invalid master fiduciary node format")
		}
		i.MasterFidNodeID = strings.TrimSpace(spl[0])
		i.MasterFidNodeAddress = strings.TrimSpace(spl[1])
	}

	return i, nil

}

func cliInput(text string) (val string, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(text + ": ")
	val, err = reader.ReadString('\n')
	val = strings.TrimSpace(val)
	return
}

func interactiveSetup(i *initOptions) error {
	var err error
	i.NodeName, err = cliInput("What is your node name?. Leave blank to generate a random name")
	if err != nil {
		return err
	}
	i.MasterFidNodeID, err = cliInput("What is your Master Fiduciary DTA’s node name? Leave blank to use this DTA as the Master Fiduciary")
	if err != nil {
		return err
	}
	if i.MasterFidNodeID != "" {
		i.MasterFidNodeAddress, err = cliInput("What is your Master Fiduciary DTA’s address?")
		if err != nil {
			return err
		}
	}

	plugin, err := cliInput("What plugin do you want to install? (B)itcoin wallet address generator or (S)afeguard secret. Leave blank for no plugin")
	if err != nil {
		return err
	}

	switch strings.ToLower(plugin) {
	case "b", "bitcoinwallet":
		i.ServicePlugin = "bitcoinwallet"
	case "s", "safeguardsecret":
		i.ServicePlugin = "safeguardsecret"
	default:
		i.ServicePlugin = "milagro"
	}

	return nil
}

func generateRandomName() string {
	b := make([]byte, 8)
	rand.Read(b)

	return hex.EncodeToString(b)
}
