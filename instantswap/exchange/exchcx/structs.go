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
	Rate     string `json:"rate"`
	RateMode string `json:"rate_mode"`
	Reserve  string `json:"reserve"`
	SvcFee   string `json:"svc_fee"`
}
