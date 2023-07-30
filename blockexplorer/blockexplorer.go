package blockexplorer

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type Config struct {
	EnableOutput bool
	Symbol       string
	ApiKey       string
	Type         NetworkType
}

var driv = driver{
	mux:    new(sync.RWMutex),
	stack:  make(map[string]NewExplorerFunc),
	layer2: make(map[NetworkType]NewExplorerFunc),
}

type NewExplorerFunc func(conf Config) (IBlockExplorer, error)

type driver struct {
	mux    *sync.RWMutex
	stack  map[string]NewExplorerFunc
	layer2 map[NetworkType]NewExplorerFunc
}

func (d *driver) registerExplorer(symbol string, networkType NetworkType, newExplorer NewExplorerFunc) {
	d.mux.Lock()
	defer d.mux.Unlock()
	if symbol != "" {
		_, ok := d.stack[symbol]
		if ok {
			log.Panicf("[%s] explorer is registered", symbol)
		}
		d.stack[symbol] = newExplorer
	}
	if networkType != "" {
		if _, ok := d.stack[symbol]; !ok {
			d.layer2[networkType] = newExplorer
		}
	}
}

func (d *driver) newExplorer(conf Config) (IBlockExplorer, error) {
	d.mux.Lock()
	defer d.mux.Unlock()
	if conf.Type == "" {
		var symbol = strings.ToLower(conf.Symbol)
		newExplorer, ok := d.stack[symbol]
		if !ok {
			return nil, fmt.Errorf("[%s] explorer is not available yet", symbol)
		}
		return newExplorer(conf)
	} else {
		newExplorer, ok := d.layer2[conf.Type]
		if !ok {
			return nil, fmt.Errorf("[%s] explorer is not available yet", conf.Type)
		}
		return newExplorer(conf)
	}
}

func RegisterExplorer(symbol string, networkType NetworkType, newDriver NewExplorerFunc) {
	driv.registerExplorer(strings.ToLower(symbol), networkType, newDriver)
}

func NewExplorer(conf Config) (IBlockExplorer, error) {
	if conf.Type == "" {
		return driv.newExplorer(conf)
	}
	return driv.newExplorer(conf)
}

type IBlockExplorer interface {
	GetTransaction(txId string) (tx *ITransaction, err error)
	GetTxsForAddress(address string, limit int, viewKey string) (tx *IRawAddrResponse, err error)
	//VerifyTransaction verifies transaction based on values passed in
	VerifyTransaction(verifier TxVerifyRequest) (tx *ITransaction, err error)
	VerifyByAddress(req AddressVerifyRequest) (vr *VerifyResult, err error)
	//PushTx pushes a raw tx hash
	PushTx(rawTxHash string) (result string, err error)
}

type TxVerifyRequest struct {
	TxId      string
	Amount    float64
	CreatedAt int64
	Address   string
	Confirms  int
	// ViewKey is used for verify monero Tx. It is corresponding with wallet address
	ViewKey string
}

type AddressVerifyRequest struct {
	Address   string
	Amount    float64
	ViewKey   string
	Confirm   int
	Timestamp int
}

type VerifyResult struct {
	Seen                bool    `json:"seen"` //tx has been seen on block explorer but not verified
	Verified            bool    `json:"verified"`
	OrderedAmount       float64 `json:"ordered_amount"`
	BlockExplorerAmount float64 `json:"block_explorer_amount"`
	MissingAmount       float64 `json:"missing_amount"`
	MissingPercent      float64 `json:"missing_percent"`
}
