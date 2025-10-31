package changenow

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vibros68/instantswap/instantswap"
)

const (
	API_BASE = "https://api.changenow.io/v1/" // API endpoint
	LIBNAME  = "changenow"
)

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// New return an ChangeNow client struct with IDExchange implement.
func New(conf instantswap.ExchangeConfig) (*ChangeNow, error) {
	if conf.ApiKey == "" {
		err := fmt.Errorf("APIKEY is blank")
		return nil, err
	}
	client := instantswap.NewClient(LIBNAME, &conf)
	return &ChangeNow{client: client, conf: &conf}, nil
}

// ChangeNow represent a ChangeNow client.
type ChangeNow struct {
	conf   *instantswap.ExchangeConfig
	client *instantswap.Client
	instantswap.IDExchange
}

func (c *ChangeNow) Name() string {
	return LIBNAME
}

// SetDebug set enable/disable http request/response dump.
func (c *ChangeNow) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *ChangeNow) GetCurrencies() (currencies []instantswap.Currency, err error) {
	r, err := c.client.Do(API_BASE, "GET", "currencies?active=true", "", false)
	if err != nil {
		return nil, err
	}
	var cnCurrencies []Currency
	err = json.Unmarshal(r, &cnCurrencies)
	if err != nil {
		return nil, err
	}
	currencies = make([]instantswap.Currency, len(cnCurrencies))
	for i, currency := range cnCurrencies {
		currencies[i] = instantswap.Currency{
			Name:     currency.Name,
			Symbol:   currency.Ticker,
			IsFiat:   currency.IsFiat,
			IsStable: currency.IsStable,
		}
	}
	return currencies, nil
}

func (c *ChangeNow) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	r, err := c.client.Do(API_BASE, "GET",
		fmt.Sprintf("currencies-to/%s", strings.ToLower(from)), "", false)
	if err != nil {
		return nil, err
	}
	var cnCurrencies []Currency
	err = json.Unmarshal(r, &cnCurrencies)
	if err != nil {
		return nil, err
	}
	currencies = make([]instantswap.Currency, len(cnCurrencies))
	for i, currency := range cnCurrencies {
		currencies[i] = instantswap.Currency{
			Name:     currency.Name,
			Symbol:   currency.Ticker,
			IsFiat:   currency.IsFiat,
			IsStable: currency.IsStable,
		}
	}
	return currencies, nil
}

// GetExchangeRateInfo get estimate on the amount for the exchange.
func (c *ChangeNow) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
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
	rate := estimate.EstimatedAmount / vars.Amount

	res = instantswap.ExchangeRateInfo{
		ExchangeRate:    rate,
		Min:             limits.Min,
		Max:             limits.Max,
		EstimatedAmount: estimate.EstimatedAmount,
	}

	return
}

// EstimateAmount get estimate on the amount for the exchange.
func (c *ChangeNow) EstimateAmount(vars instantswap.ExchangeRateRequest) (res instantswap.EstimateAmount, err error) {
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

	res = instantswap.EstimateAmount{
		EstimatedAmount:          tmpRes.EstimatedAmount,
		NetworkFee:               tmpRes.NetworkFee,
		ServiceCommission:        tmpRes.ServiceCommission,
		TransactionSpeedForecast: tmpRes.TransactionSpeedForecast,
		WarningMessage:           tmpRes.WarningMessage,
	}

	return
}

// QueryRates (list of pairs LTC-BTC, BTC-LTC, etc).
func (c *ChangeNow) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

// QueryActiveCurrencies get all active currencies.
func (c *ChangeNow) QueryActiveCurrencies(vars interface{}) (res []instantswap.ActiveCurr, err error) {
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
		tmpItem := instantswap.ActiveCurr{
			CurrencyType: currType,
			Name:         v.Ticker,
		}
		res = append(res, tmpItem)
	}
	return
}

// QueryLimits Get Exchange Rates (from, to).
func (c *ChangeNow) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	r, err := c.client.Do(API_BASE, "GET", "exchange-range/"+fromCurr+"_"+toCurr, "", false)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var tmp QueryLimits
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	res = instantswap.QueryLimits{
		Max: tmp.Max,
		Min: tmp.Min,
	}
	return
}

// CreateOrder create an instant exchange order.
func (c *ChangeNow) CreateOrder(orderInfo instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	tmpOrderInfo := CreateOrder{
		FromCurrency:      orderInfo.FromCurrency,
		ToCurrency:        orderInfo.ToCurrency,
		ToCurrencyAddress: orderInfo.Destination,
		RefundAddress:     orderInfo.RefundAddress,
		InvoicedAmount:    strconv.FormatFloat(orderInfo.InvoicedAmount, 'f', 8, 64),
		ExtraID:           orderInfo.ExtraID,
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

	var tmp CreateResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}

	res = instantswap.CreateResultInfo{
		UUID:           tmp.UUID,
		Destination:    tmp.DestinationAddress,
		ExtraID:        tmp.PayinExtraID,
		FromCurrency:   tmp.FromCurrency,
		InvoicedAmount: orderInfo.InvoicedAmount, // amount you send
		OrderedAmount:  tmp.InvoicedAmount,       // amount you get
		ToCurrency:     tmp.ToCurrency,
		DepositAddress: tmp.DepositAddress,
	}
	return
}

// UpdateOrder not available for this exchange.
func (c *ChangeNow) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

// CancelOrder not available for this exchange.
func (c *ChangeNow) CancelOrder(oId string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

// OrderInfo get information on orderid/uuid.
func (c *ChangeNow) OrderInfo(orderID string, extraIds ...string) (res instantswap.OrderInfoResult, err error) {
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
	var hash = tmp.PayoutHash
	if hash == "Internal transfer " {
		hash = instantswap.TX_HASH_INTERNAL_TRANSFER
	}

	res = instantswap.OrderInfoResult{
		LastUpdate:     tmp.UpdatedAt,
		ReceiveAmount:  amountRecv,
		TxID:           hash,
		Status:         tmp.Status,
		InternalStatus: GetLocalStatus(tmp.Status),
	}
	return
}

// GetLocalStatus translate local status to idexchange status id.
// Possible transaction statuses:
// new waiting confirming exchanging sending finished failed refunded expired
func GetLocalStatus(status string) instantswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "finished":
		return instantswap.OrderStatusCompleted
	case "waiting":
		return instantswap.OrderStatusWaitingForDeposit
	case "confirming":
		return instantswap.OrderStatusDepositReceived
	case "refunded":
		return instantswap.OrderStatusRefunded
	case "expired":
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
