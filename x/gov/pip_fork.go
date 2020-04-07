package gov

import (
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

func CheckForkPIP0_11_0(state xcom.StateDB) bool {
	if GetCurrentActiveVersion(state) >= params.FORKVERSION_0_11_0 {
		return true
	} else {
		return false
	}
}
