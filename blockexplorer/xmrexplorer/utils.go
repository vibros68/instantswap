package xmrexplorer

import (
	"encoding/json"
	"fmt"
)

const xmrErrorStatus = "fail"
const xmrSuccessStatus = "success"

func parseMoneroResponseData(data []byte, v interface{}) error {
	var xmrErr = XmrChainError{}
	var res = Response{
		Data:   &xmrErr,
		Status: "",
	}
	err := json.Unmarshal(data, &res)
	if err != nil {
		return fmt.Errorf("[%s] error: %v", LIBNAME, err)
	}
	if res.Status == xmrErrorStatus {
		return fmt.Errorf("[%s] error: %s", LIBNAME, xmrErr.Title)
	}
	res.Data = v
	err = json.Unmarshal(data, &res)
	if err != nil {
		return fmt.Errorf("[%s] error: %v", LIBNAME, err)
	}
	return nil
}
