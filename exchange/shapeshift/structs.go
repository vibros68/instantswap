package shapeshift

type MarketInfoResponse struct {
	Pair      string  `json:"pair"`
	Limit     float64 `json:"limit"`
	Rate      float64 `json:"rate"`
	Min       float64 `json:"minimum"`
	MinerFee  float64 `json:"minerFee"`
	ErrorCode string  `json:"errorCode"`
	ErrorMsg  string  `json:"error"`
}

type RateResponse struct {
	Pair      string  `json:"pair"`
	Rate      float64 `json:"rate,string"`
	ErrorCode string  `json:"errorCode"`
	ErrorMsg  string  `json:"error"`
}
type SendAmountRequest struct {
	Amount float64 `json:"amount"`
	Pair   string  `json:"pair"`
}
type SendAmountResponse struct {
	Pair             string  `json:"pair"`
	WithdrawalAmount float64 `json:"withdrawalAmount"`
	DepositAmount    float64 `json:"depositAmount"`
	Expiration       string  `json:"expiration"`
	QuotedRate       float64 `json:"quotedRate"`
	MinerFee         float64 `json:"minerFee"`
	ErrorCode        string  `json:"errorCode"`
	ErrorMsg         string  `json:"error"`
}
type OrderStatusResponse struct {
	Status          string  `json:"status"`
	Address         string  `json:"address"`
	ToAddress       string  `json:"withdraw"`
	AmountDeposited float64 `json:"incomingCoin"`
	FromCurrency    string  `json:"incomingType"`
	AmountReceiving float64 `json:"outgoingCoin,string"`
	ToCurrency      string  `json:"outgoingType"`
	TxID            string  `json:"transaction"`
	Error           string  `json:"error"`
}

//CREATE
type CreateOrder struct {
	Pair              string  `json:"pair"` //from_to format
	ToCurrencyAddress string  `json:"withdrawal"`
	RefundAddress     string  `json:"returnAddress"`
	InvoicedAmount    float64 `json:"amount"`            //amount in "from" currency
	ExtraID           string  `json:"destTag,omitempty"` //optional for some coins
	APIKEY            string  `json:"apiKey"`
	ErrorCode         string  `json:"errorCode"`
	ErrorMsg          string  `json:"error"`
}
type CreateResult struct {
	UUID           string `json:"orderId"`
	CurrencyFrom   string `json:"depositType"`
	CurrencyTo     string `json:"withdrawalType"`
	DepositAddress string `json:"deposit"`
	ToAddress      string `json:"withdrawal"`
	RefundAddress  string `json:"returnAddress"`

	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"error"`
}
