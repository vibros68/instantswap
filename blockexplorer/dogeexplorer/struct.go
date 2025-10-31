package dogeexplorer

import (
	"github.com/vibros68/instantswap/blockexplorer"
	"github.com/vibros68/instantswap/blockexplorer/global/interfaces/idaemon"
)

type Transaction struct {
	Hash          string      `json:"hash"`
	Confirmations int         `json:"confirmations"`
	Size          int         `json:"size"`
	VSize         int         `json:"vsize"`
	Weight        interface{} `json:"weight"`
	Version       int         `json:"version"`
	LockTime      int         `json:"locktime"`
	BlockHash     string      `json:"block_hash"`
	Time          int         `json:"time"`
	InputsN       int         `json:"inputs_n"`
	InputsValue   float64     `json:"inputs_value,string"`
	Inputs        []Input     `json:"inputs"`
	OutputsN      int         `json:"outputs_n"`
	OutputsValue  float64     `json:"outputs_value,string"`
	Outputs       []Output    `json:"outputs"`
	Fee           float64     `json:"fee,string"`
	Price         float64     `json:"price,string"`
}

type Output struct {
	Pos     int     `json:"pos"`
	Value   float64 `json:"value,string"`
	Type    string  `json:"type"`
	Address string  `json:"address"`
	Script  struct {
		Hex string `json:"hex"`
		Asm string `json:"asm"`
	} `json:"script"`
	Spent struct {
		Hash string `json:"hash"`
		Pos  int    `json:"pos"`
	} `json:"spent"`
}

type Input struct {
	Pos       int     `json:"pos"`
	Value     float64 `json:"value,string"`
	Address   string  `json:"address"`
	ScriptSig struct {
		Hex string `json:"hex"`
	} `json:"scriptSig"`
	PreviousOutput struct {
		Hash string `json:"hash"`
		Pos  int    `json:"pos"`
	} `json:"previous_output"`
}

type Res struct {
	Error   string `json:"error"`
	Success int    `json:"success"`
}

func (tx *Transaction) inputs() []blockexplorer.IVIN {
	var inputs []blockexplorer.IVIN
	for _, input := range tx.Inputs {
		amount, _ := idaemon.NewAmount(input.Value)
		inputs = append(inputs, blockexplorer.IVIN{
			Script:      input.ScriptSig.Hex,
			Sequence:    0,
			Witness:     "",
			TxID:        tx.Hash,
			VOUT:        0,
			Tree:        0,
			AmountIn:    amount,
			BlockIndex:  0,
			BlockHeight: 0,
		})
	}
	return inputs
}
func (tx *Transaction) outputs() []blockexplorer.IVOUT {
	var outputs []blockexplorer.IVOUT
	for _, output := range tx.Outputs {
		amount, _ := idaemon.NewAmount(output.Value)
		outputs = append(outputs, blockexplorer.IVOUT{
			Addresses: []string{
				output.Address,
			},
			AddrTag:     "",
			AddrTagLink: "",
			N:           0,
			Script:      output.Script.Hex,
			Spent:       false,
			TxIndex:     output.Pos,
			Type:        output.Type,
			Value:       amount,
		})
	}
	return outputs
}
func (tx *Transaction) tx() *blockexplorer.ITransaction {
	return &blockexplorer.ITransaction{
		BlockHeight:   0,
		DoubleSpend:   false,
		Hash:          tx.Hash,
		Inputs:        tx.inputs(),
		LockTime:      tx.LockTime,
		Outputs:       tx.outputs(),
		Rbf:           false,
		Size:          0,
		Time:          tx.Time,
		TxIndex:       0,
		Version:       tx.Version,
		VinSz:         0,
		VoutSz:        0,
		Weight:        0,
		Confirmations: tx.Confirmations,
		// verification: ignore
		Seen:                false,
		Verified:            false,
		OrderedAmount:       0,
		BlockExplorerAmount: 0,
		MissingAmount:       0,
		MissingPercent:      0,
	}
}

type TxForAddress struct {
	Hash  string  `json:"hash"`
	Value float64 `json:"value,string"`
	Time  int     `json:"time"`
	Block int     `json:"block"`
	Price float64 `json:"price,string"`
}

func convertTxs(txs []TxForAddress, iRawAddr *blockexplorer.IRawAddrResponse) {
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
	iRawAddr.Txs = txsRaw
	/*return &blockexplorer.IRawAddrResponse{
		Address:       "",
		FinalBalance:  0,
		Hash160:       "",
		NTx:           0,
		TotalReceived: 0,
		TotalSent:     0,
		Txs:           nil,
	}*/
}
