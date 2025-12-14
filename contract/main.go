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
	_ "contract-template/sdk" // ensure sdk is imported

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
