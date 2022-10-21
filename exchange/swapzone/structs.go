package swapzone

import "time"

type SwapzoneError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type ExchangeRate struct {
	Adapter     string    `json:"adapter"`
	From        string    `json:"from"`
	FromNetwork string    `json:"fromNetwork"`
	To          string    `json:"to"`
	ToNetwork   string    `json:"toNetwork"`
	AmountFrom  float64   `json:"amountFrom"`
	AmountTo    float64   `json:"amountTo"`
	MinAmount   float64   `json:"minAmount"`
	MaxAmount   float64   `json:"maxAmount"`
	QuotaId     string    `json:"quotaId"`
	ValidUntil  time.Time `json:"validUntil"`
}

type Order struct {
	Id              string    `json:"id"`
	QuotaId         string    `json:"quotaId"`
	From            string    `json:"from"`
	FromNetwork     string    `json:"fromNetwork"`
	ToNetwork       string    `json:"toNetwork"`
	To              string    `json:"to"`
	Status          string    `json:"status"`
	AddressReceive  string    `json:"addressReceive"`
	ExtraIdReceive  string    `json:"extraIdReceive"`
	AddressDeposit  string    `json:"addressDeposit"`
	AmountDeposit   string    `json:"amountDeposit"`
	AmountEstimated string    `json:"amountEstimated"`
	CreatedAt       time.Time `json:"createdAt"`
	RefundExtraId   string    `json:"refundExtraId"`
	RefundAddress   string    `json:"refundAddress"`
}

type Transaction struct {
	Transaction Order `json:"transaction"`
}
