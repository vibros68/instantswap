package simpleswap

import (
	"code.cryptopower.dev/exchange/lightningswap"
	"code.cryptopower.dev/exchange/lightningswap/utils"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	API_BASE = "https://api.simpleswap.io/v1/"
	LIBNAME  = "simpleswap"
)

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

type SimpleSwap struct {
	client *lightningswap.Client
	conf   *lightningswap.ExchangeConfig
	lightningswap.IDExchange
}

func New(conf lightningswap.ExchangeConfig) (*SimpleSwap, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := lightningswap.NewClient(LIBNAME, &conf)
	return &SimpleSwap{client: client, conf: &conf}, nil
}

//SetDebug set enable/disable http request/response dump
func (c *SimpleSwap) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *SimpleSwap) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
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
	return lightningswap.ExchangeRateInfo{
		Min:             0,
		Max:             0,
		ExchangeRate:    rate,
		EstimatedAmount: estimatedAmount,
		MaxOrder:        0,
		Signature:       "",
	}, err
}

func (c *SimpleSwap) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	return res, fmt.Errorf("not supported")
}
func (c *SimpleSwap) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	return
}
func (c *SimpleSwap) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {
	return
}
func (c *SimpleSwap) CreateOrder(vars lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
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
	res = lightningswap.CreateResultInfo{
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

//UpdateOrder accepts orderID value and more if needed per lib
func (c *SimpleSwap) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	return
}
func (c *SimpleSwap) CancelOrder(orderID string) (res string, err error) {
	return
}

//OrderInfo accepts orderID value and more if needed per lib
func (c *SimpleSwap) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
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
	return lightningswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     order.UpdatedAt,
		ReceiveAmount:  utils.StrToFloat(order.AmountTo),
		TxID:           order.TxTo,
		Status:         order.Status,
		InternalStatus: GetLocalStatus(order.Status),
		Confirmations:  "",
	}, nil
}
func (c *SimpleSwap) EstimateAmount(vars interface{}) (res lightningswap.EstimateAmount, err error) {
	return
}

func parseResponseData(data []byte, obj interface{}) error {
	var simpleSwapErr Error
	err := json.Unmarshal(data, &simpleSwapErr)
	if err != nil {
		return fmt.Errorf(string(data))
	}
	if err == nil && simpleSwapErr.Code > 0 && len(simpleSwapErr.Message) > 0 {
		return fmt.Errorf(simpleSwapErr.Message)
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}

//GetLocalStatus translate local status to lightningswap.Status
func GetLocalStatus(status string) lightningswap.Status {
	// closed, confirming, exchanging, expired, failed, finished, refunded, sending, verifying, waiting
	status = strings.ToLower(status)
	switch status {
	case "closed":
		return lightningswap.OrderStatusCanceled
	case "confirming":
		return lightningswap.OrderStatusDepositReceived
	case "exchanging":
		return lightningswap.OrderStatusExchanging
	case "expired":
		return lightningswap.OrderStatusExpired
	case "failed":
		return lightningswap.OrderStatusFailed
	case "finished":
		return lightningswap.OrderStatusCompleted
	case "refunded":
		return lightningswap.OrderStatusRefunded
	case "sending":
		return lightningswap.OrderStatusSending
	case "verifying":
		return lightningswap.OrderStatusDepositReceived
	case "waiting":
		return lightningswap.OrderStatusWaitingForDeposit
	default:
		return lightningswap.OrderStatusUnknown
	}
}
