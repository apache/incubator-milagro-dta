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
RED='\033[0;31m'
GREEN='\033[0;32m'
ORANGE='\033[0;33m'
BLUE='\033[1;34m'
NC='\033[0m' # No Color



status () {
  #Determine if an extension is running
  statusOutput=$(curl -s -X GET "$apiURL/$apiVersion/status" -H "accept: */*" -H "Content-Type: application/json")

  #printf "$apiURL/$apiVersion/status\n"

  identity=$(echo $statusOutput | jq .nodeCID)
  extensionVendor=$(echo $statusOutput | jq -r .extensionVendor)
  plugin=$(echo $statusOutput | jq -r .plugin)
  printf "Plugin ${BLUE}$plugin ${NC}\n"

  if [ -z "${extensionVendor}" ]; then
       printf "${RED} Server Not Running{NC}\n"
      exit 1
  fi
}

###############################################################################################################################

execute_bitcoin () {
  # #Run 4 Tests against the Bitcoin Extension
  echo "  Bitcoin Plugin Tests [2 Tests]"


 ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json"  -d "{\"beneficiaryIDDocumentCID\":\"\",\"extension\":{\"coin\":\"0\"}}" > ref ) &
  output1=$(fishhook $configdir $host "self" 2)
  ref=$(cat ref | jq .orderReference)
  commitment1=$(echo $output1 | jq .OrderPart2.CommitmentPublicKey)
  address1=$(echo $output1 | jq .OrderPart2.Extension.address)
  (sleep 1; curl -s -X POST "$apiURL/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref,\"beneficiaryIDDocumentCID\":$identity}" > /dev/null ) &

  output2=$(fishhook $configdir $host "self" 2)
  address2=$(echo $output2 | jq .OrderPart4.Extension.address)
  commitment2=$(echo $output2 | jq .OrderPart4.Extension.FinalPublicKey)
  #echo "Committment1 $commitment1 $address1"
  #echo "Committment2 $commitment2 $address2"
  if [ -z $commitment2 ]; then
      printf "  ${RED}FAIL${NC} Commitment is empty\n"
      exit 1
  fi
  if [ $commitment1 != $commitment2 ]; then
    printf "  ${RED}FAIL${NC}\n "
    exit 1
  fi
  if [ $address2 != $address2 ]; then
    printf "  ${RED}FAIL${NC}\n "
    exit 1
  fi
  printf "  ${GREEN}Pass${NC} - Id, Order & OrderSecret(Beneficiary)\n"


  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json"  -d "{\"beneficiaryIDDocumentCID\":$identity,\"extension\":{\"coin\":\"0\"}}" > ref ) &
  output1=$(fishhook $configdir $host "self" 2)
  ref=$(cat ref | jq .orderReference)
  commitment1=$(echo $output1 | jq .OrderPart2.CommitmentPublicKey)
  address1=$(echo $output1 | jq .OrderPart2.Extension.address)
  (sleep 1; curl -s -X POST "$apiURL/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref}" > /dev/null ) &
  output2=$(fishhook $configdir $host "self" 2)
  address2=$(echo $output2 | jq .OrderPart4.Extension.address)
  commitment2=$(echo $output2 | jq .OrderPart4.Extension.FinalPublicKey)
  #echo "Committment1 $commitment1 $address1"
  #echo "Committment2 $commitment2 $address2"
  if [ -z $commitment2 ]; then
      printf "  ${RED}FAIL${NC}  Commitment is empty\n"
      exit 1
  fi
  if [ $commitment1 != $commitment2 ]; then
    printf "  ${RED}FAIL${NC}\n "
    exit 1
  fi
  if [ $address2 != $address2 ]; then
    printf "  ${RED}FAIL${NC}\n"
    exit 1
  fi
 printf "  ${GREEN}Pass${NC} - Id, Order(Beneficiary) & OrderSecret\n"
}


###############################################################################################################################

execute_safeguardsecret () {

  inputString="This is some random test text 1234567890!"
  printf "  Encrypt a String [1 Test]\n"


  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":$identity,\"extension\":{\"plainText\":\"$inputString\"}}" > ref ) &
  output1=$(fishhook $configdir $host "self" 2)
  ref=$(cat ref | jq .orderReference)
  cipherText=$(echo $output1 | jq .OrderPart2.Extension.cypherText)

  #echo $cipherText
  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref,\"beneficiaryIDDocumentCID\":$identity,\"extension\":{\"cypherText\":$cipherText}}" > /dev/null) &
  output2=$(fishhook $configdir $host "self" 2)
  plaintext=$(echo $output2 | jq -r .OrderPart4.Extension.plainText)


  if [ -z "$plaintext" ]; then
      printf "  ${RED}FAIL${NC}  Commitment is empty\n"
      exit 1
  fi

  if [ "$inputString" == "$plaintext" ]; then
    printf "  ${GREEN}Pass ${NC}Order Create/Retrieve\n"
  else
    printf "  ${RED}FAIL ${NC}Order Create/Retrieve\n"
    exit 1
  fi

}

# #############################################################################


execute_milagro () {
  echo "  Milagro Tests [1 Test]"
  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":$identity}" > ref ) &
  output1=$(fishhook $configdir $host "self" 1)
  ref=$(cat ref | jq .orderReference)
  commitment1=$(echo $output1 | jq .OrderPart2.CommitmentPublicKey)
  #echo "Committment1 $commitment1"

  ( sleep 1; curl -s -X POST "$apiURL/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref,\"beneficiaryIDDocumentCID\":$identity}" > /dev/null) &
  output2=$(fishhook $configdir $host "self" 3)
  commitment2=$(echo $output2 | jq .OrderPart4.Extension.FinalPublicKey)
  orderIndex=0
  #echo "Committment1 $commitment1"
  #echo "Committment2 $commitment2"

  if [ -z $commitment2 ]; then
      eprintfcho "  ${RED}FAIL${NC}  Commitment is empty\n"
      exit 1
  fi

  if [ $commitment1 == $commitment2 ]; then
    printf "  ${GREEN}Pass${NC} Order Create/Retrieve\n"
  else
    printf "  ${RED}FAIL${NC} Order Create/Retrieve\n"
    exit 1
  fi
}






# #############################################################################

execute_orderlist () {
  printf "Milagro Tests [1 Test]\n"
  commitment2=$(echo $output2 | jq .commitment)
  outputList=$(curl -s -X GET "$apiURL/$apiVersion/order?page=0&perPage=2&sortBy=dateCreatedDsc" -H "accept: */*")
  orderReference=$(echo $outputList | jq -r ".orderReference | .[$orderIndex]")
  outputOrder=$(curl -s -X GET "$apiURL/$apiVersion/order/$orderReference" -H "accept: */*")

  #A simple smoke test to ensure some sort of order is returned
  hasSecret=`echo $outputOrder | grep "Secret"`

  if [ -z $hasSecret ]; then
      printf "  ${RED}FAIL${NC} Order has erro\n"
      exit 1
  else
     printf "  ${GREEN}Pass${NC} orderList & get\n"
  fi
}

# #############################################################################

status


if [ $plugin == "bitcoinwallet" ]; then
   execute_bitcoin
fi


if [ $plugin == "milagro"  ]; then
   execute_milagro
fi


if [ $plugin == "safeguardsecret" ]; then
    execute_safeguardsecret
fi

#execute_orderlist

