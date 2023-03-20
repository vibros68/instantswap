package ethplorer

import (
	"fmt"
	"strings"

	"code.cryptopower.dev/group/instantswap/blockexplorer"
)

type Tx struct {
	Hash          string        `json:"hash"`
	Timestamp     int           `json:"timestamp"`
	BlockNumber   int           `json:"blockNumber"`
	Confirmations int           `json:"confirmations"`
	Success       bool          `json:"success"`
	From          string        `json:"from"`
	To            string        `json:"to"`
	Value         int           `json:"value"`
	Input         string        `json:"input"`
	GasLimit      int           `json:"gasLimit"`
	GasUsed       int           `json:"gasUsed"`
	Logs          []TxLog       `json:"logs"`
	Operations    []TxOperation `json:"operations"`
}

type TxOperation struct {
	Timestamp       int       `json:"timestamp"`
	TransactionHash string    `json:"transactionHash"`
	Value           int       `json:"value,string"`
	IntValue        int64     `json:"intValue"`
	Type            string    `json:"type"`
	IsEth           bool      `json:"isEth"`
	Priority        int       `json:"priority"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	Addresses       []string  `json:"addresses"`
	UsdPrice        float64   `json:"usdPrice"`
	TokenInfo       TokenInfo `json:"tokenInfo"`
}

type TxLog struct {
}

type TokenInfo struct {
	Address        string `json:"address"`
	Name           string `json:"name"`
	Decimals       int    `json:"decimals,string"`
	Symbol         string `json:"symbol"`
	TotalSupply    string `json:"totalSupply"`
	Owner          string `json:"owner"`
	TxsCount       int    `json:"txsCount"`
	TransfersCount int    `json:"transfersCount"`
	LastUpdated    int    `json:"lastUpdated"`
	IssuancesCount int    `json:"issuancesCount"`
	Price          struct {
		Rate            float64 `json:"rate"`
		Diff            float64 `json:"diff"`
		Diff7D          float64 `json:"diff7d"`
		Ts              int     `json:"ts"`
		MarketCapUsd    float64 `json:"marketCapUsd"`
		AvailableSupply float64 `json:"availableSupply"`
		Volume24H       float64 `json:"volume24h"`
		VolDiff1        float64 `json:"volDiff1"`
		VolDiff7        float64 `json:"volDiff7"`
		VolDiff30       float64 `json:"volDiff30"`
		Diff30D         float64 `json:"diff30d"`
		Bid             float64 `json:"bid"`
		Currency        string  `json:"currency"`
	} `json:"price"`
	HoldersCount      int    `json:"holdersCount"`
	Description       string `json:"description"`
	Website           string `json:"website"`
	Image             string `json:"image"`
	EthTransfersCount int    `json:"ethTransfersCount"`
}

func (e *etherScan) generalTx(ethTx *Tx) (*blockexplorer.ITransaction, error) {
	var tx = &blockexplorer.ITransaction{
		BlockHeight:         ethTx.BlockNumber,
		DoubleSpend:         false,
		Hash:                ethTx.Hash,
		Inputs:              nil,
		LockTime:            0,
		Outputs:             nil,
		Rbf:                 false,
		Size:                0,
		Time:                ethTx.Timestamp,
		TxIndex:             0,
		Version:             0,
		VinSz:               0,
		VoutSz:              0,
		Weight:              0,
		Confirmations:       ethTx.Confirmations,
		Seen:                false,
		Verified:            false,
		OrderedAmount:       0,
		BlockExplorerAmount: 0,
		MissingAmount:       0,
		MissingPercent:      0,
	}
	if e.conf.Type == blockexplorer.NetworkTypeErc20 {
		var symbol = strings.ToUpper(e.conf.Symbol)
		var found bool
		for _, operation := range ethTx.Operations {
			if operation.TokenInfo.Symbol == symbol && operation.Type == "transfer" {
				found = true
				tx.Inputs = append(tx.Inputs, blockexplorer.IVIN{
					TxID:        ethTx.Hash,
					VOUT:        0,
					Tree:        0,
					AmountIn:    0,
					BlockIndex:  0,
					BlockHeight: 0,
				})
			}
		}
		if !found {
			return nil, fmt.Errorf("does not found operation for %s token", symbol)
		}
	}
	return tx, nil
}
