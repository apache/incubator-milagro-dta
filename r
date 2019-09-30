GO111MODULE=on go build -o target/milagro github.com/apache/incubator-milagro-dta/cmd/service
export MILAGRO_HOME=~/.milagro
target/milagro daemon


