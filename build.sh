set -e

GO111MODULE=on go build -o target/service github.com/apache/incubator-milagro-dta/cmd/service

target/service $@ 
