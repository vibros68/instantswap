package dcrexplorer

import (
	"encoding/json"
)

type jsonResponse struct {
	Success bool            `json:"Success"`
	Message string          `json:"Message"`
	Result  json.RawMessage `json:"Data"`
}
type jsonPrivResponse struct {
	Success bool            `json:"Success"`
	Message string          `json:"Error"`
	Result  json.RawMessage `json:"Data"`
}
type VIN struct {
	Amountin    float64 `json:"amountin"`
	Blockheight int     `json:"blockheight"`
	Blockindex  int     `json:"blockindex"`
	ScriptSig   struct {
		Asm string `json:"asm"`
		Hex string `json:"hex"`
	} `json:"scriptSig"`
	Sequence int    `json:"sequence"`
	Tree     int    `json:"tree"`
	Txid     string `json:"txid"`
	Vout     int    `json:"vout"`
}

type VOUT struct {
	N            int `json:"n"`
	ScriptPubKey struct {
		Addresses []string `json:"addresses"`
		Asm       string   `json:"asm"`
		ReqSigs   int      `json:"reqSigs"`
		Type      string   `json:"type"`
	} `json:"scriptPubKey"`
	Value   float64 `json:"value"`
	Version int     `json:"version"`
}
type DecodedTransaction struct {
	Expiry   int    `json:"expiry"`
	Locktime int    `json:"locktime"`
	Txid     string `json:"txid"`
	Version  int    `json:"version"`
	Vin      []VIN  `json:"vin"`
	Vout     []VOUT `json:"vout"`
}
type Transaction struct {
	Block struct {
		Blockhash   string `json:"blockhash"`
		Blockheight int    `json:"blockheight"`
		Blockindex  int    `json:"blockindex"`
		Blocktime   int    `json:"blocktime"`
		Time        int    `json:"time"`
	} `json:"block"`
	Confirmations int    `json:"confirmations"`
	Expiry        int    `json:"expiry"`
	Locktime      int    `json:"locktime"`
	Size          int    `json:"size"`
	Txid          string `json:"txid"`
	Version       int    `json:"version"`
	Vin           []VIN  `json:"vin"`
	Vout          []VOUT `json:"vout"`
}

type RawAddrTx struct {
	Blockhash     string          `json:"blockhash"`
	Blocktime     int             `json:"blocktime"`
	Confirmations int             `json:"confirmations"`
	Locktime      int             `json:"locktime"`
	Size          int             `json:"size"`
	Time          int             `json:"time"`
	Txid          string          `json:"txid"`
	Version       int             `json:"version"`
	Vin           []RawAddrInput  `json:"vin"`
	Vout          []RawAddrOutput `json:"vout"`
}
type RawAddrInput struct {
	Amountin    float64          `json:"amountin"`
	Blockheight int              `json:"blockheight"`
	Blockindex  int              `json:"blockindex"`
	PrevOut     RawAddrPrevOut   `json:"prevOut"`
	ScriptSig   RawAddrScriptSig `json:"scriptSig"`
	Sequence    int              `json:"sequence"`
	Tree        int              `json:"tree"`
	Txid        string           `json:"txid"`
	Vout        int              `json:"vout"`
}
type RawAddrPrevOut struct {
	Addresses []string `json:"addresses"`
	Value     float64  `json:"value"`
}
type RawAddrScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}
type RawAddrOutput struct {
	N            int                 `json:"n"`
	ScriptPubKey RawAddrScriptPubKey `json:"scriptPubKey"`
	Value        float64             `json:"value"`
	Version      int                 `json:"version"`
}
type RawAddrScriptPubKey struct {
	Addresses []string `json:"addresses"`
	Asm       string   `json:"asm"`
	ReqSigs   int      `json:"reqSigs"`
	Type      string   `json:"type"`
}
type PushTxRequest struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}
