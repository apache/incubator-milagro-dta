package status

import "encoding/json"

func UnmarshalStatusResponse(data []byte) (StatusResponse, error) {
	var r StatusResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *StatusResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type StatusResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  Result `json:"result"`
}

type Result struct {
	NodeInfo      NodeInfo      `json:"node_info"`
	SyncInfo      SyncInfo      `json:"sync_info"`
	ValidatorInfo ValidatorInfo `json:"validator_info"`
}

type NodeInfo struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version"`
	ID              string          `json:"id"`
	ListenAddr      string          `json:"listen_addr"`
	Network         string          `json:"network"`
	Version         string          `json:"version"`
	Channels        string          `json:"channels"`
	Moniker         string          `json:"moniker"`
	Other           Other           `json:"other"`
}

type Other struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

type ProtocolVersion struct {
	P2P   string `json:"p2p"`
	Block string `json:"block"`
	App   string `json:"app"`
}

type SyncInfo struct {
	LatestBlockHash   string `json:"latest_block_hash"`
	LatestAppHash     string `json:"latest_app_hash"`
	LatestBlockHeight string `json:"latest_block_height"`
	LatestBlockTime   string `json:"latest_block_time"`
	CatchingUp        bool   `json:"catching_up"`
}

type ValidatorInfo struct {
	Address     string `json:"address"`
	PubKey      PubKey `json:"pub_key"`
	VotingPower string `json:"voting_power"`
}

type PubKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
