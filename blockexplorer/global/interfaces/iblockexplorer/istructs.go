package iblockexplorer

import (
	"github.com/vibros68/instantswap/blockexplorer/global/clients/blockexplorerclient"
	"github.com/vibros68/instantswap/blockexplorer/global/interfaces/idaemon"
)

type Client struct {
	client *blockexplorerclient.Client
}

type IVIN struct {
	Script   string `json:"script"`
	Sequence int    `json:"sequence"`
	Witness  string `json:"witness"`
	//only for dcrdata
	TxID        string         `json:"txid,omitempty"`
	VOUT        int            `json:"vout,omitempty"`
	Tree        int            `json:"tree,omitempty"`
	AmountIn    idaemon.Amount `json:"amountIn,omitempty"`
	BlockIndex  int            `json:"block_index,omitempty"`
	BlockHeight int            `json:"block_height,omitempty"`
}
type IVOUT struct {
	Addresses   []string       `json:"addr"`
	AddrTag     string         `json:"addr_tag,omitempty"`
	AddrTagLink string         `json:"addr_tag_link,omitempty"`
	N           int            `json:"n"`
	Script      string         `json:"script,omitempty"`
	Spent       bool           `json:"spent,omitempty"`
	TxIndex     int            `json:"tx_index"`
	Type        string         `json:"type"`
	Value       idaemon.Amount `json:"value"`
}
type ITransaction struct {
	BlockHeight   int     `json:"block_height,omitempty"`
	DoubleSpend   bool    `json:"double_spend,omitempty"`
	Hash          string  `json:"hash"`
	Inputs        []IVIN  `json:"inputs"`
	LockTime      int     `json:"lock_time"`
	Outputs       []IVOUT `json:"outputs"`
	Rbf           bool    `json:"rbf,omitempty"`
	Size          int     `json:"size"`
	Time          int     `json:"time"`
	TxIndex       int     `json:"tx_index,omitempty"`
	Version       int     `json:"ver"`
	VinSz         int     `json:"vin_sz,omitempty"`
	VoutSz        int     `json:"vout_sz,omitempty"`
	Weight        int     `json:"weight,omitempty"`
	Confirmations int     `json:"confirmations"` //calculated from getting latest block - blockheight

	//Internal vars for verification purposes
	Seen                bool           `json:"seen"` //tx has been seen on block explorer but not verified
	Verified            bool           `json:"verified"`
	OrderedAmount       idaemon.Amount `json:"ordered_amount"`
	BlockExplorerAmount idaemon.Amount `json:"blockexplorer_amount"`
	MissingAmount       idaemon.Amount `json:"missing_amount"`
	MissingPercent      float64        `json:"missing_percent"`
}
type ILatestBlock struct {
	BlockIndex int    `json:"block_index"`
	Hash       string `json:"hash"`
	Height     int    `json:"height"`
	Time       int    `json:"time"`
	TxIndexes  []int  `json:"txIndexes"`
}

type IRawAddrResponse struct {
	Address       string       `json:"address"`
	FinalBalance  int          `json:"final_balance,omitempty"`
	Hash160       string       `json:"hash160,omitempty"`
	NTx           int          `json:"n_tx,omitempty"`
	TotalReceived int          `json:"total_received,omitempty"`
	TotalSent     int          `json:"total_sent,omitempty"`
	Txs           []IRawAddrTx `json:"txs"`
}
type IRawAddrTx struct {
	BlockHeight   int              `json:"block_height,omitempty"`
	Hash          string           `json:"hash"`
	Inputs        []IRawAddrInput  `json:"inputs"`
	LockTime      int              `json:"lock_time"`
	Outputs       []IRawAddrOutput `json:"outputs"`
	RelayedBy     string           `json:"relayed_by"`
	Result        int              `json:"result"`
	Size          int              `json:"size"`
	Time          int              `json:"time"`
	TxIndex       int              `json:"tx_index,omitempty"`
	Version       int              `json:"ver"`
	VinSz         int              `json:"vin_sz,omitempty"`
	VoutSz        int              `json:"vout_sz,omitempty"`
	Weight        int              `json:"weight,omitempty"`
	Confirmations int              `json:"confirmations"`
}
type IRawAddrInput struct {
	PrevOut  IRawAddrOutput `json:"prev_out"`
	Script   string         `json:"script,omitempty"`
	Sequence int            `json:"sequence,omitempty"`
	Witness  string         `json:"witness,omitempty"`

	//DCR Specific Structure items
	TxID string `json:"txid"`
	VOUT int    `json:"vout"`
	Tree int    `json:"tree,omitempty"`
}
type IRawAddrOutput struct {
	Addresses []string       `json:"addr"`
	N         int            `json:"n"`
	Script    string         `json:"script,omitempty"`
	Spent     bool           `json:"spent,omitempty"`
	TxIndex   int            `json:"tx_index,omitempty"`
	Type      string         `json:"type,omitempty"`
	Value     idaemon.Amount `json:"value"`
}

type IPushTxResult struct {
	Success bool
	Message string
}
