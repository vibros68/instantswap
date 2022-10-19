package coinswitch

import (
	"gitlab.com/raedah/bonzai/global/utils"
	"net/http"
	"strings"
	"time"
	//"fmt"
	//"strconv"

	"encoding/json"
	"errors"

	"gitlab.com/raedah/bonzai/global/idexchange"
)

const (
	//NOTE: when changing the address from sandboxapi to live api if the version "v1" changes you will need to change the value in
	// global/clients/exchangeclient/exchangeclient.go switch statements for auth related to "coinswitch" signature "/v1/" because
	//there api requires just the info in the url after the first forward slash /... (temporary solution)
	API_BASE                   = "https://sandboxapi.coinswitch.co/v1/" // API endpoint
	DEFAULT_HTTPCLIENT_TIMEOUT = 30                                     // HTTP client timeout
	LIBNAME                    = "coinswitch"
	waitSec                    = 3
)

func init() {
	idexchange.RegisterExchange(LIBNAME, func(config idexchange.ExchangeConfig) (idexchange.IDExchange, error) {
		return New(config)
	})
}

// New return a instanciate struct
func New(conf idexchange.ExchangeConfig) (*CoinSwitch, error) {
	client := idexchange.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
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
	client *idexchange.Client
	conf   *idexchange.ExchangeConfig
	idexchange.IDExchange
}

//SetDebug set enable/disable http request/response dump
func (c *CoinSwitch) SetDebug(enable bool) {
	c.conf.Debug = enable
}

/* func handleErr(r json.RawMessage) (err error) {
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
} */

//CalculateExchangeRate get estimate on the amount for the exchange
func (c *CoinSwitch) GetExchangeRateInfo(vars idexchange.ExchangeRateRequest) (res idexchange.ExchangeRateInfo, err error) {
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

	res = idexchange.ExchangeRateInfo{
		ExchangeRate:    rate,
		Min:             limits.Min,
		Max:             limits.Max,
		EstimatedAmount: estimate.EstimatedAmount,
		Signature:       estimate.Signature,
	}

	return
}

//EstimateAmount get estimate on the amount for the exchange
func (c *CoinSwitch) EstimateAmount(vars idexchange.ExchangeRateRequest) (res idexchange.EstimateAmount, err error) {
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

	res = idexchange.EstimateAmount{
		EstimatedAmount: tmpRes.DestinationCoinAmount,
		Signature:       tmpRes.OfferReferenceID, //used for referring to quote if needed
	}

	return
}

//QueryRates (list of pairs LTC-BTC, BTC-LTC, etc)
func (c *CoinSwitch) QueryRates(vars interface{}) (res []idexchange.QueryRate, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryActiveCurrencies get all active currencies
func (c *CoinSwitch) QueryActiveCurrencies(vars interface{}) (res []idexchange.ActiveCurr, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryLimits Get Exchange Rates (from, to)
func (c *CoinSwitch) QueryLimits(fromCurr, toCurr string) (res idexchange.QueryLimits, err error) {
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
	res = idexchange.QueryLimits{
		Min: tmp.Min,
		Max: tmp.Max,
	}
	return
}

//CreateOrder create an instant exchange order
func (c *CoinSwitch) CreateOrder(orderInfo idexchange.CreateOrder) (res idexchange.CreateResultInfo, err error) {

	//convert from interface orderInfo to local struct
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

	res = idexchange.CreateResultInfo{
		UUID: tmp.OrderID,
		/* Destination:    tmp.,
		ExtraID:        tmp.PayinExtraID,
		FromCurrency:   tmp.FromCurrency, //all commented out fields are not used for this exchange
		ToCurrency:     tmp.ToCurrency, */
		DepositAddress: tmp.ExchangeAddress.Address,
	}
	return
}

//UpdateOrder not available for this exchange
func (c *CoinSwitch) UpdateOrder(vars interface{}) (res idexchange.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

//CancelOrder not available for this exchange
func (c *CoinSwitch) CancelOrder(oId string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

//OrderInfo get information on orderid/uuid
func (c *CoinSwitch) OrderInfo(oId string) (res idexchange.OrderInfoResult, err error) {
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
	res = idexchange.OrderInfoResult{
		Expires: tmp.ValidTill,
		//LastUpdate:    tmp.UpdatedAt, //not available
		ReceiveAmount: tmp.DestinationCoinAmount,
		//Confirmations:  not available //not available
		TxID:           tmp.OutputTransactionHash,
		Status:         tmp.Status,
		InternalStatus: idexchange.OrderStatus(GetLocalStatus(tmp.Status)),
	}
	return
}

//Possible transaction statuses
//new waiting confirming exchanging sending finished failed refunded expired

//GetLocalStatus translate local status to idexchange status id
func GetLocalStatus(status string) (iStatus int) {
	status = strings.ToLower(status)
	switch status {
	case "finished", "complete":
		return 1
	case "no_deposit":
		return 2
	case "confirming":
		return 3
	case "refunded":
		return 5
	case "timeout":
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

/* func (c *CoinSwitch) CheckOrderStatus(vars interface{}) (res string, err error) {
	err = errors.New("changenow:error: not available for this exchange")
	return
} */
