package tendermint

import "encoding/json"

func UnmarshalFetchTxResponse(data []byte) (FetchTxResponse, error) {
	var r FetchTxResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FetchTxResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalAttribute(data []byte) (Attribute, error) {
	var r Attribute
	print("hello")
	err := json.Unmarshal(data, &r)
	return r, err
}

type FetchTxResponse struct {
	ID      string  `json:"id"`
	Jsonrpc string  `json:"jsonrpc"`
	TResult TResult `json:"result"`
}

type TResult struct {
	TotalCount string `json:"total_count"`
	Txs        []Tx   `json:"txs"`
}

type Tx struct {
	Hash     string   `json:"hash"`
	Height   string   `json:"height"`
	Index    int64    `json:"index"`
	Tx       string   `json:"tx"`
	TxResult TxResult `json:"tx_result"`
}

type TxResult struct {
	Code      int64       `json:"code"`
	Codespace string      `json:"codespace"`
	Data      interface{} `json:"data"`
	Events    []Event     `json:"events"`
	GasUsed   string      `json:"gasUsed"`
	GasWanted string      `json:"gasWanted"`
	Info      string      `json:"info"`
	Log       string      `json:"log"`
}

type Event struct {
	Attributes []Attribute `json:"attributes"`
	Type       string      `json:"type"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
