package xmrexplorer

import (
	"fmt"
	"github.com/vibros68/instantswap/blockexplorer/global/utils"

	"github.com/vibros68/instantswap/blockexplorer"
	"github.com/vibros68/instantswap/blockexplorer/global/clients/blockexplorerclient"
)

const (
	API_BASE                   = "https://xmrchain.net/api/" //  API endpoint
	DEFAULT_HTTPCLIENT_TIMEOUT = 30                          // HTTP client timeout
	LIBNAME                    = "monero"
)

func init() {
	blockexplorer.RegisterExplorer("XMR", "", func(conf blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New(conf), nil
	})
}

// New return a instanciate cryptopia struct
func New(conf blockexplorer.Config) *MoneroExplorer {
	client := blockexplorerclient.NewClient(API_BASE, LIBNAME, conf.EnableOutput, nil)
	return &MoneroExplorer{client: client}
}

type MoneroExplorer struct {
	client *blockexplorerclient.Client
}

func (z *MoneroExplorer) VerifyByAddress(req blockexplorer.AddressVerifyRequest) (vr *blockexplorer.VerifyResult, err error) {
	r, err := z.client.Do("GET", fmt.Sprintf("outputsblocks?address=%s&viewkey=%s&limit=%d&mempool=1",
		req.Address, req.ViewKey, 5), "", false)
	fmt.Println(string(r))
	var outputsBlocks OutputsBlocks
	if err = parseMoneroResponseData(r, &outputsBlocks); err != nil {
		return nil, err
	}
	for _, output := range outputsBlocks.Outputs {
		value := float64(output.Amount) / 1e12
		if utils.ApproximateCompare(value, req.Amount) {
			return &blockexplorer.VerifyResult{
				Seen:                true,
				Verified:            true,
				OrderedAmount:       req.Amount,
				BlockExplorerAmount: value,
				MissingAmount:       req.Amount - value,
				MissingPercent:      (req.Amount - value) / req.Amount,
			}, nil
		}
	}
	return nil, err
}

func (z *MoneroExplorer) GetTransaction(txId string) (*blockexplorer.ITransaction, error) {
	r, err := z.client.Do("GET", fmt.Sprintf("transaction/%s", txId), "", false)
	var tx Transaction
	if err = parseMoneroResponseData(r, &tx); err != nil {
		return nil, err
	}
	return tx.ITransaction(), nil
}
func (z *MoneroExplorer) GetTxsForAddress(address string, limit int, viewKey string) (account *blockexplorer.IRawAddrResponse, err error) {
	r, err := z.client.Do("GET", fmt.Sprintf("outputsblocks?address=%s&viewkey=%s&limit=%d&mempool=1", address, viewKey, limit), "", false)
	var outputsBlocks OutputsBlocks
	if err = parseMoneroResponseData(r, &outputsBlocks); err != nil {
		return nil, err
	}
	return outputsBlocks.IRawAddrResponse(), nil
}

// VerifyTransaction verifies transaction based on values passed in (params: txid, address (required), amount (required), createdAt(unix timestamp) )
func (z *MoneroExplorer) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
	r, err := z.client.Do("GET", fmt.Sprintf("outputs?txhash=%s&address=%s&viewkey=%s&txprove=0",
		verifier.TxId, verifier.Address, verifier.ViewKey), "", false)
	if err != nil {
		return nil, err
	}
	var txVerify TxVerifier
	fmt.Println(string(r))
	err = parseMoneroResponseData(r, &txVerify)
	if err != nil {
		return nil, err
	}
	return txVerify.ITransaction(verifier), nil
}

// PushTx pushes a raw tx hash
func (z *MoneroExplorer) PushTx(rawTxHash string) (result string, err error) {
	return "", fmt.Errorf("%s:error: PushTx is not supported yet... ", LIBNAME)
}
