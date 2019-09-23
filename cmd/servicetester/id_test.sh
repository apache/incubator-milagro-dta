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


#End to End Test of Services using curl/bash

apiVersion="v1"


status () {
  #Determine if an extension is running
  statusOutput=$(curl -s -X GET "http://localhost:5556/$apiVersion/status" -H "accept: */*" -H "Content-Type: application/json")
  identity=$(echo $statusOutput | jq .nodeCID)
  extensionVendor=$(echo $statusOutput | jq -r .extensionVendor)
  plugin=$(echo $statusOutput | jq -r .plugin)
  echo "ID $identity"
  echo "Plugin: $plugin"
  echo "Vendor: $extensionVendor"

  if [ -z '${extensionVendor}' ]; then
      echo "Server Not Running"
      exit
  fi
}




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


cd "$(dirname "$0")"
pushd .
cd ../..
start_server milagro

popd 

status

#Create a new Idenity
output1=$(curl -s -X POST "http://localhost:5556/$apiVersion/identity" -H "accept: */*" -H "Content-Type: application/json" -d "{\"Name\":\"AA\"}")
docid=$(echo $output1 | jq -r .idDocumentCID)

#Get the single ID
output2=$(curl -s -X GET "http://localhost:5556/$apiVersion/identity/$docid" -H "accept: */*" -H "Content-Type: application/json")



SikePublicKey=$(echo $output2 | jq -r .sikePublicKey)
BlsPublicKey=$(echo $output2 | jq -r .blsPublicKey)
BeneficiaryECPubl=$(echo $output2 | jq -r .beneficiaryECPublicKey)

#Get a list of all identities
output3=$(curl -s -X GET "http://localhost:5556/$apiVersion/identity?page=0&perPage=1&sortBy=dateCreatedDsc")
firstDoc=$(echo $output3 | jq -r '.idDocumentList[0].idDocumentCID')

#Pass test if Created ID is top of the list
if [ $firstDoc == $docid ]; then
  echo "Passed"
else
  echo "fail"
fi

kill -s int $pid
