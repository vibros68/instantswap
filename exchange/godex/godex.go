package godex

import (
	"code.cryptopower.dev/exchange/instantswap"
	"code.cryptopower.dev/exchange/instantswap/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	API_BASE = "https://api.godex.io/api/v1/"
	LIBNAME  = "godex"
)

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

type GoDEX struct {
	conf   *instantswap.ExchangeConfig
	client *instantswap.Client
	instantswap.IDExchange
}

func New(conf instantswap.ExchangeConfig) (*GoDEX, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := instantswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			r.Header.Set("public-key", conf.ApiKey)
		}
		return nil
	})
	return &GoDEX{client: client, conf: &conf}, nil
}

//SetDebug set enable/disable http request/response dump
func (c *GoDEX) SetDebug(enable bool) {
	c.conf.Debug = enable
}

func (c *GoDEX) queryRate(req InfoRequest) ([]byte, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return c.client.Do(API_BASE, "POST", "info", string(body), false)
}

func (c *GoDEX) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	var req = InfoRequest{
		From:   strings.ToUpper(vars.From),
		To:     strings.ToUpper(vars.To),
		Amount: vars.Amount,
	}
	r, err := c.queryRate(req)
	if err != nil {
		return res, err
	}
	var info InfoResponse
	err = parseResponseData(r, &info)
	if err != nil {
		return res, err
	}
	var minAmount, _ = info.MinAmount.Float64()
	var estimatedAmount, _ = info.Amount.Float64()
	if minAmount > vars.Amount {
		req.Amount = minAmount
		time.Sleep(time.Second)
		r, err := c.queryRate(req)
		if err != nil {
			return res, err
		}
		err = parseResponseData(r, &info)
		if err != nil {
			return res, err
		}
		estimatedAmount = 0
	}
	return instantswap.ExchangeRateInfo{
		Min:             utils.StrToFloat(info.MinAmount.String()),
		Max:             utils.StrToFloat(info.MaxAmount.String()),
		ExchangeRate:    1 / utils.StrToFloat(info.Rate.String()),
		EstimatedAmount: estimatedAmount,
		MaxOrder:        0,
		Signature:       "",
	}, err
}

func (c *GoDEX) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	return res, fmt.Errorf("not supported")
}

func (c *GoDEX) QueryActiveCurrencies(vars interface{}) (res []instantswap.ActiveCurr, err error) {
	return
}

func (c *GoDEX) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return
}

func (c *GoDEX) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	var txReq = TransactionReq{
		CoinFrom:          vars.FromCurrency,
		CoinTo:            vars.ToCurrency,
		DepositAmount:     vars.InvoicedAmount,
		Withdrawal:        vars.Destination,
		WithdrawalExtraId: "",
		Return:            vars.RefundAddress,
		ReturnExtraId:     vars.RefundExtraID,
		AffiliateId:       c.conf.AffiliateId,
		CoinToNetwork:     "",
		CoinFromNetwork:   "",
	}
	body, err := json.Marshal(txReq)
	if err != nil {
		return res, err
	}
	var r []byte
	r, err = c.client.Do(API_BASE, "POST", "transaction", string(body), false)
	if err != nil {
		return res, err
	}
	var tx Transaction
	err = parseResponseData(r, &tx)
	if err != nil {
		return res, err
	}
	return instantswap.CreateResultInfo{
		ChargedFee:     utils.StrToFloat(tx.Fee.String()),
		Destination:    tx.Withdrawal,
		ExchangeRate:   utils.StrToFloat(tx.Rate.String()),
		FromCurrency:   tx.CoinFrom,
		InvoicedAmount: utils.StrToFloat(tx.DepositAmount.String()),
		OrderedAmount:  utils.StrToFloat(tx.WithdrawalAmount.String()),
		ToCurrency:     tx.CoinTo,
		UUID:           tx.TransactionId,
		DepositAddress: tx.Deposit,
		Expires:        0,
		ExtraID:        "",
		PayoutExtraID:  "",
	}, err
}

//UpdateOrder accepts orderID value and more if needed per lib.
func (c *GoDEX) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}
func (c *GoDEX) CancelOrder(orderID string) (res string, err error) {
	return
}

//OrderInfo accepts orderID value and more if needed per lib.
func (c *GoDEX) OrderInfo(orderID string) (res instantswap.OrderInfoResult, err error) {
	var r []byte
	r, err = c.client.Do(API_BASE, "GET", fmt.Sprintf("transaction/%s", orderID), "", false)
	if err != nil {
		return res, err
	}
	var tx Transaction
	err = parseResponseData(r, &tx)
	if err != nil {
		return res, err
	}
	return instantswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     "",
		ReceiveAmount:  utils.StrToFloat(tx.RealWithdrawalAmount.String()),
		TxID:           tx.HashOut,
		Status:         tx.Status,
		InternalStatus: GetLocalStatus(tx.Status),
		Confirmations:  "",
	}, err
}
func (c *GoDEX) EstimateAmount(vars interface{}) (res instantswap.EstimateAmount, err error) {
	return
}

//GetLocalStatus translate local status to instantswap.Status.
func GetLocalStatus(status string) instantswap.Status {
	// closed, confirming, exchanging, expired, failed, finished, refunded, sending, verifying, waiting
	status = strings.ToLower(status)
	switch status {
	case "wait":
		return instantswap.OrderStatusNew
	case "confirmation":
		return instantswap.OrderStatusDepositReceived
	case "confirmed":
		return instantswap.OrderStatusDepositConfirmed
	case "exchanging":
		return instantswap.OrderStatusExchanging
	case "sending", "sending_confirmation":
		return instantswap.OrderStatusSending
	case "success":
		return instantswap.OrderStatusCompleted
	case "overdue":
		return instantswap.OrderStatusExpired
	case "error":
		return instantswap.OrderStatusFailed
	case "refunded":
		return instantswap.OrderStatusRefunded
	default:
		return instantswap.OrderStatusUnknown
	}
}
