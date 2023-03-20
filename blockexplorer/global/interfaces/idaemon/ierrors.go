package idaemon

import "fmt"

var (
	BumFeeMinedTxError = fmt.Errorf("transaction has been mined, or is conflicted with a mined transaction")
)
