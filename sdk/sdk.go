package sdk

import (
	"encoding/json"
	"strconv"
)

//go:wasmimport sdk console.log
func log(s *string) *string

func Log(s string) {
	log(&s)
}

//go:wasmimport sdk db.setObject
func stateSetObject(key *string, value *string) *string

//go:wasmimport sdk db.getObject
func stateGetObject(key *string) *string

//go:wasmimport sdk db.delObject
func stateDeleteObject(key *string) *string

//go:wasmimport sdk system.getEnv
func getEnv(arg *string) *string

//go:wasmimport sdk hive.getbalance
func getBalance(arg1 *string, arg2 *string) *string

//go:wasmimport sdk hive.draw
func hiveDraw(arg1 *string, arg2 *string) *string

//go:wasmimport sdk hive.transfer
func hiveTransfer(arg1 *string, arg2 *string, arg3 *string) *string

//go:wasmimport sdk hive.withdraw
func hiveWithdraw(arg1 *string, arg2 *string, arg3 *string) *string

// /TODO: this is not implemented yet
// /go:wasmimport sdk contracts.read
func contractRead(contractId *string, key *string) *string

// /TODO: this is not implemented yet
// /go:wasmimport sdk contracts.call
func contractCall(contractId *string, method *string, payload *string, options *string) *string

// var envMap = []string{
// 	"contract.id",
// 	"tx.origin",
// 	"tx.id",
// 	"tx.index",
// 	"tx.op_index",
// 	"block.id",
// 	"block.height",
// 	"block.timestamp",
// }

// Set a value by key in the contract state
func StateSetObject(key string, value string) {
	stateSetObject(&key, &value)
}

// Get a value by key from the contract state
func StateGetObject(key string) *string {
	return stateGetObject(&key)
}

// Delete or unset a value by key in the contract state
func StateDeleteObject(key string) {
	stateDeleteObject(&key)
}

// Get current execution environment variables
func GetEnv() Env {
	// Per-key lookups via system.getEnv
	get := func(key string) string { return *getEnv(&key) }
	env := Env{}
	env.ContractId = get("contract_id")
	env.TxId = get("anchor.id")
	env.BlockId = get("anchor.block")
	if v := get("anchor.height"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			env.BlockHeight = n
		}
	}
	env.Timestamp = get("anchor.timestamp")
	if v := get("anchor.tx_index"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			env.Index = uint64(n)
		}
	}
	if v := get("anchor.op_index"); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			env.OpIndex = uint64(n)
		}
	}
	// Sender/auths
	sender := get("msg.sender")
	var auths []string
	var posting []string
	if s := get("msg.required_auths"); s != "" {
		_ = json.Unmarshal([]byte(s), &auths)
	}
	if s := get("msg.required_posting_auths"); s != "" {
		_ = json.Unmarshal([]byte(s), &posting)
	}
	ra := make([]Address, 0, len(auths))
	for _, a := range auths {
		ra = append(ra, Address(a))
	}
	rpa := make([]Address, 0, len(posting))
	for _, a := range posting {
		rpa = append(rpa, Address(a))
	}
	env.Sender = Sender{Address: Address(sender), RequiredAuths: ra, RequiredPostingAuths: rpa}
	return env
}

// Get current execution environment variable by a key
// Deprecated: prefer GetEnv() which fetches keys individually
func GetEnvKey(key string) *string { return getEnv(&key) }

// Get balance of an account
func GetBalance(address Address, asset Asset) int64 {
	addr := address.String()
	as := asset.String()
	res := getBalance(&addr, &as)
	if res == nil {
		return 0
	}
	v, _ := strconv.ParseInt(*res, 10, 64)
	return v
}

// Transfer assets from caller account to the contract up to the limit specified in `intents`. The transaction must be signed using active authority for Hive accounts.
func HiveDraw(amount int64, asset Asset) {
	amt := strconv.FormatInt(amount, 10)
	as := asset.String()
	hiveDraw(&amt, &as)
}

// Transfer assets from the contract to another account.
func HiveTransfer(to Address, amount int64, asset Asset) {
	toaddr := to.String()
	amt := strconv.FormatInt(amount, 10)
	as := asset.String()
	hiveTransfer(&toaddr, &amt, &as)
}

// Unmap assets from the contract to a specified Hive account.
func HiveWithdraw(to Address, amount int64, asset Asset) {
	toaddr := to.String()
	amt := strconv.FormatInt(amount, 10)
	as := asset.String()
	hiveWithdraw(&toaddr, &amt, &as)
}
