package tendermint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/service"
	status "github.com/apache/incubator-milagro-dta/pkg/tendermint/status"
	"github.com/pkg/errors"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func catchUp(quene chan tmtypes.Tx, store *datastore.Store, logger *logger.Logger, nodeID string, listenPort string, height int) error {
	print("catch up")
	return nil
}

//Subscribe to Websocket and add to queue
func subscribeAndQueue(queueWaiting chan api.BlockChainTX, logger *logger.Logger, nodeID string, listenPort string) error {
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

			//check if this node is Sender - if so we don't need to process it
			if payload.SenderID == nodeID {
				break
			}

			//check if this node is in receipient list
			isRecipient := false
			for _, v := range payload.RecipientID {
				if v == nodeID {
					isRecipient = true
					break
				}
			}

			//If not in recipient list do nothing
			if isRecipient == false {
				logger.Info("******** Invalid Recipient - why are we receiving this TX?")
				break
			}

			//Add into the waitingQueue for later processing
			queueWaiting <- payload
			fmt.Printf("Incoming Transaction:%d \n", len(queueWaiting))

		case <-quit:
			os.Exit(0)
		}
	}
	return nil
}

func TXbyHash(TXHash string) (api.BlockChainTX, error) {
	client := tmclient.NewHTTP("tcp://"+node+"", "/websocket")
	query := fmt.Sprintf("tag.txhash='%s'", TXHash)
	result, err := client.TxSearch(query, true, 1, 1)

	if len(result.Txs) == 0 {
		return api.BlockChainTX{}, errors.New("Not found")
	}

	payload := api.BlockChainTX{}
	err = json.Unmarshal(result.Txs[0].Tx, &payload)

	_ = payload

	if err != nil {
		return payload, err
	}
	//
	// res := result.Txs[0]
	// tx := res.Tx
	return payload, nil

}

//loadAllHistoricTX - load the history for this node into a queue
func loadAllHistoricTX(start int, end int, txHistory []ctypes.ResultTx, nodeID string, listenPort string) error {
	//cycle through the historic transactions page by page
	//Get all transactions that claim to be from me - check signatures
	//Get all transactions that claim to be to me -

	client := tmclient.NewHTTP("tcp://"+node+"", "/websocket")
	currentPage := 1
	query := fmt.Sprintf("tag.recipient='%v' AND tag.sender='%v' AND tx.height>%d AND tx.height<=%d", nodeID, nodeID, start, end)
	numPerPage := 5

	for {
		result, err := client.TxSearch(query, true, currentPage, numPerPage)
		if err != nil {
			return errors.New("Failed to query chain for transaction history")
		}

		for _, tx := range result.Txs {
			txHistory = append(txHistory, *tx)
		}
		if currentPage*numPerPage > result.TotalCount {
			break
		}
		currentPage++
	}
	parseHistory(txHistory)
	return nil
}

func parseHistory(txHistory []ctypes.ResultTx) {
	txCount := len(txHistory)

	//loop backwards
	for i := txCount - 1; i >= 0; i-- {
		resTx := txHistory[i]
		tx := resTx.Tx

		//Decode TX into BlockchainTX Object
		payload := api.BlockChainTX{}
		err := json.Unmarshal(tx, &payload)
		if err != nil {
			msg := fmt.Sprintf("Invalid Transaction Hash:%v Height:%v Index:% \n", resTx.Hash, resTx.Height, resTx.Index)
			print(msg)
			continue
		}
		//Decode BlockchainTX.payload into Protobuffer Qredo
		// TODO:
		// Parse the incoming TX, check sig
		// If from self, can assume correct
		// builds transaction chains using previous transactionHash
		// Ensure every
		// Check recipient/sender in tags are correct
		//
		_ = payload
	}
	print("Finished loading - but not parsing the History\n")
}

func processTXQueue(svc service.Service, queue chan api.BlockChainTX, listenPort string) {
	print("Processing queue\n")
	for payload := range queue {
		//blockchainTX, txid, err := decodeTX(string(tx))
		//TXIDhex := hex.EncodeToString(payload.TXhash[:])
		//	logger.Info("Incoming TXHash:%s . Processor:%s", TXIDhex, payload.Processor)

		callNextTX(svc, &payload, listenPort)
	}
	print("Finished processing queue")
}

//Subscribe - Connect to the Tendermint websocket to collect events
func Subscribe(svc service.Service, store *datastore.Store, logger *logger.Logger, nodeID string, listenPort string) error {
	//Subscribe to channel
	//Get height

	latestStatus, _ := getChainStatus(node)
	currentBlockHeight, err := strconv.Atoi(latestStatus.Result.SyncInfo.LatestBlockHeight)

	if err != nil {
		return errors.New("Failed to obtain latest blockheight of Blockchain")
	}

	var processedToHeight int
	store.Get("chain", "height", &processedToHeight)

	//first catch up to Tip of chain
	var txHistory []ctypes.ResultTx
	queueWaiting := make(chan api.BlockChainTX, 1000)

	//while we are processessing the history save all new transactions in a queue for later

	go subscribeAndQueue(queueWaiting, logger, nodeID, listenPort)

	loadAllHistoricTX(processedToHeight, currentBlockHeight, txHistory, nodeID, listenPort)

	processTXQueue(svc, queueWaiting, listenPort)

	// var height int
	// store.Get("chain", "height", &height)

	// catchUp(queue, store, logger, nodeID, listenPort, height)
	// return nil

	// client := tmclient.NewHTTP("tcp://"+node+"", "/websocket")
	// //client.SetLogger(tmlogger)
	// err := client.Start()
	// if err != nil {
	// 	logger.Info("Failed to start Tendermint HTTP client %s", err)
	// 	return err
	// }
	// defer client.Stop()

	// //curl "34.246.173.153:26657/tx_search?query=\"tag.part=4%20AND%20tag.reference='579a2864-e100-11e9-aaf4-acde48001122'\""
	// query := "tag.recipient='" + nodeID + "'"
	// //query := "tm.event = 'Tx'"

	// out, err := client.Subscribe(context.Background(), "test", query, 1000)
	// if err != nil {
	// 	logger.Info("Failed to subscribe to query %s %s", query, err)
	// 	return err
	// }

	// logger.Info("Tendermint: Connected")

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// for {
	// 	select {
	// 	case result := <-out:
	// 		tx := result.Data.(tmtypes.EventDataTx).Tx
	// 		payload := api.BlockChainTX{}
	// 		err := json.Unmarshal(tx, &payload)
	// 		if err != nil {
	// 			logger.Info("******** Invalid TX - ignored")
	// 			break
	// 		}

	// 		//check if this node is Sender
	// 		if payload.SenderID == nodeID {
	// 			break
	// 		}

	// 		//check is receipient
	// 		isRecipient := false
	// 		for _, v := range payload.RecipientID {
	// 			if v == nodeID {
	// 				isRecipient = true
	// 				break
	// 			}
	// 		}

	// 		if isRecipient == false {
	// 			logger.Info("******** Invalid Recipient - why are we receiving this TX?")
	// 			break
	// 		}

	// 		//blockchainTX, txid, err := decodeTX(string(tx))
	// 		TXIDhex := hex.EncodeToString(payload.TXhash[:])
	// 		logger.Info("Incoming TXHash:%s . Processor:%s", TXIDhex, payload.Processor)

	// 		if payload.Processor == "NONE" {
	// 			DumpTX(&payload)
	// 		} else {
	// 			callNextTX(&payload, listenPort)
	// 		}

	// 		//print(blockchainTX)
	// 		// print(txid)

	// 		// print(string(xx))

	// 		// a := result.Data.(tmtypes.EventDataTx).Index
	// 		// b := result.Data.(tmtypes.EventDataTx)
	// 		// c := b.TxResult
	// 		// tx := c.Tx
	// 		// txdata := []byte(c.Tx)
	// 		// print(string(txdata))

	// 		// print(a)
	// 		// Use(c, b, tx)

	// 		//logger.Info("got tx","index", result.Data.(tmtypes.EventDataTx).Index)
	// 	case <-quit:
	// 		os.Exit(0)
	// 	}
	// }
	return nil
}

//Use - helper to remove warnings
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}

func getChainStatus(node string) (status.StatusResponse, error) {
	resp, err := http.Get("http://" + node + "/status")
	result := status.StatusResponse{}
	if err != nil {
		return result, err
	}
	json.NewDecoder(resp.Body).Decode((&result))
	return result, nil
}
