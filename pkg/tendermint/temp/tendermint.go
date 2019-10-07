package tendermint

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/service"
)

//QueryChain the blockchain for an index
func QueryChain(index string) (string, string) {
	url := "http://" + node + "/abci_query?data=\"" + index + "\""
	resp, err := http.Get(url)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanBytes)
	t := ""
	for scanner.Scan() {
		t += scanner.Text()
		///fmt.Print(scanner.Text())
	}

	res, _ := UnmarshalChainQuery([]byte(t))

	val := res.Result.Response.Value
	decodeVal, _ := base64.StdEncoding.DecodeString(val)
	return string(decodeVal), val
}

//PostToChain - send TX data to the Blockchain
func PostToChain(tx *api.BlockChainTX, method string) (string, error) {
	//Create TX Hash

	tx.RecipientID = tx.RecipientID

	TXID := sha256.Sum256(tx.Payload)
	TXIDhex := hex.EncodeToString(TXID[:])
	tx.TXhash = TXID[:]

	//serialize the whole transaction
	serializedTX, _ := json.Marshal(tx)
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX)

	body := strings.NewReader("{\"jsonrpc\":\"2.0\",\"id\":\"anything\",\"method\":\"broadcast_tx_commit\",\"params\": {\"tx\": \"" + base64EncodedTX + "\"}}")
	url := "http://" + node + ""

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		print("Error posting to Blockchain")
		return "", err
	}
	req.Header.Set("Content-Type", "text/plain;")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		print("Error posting to Blockchain")
		return "", err
	}
	defer resp.Body.Close()
	fmt.Printf("POST TO CHAIN: METHOD:%s CALLS:%s  - TXID:%s\n", method, tx.Processor, TXIDhex)
	return TXIDhex, nil
}

//DecodeChainTX - Decode the On Chain TX into a BlockChainTX object
func decodeChainTX(payload string) (*api.BlockChainTX, error) {
	base64DecodedTX, _ := base64.StdEncoding.DecodeString(payload)
	tx := &api.BlockChainTX{}

	err := json.Unmarshal(base64DecodedTX, tx)
	if err != nil {
		return &api.BlockChainTX{}, err
	}
	return tx, nil
}

// TODO: remove
//DecodeChainTX - Decode the On Chain TX into a BlockChainTX object
// func decodeTX(payload string) (*api.BlockChainTX, string, error) {
// 	tx := &api.BlockChainTX{}
// 	parts := strings.SplitN(payload, "=", 2)
// 	if len(parts) != 2 {
// 		return &api.BlockChainTX{}, "", errors.New("Invalid TX payload")
// 	}
// 	hash := string(parts[0])
// 	err := json.Unmarshal([]byte(parts[1]), tx)
// 	if err != nil {
// 		return &api.BlockChainTX{}, "", err
// 	}
// 	return tx, hash, nil
// }

func callNextTX(svc service.Service, tx *api.BlockChainTX, listenPort string) error {
	switch tx.Processor {
	case "none":
		return nil
	case "dump":
		svc.Dump(tx)
	case "v1/fulfill/order":
		svc.FulfillOrder(tx)
	case "v1/order2":
		svc.Order2(tx)
	case "v1/fulfill/order/secret":
		svc.FulfillOrderSecret(tx)
	case "v1/order/secret2":
		svc.OrderSecret2(tx)

	default:
		return errors.New("Unknown processor")
	}
	return nil
}

// TODO: remove
// func unique(stringSlice []string) []string {
// 	keys := make(map[string]bool)
// 	list := []string{}
// 	for _, entry := range stringSlice {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			list = append(list, entry)
// 		}
// 	}
// 	return list
// }
