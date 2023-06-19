package fixedfloat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/crypto-power/instantswap/instantswap"
)

const (
	API_BASE = "https://fixedfloat.com/api/v1/"
	LIBNAME  = "fixedfloat"
)

// The work on fixedfloat is pending
func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// FixedFloat represent a FixedFloat client.
type FixedFloat struct {
	conf   *instantswap.ExchangeConfig
	client *instantswap.Client
	instantswap.IDExchange
}

// New return FixedFloat client.
func New(conf instantswap.ExchangeConfig) (*FixedFloat, error) {
	if conf.ApiKey == "" || conf.ApiSecret == "" {
		return nil, fmt.Errorf("%s:error: api key and api secret must be provided", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		key := []byte(conf.ApiSecret)
		sig := hmac.New(sha256.New, key)
		sig.Write([]byte(body))
		signedMsg := hex.EncodeToString(sig.Sum(nil))
		fmt.Println(body, signedMsg, conf.ApiKey, conf.ApiSecret)
		r.Header.Set("X-API-SIGN", signedMsg)
		r.Header.Set("X-API-KEY", conf.ApiKey)
		return nil
	})
	return &FixedFloat{client: client, conf: &conf}, nil
}

// SetDebug set enable/disable http request/response dump
func (c *FixedFloat) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *FixedFloat) GetCurrencies() (currencies []instantswap.Currency, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, http.MethodGet, "getCurrencies", "", false)
	if err != nil {
		return nil, err
	}
	var ffCurrs []Currency
	err = parseResponseData(r, &ffCurrs)
	currencies = make([]instantswap.Currency, len(ffCurrs))
	for i, ffCurr := range ffCurrs {
		currencies[i] = instantswap.Currency{
			Name:   ffCurr.Currency,
			Symbol: ffCurr.Symbol,
		}
	}
	return currencies, err
}

func (c *FixedFloat) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	return nil, err
}

// GetExchangeRateInfo get estimate on the amount for the exchange.
func (c *FixedFloat) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	var form = make(url.Values)
	form.Set("fromCurrency", vars.From)
	form.Set("toCurrency", vars.To)
	form.Set("type", "fixed")
	form.Set("fromQty", fmt.Sprintf("%.8f", vars.Amount))
	var r []byte
	r, err = c.client.Do(API_BASE, "POST", "getPrice", form.Encode(), false)
	if err != nil {
		return res, err
	}
	fmt.Println(string(r), err)
	return
}

// EstimateAmount get estimate on the amount for the exchange.
func (c *FixedFloat) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	return res, fmt.Errorf("not supported")
}
func (c *FixedFloat) QueryActiveCurrencies(vars interface{}) (res []instantswap.ActiveCurr, err error) {
	return
}
func (c *FixedFloat) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return
}
func (c *FixedFloat) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	return
}

// UpdateOrder accepts orderID value and more if needed per lib.
func (c *FixedFloat) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (c *FixedFloat) CancelOrder(orderID string) (res string, err error) {
	return
}

// OrderInfo accepts string of orderID value.
func (c *FixedFloat) OrderInfo(orderID string) (res instantswap.OrderInfoResult, err error) {
	return
}

// EstimateAmount get estimate on the amount for the exchange.
func (c *FixedFloat) EstimateAmount(vars interface{}) (res instantswap.EstimateAmount, err error) {
	return
}

// GetLocalStatus translate local status to instantswap.Status.
func GetLocalStatus(status string) (iStatus int) {
	status = strings.ToLower(status)
	switch status {
	case "wait":
		return 2
	case "confirmation":
		return 3
	case "confirmed":
		return 4
	case "exchanging":
		return 9
	case "sending", "sending_confirmation":
		return 10
	case "success":
		return 1
	case "overdue":
		return 7
	case "error":
		return 11
	case "refunded":
		return 5
	default:
		return 0
	}
}
