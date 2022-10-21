package evercoin

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
type EstimateAmount struct {
	DepositCoin       string  `json:"depositCoin"`
	DestinationCoin   string  `json:"destinationCoin"`
	DepositAmount     float64 `json:"depositAmount,string"`
	DestinationAmount float64 `json:"destinationAmount,string,omitempty"`
	Signature         string  `json:"signature,omitempty"`
}
type EstimateAmountResult struct {
	Error  json.RawMessage `json:"error"`
	Result EstimateAmount  `json:"result"`
}
type ActiveCurrs struct {
	Name       string
	Currencies []ActiveCurr
}
type ActiveCurr struct {
	FromAvailable bool   `json:"fromAvailable"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	TagName       string `json:"tagName"`
	ToAvailable   bool   `json:"toAvailable"`
	Value         string `json:"value"`
}
type ErrorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type QueryLimits struct {
	Error  json.RawMessage `json:"error"`
	Result struct {
		DepositCoin     string  `json:"depositCoin"`
		DestinationCoin string  `json:"destinationCoin"`
		MaxDeposit      float64 `json:"maxDeposit,string"`
		MinDeposit      float64 `json:"minDeposit,string"`
	} `json:"result"`
}

//CREATE
type Address struct {
	MainAddress string `json:"mainAddress"`
	TagValue    string `json:"tagValue"`
}

type CreateOrder struct {
	DepositAmount      float64 `json:"depositAmount"`
	DepositCoin        string  `json:"depositCoin"`
	DestinationAddress Address `json:"destinationAddress"`
	DestinationAmount  float64 `json:"destinationAmount"`
	DestinationCoin    string  `json:"destinationCoin"`
	RefundAddress      Address `json:"refundAddress"`
	Signature          string  `json:"signature"`
}

type CreateResultInfo struct {
	DepositAddress Address `json:"depositAddress"`
	UUID           string  `json:"orderId"`
}
type CreateResult struct {
	Error   json.RawMessage  `json:"error"`
	Expires int              `json:"expires"`
	Order   CreateResultInfo `json:"result"`
}

//INFO
type OrderInfoResultInfo struct {
	CreationTime              int     `json:"creationTime"`
	DepositAddress            Address `json:"depositAddress"`
	DepositAmount             float64 `json:"depositAmount"`
	DepositCoin               string  `json:"depositCoin"`
	DepositExpectedAmount     float64 `json:"depositExpectedAmount"`
	DestinationAddress        Address `json:"destinationAddress"`
	DestinationAmount         float64 `json:"destinationAmount"`
	DestinationCoin           string  `json:"destinationCoin"`
	DestinationExpectedAmount float64 `json:"destinationExpectedAmount"`
	ExchangeStatus            int     `json:"exchangeStatus"`
	MaxDeposit                float64 `json:"maxDeposit"`
	MinDeposit                float64 `json:"minDeposit"`
	RefundAddress             Address `json:"refundAddress"`
}
type OrderInfoResult struct {
	Error json.RawMessage     `json:"error"`
	Order OrderInfoResultInfo `json:"result"`
}
type CoinInfo struct {
	FromAvailable bool   `json:"fromAvailable"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	TagName       string `json:"tagName"`
	ToAvailable   bool   `json:"toAvailable"`
	Value         string `json:"value"`
}
type Coins struct {
	Error json.RawMessage `json:"error"`
	Coins []CoinInfo      `json:"result"`
}
type UUID struct {
	UUID string `json:"uuid"`
}
