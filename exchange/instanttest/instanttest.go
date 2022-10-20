package instanttest

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"code.cryptopower.dev/exchange/lightningswap"
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
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

// New return a instanciate struct
func New(conf lightningswap.ExchangeConfig) (*InstantTest, error) {
	client := lightningswap.NewClient(LIBNAME, &conf)
	return &InstantTest{
		client: client,
		conf:   &conf,
	}, nil
}

//InstantTest represent a InstantTest client
type InstantTest struct {
	client *lightningswap.Client
	conf   *lightningswap.ExchangeConfig
	lightningswap.IDExchange
}

//CalculateExchangeRate get estimate on the amount for the exchange
func (c *InstantTest) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {

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
	var rate lightningswap.QueryRate
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
	//fmt.Printf("\nmin: %v max: %v", min, max )
	res = lightningswap.ExchangeRateInfo{
		ExchangeRate:    rateFinal,
		Min:             min,
		Max:             max,
		EstimatedAmount: (vars.Amount / rateFinal),
	}

	return
}

//EstimateAmount get estimate on the amount for the exchange
func (c *InstantTest) EstimateAmount(vars interface{}) (res lightningswap.EstimateAmount, err error) {
	//vars not used here
	err = errors.New(LIBNAME + ":error: not available for this exchange")
	return
}

//QueryRates (list of pairs LTC-BTC, BTC-LTC, etc)
func (c *InstantTest) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	randExchangeRate := randFloats(RateMin, RateMax, 1)
	rateStr := fmt.Sprintf("%.8f", (1.0 / randExchangeRate[0]))
	storedTmpExchangeRate = 0.0
	var tmpArr []lightningswap.QueryRate
	tmpArr = append(tmpArr, lightningswap.QueryRate{Name: "BTC-DCR", Value: rateStr})
	tmpArr = append(tmpArr, lightningswap.QueryRate{Name: "BTC-LTC", Value: rateStr})
	res = tmpArr
	storedTmpExchangeRate = randExchangeRate[0]
	return
}
func (c *InstantTest) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	var tmpArr []lightningswap.ActiveCurr
	tmpArr = append(tmpArr, lightningswap.ActiveCurr{Name: "Bitcoin", Code: "BTC", Precision: 8})
	tmpArr = append(tmpArr, lightningswap.ActiveCurr{Name: "Decred", Code: "DCR", Precision: 8})
	tmpArr = append(tmpArr, lightningswap.ActiveCurr{Name: "Litecoin", Code: "LTC", Precision: 8})
	res = tmpArr

	return
}

//QueryLimits Get Exchange Rates (from, to)
func (c *InstantTest) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {

	res = lightningswap.QueryLimits{
		Max: 127.1757980601,
		Min: 0.5,
	}
	return
}
func (c *InstantTest) CreateOrder(orderInfo lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
	storedTmpInvoicedAmount = 0
	res = lightningswap.CreateResultInfo{
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
func (c *InstantTest) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	orderInfo := vars.(UpdateOrder)

	randExchangeRate := randFloats(RateMin, RateMax, 1)

	res = lightningswap.UpdateOrderResultInfo{
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

//OrderInfo accepts string of orderID value
func (c *InstantTest) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
	trys++
	fmt.Printf("\ntrys %v", trys)
	if trys == 1 {
		res = lightningswap.OrderInfoResult{
			//LastUpdate:   not available
			Expires:        1440,
			ReceiveAmount:  (storedTmpInvoicedAmount / storedTmpExchangeRate),
			Confirmations:  "0",
			TxID:           "",
			Status:         "Waiting on deposit",
			InternalStatus: lightningswap.OrderStatusWaitingForDeposit,
		}
	} else {
		res = lightningswap.OrderInfoResult{
			//LastUpdate:   not available
			Expires:        1440,
			ReceiveAmount:  (storedTmpInvoicedAmount / storedTmpExchangeRate),
			Confirmations:  "6",
			TxID:           "e11525fe2e057fb19ec741ddcb972ec994f70348646368d960446a92c4d76dad",
			Status:         "Completed",
			InternalStatus: lightningswap.OrderStatusCompleted,
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
