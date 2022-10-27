package shapeshift

import (
	"code.cryptopower.dev/exchange/instantswap"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	API_BASE = "https://shapeshift.io/" // API endpoint
	LIBNAME  = "shapeshift"
)

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// New return a instanciate struct
func New(conf instantswap.ExchangeConfig) (*ShapeShift, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.ApiKey))
		return nil
	})
	return &ShapeShift{
		client: client,
		conf:   &conf,
	}, nil
}

// ShapeShift represent a ShapeShift client.
type ShapeShift struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
	instantswap.IDExchange
}

// SetDebug set enable/disable http request/response dump.
func (c *ShapeShift) SetDebug(enable bool) {
	c.conf.Debug = enable
}

// GetExchangeRateInfo get estimate on the amount for the exchange.
func (c *ShapeShift) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	pair := strings.ToLower(vars.From) + "_" + strings.ToLower(vars.To)
	r, err := c.client.Do(API_BASE, "GET", "marketinfo/"+pair, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var response MarketInfoResponse
	if err = json.Unmarshal(r, &response); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	exchangeRate := 1 / response.Rate
	estAmount := vars.Amount / exchangeRate

	res = instantswap.ExchangeRateInfo{
		ExchangeRate:    exchangeRate,
		Min:             response.Min,
		Max:             response.Limit,
		EstimatedAmount: estAmount,
	}

	return
}

// EstimateAmount get estimate on the amount for the exchange.
func (c *ShapeShift) EstimateAmount(vars interface{}) (res instantswap.EstimateAmount, err error) {
	vals := vars.(map[string]interface{})
	var to string
	var from string
	var amount float64
	for k, v := range vals {
		if k == "to" {
			to = v.(string)
		} else if k == "from" {
			from = v.(string)
		} else if k == "amount" {
			amount = v.(float64)
		}
	}
	pair := from + "_" + to

	r, err := c.client.Do(API_BASE, "GET", "rate/"+pair, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var response RateResponse
	if err = json.Unmarshal(r, &response); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	if response.ErrorMsg != "" {
		err = errors.New(LIBNAME + ":error: " + response.ErrorMsg)
		return
	}
	exchangeRate := 1 / response.Rate
	estAmount := amount / exchangeRate
	res = instantswap.EstimateAmount{
		EstimatedAmount: estAmount,
	}

	return
}

// QueryRates (list of pairs LTC-BTC, BTC-LTC, etc)
func (c *ShapeShift) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

// QueryActiveCurrencies get all active currencies.
func (c *ShapeShift) QueryActiveCurrencies(vars interface{}) (res []instantswap.ActiveCurr, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

// QueryLimits Get Exchange Rates (from, to).
func (c *ShapeShift) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	pair := strings.ToLower(fromCurr) + "_" + strings.ToLower(toCurr)
	r, err := c.client.Do(API_BASE, "GET", "marketinfo/"+pair, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var response MarketInfoResponse
	if err = json.Unmarshal(r, &response); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	if response.ErrorMsg != "" {
		err = errors.New(LIBNAME + ":error: " + response.ErrorMsg)
		return
	}
	res = instantswap.QueryLimits{
		Max: response.Limit,
		Min: response.Min,
	}
	return
}

// CreateOrder create an instant exchange order.
func (c *ShapeShift) CreateOrder(orderInfo instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {

	tmpOrderInfo := CreateOrder{
		ToCurrencyAddress: orderInfo.Destination,
		Pair:              strings.ToLower(orderInfo.FromCurrency) + "_" + strings.ToLower(orderInfo.ToCurrency),
		RefundAddress:     orderInfo.RefundAddress,
		APIKEY:            c.conf.ApiKey,
	}

	payload, err := json.Marshal(tmpOrderInfo)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	r, err := c.client.Do(API_BASE, "POST", "shift", string(payload), false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	var tmp CreateResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if tmp.ErrorMsg != "" {
		err = errors.New(LIBNAME + ":error: " + tmp.ErrorMsg)
		return
	}

	res = instantswap.CreateResultInfo{
		UUID:           tmp.UUID,
		Destination:    tmp.ToAddress,
		FromCurrency:   tmp.CurrencyFrom,
		ToCurrency:     tmp.CurrencyTo,
		DepositAddress: tmp.DepositAddress,
	}
	return
}

// UpdateOrder not available for this exchange.
func (c *ShapeShift) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

// CancelOrder not available for this exchange.
func (c *ShapeShift) CancelOrder(orderID string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

// OrderInfo get information on orderid/uuid.
func (c *ShapeShift) OrderInfo(orderID string) (res instantswap.OrderInfoResult, err error) {
	r, err := c.client.Do(API_BASE, "GET", "txStat/"+orderID, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmp OrderStatusResponse
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if tmp.Error != "" {
		err = errors.New(LIBNAME + ":error: " + tmp.Error)
		return
	}

	res = instantswap.OrderInfoResult{
		ReceiveAmount: tmp.AmountReceiving, //only shows when complete
		TxID:          tmp.TxID,            //only shows when complete

		Status:         tmp.Status,
		InternalStatus: GetLocalStatus(tmp.Status),
	}
	return
}

// GetLocalStatus converts local status to instantswap.Status.
// possible transaction statuses is:
// new waiting confirming exchanging sending finished failed refunded expired
func GetLocalStatus(status string) instantswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "complete":
		return instantswap.OrderStatusCompleted
	case "no_deposits":
		return instantswap.OrderStatusNew
	case "confirming":
		return instantswap.OrderStatusDepositReceived
	case "refunded":
		return instantswap.OrderStatusRefunded
	case "expired":
		return instantswap.OrderStatusExpired
	case "new":
		return instantswap.OrderStatusNew
	case "received":
		return instantswap.OrderStatusDepositReceived
	case "sending":
		return instantswap.OrderStatusSending
	case "failed":
		return instantswap.OrderStatusFailed
	default:
		return instantswap.OrderStatusUnknown
	}
}
