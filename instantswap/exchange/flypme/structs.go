package flypme

import (
	"encoding/json"
)

//base json structure
type jsonResponse struct {
	Errors json.RawMessage `json:"errors"`
	Result string          `json:"result"`
	Status string          `json:"status"`
}

//QUERY

type QueryRate struct {
	Name  string
	Value string
}
type ActiveCurrs struct {
	Name       string
	Currencies []ActiveCurr
}
type ActiveCurr struct {
	ChargedFee       float64 `json:"charged_fee,string"`
	Code             string  `json:"code"`
	ConfirmationTime int     `json:"confirmation_time"`
	CreatedAt        string  `json:"created_at"`
	CurrencyType     string  `json:"currency_type"`
	Default          bool    `json:"default"`
	DisplayPrecision int     `json:"display_precision"`
	Exchange         bool    `json:"exchange"`
	Name             string  `json:"name"`
	Precision        int     `json:"precision"`
	Send             bool    `json:"send"`
	UpdatedAt        string  `json:"updated_at"`
	Website          string  `json:"website"`
}

type QueryLimits struct {
	Max float64 `json:"max,string"`
	Min float64 `json:"min,string"`
}

//CREATE
type CreateOrderInfo struct {
	RefundAddress  string `json:"refund_address"`
	Destination    string `json:"destination"`
	FromCurrency   string `json:"from_currency"`
	OrderedAmount  string `json:"ordered_amount"`
	InvoicedAmount string `json:"invoiced_amount"`
	ToCurrency     string `json:"to_currency"`
}
type CreateOrder struct {
	Order CreateOrderInfo `json:"order"`
}
type CreateResultInfo struct {
	ChargedFee     float64 `json:"charged_fee,string"`
	Destination    string  `json:"destination"`
	ExchangeRate   float64 `json:"exchange_rate,string"`
	FromCurrency   string  `json:"from_currency"`
	InvoicedAmount float64 `json:"invoiced_amount,string"`
	OrderedAmount  float64 `json:"ordered_amount,string"`
	ToCurrency     string  `json:"to_currency"`
	UUID           string  `json:"uuid"`
}
type CreateResult struct {
	Errors  json.RawMessage  `json:"errors"`
	Expires int              `json:"expires"`
	Order   CreateResultInfo `json:"order"`
}

//UPDATE
type UpdateOrderInfo struct {
	Destination   string  `json:"destination"`
	OrderedAmount float64 `json:"ordered_amount,string"`
	RefundAddress string  `json:"refund_address"`
	UUID          string  `json:"uuid"`
}
type UpdateOrder struct {
	Order UpdateOrderInfo `json:"order"`
}
type UpdateOrderResultInfo struct {
	ChargedFee     float64 `json:"charged_fee,string"`
	Destination    string  `json:"destination"`
	ExchangeRate   float64 `json:"exchange_rate,string"`
	FromCurrency   string  `json:"from_currency"`
	InvoicedAmount float64 `json:"invoiced_amount,string"`
	OrderedAmount  float64 `json:"ordered_amount,string"`
	ToCurrency     string  `json:"to_currency"`
	UUID           string  `json:"uuid"`
}
type UpdateOrderResult struct {
	Errors  json.RawMessage       `json:"errors"`
	Expires int                   `json:"expires"`
	Order   UpdateOrderResultInfo `json:"order"`
}

type AcceptOrderResultInfo struct {
	ChargedFee     string `json:"charged_fee"`
	Destination    string `json:"destination"`
	ExchangeRate   string `json:"exchange_rate"`
	FromCurrency   string `json:"from_currency"`
	InvoicedAmount string `json:"invoiced_amount"`
	OrderedAmount  string `json:"ordered_amount"`
	ToCurrency     string `json:"to_currency"`
	UUID           string `json:"uuid"`
}
type AcceptOrderResult struct {
	Errors         json.RawMessage       `json:"errors"`
	DepositAddress string                `json:"deposit_address"`
	Expires        int                   `json:"expires"`
	Order          AcceptOrderResultInfo `json:"order"`
}

//INFO
type OrderInfoResultInfo struct {
	ChargedFee     float64 `json:"charged_fee,string"`
	Destination    string  `json:"destination"`
	ExchangeRate   float64 `json:"exchange_rate,string"`
	FromCurrency   string  `json:"from_currency"`
	InvoicedAmount float64 `json:"invoiced_amount,string"`
	OrderedAmount  float64 `json:"ordered_amount,string"`
	ToCurrency     string  `json:"to_currency"`
	UUID           string  `json:"uuid"`
}
type OrderInfoResult struct {
	Errors         json.RawMessage     `json:"errors"`
	DepositAddress string              `json:"deposit_address"`
	TxID           string              `json:"txid"`
	TxURL          string              `json:"txurl"`
	Expires        int                 `json:"expires"`
	Order          OrderInfoResultInfo `json:"order"`
	Status         string              `json:"status"`
	Confirmations  string              `json:"confirmations"`
}

type UUID struct {
	UUID string `json:"uuid"`
}

type MyJsonName struct {
	DepositAddress string `json:"deposit_address"`
	Expires        int    `json:"expires"`
	Order          struct {
	} `json:"order"`
	Status string `json:"status"`
}

type Currency struct {
	Code             string `json:"code"`
	Precision        int    `json:"precision"`
	DisplayPrecision int    `json:"display_precision"`
	Name             string `json:"name"`
	Website          string `json:"website"`
	ConfirmationTime int    `json:"confirmation_time"`
	Default          bool   `json:"default"`
	ChargedFee       string `json:"charged_fee"`
	CurrencyType     string `json:"currency_type"`
	Exchange         bool   `json:"exchange"`
	Send             bool   `json:"send"`
	Stake            bool   `json:"stake"`
	NewAddresses     bool   `json:"new_addresses"`
	Change24H        string `json:"change_24h"`
}
