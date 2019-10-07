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
Package main - handles config, initialisation and starts the service daemon
*/
package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/keystore"
	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/libs/transport"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/config"
	"github.com/apache/incubator-milagro-dta/pkg/defaultservice"
	"github.com/apache/incubator-milagro-dta/pkg/endpoints"
	"github.com/apache/incubator-milagro-dta/pkg/identity"
	"github.com/apache/incubator-milagro-dta/pkg/tendermint"
	"github.com/apache/incubator-milagro-dta/plugins"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/pkg/errors"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initConfig(args []string) error {
	cfg := config.DefaultConfig()
	logger, err := logger.NewLogger("text", "info")
	if err != nil {
		return err
	}

	initOptions, err := parseInitOptions(args)
	if err != nil {
		return err
	}

	if initOptions.Interactive {
		if err := interactiveSetup(initOptions); err != nil {
			return err
		}
	}

	if initOptions.NodeName != "" {
		cfg.Node.NodeName = initOptions.NodeName
	} else {
		cfg.Node.NodeName = generateRandomName()
		logger.Info("Node name not provided. Generated random name: %s", cfg.Node.NodeName)
	}
	cfg.Plugins.Service = initOptions.ServicePlugin

	// Init the config folder
	config.Init(configFolder(), cfg)

	// Init Tendermint node connector
	tmConnector, err := tendermint.NewNodeConnector(cfg.Blockchain.BroadcastNode, cfg.Node.NodeID, nil, logger)
	if err != nil {
		return errors.Wrap(err, "Blockchain Node connector")
	}

	logger.Info("Node connector address: %s", cfg.Blockchain.BroadcastNode)

	keyStore, err := keystore.NewFileStore(filepath.Join(configFolder(), keysFile))
	if err != nil {
		return err
	}
	_, rawDocID, secret, err := identity.CreateIdentity(cfg.Node.NodeName)
	if err != nil {
		return err
	}
	newID, err := identity.StoreIdentity(rawDocID, secret, tmConnector, keyStore)
	if err != nil {
		return err
	}

	logger.Info("New Node ID: %s", newID)

	cfg.Node.NodeID = newID
	if initOptions.MasterFidNodeID != "" {
		cfg.Node.MasterFiduciaryNodeID = initOptions.MasterFidNodeID
	} else {
		cfg.Node.MasterFiduciaryNodeID = newID
	}

	if cfg.Node.MasterFiduciaryNodeID == "" {
		cfg.Node.MasterFiduciaryNodeID = newID
	}

	if cfg.Node.NodeType == "" {
		cfg.Node.NodeType = "multi"
	}

	return config.SaveConfig(configFolder(), cfg)
}

func startDaemon(args []string) error {
	cfg, err := parseConfig(args)
	if err != nil {
		return err
	}

	logger, err := logger.NewLogger(
		cfg.Log.Format,
		cfg.Log.Level,
	)
	if err != nil {
		return errors.Wrap(err, "init logger")
	}

	logger.Debug("Logger in DEBUG mode!")

	// Create KV store
	logger.Info("Datastore type: %s", cfg.Node.Datastore)
	store, err := initDataStore(cfg.Node.Datastore)
	if err != nil {
		return errors.Wrap(err, "init datastore")
	}

	keyStore, err := keystore.NewFileStore(filepath.Join(configFolder(), keysFile))
	if err != nil {
		return err
	}

	// Setup Endpoint authorizer
	var authorizer transport.Authorizer
	switch cfg.HTTP.OIDCProvider {
	case "":
		authorizer = &transport.EmptyAuthorizer{}
	case "local":
		authorizer = &transport.LocalAuthorizer{}
	default:
		authorizer, err = transport.NewOIDCAuthorizer(
			cfg.HTTP.OIDCClientID,
			cfg.HTTP.OIDCProvider,
		)
		if err != nil {
			return errors.Wrap(err, "init authorizer")
		}
	}

	// Init Tendermint node connector
	tmConnector, err := tendermint.NewNodeConnector(cfg.Blockchain.BroadcastNode, cfg.Node.NodeID, store, logger)
	if err != nil {
		return errors.Wrap(err, "Blockchain Node connector")
	}

	//The Server must have a valid ID before starting up
	svcPlugin := plugins.FindServicePlugin(cfg.Plugins.Service)
	if svcPlugin == nil {
		return errors.Errorf("invalid plugin: %v", cfg.Plugins.Service)
	}

	if err := svcPlugin.Init(
		svcPlugin,
		defaultservice.WithLogger(logger),
		defaultservice.WithRng(rand.Reader),
		defaultservice.WithDataStore(store),
		defaultservice.WithKeyStore(keyStore),
		defaultservice.WithTendermint(tmConnector),
		defaultservice.WithConfig(cfg),
	); err != nil {
		return errors.Errorf("init service plugin %s", cfg.Plugins.Service)
	}
	logger.Info("Service plugin loaded: %s", svcPlugin.Name())

	if err := identity.CheckIdentity(cfg.Node.NodeID, cfg.Node.NodeName, tmConnector, keyStore); err != nil {
		return errors.Wrap(err, "Invalid node identity")
	}

	svcPlugin.SetMasterFiduciaryNodeID(cfg.Node.MasterFiduciaryNodeID)
	svcPlugin.SetNodeID(cfg.Node.NodeID)

	// Create metrics
	duration := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "milagro",
		Subsystem: "milagroservice",
		Name:      "request_duration_seconds",
		Help:      "Request duration in seconds.",
	}, []string{"method", "success"})

	// Stop chan
	ctx, cancelContext := context.WithCancel(context.Background())

	errChan := make(chan error)

	logger.Info("Node Type: %v", strings.ToLower(cfg.Node.NodeType))
	logger.Info("Node ID:  %v", svcPlugin.NodeID())
	logger.Info("Master Fiduciary: %v", svcPlugin.MasterFiduciaryNodeID())

	//Connect to Blockchain - Tendermint
	go func() {
		processFn := func(tx *api.BlockChainTX) error {
			switch tx.Processor {
			case "none":
				return nil
			case "dump":
				svcPlugin.Dump(tx)
			case "v1/fulfill/order":
				svcPlugin.FulfillOrder(tx)
			case "v1/order2":
				svcPlugin.Order2(tx)
			case "v1/fulfill/order/secret":
				svcPlugin.FulfillOrderSecret(tx)
			case "v1/order/secret2":
				svcPlugin.OrderSecret2(tx)
			default:
				return errors.New("Unknown processor")
			}
			return nil
		}

		logger.Info("Starting Blockchain listener to node: %v", cfg.Blockchain.BroadcastNode)
		errChan <- tmConnector.Subscribe(ctx, processFn)
		// errChan <- tendermint.Subscribe(svcPlugin, store, logger, cfg.Node.NodeID, cfg.HTTP.ListenAddr)
	}()

	// Start the application http server
	go func() {
		httpEndpoints := endpoints.Endpoints(svcPlugin, cfg.HTTP.CorsAllow, authorizer, logger, cfg.Node.NodeType, svcPlugin)
		httpHandler := transport.NewHTTPHandler(httpEndpoints, logger, duration)

		logger.Info("Starting HTTP listener on %v", cfg.HTTP.ListenAddr)
		errChan <- http.ListenAndServe(cfg.HTTP.ListenAddr, httpHandler)
	}()

	if cfg.HTTP.MetricsAddr != "" {
		http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
		// Start the debug and metrics http server
		go func() {
			logger.Info("Starting metrics listener on %v", cfg.HTTP.MetricsAddr)
			errChan <- http.ListenAndServe(cfg.HTTP.MetricsAddr, http.DefaultServeMux)
		}()
	}

	// Start the signal handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- errors.Errorf("received signal %s", <-c)
	}()

	stopErr := <-errChan
	if stopErr != nil {
		_ = logger.Log("exit", stopErr.Error())
	}

	cancelContext()
	return store.Close()
}

func initDataStore(ds string) (*datastore.Store, error) {
	var dsBackend datastore.Backend
	var err error
	switch ds {
	case "embedded":
		dsBackend, err = datastore.NewBoltBackend(filepath.Join(configFolder(), "datastore.dat"))
	default:
		return nil, errors.Errorf("invalid datastore: %s", ds)
	}
	if err != nil {
		return nil, err
	}

	store, err := datastore.NewStore(datastore.WithBackend(dsBackend), datastore.WithCodec(datastore.NewGOBCodec()))
	return store, err
}

func main() {
	var err error

	cmd, args := parseCommand()
	switch cmd {
	default:
		fmt.Println(printHelp())
	case cmdInit:
		err = initConfig(args)
	case cmdDaemon:
		err = startDaemon(args)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
