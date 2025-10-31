package easybit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vibros68/instantswap/instantswap"
	"github.com/vibros68/instantswap/instantswap/utils"
)

const (
	API_BASE = "https://api.easybit.com/" // API endpoint
	LIBNAME  = "easybit"
)

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// New return an EasyBit api client
func New(conf instantswap.ExchangeConfig) (*EasyBit, error) {
	if conf.ApiKey == "" {
		err := fmt.Errorf("APIKEY is blank")
		return nil, err
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		r.Header.Set("API-KEY", conf.ApiKey)
		return nil
	})
	return &EasyBit{
		client: client,
		conf:   &conf,
	}, nil
}

// EasyBit represent a EasyBit client.
type EasyBit struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
	instantswap.IDExchange
}

func (c *EasyBit) Name() string {
	return LIBNAME
}

// SetDebug set enable/disable http request/response dump.
func (c *EasyBit) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *EasyBit) GetCurrencies() (currencies []instantswap.Currency, err error) {
	r, err := c.client.Do(API_BASE, "GET", "currencyList", "", false)
	if err != nil {
		return nil, err
	}
	var ebCurrencies []Currency
	err = parseDataResponse(r, &ebCurrencies)
	if err != nil {
		return nil, err
	}
	currencies = make([]instantswap.Currency, len(ebCurrencies))
	for i, currency := range ebCurrencies {
		currencies[i] = instantswap.Currency{
			Name:   currency.Name,
			Symbol: currency.Currency,
		}
	}
	return currencies, nil
}

func (c *EasyBit) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	r, err := c.client.Do(API_BASE, "GET", "currencyList", "", false)
	if err != nil {
		return nil, err
	}
	var ebCurrencies []Currency
	err = parseDataResponse(r, &ebCurrencies)
	if err != nil {
		return nil, err
	}
	for _, currency := range ebCurrencies {
		if strings.ToLower(from) != strings.ToLower(currency.Currency) {
			currencies = append(currencies, instantswap.Currency{
				Name:   currency.Name,
				Symbol: currency.Currency,
			})
		}
	}
	return currencies, nil
}

func (c *EasyBit) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	r, err := c.client.Do(API_BASE, "GET",
		fmt.Sprintf("rate?send=%s&receive=%s&amount=%.8f", vars.From, vars.To, vars.Amount), "", false)
	if err != nil {
		return res, err
	}
	var rate ExchangeRate
	err = parseDataResponse(r, &rate)
	if err != nil {
		return res, err
	}
	pairInfo, _ := c.pairInfo(vars)
	return instantswap.ExchangeRateInfo{
		Min:             utils.StrToFloat(pairInfo.MinimumAmount),
		Max:             utils.StrToFloat(pairInfo.MaximumAmount),
		ExchangeRate:    utils.StrToFloat(rate.Rate),
		EstimatedAmount: utils.StrToFloat(rate.ReceiveAmount),
		MaxOrder:        0,
		Signature:       "",
	}, nil
}

func (c *EasyBit) pairInfo(vars instantswap.ExchangeRateRequest) (PairInfo, error) {
	r, err := c.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("pairInfo?send=%s&receive=%s", vars.From, vars.To), "", false)
	if err != nil {
		return PairInfo{}, err
	}
	var pair PairInfo
	err = parseDataResponse(r, &pair)
	return pair, err
}

func (c *EasyBit) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return
}
func (c *EasyBit) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	var orderRequest = map[string]string{
		"send":           vars.FromCurrency,
		"receive":        vars.ToCurrency,
		"amount":         fmt.Sprintf("%.8f", vars.InvoicedAmount),
		"receiveAddress": vars.Destination,
	}
	payload, err := json.Marshal(orderRequest)
	if err != nil {
		return res, err
	}
	r, err := c.client.Do(API_BASE, http.MethodPost, "order", string(payload), false)
	if err != nil {
		return res, err
	}
	var order Order
	err = parseDataResponse(r, &order)
	return instantswap.CreateResultInfo{
		ChargedFee:     0,
		Destination:    order.ReceiveAddress,
		ExchangeRate:   0,
		FromCurrency:   order.Send,
		InvoicedAmount: utils.StrToFloat(order.SendAmount),
		OrderedAmount:  utils.StrToFloat(order.ReceiveAmount),
		ToCurrency:     order.Receive,
		UUID:           order.Id,
		DepositAddress: order.SendAddress,
		Expires:        0,
		ExtraID:        "",
		PayoutExtraID:  "",
	}, err
}

// UpdateOrder accepts orderID value and more if needed per lib
func (c *EasyBit) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (c *EasyBit) CancelOrder(orderID string) (res string, err error) {
	return
}

func (c *EasyBit) OrderInfo(orderID string, extraIds ...string) (res instantswap.OrderInfoResult, err error) {
	r, err := c.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("orders?id=%s", orderID), "", false)
	if err != nil {
		return res, err
	}
	var orders []Order
	err = parseDataResponse(r, &orders)
	if err != nil {
		return res, err
	}
	for _, order := range orders {
		if order.Id == orderID {
			var txId string
			if order.HashOut != nil {
				txId = order.HashOut.(string)
			}
			return instantswap.OrderInfoResult{
				Expires:        0,
				LastUpdate:     "",
				ReceiveAmount:  utils.StrToFloat(order.ReceiveAmount),
				TxID:           txId,
				Status:         order.Status,
				InternalStatus: mapOrderStatus(order.Status),
				Confirmations:  "",
			}, nil
		}
	}
	return res, fmt.Errorf("order[%s] not found", orderID)
}

// "Refund" or "Failed" or "Volatility Protection" or "Action Request" or "Request Overdue"
func mapOrderStatus(status string) instantswap.Status {
	switch status {
	case "Awaiting Deposit":
		return instantswap.OrderStatusWaitingForDeposit
	case "Confirming Deposit":
		return instantswap.OrderStatusDepositReceived
	case "Exchanging":
		return instantswap.OrderStatusExchanging
	case "Sending":
		return instantswap.OrderStatusSending
	case "Complete":
		return instantswap.OrderStatusCompleted
	case "Refund":
		return instantswap.OrderStatusRefunded
	case "Failed":
		return instantswap.OrderStatusFailed
	case "Volatility Protection":
		return instantswap.OrderStatusExchanging
	case "Action Request":
		return instantswap.OrderStatusExchanging
	case "Request Overdue":
		return instantswap.OrderStatusExpired
	default:
		return instantswap.OrderStatusUnknown
	}
}
