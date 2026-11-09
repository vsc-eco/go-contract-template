package sdk

import (
	jlexer "github.com/CosmWasm/tinyjson/jlexer"
)

// TryResult is the structured outcome of a TryContractCall.
type TryResult struct {
	Ok        bool   // true when the callee returned normally
	Result    string // the callee's return value (when Ok)
	ErrorCode string // the callee's error code (when !Ok)
	Error     string // the callee's error message (when !Ok)
}

// UnmarshalTinyJSON decodes the host's try-call outcome envelope.
func (out *TryResult) UnmarshalTinyJSON(in *jlexer.Lexer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "ok":
			out.Ok = bool(in.Bool())
		case "result":
			out.Result = string(in.String())
		case "error_code":
			out.ErrorCode = string(in.String())
		case "error":
			out.Error = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}

// TryContractCall calls another contract in try/catch mode. If the callee
// reverts, the caller is NOT trapped: TryResult.Ok is false (with the callee's
// error), and every state/ledger write the callee made is rolled back to a
// savepoint taken just before it. On success Ok is true and Result holds the
// callee's return value. RC and gas are charged either way — the callee really
// executed, so a caught call is never a free probe.
//
// Requires chain consensus version >= 0.2.0. Below that the Try flag is ignored
// and a reverting callee traps the caller exactly like ContractCall.
func TryContractCall(contractId string, method string, payload string, options *ContractCallOptions) TryResult {
	o := ContractCallOptions{}
	if options != nil {
		o = *options
	}
	o.Try = true
	raw := ContractCall(contractId, method, payload, &o)
	var res TryResult
	if raw == nil {
		return res
	}
	l := &jlexer.Lexer{Data: []byte(*raw)}
	res.UnmarshalTinyJSON(l)
	return res
}
