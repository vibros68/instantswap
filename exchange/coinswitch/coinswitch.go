package coinswitch

import (
	"code.cryptopower.dev/exchange/lightningswap"
	"code.cryptopower.dev/exchange/lightningswap/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	// NOTE: when changing the address from sandboxapi to live api if the version "v1" changes you will need to change the value in
	// global/clients/exchangeclient/exchangeclient.go switch statements for auth related to "coinswitch" signature "/v1/" because
	// there api requires just the info in the url after the first forward slash /... (temporary solution)
	API_BASE = "https://sandboxapi.coinswitch.co/v1/" // API endpoint
	LIBNAME  = "coinswitch"
)

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

// New return a instanciate struct
func New(conf lightningswap.ExchangeConfig) (*CoinSwitch, error) {
	client := lightningswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		ipAddress, err := utils.GetPublicIP()
		if err != nil {
			return err
		}
		r.Header.Add("x-user-ip", ipAddress)
		r.Header.Add("x-api-key", conf.ApiKey)
		return nil
	})
	return &CoinSwitch{
		client: client,
		conf:   &conf,
	}, nil
}

//CoinSwitch represent a CoinSwitch client
type CoinSwitch struct {
	client *lightningswap.Client
	conf   *lightningswap.ExchangeConfig
	lightningswap.IDExchange
}

//SetDebug set enable/disable http request/response dump
func (c *CoinSwitch) SetDebug(enable bool) {
	c.conf.Debug = enable
}

//CalculateExchangeRate get estimate on the amount for the exchange
func (c *CoinSwitch) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
	limits, err := c.QueryLimits(vars.From, vars.To)
	if err != nil {
		return
	}
	time.Sleep(time.Second * 1)
	estimate, err := c.EstimateAmount(vars)
	if err != nil {
		return
	}
	rate := vars.Amount / estimate.EstimatedAmount
	res = lightningswap.ExchangeRateInfo{
		ExchangeRate:    rate,
		Min:             limits.Min,
		Max:             limits.Max,
		EstimatedAmount: estimate.EstimatedAmount,
		Signature:       estimate.Signature,
	}
	return
}

//EstimateAmount get estimate on the amount for the exchange
func (c *CoinSwitch) EstimateAmount(vars lightningswap.ExchangeRateRequest) (res lightningswap.EstimateAmount, err error) {
	estimateReq := EstimateRequest{
		DepositCoin:     strings.ToLower(vars.From),
		DestinationCoin: strings.ToLower(vars.To),
		DepositAmount:   vars.Amount,
	}

	payload, err := json.Marshal(estimateReq)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	r, err := c.client.Do(API_BASE, "POST", "offer", string(payload), false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if !response.Success {
		err = errors.New(LIBNAME + ":error:" + response.Code + ": " + response.Message)
		return
	}
	var tmpRes EstimateAmount
	if err = json.Unmarshal(response.Result, &tmpRes); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = lightningswap.EstimateAmount{
		EstimatedAmount: tmpRes.DestinationCoinAmount,
		Signature:       tmpRes.OfferReferenceID, //used for referring to quote if needed
	}

	return
}

//QueryRates (list of pairs LTC-BTC, BTC-LTC, etc)
func (c *CoinSwitch) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryActiveCurrencies get all active currencies
func (c *CoinSwitch) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryLimits Get Exchange Rates (from, to)
func (c *CoinSwitch) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {
	limitReq := QueryLimitsRequest{
		DepositCoin:     strings.ToLower(fromCurr),
		DestinationCoin: strings.ToLower(toCurr),
	}

	payload, err := json.Marshal(limitReq)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	r, err := c.client.Do(API_BASE, "POST", "limit", string(payload), false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if !response.Success {
		err = errors.New(LIBNAME + ":error:" + response.Code + ": " + response.Message)
		return
	}
	var tmp QueryLimitsResponse
	if err = json.Unmarshal(response.Result, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	//minAmount := strconv.FormatFloat(tmp.Min, 'f', 8, 64)
	res = lightningswap.QueryLimits{
		Min: tmp.Min,
		Max: tmp.Max,
	}
	return
}

//CreateOrder create an instant exchange order
func (c *CoinSwitch) CreateOrder(orderInfo lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
	destAddress := Address{Address: orderInfo.Destination}
	refundAddress := Address{Address: orderInfo.RefundAddress}
	tmpOrderInfo := CreateOrder{
		DepositCoin:        strings.ToLower(orderInfo.FromCurrency),
		DestinationCoin:    strings.ToLower(orderInfo.ToCurrency),
		DestinationAddress: destAddress,
		RefundAddress:      refundAddress,
		DepositCoinAmount:  orderInfo.InvoicedAmount,
		OfferReferenceID:   orderInfo.Signature,
		UserReferenceID:    orderInfo.UserReferenceID,
	}
	if tmpOrderInfo.DepositCoinAmount == 0.0 {
		err = errors.New(LIBNAME + ":error:createorder deposit coin amount is 0")
		return
	}
	payload, err := json.Marshal(tmpOrderInfo)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	r, err := c.client.Do(API_BASE, "POST", "order", string(payload), true)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if !response.Success {
		err = errors.New(LIBNAME + ":error:" + response.Code + ": " + response.Message)
		return
	}

	var tmp CreateResult
	if err = json.Unmarshal(response.Result, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = lightningswap.CreateResultInfo{
		UUID:           tmp.OrderID,
		DepositAddress: tmp.ExchangeAddress.Address,
	}
	return
}

//UpdateOrder not available for this exchange
func (c *CoinSwitch) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

//CancelOrder not available for this exchange
func (c *CoinSwitch) CancelOrder(oId string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

//OrderInfo get information on orderid/uuid
func (c *CoinSwitch) OrderInfo(oId string) (res lightningswap.OrderInfoResult, err error) {
	r, err := c.client.Do(API_BASE, "GET", "order/"+oId, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var response jsonResponse
	if err = json.Unmarshal(r, &response); err != nil {
		return
	}
	if !response.Success {
		err = errors.New(LIBNAME + ":error:" + response.Code + ": " + response.Message)
		return
	}

	var tmp OrderInfoResult
	if err = json.Unmarshal(response.Result, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	// Once the deposit is detected on blockchain "validTill" will be NULL which indicates there is no expiry for the order
	res = lightningswap.OrderInfoResult{
		Expires:        tmp.ValidTill,
		ReceiveAmount:  tmp.DestinationCoinAmount,
		TxID:           tmp.OutputTransactionHash,
		Status:         tmp.Status,
		InternalStatus: GetLocalStatus(tmp.Status),
	}
	return
}

//GetLocalStatus translate local status to idexchange status id
func GetLocalStatus(status string) lightningswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "finished", "complete":
		return lightningswap.OrderStatusCompleted
	case "no_deposit":
		return lightningswap.OrderStatusNew
	case "confirming":
		return lightningswap.OrderStatusDepositReceived
	case "refunded":
		return lightningswap.OrderStatusRefunded
	case "timeout":
		return lightningswap.OrderStatusExpired
	case "new":
		return lightningswap.OrderStatusNew
	case "exchanging":
		return lightningswap.OrderStatusExchanging
	case "sending":
		return lightningswap.OrderStatusSending
	case "failed":
		return lightningswap.OrderStatusFailed
	default:
		return lightningswap.OrderStatusUnknown
	}
}
