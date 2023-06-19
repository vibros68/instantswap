package ethplorer

import (
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/crypto-power/instantswap/blockexplorer"
	"github.com/crypto-power/instantswap/blockexplorer/global/clients/blockexplorerclient"
	"github.com/crypto-power/instantswap/blockexplorer/global/interfaces/idaemon"
)

const (
	API_BASE = "https://api.ethplorer.io/"
	LIBNAME  = "ethplorer"
)

func init() {
	blockexplorer.RegisterExplorer("", blockexplorer.NetworkTypeErc20, func(config blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New(config)
	})
}

type etherScan struct {
	conf   *blockexplorer.Config
	client *blockexplorerclient.Client
}

func New(conf blockexplorer.Config) (*etherScan, error) {
	client := blockexplorerclient.NewClient(API_BASE, LIBNAME, conf.EnableOutput, func(r *http.Request) {

	})
	return &etherScan{
		client: client,
		conf:   &conf,
	}, nil
}

func (e *etherScan) getTx(txId string) (*Tx, error) {
	r, err := e.client.Do("GET", fmt.Sprintf("getTxInfo/%s?apiKey=freekey", txId), "", false)
	if err != nil {
		return nil, err
	}
	var ethTx Tx
	err = parse(r, &ethTx)
	return &ethTx, err
}

func (e *etherScan) GetTransaction(txId string) (tx *blockexplorer.ITransaction, err error) {
	ethTx, err := e.getTx(txId)
	if err != nil {
		return nil, err
	}
	fmt.Println(ethTx.Operations[0].TokenInfo.Decimals)
	return e.generalTx(ethTx)
}
func (e *etherScan) GetTxsForAddress(address string, limit int, viewKey string) (tx *blockexplorer.IRawAddrResponse, err error) {
	r, err := e.client.Do("GET", fmt.Sprintf("getAddressHistory/%s?apiKey=freekey", address), "", false)
	if err != nil {
		return nil, err
	}
	var addrInfo struct {
		Operations []TxOperation `json:"operations"`
	}
	err = parse(r, &addrInfo)
	if err != nil {
		return nil, err
	}
	tx = &blockexplorer.IRawAddrResponse{
		Address:       address,
		FinalBalance:  0,
		Hash160:       "",
		NTx:           0,
		TotalReceived: 0,
		TotalSent:     0,
		Txs:           nil,
	}
	if e.conf.Type == blockexplorer.NetworkTypeErc20 {
		var symbol = strings.ToUpper(e.conf.Symbol)
		for _, operation := range addrInfo.Operations {
			if operation.TokenInfo.Symbol == symbol {
				explorerAmount := float64(operation.IntValue) / math.Pow(10, float64(operation.TokenInfo.Decimals))
				amount, _ := idaemon.NewAmount(explorerAmount)
				tx.Txs = append(tx.Txs, blockexplorer.IRawAddrTx{
					BlockHeight: 0,
					Hash:        operation.TransactionHash,
					Inputs:      []blockexplorer.IRawAddrInput{},
					LockTime:    0,
					Outputs: []blockexplorer.IRawAddrOutput{
						{
							Value:     amount,
							Addresses: []string{operation.To},
						},
					},
					RelayedBy:     "",
					Result:        0,
					Size:          0,
					Time:          operation.Timestamp,
					TxIndex:       0,
					Version:       0,
					VinSz:         0,
					VoutSz:        0,
					Weight:        0,
					Confirmations: 0,
				})
			}
		}
	}
	return nil, nil
}
func (e *etherScan) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
	ethTx, err := e.getTx(verifier.TxId)
	if err != nil {
		return nil, err
	}
	tx, err = e.generalTx(ethTx)
	if err != nil {
		return nil, err
	}
	orderedAmount, err := idaemon.NewAmount(verifier.Amount)
	if err != nil {
		return nil, err
	}
	tx.Seen = verifier.Address == ethTx.To
	tx.Verified = verifier.Address != ethTx.To
	if e.conf.Type == blockexplorer.NetworkTypeErc20 {
		var symbol = strings.ToUpper(e.conf.Symbol)
		var found bool
		for _, operation := range ethTx.Operations {
			if operation.TokenInfo.Symbol == symbol && operation.Type == "transfer" {
				found = true
				tx.OrderedAmount = orderedAmount
				explorerAmount := float64(operation.IntValue) / math.Pow(10, float64(operation.TokenInfo.Decimals))
				tx.BlockExplorerAmount, _ = idaemon.NewAmount(explorerAmount)
				tx.MissingAmount = tx.OrderedAmount - tx.BlockExplorerAmount
				tx.MissingPercent = 100 * float64(tx.MissingAmount) / float64(tx.OrderedAmount)
			}
		}
		tx.Verified = found
	}
	return tx, nil
}
func (e *etherScan) PushTx(rawTxHash string) (result string, err error) {
	return "", fmt.Errorf("not supported")
}
