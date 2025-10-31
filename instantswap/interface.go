package instantswap

import (
	"fmt"
	"log"
	"sync"
)

const TX_HASH_INTERNAL_TRANSFER = "Internal transfer"

type IDExchange interface {
	Name() string
	// GetCurrencies returns current active currencies in the exchange
	GetCurrencies() (currencies []Currency, err error)
	// GetCurrenciesToPair return a list of available currencies for a specific currency
	GetCurrenciesToPair(from string) (currencies []Currency, err error)
	QueryLimits(fromCurr, toCurr string) (res QueryLimits, err error)
	CreateOrder(vars CreateOrder) (res CreateResultInfo, err error)
	//UpdateOrder accepts orderID value and more if needed per lib
	UpdateOrder(vars interface{}) (res UpdateOrderResultInfo, err error)
	CancelOrder(orderID string) (res string, err error)

	//OrderInfo accepts orderID value and more if needed per lib
	OrderInfo(orderID string, extraIds ...string) (res OrderInfoResult, err error)

	GetExchangeRateInfo(vars ExchangeRateRequest) (res ExchangeRateInfo, err error)
}

type ExchangeRateRequest struct {
	From        string
	FromNetwork string
	To          string
	ToNetwork   string
	Amount      float64
}

var driv = driver{
	mux:   new(sync.RWMutex),
	stack: make(map[string]NewExchangeFunc),
}

type NewExchangeFunc func(config ExchangeConfig) (IDExchange, error)

type driver struct {
	mux   *sync.RWMutex
	stack map[string]NewExchangeFunc
}

func (d *driver) registerExchange(symbol string, newExchange NewExchangeFunc) {
	d.mux.Lock()
	defer d.mux.Unlock()
	_, ok := d.stack[symbol]
	if ok {
		log.Panicf("[%s] explorer is registered", symbol)
	}
	d.stack[symbol] = newExchange
}

func (d *driver) newExchange(name string, config ExchangeConfig) (IDExchange, error) {
	d.mux.Lock()
	defer d.mux.Unlock()
	newExplorer, ok := d.stack[name]
	if !ok {
		return nil, fmt.Errorf("[%s] exchange is not registered yet", name)
	}
	return newExplorer(config)
}

func RegisterExchange(symbol string, newExchange NewExchangeFunc) {
	driv.registerExchange(symbol, newExchange)
}

func NewExchange(symbol string, config ExchangeConfig) (IDExchange, error) {
	return driv.newExchange(symbol, config)
}
