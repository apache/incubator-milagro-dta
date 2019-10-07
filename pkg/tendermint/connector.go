package tendermint

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-milagro-dta/libs/datastore"
	"github.com/apache/incubator-milagro-dta/libs/logger"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	status "github.com/apache/incubator-milagro-dta/pkg/tendermint/status"
	"github.com/pkg/errors"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	nodeConnectionTimeout = time.Second * 10
	txChanSize            = 1000
)

// ProcessTXFunc is executed on each incoming TX
type ProcessTXFunc func(tx *api.BlockChainTX) error

// NodeConnector is using external tendermint node to post and get transactions
type NodeConnector struct {
	nodeID     string
	tmNodeAddr string
	httpClient *http.Client
	tmClient   *tmclient.HTTP
	log        *logger.Logger
	store      *datastore.Store
}

// NewNodeConnector constructs a new Tendermint NodeConnector
func NewNodeConnector(tmNodeAddr string, nodeID string, store *datastore.Store, log *logger.Logger) (conn *NodeConnector, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("Initialize tendermint node connector: %v", r)
		}
	}()

	tmNodeAddr = strings.TrimRight(tmNodeAddr, "/")
	tmClient := tmclient.NewHTTP(fmt.Sprintf("tcp://%s", tmNodeAddr), "/websocket")
	if err := tmClient.Start(); err != nil {
		return nil, errors.Wrap(err, "Start tendermint client")
	}

	return &NodeConnector{
		tmNodeAddr: tmNodeAddr,
		nodeID:     nodeID,
		log:        log,
		store:      store,
		httpClient: &http.Client{
			Timeout: nodeConnectionTimeout,
		},
		tmClient: tmClient,
	}, nil

}

// Stop is performing clean-up
func (nc *NodeConnector) Stop() error {
	return nc.tmClient.Stop()
}

// GetTx retreives a transaction by hash
func (nc *NodeConnector) GetTx(txHash string) (*api.BlockChainTX, error) {
	query := fmt.Sprintf("tag.txhash='%s'", txHash)
	result, err := nc.tmClient.TxSearch(query, true, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(result.Txs) == 0 {
		return nil, errors.New("Transaction not found")
	}

	payload := &api.BlockChainTX{}
	if err := json.Unmarshal(result.Txs[0].Tx, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// PostTx posts a transaction to the chain and returns the transaction ID
func (nc *NodeConnector) PostTx(tx *api.BlockChainTX, method string) (txID string, err error) {
	txID = tx.CalcHash()

	//serialize the whole transaction
	serializedTX, err := json.Marshal(tx)
	if err != nil {
		return
	}
	base64EncodedTX := base64.StdEncoding.EncodeToString(serializedTX)

	// TODO: use net/rpc
	body := strings.NewReader(`{
		"jsonrpc": "2.0",
		"id": "anything",
		"method": "broadcast_tx_commit",
		"params": {
			"tx": "` + base64EncodedTX + `"}
	}`)
	url := "http://" + nc.tmNodeAddr

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", errors.Wrap(err, "post to blockchain node")
	}
	req.Header.Set("Content-Type", "text/plain;")

	resp, err := nc.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "post to blockchain node")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var respErr string
		if b, err := ioutil.ReadAll(resp.Body); err != nil {
			respErr = resp.Status
		} else {
			respErr = string(b)
		}

		return "", errors.Errorf("Post to blockchain node status %v: %v", resp.StatusCode, respErr)
	}

	nc.log.Debug("POST TO CHAIN: METHOD: %s CALLS: %s  - TXID: %s", method, tx.Processor, txID)

	return
}

// Subscribe connects to the Tendermint node and collect the events
func (nc *NodeConnector) Subscribe(ctx context.Context, processFn ProcessTXFunc) error {
	chainStatus, err := nc.getChainStatus()
	if err != nil {
		return err
	}

	currentBlockHeight, err := strconv.Atoi(chainStatus.Result.SyncInfo.LatestBlockHeight)
	if err != nil {
		return errors.Wrap(err, "Failed to obtain latest blockheight of Blockchain")
	}

	var processedToHeight int
	if err := nc.store.Get("chain", "height", &processedToHeight); err != nil {
		if err != datastore.ErrKeyNotFound {
			return errors.Wrap(err, "Get last processed block height")
		}
	}

	nc.log.Debug("Block height: Current: %v; Processed: %v", currentBlockHeight, processedToHeight)

	// create the transaction queue chan
	txQueue := make(chan *api.BlockChainTX, txChanSize)

	// Collect events
	if err := nc.subscribeAndQueue(ctx, txQueue); err != nil {
		return err
	}

	// TODO: load historicTX

	// Process events
	return nc.processTXQueue(ctx, txQueue, processFn)
}

func (nc *NodeConnector) subscribeAndQueue(ctx context.Context, txQueue chan *api.BlockChainTX) error {
	query := "tag.recipient='" + nc.nodeID + "'"

	out, err := nc.tmClient.Subscribe(context.Background(), "test", query, 1000)
	if err != nil {
		return errors.Wrapf(err, "Failed to subscribe to query %s", query)
	}

	go func() {
		for {
			select {
			case result := <-out:
				tx := result.Data.(tmtypes.EventDataTx).Tx
				payload := &api.BlockChainTX{}
				err := json.Unmarshal(tx, payload)
				if err != nil {
					nc.log.Debug("IGNORED TX - Invalid!")
					break
				}

				//check if this node is in receipient list
				if payload.RecipientID != nc.nodeID {
					nc.log.Debug("IGNORED TX! Recipient not match the query! (%v != %v)", payload.RecipientID, nc.nodeID)
					break
				}

				//Add into the waitingQueue for later processing
				txQueue <- payload
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (nc *NodeConnector) processTXQueue(ctx context.Context, txQueue chan *api.BlockChainTX, processFn ProcessTXFunc) error {
	for {
		select {
		case tx := <-txQueue:
			if err := processFn(tx); err != nil {
				// TODO: errors block processing the queue
				return err
			}
			// TODO: store the last block height
		case <-ctx.Done():
			return nil
		}
	}
}

func (nc *NodeConnector) getChainStatus() (*status.StatusResponse, error) {
	url := fmt.Sprintf("http://%s/status", nc.tmNodeAddr)
	resp, err := nc.httpClient.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Get node status")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Get node status status code: %v", resp.StatusCode)
	}

	status := &status.StatusResponse{}
	if err := json.NewDecoder(resp.Body).Decode((&status)); err != nil {
		return nil, errors.Wrap(err, "Invalid node status response")
	}

	return status, nil
}
