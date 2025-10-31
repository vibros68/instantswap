package zecexplorer

import (
	"github.com/vibros68/instantswap/blockexplorer"
	"github.com/vibros68/instantswap/blockexplorer/global/interfaces/idaemon"
)

type Transaction struct {
	Hash            string        `json:"hash"`
	MainChain       bool          `json:"mainChain"`
	Fee             float64       `json:"fee"`
	Type            string        `json:"type"`
	Shielded        bool          `json:"shielded"`
	Index           int           `json:"index"`
	BlockHash       string        `json:"blockHash"`
	BlockHeight     int           `json:"blockHeight"`
	Version         int           `json:"version"`
	LockTime        int           `json:"lockTime"`
	Timestamp       int           `json:"timestamp"`
	Time            int           `json:"time"`
	Vin             []Vin         `json:"vin"`
	Vout            []Vout        `json:"vout"`
	Vjoinsplit      []interface{} `json:"vjoinsplit"`
	VShieldedOutput int           `json:"vShieldedOutput"`
	VShieldedSpend  int           `json:"vShieldedSpend"`
	ValueBalance    float64       `json:"valueBalance"`
	Value           float64       `json:"value"`
	OutputValue     float64       `json:"outputValue"`
	ShieldedValue   int           `json:"shieldedValue"`
	Overwintered    bool          `json:"overwintered"`
}

type Vin struct {
	Coinbase      string `json:"coinbase"`
	RetrievedVout *Vout  `json:"retrievedVout"`
	ScriptSig     struct {
		Asm string `json:"asm"`
		Hex string `json:"hex"`
	} `json:"scriptSig"`
	Sequence int    `json:"sequence"`
	Txid     string `json:"txid"`
	Vout     int    `json:"vout"`
}

type Vout struct {
	N            int `json:"n"`
	ScriptPubKey struct {
		Addresses []string `json:"addresses"`
		Asm       string   `json:"asm"`
		Hex       string   `json:"hex"`
		ReqSigs   int      `json:"reqSigs"`
		Type      string   `json:"type"`
	} `json:"scriptPubKey"`
	Value    float64 `json:"value"`
	ValueZat int     `json:"valueZat"`
}

func (v *Vout) IRawOutput() blockexplorer.IRawAddrOutput {
	return blockexplorer.IRawAddrOutput{
		Addresses: v.ScriptPubKey.Addresses,
		N:         v.N,
		Script:    "",
		Spent:     false,
		TxIndex:   0,
		Type:      v.ScriptPubKey.Type,
		Value:     idaemon.Amount(v.ValueZat),
	}
}

type Network struct {
	Name            string  `json:"name"`
	Accounts        int     `json:"accounts"`
	Transactions    int     `json:"transactions"`
	BlockHash       string  `json:"blockHash"`
	BlockNumber     int     `json:"blockNumber"`
	Difficulty      float64 `json:"difficulty"`
	HashRate        int64   `json:"hashrate"`
	MeanBlockTime   float64 `json:"meanBlockTime"`
	PeerCount       int     `json:"peerCount"`
	ProtocolVersion int     `json:"protocolVersion"`
	RelayFee        float64 `json:"relayFee"`
	Version         int     `json:"version"`
	SubVersion      string  `json:"subVersion"`
	TotalAmount     float64 `json:"totalAmount"`
	SproutPool      float64 `json:"sproutPool"`
	SaplingPool     float64 `json:"saplingPool"`
}

type Account struct {
	Address    string  `json:"address"`
	Balance    float64 `json:"balance"`
	FirstSeen  int     `json:"firstSeen"`
	LastSeen   int     `json:"lastSeen"`
	SentCount  int     `json:"sentCount"`
	RecvCount  int     `json:"recvCount"`
	MinedCount int     `json:"minedCount"`
	TotalSent  float64 `json:"totalSent"`
	TotalRecv  float64 `json:"totalRecv"`
}

func (a *Account) acount() *blockexplorer.IRawAddrResponse {
	var iRaw = blockexplorer.IRawAddrResponse{
		Address: a.Address,
		//FinalBalance:  a.Balance,
		Hash160: "",
		NTx:     0,
		//TotalReceived: a.TotalRecv,
		//TotalSent:     a.TotalSent,
		Txs: nil,
	}
	return &iRaw
}

func convertTxs(txs []Transaction) []blockexplorer.IRawAddrTx {
	var iRaws = make([]blockexplorer.IRawAddrTx, len(txs))
	for i, tx := range txs {
		iraw := blockexplorer.IRawAddrTx{
			BlockHeight:   tx.BlockHeight,
			Hash:          tx.Hash,
			Inputs:        tx.iRawInputs(),
			LockTime:      tx.LockTime,
			Outputs:       tx.iRawOutputs(),
			RelayedBy:     "",
			Result:        0,
			Size:          0,
			Time:          tx.Time,
			TxIndex:       tx.Index,
			Version:       tx.Version,
			VinSz:         0,
			VoutSz:        0,
			Weight:        0,
			Confirmations: 0,
		}
		iRaws[i] = iraw
	}
	return iRaws
}
