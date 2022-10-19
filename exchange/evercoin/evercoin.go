package evercoin

import (
	"fmt"
	"time"
	//"strconv"
	"encoding/json"
	"errors"

	"code.cryptopower.dev/exchange/lightningswap"
)

const (
	API_BASE = "https://test.evercoin.com/v1/" // TEST API endpoint
	//API_BASE                   = "https://api.evercoin.com/v1/" // LIVE API endpoint

	//NOTE:
	//EverCoin exchange requires you to email support@evercoin.com to get a production api key

	DEFAULT_HTTPCLIENT_TIMEOUT = 30 // HTTP client timeout
	LIBNAME                    = "evercoin"
)

func init() {
	lightningswap.RegisterExchange(LIBNAME, func(config lightningswap.ExchangeConfig) (lightningswap.IDExchange, error) {
		return New(config)
	})
}

// New return a instanciate struct
func New(conf lightningswap.ExchangeConfig) (*EverCoin, error) {
	return nil, fmt.Errorf("evercoin is added time ago. Now I can not find any link or docs to this exchange. It should be removed")
	/*client := lightningswap.NewClient(LIBNAME, &conf)
	return &EverCoin{
		client: client,
		conf:   &conf,
	}, nil*/
}

//EverCoin represent a EverCoin client
type EverCoin struct {
	client *lightningswap.Client
	conf   *lightningswap.ExchangeConfig
	lightningswap.IDExchange
}

//SetDebug set enable/disable http request/response dump
func (c *EverCoin) SetDebug(enable bool) {
	c.conf.Debug = enable
}
func handleErr(r json.RawMessage) (err error) {
	var errorVals ErrorMsg
	if err = json.Unmarshal(r, &errorVals); err != nil {
		return err
	}
	if errorVals != (ErrorMsg{}) {
		err = errors.New(errorVals.Message)
		return err
	}
	return nil
}

const (
	waitSec = 3
)

//CalculateExchangeRate get estimate on the amount for the exchange
func (c *EverCoin) GetExchangeRateInfo(vars lightningswap.ExchangeRateRequest) (res lightningswap.ExchangeRateInfo, err error) {
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

	rate := vars.Amount / estimate.EstimatedAmount

	res = lightningswap.ExchangeRateInfo{
		ExchangeRate: rate,
		Min:          limits.Min,
		Max:          limits.Max,
		Signature:    estimate.Signature,
	}

	return
}

//EstimateAmount get estimate on the amount for the exchange
func (c *EverCoin) EstimateAmount(vars lightningswap.ExchangeRateRequest) (res lightningswap.EstimateAmount, err error) {
	//amountStr := strconv.FormatFloat(amount, 'f', 8, 64)
	estAmount := EstimateAmount{
		DepositAmount:   vars.Amount,
		DestinationCoin: vars.To,
		DepositCoin:     vars.From,
	}

	payload, err := json.Marshal(estAmount)
	if err != nil {
		return
	}
	fmt.Printf("\n %v", string(payload))
	r, err := c.client.Do(API_BASE, "POST", "price", string(payload), false)
	if err != nil {
		return
	}

	var tmp EstimateAmountResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		return
	}
	err = handleErr(tmp.Error)
	if err != nil {
		return
	}

	res = lightningswap.EstimateAmount{
		DepositAmount:   tmp.Result.DepositAmount,
		EstimatedAmount: tmp.Result.DestinationAmount,
		ToCurrency:      tmp.Result.DestinationCoin,
		FromCurrency:    tmp.Result.DepositCoin,
		Signature:       tmp.Result.Signature,
	}
	return
}

//QueryRates (list of pairs LTC-BTC, BTC-LTC, etc)
func (c *EverCoin) QueryRates(vars interface{}) (res []lightningswap.QueryRate, err error) {
	err = errors.New(LIBNAME + ":error:queryrates not available for this exchange")
	return
}

//QueryActiveCurrencies get all active currencies
func (c *EverCoin) QueryActiveCurrencies(vars interface{}) (res []lightningswap.ActiveCurr, err error) {
	err = errors.New(LIBNAME + ":error:queryactive not available for this exchange")
	return
}

//QueryLimits Get Exchange Rates (from, to)
func (c *EverCoin) QueryLimits(fromCurr, toCurr string) (res lightningswap.QueryLimits, err error) {
	r, err := c.client.Do(API_BASE, "GET", "limit/"+fromCurr+"-"+toCurr, "", false)
	if err != nil {
		return
	}
	var tmp QueryLimits
	if err = json.Unmarshal(r, &tmp); err != nil {
		return
	}
	err = handleErr(tmp.Error)
	if err != nil {
		return
	}

	res = lightningswap.QueryLimits{
		Max: tmp.Result.MaxDeposit,
		Min: tmp.Result.MinDeposit,
	}
	return
}

//GetCoins get coin list
func (c *EverCoin) GetCoins() (res []CoinInfo, err error) {
	r, err := c.client.Do(API_BASE, "GET", "coins", "", false)
	if err != nil {
		return
	}
	var tmp Coins
	if err = json.Unmarshal(r, &tmp); err != nil {
		return
	}
	err = handleErr(tmp.Error)
	if err != nil {
		return
	}

	res = tmp.Coins
	return
}

func (c *EverCoin) CreateOrder(orderInfo lightningswap.CreateOrder) (res lightningswap.CreateResultInfo, err error) {
	//TODO: call estimateamount first then send that info to api
	estAmount, err := c.EstimateAmount(lightningswap.ExchangeRateRequest{
		From:   orderInfo.FromCurrency,
		To:     orderInfo.ToCurrency,
		Amount: orderInfo.InvoicedAmount,
	})
	if err != nil {
		return
	}

	depositCoin := Address{
		MainAddress: orderInfo.RefundAddress,
		TagValue:    "",
	}
	destinationCoin := Address{
		MainAddress: orderInfo.Destination,
		TagValue:    "",
	}
	//translate from interface struct to local struct for POST
	newOrder := CreateOrder{
		DepositAmount:      estAmount.DepositAmount,
		DepositCoin:        orderInfo.FromCurrency,
		DestinationAddress: destinationCoin,
		DestinationAmount:  estAmount.EstimatedAmount,
		DestinationCoin:    orderInfo.ToCurrency,
		RefundAddress:      depositCoin,
		Signature:          estAmount.Signature,
	}
	payload, err := json.Marshal(newOrder)
	if err != nil {
		return
	}
	r, err := c.client.Do(API_BASE, "POST", "order", string(payload), false)
	if err != nil {
		return
	}
	var tmp CreateResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		return
	}
	err = handleErr(tmp.Error)
	if err != nil {
		return
	}

	res = lightningswap.CreateResultInfo{
		UUID:           tmp.Order.UUID,
		DepositAddress: tmp.Order.DepositAddress.MainAddress,
	}

	return
}
func (c *EverCoin) UpdateOrder(vars interface{}) (res lightningswap.UpdateOrderResultInfo, err error) {
	err = errors.New(LIBNAME + ":error:update not available for this exchange")
	return
}

func (c *EverCoin) CancelOrder(oId string) (res string, err error) {
	err = errors.New(LIBNAME + ":error:cancel not available for this exchange")
	return
}

func (c *EverCoin) OrderInfo(orderID string) (res lightningswap.OrderInfoResult, err error) {
	r, err := c.client.Do(API_BASE, "GET", "status/"+orderID, "", false)
	if err != nil {
		return
	}
	var tmp OrderInfoResult
	if err = json.Unmarshal(r, &tmp); err != nil {
		return
	}
	err = handleErr(tmp.Error)
	if err != nil {
		return
	}

	return
}
func GetLocalStatus(statusID int) (iStatus int) {

	switch statusID {
	case 5: // All_Done
		return 1
	case 1: // Awaiting_Deposit
		return 2
	case 2: // Awaiting_Confirm
		return 3
	case 3: // Awaiting_Exchange
		return 4
	case 6:
		return 5
	case 10: // Canceled
		return 6
	case 9: // Expire_Exchange
		return 7
	default:
		return 0
	}
}
func GetOrderStatusString(statusID int) (res string) {
	switch statusID {
	case 1:
		return "Awaiting Deposit"
	case 2:
		return "Awaiting Confirm"
	case 3:
		return "Awaiting Exchange"
	case 4:
		return "Awaiting Refund"
	case 5:
		return "All Done"
	case 6:
		return "Refund Done"
	case 7:
		return "Minimum Cancel"
	case 8:
		return "Send Money Error"
	case 9:
		return "Expire Exchange"
	case 10:
		return "Cancel"
	}
	return "status n/a"
}
