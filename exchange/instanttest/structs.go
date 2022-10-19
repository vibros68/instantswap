package instanttest

type UpdateOrderInfo struct {
	Destination   string  `json:"destination"`
	OrderedAmount float64 `json:"ordered_amount,string"`
	RefundAddress string  `json:"refund_address"`
	UUID          string  `json:"uuid"`
}
type UpdateOrder struct {
	Order UpdateOrderInfo `json:"order"`
}
