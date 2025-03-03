package exchcx

type Error struct {
	Error string `json:"error"`
}

type Volume struct {
	Volume string `json:"volume"`
}

type Rate struct {
	NetworkFee struct {
		F string `json:"f"`
		M string `json:"m"`
		S string `json:"s"`
	} `json:"network_fee"`
	Rate     float64 `json:"rate,string"`
	RateMode string  `json:"rate_mode"`
	Reserve  any     `json:"reserve"`
	SvcFee   string  `json:"svc_fee"`
}

type Order struct {
	Created               int      `json:"created"`
	FromAddr              string   `json:"from_addr"`
	FromAmountReceived    *float64 `json:"from_amount_received"`
	FromCurrency          string   `json:"from_currency"`
	MaxInput              string   `json:"max_input"`
	MinInput              string   `json:"min_input"`
	NetworkFee            int      `json:"network_fee"`
	OrderId               string   `json:"orderid"`
	Rate                  float64  `json:"rate,string"`
	RateMode              string   `json:"rate_mode"`
	State                 string   `json:"state"`
	SvcFee                string   `json:"svc_fee"`
	ToAddress             string   `json:"to_address"`
	ToAmount              *float64 `json:"to_amount"`
	ToCurrency            string   `json:"to_currency"`
	TransactionIdReceived *string  `json:"transaction_id_received"`
	TransactionIdSent     *string  `json:"transaction_id_sent"`
}
