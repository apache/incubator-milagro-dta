package tendermint

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

func catchUp(store *datastore.Store, logger *logger.Logger, nodeID string, listenPort string) error {

	return nil
}

//Subscribe - Connect to the Tendermint websocket to collect events
func Subscribe(store *datastore.Store, logger *logger.Logger, nodeID string, listenPort string) error {

	//first catch up to Tip of chain
	//catchUp(store, logger, nodeID, listenPort)

	client := tmclient.NewHTTP("tcp://"+node+"", "/websocket")
	//client.SetLogger(tmlogger)
	err := client.Start()
	if err != nil {
		logger.Info("Failed to start Tendermint HTTP client %s", err)
		return err
	}
	defer client.Stop()

	//curl "34.246.173.153:26657/tx_search?query=\"tag.part=4%20AND%20tag.reference='579a2864-e100-11e9-aaf4-acde48001122'\""
	query := "tag.recipient='" + nodeID + "'"
	//query := "tm.event = 'Tx'"

	out, err := client.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		logger.Info("Failed to subscribe to query %s %s", query, err)
		return err
	}

	logger.Info("Tendermint: Connected")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case result := <-out:
			tx := result.Data.(tmtypes.EventDataTx).Tx
			payload := api.BlockChainTX{}
			err := json.Unmarshal(tx, &payload)
			if err != nil {
				logger.Info("******** Invalid TX - ignored")
				break
			}

			//check if this node is Sender
			if payload.SenderID == nodeID {
				break
			}

			//check is receipient
			isRecipient := false
			for _, v := range payload.RecipientID {
				if v == nodeID {
					isRecipient = true
					break
				}
			}

			if isRecipient == false {
				logger.Info("******** Invalid Recipient - why are we receiving this TX?")
				break
			}

			//blockchainTX, txid, err := decodeTX(string(tx))
			TXIDhex := hex.EncodeToString(payload.TXhash[:])
			logger.Info("Incoming TXHash:%s . Processor:%s", TXIDhex, payload.Processor)

			if payload.Processor == "NONE" {
				DumpTX(&payload)
			} else {
				callNextTX(&payload, listenPort)
			}

			//print(blockchainTX)
			// print(txid)

			// print(string(xx))

			// a := result.Data.(tmtypes.EventDataTx).Index
			// b := result.Data.(tmtypes.EventDataTx)
			// c := b.TxResult
			// tx := c.Tx
			// txdata := []byte(c.Tx)
			// print(string(txdata))

			// print(a)
			// Use(c, b, tx)

			//logger.Info("got tx","index", result.Data.(tmtypes.EventDataTx).Index)
		case <-quit:
			os.Exit(0)
		}
	}
}

//Use - helper to remove warnings
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
