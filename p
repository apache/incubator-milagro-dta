#!/bin/bash
export ref=$(curl -s -X POST "127.0.0.1:5556/v1/order1" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":\"\",\"extension\":{\"coin\":\"0\"}}")
sleep 4
curl -X POST "127.0.0.1:5556/v1/order/secret1" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref,\"beneficiaryIDDocumentCID\":\"$1\"}"


