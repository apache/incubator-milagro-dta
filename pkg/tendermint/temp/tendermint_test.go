package tendermint

import "testing"

var (
	nodeID = "QmT4y4MtV5mvPHkFjfUQYQ7h1WvAagMy2GTJCn2bF8DQb7"
)

func Test_Order1(t *testing.T) {
	a := "eyJQcm9jZXNzb3IiOiJ2MS9mdWxmaWxsL29yZGVyIiwiU2VuZGVySUQiOiJRbVQ0eTRNdFY1bXZQSGtGamZVUVlRN2gxV3ZBYWdNeTJHVEpDbjJiRjhEUWI3IiwiUmVjaXBpZW50SUQiOiJRbVQ0eTRNdFY1bXZQSGtGamZVUVlRN2gxV3ZBYWdNeTJHVEpDbjJiRjhEUWI3IiwiUGF5bG9hZCI6ImV5SnZjbVJsY2xCaGNuUXhRMGxFSWpvaVVXMVpVRU5xVEVGME1tbzVVbWhxU0U1TVkwRnVObEF5WTJseVJHWjZTRlpFWTBwMFkzbGtUVFZ5VWxoM1V5SXNJbVJ2WTNWdFpXNTBRMGxFSWpvaVVXMVVOSGswVFhSV05XMTJVRWhyUm1wbVZWRlpVVGRvTVZkMlFXRm5UWGt5UjFSS1EyNHlZa1k0UkZGaU55SjkifQ=="
	err := HandleChainTX(nodeID, a)
	if err != nil {
		panic(err)
	}
}

func Test_FullFill(t *testing.T) {
	a := "eyJQcm9jZXNzb3IiOiJPUkRFUl9SRVNQT05TRSIsIlNlbmRlcklEIjoiUW1UNHk0TXRWNW12UEhrRmpmVVFZUTdoMVd2QWFnTXkyR1RKQ24yYkY4RFFiNyIsIlJlY2lwaWVudElEIjoiUW1UNHk0TXRWNW12UEhrRmpmVVFZUTdoMVd2QWFnTXkyR1RKQ24yYkY4RFFiNyIsIlBheWxvYWQiOiJleUp2Y21SbGNsQmhjblF5UTBsRUlqb2lVVzFVZUZka1ltZEdhRGxHYWpGMlJIbFhlazVCWkROVmFuRjNlVEYyTkRsRlFtVjJhRzUyTVVWdk5HVllSaUo5In0="
	err := HandleChainTX(nodeID, a)
	if err != nil {
		panic(err)
	}

}

func Test_DumpTXID(t *testing.T) {
	a := "5fe5823c0d8b6d49f2ac99c90575566962ac3a14a6b2f1e7fe7ea1099b7b3bbd"
	value, raw := QueryChain(a)
	println(value)
	bc, _ := decodeChainTX(raw)
	print(string(bc.Payload))
}

//Use this to generate Order1
//curl -s -X POST "http://localhost:5556/v1/order1" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":\"\",\"extension\":{\"coin\":\"0\"}}"

func Test_All(t *testing.T) {
	DumpTXID("dea1396bce7890f85da7dc86b4ece5c4d372886ed08948eca6a0beee36c412e0")

}

func Test_1(t *testing.T) {
	txid := "473407b069ff917b110f38c36d5b9e5246b5ace5d82df38c5a188d5ac868cfec"
	DumpTXID(txid)
	ProcessTransactionID(txid)
}

func Test_2(t *testing.T) {
	txid := "586bc14b15a31999571c8188241beef046d3b78a9481ecee984e7c76a1d95112"
	DumpTXID(txid)
	ProcessTransactionID(txid)
}

func Test_3(t *testing.T) {
	txid := "5a48129fd272f2a8c57fdd96716a78c3be55a3cf811b179e82e54221d95ccbc4"
	DumpTXID(txid)
	ProcessTransactionID(txid)
}

//curl -s -X POST "http://localhost:5556/v1/order1" -H "accept: */*" -H "Content-Type: application/json" -d "{\"beneficiaryIDDocumentCID\":\"\",\"extension\":{\"coin\":\"0\"}}"

//DumpTXID -
func DumpTXID(txid string) {
	value, raw := QueryChain(txid)
	println(value)
	bc, _ := decodeChainTX(raw)
	println(string(bc.Payload))
	println()
}

//ProcessTransactionID -
func ProcessTransactionID(txid string) error {
	_, payload := QueryChain((txid))
	return HandleChainTX("", payload)
}

//HandleChainTX -
func HandleChainTX(myID string, tx string) error {
	blockChainTX, err := decodeChainTX(tx)
	if err != nil {
		return err
	}
	if err := callNextTX(nil, blockChainTX, "5556"); err != nil {
		return err
	}

	return nil
}
