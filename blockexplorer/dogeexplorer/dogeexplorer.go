package dogeexplorer

import (
	"encoding/json"
	"fmt"
	"github.com/vibros68/instantswap/blockexplorer"
	"github.com/vibros68/instantswap/blockexplorer/global/clients/blockexplorerclient"
	"github.com/vibros68/instantswap/blockexplorer/global/interfaces/idaemon"
)

func init() {
	blockexplorer.RegisterExplorer("DOGE", "", func(config blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New(config), nil
	})
}

const (
	API_BASE = "https://dogechain.info/api/v1/"
	LIBNAME  = "doge"
)

type dogeExplorer struct {
	conf      blockexplorer.Config
	client    *blockexplorerclient.Client
	apiKey    string
	apiSecret string
}

// New return an IBlockExplorer interface
func New(config blockexplorer.Config) *dogeExplorer {
	client := blockexplorerclient.NewClient(API_BASE, LIBNAME, config.EnableOutput, nil)
	return &dogeExplorer{client: client, conf: config}
}
func (d *dogeExplorer) VerifyByAddress(req blockexplorer.AddressVerifyRequest) (vr *blockexplorer.VerifyResult, err error) {
	txs, err := d.getTxsForAddress(req.Address)
	for _, tx := range txs {
		if tx.Value == req.Amount {
			return &blockexplorer.VerifyResult{
				Seen:                true,
				Verified:            true,
				OrderedAmount:       req.Amount,
				BlockExplorerAmount: tx.Value,
				MissingAmount:       0,
				MissingPercent:      0,
			}, nil
		}
	}
	return nil, fmt.Errorf("not found")
}
func (d *dogeExplorer) GetTransaction(txId string) (tx *blockexplorer.ITransaction, err error) {
	var response = struct {
		Res
		Tx Transaction `json:"transaction"`
	}{}
	r, err := d.client.Do("GET", fmt.Sprintf("transaction/%s", txId), "", false)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &response)
	if err != nil {
		return nil, err
	}
	if response.Success == 0 {
		return nil, fmt.Errorf(response.Error)
	}
	return response.Tx.tx(), nil
}
func (d *dogeExplorer) getTxsForAddress(address string) (txs []TxForAddress, err error) {
	var response = struct {
		Res
		Txs []TxForAddress `json:"transactions"`
	}{}
	r, err := d.client.Do("GET", fmt.Sprintf("address/transactions/%s", address), "", false)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(r, &response)
	if response.Success == 0 {
		return nil, fmt.Errorf(response.Error)
	}
	return response.Txs, err
}
func (d *dogeExplorer) GetTxsForAddress(address string, limit int, viewKey string) (tx *blockexplorer.IRawAddrResponse, err error) {
	tx = &blockexplorer.IRawAddrResponse{}
	txs, err := d.getTxsForAddress(address)
	if err != nil {
		return nil, err
	}

	var txsRaw = make([]blockexplorer.IRawAddrTx, len(txs))
	for i, tx := range txs {
		txsRaw[i] = blockexplorer.IRawAddrTx{
			BlockHeight:   tx.Block,
			Hash:          tx.Hash,
			Inputs:        nil,
			LockTime:      0,
			Outputs:       nil,
			RelayedBy:     "",
			Result:        0,
			Size:          0,
			Time:          tx.Time,
			TxIndex:       0,
			Version:       0,
			VinSz:         0,
			VoutSz:        0,
			Weight:        0,
			Confirmations: 0,
		}
	}
	tx.Txs = txsRaw
	return nil, err
}

// VerifyTransaction verifies transaction based on values passed in
func (d *dogeExplorer) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
	tx, err = d.GetTransaction(verifier.TxId)
	if err != nil {
		return nil, err
	}
	for _, output := range tx.Outputs {
		if output.Addresses[0] == verifier.Address {
			tx.Seen = true
			tx.Verified = tx.Confirmations > verifier.Confirms
			orderedAmount, _ := idaemon.NewAmount(verifier.Amount)
			tx.OrderedAmount = orderedAmount
			tx.BlockExplorerAmount = output.Value
			tx.MissingAmount = orderedAmount - output.Value
			tx.MissingPercent = (tx.MissingAmount.ToCoin() / orderedAmount.ToCoin()) * 100
		}
	}
	return tx, err
}

// PushTx pushes a raw tx hash
func (d *dogeExplorer) PushTx(rawTxHash string) (result string, err error) {
	return "", fmt.Errorf("not supported")
}
