package shapeshift

import (
	"code.cryptopower.dev/exchange/lightningswap"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	API_BASE                   = "https://shapeshift.io/" // API endpoint
	DEFAULT_HTTPCLIENT_TIMEOUT = 30                       // HTTP client timeout
	LIBNAME                    = "shapeshift"
	waitSec                    = 3
)

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

// New return a instanciate struct
func New(conf lightningswap.ExchangeConfig) (*ShapeShift, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := lightningswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", conf.ApiKey))
		return nil
	})
	return &ShapeShift{
		client: client,
		conf:   &conf,
	}, nil
}

//ShapeShift represent a ShapeShift client
type ShapeShift struct {
	client *lightningswap.Client
	conf   *lightningswap.ExchangeConfig
	lightningswap.IDExchange
}

//SetDebug set enable/disable http request/response dump
func (c *ShapeShift) SetDebug(enable bool) {
	c.conf.Debug = enable
}

//CalculateExchangeRate get estimate on the amount for the exchange
func (c *ShapeShift) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
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

	res = lightningswap.ExchangeRateInfo{
		ExchangeRate:    exchangeRate,
		Min:             response.Min,
		Max:             response.Limit,
		EstimatedAmount: estAmount,
	}

	return
}

//EstimateAmount get estimate on the amount for the exchange
func (c *ShapeShift) EstimateAmount(vars interface{}) (res lightningswap.EstimateAmount, err error) {
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
	//amountStr := strconv.FormatFloat(amount, 'f', 8, 64)

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
	res = lightningswap.EstimateAmount{
		EstimatedAmount: estAmount,
	}

	return
}

//QueryRates (list of pairs LTC-BTC, BTC-LTC, etc)
func (c *ShapeShift) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryActiveCurrencies get all active currencies
func (c *ShapeShift) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryLimits Get Exchange Rates (from, to)
func (c *ShapeShift) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {
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
	res = lightningswap.QueryLimits{
		Max: response.Limit,
		Min: response.Min,
	}
	return
}

//CreateOrder create an instant exchange order
func (c *ShapeShift) CreateOrder(orderInfo lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {

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

	res = lightningswap.CreateResultInfo{
		UUID:           tmp.UUID,
		Destination:    tmp.ToAddress,
		FromCurrency:   tmp.CurrencyFrom,
		ToCurrency:     tmp.CurrencyTo,
		DepositAddress: tmp.DepositAddress,
		/* ChargedFee:     tmp.APIExtraFee,
		ExtraID:        tmp.PayinExtraID,
		PayoutExtraID:  tmp.PayoutExtraID, */
	}
	return
}

//UpdateOrder not available for this exchange
func (c *ShapeShift) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

//CancelOrder not available for this exchange
func (c *ShapeShift) CancelOrder(orderID string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

//OrderInfo get information on orderid/uuid
func (c *ShapeShift) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
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

	// maybe use later, tmp.NetworkFee
	res = lightningswap.OrderInfoResult{
		ReceiveAmount: tmp.AmountReceiving, //only shows when complete
		TxID:          tmp.TxID,            //only shows when complete

		Status:         tmp.Status,
		InternalStatus: lightningswap.OrderStatus(GetLocalStatus(tmp.Status)),
	}
	return
}

//Possible transaction statuses
//new waiting confirming exchanging sending finished failed refunded expired

//GetLocalStatus translate local status to idexchange status id
func GetLocalStatus(status string) (iStatus int) {
	status = strings.ToLower(status)
	switch status {
	case "complete":
		return 1
	case "no_deposits":
		return 2
	case "confirming":
		return 3
	case "refunded":
		return 5
	case "expired":
		return 7
	case "new":
		return 8
	case "received":
		return 9
	case "sending":
		return 10
	case "failed":
		return 11
	default:
		return 0
	}
}

/* func (c *ShapeShift) CheckOrderStatus(vars interface{}) (res string, err error) {
	err = errors.New("changenow:error: not available for this exchange")
	return
} */
