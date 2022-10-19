package changenow

import (
	"code.cryptopower.dev/exchange/lightningswap"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE                   = "https://changenow.io/api/v1/" // API endpoint
	DEFAULT_HTTPCLIENT_TIMEOUT = 30                             // HTTP client timeout
	LIBNAME                    = "changenow"
	waitSec                    = 3
)

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

// New return an ChangeNow client struct with IDExchange implement
func New(conf lightningswap.ExchangeConfig) (*ChangeNow, error) {
	if conf.ApiKey == "" {
		err := fmt.Errorf("APIKEY is blank")
		return nil, err
	}
	client := lightningswap.NewClient(LIBNAME, &conf)
	return &ChangeNow{client: client, conf: &conf}, nil
}

//ChangeNow represent a ChangeNow client
type ChangeNow struct {
	conf   *lightningswap.ExchangeConfig
	client *lightningswap.Client
	lightningswap.IDExchange
}

//SetDebug set enable/disable http request/response dump
func (c *ChangeNow) SetDebug(enable bool) {
	c.conf.Debug = enable
}
func handleErr(r json.RawMessage) (err error) {
	var errorVals map[string][]string
	if err = json.Unmarshal(r, &errorVals); err != nil {
		return err
	}
	if len(errorVals) > 0 {
		var errorStr string
		errorStr = LIBNAME + " error(s): "
		for k, v := range errorVals {
			errorStr += k + ": "
			for i := 0; i < len(v); i++ {
				errorStr += v[i] + ", "
			}
		}
		err = errors.New(errorStr)
		return
	}
	return nil
}

//CalculateExchangeRate get estimate on the amount for the exchange
func (c *ChangeNow) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
	limits, err := c.QueryLimits(vars.From, vars.To)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	time.Sleep(time.Second * 1)
	estimate, err := c.EstimateAmount(vars)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	rate := vars.Amount / estimate.EstimatedAmount

	res = lightningswap.ExchangeRateInfo{
		ExchangeRate:    rate,
		Min:             limits.Min,
		Max:             limits.Max,
		EstimatedAmount: estimate.EstimatedAmount,
	}

	return
}

//EstimateAmount get estimate on the amount for the exchange
func (c *ChangeNow) EstimateAmount(vars lightningswap.ExchangeRateRequest) (res lightningswap.EstimateAmount, err error) {
	amountStr := strconv.FormatFloat(vars.Amount, 'f', 8, 64)
	r, err := c.client.Do(API_BASE, "GET",
		fmt.Sprintf("exchange-amount/%s/%s_%s?api_key=%s", amountStr, vars.From, vars.To, c.conf.ApiKey), "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmpRes EstimateAmount
	if err = json.Unmarshal(r, &tmpRes); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = lightningswap.EstimateAmount{
		EstimatedAmount:          tmpRes.EstimatedAmount,
		NetworkFee:               tmpRes.NetworkFee,
		ServiceCommission:        tmpRes.ServiceCommission,
		TransactionSpeedForecast: tmpRes.TransactionSpeedForecast,
		WarningMessage:           tmpRes.WarningMessage,
	}

	return
}

//QueryRates (list of pairs LTC-BTC, BTC-LTC, etc)
func (c *ChangeNow) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryActiveCurrencies get all active currencies
func (c *ChangeNow) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	//vars not used here
	r, err := c.client.Do(API_BASE, "GET", "currencies?active=true", "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmpArr []ActiveCurr
	if err = json.Unmarshal(r, &tmpArr); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	for _, v := range tmpArr {
		currType := "CRYPTO"
		if v.IsFiat {
			currType = "FIAT"
		}
		tmpItem := lightningswap.ActiveCurr{
			CurrencyType: currType,
			Name:         v.Ticker,
		}
		res = append(res, tmpItem)
	}
	return
}

//QueryLimits Get Exchange Rates (from, to)
func (c *ChangeNow) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {

	r, err := c.client.Do(API_BASE, "GET", "min-amount/"+fromCurr+"_"+toCurr, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmp QueryLimits
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	//minAmount := strconv.FormatFloat(tmp.Min, 'f', 8, 64)
	res = lightningswap.QueryLimits{
		Min: tmp.Min,
	}
	return
}

//CreateOrder create an instant exchange order
func (c *ChangeNow) CreateOrder(orderInfo lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {

	//convert from interface orderInfo to local struct
	tmpOrderInfo := CreateOrder{
		FromCurrency:      orderInfo.FromCurrency,
		ToCurrency:        orderInfo.ToCurrency,
		ToCurrencyAddress: orderInfo.Destination,
		InvoicedAmount:    orderInfo.InvoicedAmount,
		ExtraID:           orderInfo.ExtraID,
	}
	if tmpOrderInfo.InvoicedAmount == 0.0 {
		err = errors.New(LIBNAME + ":error:createorder invoiced amount is 0")
		return
	}
	payload, err := json.Marshal(tmpOrderInfo)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	r, err := c.client.Do(API_BASE, "POST", "transactions/"+c.conf.ApiKey, string(payload), false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tempItem interface{}
	err = json.Unmarshal(r, &tempItem)
	//fmt.Println(tempItem)
	var tmp CreateResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = lightningswap.CreateResultInfo{
		UUID:           tmp.UUID,
		Destination:    tmp.DestinationAddress,
		ExtraID:        tmp.PayinExtraID,
		FromCurrency:   tmp.FromCurrency,
		ToCurrency:     tmp.ToCurrency,
		DepositAddress: tmp.DepositAddress,
	}
	return
}

//UpdateOrder not available for this exchange
func (c *ChangeNow) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

//CancelOrder not available for this exchange
func (c *ChangeNow) CancelOrder(oId string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

//OrderInfo get information on orderid/uuid
func (c *ChangeNow) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
	r, err := c.client.Do(API_BASE, "GET", "transactions/"+orderID+"/"+c.conf.ApiKey, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmp OrderInfoResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var amountRecv float64
	if tmp.Status != "finished" {
		amountRecv = tmp.ExpectedAmountReceive
	} else {
		amountRecv = tmp.AmountReceive
	}
	// maybe use later, tmp.NetworkFee
	res = lightningswap.OrderInfoResult{
		//Expires:      not available
		LastUpdate:    tmp.UpdatedAt,
		ReceiveAmount: amountRecv,
		//Confirmations:  not available
		TxID:           tmp.PayoutHash,
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
	case "finished":
		return 1
	case "waiting":
		return 2
	case "confirming":
		return 3
	case "refunded":
		return 5
	case "expired":
		return 7
	case "new":
		return 8
	case "exchanging":
		return 9
	case "sending":
		return 10
	case "failed":
		return 11
	default:
		return 0
	}
}

/* func (c *ChangeNow) CheckOrderStatus(vars interface{}) (res string, err error) {
	err = errors.New("changenow:error: not available for this exchange")
	return
} */
