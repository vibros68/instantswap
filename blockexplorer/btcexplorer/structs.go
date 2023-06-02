package btcexplorer

import (
	"encoding/json"

	"gitlab.com/cryptopower/instantswap/blockexplorer/global/interfaces/idaemon"
)

type jsonResponse struct {
	Success bool            `json:"Success"`
	Message string          `json:"Message"`
	Result  json.RawMessage `json:"Data"`
}
type jsonPrivResponse struct {
	Success bool            `json:"Success"`
	Message string          `json:"Error"`
	Result  json.RawMessage `json:"Data"`
}

type VIN struct {
	Script   string `json:"script"`
	Sequence int    `json:"sequence"`
	Witness  string `json:"witness"`
}
type VOUT struct {
	Addr        string         `json:"addr"`
	AddrTag     string         `json:"addr_tag"`
	AddrTagLink string         `json:"addr_tag_link"`
	N           int            `json:"n"`
	Script      string         `json:"script"`
	Spent       bool           `json:"spent"`
	TxIndex     int            `json:"tx_index"`
	Type        int            `json:"type"`
	Value       idaemon.Amount `json:"value"`
}
type Transaction struct {
	BlockHeight   int    `json:"block_height"`
	DoubleSpend   bool   `json:"double_spend"`
	Hash          string `json:"hash"`
	Inputs        []VIN  `json:"inputs"`
	LockTime      int    `json:"lock_time"`
	Out           []VOUT `json:"out"`
	Rbf           bool   `json:"rbf"`
	RelayedBy     string `json:"relayed_by"`
	Size          int    `json:"size"`
	Time          int    `json:"time"`
	TxIndex       int    `json:"tx_index"`
	Ver           int    `json:"ver"`
	VinSz         int    `json:"vin_sz"`
	VoutSz        int    `json:"vout_sz"`
	Weight        int    `json:"weight"`
	Confirmations int    `json:"confirmations"` //calculated from getting latest block - blockheight
}
type LatestBlock struct {
	BlockIndex int    `json:"block_index"`
	Hash       string `json:"hash"`
	Height     int    `json:"height"`
	Time       int    `json:"time"`
	TxIndexes  []int  `json:"txIndexes"`
}

type RawAddrResponse struct {
	Address       string       `json:"address"`
	FinalBalance  int          `json:"final_balance"`
	Hash160       string       `json:"hash160"`
	NTx           int          `json:"n_tx"`
	TotalReceived int          `json:"total_received"`
	TotalSent     int          `json:"total_sent"`
	Txs           []RawAddrTxs `json:"txs"`
}
type RawAddrTxs struct {
	BlockHeight int             `json:"block_height"`
	Hash        string          `json:"hash"`
	Inputs      []RawAddrInput  `json:"inputs"`
	LockTime    int             `json:"lock_time"`
	Outputs     []RawAddrOutput `json:"out"`
	RelayedBy   string          `json:"relayed_by"`
	Result      int             `json:"result"`
	Size        int             `json:"size"`
	Time        int             `json:"time"`
	TxIndex     int             `json:"tx_index"`
	Version     int             `json:"ver"`
	VinSz       int             `json:"vin_sz"`
	VoutSz      int             `json:"vout_sz"`
	Weight      int             `json:"weight"`
}
type RawAddrInput struct {
	PrevOut struct {
		Addr    string `json:"addr"`
		N       int    `json:"n"`
		Script  string `json:"script"`
		Spent   bool   `json:"spent"`
		TxIndex int    `json:"tx_index"`
		Type    int    `json:"type"`
		Value   int    `json:"value"`
	} `json:"prev_out"`
	Script   string `json:"script"`
	Sequence int    `json:"sequence"`
	Witness  string `json:"witness"`
}
type RawAddrOutput struct {
	Address string `json:"addr"`
	N       int    `json:"n"`
	Script  string `json:"script"`
	Spent   bool   `json:"spent"`
	TxIndex int    `json:"tx_index"`
	Type    int    `json:"type"`
	Value   int    `json:"value"`
}
