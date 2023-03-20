package blockcypher

import (
	"fmt"
	"time"

	"code.cryptopower.dev/group/instantswap/blockexplorer"
	"code.cryptopower.dev/group/instantswap/blockexplorer/global/interfaces/idaemon"
)

type Err struct {
	ErrorMsg string `json:"error"`
}

func (err Err) Error() string {
	return err.ErrorMsg
}

type Tx struct {
	BlockHash     string    `json:"block_hash"`
	BlockHeight   int       `json:"block_height"`
	BlockIndex    int       `json:"block_index"`
	Hash          string    `json:"hash"`
	Addresses     []string  `json:"addresses"`
	Total         int64     `json:"total"`
	Fees          int       `json:"fees"`
	Size          int       `json:"size"`
	Vsize         int       `json:"vsize"`
	Preference    string    `json:"preference"`
	RelayedBy     string    `json:"relayed_by"`
	Confirmed     time.Time `json:"confirmed"`
	Received      time.Time `json:"received"`
	Ver           int       `json:"ver"`
	DoubleSpend   bool      `json:"double_spend"`
	VinSz         int       `json:"vin_sz"`
	VoutSz        int       `json:"vout_sz"`
	Confirmations int       `json:"confirmations"`
	Confidence    int       `json:"confidence"`
	Inputs        []TxInput `json:"inputs"`
	Outputs       []TxOuput `json:"outputs"`
}

type TxInput struct {
	PrevHash    string   `json:"prev_hash"`
	OutputIndex int      `json:"output_index"`
	OutputValue int      `json:"output_value"`
	Sequence    int      `json:"sequence"`
	Addresses   []string `json:"addresses"`
	ScriptType  string   `json:"script_type"`
	Age         int      `json:"age"`
	Witness     []string `json:"witness"`
}

type TxOuput struct {
	Value      int      `json:"value"`
	Script     string   `json:"script"`
	SpentBy    string   `json:"spent_by"`
	Addresses  []string `json:"addresses"`
	ScriptType string   `json:"script_type"`
}

func (t *Tx) generalTx(c *chainzCryptoid) (tx *blockexplorer.ITransaction, err error) {
	if t.Hash == "" {
		return nil, fmt.Errorf("tx not found")
	}
	tx = &blockexplorer.ITransaction{
		BlockHeight:         t.BlockHeight,
		DoubleSpend:         false,
		Hash:                c.ethId(t.Hash),
		Inputs:              t.inputs(c),
		LockTime:            0,
		Outputs:             t.outputs(c),
		Rbf:                 false,
		Size:                t.Size,
		Time:                int(t.Received.Unix()),
		TxIndex:             t.BlockIndex,
		Version:             t.Ver,
		VinSz:               t.VinSz,
		VoutSz:              t.VoutSz,
		Weight:              0,
		Confirmations:       t.Confirmations,
		Seen:                true,
		Verified:            true,
		OrderedAmount:       0,
		BlockExplorerAmount: 0,
		MissingAmount:       0,
		MissingPercent:      0,
	}
	return tx, err
}

func (t *Tx) inputs(c *chainzCryptoid) []blockexplorer.IVIN {
	var inputs = make([]blockexplorer.IVIN, len(t.Inputs))
	for i, input := range t.Inputs {
		inputs[i] = blockexplorer.IVIN{
			Script:      input.ScriptType,
			Sequence:    input.Sequence,
			Witness:     "",
			TxID:        c.ethId(input.PrevHash),
			VOUT:        input.OutputIndex,
			Tree:        0,
			AmountIn:    idaemon.Amount(input.OutputValue),
			BlockIndex:  0,
			BlockHeight: 0,
		}
	}
	return inputs
}

func (t *Tx) outputs(c *chainzCryptoid) []blockexplorer.IVOUT {
	var outputs = make([]blockexplorer.IVOUT, len(t.Outputs))
	for i, output := range t.Outputs {
		outputs[i] = blockexplorer.IVOUT{
			Addresses:   c.ethArrayId(output.Addresses),
			AddrTag:     "",
			AddrTagLink: "",
			N:           0,
			Script:      output.Script,
			Spent:       false,
			TxIndex:     0,
			Type:        output.ScriptType,
			Value:       idaemon.Amount(output.Value),
		}
	}
	return outputs
}

type Address struct {
	Address            string      `json:"address"`
	TotalReceived      int         `json:"total_received"`
	TotalSent          int         `json:"total_sent"`
	Balance            int         `json:"balance"`
	UnconfirmedBalance int         `json:"unconfirmed_balance"`
	FinalBalance       int         `json:"final_balance"`
	NTx                int         `json:"n_tx"`
	UnconfirmedNTx     int         `json:"unconfirmed_n_tx"`
	FinalNTx           int         `json:"final_n_tx"`
	Txrefs             []CompactTx `json:"txrefs"`
	HasMore            bool        `json:"hasMore"`
	TxUrl              string      `json:"tx_url"`
}

type CompactTx struct {
	TxHash        string    `json:"tx_hash"`
	BlockHeight   int       `json:"block_height"`
	TxInputN      int       `json:"tx_input_n"`
	TxOutputN     int       `json:"tx_output_n"`
	Value         int64     `json:"value"`
	RefBalance    int64     `json:"ref_balance"`
	Confirmations int       `json:"confirmations"`
	Confirmed     time.Time `json:"confirmed"`
	DoubleSpend   bool      `json:"double_spend"`
	Spent         bool      `json:"spent,omitempty"`
	SpentBy       string    `json:"spent_by,omitempty"`
}

func (a *Address) getIRawAddrResponse(c *chainzCryptoid) (*blockexplorer.IRawAddrResponse, error) {
	var iTx = blockexplorer.IRawAddrResponse{
		Address:       a.Address,
		FinalBalance:  a.FinalBalance,
		Hash160:       "",
		NTx:           a.NTx,
		TotalReceived: a.TotalReceived,
		TotalSent:     a.TotalSent,
		Txs:           make([]blockexplorer.IRawAddrTx, len(a.Txrefs)),
	}
	for i, tx := range a.Txrefs {
		iTx.Txs[i] = blockexplorer.IRawAddrTx{
			BlockHeight:   tx.BlockHeight,
			Hash:          c.ethId(tx.TxHash),
			Inputs:        nil,
			LockTime:      0,
			Outputs:       nil,
			RelayedBy:     "",
			Result:        0,
			Size:          0,
			Time:          int(tx.Confirmed.Unix()),
			TxIndex:       0,
			Version:       0,
			VinSz:         0,
			VoutSz:        0,
			Weight:        0,
			Confirmations: tx.Confirmations,
		}
	}
	return &iTx, nil
}

func (c *chainzCryptoid) ethId(id string) string {
	if c.coinName == "eth" {
		return fmt.Sprintf("0x%s", id)
	}
	return id
}

func (c *chainzCryptoid) ethArrayId(ids []string) []string {
	var outIds = make([]string, len(ids))
	for i, id := range ids {
		outIds[i] = c.ethId(id)
	}
	return outIds
}
