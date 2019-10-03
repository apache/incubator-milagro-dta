module github.com/apache/incubator-milagro-dta

require (
	github.com/TylerBrock/colorjson v0.0.0-20180527164720-95ec53f28296
	github.com/btcsuite/btcd v0.0.0-20190427004231-96897255fd17
	github.com/VividCortex/gohistogram v1.0.0 // indirect
	github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/coreos/go-oidc v2.0.0+incompatible
	github.com/go-kit/kit v0.9.0
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-test/deep v1.0.2
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gogo/protobuf v1.3.0
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/hokaccha/go-prettyjson v0.0.0-20190818114111-108c894c2c0e
	github.com/ipfs/go-datastore v0.0.5
	github.com/ipfs/go-ds-leveldb v0.0.2
	github.com/ipfs/go-ipfs v0.4.22
	github.com/ipfs/go-ipfs-api v0.0.1
	github.com/ipfs/go-ipfs-config v0.0.3
	github.com/ipfs/go-ipfs-files v0.0.3
	github.com/ipfs/interface-go-ipfs-core v0.0.8
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/libp2p/go-libp2p-crypto v0.0.2
	github.com/libp2p/go-libp2p-peer v0.1.1
	github.com/multiformats/go-multihash v0.0.5
	github.com/mwitkow/go-proto-validators v0.1.0
	github.com/pkg/errors v0.8.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/prometheus/client_golang v0.9.3
	github.com/stretchr/testify v1.4.0
	github.com/stumble/gorocksdb v0.0.3 // indirect
	github.com/tendermint/tendermint v0.32.4
	github.com/tyler-smith/go-bip39 v1.0.0
	github.com/urfave/cli v1.22.1
	go.etcd.io/bbolt v1.3.3
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
)

replace github.com/golangci/golangci-lint => github.com/golangci/golangci-lint v1.18.0

replace github.com/go-critic/go-critic v0.0.0-20181204210945-ee9bf5809ead => github.com/go-critic/go-critic v0.3.5-0.20190526074819-1df300866540

go 1.13
