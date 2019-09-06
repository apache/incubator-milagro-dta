package tendermint

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/incubator-milagro-dta/libs/logger"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
)

//Subscribe - Connect to the Tendermint websocket to collect events
func Subscribe(logger *logger.Logger) error {

	//tmlogger := log2.NewTMLogger(log.NewSyncWriter(os.Stdout))

	client := tmclient.NewHTTP("tcp://localhost:26657", "/websocket")
	//client.SetLogger(tmlogger)
	err := client.Start()
	if err != nil {
		logger.Info("Failed to start Tendermint HTTP client %s", err)
		return err
	}
	defer client.Stop()

	query := "tm.event = 'Tx'"
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

			blockchainTX, txid, err := decodeTX(string(tx))

			logger.Info("Incoming TX %s", txid)

			if err != nil {
				logger.Info("Invalid Incoming Transaction %s - %s:", err, string(tx))

			}

			if blockchainTX.Processor == "NONE" {
				DumpTX(blockchainTX)
			} else {
				callNextTX(blockchainTX)
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
