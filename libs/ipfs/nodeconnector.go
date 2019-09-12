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

package ipfs

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	memoryds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	levelds "github.com/ipfs/go-ds-leveldb"
	cfg "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	core "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	repo "github.com/ipfs/go-ipfs/repo"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	path "github.com/ipfs/interface-go-ipfs-core/path"
	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"

	"github.com/pkg/errors"
)

const (
	defaultShardFun = "prefix"
)

// NodeConnectorBuilder for building the IPFS embedded node connector
type NodeConnectorBuilder struct {
	ctx            context.Context
	localAddresses []string
	bootstrapPeers []string
	dstore         repo.Datastore
}

// NodeConnectorOption function
type NodeConnectorOption func(*NodeConnectorBuilder) error

// NodeConnector is IPFS embedded node
type NodeConnector struct {
	id   string
	ctx  context.Context
	node *core.IpfsNode
	api  coreiface.CoreAPI
}

// NewNodeConnector inisialises and runs a new IPFS Node
func NewNodeConnector(options ...NodeConnectorOption) (Connector, error) {
	cb := &NodeConnectorBuilder{
		localAddresses: []string{},
		bootstrapPeers: []string{},
		ctx:            context.Background(),
		dstore:         nil,
	}

	for _, option := range options {
		if err := option(cb); err != nil {
			return nil, err
		}
	}

	if len(cb.localAddresses) == 0 {
		cb.localAddresses = []string{"/ip4/0.0.0.0/tcp/4001"}
	}

	if cb.dstore == nil {
		return nil, errors.New("IPFS datastore not initialized")
	}

	// Generate Node keys
	priv, _, err := ci.GenerateKeyPairWithReader(ci.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}
	privkeyb, err := priv.Bytes()
	if err != nil {
		return nil, err
	}
	pid, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	conf := cfg.Config{
		Bootstrap: cb.bootstrapPeers,
	}
	conf.Addresses.Swarm = cb.localAddresses
	conf.Identity.PeerID = pid.Pretty()
	conf.Identity.PrivKey = base64.StdEncoding.EncodeToString(privkeyb)
	conf.Routing.Type = "dht"
	conf.Swarm.ConnMgr = cfg.ConnMgr{
		Type:        "basic",
		LowWater:    20,
		HighWater:   100,
		GracePeriod: "30s",
	}

	appRepo := &repo.Mock{
		D: dsync.MutexWrap(cb.dstore),
		C: conf,
	}

	node, err := core.NewNode(cb.ctx, &core.BuildCfg{
		Repo:   appRepo,
		Online: true,
	})
	if err != nil {
		return nil, err
	}

	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	return &NodeConnector{
		id:   pid.Pretty(),
		ctx:  cb.ctx,
		node: node,
		api:  api,
	}, nil
}

// Add adds a data to the IPFS network and returns the ipfs path
func (c *NodeConnector) Add(data []byte) (string, error) {
	f := files.NewBytesFile(data)
	defer f.Close()
	r, err := c.api.Unixfs().Add(c.ctx, f)
	if err != nil {
		return "", err
	}
	return r.Cid().String(), nil
}

// Get gets a data from ipfs path
func (c *NodeConnector) Get(ipfsPath string) ([]byte, error) {
	p := path.New(ipfsPath)

	node, err := c.api.Unixfs().Get(c.ctx, p)
	if err != nil {
		return nil, err
	}
	defer node.Close()

	f := files.ToFile(node)
	if f == nil {
		return nil, errors.New("not a file")
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// TODO: Pin the document
	return b, nil
}

// AddJSON encodes the data in JSON and adds it to the IPFS
func (c *NodeConnector) AddJSON(data interface{}) (string, error) {
	jd, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return c.Add(jd)
}

// GetJSON gets data from IPFS and decodes it from JSON
func (c *NodeConnector) GetJSON(path string, data interface{}) error {
	jd, err := c.Get(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jd, data); err != nil {
		return errors.Wrap(ErrDocumentNotValid, err.Error())
	}

	return nil
}

// GetID returns the local id
func (c *NodeConnector) GetID() string {
	return c.id
}

// AddLocalAddress adds a new local IPFS Swarm address
func AddLocalAddress(addr string) NodeConnectorOption {
	return func(cb *NodeConnectorBuilder) error {
		cb.localAddresses = append(cb.localAddresses, addr)
		return nil
	}
}

// AddBootstrapPeer adds a IPFS peer to the list of IPFS bootatrap addresses
func AddBootstrapPeer(peerAddr ...string) NodeConnectorOption {
	return func(cb *NodeConnectorBuilder) error {
		cb.bootstrapPeers = append(cb.bootstrapPeers, peerAddr...)
		return nil
	}
}

// AddDefaultBootstrapPeers adds the default IPFS bootstrap addresses
// The default IPFS addresses are maintained by the IPFS team
func AddDefaultBootstrapPeers() NodeConnectorOption {
	return func(cb *NodeConnectorBuilder) error {
		cb.bootstrapPeers = append(cb.bootstrapPeers, cfg.DefaultBootstrapAddresses...)
		return nil
	}
}

// WithLevelDatastore initialize IPFS node with LevelDB datastore
func WithLevelDatastore(path string) NodeConnectorOption {
	return func(cb *NodeConnectorBuilder) (err error) {
		cb.dstore, err = levelds.NewDatastore(path, &levelds.Options{
			Compression: 0,
		})

		return nil
	}
}

// WithMemoryDatastore initialize IPFS node with memory datastore
func WithMemoryDatastore() NodeConnectorOption {
	return func(cb *NodeConnectorBuilder) (err error) {
		cb.dstore = memoryds.NewMapDatastore()

		return nil
	}
}

// WithContext sets the context
func WithContext(ctx context.Context) NodeConnectorOption {
	return func(cb *NodeConnectorBuilder) error {
		cb.ctx = ctx
		return nil
	}
}
