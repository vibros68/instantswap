package ethplorer

import (
	"encoding/json"
	"fmt"
)

func parse(r []byte, obj interface{}) error {
	var ethErr ethError
	json.Unmarshal(r, &ethErr)
	if ethErr.Error.Code > 0 && len(ethErr.Error.Message) > 0 {
		return fmt.Errorf(ethErr.Error.Message)
	}
	return json.Unmarshal(r, obj)
}

type ethError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
