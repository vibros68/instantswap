package fixedfloat

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vibros68/instantswap/instantswap"
)

const (
	API_BASE = "https://ff.io/api/v2/"
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
		r.Header.Set("X-API-SIGN", signedMsg)
		r.Header.Set("X-API-KEY", conf.ApiKey)
		return nil
	})
	return &FixedFloat{client: client, conf: &conf}, nil
}

func (c *FixedFloat) Name() string {
	return LIBNAME
}

// SetDebug set enable/disable http request/response dump
func (c *FixedFloat) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *FixedFloat) GetCurrencies() (currencies []instantswap.Currency, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, http.MethodPost, "ccies", "", false)
	if err != nil {
		return nil, err
	}
	var ffCurrs []Currency
	err = parseResponseData(r, &ffCurrs)
	currencies = make([]instantswap.Currency, len(ffCurrs))
	for i, ffCurr := range ffCurrs {
		currencies[i] = instantswap.Currency{
			Name:   ffCurr.Name,
			Symbol: ffCurr.Code,
		}
	}
	return currencies, err
}

func (c *FixedFloat) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, http.MethodPost, "ccies", "", false)
	if err != nil {
		return nil, err
	}
	var ffCurrs []Currency
	err = parseResponseData(r, &ffCurrs)
	for _, ffCurr := range ffCurrs {
		if strings.ToLower(from) != strings.ToLower(ffCurr.Code) {
			currencies = append(currencies, instantswap.Currency{
				Name:   ffCurr.Name,
				Symbol: ffCurr.Code,
			})
		}
	}
	return currencies, err
}

func buildBody(data interface{}) string {
	if data == nil {
		return ""
	}
	body, _ := json.Marshal(data)
	return string(body)
}

// GetExchangeRateInfo get estimate on the amount for the exchange.
func (c *FixedFloat) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	f := PriceReq{
		FromCcy:   vars.From,
		ToCcy:     vars.To,
		Amount:    vars.Amount,
		Direction: "from",
		Type:      "fixed",
	}
	var r []byte
	r, err = c.client.Do(API_BASE, http.MethodPost, "price", buildBody(f), false)
	if err != nil {
		return res, err
	}
	var priceRes PriceResult
	err = parseResponseData(r, &priceRes)
	if err != nil {
		return res, err
	}
	return instantswap.ExchangeRateInfo{
		Min:             priceRes.From.Min,
		Max:             priceRes.From.Max,
		ExchangeRate:    priceRes.From.Rate,
		EstimatedAmount: priceRes.To.Amount,
		MaxOrder:        0,
		Signature:       "",
	}, nil
}

func (c *FixedFloat) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	var f = CreateOrderRequest{
		FromCcy:   vars.FromCurrency,
		ToCcy:     vars.ToCurrency,
		Amount:    vars.InvoicedAmount,
		Direction: "from",
		Type:      "fixed",
		ToAddress: vars.Destination,
	}
	var r []byte
	r, err = c.client.Do(API_BASE, http.MethodPost, "create", buildBody(f), false)
	if err != nil {
		return res, err
	}
	var orderRes OrderResponse
	err = parseResponseData(r, &orderRes)
	if err != nil {
		return res, err
	}
	return instantswap.CreateResultInfo{
		ChargedFee:     0,
		Destination:    orderRes.From.Address,
		ExchangeRate:   orderRes.From.Amount / orderRes.To.Amount,
		FromCurrency:   orderRes.From.Code,
		InvoicedAmount: orderRes.From.Amount,
		OrderedAmount:  orderRes.To.Amount,
		ToCurrency:     orderRes.To.Code,
		UUID:           orderRes.Id,
		DepositAddress: orderRes.From.Address,
		Expires:        0,
		ExtraID:        orderRes.Token,
		PayoutExtraID:  "",
	}, nil
}

// UpdateOrder accepts orderID value and more if needed per lib.
func (c *FixedFloat) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (c *FixedFloat) CancelOrder(orderID string) (res string, err error) {
	return
}

// OrderInfo accepts string of orderID value.
func (c *FixedFloat) OrderInfo(req instantswap.TrackingRequest) (res instantswap.OrderInfoResult, err error) {
	if len(req.ExtraId) == 0 {
		return res, fmt.Errorf("fetching fixedfloat order require order token")
	}
	var f = struct {
		Id    string `json:"id"`
		Token string `json:"token"`
	}{
		Id:    req.OrderId,
		Token: req.ExtraId,
	}
	var r []byte
	r, err = c.client.Do(API_BASE, http.MethodPost, "order", buildBody(f), false)
	if err != nil {
		return res, err
	}
	var orderRes OrderResponse
	err = parseResponseData(r, &orderRes)
	if err != nil {
		return res, err
	}
	var txId string
	if orderRes.To.Tx.Id != nil {
		txId, _ = orderRes.To.Tx.Id.(string)
	}
	return instantswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     "",
		ReceiveAmount:  orderRes.To.Amount,
		TxID:           txId,
		Status:         orderRes.Status,
		InternalStatus: GetLocalStatus(orderRes.Status),
		Confirmations:  "",
	}, nil
}

// GetLocalStatus translate local status to instantswap.Status.
func GetLocalStatus(status string) (iStatus instantswap.Status) {
	switch status {
	case "NEW":
		return instantswap.OrderStatusWaitingForDeposit
	case "PENDING":
		return instantswap.OrderStatusDepositReceived
	case "EXCHANGE":
		return instantswap.OrderStatusExchanging
	case "WITHDRAW":
		return instantswap.OrderStatusSending
	case "DONE":
		return instantswap.OrderStatusCompleted
	case "EXPIRED":
		return instantswap.OrderStatusExpired
	case "EMERGENCY":
		return instantswap.OrderStatusFailed
	default:
		return instantswap.OrderStatusUnknown
	}
}
