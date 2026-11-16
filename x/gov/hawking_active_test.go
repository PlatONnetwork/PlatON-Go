package gov

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/params"
)

func TestHawkingActiveConfigForkWithoutState(t *testing.T) {
	cfg := &params.ChainConfig{
		PauliBlock:   big.NewInt(0),
		HawkingBlock: big.NewInt(0),
	}
	if !HawkingActive(cfg, big.NewInt(1), nil) {
		t.Fatal("HawkingBlock must enable rules without governance state")
	}
	before := &params.ChainConfig{PauliBlock: big.NewInt(0)}
	if HawkingActive(before, big.NewInt(1), nil) {
		t.Fatal("missing HawkingBlock and state must stay inactive")
	}
}
