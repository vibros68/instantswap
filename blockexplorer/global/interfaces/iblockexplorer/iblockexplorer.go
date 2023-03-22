package iblockexplorer

type IBlockExplorer interface {
	New(apiKey, apiSecret string, enableOutput bool) *Client
	GetTransaction(txid string) (tx ITransaction, err error)
	GetTxsForAddress(address string, limit int) (tx IRawAddrResponse, err error)

	//VerifyTransaction verifies transaction based on values passed in (params: txid, address (required), amount (required), createdAt(unix timestamp) )
	VerifyTransaction(vars interface{}) (tx ITransaction, err error)

	//PushTx pushes a raw tx hash
	PushTx(rawtxhash string) (result string, err error)
}
