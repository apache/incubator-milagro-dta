// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    fetchTxResponse, err := UnmarshalFetchTxResponse(bytes)
//    bytes, err = fetchTxResponse.Marshal()

package tendermint

import "encoding/json"

//UnmarshalChainQuery -
func UnmarshalChainQuery(data []byte) (ChainQuery, error) {
	var r ChainQuery
	err := json.Unmarshal(data, &r)
	return r, err
}

//Marshal -
func (r *ChainQuery) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

//ChainQuery -
type ChainQuery struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  Result `json:"result"`
}

//Result -
type Result struct {
	Response Response `json:"response"`
}

//Response -
type Response struct {
	Log   string `json:"log"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
