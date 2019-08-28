#!/bin/bash
#End to End Test of Services using curl/bash

apiVersion="v1"
simplePolicy='{ "example":"single fiduciary","walletRef": "CD-2367-227","beneficiarydocid":"","nodeId": "123","coin": 0,"sharingGroups": [{"groupref": "Back Office Team","threshold":1,"ids": [{"id": "alice","idType": "oidc","idRef": "Alice"}]}],"participantCount": 7}'

#simplePolicy='policy:{"example":"single fiduciary"}'

status () {
  #Determine if an extension is running
  statusOutput=$(curl -s -X GET "http://localhost:5556/$apiVersion/status" -H "accept: */*" -H "Content-Type: application/json")
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
  #Run 2 Tests against the Bitcoin Extension
  echo "Bitcoin Plugin Tests [2 Tests]"

#  echo "http://localhost:5556/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d                "{\"policy\":$simplePolicy,\"beneficiaryIDDocumentCID\":\"\",\"extension\":{\"coin\":\"0\"}}"

  output1=$(curl -s -X POST "http://localhost:5556/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d "{\"policy\":$simplePolicy,\"beneficiaryIDDocumentCID\":\"\",\"extension\":{\"coin\":\"0\"}}")
  echo $output1
  op1=$(echo $output1 | jq .orderReference)
  commitment1=$(echo $output1 | jq .commitment)
  address1=$(echo $output1 | jq .extension.address)
  output2=$(curl -s -X POST "http://localhost:5556/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$op1,\"beneficiaryIDDocumentCID\":$identity}")
  address2=$(echo $output2 | jq .extension.address)
  commitment2=$(echo $output2 | jq .commitment)

  echo "Committment1 $commitment1 $address1"
  echo "Committment2 $commitment2 $address2"

  if [ -z $commitment2 ]; then
      echo "Failed Commitment is empty"
      exit 1
  fi

  if [ $commitment1 == $commitment2 ]; then
    echo "Pass - Id, Order & OrderSecret(Beneficiary)"
  else
    echo "Fail"
    exit 1
  fi

  output3=$(curl -s -X POST "http://localhost:5556/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":$identity,\"extension\":{\"coin\":\"0\"}}")

  op3=$(echo $output3 | jq .orderReference)
  commitment3=$(echo $output3 | jq .commitment)
  address3=$(echo $output3 | jq .extension.address)
  output4=$(curl -s -X POST "http://localhost:5556/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$op3}")
  commitment4=$(echo $output4 | jq .commitment)
  address4=$(echo $output4 | jq .extension.address)
  orderReference=$(echo $output4 | jq .orderReference)
  orderIndex=1

  echo "Committment3 $commitment3 $address3"
  echo "Committment4 $commitment4 $address4"

  if [ -z $commitment4 ]; then
      echo "Failed Commitment is empty"
      exit 1
  fi

  if [ $commitment3 == $commitment4 ]; then
    echo "Pass - Id, Order(Beneficiary) & OrderSecret"
  else
      echo "Fail"
      exit 1
  fi
}


###############################################################################################################################

execute_safeguardsecret () {
  inputString="This is some random test text 1234567890!"
  echo "Encrypt a String [1 Test]"
  echo $output1
  output1=$(curl -s -X POST "http://localhost:5556/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":$identity,\"extension\":{\"plainText\":\"$inputString\"}}")
  echo $output1
  op1=$(echo $output1 | jq .orderReference)
  cipherText=$(echo $output1 | jq .extension.cypherText)
  tvalue=$(echo $output1 | jq .extension.t)
  vvalue=$(echo $output1 | jq .extension.v)
  commitment1=$(echo $output1 | jq .commitment)
  output2=$(curl -s -X POST "http://localhost:5556/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$op1,\"beneficiaryIDDocumentCID\":$identity,\"extension\":{\"cypherText\":$cipherText,\"t\":$tvalue,\"v\":$vvalue}}")
  result=$(echo $output2 | jq -r .extension.plainText)

  orderReference=$(echo $output2 | jq .orderReference)
  orderIndex=0


  if [ "$inputString" == "$result" ]; then
    echo "Pass"
  else
    echo "Fail"
    exit 1
  fi
}

# #############################################################################


execute_milagro () {
  echo "Milagro Tests [1 Test]"
  output1=$(curl -s -X POST "http://localhost:5556/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":$identity}")
  op1=$(echo $output1 | jq .orderReference)

  commitment1=$(echo $output1 | jq .commitment)
  output2=$(curl -s -X POST "http://localhost:5556/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$op1,\"beneficiaryIDDocumentCID\":$identity}")
  commitment2=$(echo $output2 | jq .commitment)

  orderReference=$(echo $output2 | jq .orderReference)
  orderIndex=0


  echo "Committment1 $commitment1"
  echo "Committment2 $commitment2"

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
  outputList=$(curl -s -X GET "http://localhost:5556/$apiVersion/order?page=0&perPage=2&sortBy=dateCreatedDsc" -H "accept: */*")
  orderReference=$(echo $outputList | jq -r ".orderReference | .[$orderIndex]")
  outputOrder=$(curl -s -X GET "http://localhost:5556/$apiVersion/order/$orderReference" -H "accept: */*")

  #A simple smoke test to ensure some sort of order is returned
  hasSecret=`echo $outputOrder | grep "Secret"`

  if [ -z "$hasSecret" ]; then
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

if [ $plugin == "milagro"  ]; then
   execute_milagro
fi

if [ $plugin == "safeguardsecret" ]; then
    execute_safeguardsecret
fi
execute_orderlist

