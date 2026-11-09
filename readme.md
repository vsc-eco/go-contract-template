# VSC Golang Contracts Template



### Structure

```golang 
./contract-template
├── artifacts/  //Contains 
├── contract/
│   └── main.go //This is where your contract code will go
├── deploy.sh
├── mock_test.sh //Test your golang smart contract 
├── readme.md
├── runtime/
│   └── gc_leaking_exported.go
├── scripts/
│   └── sync-sdk.sh //Copy sdk/ into a downstream contract repo
├── sdk/ //SDK implementation. Do NOT modify
│   └── sdk.go
└── testing/
```

### The SDK lives here

This repo holds the **primary** copy of the Magi contract SDK. Every contract
repo (`magi_nft-contract`, `magi_token-contract`, ...) vendors a
**byte-identical** `sdk/` directory, so a host function added here reaches all
of them by copying, never by hand-editing a downstream copy:

```bash
scripts/sync-sdk.sh /path/to/magi_nft-contract           # copy + verify
scripts/sync-sdk.sh --check /path/to/magi_nft-contract   # verify only (CI)
```

That works because `sdk/` imports nothing from its own module: the custom-GC
`runtime` package is blank-imported by the contract's `main` package
(`contract/main.go`), not by the SDK.

The SDK mirrors the host functions the node exports from
`modules/wasm/sdk/sdk.go` (`SdkNamespaces`). A contract can only call what the
node it runs on actually provides, and what this repo's tests can exercise is
set by the `replace vsc-node` pin in `go.mod` — if a wrapper's host function is
missing there, the module fails to instantiate rather than returning an error.
`test/sdk_probe_test.go` runs the newer host functions end to end against that
pinned node.