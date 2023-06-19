package coinswitch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/crypto-power/instantswap/instantswap"
	"github.com/crypto-power/instantswap/instantswap/utils"
)

const (
	API_BASE = "https://sandboxapi.coinswitch.co/v1/" // API endpoint
	LIBNAME  = "coinswitch"
)

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// New return a CoinSwitch client.
func New(conf instantswap.ExchangeConfig) (*CoinSwitch, error) {
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
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

// CoinSwitch represent a CoinSwitch client.
type CoinSwitch struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
	instantswap.IDExchange
}

// SetDebug set enable/disable http request/response dump.
func (c *CoinSwitch) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *CoinSwitch) GetCurrencies() (currencies []instantswap.Currency, err error) {
	r, err := c.client.Do(API_BASE, "GET", "coins", "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var csCurrencies []Currency
	err = parseResponseData(r, &csCurrencies)
	if err != nil {
		return nil, err
	}
	for _, currency := range csCurrencies {
		if currency.IsActive {
			currencies = append(currencies, instantswap.Currency{
				Name:   currency.Name,
				Symbol: currency.Symbol,
			})
		}
	}
	return
}

func (c *CoinSwitch) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	r, err := c.client.Do(API_BASE, "GET", fmt.Sprintf("coins/%s/destination-coins", strings.ToLower(from)), "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var csCurrencies []Currency
	err = parseResponseData(r, &csCurrencies)
	if err != nil {
		return nil, err
	}
	currencies = []instantswap.Currency{}
	for _, currency := range csCurrencies {
		if currency.IsActive {
			currencies = append(currencies, instantswap.Currency{
				Name:   currency.Name,
				Symbol: currency.Symbol,
			})
		}
	}
	return
}

// GetExchangeRateInfo get estimate on the amount for the exchange.
func (c *CoinSwitch) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
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
	res = instantswap.ExchangeRateInfo{
		ExchangeRate:    rate,
		Min:             limits.Min,
		Max:             limits.Max,
		EstimatedAmount: estimate.EstimatedAmount,
		Signature:       estimate.Signature,
	}
	return
}

// EstimateAmount get estimate on the amount for the exchange.
func (c *CoinSwitch) EstimateAmount(vars instantswap.ExchangeRateRequest) (res instantswap.EstimateAmount, err error) {
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

	res = instantswap.EstimateAmount{
		EstimatedAmount: tmpRes.DestinationCoinAmount,
		Signature:       tmpRes.OfferReferenceID, //used for referring to quote if needed
	}

	return
}

// QueryRates (list of pairs LTC-BTC, BTC-LTC, etc).
func (c *CoinSwitch) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

// QueryLimits Get Exchange Rates (from, to).
func (c *CoinSwitch) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
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
	res = instantswap.QueryLimits{
		Min: tmp.Min,
		Max: tmp.Max,
	}
	return
}

// CreateOrder create an instant exchange order.
func (c *CoinSwitch) CreateOrder(orderInfo instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	destAddress := Address{Address: orderInfo.Destination}
	refundAddress := Address{Address: orderInfo.RefundAddress}
	tmpOrderInfo := CreateOrder{
		DepositCoin:        strings.ToLower(orderInfo.FromCurrency),
		DestinationCoin:    strings.ToLower(orderInfo.ToCurrency),
		DestinationAddress: destAddress,
		RefundAddress:      refundAddress,
		DepositCoinAmount:  orderInfo.InvoicedAmount,
		OfferReferenceID:   orderInfo.Signature,
		UserReferenceID:    c.conf.AffiliateId,
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

	res = instantswap.CreateResultInfo{
		UUID:           tmp.OrderID,
		DepositAddress: tmp.ExchangeAddress.Address,
	}
	return
}

// UpdateOrder not available for this exchange.
func (c *CoinSwitch) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

// CancelOrder not available for this exchange.
func (c *CoinSwitch) CancelOrder(oId string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

// OrderInfo get information on orderid/uuid.
func (c *CoinSwitch) OrderInfo(oId string) (res instantswap.OrderInfoResult, err error) {
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

	// Once the deposit is detected on blockchain "validTill" will be NULL which indicates
	// there is no expiry for the order
	res = instantswap.OrderInfoResult{
		Expires:        tmp.ValidTill,
		ReceiveAmount:  tmp.DestinationCoinAmount,
		TxID:           tmp.OutputTransactionHash,
		Status:         tmp.Status,
		InternalStatus: GetLocalStatus(tmp.Status),
	}
	return
}

// GetLocalStatus translate local status to instantswap.Status.
func GetLocalStatus(status string) instantswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "finished", "complete":
		return instantswap.OrderStatusCompleted
	case "no_deposit":
		return instantswap.OrderStatusNew
	case "confirming":
		return instantswap.OrderStatusDepositReceived
	case "refunded":
		return instantswap.OrderStatusRefunded
	case "timeout":
		return instantswap.OrderStatusExpired
	case "new":
		return instantswap.OrderStatusNew
	case "exchanging":
		return instantswap.OrderStatusExchanging
	case "sending":
		return instantswap.OrderStatusSending
	case "failed":
		return instantswap.OrderStatusFailed
	default:
		return instantswap.OrderStatusUnknown
	}
}
