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
	kind, content, _, err := rlp.Split(input)
	switch {
	case err != nil:
		return "", nil, err
	case kind != rlp.List:
		return "", nil, fmt.Errorf("input type error")
	}

	_, funcName, params, err := rlp.Split(content)
	switch {
	case err != nil:
		return "", nil, err
		//case len(funcName) != 1:
		//	return "", nil, fmt.Errorf("funcName type error")
	}
	return string(funcName), params, nil

}
