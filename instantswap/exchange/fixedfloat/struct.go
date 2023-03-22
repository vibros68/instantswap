package fixedfloat

import (
	"encoding/json"
	"fmt"
)

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data json.RawMessage
}

func parseResponseData(r []byte, obj interface{}) error {
	var res response
	err := json.Unmarshal(r, &res)
	if err != nil {
		return err
	}
	if err == nil && res.Code > 0 {
		return fmt.Errorf(res.Msg)
	}
	err = json.Unmarshal(res.Data, obj)
	return err
}

type Currency struct {
	Currency  string `json:"currency"`
	Symbol    string `json:"symbol"`
	Network   string `json:"network"`
	Sub       string `json:"sub"`
	Name      string `json:"name"`
	Alias     string `json:"alias"`
	Type      string `json:"type"`
	Precision string `json:"precision"`
	Send      string `json:"send"`
	Recv      string `json:"recv"`
}
