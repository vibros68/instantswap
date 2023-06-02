package aptexplorer

import (
	"fmt"
	"strings"

	"gitlab.com/cryptopower/instantswap/blockexplorer"
	"gitlab.com/cryptopower/instantswap/blockexplorer/global/clients/blockexplorerclient"
	"gitlab.com/cryptopower/instantswap/blockexplorer/global/interfaces/idaemon"
	"gitlab.com/cryptopower/instantswap/blockexplorer/global/utils"
)

const (
	API_BASE = "https://fullnode.mainnet.aptoslabs.com/v1/"
	LIBNAME  = "aptoslabs"
)

func init() {
	blockexplorer.RegisterExplorer("APT", "", func(config blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New(config), nil
	})
}

type aptExplorer struct {
	conf      blockexplorer.Config
	client    *blockexplorerclient.Client
	apiKey    string
	apiSecret string
}

// New return an IBlockExplorer interface
func New(config blockexplorer.Config) *aptExplorer {
	client := blockexplorerclient.NewClient(API_BASE, LIBNAME, config.EnableOutput, nil)
	return &aptExplorer{client: client, conf: config}
}

func (e *aptExplorer) blockchainInfo() (*Blockchain, error) {
	r, err := e.client.Do("GET", "", "", false)
	if err != nil {
		return nil, err
	}
	var b Blockchain
	err = parseResponseData(r, &b)
	return &b, err
}

func (e *aptExplorer) getTxByHash(hash string) (*Transaction, error) {
	r, err := e.client.Do("GET", fmt.Sprintf("transactions/by_hash/%s", hash), "", false)
	if err != nil {
		return nil, err
	}
	var aptTx Transaction
	err = parseResponseData(r, &aptTx)
	if err != nil {
		return nil, err
	}
	return &aptTx, err
}

func (e *aptExplorer) GetTransaction(txId string) (tx *blockexplorer.ITransaction, err error) {
	aptTx, err := e.getTxByHash(txId)
	if err != nil {
		return nil, err
	}
	var blockHeight, confirmations int
	block, _ := e.getBlockByVersion(aptTx.Version)
	if block != nil {
		blockHeight = utils.StrToInt(block.BlockHeight)
	}
	blockchain, _ := e.blockchainInfo()
	if blockchain != nil {
		confirmations = utils.StrToInt(blockchain.BlockHeight) - blockHeight
	}
	vIns, vOuts := aptTx.getInOutPuts()
	return &blockexplorer.ITransaction{
		BlockHeight:         blockHeight,
		DoubleSpend:         false,
		Hash:                aptTx.Hash,
		Inputs:              vIns,
		LockTime:            0,
		Outputs:             vOuts,
		Rbf:                 false,
		Size:                0,
		Time:                0,
		TxIndex:             0,
		Version:             0,
		VinSz:               0,
		VoutSz:              0,
		Weight:              0,
		Confirmations:       confirmations,
		Seen:                true,
		Verified:            true,
		OrderedAmount:       0,
		BlockExplorerAmount: 0,
		MissingAmount:       0,
		MissingPercent:      0,
	}, err
}

func (e *aptExplorer) GetTxsForAddress(address string, limit int, viewKey string) (tx *blockexplorer.IRawAddrResponse, err error) {
	//r, err := e.client.Do("GET", fmt.Sprintf("accounts/%s/resources", address), "", false)
	return nil, fmt.Errorf("%s:not supported", LIBNAME)
}

func (e *aptExplorer) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
	tx = &blockexplorer.ITransaction{}
	aptTx, err := e.getTxByHash(verifier.TxId)
	if err != nil {
		return nil, err
	}
	tx.Hash = aptTx.Hash
	block, _ := e.getBlockByVersion(aptTx.Version)
	tx.BlockHeight = utils.StrToInt(block.BlockHeight)
	blockchain, _ := e.blockchainInfo()
	if blockchain != nil {
		tx.Confirmations = utils.StrToInt(blockchain.BlockHeight) - tx.BlockHeight
	}
	for _, event := range aptTx.Events {
		if event.Guid.AccountAddress == verifier.Address {
			tx.Seen = true
			tx.Verified = true
			tArr := strings.Split(event.Type, "::")
			if len(tArr) != 3 {
				continue
			}
			if tArr[1] == "coin" && (tArr[2] == "WithdrawEvent" || tArr[2] == "DepositEvent") {
				tx.BlockExplorerAmount = idaemon.Amount(utils.StrToInt(event.Data.Amount))
				coinAmount := tx.BlockExplorerAmount.ToCoin()
				tx.MissingAmount, _ = idaemon.NewAmount(verifier.Amount - coinAmount)
				tx.MissingPercent = 100 * (verifier.Amount - coinAmount) / verifier.Amount
				break
			}
		}
	}
	return
}

func (e *aptExplorer) PushTx(rawTxHash string) (result string, err error) {
	return "", fmt.Errorf("%s:not supported", LIBNAME)
}

func (e *aptExplorer) getBlockByVersion(version string) (*BlockInfo, error) {
	r, err := e.client.Do("GET", fmt.Sprintf("blocks/by_version/%s", version), "", false)
	if err != nil {
		return nil, err
	}
	var b BlockInfo
	err = parseResponseData(r, &b)
	return &b, err
}
