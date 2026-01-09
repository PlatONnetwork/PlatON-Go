package params

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultEMConfig(t *testing.T) {
	t.Run("DefaultMainNet", func(t *testing.T) {
		if getDefaultEMConfig(DefaultMainNet) == nil {
			t.Error("DefaultMainNet can't be nil config")
		}
		if err := CheckEconomicModel(FORKVERSION_1_3_0); nil != err {
			t.Error(err)
		}
	})
	t.Run("DefaultTestNet", func(t *testing.T) {
		if getDefaultEMConfig(DefaultTestNet) == nil {
			t.Error("DefaultTestNet can't be nil config")
		}
		if err := CheckEconomicModel(FORKVERSION_1_3_0); nil != err {
			t.Error(err)
		}
	})
	t.Run("DefaultUnitTestNet", func(t *testing.T) {
		if getDefaultEMConfig(DefaultUnitTestNet) == nil {
			t.Error("DefaultUnitTestNet can't be nil config")
		}
		if err := CheckEconomicModel(FORKVERSION_1_3_0); nil != err {
			t.Error(err)
		}
	})
	if getDefaultEMConfig(10) != nil {
		x := getDefaultEMConfig(10)
		t.Log(x.Reward)
		t.Error("the chain config not support")
	}
}

func TestMainNetHash(t *testing.T) {
	tempEc := getDefaultEMConfig(DefaultMainNet)
	bytes, err := rlp.EncodeToBytes(tempEc)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, common.RlpHash(bytes).Hex() == MainNetECHash)
}
