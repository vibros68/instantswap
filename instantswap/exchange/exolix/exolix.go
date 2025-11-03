package exolix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vibros68/instantswap/instantswap"
)

const (
	API_BASE = "https://exolix.com/api/v2/"
	LIBNAME  = "exolix"
)

type Exolix struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
	instantswap.IDExchange

	cache struct {
		mu              sync.RWMutex
		currencies      []instantswap.Currency
		lastUpdate      time.Time
		effectivePeriod time.Duration
	}
}

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

func (e *Exolix) Name() string {
	return LIBNAME
}

// SetDebug set enable/disable http request/response dump.
func (e *Exolix) SetDebug(enable bool) {
	e.conf.Debug = enable
}

// New return a exolix client.
func New(conf instantswap.ExchangeConfig) (*Exolix, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		r.Header.Set("Authorization", conf.ApiKey)
		return nil
	})
	exolixObj := &Exolix{client: client, conf: &conf}
	// Set the currencies cache validity time to 30 days
	exolixObj.cache.effectivePeriod = 720 * time.Hour
	return exolixObj, nil
}

func (e *Exolix) fetchCurrencies() (currencies []instantswap.Currency, err error) {
	pageSize := "100" // maximum allowable of pagesize
	currentPage := 1
	var exoCurrencies []Currency
	for {
		params := url.Values{}
		params.Add("page", strconv.Itoa(currentPage))
		params.Add("size", pageSize)
		params.Add("withNetworks", "true")
		query := params.Encode()
		// handler get currencies
		r, err := e.client.Do(API_BASE, http.MethodGet,
			fmt.Sprintf("currencies?%s", query), "", false)
		if err != nil {
			return nil, err
		}
		var pageCurrenciesRes CurrencyResponse
		err = parseResponseData(r, &pageCurrenciesRes)
		if err != nil {
			return nil, err
		}
		// if res length is 0, break
		if len(pageCurrenciesRes.Data) == 0 {
			break
		}
		exoCurrencies = append(exoCurrencies, pageCurrenciesRes.Data...)
		// set for next page
		currentPage++
	}
	for _, curr := range exoCurrencies {
		if len(curr.Networks) > 0 {
			for _, net := range curr.Networks {
				currencies = append(currencies, instantswap.Currency{
					Name:    curr.Name,
					Symbol:  strings.ToLower(curr.Code),
					Network: net.Network,
				})
			}
		} else {
			currencies = append(currencies, instantswap.Currency{
				Name:   curr.Name,
				Symbol: strings.ToLower(curr.Code),
			})
		}
	}
	return currencies, nil
}

func (e *Exolix) GetCurrencies() (currencies []instantswap.Currency, err error) {
	e.cache.mu.RLock()
	// if cache is valid and element exists in array
	if time.Since(e.cache.lastUpdate) < e.cache.effectivePeriod && len(e.cache.currencies) > 0 {
		defer e.cache.mu.RUnlock()
		return e.cache.currencies, nil
	}
	e.cache.mu.RUnlock()

	// Fetch currencies
	currencies, err = e.fetchCurrencies()
	if err != nil {
		return nil, err
	}

	// Update cache
	e.cache.mu.Lock()
	e.cache.currencies = currencies
	e.cache.lastUpdate = time.Now()
	e.cache.mu.Unlock()

	return currencies, nil
}

func (e *Exolix) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	// get all currencies
	allCurrencies, err := e.GetCurrencies()
	if err != nil {
		return nil, err
	}
	for _, exoCurr := range allCurrencies {
		if !strings.EqualFold(from, exoCurr.Symbol) {
			currencies = append(currencies, exoCurr)
		}
	}
	return currencies, nil
}

func (e *Exolix) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	var rateResponse RateResponse
	params := url.Values{}
	params.Add("coinFrom", vars.From)
	params.Add("coinTo", vars.To)
	params.Add("amount", fmt.Sprintf("%f", vars.Amount))
	params.Add("rateType", "fixed")
	if vars.FromNetwork != "" {
		params.Add("networkFrom", vars.FromNetwork)
	}
	if vars.ToNetwork != "" {
		params.Add("networkTo", vars.ToNetwork)
	}
	query := params.Encode()
	endpoint := fmt.Sprintf("rate?%s", query)
	r, rerr := e.client.Do(API_BASE, http.MethodGet, endpoint, "", false)
	if rerr != nil {
		err = rerr
		return
	}
	err = parseResponseData(r, &rateResponse)
	if err != nil {
		return
	}
	res = instantswap.ExchangeRateInfo{
		ExchangeRate:    rateResponse.Rate,
		Min:             rateResponse.MinAmount,
		Max:             rateResponse.MaxAmount,
		EstimatedAmount: rateResponse.ToAmount,
	}
	return
}

func (e *Exolix) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return res, fmt.Errorf("not supported")
}

func (e *Exolix) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	var req = OrderRequest{
		CoinFrom:          vars.FromCurrency,
		CoinTo:            vars.ToCurrency,
		NetworkFrom:       vars.FromNetwork,
		NetworkTo:         vars.ToNetwork,
		Amount:            vars.InvoicedAmount,
		WithdrawalAddress: vars.Destination,
		WithdrawalAmount:  vars.OrderedAmount,
		WithdrawalExtraId: vars.ExtraID,
		RefundAddress:     vars.RefundAddress,
		RefundExtraId:     vars.RefundExtraID,
	}
	body, _ := json.Marshal(req)
	r, err := e.client.Do(API_BASE, http.MethodPost, "transactions", string(body), false)
	if err != nil {
		return res, err
	}
	var order Order
	err = parseResponseData(r, &order)
	if err != nil {
		return res, err
	}
	res = instantswap.CreateResultInfo{
		ChargedFee:     0,
		Destination:    order.WithdrawalAddress,
		ExchangeRate:   order.Rate,
		FromCurrency:   order.CoinFrom.CoinCode,
		InvoicedAmount: order.Amount,
		OrderedAmount:  order.AmountTo,
		ToCurrency:     order.CoinTo.CoinCode,
		UUID:           order.Id,
		DepositAddress: order.DepositAddress,
		Expires:        0,
		ExtraID:        "",
		PayoutExtraID:  "",
	}
	return res, nil
}

func (e *Exolix) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (e *Exolix) CancelOrder(orderID string) (res string, err error) {
	return
}

func (e *Exolix) OrderInfo(req instantswap.TrackingRequest) (res instantswap.OrderInfoResult, err error) {
	r, err := e.client.Do(API_BASE, http.MethodGet, fmt.Sprintf("transactions/%s", req.OrderId), "", false)
	if err != nil {
		return res, err
	}
	var order Order
	err = parseResponseData(r, &order)
	if err != nil {
		return res, err
	}
	var txId string
	if order.HashOut.Hash != nil {
		txId = *order.HashOut.Hash
	}
	res = instantswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     "",
		ReceiveAmount:  order.AmountTo,
		TxID:           txId,
		Status:         order.Status,
		InternalStatus: parseStatus(order.Status),
		Confirmations:  "",
	}
	return
}

func parseResponseData(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}

// wait, confirmation, confirmed, exchanging, sending, success, overdue, refunded
func parseStatus(status string) instantswap.Status {
	switch status {
	case "wait":
		return instantswap.OrderStatusWaitingForDeposit
	case "confirmation":
		return instantswap.OrderStatusDepositReceived
	case "confirmed":
		return instantswap.OrderStatusDepositConfirmed
	case "exchanging":
		return instantswap.OrderStatusExchanging
	case "sending":
		return instantswap.OrderStatusSending
	case "success":
		return instantswap.OrderStatusCompleted
	case "overdue":
		return instantswap.OrderStatusFailed
	case "refunded":
		return instantswap.OrderStatusRefunded
	}
	return instantswap.OrderStatusUnknown
}
