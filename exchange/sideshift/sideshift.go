package sideshift

import (
	"code.cryptopower.dev/group/instantswap"
	"code.cryptopower.dev/group/instantswap/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	API_BASE = "https://sideshift.ai/api/v2/"
	LIBNAME  = "sideshift"
)

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

type SideShift struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
	instantswap.IDExchange
}

func New(conf instantswap.ExchangeConfig) (*SideShift, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: api key is blank, it is account id on sideshift", LIBNAME)
	}
	if conf.ApiSecret == "" {
		return nil, fmt.Errorf("%s:error: api secret is blank", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		if r.Method == http.MethodPost {
			ipAddress, err := utils.GetPublicIP()
			if err != nil {
				return err
			}
			r.Header.Add("x-user-ip", ipAddress)
			r.Header.Add("x-sideshift-secret", conf.ApiSecret)
		}
		return nil
	})
	return &SideShift{client: client, conf: &conf}, nil
}

func (s *SideShift) GetCurrencies() (currencies []instantswap.Currency, err error) {
	r, err := s.client.Do(API_BASE, "GET", "coins", "", false)
	if err != nil {
		err = fmt.Errorf("%s:error:%v", LIBNAME, err)
		return
	}
	var csCurrencies []Currency
	err = parseResponseData(r, &csCurrencies)
	if err != nil {
		return nil, err
	}
	currencies = make([]instantswap.Currency, len(csCurrencies))
	for i, currency := range csCurrencies {
		currencies[i] = instantswap.Currency{
			Name:     currency.Name,
			Symbol:   currency.Coin,
			Networks: currency.Networks,
		}
	}
	return
}

func (s *SideShift) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	r, err := s.client.Do(API_BASE, "GET", "coins", "", false)
	if err != nil {
		err = fmt.Errorf("%s:error:%v", LIBNAME, err)
		return
	}
	var csCurrencies []Currency
	err = parseResponseData(r, &csCurrencies)
	if err != nil {
		return nil, err
	}
	from = strings.ToUpper(from)
	for _, currency := range csCurrencies {
		if currency.Coin != from {
			currencies = append(currencies, instantswap.Currency{
				Name:     currency.Name,
				Symbol:   currency.Coin,
				Networks: currency.Networks,
			})
		}
	}
	return
}

func (s *SideShift) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	return
}

func (s *SideShift) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return
}

func (s *SideShift) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	req := createFixedShift{
		SettleAddress: vars.Destination,
		AffiliateId:   s.conf.ApiKey,
		QuoteId:       vars.Signature,
		RefundAddress: vars.RefundAddress,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return res, err
	}
	r, err := s.client.Do(API_BASE, http.MethodPost, "shifts/fixed", string(body), false)
	if err != nil {
		return res, err
	}
	var shift FixedShift
	err = parseResponseData(r, &shift)
	if err != nil {
		return res, err
	}
	return instantswap.CreateResultInfo{
		ChargedFee:     0,
		Destination:    shift.SettleAddress,
		ExchangeRate:   utils.StrToFloat(shift.Rate),
		FromCurrency:   shift.DepositCoin,
		InvoicedAmount: utils.StrToFloat(shift.DepositAmount),
		OrderedAmount:  utils.StrToFloat(shift.SettleAmount),
		ToCurrency:     shift.SettleCoin,
		UUID:           shift.Id,
		DepositAddress: shift.DepositAddress,
		Expires:        int(shift.ExpiresAt.Unix()),
		ExtraID:        "",
		PayoutExtraID:  "",
	}, nil
}

func (s *SideShift) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return res, fmt.Errorf("not supported")
}

func (s *SideShift) CancelOrder(orderID string) (res string, err error) {
	return res, fmt.Errorf("not supported")
}

func (s *SideShift) OrderInfo(orderID string) (res instantswap.OrderInfoResult, err error) {
	r, err := s.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("shifts/%s", orderID), "", false)
	if err != nil {
		return res, err
	}
	var shift FixedShift
	err = parseResponseData(r, &shift)
	if err != nil {
		return res, err
	}
	return instantswap.OrderInfoResult{
		Expires:        int(shift.ExpiresAt.Unix()),
		LastUpdate:     "",
		ReceiveAmount:  0,
		TxID:           shift.SettleHash,
		Status:         "",
		InternalStatus: GetLocalStatus(shift.Status),
		Confirmations:  "",
	}, nil
}

func (s *SideShift) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	var req = ExchangeRateRequest{
		DepositCoin:    strings.ToLower(vars.From),
		DepositNetwork: vars.FromNetwork,
		SettleCoin:     strings.ToLower(vars.To),
		SettleNetwork:  vars.ToNetwork,
		DepositAmount:  fmt.Sprintf("%f", vars.Amount),
		SettleAmount:   "",
		AffiliateId:    s.conf.ApiKey,
		CommissionRate: "0",
	}
	body, err := json.Marshal(req)
	if err != nil {
		return res, err
	}
	r, err := s.client.Do(API_BASE, http.MethodPost, "quotes", string(body), false)
	if err != nil {
		err = fmt.Errorf("%s:error:%v", LIBNAME, err)
		return
	}
	var quote Quote
	err = parseResponseData(r, &quote)
	if err != nil {
		return res, err
	}
	pair, _ := s.pair(vars)
	return instantswap.ExchangeRateInfo{
		Min:             utils.StrToFloat(pair.Min),
		Max:             utils.StrToFloat(pair.Max),
		ExchangeRate:    utils.StrToFloat(quote.Rate),
		EstimatedAmount: utils.StrToFloat(quote.SettleAmount),
		MaxOrder:        0,
		Signature:       quote.Id,
	}, nil
}

func (s *SideShift) pair(vars instantswap.ExchangeRateRequest) (pair PairResponse, err error) {
	r, err := s.client.Do(API_BASE, http.MethodGet,
		fmt.Sprintf("pair/%s-%s/%s-%s",
			strings.ToLower(vars.From), vars.FromNetwork,
			strings.ToLower(vars.To), vars.ToNetwork), "", false)
	if err != nil {
		return pair, err
	}
	err = parseResponseData(r, &pair)
	return pair, err
}

// GetLocalStatus translate local status to instantswap.Status.
func GetLocalStatus(status string) instantswap.Status {
	status = strings.ToLower(status)
	switch status {
	case "waiting":
		return instantswap.OrderStatusWaitingForDeposit
	case "pending":
		return instantswap.OrderStatusDepositReceived
	case "processing":
		return instantswap.OrderStatusDepositConfirmed
	case "review":
		return instantswap.OrderStatusUnknown
	case "settling":
		return instantswap.OrderStatusExchanging
	case "settled":
		return instantswap.OrderStatusCompleted
	case "refund":
		return instantswap.OrderStatusFailed
	case "refunding", "refunded":
		return instantswap.OrderStatusRefunded
	case "multiple":
		return instantswap.OrderStatusUnknown
	default:
		return instantswap.OrderStatusUnknown
	}
}
