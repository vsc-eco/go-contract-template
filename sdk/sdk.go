package sdk

import (
	_ "contract-template/runtime"
	"encoding/hex"
	"strconv"

	"github.com/CosmWasm/tinyjson"
)

// ===================================
// Event logs
// ===================================

// Emit a log to the contract output
func Log(s string) {
	log(&s)
}

// ===================================
// Abort/Revert Execution
// ===================================

// Aborts the contract execution
func Abort(msg string) {
	ln := int32(0)
	abort(&msg, nil, &ln, &ln)
	panic(msg)
}

// Reverts the transaction and abort execution in the same way as Abort().
func Revert(msg string, symbol string) {
	revert(&msg, &symbol)
}

// ===================================
// Persistent State
// ===================================

// Set a value by key in the contract state
func StateSetObject(key string, value string) {
	stateSetObject(&key, &value)
}

// Set a uint64 value by key in the contract state
func StateSetU64(key string, value uint64) {
	StateSetObject(key, string(u64ToBytes(value)))
}

// Get a value by key from the contract state
func StateGetObject(key string) *string {
	return stateGetObject(&key)
}

// Get a uint64 value by key from the contract state
func StateGetU64(key string, value uint64) uint64 {
	val := StateGetObject(key)
	if val == nil || *val == "" {
		return 0
	}
	return bytesToU64([]byte(*val))
}

// Delete or unset a value by key in the contract state
func StateDeleteObject(key string) {
	stateDeleteObject(&key)
}

// ===================================
// Ephemeral State
// ===================================

// Set a value by key in the ephemeral contract state
func EphemStateSetObject(key string, value string) {
	ephemStateSetObject(&key, &value)
}

// Get a value by key from the ephemeral contract state
func EphemStateGetObject(contractId string, key string) *string {
	return ephemStateGetObject(&contractId, &key)
}

// Delete or unset a value by key in the ephemeral contract state
func EphemStateDeleteObject(key string) {
	ephemStateDeleteObject(&key)
}

// ===================================
// Contract Env
// ===================================

// Get current execution environment variables
func GetEnv() Env {
	envStr := *getEnv(nil)
	env := Env{}
	tinyjson.Unmarshal([]byte(envStr), &env)
	env2 := Env2{}
	tinyjson.Unmarshal([]byte(envStr), &env2)

	requiredAuths := make([]Address, 0)
	for _, addr := range env2.Auths {
		requiredAuths = append(requiredAuths, Address(addr))
	}
	requiredPostingAuths := make([]Address, 0)
	for _, addr := range env2.PostingAuths {
		requiredPostingAuths = append(requiredPostingAuths, Address(addr))
	}

	env.Sender = Sender{
		Address:              Address(env2.Sender),
		RequiredAuths:        requiredAuths,
		RequiredPostingAuths: requiredPostingAuths,
	}

	return env
}

// Get current execution environment variables as json string
func GetEnvStr() string {
	return *getEnv(nil)
}

// Get current execution environment variable by a key
func GetEnvKey(key string) *string {
	return getEnvKey(&key)
}

// ===================================
// Ledgers
// ===================================

// Get balance of an account
func GetBalance(address Address, asset Asset) int64 {
	addr := address.String()
	as := asset.String()
	balStr := *getBalance(&addr, &as)
	bal, err := strconv.ParseInt(balStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return bal
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

// ===================================
// Intercontract
// ===================================

// Get a value by key from the contract state of another contract
func ContractStateGet(contractId string, key string) *string {
	return contractRead(&contractId, &key)
}

// Call another contract
func ContractCall(contractId string, method string, payload string, options *ContractCallOptions) *string {
	optStr := ""
	if options != nil {
		optByte, err := tinyjson.Marshal(options)
		if err != nil {
			Revert("could not serialize options", "sdk_error")
		}
		optStr = string(optByte)
	}
	return contractCall(&contractId, &method, &payload, &optStr)
}

// ===================================
// TSS
// ===================================

// TssCreateKey creates a key with the maximum epoch lifespan of 365 epochs.
// Deprecated: use TssCreateKeyForEpochs to specify a lifespan.
func TssCreateKey(keyId string, algo string) string {
	if algo != "ecdsa" && algo != "eddsa" {
		Abort("algo must be ecdsa or eddsa")
	}

	return *tssCreateKey(&keyId, &algo)
}

// TssCreateKeyForEpochs creates a key that expires after the given number of epochs.
func TssCreateKeyForEpochs(keyId string, algo string, epochs uint64) string {
	if algo != "ecdsa" && algo != "eddsa" {
		Abort("algo must be ecdsa or eddsa")
	}
	epochsStr := strconv.FormatUint(epochs, 10)
	return *tssV2CreateKey(&keyId, &algo, &epochsStr)
}

// TssRenewKey extends the expiry of a key owned by this contract by additionalEpochs.
// The key must be active or deprecated. Returns "renewed" on success.
func TssRenewKey(keyId string, additionalEpochs uint64) string {
	epochsStr := strconv.FormatUint(additionalEpochs, 10)
	return *tssRenewKey(&keyId, &epochsStr)
}

// Get details of a TSS key. Returns a comma-separated string consists of status, public key and algo.
func TssGetKey(keyId string) string {
	return *tssGetKey(&keyId)
}

// Request a digest to be signed by the TSS key.
func TssSignKey(keyId string, bytes []byte) string {
	byteStr := hex.EncodeToString(bytes)
	return *tssSignKey(&keyId, &byteStr)
}
