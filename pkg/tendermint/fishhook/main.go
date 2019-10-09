package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/apache/incubator-milagro-dta/libs/documents"
	"github.com/apache/incubator-milagro-dta/libs/keystore"
	"github.com/apache/incubator-milagro-dta/pkg/api"
	"github.com/apache/incubator-milagro-dta/pkg/config"
	"github.com/apache/incubator-milagro-dta/pkg/identity"
	"github.com/hokaccha/go-prettyjson"
	tmclient "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/urfave/cli"
)

const (
	envMilagroHome = "MILAGRO_HOME"
	keysFile       = "keys"

	cmdInit   = "init"
	cmdDaemon = "daemon"
)

func main() {
	app := cli.NewApp()
	app.Name = "tmget"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Chris Morris",
			Email: "chris@morris.net",
		},
	}
	app.Copyright = "(c) 2019 Chris Morris"
	app.UsageText = `fishhook configdir nodeurl query skip
eg. fishhook /Users/john/.milagro 10.10,10,10:26657 "tag.recipient='Au1WipqVeTx9i2PV4UcCxmY6iQvA9RZXy88xJLRzafwc'" 3

configdir - the local directory where the DT-A configuration (eg. config.yaml, keys) are stored
nodeurl   - the host:port of a member Node of the Tendermint Network
query     - A query to filter the results by (enclosed query in double quotes and values in single quotes)
skip      - number of matches to skip before showing match and terminating
`
	app.Usage = `retrieve and parse a transaction in the Qredo DT-A Format from a Tendermint Blockchain
Note tags are case sensistive
`

	app.Action = func(c *cli.Context) error {
		folder := c.Args().Get(0)
		host := c.Args().Get(1)
		query := c.Args().Get(2)
		skip, err := strconv.Atoi(c.Args().Get(3))

		if err != nil {
			print("Invalid skip value\n")
			os.Exit(1)
		}

		if len(c.Args()) != 4 {
			print(app.UsageText)
			os.Exit(1)
			return nil
		}

		cfg, err := parseConfig(folder)
		if err != nil {
			print("Failed to open config")
			os.Exit(1)
		}

		keyStore, err := keystore.NewFileStore(filepath.Join(folder, keysFile))
		if err != nil {
			print("Fail to open keystore")
			os.Exit(1)
		}

		keyseed, err := keyStore.Get("seed")
		if err != nil {
			print("Fail to retrieve keyseed")
			os.Exit(1)
		}

		_, sikeSK, err := identity.GenerateSIKEKeys(keyseed)
		if err != nil {
			print("Fail to retrieve sikeSK")
			os.Exit(1)
		}
		blsPk, _, err := identity.GenerateBLSKeys(keyseed)
		if err != nil {
			print("Fail to retrieve blsSK")
			os.Exit(1)
		}

		//connect to Node
		tmClient := tmclient.NewHTTP(fmt.Sprintf("tcp://%s", host), "/websocket")
		if err := tmClient.Start(); err != nil {
			print("Failed to open websocket")
			os.Exit(1)
		}

		out, err := tmClient.Subscribe(context.Background(), "", query, 1000)
		if err != nil {
			print("Failed to subscribe to node")
			os.Exit(1)
		}

		matchCount := 0
		for {
			select {
			case result := <-out:
				matchCount++
				if matchCount != skip {
					continue
				}

				print("result")
				tx := result.Data.(tmtypes.EventDataTx).Tx
				nodeID := cfg.Node.NodeID
				payload := &api.BlockChainTX{}
				err := json.Unmarshal(tx, payload)
				if err != nil {
					print("Invalid Transaction received")
					os.Exit(1)
				}

				//dump
				order := &documents.OrderDoc{}
				err = documents.DecodeOrderDocument(payload.Payload, "", order, sikeSK, nodeID, blsPk)

				pp, _ := prettyjson.Marshal(order)
				fmt.Println(string(pp))
				os.Exit(0)
			}
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(name, defaultValue string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return v
}

func parseConfig(folder string) (*config.Config, error) {
	cfg, err := config.ParseConfig(folder)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

//Use - helper to remove warnings
func Use(vals ...interface{}) {
	for _, val := range vals {
		_ = val
	}
}
