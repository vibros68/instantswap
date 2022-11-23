package easybit

import (
	"encoding/json"
	"fmt"
)

type general struct {
	Success      int             `json:"success"`
	ErrorMessage string          `json:"errorMessage"`
	ErrorCode    int             `json:"errorCode"`
	Data         json.RawMessage `json:"data"`
}

func parseDataResponse(r []byte, obj interface{}) error {
	var res general
	err := json.Unmarshal(r, &res)
	if err != nil {
		return err
	}
	if res.Success == 0 {
		return fmt.Errorf("error[%d]: %s", res.ErrorCode, res.ErrorMessage)
	}
	return json.Unmarshal(res.Data, obj)
}

type Currency struct {
	Currency         string    `json:"currency"`
	Name             string    `json:"name"`
	SendStatusAll    bool      `json:"sendStatusAll"`
	ReceiveStatusAll bool      `json:"receiveStatusAll"`
	NetworkList      []Network `json:"networkList"`
}

type Network struct {
	Network              string      `json:"network"`
	Name                 string      `json:"name"`
	IsDefault            bool        `json:"isDefault"`
	SendStatus           bool        `json:"sendStatus"`
	ReceiveStatus        bool        `json:"receiveStatus"`
	HasTag               bool        `json:"hasTag"`
	TagName              interface{} `json:"tagName"`
	ReceiveDecimals      int         `json:"receiveDecimals"`
	ConfirmationsMinimum int         `json:"confirmationsMinimum"`
	ConfirmationsMaximum int         `json:"confirmationsMaximum"`
	Explorer             string      `json:"explorer"`
	ExplorerHash         string      `json:"explorerHash"`
	ExplorerAddress      string      `json:"explorerAddress"`
}

type ExchangeRate struct {
	Rate           string `json:"rate"`
	SendAmount     string `json:"sendAmount"`
	ReceiveAmount  string `json:"receiveAmount"`
	NetworkFee     string `json:"networkFee"`
	Confirmations  int    `json:"confirmations"`
	ProcessingTime string `json:"processingTime"`
}

type PairInfo struct {
	MinimumAmount  string `json:"minimumAmount"`
	MaximumAmount  string `json:"maximumAmount"`
	NetworkFee     string `json:"networkFee"`
	Confirmations  int    `json:"confirmations"`
	ProcessingTime string `json:"processingTime"`
}

type Order struct {
	Id             string      `json:"id"`
	Send           string      `json:"send"`
	Receive        string      `json:"receive"`
	SendNetwork    string      `json:"sendNetwork"`
	ReceiveNetwork string      `json:"receiveNetwork"`
	SendAmount     string      `json:"sendAmount"`
	ReceiveAmount  string      `json:"receiveAmount"`
	SendAddress    string      `json:"sendAddress"`
	SendTag        interface{} `json:"sendTag"`
	ReceiveAddress string      `json:"receiveAddress"`
	ReceiveTag     interface{} `json:"receiveTag"`
	RefundAddress  interface{} `json:"refundAddress"`
	RefundTag      interface{} `json:"refundTag"`
	Vpm            string      `json:"vpm"`
	CreatedAt      int64       `json:"createdAt"`
	Status         string      `json:"status"`
	HashIn         interface{} `json:"hashIn"`
	HashOut        interface{} `json:"hashOut"`
	NetworkFee     string      `json:"networkFee"`
	Earned         string      `json:"earned"`
	UpdatedAt      int64       `json:"updatedAt"`
}
