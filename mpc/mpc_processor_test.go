package mpc

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"testing"
)

func TestMPCExec(t *testing.T) {
	param := MPCParams {
		TaskId : 	"ddc1ce5ccf0fac429a3aa075f28f73de9831dc60829bebeae08372e641a0da4b",
		Pubkey : "a363d1243646b6eabf1d4851f646b523f5707d053caab95022f1682605aca0537ee0c5c14b4dfa76dcbce264b7e68d59de79a42b7cda059e9d358336a9ab8d80",
		From  : common.HexToAddress("0x60Ceca9c1290EE56b98d4E160EF0453F7C40d219"),
		IRAddr : common.HexToAddress("0xC1FB0780933718Ccb12DF862726776C840F22C33"),
		Method : "start_calc",
		Extra : "",
	}
	err := ExecuteMPCTx(param)
	if err != nil {
		t.Error("exec fail.")
	}
}