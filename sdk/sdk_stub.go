package sdk

//go:wasmimport sdk console.log
func log(s *string) *string

//go:wasmimport sdk db.set_object
func stateSetObject(key *string, value *string) *string

//go:wasmimport sdk db.get_object
func stateGetObject(key *string) *string

//go:wasmimport sdk db.rm_object
func stateDeleteObject(key *string) *string

//go:wasmimport sdk ephem_db.set_object
func ephemStateSetObject(key *string, value *string) *string

//go:wasmimport sdk ephem_db.get_object
func ephemStateGetObject(contractId *string, key *string) *string

//go:wasmimport sdk ephem_db.rm_object
func ephemStateDeleteObject(key *string) *string

//go:wasmimport sdk system.get_env
func getEnv(arg *string) *string

//go:wasmimport sdk system.get_env_key
func getEnvKey(arg *string) *string

//go:wasmimport sdk system.verify_address
func verifyAddress(addr *string) *string

//go:wasmimport sdk system.call
func systemCall(name *string, args *string) *string

//go:wasmimport sdk system.pendulum_apply_swap_fees
func pendulumApplySwapFees(arg *string) *string

//go:wasmimport sdk hive.get_balance
func getBalance(arg1 *string, arg2 *string) *string

//go:wasmimport sdk hive.draw
func hiveDraw(arg1 *string, arg2 *string) *string

//go:wasmimport sdk hive.transfer
func hiveTransfer(arg1 *string, arg2 *string, arg3 *string) *string

//go:wasmimport sdk hive.withdraw
func hiveWithdraw(arg1 *string, arg2 *string, arg3 *string) *string

//go:wasmimport sdk contracts.read
func contractRead(contractId *string, key *string) *string

//go:wasmimport sdk contracts.call
func contractCall(contractId *string, method *string, payload *string, options *string) *string

//go:wasmimport sdk tss.create_key
func tssCreateKey(keyId *string, algo *string) *string

//go:wasmimport sdk tss_v2.create_key
func tssV2CreateKey(keyId *string, algo *string, epochs *string) *string

//go:wasmimport sdk tss_v2.renew_key
func tssRenewKey(keyId *string, epochs *string) *string

//go:wasmimport sdk tss.sign_key
func tssSignKey(keyId *string, msgId *string) *string

//go:wasmimport sdk tss.get_key
func tssGetKey(keyId *string) *string

//go:wasmimport sdk crypto.keccak256
func cryptoKeccak256(hexData *string) *string

//go:wasmimport sdk crypto.ecrecover
func cryptoEcrecover(hashHex *string, sigHex *string) *string

//go:wasmimport sdk crypto.rlp_decode
func cryptoRlpDecode(hexData *string) *string

//go:wasmimport sdk crypto.sp1_verify_groth16
func cryptoSp1VerifyGroth16(proof *string, publicInputs *string, sp1VkeyHash *string, groth16Vk *string, vkRoot *string) *string

//go:wasmimport env abort
func abort(msg, file *string, line, column *int32)

//go:wasmimport env revert
func revert(msg, symbol *string)
