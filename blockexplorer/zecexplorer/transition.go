package zecexplorer

import (
	"code.cryptopower.dev/group/instantswap/blockexplorer"
	"code.cryptopower.dev/group/instantswap/blockexplorer/global/interfaces/idaemon"
)

func (t *Transaction) valueZat() int {
	var amount int
	for _, vout := range t.Vout {
		amount += vout.ValueZat
	}
	return amount
}

func (t *Transaction) generalTx(network *Network) *blockexplorer.ITransaction {
	var iTx = &blockexplorer.ITransaction{
		BlockHeight:         t.BlockHeight,
		DoubleSpend:         false,
		Hash:                t.Hash,
		Inputs:              t.inputs(),
		LockTime:            t.LockTime,
		Outputs:             t.outputs(),
		Rbf:                 false,
		Size:                0,
		Time:                t.Time,
		TxIndex:             t.Index,
		Version:             t.Version,
		VinSz:               0,
		VoutSz:              0,
		Weight:              0,
		Confirmations:       0,
		Seen:                false,
		Verified:            false,
		OrderedAmount:       0,
		BlockExplorerAmount: idaemon.Amount(t.valueZat()),
		MissingAmount:       0,
		MissingPercent:      0,
	}
	if network != nil {
		iTx.Confirmations = network.BlockNumber - t.BlockHeight + 1
		if iTx.Confirmations != 0 {
			iTx.Verified = true
		}
	}
	return iTx
}

func (t *Transaction) inputs() []blockexplorer.IVIN {
	var iVin = make([]blockexplorer.IVIN, len(t.Vin))
	for i, vin := range t.Vin {
		iVin[i] = blockexplorer.IVIN{
			//Script:      vin.ScriptSig,
			Sequence:    vin.Sequence,
			Witness:     "",
			TxID:        vin.Txid,
			VOUT:        vin.Vout,
			Tree:        vin.RetrievedVout.N,
			AmountIn:    idaemon.Amount(vin.RetrievedVout.ValueZat),
			BlockIndex:  t.Index,
			BlockHeight: t.BlockHeight,
		}
	}
	return iVin
}

func (t *Transaction) outputs() []blockexplorer.IVOUT {
	var iVout = make([]blockexplorer.IVOUT, len(t.Vout))
	for i, vout := range t.Vout {
		iVout[i] = blockexplorer.IVOUT{
			Addresses:   vout.ScriptPubKey.Addresses,
			AddrTag:     "",
			AddrTagLink: "",
			N:           vout.N,
			Script:      "",
			Spent:       false,
			TxIndex:     t.Index,
			Type:        vout.ScriptPubKey.Type,
			Value:       idaemon.Amount(vout.ValueZat),
		}
	}
	return iVout
}

func (t *Transaction) iRawInputs() []blockexplorer.IRawAddrInput {
	var iVin = make([]blockexplorer.IRawAddrInput, len(t.Vin))
	for i, vin := range t.Vin {
		iVin[i] = blockexplorer.IRawAddrInput{
			PrevOut:  vin.RetrievedVout.IRawOutput(),
			Script:   "",
			Sequence: vin.Sequence,
			Witness:  "",
			TxID:     vin.Txid,
			VOUT:     vin.Vout,
			Tree:     vin.RetrievedVout.N,
		}
	}
	return iVin
}

func (t *Transaction) iRawOutputs() []blockexplorer.IRawAddrOutput {
	var iVout = make([]blockexplorer.IRawAddrOutput, len(t.Vout))
	for i, vout := range t.Vout {
		iVout[i] = vout.IRawOutput()
	}
	return iVout
}
