package xmrexplorer

import (
	"github.com/vibros68/instantswap/blockexplorer"
	"github.com/vibros68/instantswap/blockexplorer/global/interfaces/idaemon"
)

const XMR_EXTRA_UNIT = 10000

type Response struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}

type XmrChainError struct {
	Title string `json:"title"`
}

type Transaction struct {
	BlockHeight   int      `json:"block_height"`
	Coinbase      bool     `json:"coinbase"`
	Confirmations int      `json:"confirmations"`
	CurrentHeight int      `json:"current_height"`
	Extra         string   `json:"extra"`
	Inputs        []Input  `json:"inputs"`
	Mixin         int      `json:"mixin"`
	Outputs       []Output `json:"outputs"`
	PaymentId     string   `json:"payment_id"`
	PaymentId8    string   `json:"payment_id8"`
	RctType       int      `json:"rct_type"`
	Timestamp     int      `json:"timestamp"`
	TimestampUtc  string   `json:"timestamp_utc"`
	TxFee         int      `json:"tx_fee"`
	TxHash        string   `json:"tx_hash"`
	TxSize        int      `json:"tx_size"`
	TxVersion     int      `json:"tx_version"`
	XmrInputs     int      `json:"xmr_inputs"`
	XmrOutputs    int      `json:"xmr_outputs"`
}

type Input struct {
	Amount   int    `json:"amount"`
	KeyImage string `json:"key_image"`
	Mixins   []struct {
		BlockNo   int    `json:"block_no"`
		PublicKey string `json:"public_key"`
	} `json:"mixins"`
}

type Output struct {
	Amount    int    `json:"amount"`
	PublicKey string `json:"public_key"`
}

func (i *Input) IVIN() blockexplorer.IVIN {
	return blockexplorer.IVIN{
		Script:      "",
		Sequence:    0,
		Witness:     "",
		TxID:        "",
		VOUT:        0,
		Tree:        0,
		AmountIn:    idaemon.Amount(i.Amount / XMR_EXTRA_UNIT),
		BlockIndex:  0,
		BlockHeight: 0,
	}
}

func (o *Output) IVOUT() blockexplorer.IVOUT {
	return blockexplorer.IVOUT{
		Addresses:   nil,
		AddrTag:     "",
		AddrTagLink: "",
		N:           0,
		Script:      "",
		Spent:       false,
		TxIndex:     0,
		Type:        "",
		Value:       idaemon.Amount(o.Amount / XMR_EXTRA_UNIT),
	}
}

func (t *Transaction) inputs() []blockexplorer.IVIN {
	var ivins = make([]blockexplorer.IVIN, len(t.Inputs))
	for i, in := range t.Inputs {
		ivins[i] = in.IVIN()
	}
	return ivins
}

func (t *Transaction) outputs() []blockexplorer.IVOUT {
	var iVouts = make([]blockexplorer.IVOUT, len(t.Outputs))
	for i, in := range t.Outputs {
		iVouts[i] = in.IVOUT()
	}
	return iVouts
}

func (t *Transaction) ITransaction() *blockexplorer.ITransaction {
	return &blockexplorer.ITransaction{
		BlockHeight:         t.BlockHeight,
		DoubleSpend:         false,
		Hash:                t.TxHash,
		Inputs:              t.inputs(),
		LockTime:            0,
		Outputs:             t.outputs(),
		Rbf:                 false,
		Size:                0,
		Time:                t.Timestamp,
		TxIndex:             0,
		Version:             t.TxVersion,
		VinSz:               0,
		VoutSz:              0,
		Weight:              0,
		Confirmations:       t.Confirmations,
		Seen:                false,
		Verified:            t.Confirmations != 0,
		OrderedAmount:       0,
		BlockExplorerAmount: 0,
		MissingAmount:       0,
		MissingPercent:      0,
	}
}

type VerifyOutput struct {
	Amount       int64  `json:"amount"`
	Match        bool   `json:"match"`
	OutputIdx    int    `json:"output_idx"`
	OutputPubkey string `json:"output_pubkey"`
}

type TxVerifier struct {
	Address         string         `json:"address"`
	Outputs         []VerifyOutput `json:"outputs"`
	TxConfirmations int            `json:"tx_confirmations"`
	TxHash          string         `json:"tx_hash"`
	TxProve         bool           `json:"tx_prove"`
	TxTimestamp     int            `json:"tx_timestamp"`
	Viewkey         string         `json:"viewkey"`
}

func (v *TxVerifier) ITransaction(verifier blockexplorer.TxVerifyRequest) *blockexplorer.ITransaction {
	var amount int64
	var seen bool
	for _, output := range v.Outputs {
		if output.Match {
			seen = true
			amount += output.Amount
		}
	}
	orderedAmount, _ := idaemon.NewAmount(verifier.Amount)
	explorerAmount := idaemon.Amount(amount / XMR_EXTRA_UNIT)
	return &blockexplorer.ITransaction{
		BlockHeight:         0,
		DoubleSpend:         false,
		Hash:                v.TxHash,
		Inputs:              nil,
		LockTime:            0,
		Outputs:             nil,
		Rbf:                 false,
		Size:                0,
		Time:                v.TxTimestamp,
		TxIndex:             0,
		Version:             0,
		VinSz:               0,
		VoutSz:              0,
		Weight:              0,
		Confirmations:       v.TxConfirmations,
		Seen:                seen,
		Verified:            seen,
		OrderedAmount:       orderedAmount,
		BlockExplorerAmount: idaemon.Amount(amount / XMR_EXTRA_UNIT),
		MissingAmount:       orderedAmount - explorerAmount,
		MissingPercent:      100 * float64(orderedAmount-explorerAmount) / float64(orderedAmount),
	}
}

type OutputsBlocks struct {
	Address string        `json:"address"`
	Height  int           `json:"height"`
	Limit   string        `json:"limit"`
	Mempool bool          `json:"mempool"`
	Outputs []OutputBlock `json:"outputs"`
	Viewkey string        `json:"viewkey"`
}

type OutputBlock struct {
	Amount       int64  `json:"amount"`
	BlockNo      int    `json:"block_no"`
	InMempool    bool   `json:"in_mempool"`
	OutputIdx    int    `json:"output_idx"`
	OutputPubkey string `json:"output_pubkey"`
	PaymentId    string `json:"payment_id"`
	TxHash       string `json:"tx_hash"`
}

func (o *OutputsBlocks) IRawAddrResponse() *blockexplorer.IRawAddrResponse {
	return &blockexplorer.IRawAddrResponse{
		Address:       o.Address,
		FinalBalance:  0,
		Hash160:       "",
		NTx:           0,
		TotalReceived: 0,
		TotalSent:     0,
		Txs:           convertIRawAddrTx(o.Outputs, o.Address),
	}
}

func convertIRawAddrTx(outputs []OutputBlock, address string) []blockexplorer.IRawAddrTx {
	var addrTxs = make([]blockexplorer.IRawAddrTx, len(outputs))
	for i, output := range outputs {
		addrTxs[i] = blockexplorer.IRawAddrTx{
			BlockHeight: output.BlockNo,
			Hash:        output.TxHash,
			Inputs:      nil,
			LockTime:    0,
			Outputs: []blockexplorer.IRawAddrOutput{
				{
					Addresses: []string{address},
					Value:     idaemon.Amount(output.Amount / XMR_EXTRA_UNIT),
				},
			},
			RelayedBy:     "",
			Result:        0,
			Size:          0,
			Time:          0,
			TxIndex:       0,
			Version:       0,
			VinSz:         0,
			VoutSz:        0,
			Weight:        0,
			Confirmations: 0,
		}
	}
	return addrTxs
}
