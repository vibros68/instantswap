package exchcx

import (
	"encoding/json"
	"fmt"
	"github.com/crypto-power/instantswap/instantswap"
	"io"
	"net/http"
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
	return
}

func (e *exchCx) UpdateOrder(vars interface{}) (res instantswap.UpdateOrderResultInfo, err error) {
	return
}

func (e *exchCx) CancelOrder(orderID string) (res string, err error) {
	return
}

func (e *exchCx) OrderInfo(orderID string, extraIds ...string) (res instantswap.OrderInfoResult, err error) {
	return
}

func (e *exchCx) GetExchangeRateInfo(vars instantswap.ExchangeRateRequest) (res instantswap.ExchangeRateInfo, err error) {
	return
}
