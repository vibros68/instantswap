package aptexplorer

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/vibros68/instantswap/blockexplorer"
	"github.com/vibros68/instantswap/blockexplorer/global/clients/blockexplorerclient"
	"github.com/vibros68/instantswap/blockexplorer/global/interfaces/idaemon"
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

func (a *aptExplorer) VerifyByAddress(req blockexplorer.AddressVerifyRequest) (vr *blockexplorer.VerifyResult, err error) {
	txs, err := a.getTxsForAddress(req.Address, 25, "")
	if err != nil {
		return nil, err
	}
	for _, tx := range txs {

		for _, event := range tx.Events {
			if event.Guid.AccountAddress == req.Address {
				tArr := strings.Split(event.Type, "::")
				if len(tArr) != 3 {
					continue
				}
				if tArr[1] == "coin" && (tArr[2] == "WithdrawEvent" || tArr[2] == "DepositEvent") {
					blockExplorerAmount := idaemon.Amount(event.Data.Amount)
					if blockExplorerAmount.ToCoin() == req.Amount {
						return &blockexplorer.VerifyResult{
							Seen:                true,
							Verified:            true,
							OrderedAmount:       req.Amount,
							BlockExplorerAmount: blockExplorerAmount.ToCoin(),
							MissingAmount:       0,
							MissingPercent:      0,
						}, nil
					}
				}
			}
		}
	}
	return nil, fmt.Errorf("tx not found")
}

func (a *aptExplorer) blockchainInfo() (*Blockchain, error) {
	r, err := a.client.Do("GET", "", "", false)
	if err != nil {
		return nil, err
	}
	var b Blockchain
	err = parseResponseData(r, &b)
	return &b, err
}

func (a *aptExplorer) getTxByHash(hash string) (*Transaction, error) {
	r, err := a.client.Do("GET", fmt.Sprintf("transactions/by_hash/%s", hash), "", false)
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

func (a *aptExplorer) getTxByVersion(version string) (*Transaction, error) {
	r, err := a.client.Do("GET", fmt.Sprintf("transactions/by_version/%s", version), "", false)
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

func (a *aptExplorer) GetTransaction(txId string) (tx *blockexplorer.ITransaction, err error) {
	aptTx, err := a.getTxByHash(txId)
	if err != nil {
		return nil, err
	}
	var blockHeight, confirmations int
	block, _ := a.getBlockByVersion(aptTx.Version)
	if block != nil {
		blockHeight = block.BlockHeight
	}
	blockchain, _ := a.blockchainInfo()
	if blockchain != nil {
		confirmations = blockchain.BlockHeight - blockHeight
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

func (a *aptExplorer) GetTxsForAddress(address string, limit int, viewKey string) (tx *blockexplorer.IRawAddrResponse, err error) {
	//r, err := e.client.Do("GET", fmt.Sprintf("accounts/%s/resources", address), "", false)
	return nil, fmt.Errorf("%s:not supported", LIBNAME)
}

func (a *aptExplorer) getTxsForAddress(address string, limit int, viewKey string) ([]*Transaction, error) {
	query := fmt.Sprintf(`{
	"operationName":"AccountTransactionsData",
	"variables":{"address":"%s","limit":%d,"offset":0},
	"query":"query AccountTransactionsData($address: String, $limit: Int, $offset: Int) {\n  address_version_from_move_resources(\n    where: {address: {_eq: $address}}\n    order_by: {transaction_version: desc}\n    limit: $limit\n    offset: $offset\n  ) {\n    transaction_version\n    __typename\n  }\n}"}`,
		address, limit)
	r, err := http.NewRequest("POST", "https://indexer.mainnet.aptoslabs.com/v1/graphql", bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	var obj struct {
		AddressVersionFromMoveResources []TxVersionResponse `json:"address_version_from_move_resources"`
	}
	err = parseDgraph(res.Body, &obj)
	if err != nil {
		return nil, err
	}
	var txs []*Transaction
	for _, txVer := range obj.AddressVersionFromMoveResources {
		aptTx, err := a.getTxByVersion(fmt.Sprintf("%d", txVer.TransactionVersion))
		if err == nil {
			txs = append(txs, aptTx)
		}
	}
	return txs, nil
}

func (a *aptExplorer) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
	tx = &blockexplorer.ITransaction{}
	aptTx, err := a.getTxByHash(verifier.TxId)
	if err != nil {
		return nil, err
	}
	tx.Hash = aptTx.Hash
	block, _ := a.getBlockByVersion(aptTx.Version)
	tx.BlockHeight = block.BlockHeight
	blockchain, _ := a.blockchainInfo()
	if blockchain != nil {
		tx.Confirmations = blockchain.BlockHeight - tx.BlockHeight
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
				tx.BlockExplorerAmount = idaemon.Amount(event.Data.Amount)
				coinAmount := tx.BlockExplorerAmount.ToCoin()
				tx.MissingAmount, _ = idaemon.NewAmount(verifier.Amount - coinAmount)
				tx.MissingPercent = 100 * (verifier.Amount - coinAmount) / verifier.Amount
				break
			}
		}
	}
	return
}

func (a *aptExplorer) PushTx(rawTxHash string) (result string, err error) {
	return "", fmt.Errorf("%s:not supported", LIBNAME)
}

func (a *aptExplorer) getBlockByVersion(version string) (*BlockInfo, error) {
	r, err := a.client.Do("GET", fmt.Sprintf("blocks/by_version/%s", version), "", false)
	if err != nil {
		return nil, err
	}
	var b BlockInfo
	err = parseResponseData(r, &b)
	return &b, err
}
