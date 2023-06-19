package changelly

import (
	"encoding/json"
)

// base json structure
type jsonRequest struct {
	ID      string      `json:"id"`
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}
type jsonResponse struct {
	ID      string          `json:"id"`
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   json.RawMessage `json:"error"`
}
type jsonError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
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
	Featured      bool   `json:"featured"`
	HasExternalID bool   `json:"hasExternalId"`
	Image         string `json:"image"`
	IsFiat        bool   `json:"isFiat"`
	Ticker        string `json:"ticker"`
}

type QueryLimits struct {
	//Max string `json:"max"`
	Min float64 `json:"minAmount"`
}

// CREATE
type CreateOrder struct {
	FromCurrency      string  `json:"from"`
	ToCurrency        string  `json:"to"`
	ToCurrencyAddress string  `json:"address"`
	InvoicedAmount    float64 `json:"amount"`            //amount in "from" currency
	ExtraID           string  `json:"extraID,omitempty"` //optional for some coins
}
type CreateResult struct {
	UUID          string  `json:"id"`
	AmountTo      float64 `json:"amountTo"` //0 until amount has been deposited based on api docs
	APIExtraFee   float64 `json:"apiExtraFee,string"`
	ChangellyFee  float64 `json:"changellyFee,string"`
	CreatedAt     string  `json:"createdAt"`
	CurrencyFrom  string  `json:"currencyFrom"`
	CurrencyTo    string  `json:"currencyTo"`
	PayinAddress  string  `json:"payinAddress"`
	PayinExtraID  string  `json:"payinExtraId"`
	PayoutAddress string  `json:"payoutAddress"`
	PayoutExtraID string  `json:"payoutExtraId"`
	RefundAddress string  `json:"refundAddress"`
	RefundExtraID string  `json:"refundExtraId"`
	Status        string  `json:"status"`
}

//INFO

type UUID struct {
	UUID   string `json:"id"`
	APIKEY string `json:"api_key"`
}

type OrderInfoResult struct {
	AmountFrom         string      `json:"amountFrom"`
	AmountTo           float64     `json:"amountTo,string"`
	APIExtraFee        float64     `json:"apiExtraFee,string"`
	ChangellyFee       float64     `json:"changellyFee,string"`
	CreatedAt          int         `json:"createdAt"`
	CurrencyFrom       string      `json:"currencyFrom"`
	CurrencyTo         string      `json:"currencyTo"`
	UUID               string      `json:"id"`
	NetworkFee         interface{} `json:"networkFee"`
	PayinAddress       string      `json:"payinAddress"`
	PayinConfirmations string      `json:"payinConfirmations"`
	PayinExtraID       string      `json:"payinExtraId"`
	PayinHash          string      `json:"payinHash"`
	PayoutAddress      string      `json:"payoutAddress"`
	PayoutExtraID      string      `json:"payoutExtraId"`
	PayoutHash         string      `json:"payoutHash"`
	Status             string      `json:"status"`
}
type EstimateAmount struct {
	EstimatedAmount          float64     `json:"estimatedAmount"` //destinationCurrency
	NetworkFee               float64     `json:"networkFee"`
	ServiceCommission        float64     `json:"serviceCommission"`
	TransactionSpeedForecast string      `json:"transactionSpeedForecast"`
	WarningMessage           interface{} `json:"warningMessage"`
}
