package query

import "encoding/json"

func UnmarshalQueryResponse(data []byte) (QueryResponse, error) {
	var r QueryResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *QueryResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type QueryResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  Result `json:"result"`
}

type Result struct {
	Txs        []Tx   `json:"txs"`
	TotalCount string `json:"total_count"`
}

type Tx struct {
	Hash     string   `json:"hash"`
	Height   string   `json:"height"`
	Index    int64    `json:"index"`
	TxResult TxResult `json:"tx_result"`
	Tx       string   `json:"tx"`
}

type TxResult struct {
	Code      int64       `json:"code"`
	Data      interface{} `json:"data"`
	Log       string      `json:"log"`
	Info      string      `json:"info"`
	GasWanted string      `json:"gasWanted"`
	GasUsed   string      `json:"gasUsed"`
	Events    []Event     `json:"events"`
	Codespace string      `json:"codespace"`
}

type Event struct {
	Type       string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
