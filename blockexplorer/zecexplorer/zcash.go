package zecexplorer

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/vibros68/instantswap/blockexplorer"
	"github.com/vibros68/instantswap/blockexplorer/global/clients/blockexplorerclient"
	"github.com/vibros68/instantswap/blockexplorer/global/interfaces/idaemon"
)

const (
	API_BASE                   = "https://api.zcha.in/v2/" //  API endpoint
	DEFAULT_HTTPCLIENT_TIMEOUT = 30                        // HTTP client timeout
	LIBNAME                    = "zcha"
)

func init() {
	blockexplorer.RegisterExplorer("ZEC", "", func(config blockexplorer.Config) (blockexplorer.IBlockExplorer, error) {
		return New(config), nil
	})
}

// New return a instanciate cryptopia struct
func New(conf blockexplorer.Config) *ZcashExplorer {
	client := blockexplorerclient.NewClient(API_BASE, LIBNAME, conf.EnableOutput, nil)
	return &ZcashExplorer{client: client}
}

type ZcashExplorer struct {
	client *blockexplorerclient.Client
}

func (z *ZcashExplorer) VerifyByAddress(req blockexplorer.AddressVerifyRequest) (vr *blockexplorer.VerifyResult, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (z *ZcashExplorer) getNetwork() (*Network, error) {
	r, err := z.client.Do("GET", "mainnet/network", "", false)
	if err != nil {
		return nil, err
	}
	var network Network
	err = json.Unmarshal(r, &network)
	return &network, err
}

func (z *ZcashExplorer) GetTransaction(txId string) (*blockexplorer.ITransaction, error) {
	r, err := z.client.Do("GET", fmt.Sprintf("mainnet/transactions/%s", txId), "", false)
	if err != nil {
		return nil, err
	}
	var tx Transaction
	if err = json.Unmarshal(r, &tx); err != nil {
		return nil, err
	}
	network, _ := z.getNetwork()
	return tx.generalTx(network), nil
}
func (z *ZcashExplorer) GetTxsForAddress(address string, limit int, viewKey string) (account *blockexplorer.IRawAddrResponse, err error) {
	if limit > 20 || limit < 1 {
		limit = 20
	}
	var zcashAccount Account
	r, err := z.client.Do("GET", fmt.Sprintf("mainnet/accounts/%s", address), "", false)
	if err = json.Unmarshal(r, &zcashAccount); err != nil {
		return nil, err
	}
	account = zcashAccount.acount()
	var recvTxs []Transaction
	r, err = z.client.Do("GET",
		fmt.Sprintf("mainnet/accounts/%s/recv?limit=%d&offset=0&sort=timestamp&direction=descending", address, limit), "", false)
	if err = json.Unmarshal(r, &recvTxs); err != nil {
		return nil, err
	}
	var sendTxs []Transaction
	r, err = z.client.Do("GET",
		fmt.Sprintf("mainnet/accounts/%s/sent?limit=%d&offset=0&sort=timestamp&direction=descending", address, limit), "", false)
	if err = json.Unmarshal(r, &sendTxs); err != nil {
		return nil, err
	}
	var txs = append(recvTxs, sendTxs...)
	sort.SliceStable(txs, func(i, j int) bool {
		return txs[i].Timestamp < txs[j].Timestamp
	})
	account.Txs = convertTxs(txs)
	return account, nil
}

// VerifyTransaction verifies transaction based on values passed in (params: txid, address (required), amount (required), createdAt(unix timestamp) )
func (z *ZcashExplorer) VerifyTransaction(verifier blockexplorer.TxVerifyRequest) (tx *blockexplorer.ITransaction, err error) {
	if verifier.Address == "" {
		return nil, fmt.Errorf(LIBNAME + ":error: address is blank so tx cannot be verified")
	}
	if verifier.Amount == 0 {
		return nil, fmt.Errorf(LIBNAME+":error: amount is %.8f so tx cannot be verified", verifier.Amount)
	}
	tx = new(blockexplorer.ITransaction)
	if verifier.TxId != "" && verifier.Address != "" { //verify tx if txid is available
		txInfo, err := z.GetTransaction(verifier.TxId)
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
						return tx, fmt.Errorf("seen, waiting for confirms (%v/%v)", txInfo.Confirmations, verifier.Confirms)
					}

					orderedAmount, err := idaemon.NewAmount(verifier.Amount)
					if err != nil {
						return tx, fmt.Errorf(LIBNAME+":error: orderedAmount %s", err.Error())
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
		txInfo, err := z.GetTxsForAddress(verifier.Address, 10, "")
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
							return tx, fmt.Errorf("seen, waiting for confirms (%v/%v)", u.Confirmations, verifier.Confirms)
						}

						orderedAmount, err := idaemon.NewAmount(verifier.Amount)
						if err != nil {
							return tx, fmt.Errorf(LIBNAME+":error: orderedAmount %s", err.Error())
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
		return tx, fmt.Errorf(LIBNAME+":error: vars passed for verification cannot be checked: \ntxid %s address: %s amount %.8f createdAt %v",
			verifier.TxId, verifier.Address, verifier.Amount, verifier.CreatedAt)
	}
	return nil, nil
}

// PushTx pushes a raw tx hash
func (z *ZcashExplorer) PushTx(rawTxHash string) (result string, err error) {
	return "", fmt.Errorf("%s:error: PushTx is not supported yet... ", LIBNAME)
}
