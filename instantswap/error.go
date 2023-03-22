package instantswap

import "fmt"

var (
	TooManyRequestsError = fmt.Errorf("exchangeclient:error:429 Too Many Requests")
)
