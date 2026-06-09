package sdk

import "encoding/hex"

// ===================================
// Crypto — pure cryptographic operations run on the host, not in WASM.
// ===================================

// Keccak256 computes the Keccak-256 hash of data via the host runtime.
// Input and output are raw bytes (the host marshalling is hex-encoded for you).
func Keccak256(data []byte) []byte {
	input := hex.EncodeToString(data)
	out := cryptoKeccak256(&input)
	if out == nil {
		return nil
	}
	decoded, _ := hex.DecodeString(*out)
	return decoded
}

// Ecrecover recovers the signer's 20-byte Ethereum address from a hash and
// signature. hash must be 32 bytes; sig must be 65 bytes (r[32] + s[32] + v[1]).
// Returns the recovered address, or an error if recovery fails.
func Ecrecover(hash []byte, sig []byte) ([]byte, error) {
	hashHex := hex.EncodeToString(hash)
	sigHex := hex.EncodeToString(sig)
	out := cryptoEcrecover(&hashHex, &sigHex)
	if out == nil {
		return nil, nil
	}
	return hex.DecodeString(*out)
}

// RlpDecode decodes RLP-encoded data via the host runtime, returning the JSON
// string representation of the decoded items.
func RlpDecode(data []byte) string {
	input := hex.EncodeToString(data)
	out := cryptoRlpDecode(&input)
	if out == nil {
		return ""
	}
	return *out
}

// Sp1VerifyGroth16 verifies an SP1 Groth16 proof on the host. All arguments are
// hex-encoded strings (proof, public inputs, the SP1 vkey hash, the Groth16
// verifying key, and the vk root). Returns true iff the proof verifies.
func Sp1VerifyGroth16(proof, publicInputs, sp1VkeyHash, groth16Vk, vkRoot string) bool {
	out := cryptoSp1VerifyGroth16(&proof, &publicInputs, &sp1VkeyHash, &groth16Vk, &vkRoot)
	return out != nil && *out == "true"
}
