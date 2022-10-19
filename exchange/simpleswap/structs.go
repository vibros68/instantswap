package simpleswap

type CreateExchange struct {
	CurrencyFrom      string  `json:"currency_from"`
	CurrencyTo        string  `json:"currency_to"`
	Fixed             bool    `json:"fixed"`
	Amount            float64 `json:"amount"`
	AddressTo         string  `json:"address_to"`
	ExtraIdTo         string  `json:"extraIdTo"`
	UserRefundAddress string  `json:"userRefundAddress"`
	UserRefundExtraId string  `json:"userRefundExtraId"`
	Referral          string  `json:"referral"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Order struct {
	Id                string `json:"id"`
	Type              string `json:"type"`
	Timestamp         string `json:"timestamp"`
	UpdatedAt         string `json:"updated_at"`
	CurrencyFrom      string `json:"currency_from"`
	CurrencyTo        string `json:"currency_to"`
	AmountFrom        string `json:"amount_from"`
	ExpectedAmount    string `json:"expected_amount"`
	AmountTo          string `json:"amount_to"`
	AddressFrom       string `json:"address_from"`
	AddressTo         string `json:"address_to"`
	ExtraIdFrom       string `json:"extra_id_from"`
	ExtraIdTo         string `json:"extra_id_to"`
	UserRefundAddress string `json:"user_refund_address"`
	UserRefundExtraId string `json:"user_refund_extra_id"`
	TxFrom            string `json:"tx_from"`
	TxTo              string `json:"tx_to"`
	Status            string `json:"status"`
	Currencies        struct {
		CurrencyFromTicker string `json:"currency_from_ticker"`
		CurrencyToTicker   string `json:"currency_to_ticker"`
	} `json:"currencies"`
}
