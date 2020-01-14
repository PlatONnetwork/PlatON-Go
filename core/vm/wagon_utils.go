package vm

import (
	"bytes"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/wagon/exec"
	"github.com/PlatONnetwork/wagon/validate"
	"github.com/PlatONnetwork/wagon/wasm"
)

func ReadWasmModule(Code []byte, verify bool) (*exec.CompiledModule, error) {
	m, err := wasm.ReadModule(bytes.NewReader(Code), func(name string) (*wasm.Module, error) {
		switch name {
		case "env":
			return NewHostModule(), nil
		}
		return nil, fmt.Errorf("module %q unknown", name)
	})
	if err != nil {
		return nil, err
	}

	if verify {
		err = validate.VerifyModule(m)
		if err != nil {
			return nil, err
		}
	}

	compiled, err := exec.CompileModule(m)

	if err != nil {
		return nil, err
	}

	return compiled, nil
}

func decodeFuncAndParams(input []byte) (string, []byte, error) {
	content, _, err := rlp.SplitList(input)
	if nil != err {
		return "", nil, fmt.Errorf("failed to decode input funcName and params: %v", err)
	}

	funcName, params, err := rlp.SplitString(content)
	if nil != err {
		return "", nil, fmt.Errorf("failed to decode input funcName and params: %v", err)
	}
	return string(funcName), params, nil

}
