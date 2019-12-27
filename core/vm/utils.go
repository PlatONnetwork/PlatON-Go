package vm

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func decodeInput(input []byte) (byte, []byte, error) {
	kind, content, _, err := rlp.Split(input)

	switch {
	case err != nil:
		return 0, nil, err
	case kind != rlp.List:
		return 0, nil, fmt.Errorf("input type error")
	}

	_, vmType, rest, err := rlp.Split(content)
	switch {
	case err != nil:
		return 0, nil, err
	case len(vmType) != 1:
		return 0, nil, fmt.Errorf("vm type error")
	}
	return vmType[0], rest, nil

}
