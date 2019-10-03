#!/bin/bash

# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.


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



