package blockchair

import (
	"encoding/json"
	"time"

	"code.cryptopower.dev/group/instantswap/blockexplorer"
	"code.cryptopower.dev/group/instantswap/blockexplorer/global/interfaces/idaemon"
)

const timeFormat = "2006-01-02 15:04:05"

type jsonResponse struct {
	Context Context         `json:"context"`
	Data    json.RawMessage `json:"data"`
}

type Context struct {
	Code           int     `json:"code"`
	Source         string  `json:"source"`
	Results        int     `json:"results"`
	State          int     `json:"state"`
	MarketPriceUsd float64 `json:"market_price_usd"`
	Cache          BCCache `json:"cache"`
	Api            BCApi   `json:"api"`
	Servers        string  `json:"servers"`
	Time           float64 `json:"time"`
	RenderTime     float64 `json:"render_time"`
	FullTime       float64 `json:"full_time"`
	RequestCost    int     `json:"request_cost"`
}

type BCApi struct {
	Version         string      `json:"version"`
	LastMajorUpdate string      `json:"last_major_update"`
	NextMajorUpdate interface{} `json:"next_major_update"`
	Documentation   string      `json:"documentation"`
	Notice          string      `json:"notice"`
}

type BCCache struct {
	Live     bool        `json:"live"`
	Duration int         `json:"duration"`
	Since    string      `json:"since"`
	Until    string      `json:"until"`
	Time     interface{} `json:"time"`
}

type TxWrapper struct {
	Transaction Transaction `json:"transaction"`
	Inputs      []TxInput   `json:"inputs"`
	Outputs     []TxOutput  `json:"outputs"`
}

type Transaction struct {
	BlockId            int           `json:"block_id"`
	Id                 int           `json:"id"`
	Hash               string        `json:"hash"`
	Date               string        `json:"date"`
	Time               string        `json:"time"`
	Size               int           `json:"size"`
	IsOverwintered     bool          `json:"is_overwintered"`
	Version            int           `json:"version"`
	VersionGroupId     string        `json:"version_group_id"`
	LockTime           int           `json:"lock_time"`
	ExpiryHeight       int           `json:"expiry_height"`
	IsCoinbase         bool          `json:"is_coinbase"`
	InputCount         int           `json:"input_count"`
	OutputCount        int           `json:"output_count"`
	InputTotal         int64         `json:"input_total"`
	InputTotalUsd      float64       `json:"input_total_usd"`
	OutputTotal        int64         `json:"output_total"`
	OutputTotalUsd     float64       `json:"output_total_usd"`
	Fee                int           `json:"fee"`
	FeeUsd             float64       `json:"fee_usd"`
	FeePerKb           float64       `json:"fee_per_kb"`
	FeePerKbUsd        float64       `json:"fee_per_kb_usd"`
	CddTotal           float64       `json:"cdd_total"`
	ShieldedValueDelta int           `json:"shielded_value_delta"`
	JoinSplitRaw       []interface{} `json:"join_split_raw"`
	ShieldedInputRaw   []interface{} `json:"shielded_input_raw"`
	ShieldedOutputRaw  []interface{} `json:"shielded_output_raw"`
	BindingSignature   interface{}   `json:"binding_signature"`
}

type TxInput struct {
	BlockId                 int         `json:"block_id"`
	TransactionId           int         `json:"transaction_id"`
	Index                   int         `json:"index"`
	TransactionHash         string      `json:"transaction_hash"`
	Date                    string      `json:"date"`
	Time                    string      `json:"time"`
	Value                   int         `json:"value"`
	ValueUsd                float64     `json:"value_usd"`
	Recipient               string      `json:"recipient"`
	Type                    string      `json:"type"`
	ScriptHex               string      `json:"script_hex"`
	IsFromCoinbase          bool        `json:"is_from_coinbase"`
	IsSpendable             interface{} `json:"is_spendable"`
	IsSpent                 bool        `json:"is_spent"`
	SpendingBlockId         int         `json:"spending_block_id"`
	SpendingTransactionId   int         `json:"spending_transaction_id"`
	SpendingIndex           int         `json:"spending_index"`
	SpendingTransactionHash string      `json:"spending_transaction_hash"`
	SpendingDate            string      `json:"spending_date"`
	SpendingTime            string      `json:"spending_time"`
	SpendingValueUsd        float64     `json:"spending_value_usd"`
	SpendingSequence        int         `json:"spending_sequence"`
	SpendingSignatureHex    string      `json:"spending_signature_hex"`
	Lifespan                int         `json:"lifespan"`
	Cdd                     float64     `json:"cdd"`
}

type TxOutput struct {
	BlockId                 int         `json:"block_id"`
	TransactionId           int         `json:"transaction_id"`
	Index                   int         `json:"index"`
	TransactionHash         string      `json:"transaction_hash"`
	Date                    string      `json:"date"`
	Time                    string      `json:"time"`
	Value                   int64       `json:"value"`
	ValueUsd                float64     `json:"value_usd"`
	Recipient               string      `json:"recipient"`
	Type                    string      `json:"type"`
	ScriptHex               string      `json:"script_hex"`
	IsFromCoinbase          bool        `json:"is_from_coinbase"`
	IsSpendable             interface{} `json:"is_spendable"`
	IsSpent                 bool        `json:"is_spent"`
	SpendingBlockId         int         `json:"spending_block_id"`
	SpendingTransactionId   int         `json:"spending_transaction_id"`
	SpendingIndex           int         `json:"spending_index"`
	SpendingTransactionHash string      `json:"spending_transaction_hash"`
	SpendingDate            string      `json:"spending_date"`
	SpendingTime            string      `json:"spending_time"`
	SpendingValueUsd        float64     `json:"spending_value_usd"`
	SpendingSequence        int64       `json:"spending_sequence"`
	SpendingSignatureHex    string      `json:"spending_signature_hex"`
	Lifespan                int         `json:"lifespan"`
	Cdd                     float64     `json:"cdd"`
}

func (b *BlockChair) generalTx(txW *TxWrapper, ctx *Context) (tx *blockexplorer.ITransaction, err error) {
	var t, _ = time.Parse(timeFormat, txW.Transaction.Time)
	tx = &blockexplorer.ITransaction{
		BlockHeight:   txW.Transaction.BlockId,
		DoubleSpend:   false,
		Hash:          txW.Transaction.Hash,
		Inputs:        nil,
		LockTime:      txW.Transaction.LockTime,
		Outputs:       nil,
		Rbf:           false,
		Size:          0,
		Time:          int(t.Unix()),
		TxIndex:       0,
		Version:       txW.Transaction.Version,
		VinSz:         0,
		VoutSz:        0,
		Weight:        0,
		Confirmations: ctx.State - txW.Transaction.BlockId + 1,
	}
	for _, txIn := range txW.Inputs {
		tx.Inputs = append(tx.Inputs, blockexplorer.IVIN{
			Script:      txIn.ScriptHex,
			Sequence:    txIn.SpendingSequence,
			Witness:     "",
			TxID:        txIn.TransactionHash,
			VOUT:        0,
			Tree:        0,
			AmountIn:    idaemon.Amount(txIn.Value),
			BlockIndex:  txIn.Index,
			BlockHeight: txIn.BlockId,
		})
	}
	for _, txOut := range txW.Outputs {
		tx.Outputs = append(tx.Outputs, blockexplorer.IVOUT{
			Addresses:   []string{txOut.Recipient},
			AddrTag:     "",
			AddrTagLink: "",
			N:           0,
			Script:      txOut.ScriptHex,
			Spent:       false,
			TxIndex:     txOut.Index,
			Type:        txOut.Type,
			Value:       idaemon.Amount(txOut.Value),
		})
	}
	return
}

type SimpleTx struct {
	BlockId       int    `json:"block_id"`
	Hash          string `json:"hash"`
	Time          string `json:"time"`
	BalanceChange int    `json:"balance_change"`
}
type Utxo struct {
	BlockId         int    `json:"block_id"`
	TransactionHash string `json:"transaction_hash"`
	Index           int    `json:"index"`
	Value           int    `json:"value"`
}
type Address struct {
	Type               string      `json:"type"`
	ScriptHex          string      `json:"script_hex"`
	Balance            int         `json:"balance"`
	BalanceUsd         float64     `json:"balance_usd"`
	Received           int         `json:"received"`
	ReceivedUsd        float64     `json:"received_usd"`
	Spent              int         `json:"spent"`
	SpentUsd           float64     `json:"spent_usd"`
	OutputCount        int         `json:"output_count"`
	UnspentOutputCount int         `json:"unspent_output_count"`
	FirstSeenReceiving string      `json:"first_seen_receiving"`
	LastSeenReceiving  string      `json:"last_seen_receiving"`
	FirstSeenSpending  interface{} `json:"first_seen_spending"`
	LastSeenSpending   interface{} `json:"last_seen_spending"`
	ScripthashType     interface{} `json:"scripthash_type"`
	TransactionCount   int         `json:"transaction_count"`
}

func (b *BlockChair) generalAddr(address string, addr *AddrWrapper, ctx *Context) (txs *blockexplorer.IRawAddrResponse) {
	txs = &blockexplorer.IRawAddrResponse{
		Address:       address,
		FinalBalance:  addr.Address.Balance,
		Hash160:       "",
		NTx:           addr.Address.TransactionCount,
		TotalReceived: addr.Address.Received,
		TotalSent:     addr.Address.Spent,
		Txs:           nil,
	}
	for _, tx := range addr.Transactions {
		txs.Txs = append(txs.Txs, blockexplorer.IRawAddrTx{
			BlockHeight:   tx.BlockId,
			Hash:          tx.Hash,
			LockTime:      0,
			RelayedBy:     "",
			Result:        tx.BalanceChange,
			Size:          0,
			Time:          0,
			TxIndex:       0,
			Version:       0,
			VinSz:         0,
			VoutSz:        0,
			Weight:        0,
			Confirmations: ctx.State - tx.BlockId + 1,
		})
	}
	return txs
}

type AddrWrapper struct {
	Address      Address    `json:"address"`
	Transactions []SimpleTx `json:"transactions"`
	Utxo         []Utxo     `json:"utxo"`
}
