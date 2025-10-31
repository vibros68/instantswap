package btcexplorer

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vibros68/instantswap/blockexplorer"
	"github.com/vibros68/instantswap/blockexplorer/global/interfaces/idaemon"

	"github.com/vibros68/instantswap/blockexplorer/global/clients/blockexplorerclient"
)

func init() {
	blockexplorer.RegisterExplorer("BTC", "", func(conf blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New(conf), nil
	})
}

const (
	API_BASE                   = "https://blockchain.info/" //  API endpoint
	DEFAULT_HTTPCLIENT_TIMEOUT = 30                         // HTTP client timeout
	LIBNAME                    = "btcexplorer"
)

// New return a instanciate cryptopia struct
func New(conf blockexplorer.Config) *BlockChainInfo {
	client := blockexplorerclient.NewClient(API_BASE, LIBNAME, conf.EnableOutput, nil)
	return &BlockChainInfo{client: client}
}

// handleErr gets JSON response from the API and deal with error
func handleErr(r jsonResponse) error {
	if !r.Success {
		return errors.New(r.Message)
	}
	return nil
}
func handlePrivErr(r jsonPrivResponse) error {
	if !r.Success {
		return errors.New(r.Message)
	}
	return nil
}

// represent a * client
type BlockChainInfo struct {
	client    *blockexplorerclient.Client
	apiKey    string
	apiSecret string
}

// SetDebug set enable/disable http request/response dump
func (c *BlockChainInfo) SetDebug(enable bool) {
	c.client.Debug = enable
}

func (c *BlockChainInfo) VerifyByAddress(req blockexplorer.AddressVerifyRequest) (vr *blockexplorer.VerifyResult, err error) {
	rawAddr, err := c.getTxsForAddress(req.Address, 25)
	if err != nil {
		return nil, err
	}
	for _, tx := range rawAddr.Txs {
		for _, out := range tx.Outputs {
			value := idaemon.Amount(out.Value)
			if value.ToCoin() == req.Amount {
				return &blockexplorer.VerifyResult{
					Seen:                true,
					Verified:            true,
					OrderedAmount:       req.Amount,
					BlockExplorerAmount: value.ToCoin(),
					MissingAmount:       0,
					MissingPercent:      0,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}

// GetTransaction returns decoded transaction from api
func (c *BlockChainInfo) GetTransaction(txid string) (tx *blockexplorer.ITransaction, err error) {
	r, err := c.client.Do("GET", "rawtx/"+txid, "", false)
	if err != nil {
		return
	}
	var tmp Transaction
	if err = json.Unmarshal(r, &tmp); err != nil {
		return
	}

	//get latest block to get our confirmations
	latestBlock, err := c.GetLatestBlock()
	if err != nil {
		return
	}

	//confirmations for this tx
	tmp.Confirmations = latestBlock.Height - tmp.BlockHeight

	tx = &blockexplorer.ITransaction{
		Confirmations: tmp.Confirmations,
		BlockHeight:   tmp.BlockHeight,
		DoubleSpend:   tmp.DoubleSpend,
		Hash:          tmp.Hash,
		LockTime:      tmp.LockTime,
		Rbf:           tmp.Rbf,
		Size:          tmp.Size,
		Time:          tmp.Time,
		TxIndex:       tmp.TxIndex,
		Version:       tmp.Ver,
		VinSz:         tmp.VinSz,
		VoutSz:        tmp.VoutSz,
		Weight:        tmp.Weight,
	}
	for _, v := range tmp.Inputs {
		tmpIn := blockexplorer.IVIN{
			Script:   v.Script,
			Sequence: v.Sequence,
			Witness:  v.Witness,
		}
		tx.Inputs = append(tx.Inputs, tmpIn)
	}
	for _, v := range tmp.Out {
		addresses := []string{v.Addr}
		tmpOut := blockexplorer.IVOUT{
			Script:      v.Script,
			Addresses:   addresses,
			AddrTag:     v.AddrTag,
			AddrTagLink: v.AddrTagLink,
			N:           v.N,
			Spent:       v.Spent,
			TxIndex:     v.TxIndex,
			Type:        fmt.Sprintf("%v", v.Type),
			Value:       v.Value,
		}
		tx.Outputs = append(tx.Outputs, tmpOut)
	}
	//err = errors.New("test failure check error")
	return
}

func (c *BlockChainInfo) getTxsForAddress(address string, limit int) (txs *RawAddrResponse, err error) {
	r, err := c.client.Do("GET", fmt.Sprintf("rawaddr/%s?&limit=%v", address, limit), "", false)
	if err != nil {
		return
	}
	var rawaddr RawAddrResponse
	if err = json.Unmarshal(r, &rawaddr); err != nil {
		errStr := fmt.Sprintf(LIBNAME+":error: could not parse response for address %s msg: %s", address, err.Error())
		err = errors.New(errStr)
		return
	}
	return &rawaddr, nil
}

// GetTransactionsForAddress
func (c *BlockChainInfo) GetTxsForAddress(address string, limit int, viewKey string) (txs *blockexplorer.IRawAddrResponse, err error) {
	tmp, err := c.getTxsForAddress(address, limit)

	if err != nil {
		return nil, err
	}

	//get latest block to get our confirmations
	latestBlock, err := c.GetLatestBlock()
	if err != nil {
		return
	}

	txs = &blockexplorer.IRawAddrResponse{
		Address:       tmp.Address,
		FinalBalance:  tmp.FinalBalance,
		Hash160:       tmp.Hash160,
		NTx:           tmp.NTx,
		TotalReceived: tmp.TotalReceived,
		TotalSent:     tmp.TotalSent,
	}
	var allTxs []blockexplorer.IRawAddrTx
	//gather txs for this address and format them for interface
	for _, v := range tmp.Txs {

		tmpTx := blockexplorer.IRawAddrTx{
			BlockHeight:   v.BlockHeight,
			Hash:          v.Hash,
			LockTime:      v.LockTime,
			RelayedBy:     v.RelayedBy,
			Result:        v.Result,
			Size:          v.Size,
			Time:          v.Time,
			TxIndex:       v.TxIndex,
			Version:       v.Version,
			VinSz:         v.VinSz,
			VoutSz:        v.VoutSz,
			Weight:        v.Weight,
			Confirmations: latestBlock.Height - v.BlockHeight,
		}
		var tmpInputs []blockexplorer.IRawAddrInput
		for _, w := range v.Inputs {

			tmpIn := blockexplorer.IRawAddrInput{
				Script:   w.Script,
				Sequence: w.Sequence,
				Witness:  w.Witness,
			}

			tmpAmount := idaemon.Amount(w.PrevOut.Value)

			addresses := []string{w.PrevOut.Addr}
			tmpPrevOutput := blockexplorer.IRawAddrOutput{
				Addresses: addresses,
				N:         w.PrevOut.N,
				Script:    w.PrevOut.Script,
				Spent:     w.PrevOut.Spent,
				TxIndex:   w.PrevOut.TxIndex,
				Type:      fmt.Sprintf("%v", w.PrevOut.Type),
				Value:     tmpAmount,
			}
			tmpIn.PrevOut = tmpPrevOutput

			tmpInputs = append(tmpInputs, tmpIn)

		}
		tmpTx.Inputs = tmpInputs //assign inputs to current tx item

		var tmpOuputs []blockexplorer.IRawAddrOutput
		for _, w := range v.Outputs {

			tmpAmount := idaemon.Amount(w.Value)

			addresses := []string{w.Address}
			tmpOut := blockexplorer.IRawAddrOutput{
				Addresses: addresses,
				N:         w.N,
				Script:    w.Script,
				Spent:     w.Spent,
				TxIndex:   w.TxIndex,
				Type:      fmt.Sprintf("%v", w.Type),
				Value:     tmpAmount,
			}
			tmpOuputs = append(tmpOuputs, tmpOut)
		}
		tmpTx.Outputs = tmpOuputs //assign ouputs to current tx item

		allTxs = append(allTxs, tmpTx) //append additional tx to list of all txs
	}

	txs.Txs = allTxs //assign all txs to main return
	return
}

// PushTx
func (c *BlockChainInfo) PushTx(txhash string) (res string, err error) {
	r, err := c.client.Do("POST", fmt.Sprintf("pushtx?tx=%s", txhash), "", false)
	if err != nil {
		return
	}

	res = fmt.Sprintf("%s", r)

	return
}

// VerifyTransaction verifies transaction based on values passed in (params: txid, address, amount, createdAt )
func (c *BlockChainInfo) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
	tx = new(blockexplorer.ITransaction)
	if verifier.Address == "" {
		err = errors.New(LIBNAME + ":error: address is blank so tx cannot be verified")
		return
	}
	if verifier.Amount == 0 {
		errStr := fmt.Sprintf(LIBNAME+":error: amount is %.8f so tx cannot be verified", verifier.Amount)
		err = errors.New(errStr)
		return
	}

	if verifier.TxId != "" && verifier.Address != "" { //verify tx if txid is available
		txInfo, err := c.GetTransaction(verifier.TxId)
		if err != nil {
			return tx, err
		}

		for _, v := range txInfo.Outputs {
			for _, w := range v.Addresses {

				if w == verifier.Address {
					//fmt.Printf("seen")

					//tx has been seen on block explorer but still only has 0 confirmations
					if txInfo.Confirmations < verifier.Confirms {
						tx.Seen = true
						errStr := fmt.Sprintf("seen, waiting for confirms (%v/%v)", txInfo.Confirmations, verifier.Confirms)
						err = errors.New(errStr)
						return tx, err
					}

					orderedAmount, err := idaemon.NewAmount(verifier.Amount)
					if err != nil {
						errStr := fmt.Sprintf(LIBNAME+":error: orderedAmount %s", err.Error())
						err = errors.New(errStr)
						return tx, err
					}
					var missingAmount idaemon.Amount
					missingAmount = v.Value - orderedAmount
					if v.Value == orderedAmount {
						missingAmount = 0
					}
					missingPercent := (missingAmount.ToCoin() / v.Value.ToCoin()) * 100

					tx = txInfo
					tx.MissingAmount = missingAmount
					tx.MissingPercent = missingPercent
					tx.Verified = true
					tx.OrderedAmount = orderedAmount
					tx.BlockExplorerAmount = v.Value
					tx.Seen = true
					return tx, err

				} else {
					//fmt.Printf("\nadditional vout amount: %v to: %v", v.Value, w)
				}
			}
		}
	} else if verifier.Address != "" { //verify tx on blockchain based on address history for address var
		txInfo, err := c.GetTxsForAddress(verifier.Address, 10, "")
		if err != nil {
			return tx, err
		}

		for _, u := range txInfo.Txs {
			for _, v := range u.Outputs {
				for _, w := range v.Addresses {
					//fmt.Printf("\n%s", w)
					if w == verifier.Address && int64(u.Time) >= int64(verifier.CreatedAt) {
						//fmt.Printf("seen")

						//tx has been seen on block explorer but still only has 0 confirmations
						if u.Confirmations < verifier.Confirms {
							tx.Seen = true
							errStr := fmt.Sprintf("seen, waiting for confirms (%v/%v)", u.Confirmations, verifier.Confirms)
							err = errors.New(errStr)
							return tx, err
						}

						orderedAmount, err := idaemon.NewAmount(verifier.Amount)
						if err != nil {
							errStr := fmt.Sprintf(LIBNAME+":error: orderedAmount %s", err.Error())
							err = errors.New(errStr)
							return tx, err
						}

						var missingAmount idaemon.Amount
						missingAmount = v.Value - orderedAmount
						if v.Value == orderedAmount {
							missingAmount = 0
						}
						missingPercent := (missingAmount.ToCoin() / v.Value.ToCoin()) * 100

						tx = &blockexplorer.ITransaction{
							BlockHeight:   u.BlockHeight,
							Hash:          u.Hash,
							LockTime:      u.LockTime,
							Size:          u.Size,
							Version:       u.Version,
							Confirmations: u.Confirmations,
							Time:          u.Time,

							//custom fields for validation result
							OrderedAmount:       orderedAmount,
							MissingAmount:       missingAmount,
							MissingPercent:      missingPercent,
							Verified:            true,
							BlockExplorerAmount: v.Value,
						}
						var vins []blockexplorer.IVIN
						for _, x := range u.Inputs {
							tmpIn := blockexplorer.IVIN{
								Script:   x.Script,
								Sequence: x.Sequence,
								Witness:  x.Witness,
								TxID:     x.TxID,
								VOUT:     x.VOUT,
								Tree:     x.Tree,
								AmountIn: x.PrevOut.Value,
							}
							vins = append(vins, tmpIn)
						}
						tx.Inputs = vins

						var vouts []blockexplorer.IVOUT
						for _, x := range u.Outputs {
							tmpOut := blockexplorer.IVOUT{
								Addresses: x.Addresses,
								N:         x.N,
								Type:      x.Type,
								Script:    x.Script,
								Value:     x.Value,
							}
							vouts = append(vouts, tmpOut)
						}
						tx.Outputs = vouts
						tx.Seen = true
						return tx, err

					} else {
						//fmt.Printf("\nadditional vout amount: %v to: %v", v.Value, w)
					}
				}
			}
		}

	} else {
		errStr := fmt.Sprintf(LIBNAME+":error: vars passed for verification cannot be checked: \ntxid %s address: %s amount %.8f createdAt %v",
			verifier.TxId, verifier.Address, verifier.Amount, verifier.CreatedAt)
		err = errors.New(errStr)
		return tx, err
	}

	return
}

// GetLatestBlock returns decoded transaction from api
func (c *BlockChainInfo) GetLatestBlock() (latestBlock LatestBlock, err error) {
	r, err := c.client.Do("GET", "latestblock", "", false)
	if err != nil {
		return
	}
	var response LatestBlock
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	latestBlock = response
	//err = errors.New("test failure check error")
	return
}
