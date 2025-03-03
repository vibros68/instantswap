package exchcx

import (
	"encoding/json"
	"fmt"
	"github.com/crypto-power/instantswap/instantswap"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type exchCx struct {
}

func (e *exchCx) Do(req *http.Request, resObj any) error {
	client := &http.Client{}
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	var exchErr Error
	_ = json.Unmarshal(body, &exchErr)
	if exchErr.Error != "" {
		return fmt.Errorf(exchErr.Error)
	}
	return json.Unmarshal(body, resObj)
}

func (e *exchCx) path(path string) string {
	return fmt.Sprintf("https://exch.cx/api/%s", path)
}

func (e *exchCx) GetCurrencies() (currencies []instantswap.Currency, err error) {
	req, err := http.NewRequest(http.MethodGet, e.path("volume"), nil)
	if err != nil {
		return
	}
	var volumnMap map[string]Volume
	err = e.Do(req, &volumnMap)
	if err != nil {
		return
	}
	for currency, _ := range volumnMap {
		currencies = append(currencies, instantswap.Currency{
			Symbol: currency,
		})
	}
	return
}

func (e *exchCx) GetCurrenciesToPair(from string) (currencies []instantswap.Currency, err error) {
	req, err := http.NewRequest(http.MethodGet, e.path("rates"), nil)
	if err != nil {
		return
	}
	from = strings.ToUpper(from)
	var rateMap map[string]Rate
	err = e.Do(req, &rateMap)
	if err != nil {
		return
	}
	for currencyPair, _ := range rateMap {
		var pair = strings.Split(currencyPair, "_")
		if len(pair) == 2 && pair[0] == from {
			currencies = append(currencies, instantswap.Currency{
				Symbol: pair[1],
			})
		}
	}
	return
}

func (e *exchCx) QueryLimits(fromCurr, toCurr string) (res instantswap.QueryLimits, err error) {
	return
}

func (e *exchCx) CreateOrder(vars instantswap.CreateOrder) (res instantswap.CreateResultInfo, err error) {
	var params = url.Values{}
	params.Set("from_currency", vars.FromCurrency)
	params.Set("to_currency", vars.ToCurrency)
	params.Set("to_address", vars.Destination)
	params.Set("refund_address", vars.RefundAddress)
	params.Set("rate_mode", "flat")
	params.Set("fee_option", "s")
	req, err := http.NewRequest(http.MethodGet, e.path("create?"+params.Encode()), nil)
	if err != nil {
		return res, err
	}
	var createResponse struct {
		OrderId string `json:"orderid"`
	}
	err = e.Do(req, &createResponse)
	if err != nil {
		return res, err
	}
	order, err := e.getOrder(createResponse.OrderId)
	if err != nil {
		return res, err
	}
	res = instantswap.CreateResultInfo{
		ChargedFee:     0,
		Destination:    order.ToAddress,
		ExchangeRate:   order.Rate,
		FromCurrency:   order.FromCurrency,
		InvoicedAmount: 0,
		OrderedAmount:  0,
		ToCurrency:     order.ToCurrency,
		UUID:           "",
		DepositAddress: order.FromAddr,
		Expires:        0,
		ExtraID:        "",
		PayoutExtraID:  "",
	}
	return
}

func (e *exchCx) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}

func (e *exchCx) CancelOrder(orderID string) (res string, err error) {
	return
}

func (e *exchCx) getOrder(orderId string) (*Order, error) {
	req, err := http.NewRequest(http.MethodGet, e.path("order?orderid="+orderId), nil)
	if err != nil {
		return nil, err
	}
	var order Order
	err = e.Do(req, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (e *exchCx) OrderInfo(orderID string, extraIds ...string) (res instantswap.OrderInfoResult, err error) {
	order, err := e.getOrder(orderID)
	if err != nil {
		return
	}
	res = instantswap.OrderInfoResult{
		Expires:        0,
		LastUpdate:     "",
		ReceiveAmount:  *order.ToAmount,
		TxID:           *order.TransactionIdSent,
		Status:         order.State,
		InternalStatus: statusMap[order.State],
		Confirmations:  "",
	}
	return
}

func (e *exchCx) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	req, err := http.NewRequest(http.MethodGet, e.path("rates"), nil)
	if err != nil {
		return
	}
	var rateMap map[string]Rate
	err = e.Do(req, &rateMap)
	if err != nil {
		return
	}
	var pair = strings.ToUpper(vars.From) + "_" + strings.ToUpper(vars.To)
	var rate, ok = rateMap[pair]
	if !ok {
		err = fmt.Errorf("exchange rate info not found")
		return
	}
	res.ExchangeRate = rate.Rate
	return
}

var statusMap = map[string]instantswap.Status{
	"CREATED":           instantswap.OrderStatusNew,
	"CANCELLED":         instantswap.OrderStatusCanceled,
	"AWAITING_INPUT":    instantswap.OrderStatusWaitingForDeposit,
	"CONFIRMING_INPUT":  instantswap.OrderStatusWaitingForDeposit,
	"EXCHANGING":        instantswap.OrderStatusExchanging,
	"FUNDED":            instantswap.OrderStatusDepositConfirmed,
	"BRIDGING":          instantswap.OrderStatusSending,
	"CONFIRMING_SEND":   instantswap.OrderStatusSending,
	"COMPLETE":          instantswap.OrderStatusCompleted,
	"REFUND_REQUEST":    instantswap.OrderStatusRefunded,
	"REFUND_PENDING":    instantswap.OrderStatusRefunded,
	"CONFIRMING_REFUND": instantswap.OrderStatusRefunded,
	"REFUNDED":          instantswap.OrderStatusRefunded,
}
