package exolix

import "time"

// struct for get currencies
type CurrencyResponse struct {
	Data  []Currency `json:"data"`
	Count int        `json:"count"`
}

type Currency struct {
	Code     string    `json:"code"`
	Name     string    `json:"name"`
	Icon     string    `json:"icon"`
	Notes    string    `json:"notes"`
	Networks []Network `json:"networks"`
}

type Network struct {
	CoinNetworkID int     `json:"coinNetworkId"`
	Network       string  `json:"network"`
	Name          string  `json:"name"`
	ShortName     string  `json:"shortName"`
	Notes         string  `json:"notes"`
	AddressRegex  *string `json:"addressRegex"`
	IsDefault     bool    `json:"isDefault"`
	BlockExplorer *string `json:"blockExplorer"`
	MemoNeeded    bool    `json:"memoNeeded"`
	MemoName      *string `json:"memoName"`
	MemoRegex     *string `json:"memoRegex"`
	Precision     int     `json:"precision"`
	Contract      *string `json:"contract"`
	Icon          *string `json:"icon"`
}

// struct for get exchange rate
type RateResponse struct {
	FromAmount  float64 `json:"fromAmount"`
	ToAmount    float64 `json:"toAmount"`
	Rate        float64 `json:"rate"`
	Message     *string `json:"message"`
	MinAmount   float64 `json:"minAmount"`
	WithdrawMin float64 `json:"withdrawMin"`
	MaxAmount   float64 `json:"maxAmount"`
}

// struct for create order request
type OrderRequest struct {
	CoinFrom          string  `json:"coinFrom"`
	CoinTo            string  `json:"coinTo"`
	NetworkFrom       string  `json:"networkFrom"`
	NetworkTo         string  `json:"networkTo"`
	Amount            float64 `json:"amount"`
	WithdrawalAmount  float64 `json:"withdrawalAmount,omitempty"`
	WithdrawalAddress string  `json:"withdrawalAddress"`
	WithdrawalExtraId string  `json:"withdrawalExtraId,omitempty"`
	RateType          string  `json:"rateType,omitempty"`
	RefundAddress     string  `json:"refundAddress,omitempty"`
	RefundExtraId     string  `json:"refundExtraId,omitempty"`
	Slippage          float64 `json:"slippage,omitempty"`
}

// struct for order info
type Order struct {
	Id                string       `json:"id"`
	Amount            float64      `json:"amount"`
	AmountTo          float64      `json:"amountTo"`
	CoinFrom          CoinContract `json:"coinFrom"`
	CoinTo            CoinContract `json:"coinTo"`
	Comment           *string      `json:"comment"`
	CreatedAt         time.Time    `json:"createdAt"`
	DepositAddress    string       `json:"depositAddress"`
	DepositExtraId    *string      `json:"depositExtraId"`
	WithdrawalAddress string       `json:"withdrawalAddress"`
	WithdrawalExtraId string       `json:"withdrawalExtraId"`
	HashIn            HashLink     `json:"hashIn"`
	HashOut           HashLink     `json:"hashOut"`
	Rate              float64      `json:"rate"`
	RateType          string       `json:"rateType"`
	RefundAddress     *string      `json:"refundAddress"`
	RefundExtraId     *string      `json:"refundExtraId"`
	Status            string       `json:"status"`
}

type CoinContract struct {
	CoinCode         string  `json:"coinCode"`
	CoinName         string  `json:"coinName"`
	Network          string  `json:"network"`
	NetworkName      string  `json:"networkName"`
	NetworkShortName *string `json:"networkShortName"`
	Icon             string  `json:"icon"`
	MemoName         *string `json:"memoName"`
	Contract         *string `json:"contract"`
}

type HashLink struct {
	Hash *string `json:"hash"`
	Link *string `json:"link"`
}
