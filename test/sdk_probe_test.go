package test

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"vsc-node/lib/test_utils"
	ledgerDb "vsc-node/modules/db/vsc/ledger"
	stateEngine "vsc-node/modules/state-processing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"
)

// TestSdkProbe runs the contract's sdkprobe entrypoint, which calls the host
// functions added to the SDK after the original template: crypto.keccak256,
// system.verify_address, the db uint64 helpers, and ephem_db. A missing host
// function would stop the module instantiating at all, and a mis-marshalled
// argument would come back wrong — so this is the proof the wrappers work
// against the node pinned in go.mod, not just that they compile.
func TestSdkProbe(t *testing.T) {
	testCode, err := os.ReadFile("../artifacts/main.wasm")
	assert.NoError(t, err)

	ct := test_utils.NewContractTest()
	ct.Deposit("hive:someone", 5000, ledgerDb.AssetHive)
	ct.RegisterContract("vscmycontract", "hive:someowner", testCode)

	callResult := ct.Call(stateEngine.TxVscCallContract{
		Self: stateEngine.TxSelf{
			TxId:                 "probetxid",
			BlockId:              "abcdef",
			Index:                0,
			OpIndex:              0,
			Timestamp:            "2025-09-03T00:00:00",
			RequiredAuths:        []string{"hive:someone"},
			RequiredPostingAuths: []string{},
		},
		ContractId: "vscmycontract",
		Action:     "sdkprobe",
		Payload:    json.RawMessage([]byte(`""`)),
		RcLimit:    100000,
	})

	assert.True(t, callResult.Success, "sdkprobe call failed: %s", callResult.Ret)

	parts := strings.Split(callResult.Ret, "|")
	if !assert.Len(t, parts, 4, "unexpected probe output: %q", callResult.Ret) {
		return
	}

	// crypto.keccak256 — cross-checked against a local Keccak-256, not a
	// hardcoded digest, so the test still means something if the vector changes.
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte("abc"))
	assert.Equal(t, hex.EncodeToString(h.Sum(nil)), parts[0], "crypto.keccak256 mismatch")

	// system.verify_address
	assert.Equal(t, "user:hive", parts[1], "system.verify_address mismatch")

	// db.set_object/get_object through StateSetU64/StateGetU64
	assert.Equal(t, "4096", parts[2], "StateSetU64/StateGetU64 round-trip mismatch")

	// ephem_db.set_object/get_object
	assert.Equal(t, "ok", parts[3], "ephem_db round-trip mismatch")
}
