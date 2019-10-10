package main

import (
	"context"
	"encoding/json"
	"os"

	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

type BlockChainTX struct {
	Processor              string
	SenderID               string
	RecipientID            string
	AdditionalRecipientIDs []string
	Payload                []byte
	TXhash                 []byte
	Tags                   map[string]string
}

func main() {

	subscribe(3)
}

//Subscribe to Websocket and add to queue
func subscribe(ignore int) error {
	client := tmclient.NewHTTP("tcp://34.246.173.153:26657", "/websocket")
	err := client.Start()
	if err != nil {
		print("Failed to start Tendermint HTTP client %s", err)
		return err
	}
	defer client.Stop()
	query := "tag.recipient='QmNpPp9wbBFUfBYCABrmKcaxd3eVxd2bGHfVQsEhRJPHjD'"
	//query := "tm.event = 'Tx'"

	out, err := client.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		print("Failed to subscribe to query %s %s", query, err)
		return err
	}

	print("Tendermint: Connected")

	for {
		select {
		case result := <-out:
			tx := result.Data.(tmtypes.EventDataTx).Tx
			payload := BlockChainTX{}
			err := json.Unmarshal(tx, &payload)
			if err != nil {
				print("******** Invalid TX - ignored")
				break
			}
			print(".")
			ignore--
			if ignore == 0 {
				print("Complete")
				os.Exit(0)
			}
		}

	}
	return nil
}



//Dump - used for debugging purpose, print the entire Encrypted Transaction
func  Dump(tx *BlockChainTX) error {
	nodeID := s.NodeID()
	txHashString := hex.EncodeToString(tx.TXhash)

	localIDDoc, err := common.RetrieveIDDocFromIPFS(s.Ipfs, nodeID)
	if err != nil {
		return err
	}

	_, _, _, sikeSK, err := common.RetrieveIdentitySecrets(s.Store, nodeID)
	if err != nil {
		return err
	}

	order := &documents.OrderDoc{}
	err = documents.FinalPrivateKey(tx.Payload, txHashString, order, sikeSK, nodeID, localIDDoc.BLSPublicKey)

	pp, _ := prettyjson.Marshal(order)
	fmt.Println(string(pp))

	return nil
}