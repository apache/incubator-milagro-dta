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
	"os"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/apache/incubator-milagro-dta/pkg/api"
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

	tx.RecipientID = unique(tx.RecipientID)

	TXID := sha256.Sum256(tx.Payload)
	TXIDhex := hex.EncodeToString(TXID[:])
	tx.TXhash = TXID[:]

	//serialize the whole transaction
	serializedTX, _ := json.Marshal(tx)
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX)

	body := strings.NewReader("{\"jsonrpc\":\"2.0\",\"id\":\"anything\",\"method\":\"broadcast_tx_commit\",\"params\": {\"tx\": \"" + base64EncodedTX + "\"}}")
	req, err := http.NewRequest("POST", "http://"+node+"", body)
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

//HandleChainTX -
func HandleChainTX(myID string, tx string) error {
	blockChainTX, err := decodeChainTX(tx)
	if err != nil {
		return err
	}
	panic(nil)
	err = callNextTX(blockChainTX, "5556")
	if err != nil {
		return err
	}
	return nil
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

//DecodeChainTX - Decode the On Chain TX into a BlockChainTX object
func decodeTX(payload string) (*api.BlockChainTX, string, error) {
	tx := &api.BlockChainTX{}
	parts := strings.SplitN(payload, "=", 2)
	if len(parts) != 2 {
		return &api.BlockChainTX{}, "", errors.New("Invalid TX payload")
	}
	hash := string(parts[0])
	err := json.Unmarshal([]byte(parts[1]), tx)
	if err != nil {
		return &api.BlockChainTX{}, "", err
	}
	return tx, hash, nil
}

func callNextTX(tx *api.BlockChainTX, listenPort string) error {
	// recipient := tx.RecipientID
	// sender := tx.SenderID
	//payloadJSON := tx.Payload
	payloadString := string(tx.Payload)

	if tx.Processor == "NONE" {
		//The TX is information and doesn't require any further processing
		return nil
	}

	desintationURL := fmt.Sprintf("http://localhost"+listenPort+"/%s", tx.Processor)

	body := strings.NewReader(payloadString)
	req, err := http.NewRequest("POST", os.ExpandEnv(desintationURL), body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanBytes)
	t := ""
	for scanner.Scan() {
		t += scanner.Text()
		///fmt.Print(scanner.Text())
	}
	return nil
}

//DumpTX - Decode the Payload into JSON and displays the entire Blockchain TX unencoded
func DumpTX(bctx *api.BlockChainTX) {
	f := colorjson.NewFormatter()
	f.Indent = 4
	var payloadObj map[string]interface{}
	payload := bctx.Payload
	json.Unmarshal([]byte(payload), &payloadObj)
	jsonstring, _ := json.Marshal(bctx)
	var obj map[string]interface{}
	json.Unmarshal([]byte(jsonstring), &obj)
	obj["Payload"] = payloadObj
	s, _ := f.Marshal(obj)
	fmt.Println(string(s))
}

//DumpTXID -
func DumpTXID(txid string) {
	value, raw := QueryChain(txid)
	println(value)
	bc, _ := decodeChainTX(raw)
	println(string(bc.Payload))
	println()
}

//ProcessTransactionID -
func ProcessTransactionID(txid string) {
	_, payload := QueryChain((txid))
	err := HandleChainTX("", payload)
	if err != nil {
		panic(err)
	}
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}