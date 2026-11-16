package gov

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/params"
)

func TestActiveBlockOf(t *testing.T) {
	av := func(version uint32, block uint64) ActiveVersionValue {
		return ActiveVersionValue{ActiveVersion: version, ActiveBlock: block}
	}
	v140 := uint32(1<<16 | 4<<8)
	v153 := uint32(1<<16 | 5<<8 | 3)
	for _, tc := range []struct {
		name string
		list []ActiveVersionValue
		want int64 // -1 means not found
	}{
		{"empty", nil, -1},
		{"before 1.5.0", []ActiveVersionValue{av(v140, 100), av(params.GenesisVersion, 0)}, -1},
		{"1.5.0 is current", []ActiveVersionValue{av(params.FORKVERSION_1_5_0, 500), av(v140, 100), av(params.GenesisVersion, 0)}, 500},
		{"1.6.0 activated after 1.5.0", []ActiveVersionValue{av(params.FORKVERSION_1_6_0, 900), av(params.FORKVERSION_1_5_0, 500), av(v140, 100)}, 500},
		// A 1.5.3 genesis never activated 1.5.0 by governance.
		{"1.5.3 genesis", []ActiveVersionValue{av(v153, 0)}, -1},
	} {
		block, ok := ActiveBlockOf(tc.list, params.FORKVERSION_1_5_0)
		switch {
		case tc.want < 0 && ok:
			t.Errorf("%s: got %d, want not found", tc.name, block)
		case tc.want >= 0 && (!ok || block != uint64(tc.want)):
			t.Errorf("%s: got %d (found %v), want %d", tc.name, block, ok, tc.want)
		}
	}
}
