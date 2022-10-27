package instanttest

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"code.cryptopower.dev/exchange/instantswap"
)

const (
	LIBNAME = "instanttest"
	RateMin = 0.00500000
	RateMax = 0.00530000
)

var (
	trys                    = 0 //amount of test trys for returning test data
	storedTmpExchangeRate   = 0.0
	storedTmpInvoicedAmount = 0.0
)

func init() {
	instantswap.RegisterExchange(LIBNAME, func(config instantswap.ExchangeConfig) (instantswap.IDExchange, error) {
		return New(config)
	})
}

// New return an InstantTest struct.
func New(conf instantswap.ExchangeConfig) (*InstantTest, error) {
	client := instantswap.NewClient(LIBNAME, &conf)
	return &InstantTest{
		client: client,
		conf:   &conf,
	}, nil
}

// InstantTest represent a InstantTest client.
type InstantTest struct {
	client *instantswap.Client
	conf   *instantswap.ExchangeConfig
	instantswap.IDExchange
}

// GetExchangeRateInfo get estimate on the amount for the exchange.
func (c *InstantTest) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {

	limits, err := c.QueryLimits(vars.From, vars.To)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	time.Sleep(time.Second * 1)
	exchangeRates, err := c.QueryRates(nil)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	var rate instantswap.QueryRate
	var pair = fmt.Sprintf("%s-%s", vars.From, vars.To)
	for _, v := range exchangeRates {
		if v.Name == pair {
			rate = v
		}
	}
	if rate.Name == "" || rate.Value == "" {
		err = errors.New(LIBNAME + ":error: rate not found for " + pair + " pair")
		return
	}

	exchangeRate, err := strconv.ParseFloat(rate.Value, 64)
	if err != nil {
		err = errors.New(LIBNAME + ":error: " + err.Error())
		return
	}
	rateFinal := 1 / exchangeRate
	min := limits.Min * rateFinal
	max := limits.Max * rateFinal
	res = instantswap.ExchangeRateInfo{
		ExchangeRate:    rateFinal,
		Min:             min,
		Max:             max,
		EstimatedAmount: (vars.Amount / rateFinal),
	}

	return
}

// EstimateAmount get estimate on the amount for the exchange.
func (c *InstantTest) EstimateAmount(vars interface{}) (res instantswap.EstimateAmount, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

// QueryRates (list of pairs LTC-BTC, BTC-LTC, etc).
func (c *InstantTest) QueryRates(vars interface{}) (res []instantswap.QueryRate, err error) {
	randExchangeRate := randFloats(RateMin, RateMax, 1)
	rateStr := fmt.Sprintf("%.8f", (1.0 / randExchangeRate[0]))
	storedTmpExchangeRate = 0.0
	var tmpArr []instantswap.QueryRate
	tmpArr = append(tmpArr, instantswap.QueryRate{Name: "BTC-DCR", Value: rateStr})
	tmpArr = append(tmpArr, instantswap.QueryRate{Name: "BTC-LTC", Value: rateStr})
	res = tmpArr
	storedTmpExchangeRate = randExchangeRate[0]
	return
}
func (c *InstantTest) QueryActiveCurrencies(vars interface{}) (res []instantswap.ActiveCurr, err error) {
	var tmpArr []instantswap.ActiveCurr
	tmpArr = append(tmpArr, instantswap.ActiveCurr{Name: "Bitcoin", Code: "BTC", Precision: 8})
	tmpArr = append(tmpArr, instantswap.ActiveCurr{Name: "Decred", Code: "DCR", Precision: 8})
	tmpArr = append(tmpArr, instantswap.ActiveCurr{Name: "Litecoin", Code: "LTC", Precision: 8})
	res = tmpArr

	return
}

// QueryLimits Get Exchange Rates (from, to).
func (c *InstantTest) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {

	res = instantswap.QueryLimits{
		Max: 127.1757980601,
		Min: 0.5,
	}
	return
}

func (c *InstantTest) CreateOrder(orderInfo instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	storedTmpInvoicedAmount = 0
	res = instantswap.CreateResultInfo{
		ChargedFee:     0.001,
		Destination:    orderInfo.Destination,
		ExchangeRate:   storedTmpExchangeRate,
		FromCurrency:   orderInfo.FromCurrency,
		InvoicedAmount: orderInfo.InvoicedAmount,
		OrderedAmount:  0.0,
		ToCurrency:     orderInfo.ToCurrency,
		UUID:           "123456789",
		DepositAddress: "notadepositaddress", //from accept order result
		Expires:        1440,                 //from accept order result
	}
	storedTmpInvoicedAmount = orderInfo.InvoicedAmount
	return
}

func (c *InstantTest) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	orderInfo := vars.(UpdateOrder)

	randExchangeRate := randFloats(RateMin, RateMax, 1)

	res = instantswap.UpdateOrderResultInfo{
		ChargedFee:     0.001,
		Destination:    orderInfo.Order.Destination,
		ExchangeRate:   randExchangeRate[0],
		FromCurrency:   "BTC",
		InvoicedAmount: orderInfo.Order.OrderedAmount,
		OrderedAmount:  0.0,
		ToCurrency:     "DCR",
		UUID:           "123456789",
	}

	return
}

func (c *InstantTest) CancelOrder(orderID string) (res string, err error) {
	res = "success"
	return
}

// OrderInfo accepts string of orderID value.
func (c *InstantTest) OrderInfo(orderID string) (res instantswap.OrderInfoResult, err error) {
	trys++
	if trys == 1 {
		res = instantswap.OrderInfoResult{
			Expires:        1440,
			ReceiveAmount:  (storedTmpInvoicedAmount / storedTmpExchangeRate),
			Confirmations:  "0",
			TxID:           "",
			Status:         "Waiting on deposit",
			InternalStatus: instantswap.OrderStatusWaitingForDeposit,
		}
	} else {
		res = instantswap.OrderInfoResult{
			Expires:        1440,
			ReceiveAmount:  (storedTmpInvoicedAmount / storedTmpExchangeRate),
			Confirmations:  "6",
			TxID:           "e11525fe2e057fb19ec741ddcb972ec994f70348646368d960446a92c4d76dad",
			Status:         "Completed",
			InternalStatus: instantswap.OrderStatusCompleted,
		}
	}

	return
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
		precConvStr := fmt.Sprintf("%.8f", res[i])
		precConvF, err := strconv.ParseFloat(precConvStr, 64)
		if err != nil {
			fmt.Printf("%s:error: precision float conversion error, err: %s", err.Error())
		}
		res[i] = precConvF
	}
	return res
}
