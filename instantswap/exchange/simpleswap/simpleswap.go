package simpleswap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"code.cryptopower.dev/group/instantswap/instantswap"
	"code.cryptopower.dev/group/instantswap/instantswap/utils"
)

const (
	API_BASE = "https://api.simpleswap.io/v1/"
	LIBNAME  = "simpleswap"
)

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

type SimpleSwap struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
	instantswap.IDExchange
}

func New(conf instantswap.ExchangeConfig) (*SimpleSwap, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf)
	return &SimpleSwap{client: client, conf: &conf}, nil
}

// SetDebug set enable/disable http request/response dump
func (c *SimpleSwap) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *SimpleSwap) GetCurrencies() (currencies []instantswap.Currency, err error) {
	r, err := c.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("get_all_currencies?api_key=%s", c.conf.ApiKey),
		"", false)
	if err != nil {
		return
	}
	var ssCurrencies []Currency
	err = parseResponseData(r, &ssCurrencies)
	if err != nil {
		return
	}
	currencies = make([]instantswap.Currency, len(ssCurrencies))
	for i, curr := range ssCurrencies {
		currencies[i] = instantswap.Currency{
			Name:   curr.Name,
			Symbol: curr.Symbol,
		}
	}
	return currencies, nil
}

func (c *SimpleSwap) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	r, err := c.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("get_pairs?api_key=%s&fixed=true&symbol=%s", c.conf.ApiKey, strings.ToLower(from)),
		"", false)
	if err != nil {
		return
	}
	var ssCurrencies []string
	err = parseResponseData(r, &ssCurrencies)
	if err != nil {
		return
	}
	currencies = make([]instantswap.Currency, len(ssCurrencies))
	for i, curr := range ssCurrencies {
		currencies[i] = instantswap.Currency{
			Name:   curr,
			Symbol: curr,
		}
	}
	return
}

func (c *SimpleSwap) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, "GET",
		fmt.Sprintf("get_estimated?api_key=%s&currency_from=%s&currency_to=%s&fixed=true&amount=%.8f",
			c.conf.ApiKey, strings.ToLower(vars.From), strings.ToLower(vars.To), vars.Amount),
		"", false)
	if err != nil {
		return
	}
	var response = string(r)
	if response == "null" {
		return res, fmt.Errorf("invalid request")
	}
	estimatedAmount := utils.StrToFloat(response)
	rate := vars.Amount / estimatedAmount
	return instantswap.ExchangeRateInfo{
		Min:             0,
		Max:             0,
		ExchangeRate:    rate,
		EstimatedAmount: estimatedAmount,
		MaxOrder:        0,
		Signature:       "",
	}, err
}

func (c *SimpleSwap) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	return res, fmt.Errorf("not supported")
}

func (c *SimpleSwap) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	var form = CreateExchange{
		CurrencyFrom:      strings.ToLower(vars.FromCurrency),
		CurrencyTo:        strings.ToLower(vars.ToCurrency),
		Fixed:             false,
		Amount:            vars.OrderedAmount,
		AddressTo:         vars.Destination,
		ExtraIdTo:         "",
		UserRefundAddress: vars.RefundAddress,
		UserRefundExtraId: "",
		Referral:          vars.Signature,
	}
	payload, err := json.Marshal(form)
	if err != nil {
		return res, err
	}
	// do request
	var r []byte
	r, err = c.client.Do(API_BASE, "POST", fmt.Sprintf("create_exchange?api_key=%s", c.conf.ApiKey),
		string(payload), false)
	if err != nil {
		return
	}
	var order Order
	err = parseResponseData(r, &order)
	if err != nil {
		return
	}
	var invoicedAmount = utils.StrToFloat(order.AmountFrom)
	var orderedAmount = utils.StrToFloat(order.AmountTo)
	res = instantswap.CreateResultInfo{
		ChargedFee:     0,
		Destination:    order.AddressTo,
		ExchangeRate:   invoicedAmount / orderedAmount,
		FromCurrency:   order.CurrencyFrom,
		InvoicedAmount: invoicedAmount,
		OrderedAmount:  orderedAmount,
		ToCurrency:     order.CurrencyTo,
		UUID:           order.Id,
		DepositAddress: order.UserRefundAddress,
		Expires:        0,
		ExtraID:        "",
		PayoutExtraID:  "",
	}
	return
}

// UpdateOrder accepts orderID value and more if needed per lib
func (c *SimpleSwap) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}

func (c *SimpleSwap) CancelOrder(orderID string) (res string, err error) {
	return
}

// OrderInfo accepts orderID value and more if needed per lib.
func (c *SimpleSwap) OrderInfo(orderID string) (res instantswap.OrderInfoResult, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, "GET",
		fmt.Sprintf("get_exchange?id=%s&api_key=%s", orderID, c.conf.ApiKey),
		"", false)
	if err != nil {
		return
	}
	var order Order
	err = parseResponseData(r, &order)
	if err != nil {
		return
	}
	return instantswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     order.UpdatedAt,
		ReceiveAmount:  utils.StrToFloat(order.AmountTo),
		TxID:           order.TxTo,
		Status:         order.Status,
		InternalStatus: GetLocalStatus(order.Status),
		Confirmations:  "",
	}, nil
}

func (c *SimpleSwap) EstimateAmount(vars interface{}) (res instantswap.EstimateAmount, err error) {
	return
}

func parseResponseData(data []byte, obj interface{}) error {
	var simpleSwapErr Error
	err := json.Unmarshal(data, &simpleSwapErr)
	/*if err != nil {
		return fmt.Errorf(string(data))
	}*/
	if err == nil && simpleSwapErr.Code > 0 && len(simpleSwapErr.Message) > 0 {
		return fmt.Errorf(simpleSwapErr.Message)
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}

// GetLocalStatus translate local status to instantswap.Status.
func GetLocalStatus(status string) instantswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "closed":
		return instantswap.OrderStatusCanceled
	case "confirming":
		return instantswap.OrderStatusDepositReceived
	case "exchanging":
		return instantswap.OrderStatusExchanging
	case "expired":
		return instantswap.OrderStatusExpired
	case "failed":
		return instantswap.OrderStatusFailed
	case "finished":
		return instantswap.OrderStatusCompleted
	case "refunded":
		return instantswap.OrderStatusRefunded
	case "sending":
		return instantswap.OrderStatusSending
	case "verifying":
		return instantswap.OrderStatusDepositReceived
	case "waiting":
		return instantswap.OrderStatusWaitingForDeposit
	default:
		return instantswap.OrderStatusUnknown
	}
}
