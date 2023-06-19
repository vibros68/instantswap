package blockchair

import (
	"encoding/json"
	"fmt"

	"github.com/crypto-power/instantswap/blockexplorer"
	"github.com/crypto-power/instantswap/blockexplorer/global/clients/blockexplorerclient"
	"github.com/crypto-power/instantswap/blockexplorer/global/interfaces/idaemon"
)

const (
	API_BASE = "https://api.blockchair.com"
	LIBNAME  = "blockchair"
)

func init() {
	blockexplorer.RegisterExplorer("ZEC", "", func(conf blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New("zec", "zcash", conf), nil
	})
}

// New return a ClockChair client
func New(coinName, network string, conf blockexplorer.Config) *BlockChair {
	apiBase := fmt.Sprintf("%s/%s/dashboards/", API_BASE, network)
	client := blockexplorerclient.NewClient(apiBase, LIBNAME, conf.EnableOutput, nil)
	return &BlockChair{
		client:   client,
		coinName: coinName,
		network:  network,
		conf:     conf,
	}
}

type BlockChair struct {
	client    *blockexplorerclient.Client
	conf      blockexplorer.Config
	apiKey    string
	apiSecret string
	coinName  string
	network   string
}

func (b *BlockChair) getTx(txid string) (*TxWrapper, *Context, error) {
	r, err := b.client.Do("GET", fmt.Sprintf("transaction/%s", txid), "", false)
	if err != nil {
		return nil, nil, err
	}
	var txWrapperMap map[string]TxWrapper
	ctx, err := parseData(r, &txWrapperMap)
	if err != nil {
		return nil, nil, err
	}
	if txWrapperMap == nil {
		return nil, nil, fmt.Errorf("not found")
	}
	if txWrapper, ok := txWrapperMap[txid]; ok {
		return &txWrapper, ctx, nil
	}
	return nil, nil, fmt.Errorf("not found")
}

func (b *BlockChair) GetTransaction(txid string) (tx *blockexplorer.ITransaction, err error) {
	txW, ctx, err := b.getTx(txid)
	if err != nil {
		return nil, err
	}
	return b.generalTx(txW, ctx)
}

func (b *BlockChair) GetTxsForAddress(address string, limit int, viewKey string) (txs *blockexplorer.IRawAddrResponse, err error) {
	r, err := b.client.Do("GET", fmt.Sprintf("address/%s?transaction_details=true&omni=true", address), "", false)
	fmt.Println(string(r))
	if err != nil {
		return nil, err
	}
	var addrWrapperMap map[string]AddrWrapper
	ctx, err := parseData(r, &addrWrapperMap)
	if err != nil {
		return nil, err
	}
	if addrWrapperMap == nil {
		return nil, fmt.Errorf("not found")
	}
	if addrWrapper, ok := addrWrapperMap[address]; ok {
		return b.generalAddr(address, &addrWrapper, ctx), nil
	}
	return nil, fmt.Errorf("not found")
}

func (b *BlockChair) PushTx(txhash string) (res string, err error) {
	return "", fmt.Errorf("does not support PushTx")
}

func (b *BlockChair) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
	txW, ctx, err := b.getTx(verifier.TxId)
	if err != nil {
		return nil, err
	}
	tx, err = b.generalTx(txW, ctx)
	if err != nil {
		return nil, err
	}
	ordered, _ := idaemon.NewAmount(verifier.Amount)
	tx.OrderedAmount = ordered
	for _, out := range tx.Outputs {
		if out.Addresses[0] == verifier.Address {
			tx.Seen = true
			tx.Verified = true
			tx.BlockExplorerAmount = out.Value
			tx.MissingAmount = tx.OrderedAmount - tx.BlockExplorerAmount
			tx.MissingPercent = tx.MissingAmount.ToCoin() / tx.OrderedAmount.ToCoin() * 100
		}
	}
	return tx, nil
}

func parseData(data []byte, destination interface{}) (*Context, error) {
	var res jsonResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(res.Data, destination); err != nil {
		return nil, err
	}
	return &res.Context, nil
}
