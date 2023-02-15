package changenow

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
	Featured      bool   `json:"featured"`
	HasExternalID bool   `json:"hasExternalId"`
	Image         string `json:"image"`
	IsFiat        bool   `json:"isFiat"`
	Ticker        string `json:"ticker"`
}

type QueryLimits struct {
	Max float64 `json:"max"`
	Min float64 `json:"minAmount"`
}

//CREATE
type CreateOrder struct {
	FromCurrency      string  `json:"from"`
	ToCurrency        string  `json:"to"`
	ToCurrencyAddress string  `json:"address"`
	InvoicedAmount    float64 `json:"amount"`            //amount in "from" currency
	ExtraID           string  `json:"extraID,omitempty"` //optional for some coins
}
type CreateResult struct {
	UUID               string `json:"id"`
	DepositAddress     string `json:"payinAddress"`
	DestinationAddress string `json:"payoutAddress"`
	PayinExtraID       string `json:"payinExtraId"`
	FromCurrency       string `json:"fromCurrency"`
	ToCurrency         string `json:"toCurrency"`
}

//INFO

type UUID struct {
	UUID   string `json:"id"`
	APIKEY string `json:"api_key"`
}

type OrderInfoResult struct {
	AmountReceive         float64 `json:"amountReceive"`
	AmountSend            float64 `json:"amountSend"`
	ExpectedAmountReceive float64 `json:"expectedReceiveAmount"`
	ExpectedAmountSend    float64 `json:"expectedSendAmount"`
	FromCurrency          string  `json:"fromCurrency"`
	Hash                  string  `json:"hash"`
	ID                    string  `json:"id"`
	NetworkFee            float64 `json:"networkFee,string"`
	PayinAddress          string  `json:"payinAddress"`
	PayinExtraID          string  `json:"payinExtraId"`
	PayinHash             string  `json:"payinHash"`
	PayoutAddress         string  `json:"payoutAddress"`
	PayoutExtraID         string  `json:"payoutExtraId"`
	PayoutHash            string  `json:"payoutHash"`
	Status                string  `json:"status"`
	ToCurrency            string  `json:"toCurrency"`
	UpdatedAt             string  `json:"updatedAt"`
}
type EstimateAmount struct {
	EstimatedAmount          float64     `json:"estimatedAmount"` //destinationCurrency
	NetworkFee               float64     `json:"networkFee"`
	ServiceCommission        float64     `json:"serviceCommission"`
	TransactionSpeedForecast string      `json:"transactionSpeedForecast"`
	WarningMessage           interface{} `json:"warningMessage"`
}

type Currency struct {
	Ticker            string `json:"ticker"`
	Name              string `json:"name"`
	Image             string `json:"image"`
	HasExternalId     bool   `json:"hasExternalId"`
	IsFiat            bool   `json:"isFiat"`
	Featured          bool   `json:"featured"`
	IsStable          bool   `json:"isStable"`
	SupportsFixedRate bool   `json:"supportsFixedRate"`
}
