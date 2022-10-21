package flypme

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
	API_BASE = "https://flyp.me/api/v1/" // API endpoint
	LIBNAME  = "flypme"
)

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

// New return a FlypMe struct.
func New(conf lightningswap.ExchangeConfig) (*FlypMe, error) {
	client := lightningswap.NewClient(LIBNAME, &conf)
	return &FlypMe{
		client: client,
		conf:   &conf,
	}, nil
}

// FlypMe represent a flyp.me exchange client.
type FlypMe struct {
	client *lightningswap.Client
	conf   *lightningswap.ExchangeConfig
	lightningswap.IDExchange
}

// SetDebug set enable/disable http request/response dump.
func (c *FlypMe) SetDebug(enable bool) {
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

// GetExchangeRateInfo get estimate on the amount for the exchange.
func (c *FlypMe) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
	limits, err := c.QueryLimits(vars.From, vars.To)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	time.Sleep(time.Second * 1)
	exchangeRates, err := c.QueryRates(nil)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var rate lightningswap.QueryRate
	var pair = fmt.Sprintf("%s-%s", vars.From, vars.To)
	for _, v := range exchangeRates {
		if v.Name == pair {
			rate = v
		}
	}
	if rate.Name == "" || rate.Value == "" {
		err = errors.New(LIBNAME + ":error: rate not found for " + pair + " pair")
		return
	}

	exchangeRate, err := strconv.ParseFloat(rate.Value, 64)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	rateFinal := 1 / exchangeRate
	min := limits.Min * rateFinal * 1.5
	max := limits.Max * rateFinal

	res = lightningswap.ExchangeRateInfo{
		ExchangeRate:    rateFinal,
		Min:             min,
		Max:             max,
		EstimatedAmount: (vars.Amount / rateFinal),
	}

	return
}

// EstimateAmount get estimate on the amount for the exchange.
func (c *FlypMe) EstimateAmount(vars interface{}) (res lightningswap.EstimateAmount, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

// QueryRates (list of pairs LTC-BTC, BTC-LTC, etc).
func (c *FlypMe) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	//vars not used here
	r, err := c.client.Do(API_BASE, "GET", "data/exchange_rates", "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	tmpArr := []lightningswap.QueryRate{}
	var v interface{}
	if err = json.Unmarshal(r, &v); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	data := v.(map[string]interface{})
	for k, v := range data {
		val := (v).(string)
		tmpQ := lightningswap.QueryRate{Name: k, Value: val}
		tmpArr = append(tmpArr, tmpQ)
	}
	res = tmpArr

	return
}

// QueryActiveCurrencies returns Flypme's supported currencies
func (c *FlypMe) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	//vars not used here
	r, err := c.client.Do(API_BASE, "GET", "currencies", "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	tmpArr := []lightningswap.ActiveCurr{}
	var v interface{}
	if err = json.Unmarshal(r, &v); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	data := v.(map[string]interface{})
	for _, v := range data {

		curr := (v).(map[string]interface{})
		currMarsh, err := json.Marshal(curr)
		if err != nil {
			err = errors.New(LIBNAME + ":error: " + err.Error())
			return tmpArr, err
		}

		var activeCurr lightningswap.ActiveCurr
		err = json.Unmarshal(currMarsh, &activeCurr)
		if err != nil {
			err = errors.New(LIBNAME + ":error: " + err.Error())
			return tmpArr, err
		}

		tmpArr = append(tmpArr, activeCurr)
	}
	res = tmpArr
	return
}

// QueryLimits Get Exchange Rates (from, to).
func (c *FlypMe) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {

	r, err := c.client.Do(API_BASE, "GET", "order/limits/"+fromCurr+"/"+toCurr, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmp QueryLimits
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = lightningswap.QueryLimits{
		Max: tmp.Max,
		Min: tmp.Min,
	}
	return
}

func (c *FlypMe) CreateOrder(orderInfo lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
	newOrder := CreateOrder{
		Order: CreateOrderInfo{
			FromCurrency:   orderInfo.FromCurrency,
			ToCurrency:     orderInfo.ToCurrency,
			InvoicedAmount: strconv.FormatFloat(orderInfo.InvoicedAmount, 'f', 8, 64), //amount in "from" currency
			OrderedAmount:  "",                                                        //amount in "to" currency
			Destination:    orderInfo.Destination,
			RefundAddress:  orderInfo.RefundAddress,
		},
	}
	payload, err := json.Marshal(newOrder)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	r, err := c.client.Do(API_BASE, "POST", "order/new", string(payload), false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmp CreateResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	if len(tmp.Errors) > 0 {
		err = handleErr(tmp.Errors)
		if err != nil {
			err = errors.New(LIBNAME + ":error: " + err.Error())
			return
		}
	}
	acceptOrder := UUID{
		UUID: tmp.Order.UUID,
	}
	acceptPayload, err := json.Marshal(acceptOrder)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	acceptRes, err := c.client.Do(API_BASE, "POST", "order/accept", string(acceptPayload), false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmpAccept AcceptOrderResult
	if err = json.Unmarshal(acceptRes, &tmpAccept); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if len(tmp.Errors) > 0 {
		err = handleErr(tmp.Errors)
		if err != nil {
			err = errors.New(LIBNAME + ":error: " + err.Error())
			return
		}
	}

	res = lightningswap.CreateResultInfo{
		ChargedFee:     tmp.Order.ChargedFee,
		Destination:    tmp.Order.Destination,
		ExchangeRate:   tmp.Order.ExchangeRate,
		FromCurrency:   tmp.Order.FromCurrency,
		InvoicedAmount: tmp.Order.InvoicedAmount,
		OrderedAmount:  tmp.Order.OrderedAmount,
		ToCurrency:     tmp.Order.ToCurrency,
		UUID:           tmp.Order.UUID,
		DepositAddress: tmpAccept.DepositAddress, //from accept order result
		Expires:        tmpAccept.Expires,        //from accept order result
	}

	return
}

// UpdateOrder update the information of an order.
func (c *FlypMe) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	orderInfo := vars.(UpdateOrder)
	payload, err := json.Marshal(orderInfo)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	r, err := c.client.Do(API_BASE, "POST", "order/update", string(payload), false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmp UpdateOrderResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if len(tmp.Errors) > 0 {
		err = handleErr(tmp.Errors)
		if err != nil {
			err = errors.New(LIBNAME + ":error: " + err.Error())
			return
		}
	}

	res = lightningswap.UpdateOrderResultInfo{
		ChargedFee:     tmp.Order.ChargedFee,
		Destination:    tmp.Order.Destination,
		ExchangeRate:   tmp.Order.ExchangeRate,
		FromCurrency:   tmp.Order.FromCurrency,
		InvoicedAmount: tmp.Order.InvoicedAmount,
		OrderedAmount:  tmp.Order.OrderedAmount,
		ToCurrency:     tmp.Order.ToCurrency,
		UUID:           tmp.Order.UUID,
	}

	return
}

// CancelOrder will delete an order based on its id.
func (c *FlypMe) CancelOrder(orderId string) (res string, err error) {
	cancelOrder := UUID{
		UUID: orderId,
	}
	payload, err := json.Marshal(cancelOrder)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	r, err := c.client.Do(API_BASE, "POST", "order/cancel", string(payload), false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var result jsonResponse
	if err = json.Unmarshal(r, &result); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	if len(result.Errors) > 0 {
		err = handleErr(result.Errors)
		if err != nil {
			err = errors.New(LIBNAME + ":error: " + err.Error())
			return
		}
	}
	res = result.Result

	return
}

// OrderInfo accepts string of orderID value and return
// its information
func (c *FlypMe) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
	getOrderInfo := UUID{
		UUID: orderID,
	}
	payload, err := json.Marshal(getOrderInfo)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	r, err := c.client.Do(API_BASE, "POST", "order/info", string(payload), false)
	if err != nil {
		return
	}
	var tmp OrderInfoResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	if len(tmp.Errors) > 0 {
		err = handleErr(tmp.Errors)
		if err != nil {
			err = errors.New(LIBNAME + ":error: " + err.Error())
			return
		}
	}
	res = lightningswap.OrderInfoResult{
		//LastUpdate:   not available
		Expires:        tmp.Expires,
		ReceiveAmount:  tmp.Order.OrderedAmount,
		Confirmations:  tmp.Confirmations,
		TxID:           tmp.TxID,
		Status:         tmp.Status,
		InternalStatus: GetLocalStatus(tmp.Status),
	}
	// flypme will return a pending txID like: pending_b1fdc5a8-e470-63c1-a034-eddf78c8fdf6
	// while status still be completed. In this case we will return pending status in our system
	if strings.Index(tmp.TxID, "_") != -1 {
		res.TxID = ""
		res.InternalStatus = lightningswap.OrderStatusExchanging
	}
	return
}

// GetLocalStatus translate local status to lightningswap.Status
// Possible statuses are: WAITING_FOR_DEPOSIT, DEPOSIT_RECEIVED, DEPOSIT_CONFIRMED, EXECUTED, REFUNDED, CANCELED and EXPIRED
func GetLocalStatus(status string) lightningswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "executed":
		return lightningswap.OrderStatusCompleted
	case "waiting_for_deposit":
		return lightningswap.OrderStatusNew
	case "deposit_received":
		return lightningswap.OrderStatusDepositReceived
	case "deposit_confirmed":
		return lightningswap.OrderStatusDepositConfirmed
	case "refunded":
		return lightningswap.OrderStatusRefunded
	case "canceled":
		return lightningswap.OrderStatusCanceled
	case "expired":
		return lightningswap.OrderStatusExpired
	default:
		return lightningswap.OrderStatusUnknown
	}
}
