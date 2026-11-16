package plugin

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
)

// Shares normally covers Released+RestrictingPlan, so the candidate here is
// deliberately inconsistent to reach the shortfall branch.
func TestToSlash_SharesShortfallGatedByHawking(t *testing.T) {
	params.GetEc(params.DefaultTestNet)

	for _, tc := range []struct {
		name         string
		hawkingBlock *big.Int
		wantSlashed  bool
	}{
		{"pre-Hawking fails after the balance slash", nil, true},
		{"Hawking fails before any write", big.NewInt(0), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			chain := mock.NewChain()
			defer chain.SnapDB.Clear()

			stkDB := staking.NewStakingDBWithDB(chain.SnapDB)
			sk := &StakingPlugin{
				db:          stkDB,
				chainConfig: &params.ChainConfig{PauliBlock: big.NewInt(0), HawkingBlock: tc.hawkingBlock},
			}

			can, addr := newZeroDelegateCandidate(new(big.Int))
			can.Released = big.NewInt(1000)
			can.Shares = big.NewInt(10)
			chain.StateDB.AddBalance(vm.StakingContractAddr, big.NewInt(1000))

			hash := common.HexToHash("0x01")
			if err := chain.SnapDB.NewBlock(big.NewInt(1), chain.CurrentHeader().ParentHash, hash); err != nil {
				t.Fatal(err)
			}
			if err := stkDB.SetCanBaseStore(hash, addr, can.CandidateBase); err != nil {
				t.Fatal(err)
			}
			if err := stkDB.SetCanMutableStore(hash, addr, can.CandidateMutable); err != nil {
				t.Fatal(err)
			}

			reporter := common.HexToAddress("0x1000000000000000000000000000000000000001")
			_, err := sk.toSlash(chain.StateDB, 1, hash, &staking.SlashNodeItem{
				NodeId:      can.NodeId,
				Amount:      big.NewInt(100),
				SlashType:   staking.DuplicateSign,
				BenefitAddr: reporter,
			})
			if err != errCandidateSharesNotEnough {
				t.Fatalf("err = %v, want %v", err, errCandidateSharesNotEnough)
			}
			if slashed := chain.StateDB.GetBalance(reporter).Sign() != 0; slashed != tc.wantSlashed {
				t.Fatalf("balance slashed = %v, want %v", slashed, tc.wantSlashed)
			}
		})
	}
}
