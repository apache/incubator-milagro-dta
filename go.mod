module github.com/apache/incubator-milagro-dta

require (
	github.com/btcsuite/btcd v0.22.1
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/go-kit/kit v0.10.0
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-test/deep v1.0.2
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/ipfs/go-datastore v0.5.1
	github.com/ipfs/go-ds-leveldb v0.5.0
	github.com/ipfs/go-ipfs v0.13.1
	github.com/ipfs/go-ipfs-api v0.0.1
	github.com/ipfs/go-ipfs-config v0.0.3
	github.com/ipfs/go-ipfs-files v0.1.1
	github.com/ipfs/interface-go-ipfs-core v0.7.0
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/libp2p/go-libp2p-crypto v0.1.0
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/multiformats/go-multihash v0.1.0
	github.com/mwitkow/go-proto-validators v0.1.0
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_golang v1.12.1
	github.com/stretchr/testify v1.7.1
	github.com/tyler-smith/go-bip39 v1.0.0
	go.etcd.io/bbolt v1.3.3
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
)

replace github.com/golangci/golangci-lint => github.com/golangci/golangci-lint v1.18.0

replace github.com/go-critic/go-critic v0.0.0-20181204210945-ee9bf5809ead => github.com/go-critic/go-critic v0.3.5-0.20190526074819-1df300866540

go 1.13
