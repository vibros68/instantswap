package aptexplorer

import (
	"encoding/json"
	"strings"

	"code.cryptopower.dev/group/instantswap/blockexplorer"
	"code.cryptopower.dev/group/instantswap/blockexplorer/global/interfaces/idaemon"
	"code.cryptopower.dev/group/instantswap/blockexplorer/global/utils"
)

func parseResponseData(r []byte, obj interface{}) error {
	return json.Unmarshal(r, obj)
}

type T struct {
	Type string `json:"type"`
	Data struct {
		Coin struct {
			Value string `json:"value"`
		} `json:"coin"`
		DepositEvents struct {
			Counter string `json:"counter"`
			Guid    struct {
				Id struct {
					Addr        string `json:"addr"`
					CreationNum string `json:"creation_num"`
				} `json:"id"`
			} `json:"guid"`
		} `json:"deposit_events"`
		Frozen         bool `json:"frozen"`
		WithdrawEvents struct {
			Counter string    `json:"counter"`
			Guid    EventGuid `json:"guid"`
		} `json:"withdraw_events"`
	} `json:"data"`
}

type EventGuid struct {
	Id struct {
		Addr        string `json:"addr"`
		CreationNum string `json:"creation_num"`
	} `json:"id"`
}

type Event struct {
	Counter string    `json:"counter"`
	Guid    EventGuid `json:"guid"`
}

type Transaction struct {
	Version                 string      `json:"version"`
	Hash                    string      `json:"hash"`
	StateChangeHash         string      `json:"state_change_hash"`
	EventRootHash           string      `json:"event_root_hash"`
	StateCheckpointHash     interface{} `json:"state_checkpoint_hash"`
	GasUsed                 string      `json:"gas_used"`
	Success                 bool        `json:"success"`
	VmStatus                string      `json:"vm_status"`
	AccumulatorRootHash     string      `json:"accumulator_root_hash"`
	Changes                 []Change    `json:"changes"`
	Sender                  string      `json:"sender"`
	SequenceNumber          string      `json:"sequence_number"`
	MaxGasAmount            string      `json:"max_gas_amount"`
	GasUnitPrice            string      `json:"gas_unit_price"`
	ExpirationTimestampSecs string      `json:"expiration_timestamp_secs"`
	Payload                 *TxPayload  `json:"payload"`
	Signature               struct {
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
		Type      string `json:"type"`
	} `json:"signature"`
	Events    []TxEvent `json:"events"`
	Timestamp string    `json:"timestamp"`
	Type      string    `json:"type"`
}

type Change struct {
	Address      string `json:"address,omitempty"`
	StateKeyHash string `json:"state_key_hash"`
	Data         *Data  `json:"data"`
	Type         string `json:"type"`
	Handle       string `json:"handle,omitempty"`
	Key          string `json:"key,omitempty"`
	Value        string `json:"value,omitempty"`
}

type CapabilityOffer struct {
	For struct {
		Vec []interface{} `json:"vec"`
	} `json:"for"`
}

type Data struct {
	Type string `json:"type"`
	Data struct {
		Coin struct {
			Value string `json:"value"`
		} `json:"coin,omitempty"`
		DepositEvents           *Event           `json:"deposit_events,omitempty"`
		Frozen                  bool             `json:"frozen,omitempty"`
		WithdrawEvents          *Event           `json:"withdraw_events,omitempty"`
		AuthenticationKey       string           `json:"authentication_key,omitempty"`
		CoinRegisterEvents      *Event           `json:"coin_register_events,omitempty"`
		GuidCreationNum         string           `json:"guid_creation_num,omitempty"`
		KeyRotationEvents       *Event           `json:"key_rotation_events,omitempty"`
		RotationCapabilityOffer *CapabilityOffer `json:"rotation_capability_offer,omitempty"`
		SequenceNumber          string           `json:"sequence_number,omitempty"`
		SignerCapabilityOffer   *CapabilityOffer `json:"signer_capability_offer,omitempty"`
	} `json:"data"`
}

type TxEvent struct {
	Guid struct {
		CreationNumber string `json:"creation_number"`
		AccountAddress string `json:"account_address"`
	} `json:"guid"`
	SequenceNumber string `json:"sequence_number"`
	Type           string `json:"type"`
	Data           struct {
		Amount string `json:"amount"`
	} `json:"data"`
}

type TxPayload struct {
	Function      string   `json:"function"`
	TypeArguments []string `json:"type_arguments"`
	Arguments     []string `json:"arguments"`
	Type          string   `json:"type"`
}

type Blockchain struct {
	ChainId             int    `json:"chain_id"`
	Epoch               string `json:"epoch"`
	LedgerVersion       string `json:"ledger_version"`
	OldestLedgerVersion string `json:"oldest_ledger_version"`
	LedgerTimestamp     string `json:"ledger_timestamp"`
	NodeRole            string `json:"node_role"`
	OldestBlockHeight   string `json:"oldest_block_height"`
	BlockHeight         string `json:"block_height"`
	GitHash             string `json:"git_hash"`
}

type BlockInfo struct {
	BlockHeight    string `json:"block_height"`
	BlockHash      string `json:"block_hash"`
	BlockTimestamp string `json:"block_timestamp"`
	FirstVersion   string `json:"first_version"`
	LastVersion    string `json:"last_version"`
}

func (tx *Transaction) getInOutPuts() (vIns []blockexplorer.IVIN, vOuts []blockexplorer.IVOUT) {
	for _, event := range tx.Events {
		if getTypeName(event.Type) == "WithdrawEvent" {
			vIns = append(vIns, blockexplorer.IVIN{
				Script:      "",
				Sequence:    0,
				Witness:     "",
				TxID:        "",
				VOUT:        0,
				Tree:        0,
				AmountIn:    idaemon.Amount(utils.StrToInt(event.Data.Amount)),
				BlockIndex:  0,
				BlockHeight: 0,
			})
		}
		if getTypeName(event.Type) == "DepositEvent" {
			vOuts = append(vOuts, blockexplorer.IVOUT{
				Addresses: []string{
					event.Guid.AccountAddress,
				},
				AddrTag:     "",
				AddrTagLink: "",
				N:           0,
				Script:      "",
				Spent:       false,
				TxIndex:     0,
				Type:        "",
				Value:       idaemon.Amount(utils.StrToInt(event.Data.Amount)),
			})
		}
	}
	return
}

func getTypeName(t string) string {
	tPiece := strings.Split(t, "::")
	if len(tPiece) == 3 {
		return tPiece[2]
	}
	return ""
}
