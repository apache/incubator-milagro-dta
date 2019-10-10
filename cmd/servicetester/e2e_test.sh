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
defaultURL="http://localhost:5556"
apiURL="${1:-$defaultURL}"
configdir=~/.milagro
host="34.246.173.153:26657"

status () {
  #Determine if an extension is running
  statusOutput=$(curl -s -X GET "$apiURL/$apiVersion/status" -H "accept: */*" -H "Content-Type: application/json")

  echo "$apiURL/$apiVersion/status"

  identity=$(echo $statusOutput | jq .nodeCID)
  extensionVendor=$(echo $statusOutput | jq -r .extensionVendor)
  plugin=$(echo $statusOutput | jq -r .plugin)
  echo "Plugin $plugin"

  if [ -z "${extensionVendor}" ]; then
      echo "Server Not Running"
      exit 1
  fi
}

###############################################################################################################################

execute_bitcoin () {
  # #Run 4 Tests against the Bitcoin Extension
  echo "Bitcoin Plugin Tests [2 Tests]"


 ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json"  -d "{\"beneficiaryIDDocumentCID\":\"\",\"extension\":{\"coin\":\"0\"}}" > ref ) &
  output1=$(fishhook $configdir $host "self" 2)
  ref=$(cat ref)
  commitment1=$(echo $output1 | jq .OrderPart2.CommitmentPublicKey)
  address1=$(echo $output1 | jq .OrderPart2.Extension.address)
  (sleep 1; curl -s -X POST "$apiURL/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref,\"beneficiaryIDDocumentCID\":$identity}" > /dev/null ) &
  output2=$(fishhook $configdir $host "self" 2)
  address2=$(echo $output2 | jq .OrderPart4.Extension.address)
  commitment2=$(echo $output2 | jq .OrderPart4.Extension.FinalPublicKey)
  #echo "Committment1 $commitment1 $address1"
  #echo "Committment2 $commitment2 $address2"
  if [ -z $commitment2 ]; then
      echo "Failed Commitment is empty"
      exit 1
  fi
  if [ $commitment1 != $commitment2 ]; then
    echo "Fail"
    exit 1
  fi
  if [ $address2 != $address2 ]; then
    echo "Fail"
    exit 1
  fi
  echo "Pass - Id, Order & OrderSecret(Beneficiary)"


  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json"  -d "{\"beneficiaryIDDocumentCID\":$identity,\"extension\":{\"coin\":\"0\"}}" > ref ) &
  output1=$(fishhook $configdir $host "self" 2)
  ref=$(cat ref)
  commitment1=$(echo $output1 | jq .OrderPart2.CommitmentPublicKey)
  address1=$(echo $output1 | jq .OrderPart2.Extension.address)
  (sleep 1; curl -s -X POST "$apiURL/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref}" > /dev/null ) &
  output2=$(fishhook $configdir $host "self" 2)
  address2=$(echo $output2 | jq .OrderPart4.Extension.address)
  commitment2=$(echo $output2 | jq .OrderPart4.Extension.FinalPublicKey)
  #echo "Committment1 $commitment1 $address1"
  #echo "Committment2 $commitment2 $address2"
  if [ -z $commitment2 ]; then
      echo "Failed Commitment is empty"
      exit 1
  fi
  if [ $commitment1 != $commitment2 ]; then
    echo "Fail"
    exit 1
  fi
  if [ $address2 != $address2 ]; then
    echo "Fail"
    exit 1
  fi
 echo "Pass - Id, Order(Beneficiary) & OrderSecret"
}


###############################################################################################################################

execute_safeguardsecret () {

  inputString="This is some random test text 1234567890!"
  echo "Encrypt a String [1 Test]"


  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":$identity,\"extension\":{\"plainText\":\"$inputString\"}}" > ref ) &
  output1=$(fishhook $configdir $host "self" 2)
  ref=$(cat ref)
  cipherText=$(echo $output1 | jq .OrderPart2.Extension.cypherText)

  #echo $cipherText
  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref,\"beneficiaryIDDocumentCID\":$identity,\"extension\":{\"cypherText\":$cipherText}}" > /dev/null) &
  output2=$(fishhook $configdir $host "self" 2)
  plaintext=$(echo $output2 | jq -r .OrderPart4.Extension.plainText)


  if [ -z "$plaintext" ]; then
      echo "Failed Commitment is empty"
      exit 1
  fi

  if [ "$inputString" == "$plaintext" ]; then
    echo "Order Create/Retrieve Pass"
  else
    echo "Order Create/Retrieve Fail"
    exit 1
  fi

}

# #############################################################################


execute_milagro () {
  echo "Milagro Tests [1 Test]"
  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":$identity}" > ref ) &
  output1=$(fishhook $configdir $host "self" 1)
  ref=$(cat ref)
  commitment1=$(echo $output1 | jq .OrderPart2.CommitmentPublicKey)
  #echo "Committment1 $commitment1"

  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref,\"beneficiaryIDDocumentCID\":$identity}" > /dev/null) &
  output2=$(fishhook $configdir $host "self" 3)
  commitment2=$(echo $output2 | jq .OrderPart4.Extension.FinalPublicKey)
  orderIndex=0
  #echo "Committment1 $commitment1"
  #echo "Committment2 $commitment2"

  if [ -z $commitment2 ]; then
      echo "Failed Commitment is empty"
      exit 1
  fi

  if [ $commitment1 == $commitment2 ]; then
    echo "Order Create/Retrieve Pass"
  else
    echo "Order Create/Retrieve Fail"
    exit 1
  fi
}






# #############################################################################

execute_orderlist () {
  echo "Milagro Tests [1 Test]"
  commitment2=$(echo $output2 | jq .commitment)
  outputList=$(curl -s -X GET "$apiURL/$apiVersion/order?page=0&perPage=2&sortBy=dateCreatedDsc" -H "accept: */*")
  orderReference=$(echo $outputList | jq -r ".orderReference | .[$orderIndex]")
  outputOrder=$(curl -s -X GET "$apiURL/$apiVersion/order/$orderReference" -H "accept: */*")

  #A simple smoke test to ensure some sort of order is returned
  hasSecret=`echo $outputOrder | grep "Secret"`

  if [ -z $hasSecret ]; then
      echo "Failed Order has error"
      exit 1
  else
     echo "Passed orderList & get"
  fi
}

# #############################################################################

status

if [ $plugin == "bitcoinwallet" ]; then
   execute_bitcoin
fi

if [ $plugin == "qredoplugin" ]; then
   execute_bitcoin
fi

if [ $plugin == "milagro"  ]; then
   execute_milagro
fi

if [ $plugin == "safeguardsecret" ]; then
    execute_safeguardsecret
fi
#execute_orderlist

