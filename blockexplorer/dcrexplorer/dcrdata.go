package dcrexplorer

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vibros68/instantswap/blockexplorer"
	"github.com/vibros68/instantswap/blockexplorer/global/clients/blockexplorerclient"
	"github.com/vibros68/instantswap/blockexplorer/global/interfaces/idaemon"
)

const (
	API_BASE                   = "https://explorer.dcrdata.org/api/" //  API endpoint
	DEFAULT_HTTPCLIENT_TIMEOUT = 30                                  // HTTP client timeout
	LIBNAME                    = "dcrdata"
)

func init() {
	blockexplorer.RegisterExplorer("DCR", "", func(conf blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New(conf), nil
	})
}

// New return a instanciate cryptopia struct
func New(conf blockexplorer.Config) *DCRData {
	client := blockexplorerclient.NewClient(API_BASE, LIBNAME, conf.EnableOutput, nil)
	return &DCRData{client: client}
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
type DCRData struct {
	client    *blockexplorerclient.Client
	apiKey    string
	apiSecret string
}

// set enable/disable http request/response dump
func (c *DCRData) SetDebug(enable bool) {
	c.client.Debug = enable
}

func (c *DCRData) VerifyByAddress(req blockexplorer.AddressVerifyRequest) (vr *blockexplorer.VerifyResult, err error) {
	txs, err := c.getTxsForAddress(req.Address, 25)
	if err != nil {
		return nil, err
	}
	for _, tx := range txs {
		for _, out := range tx.Vout {
			if len(out.ScriptPubKey.Addresses) == 1 && out.ScriptPubKey.Addresses[0] == req.Address {
				if out.Value == req.Amount {
					return &blockexplorer.VerifyResult{
						Seen:                true,
						Verified:            true,
						OrderedAmount:       req.Amount,
						BlockExplorerAmount: out.Value,
						MissingAmount:       0,
						MissingPercent:      0,
					}, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("not found")
}

// GetTransaction returns decoded transaction from explorer.dcrdata.org/api
func (c *DCRData) GetTransaction(txid string) (tx *blockexplorer.ITransaction, err error) {
	r, err := c.client.Do("GET", "tx/"+txid, "", false)
	if err != nil {
		return
	}
	var tmp Transaction
	if err = json.Unmarshal(r, &tmp); err != nil {
		return
	}
	tx = &blockexplorer.ITransaction{
		Confirmations: tmp.Confirmations,
		BlockHeight:   tmp.Block.Blockheight,
		//DoubleSpend: tmp.DoubleSpend,
		Hash:     tmp.Txid,
		LockTime: tmp.Locktime,
		//Rbf:tmp.Rbf,
		Size:    tmp.Size,
		Time:    tmp.Block.Time,
		TxIndex: tmp.Block.Blockindex,
		Version: tmp.Version,
		//VinSz:   tmp.VinSz,
		//VoutSz:  tmp.VoutSz,
		//Weight:  tmp.Weight,
	}
	for _, v := range tmp.Vin {
		amountIn, err := idaemon.NewAmount(v.Amountin)
		if err != nil {
			return tx, err
		}
		tmpIn := blockexplorer.IVIN{
			Script:      v.ScriptSig.Hex,
			Sequence:    v.Sequence,
			VOUT:        v.Vout,
			Tree:        v.Tree,
			AmountIn:    amountIn,
			BlockHeight: v.Blockheight,
			BlockIndex:  v.Blockindex,
			TxID:        v.Txid,
			//Witness:  v.,
		}
		tx.Inputs = append(tx.Inputs, tmpIn)
	}
	for _, v := range tmp.Vout {
		valueAmount, err := idaemon.NewAmount(v.Value)
		if err != nil {
			return tx, err
		}
		tmpOut := blockexplorer.IVOUT{
			//Script:      v.ScriptPubKey.,
			Addresses: v.ScriptPubKey.Addresses,
			//AddrTag:     v.AddrTag,
			//AddrTagLink: v.AddrTagLink,
			N: v.N,
			//Spent:       v.Spent,
			//TxIndex:     v.TxIndex,
			Type:  v.ScriptPubKey.Type,
			Value: valueAmount,
		}
		tx.Outputs = append(tx.Outputs, tmpOut)
	}
	//err = errors.New("test failure check error")
	return
}

func (c *DCRData) getTxsForAddress(address string, limit int) (txs []RawAddrTx, err error) {
	r, err := c.client.Do("GET", fmt.Sprintf("address/%s/count/%v/raw", address, limit), "", false)
	if err != nil {
		return nil, fmt.Errorf(" could not find/parse address %s msg: %s", address, err.Error())
	}
	err = json.Unmarshal(r, &txs)
	return
}

// GetTransactionsForAddress
func (c *DCRData) GetTxsForAddress(address string, limit int, viewKey string) (txs *blockexplorer.IRawAddrResponse, err error) {
	tmp, err := c.getTxsForAddress(address, limit)
	if err != nil {
		return nil, err
	}

	txs = &blockexplorer.IRawAddrResponse{
		Address: address,
		/* FinalBalance: tmp.FinalBalance,
		Hash160: tmp.Hash160,
		NTx: tmp.NTx,
		TotalReceived: tmp.TotalReceived,
		TotalSent: tmp.TotalSent, */
	}
	var allTxs []blockexplorer.IRawAddrTx
	//gather txs for this address and format them for interface
	for _, v := range tmp {
		tmpTx := blockexplorer.IRawAddrTx{
			Hash:          v.Txid,
			LockTime:      v.Locktime,
			Size:          v.Size,
			Time:          v.Time,
			Version:       v.Version,
			Confirmations: v.Confirmations,
		}
		var tmpInputs []blockexplorer.IRawAddrInput
		for _, w := range v.Vin {

			tmpIn := blockexplorer.IRawAddrInput{
				//Script: w.ScriptSig.,
				Sequence: w.Sequence,
				TxID:     w.Txid,
				VOUT:     w.Vout,
				Tree:     w.Tree,
				//Witness: w.Witness,
			}

			tmpAmount, err := idaemon.NewAmount(w.Amountin)
			if err != nil {
				return nil, err
			}

			tmpPrevOutput := blockexplorer.IRawAddrOutput{
				Addresses: w.PrevOut.Addresses,
				//Script: w.ScriptSig.,
				TxIndex: w.Blockindex,
				Value:   tmpAmount,
			}
			tmpIn.PrevOut = tmpPrevOutput

			tmpInputs = append(tmpInputs, tmpIn)

		}
		tmpTx.Inputs = tmpInputs //assign inputs to current tx item

		var tmpOuputs []blockexplorer.IRawAddrOutput
		for _, w := range v.Vout {

			tmpAmount, err := idaemon.NewAmount(w.Value)
			if err != nil {
				return nil, err
			}

			tmpOut := blockexplorer.IRawAddrOutput{
				Addresses: w.ScriptPubKey.Addresses,
				N:         w.N,
				Type:      w.ScriptPubKey.Type,
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

// PushTx pushed a raw tx hash to mainnet
func (c *DCRData) PushTx(txhash string) (res string, err error) {
	err = errors.New("dcrdata:error: pushtx is not available yet... ")
	/* payload, err := json.Marshal(PushTxRequest{Event: "sendtx", Message: txhash})
	if err != nil {
		return
	}
	//figure out endpoint for dcrdata send/broadcast raw tx
	r, err := c.client.Do(LIBNAME, "https://explorer.dcrdata.org/insight/api/", "POST", "broadcastRawTransaction", string(payload), false)
	if err != nil {
		return
	}

	res = fmt.Sprintf("%s", r) */

	return
}

// VerifyTransaction verifies transaction based on values passed in (params: txid, address, amount, createdAt )
func (c *DCRData) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
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
			errStr := fmt.Sprintf("%s", err.Error())
			err = errors.New(errStr)
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
	} else if verifier.Address != "" && verifier.CreatedAt > 0 { //verify tx on blockchain based on address history for address var
		txInfo, err := c.GetTxsForAddress(verifier.Address, 10, "")
		if err != nil {
			errStr := fmt.Sprintf(LIBNAME+":error: %s", err.Error())
			err = errors.New(errStr)
			return tx, err
		}

		for _, u := range txInfo.Txs {
			for _, v := range u.Outputs {
				for _, w := range v.Addresses {

					if w == verifier.Address && int64(u.Time) >= int64(verifier.CreatedAt) {
						/* fmt.Printf("\ncreatedAt: %v tx time: %v", int64(createdAt), int64(u.Time))
						fmt.Printf("seen") */

						//tx has been seen on block explorer but still only has 0 confirmations
						if u.Confirmations < verifier.Confirms {
							tx.Seen = true
							err = fmt.Errorf("seen, waiting for confirms (%v/%v)", u.Confirmations, verifier.Confirms)
							return tx, err
						}

						orderedAmount, err := idaemon.NewAmount(verifier.Amount)
						if err != nil {
							err = fmt.Errorf(LIBNAME+":error: orderedAmount %s", err.Error())
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
		err := fmt.Errorf(LIBNAME+":error: vars passed for verification cannot be checked: \ntxid %s address: %s amount %.8f createdAt %v",
			verifier.TxId, verifier.Address, verifier.Amount, verifier.CreatedAt)
		return tx, err
	}

	return
}

// GetTransaction returns decoded transaction from explorer.dcrdata.org/api
func (c *DCRData) GetDecodedTransaction(txid string) (tx DecodedTransaction, err error) {
	r, err := c.client.Do("GET", "tx/decoded/"+txid, "", false)
	if err != nil {
		return
	}
	var response DecodedTransaction
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	tx = response
	//err = errors.New("test failure check error")
	return
}
