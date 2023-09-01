package fixedfloat

import (
	"encoding/json"
	"fmt"
)

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data json.RawMessage
}

func parseResponseData(r []byte, obj interface{}) error {
	var res response
	err := json.Unmarshal(r, &res)
	if err != nil {
		return err
	}
	if err == nil && res.Code > 0 {
		return fmt.Errorf(res.Msg)
	}
	err = json.Unmarshal(res.Data, obj)
	return err
}

type Currency struct {
	Code     string      `json:"code"`
	Coin     string      `json:"coin"`
	Network  string      `json:"network"`
	Name     string      `json:"name"`
	Recv     interface{} `json:"recv"`
	Send     interface{} `json:"send"`
	Tag      interface{} `json:"tag"`
	Logo     string      `json:"logo"`
	Color    string      `json:"color"`
	Priority int         `json:"priority"`
}

// {"fromCcy":"BTC","toCcy":"USDTTRC","amount":0.5,"direction":"from","type":"float"}
type PriceReq struct {
	FromCcy   string  `json:"fromCcy"`
	ToCcy     string  `json:"toCcy"`
	Amount    float64 `json:"amount"`
	Direction string  `json:"direction"`
	Type      string  `json:"type"`
}

type PriceResult struct {
	From struct {
		Code      string  `json:"code"`
		Network   string  `json:"network"`
		Coin      string  `json:"coin"`
		Amount    float64 `json:"amount,string"`
		Rate      float64 `json:"rate,string"`
		Precision int     `json:"precision"`
		Min       float64 `json:"min,string"`
		Max       float64 `json:"max,string"`
		Usd       float64 `json:"usd,string"`
		Btc       float64 `json:"btc,string"`
	} `json:"from"`
	To struct {
		Code      string  `json:"code"`
		Network   string  `json:"network"`
		Coin      string  `json:"coin"`
		Amount    float64 `json:"amount,string"`
		Rate      float64 `json:"rate,string"`
		Precision int     `json:"precision"`
		Min       float64 `json:"min,string"`
		Max       float64 `json:"max,string"`
		Usd       float64 `json:"usd,string"`
	} `json:"to"`
	Errors []interface{} `json:"errors"`
}

type CreateOrderRequest struct {
	FromCcy   string  `json:"fromCcy"`
	ToCcy     string  `json:"toCcy"`
	Amount    float64 `json:"amount"`
	Direction string  `json:"direction"`
	Type      string  `json:"type"`
	ToAddress string  `json:"toAddress"`
}

type OrderResponse struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Email  string `json:"email"`
	Status string `json:"status"`
	Time   struct {
		Reg        int         `json:"reg"`
		Start      interface{} `json:"start"`
		Finish     interface{} `json:"finish"`
		Update     int         `json:"update"`
		Expiration int         `json:"expiration"`
		Left       int         `json:"left"`
	} `json:"time"`
	From struct {
		Code             string      `json:"code"`
		Coin             string      `json:"coin"`
		Network          string      `json:"network"`
		Name             string      `json:"name"`
		Alias            string      `json:"alias"`
		Amount           float64     `json:"amount,string"`
		Address          string      `json:"address"`
		AddressAlt       interface{} `json:"addressAlt"`
		Tag              interface{} `json:"tag"`
		TagName          interface{} `json:"tagName"`
		ReqConfirmations int         `json:"reqConfirmations"`
		MaxConfirmations int         `json:"maxConfirmations"`
		Tx               struct {
			Id            interface{} `json:"id"`
			Amount        interface{} `json:"amount"`
			Fee           interface{} `json:"fee"`
			Ccyfee        interface{} `json:"ccyfee"`
			TimeReg       interface{} `json:"timeReg"`
			TimeBlock     interface{} `json:"timeBlock"`
			Confirmations interface{} `json:"confirmations"`
		} `json:"tx"`
	} `json:"from"`
	To struct {
		Code    string      `json:"code"`
		Coin    string      `json:"coin"`
		Network string      `json:"network"`
		Name    string      `json:"name"`
		Alias   string      `json:"alias"`
		Amount  float64     `json:"amount,string"`
		Address string      `json:"address"`
		Tag     interface{} `json:"tag"`
		TagName interface{} `json:"tagName"`
		Tx      struct {
			Id            interface{} `json:"id"`
			Amount        interface{} `json:"amount"`
			Fee           interface{} `json:"fee"`
			Ccyfee        interface{} `json:"ccyfee"`
			TimeReg       interface{} `json:"timeReg"`
			TimeBlock     interface{} `json:"timeBlock"`
			Confirmations interface{} `json:"confirmations"`
		} `json:"tx"`
	} `json:"to"`
	Back struct {
		Code    interface{} `json:"code"`
		Coin    interface{} `json:"coin"`
		Network interface{} `json:"network"`
		Name    interface{} `json:"name"`
		Alias   interface{} `json:"alias"`
		Amount  interface{} `json:"amount"`
		Address interface{} `json:"address"`
		Tag     interface{} `json:"tag"`
		TagName interface{} `json:"tagName"`
		Tx      struct {
			Id            interface{} `json:"id"`
			Amount        interface{} `json:"amount"`
			Fee           interface{} `json:"fee"`
			Ccyfee        interface{} `json:"ccyfee"`
			TimeReg       interface{} `json:"timeReg"`
			TimeBlock     interface{} `json:"timeBlock"`
			Confirmations interface{} `json:"confirmations"`
		} `json:"tx"`
	} `json:"back"`
	Emergency struct {
		Status []interface{} `json:"status"`
		Choice string        `json:"choice"`
		Repeat string        `json:"repeat"`
	} `json:"emergency"`
	Token string `json:"token"`
}
