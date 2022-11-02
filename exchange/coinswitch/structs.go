package coinswitch

import (
	"encoding/json"
	"fmt"
)

//base json structure
type jsonResponse struct {
	Code    string          `json:"code"`
	Message string          `json:"msg"`
	Result  json.RawMessage `json:"data"`
	Success bool            `json:"success"`
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
type QueryLimitsRequest struct {
	DepositCoin     string `json:"depositCoin"`
	DestinationCoin string `json:"destinationCoin"`
}
type QueryLimitsResponse struct {
	Max float64 `json:"depositCoinMaxAmount"`
	Min float64 `json:"depositCoinMinAmount"`
}
type EstimateRequest struct {
	DepositCoin     string  `json:"depositCoin"`
	DestinationCoin string  `json:"destinationCoin"`
	DepositAmount   float64 `json:"depositCoinAmount"`
}

type Address struct {
	Address string `json:"address"`
	Tag     string `json:"tag"`
}
type CreateOrder struct {
	DepositCoin        string  `json:"depositCoin"`
	DepositCoinAmount  float64 `json:"depositCoinAmount"`
	DestinationAddress Address `json:"destinationAddress"`
	DestinationCoin    string  `json:"destinationCoin"`
	OfferReferenceID   string  `json:"offerReferenceId"`
	RefundAddress      Address `json:"refundAddress"`
	UserReferenceID    string  `json:"userReferenceId"`
}
type CreateResult struct {
	ExchangeAddress struct {
		Address string `json:"address"`
		Tag     string `json:"tag"`
	} `json:"exchangeAddress"`
	OrderID string `json:"orderId"`
}

type UUID struct {
	UUID   string `json:"id"`
	APIKEY string `json:"api_key"`
}

type OrderInfoResult struct {
	CreatedAt             int     `json:"createdAt"`
	DepositCoin           string  `json:"depositCoin"`
	DepositCoinAmount     float64 `json:"depositCoinAmount"`
	DestinationAddress    Address `json:"destinationAddress"`
	DestinationCoin       string  `json:"destinationCoin"`
	DestinationCoinAmount float64 `json:"destinationCoinAmount"`
	ExchangeAddress       Address `json:"exchangeAddress"`
	InputTransactionHash  string  `json:"inputTransactionHash"`
	OrderID               string  `json:"orderId"`
	OutputTransactionHash string  `json:"outputTransactionHash"`
	Status                string  `json:"status"`
	ValidTill             int     `json:"validTill,string"`
}

type EstimateAmount struct {
	DepositCoin           string  `json:"depositCoin"`
	DepositCoinAmount     float64 `json:"depositCoinAmount"`
	DestinationCoin       string  `json:"destinationCoin"`
	DestinationCoinAmount float64 `json:"destinationCoinAmount"`
	OfferReferenceID      string  `json:"offerReferenceId"`
}

type Currency struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

func parseResponseData(res []byte, obj interface{}) error {
	var response jsonResponse
	if err := json.Unmarshal(res, &response); err != nil {
		return err
	}
	if !response.Success {
		return fmt.Errorf("%s:error:%s:%s", LIBNAME, response.Code, response.Message)
	}
	return json.Unmarshal(response.Result, obj)
}
