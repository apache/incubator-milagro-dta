#!/bin/bash

pushd () {
    command pushd "$@" > /dev/null
}

popd () {
    command popd "$@" > /dev/null
}

start_server () {
    GO111MODULE=on go build -o target/service github.com/apache/incubator-milagro-dta/cmd/service
    target/service daemon -service=$1 > /dev/null &
    pid=$!
    sleep 3
}

report () {
    if [ $2 -eq 0 ]; then
        echo "PASSED $1"
    else
        echo "FAILED $1"
    fi
}

test_plugin () {
    pushd .
    cd ../..
    start_server $1
    popd 
    ./e2e_test.sh #> /dev/null
    res=$?
    report "$1" $res 
    
    kill -s int $pid
}

cd "$(dirname "$0")"
test_plugin bitcoinwallet
test_plugin milagro
test_plugin safeguardsecret



