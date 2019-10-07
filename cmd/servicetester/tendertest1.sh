ref=$(curl -s -X POST "127.0.0.1:5556/v1/order1" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":\"QmecffRZTSJDETCRLcjcPe7ynmYNyYpnh7WKzKTdmX1GBZ\",\"extension\":{\"coin\":\"0\"}}")

#sleep long enough for blockchain to catch up
sleep 4

curl -X POST "127.0.0.1:5556/v1/order/secret1" -H "accept: */*" -H "Content-Type: application/json" -d "{\"orderReference\":$ref,\"beneficiaryIDDocumentCID\":\"QmecffRZTSJDETCRLcjcPe7ynmYNyYpnh7WKzKTdmX1GBZ\"}"


