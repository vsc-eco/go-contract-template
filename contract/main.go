// Example Magi contract in Golang
//
// Build command: tinygo build -gc=custom -scheduler=none -panic=trap -no-debug -target=wasm-unknown -o artifacts/main.wasm ./contract
// Docker build command: docker run --rm -v $(pwd):/home/tinygo tinygo/tinygo:0.39.0 tinygo build -gc=custom -scheduler=none -panic=trap -no-debug -target=wasm-unknown -o artifacts/main.wasm ./contract
// Inspect Output: wasmer inspect main.wasm
//
// Caveats:
// - Go routines, channels, and defer are disabled
// - panic() always halts the program, since you can't recover in a deferred function call
// - must import sdk or build fails
// - to mark a function as a valid entrypoint, it must be manually exported (//go:wasmexport <entrypoint-name>)
//
// TODO:
// - when panic()ing, call `env.abort()` instead of executing the unreachable WASM instruction
// - Remove _initalize() export & double check not necessary

package main

import (
	"encoding/hex"
	"strconv"

	// the custom-GC runtime must be linked in by the main package (the sdk
	// package stays free of module-local imports so it can be copied verbatim)
	_ "contract-template/runtime"

	"contract-template/sdk"
)

func main() {

}

//go:wasmexport entrypoint
func Entrypoint(a *string) *string {
	sdk.Log(*a)
	// panic("test")
	return a
}

// sdkprobe exercises the host functions that were added to the SDK after the
// original template, so a `go test ./test/...` run proves the wrappers really
// reach the node (a wrapper whose host function is missing would fail to
// instantiate, and a mis-marshalled one returns the wrong bytes). It is a
// template demo, not a required entrypoint — delete it in a real contract.
//
//go:wasmexport sdkprobe
func SdkProbe(a *string) *string {
	// crypto.keccak256 — known-answer test on the empty-input vector
	hash := hex.EncodeToString(sdk.Keccak256([]byte("abc")))

	// system.verify_address — host-side address classification
	kind := sdk.VerifyAddress("hive:someone")

	// db.set_object / get_object through the uint64 helpers
	sdk.StateSetU64("probe:counter", 4096)
	counter := sdk.StateGetU64("probe:counter", 0)

	// ephem_db — per-block scratch state, dropped at the end of the block
	sdk.EphemStateSetObject("probe:ephem", "ok")
	ephem := ""
	if v := sdk.EphemStateGetObject(sdk.GetEnv().ContractId, "probe:ephem"); v != nil {
		ephem = *v
	}

	out := hash + "|" + kind + "|" + strconv.FormatUint(counter, 10) + "|" + ephem
	return &out
}
