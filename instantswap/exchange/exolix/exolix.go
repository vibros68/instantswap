package exolix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
}

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

func (s *Exolix) Name() string {
	return LIBNAME
}

// SetDebug set enable/disable http request/response dump.
func (s *Exolix) SetDebug(enable bool) {
	s.conf.Debug = enable
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
	return &Exolix{client: client, conf: &conf}, nil
}

func (s *Exolix) GetCurrencies() (currencies []instantswap.Currency, err error) {
	getCurrenciesPath := `currencies?page=%d&size=%d&withNetworks=true`
	var exoCurrencies []Currency
	pageSize := 100 // maximum allowable of pagesize
	currentPage := 1
	for {
		// handler get currencies
		r, err := s.client.Do(API_BASE, http.MethodGet,
			fmt.Sprintf(getCurrenciesPath, currentPage, pageSize), "", false)
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
	currencies = make([]instantswap.Currency, len(exoCurrencies))
	for i, curr := range exoCurrencies {
		currencies[i] = instantswap.Currency{
			Name:   curr.Name,
			Symbol: strings.ToLower(curr.Code),
		}
		// set networks
		currencies[i].Networks = make([]string, len(curr.Networks))
		for _, net := range curr.Networks {
			currencies[i].Networks = append(currencies[i].Networks, net.Network)
		}
	}
	return currencies, nil
}

func (s *Exolix) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	// get all currencies
	allCurrencies, err := s.GetCurrencies()
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

func (s *Exolix) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	var rateResponse RateResponse
	fromNetworkParam := ""
	toNetworkParam := ""
	if vars.FromNetwork != "" {
		fromNetworkParam = fmt.Sprintf("&networkFrom=%s", vars.FromNetwork)
	}
	if vars.ToNetwork != "" {
		toNetworkParam = fmt.Sprintf("&networkTo=%s", vars.ToNetwork)
	}

	r, rerr := s.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("rate?coinFrom=%s&coinTo=%s%s%s&amount=%f&rateType=fixed", vars.From, vars.To, fromNetworkParam, toNetworkParam, vars.Amount), "", false)
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

func (s *Exolix) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return res, fmt.Errorf("not supported")
}

func (s *Exolix) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
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
	r, err := s.client.Do(API_BASE, http.MethodPost, "transactions", string(body), false)
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

func (s *Exolix) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (s *Exolix) CancelOrder(orderID string) (res string, err error) {
	return
}

func (s *Exolix) OrderInfo(req instantswap.TrackingRequest) (res instantswap.OrderInfoResult, err error) {
	r, err := s.client.Do(API_BASE, http.MethodGet, fmt.Sprintf("transactions/%s", req.OrderId), "", false)
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
