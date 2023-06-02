package swapzone

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gitlab.com/cryptopower/instantswap/instantswap"
	"gitlab.com/cryptopower/instantswap/instantswap/utils"
)

const (
	API_BASE = "https://api.swapzone.io/v1/"
	LIBNAME  = "swapzone"
)

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// New return a SwapZone client.
func New(conf instantswap.ExchangeConfig) (*SwapZone, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		r.Header.Set("x-api-key", conf.ApiKey)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return nil
	})
	return &SwapZone{client: client, conf: &conf}, nil
}

// SwapZone represent a SwapZone client.
type SwapZone struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
	instantswap.IDExchange
}

// SetDebug set enable/disable http request/response dump.
func (c *SwapZone) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *SwapZone) GetCurrencies() (currencies []instantswap.Currency, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, "GET", "exchange/currencies", "", false)
	if err != nil {
		return
	}
	var szCurrencies []Currency
	err = parseResponseData(r, &szCurrencies)
	if err != nil {
		return
	}
	currencies = make([]instantswap.Currency, len(szCurrencies))
	for i, curr := range szCurrencies {
		currencies[i] = instantswap.Currency{
			Name:   curr.Name,
			Symbol: curr.Ticker,
		}
	}
	return
}
func (c *SwapZone) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, "GET", "exchange/currencies", "", false)
	if err != nil {
		return
	}
	var szCurrencies []Currency
	err = parseResponseData(r, &szCurrencies)
	if err != nil {
		return
	}
	currencies = []instantswap.Currency{}
	for _, curr := range szCurrencies {
		if strings.ToLower(curr.Ticker) == strings.ToLower(from) {
			continue
		}
		currencies = append(currencies, instantswap.Currency{
			Name:   curr.Name,
			Symbol: curr.Ticker,
		})
	}
	return
}

func (c *SwapZone) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
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

func (c *SwapZone) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	return res, fmt.Errorf("not supported")
}

func (c *SwapZone) QueryActiveCurrencies(vars interface{}) (res []instantswap.ActiveCurr, err error) {
	return
}

func (c *SwapZone) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return
}

func (c *SwapZone) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
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
	res = instantswap.CreateResultInfo{
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
func (c *SwapZone) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (c *SwapZone) CancelOrder(orderID string) (res string, err error) {
	return
}

// OrderInfo accepts orderID value and more if needed per lib.
func (c *SwapZone) OrderInfo(orderID string) (res instantswap.OrderInfoResult, err error) {
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
	res = instantswap.OrderInfoResult{
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

func (c *SwapZone) EstimateAmount(vars interface{}) (res instantswap.EstimateAmount, err error) {
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

// GetLocalStatus translate local status to instantswap.Status.
func GetLocalStatus(status string) instantswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "waiting":
		return instantswap.OrderStatusNew
	case "confirming":
		return instantswap.OrderStatusDepositReceived
	case "exchanging":
		return instantswap.OrderStatusExchanging
	case "sending":
		return instantswap.OrderStatusSending
	case "finished":
		return instantswap.OrderStatusCompleted
	case "refunded":
		return instantswap.OrderStatusRefunded
	case "failed":
		return instantswap.OrderStatusFailed
	case "overdue":
		return instantswap.OrderStatusExpired
	default:
		return instantswap.OrderStatusUnknown
	}
}
