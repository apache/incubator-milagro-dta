#! /bin/bash

apiVersion="v1"
defaultURL="http://localhost:5556"
apiURL="${1:-$defaultURL}"

echo "DTA URL: $apiURL"

statusOutput=$(curl -s -X GET "$apiURL/$apiVersion/status")
identity=$(echo $statusOutput | jq -r .nodeCID)
plugin=$(echo $statusOutput | jq -r .plugin)

if [ -z "${identity}" ]; then
  echo "Server Not Running"
  exit 1
fi

benID="${2:-$identity}"

echo "DTA Plugin: $plugin"
echo "DTA ID: $identity"
echo "BeneficiaryID: $benID"


respOrder=$(curl -s -X POST "$apiURL/$apiVersion/order" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":\"$benID\",\"extension\":{\"coin\":\"0\"}}")

orderRef=$(echo $respOrder | jq '.orderReference')
if [ -z "${orderRef}" ]; then
	echo "Create order invalid response"
	exit 1
fi


#sleep long enough for blockchain to catch up
sleep 4

curl -X POST "$apiURL/$apiVersion/order/secret" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$orderRef,\"beneficiaryIDDocumentCID\":\"$benID\"}"


