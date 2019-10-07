ref=$(curl -s -X POST "127.0.0.1:5556/v1/order1" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":\"90bac919210be1ef29cd6da22d512c2f8a04693544fe6e474cb5d90c6fbe4645\",\"extension\":{\"coin\":\"0\"}}")

#sleep long enough for blockchain to catch up
sleep 4

curl -X POST "127.0.0.1:5556/v1/order/secret1" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref,\"beneficiaryIDDocumentCID\":\"90bac919210be1ef29cd6da22d512c2f8a04693544fe6e474cb5d90c6fbe4645\"}"


