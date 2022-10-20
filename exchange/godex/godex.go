package godex

import (
	"code.cryptopower.dev/exchange/lightningswap"
	"code.cryptopower.dev/exchange/lightningswap/utils"
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
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

type GoDEX struct {
	conf   *lightningswap.ExchangeConfig
	client *lightningswap.Client
	lightningswap.IDExchange
}

func New(conf lightningswap.ExchangeConfig) (*GoDEX, error) {
	if conf.ApiKey == "" {
		return nil, fmt.Errorf("%s:error: APIKEY is blank", LIBNAME)
	}
	client := lightningswap.NewClient(LIBNAME, &conf, func(r *http.Request, body string) error {
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

func (c *GoDEX) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
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
	return lightningswap.ExchangeRateInfo{
		Min:             utils.StrToFloat(info.MinAmount.String()),
		Max:             utils.StrToFloat(info.MaxAmount.String()),
		ExchangeRate:    1 / utils.StrToFloat(info.Rate.String()),
		EstimatedAmount: estimatedAmount,
		MaxOrder:        0,
		Signature:       "",
	}, err
}

func (c *GoDEX) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	return res, fmt.Errorf("not supported")
}
func (c *GoDEX) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	return
}
func (c *GoDEX) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {
	return
}
func (c *GoDEX) CreateOrder(vars lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
	fmt.Println("vars.OrderedAmount:", vars.OrderedAmount)
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
	return lightningswap.CreateResultInfo{
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

//UpdateOrder accepts orderID value and more if needed per lib
func (c *GoDEX) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	return
}
func (c *GoDEX) CancelOrder(orderID string) (res string, err error) {
	return
}

//OrderInfo accepts orderID value and more if needed per lib
func (c *GoDEX) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
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
	return lightningswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     "",
		ReceiveAmount:  utils.StrToFloat(tx.RealWithdrawalAmount.String()),
		TxID:           tx.HashOut,
		Status:         tx.Status,
		InternalStatus: GetLocalStatus(tx.Status),
		Confirmations:  "",
	}, err
}
func (c *GoDEX) EstimateAmount(vars interface{}) (res lightningswap.EstimateAmount, err error) {
	return
}

//GetLocalStatus translate local status to idexchange status id
func GetLocalStatus(status string) lightningswap.Status {
	// closed, confirming, exchanging, expired, failed, finished, refunded, sending, verifying, waiting
	status = strings.ToLower(status)
	switch status {
	case "wait":
		return lightningswap.OrderStatusNew
	case "confirmation":
		return lightningswap.OrderStatusDepositReceived
	case "confirmed":
		return lightningswap.OrderStatusDepositConfirmed
	case "exchanging":
		return lightningswap.OrderStatusExchanging
	case "sending", "sending_confirmation":
		return lightningswap.OrderStatusSending
	case "success":
		return lightningswap.OrderStatusCompleted
	case "overdue":
		return lightningswap.OrderStatusExpired
	case "error":
		return lightningswap.OrderStatusFailed
	case "refunded":
		return lightningswap.OrderStatusRefunded
	default:
		return lightningswap.OrderStatusUnknown
	}
}
