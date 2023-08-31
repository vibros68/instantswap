package sideshift

import (
	"encoding/json"
	"time"
)

func parseResponseData(r []byte, obj interface{}) error {
	return json.Unmarshal(r, obj)
}

type Currency struct {
	Coin           string      `json:"coin"`
	Networks       []string    `json:"networks"`
	Name           string      `json:"name"`
	HasMemo        bool        `json:"hasMemo"`
	FixedOnly      interface{} `json:"fixedOnly"`
	VariableOnly   interface{} `json:"variableOnly"`
	DepositOffline interface{} `json:"depositOffline"`
	SettleOffline  interface{} `json:"settleOffline"`
}

type ExchangeRateRequest struct {
	DepositCoin    string `json:"depositCoin"`
	DepositNetwork string `json:"depositNetwork"`
	SettleCoin     string `json:"settleCoin"`
	SettleNetwork  string `json:"settleNetwork"`
	DepositAmount  string `json:"depositAmount,omitempty"`
	SettleAmount   string `json:"settleAmount,omitempty"`
	AffiliateId    string `json:"affiliateId"`
	CommissionRate string `json:"commissionRate,omitempty"`
}

type Quote struct {
	Id             string    `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	DepositCoin    string    `json:"depositCoin"`
	SettleCoin     string    `json:"settleCoin"`
	DepositNetwork string    `json:"depositNetwork"`
	SettleNetwork  string    `json:"settleNetwork"`
	ExpiresAt      time.Time `json:"expiresAt"`
	DepositAmount  string    `json:"depositAmount"`
	SettleAmount   string    `json:"settleAmount"`
	Rate           string    `json:"rate"`
	AffiliateId    string    `json:"affiliateId"`
}

type PairResponse struct {
	Min            string `json:"min"`
	Max            string `json:"max"`
	Rate           string `json:"rate"`
	DepositCoin    string `json:"depositCoin"`
	SettleCoin     string `json:"settleCoin"`
	DepositNetwork string `json:"depositNetwork"`
	SettleNetwork  string `json:"settleNetwork"`
}

type createFixedShift struct {
	SettleAddress string `json:"settleAddress"`
	AffiliateId   string `json:"affiliateId"`
	QuoteId       string `json:"quoteId"`
	RefundAddress string `json:"refundAddress"`
}

type FixedShift struct {
	Id             string    `json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	DepositCoin    string    `json:"depositCoin"`
	SettleCoin     string    `json:"settleCoin"`
	DepositNetwork string    `json:"depositNetwork"`
	SettleNetwork  string    `json:"settleNetwork"`
	DepositAddress string    `json:"depositAddress"`
	SettleAddress  string    `json:"settleAddress"`
	DepositMin     string    `json:"depositMin"`
	DepositMax     string    `json:"depositMax"`
	RefundAddress  string    `json:"refundAddress"`
	Type           string    `json:"type"`
	QuoteId        string    `json:"quoteId"`
	DepositAmount  string    `json:"depositAmount"`
	SettleAmount   string    `json:"settleAmount"`
	ExpiresAt      time.Time `json:"expiresAt"`
	Status         string    `json:"status"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Rate           string    `json:"rate"`
	// information of get request
	DepositHash       string    `json:"depositHash"`
	SettleHash        string    `json:"settleHash"`
	DepositReceivedAt time.Time `json:"depositReceivedAt"`
}
