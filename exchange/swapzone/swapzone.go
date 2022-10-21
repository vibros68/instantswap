package swapzone

import (
	"code.cryptopower.dev/exchange/lightningswap"
	"code.cryptopower.dev/exchange/lightningswap/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	API_BASE = "https://api.swapzone.io/v1/"
	LIBNAME  = "swapzone"
)

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

// New return a SwapZone client.
func New(conf lightningswap.ExchangeConfig) (*SwapZone, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := lightningswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		r.Header.Set("x-api-key", conf.ApiKey)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return nil
	})
	return &SwapZone{client: client, conf: &conf}, nil
}

// SwapZone represent a SwapZone client.
type SwapZone struct {
	client *lightningswap.Client
	conf   *lightningswap.ExchangeConfig
	lightningswap.IDExchange
}

// SetDebug set enable/disable http request/response dump.
func (c *SwapZone) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *SwapZone) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, "GET",
		fmt.Sprintf("exchange/get-rate?from=%s&to=%s&amount=%.8f&rateType=all&availableInUSA=false&chooseRate=best&noRefundAddress=false",
			strings.ToLower(vars.From), strings.ToLower(vars.To), vars.Amount),
		"", false)
	if err != nil {
		return
	}
	var exchangeRate ExchangeRate
	err = parseResponseData(r, &exchangeRate)
	if err != nil {
		return
	}
	res.Min = exchangeRate.MinAmount
	res.Max = exchangeRate.MaxAmount
	res.EstimatedAmount = exchangeRate.AmountTo
	res.ExchangeRate = exchangeRate.AmountFrom / exchangeRate.AmountTo
	res.Signature = exchangeRate.QuotaId
	return
}

func (c *SwapZone) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	return res, fmt.Errorf("not supported")
}

func (c *SwapZone) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	return
}

func (c *SwapZone) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {
	return
}

func (c *SwapZone) CreateOrder(vars lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
	var form = make(url.Values)
	form.Set("from", strings.ToLower(vars.FromCurrency))
	form.Set("to", strings.ToLower(vars.ToCurrency))
	form.Set("amountDeposit", fmt.Sprintf("%.8f", vars.InvoicedAmount))
	form.Set("addressReceive", vars.Destination)
	form.Set("extraIdReceive", "") // Memo tag (optional)
	form.Set("refundAddress", vars.RefundAddress)
	form.Set("refundExtraId", "") // Memo tag for refund address (optional)
	if len(vars.Signature) > 0 {
		form.Set("quotaId", vars.Signature)
	}

	var r []byte
	r, err = c.client.Do(API_BASE, "POST", "exchange/create",
		form.Encode(), false)
	if err != nil {
		return
	}
	var tx Transaction
	err = parseResponseData(r, &tx)
	if err != nil {
		return
	}
	var order = tx.Transaction
	var invoicedAmount = utils.StrToFloat(order.AmountDeposit)
	var orderedAmount = utils.StrToFloat(order.AmountEstimated)
	res = lightningswap.CreateResultInfo{
		ChargedFee:     0,
		Destination:    order.AddressReceive,
		ExchangeRate:   invoicedAmount / orderedAmount,
		FromCurrency:   order.From,
		InvoicedAmount: invoicedAmount,
		OrderedAmount:  orderedAmount,
		ToCurrency:     order.To,
		UUID:           order.Id,
		DepositAddress: order.AddressDeposit,
		Expires:        0,
		ExtraID:        "",
		PayoutExtraID:  "",
	}
	return
}

// UpdateOrder accepts orderID value and more if needed per lib.
func (c *SwapZone) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	return
}
func (c *SwapZone) CancelOrder(orderID string) (res string, err error) {
	return
}

// OrderInfo accepts orderID value and more if needed per lib.
func (c *SwapZone) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, "GET",
		fmt.Sprintf("exchange/tx?id=%s", orderID),
		"", false)
	if err != nil {
		return
	}
	var tx Transaction
	err = parseResponseData(r, &tx)
	if err != nil {
		return
	}
	var order = tx.Transaction
	res = lightningswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     "",
		ReceiveAmount:  0,
		TxID:           "",
		Status:         order.Status,
		InternalStatus: GetLocalStatus(order.Status),
		Confirmations:  "",
	}
	if order.Status == "" {
		res.ReceiveAmount = utils.StrToFloat(order.AmountEstimated)
	}
	return
}

func (c *SwapZone) EstimateAmount(vars interface{}) (res lightningswap.EstimateAmount, err error) {
	return
}

func parseResponseData(data []byte, obj interface{}) error {
	var swapzoneErr SwapzoneError
	err := json.Unmarshal(data, &swapzoneErr)
	if err == nil && swapzoneErr.Error {
		return fmt.Errorf(swapzoneErr.Message)
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}

// GetLocalStatus translate local status to lightningswap.Status.
func GetLocalStatus(status string) lightningswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "waiting":
		return lightningswap.OrderStatusNew
	case "confirming":
		return lightningswap.OrderStatusDepositReceived
	case "exchanging":
		return lightningswap.OrderStatusExchanging
	case "sending":
		return lightningswap.OrderStatusSending
	case "finished":
		return lightningswap.OrderStatusCompleted
	case "refunded":
		return lightningswap.OrderStatusRefunded
	case "failed":
		return lightningswap.OrderStatusFailed
	case "overdue":
		return lightningswap.OrderStatusExpired
	default:
		return lightningswap.OrderStatusUnknown
	}
}
