package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/apache/incubator-milagro-dta/pkg/tendermint"
	"github.com/urfave/cli"
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
	app.UsageText = `tmget tag search [url]
eg. tmget tag.recipient=Nzw3127EaxPiiZahOH592sGhPnCPaYkzOSqEk 127.0.0.1:5556
    tmget tx Nzw3127EaxPiiZahOH592sGhPnCPaYkzOSqEk`

	app.Usage = `retrieve and parse a transaction in the Qredo DT-A Format from a Tendermint Blockchain
Note tags are case sensistive
Qredo DT-A uses:
  tag.recipient
  tag.senderid
  tag.reference
  tx

`

	app.ArgsUsage = "tx"

	app.Action = func(c *cli.Context) error {
		tag := c.Args().Get(0)
		lookup := "'" + c.Args().Get(1) + "'"

		url := c.Args().Get(2)
		//curl "localhost:26657/tx_search?query=\"tag.name='matts'\"&prove=true"
		if url == "" {
			url = "localhost:26657"
		}
		// if len(c.Args()) == 0 {
		// 	print(app.UsageText)
		// 	return nil
		// }

		fullUrl := "http://" + url + "/tx_search?query=\"" + tag + "=" + lookup + "\""
		print(fullUrl)

		resp, err := http.Get(fullUrl)
		if err != nil {
			// handle error
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		txResponse, err := tendermint.UnmarshalFetchTxResponse([]byte(body))

		txs := txResponse.TResult.Txs
		for r1, v := range txs {
			for r2, v1 := range v.TxResult.Events {
				for r3, v2 := range v1.Attributes {
					newkey, _ := base64.StdEncoding.DecodeString(v2.Key)
					txResponse.TResult.Txs[r1].TxResult.Events[r2].Attributes[r3].Key = string(newkey)
					newval, _ := base64.StdEncoding.DecodeString(v2.Value)
					txResponse.TResult.Txs[r1].TxResult.Events[r2].Attributes[r3].Value = string(newval)

					if string(newkey) == "key" {
						txResponse.TResult.Txs[r1].TxResult.Events[r2].Attributes[r3].Value = hex.EncodeToString(newval)
					}
				}
			}
		}

		x, err := txResponse.Marshal()

		var obj map[string]interface{}
		json.Unmarshal(x, &obj)

		//print(string(x))

		// Make a custom formatter with indent set
		f := colorjson.NewFormatter()
		f.Indent = 4

		// Marshall the Colorized JSON
		s, _ := f.Marshal(obj)
		fmt.Println(string(s))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
