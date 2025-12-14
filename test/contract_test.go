package test

import (
	"encoding/json"
	"os"
	"testing"

	"vsc-node/lib/test_utils"
	"vsc-node/modules/db/vsc/contracts"
	ledgerDb "vsc-node/modules/db/vsc/ledger"
	stateEngine "vsc-node/modules/state-processing"

	"github.com/stretchr/testify/assert"
)

func TestEntrypoint(t *testing.T) {
	testCode, err := os.ReadFile("../artifacts/main.wasm")
	assert.NoError(t, err)

	ct := test_utils.NewContractTest()
	ct.Deposit("hive:someone", 5000, ledgerDb.AssetHive)
	ct.Deposit("hive:someone", 5000, ledgerDb.AssetHbd)
	ct.RegisterContract("vscmycontract", "hive:someowner", testCode)

	txSelf := stateEngine.TxSelf{
		TxId:                 "sometxid",
		BlockId:              "abcdef",
		Index:                69,
		OpIndex:              0,
		Timestamp:            "2025-09-03T00:00:00",
		RequiredAuths:        []string{"did:key:someone"},
		RequiredPostingAuths: []string{},
	}

	callResult, _, _ := ct.Call(stateEngine.TxVscCallContract{
		Self:       txSelf,
		ContractId: "vscmycontract",
		Action:     "entrypoint",
		Payload:    json.RawMessage([]byte("abc")),
		RcLimit:    1000,
		Intents: []contracts.Intent{{
			Type: "transfer.allow",
			Args: map[string]string{
				"limit": "1.000",
				"token": "hive",
			},
		}},
	})
	assert.True(t, callResult.Success)
	assert.Equal(t, callResult.Ret, "abc")
}
