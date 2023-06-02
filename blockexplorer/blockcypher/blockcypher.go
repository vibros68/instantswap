package blockcypher

import (
	"encoding/json"
	"fmt"

	"gitlab.com/cryptopower/instantswap/blockexplorer"
	"gitlab.com/cryptopower/instantswap/blockexplorer/global/clients/blockexplorerclient"
	"gitlab.com/cryptopower/instantswap/blockexplorer/global/interfaces/idaemon"
)

const (
	API_BASE = "https://api.blockcypher.com/v1"
	LIBNAME  = "blockcypher"
)

func init() {
	blockexplorer.RegisterExplorer("LTC", "", func(conf blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New("ltc", "main", conf), nil
	})
	blockexplorer.RegisterExplorer("ETH", "", func(conf blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New("eth", "main", conf), nil
	})
}

// represent a * client
type chainzCryptoid struct {
	client    *blockexplorerclient.Client
	conf      blockexplorer.Config
	apiKey    string
	apiSecret string
	coinName  string
	network   string
}

// New return a blockcypher instance
func New(coinName, network string, conf blockexplorer.Config) *chainzCryptoid {
	apiBase := fmt.Sprintf("%s/%s/%s/", API_BASE, coinName, network)
	client := blockexplorerclient.NewClient(apiBase, LIBNAME, conf.EnableOutput, nil)
	return &chainzCryptoid{
		client:   client,
		coinName: coinName,
		network:  network,
		conf:     conf,
	}
}

// SetDebug set enable/disable http request/response dump
func (c *chainzCryptoid) SetDebug(enable bool) {
	c.client.Debug = enable
}

// GetTransaction returns decoded transaction from api
func (c *chainzCryptoid) GetTransaction(txid string) (tx *blockexplorer.ITransaction, err error) {
	r, err := c.client.Do("GET", fmt.Sprintf("txs/%s", txid), "", false)
	if err != nil {
		return nil, err
	}
	var ltcTx Tx
	err = parseData(r, &ltcTx)
	if err != nil {
		return nil, err
	}
	return ltcTx.generalTx(c)
}

// GetTransactionsForAddress
func (c *chainzCryptoid) GetTxsForAddress(address string, limit int, viewKey string) (txs *blockexplorer.IRawAddrResponse, err error) {
	r, err := c.client.Do("GET", fmt.Sprintf("addrs/%s", address), "", false)
	if err != nil {
		return nil, err
	}
	var addr Address
	err = parseData(r, &addr)
	if err != nil {
		return nil, err
	}
	return addr.getIRawAddrResponse(c)
}

// PushTx
func (c *chainzCryptoid) PushTx(txhash string) (res string, err error) {
	return "", fmt.Errorf("ltc is not support PushTx yet")
}

// VerifyTransaction verifies transaction based on values passed in (params: txid, address, amount, createdAt )
func (c *chainzCryptoid) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
	tx = new(blockexplorer.ITransaction)
	if verifier.Address == "" {
		err = fmt.Errorf("%s:error: address is blank so tx cannot be verified", LIBNAME)
		return
	}
	if verifier.Amount == 0 {
		err = fmt.Errorf("%s:error: amount is %.8f so tx cannot be verified", LIBNAME, verifier.Amount)
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
					//tx has been seen on block explorer but still only has 0 confirmations
					if txInfo.Confirmations < verifier.Confirms {
						tx.Seen = true
						err = fmt.Errorf("seen, waiting for confirms (%v/%v)", txInfo.Confirmations, verifier.Confirms)
						return tx, err
					}

					orderedAmount, err := idaemon.NewAmount(verifier.Amount)
					if err != nil {
						err = fmt.Errorf("%s:error: orderedAmount %s", LIBNAME, err.Error())
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
							err = fmt.Errorf("seen, waiting for confirms (%v/%v)", u.Confirmations, verifier.Confirms)
							return tx, err
						}

						orderedAmount, err := idaemon.NewAmount(verifier.Amount)
						if err != nil {
							err = fmt.Errorf("%s:error: orderedAmount %s", LIBNAME, err.Error())
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
		err = fmt.Errorf("%s:error: vars passed for verification cannot be checked: \ntxid %s address: %s amount %.8f createdAt %v",
			LIBNAME, verifier.TxId, verifier.Address, verifier.Amount, verifier.CreatedAt)
		return tx, err
	}

	return
}

func parseData(data []byte, destination interface{}) error {
	var err Err
	if json.Unmarshal(data, &err) == nil {
		if err.ErrorMsg != "" {
			return err
		}
	}
	return json.Unmarshal(data, destination)
}
