package godex

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type InfoRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type RevertResponse struct {
	MinAmount float64 `json:"min_amount"`
	MaxAmount int     `json:"max_amount"`
	Amount    float64 `json:"amount"`
	Fee       int     `json:"fee"`
	Rate      float64 `json:"rate"`
}

func parseResponseData(data []byte, obj interface{}) error {
	var godexErr Error
	err := json.Unmarshal(data, &godexErr)
	if err != nil {
		return fmt.Errorf(string(data))
	}
	if err == nil && len(godexErr.Error) > 0 {
		return fmt.Errorf(godexErr.Error)
	}
	err = json.Unmarshal(data, obj)
	return err
}

type InfoResponse struct {
	MinAmount    json.Number `json:"min_amount,omitempty"`
	MaxAmount    json.Number `json:"max_amount,omitempty"`
	Amount       json.Number `json:"amount,omitempty"`
	Fee          json.Number `json:"fee,omitempty"`
	Rate         json.Number `json:"rate,omitempty"`
	NetworksFrom []Network   `json:"networks_from"`
	NetworksTo   []Network   `json:"networks_to"`
}

type Network struct {
	Network string `json:"network"`
	HasTag  int    `json:"has_tag"`
}

type TransactionReq struct {
	CoinFrom          string  `json:"coin_from"`
	CoinTo            string  `json:"coin_to"`
	DepositAmount     float64 `json:"deposit_amount"`
	Withdrawal        string  `json:"withdrawal"`
	WithdrawalExtraId string  `json:"withdrawal_extra_id"`
	Return            string  `json:"return"`
	ReturnExtraId     string  `json:"return_extra_id"`
	AffiliateId       string  `json:"affiliate_id"`
	CoinToNetwork     string  `json:"coin_to_network"`
	CoinFromNetwork   string  `json:"coin_from_network"`
}

type Transaction struct {
	Status               string      `json:"status"`
	CoinFrom             string      `json:"coin_from"`
	CoinTo               string      `json:"coin_to"`
	DepositAmount        json.Number `json:"deposit_amount"`
	Withdrawal           string      `json:"withdrawal"`
	WithdrawalExtraId    string      `json:"withdrawal_extra_id"`
	Return               string      `json:"return"`
	ReturnExtraId        string      `json:"return_extra_id"`
	WithdrawalAmount     json.Number `json:"withdrawal_amount"`
	Deposit              string      `json:"deposit"`
	DepositExtraId       string      `json:"deposit_extra_id"`
	Rate                 json.Number `json:"rate"`
	Fee                  json.Number `json:"fee"`
	TransactionId        string      `json:"transaction_id"`
	HashIn               string      `json:"hash_in"`
	HashOut              string      `json:"hash_out"`
	RealDepositAmount    json.Number `json:"real_deposit_amount"`
	RealWithdrawalAmount json.Number `json:"real_withdrawal_amount"`
}

type Currency struct {
	Code      string      `json:"code"`
	Name      string      `json:"name"`
	Disabled  int         `json:"disabled"`
	Icon      string      `json:"icon"`
	HasExtra  int         `json:"has_extra"`
	ExtraName interface{} `json:"extra_name"`
	Explorer  interface{} `json:"explorer"`
}
