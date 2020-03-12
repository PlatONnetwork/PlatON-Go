package vm

type WasmParams struct {
	FuncName []byte
	Args     [][]byte
}

type WasmDeploy struct {
	VM   byte
	Args [][]byte
}

type WasmInoke struct {
	VM   byte
	Args *WasmParams
}
