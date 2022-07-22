package staking

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

func NewDelegation() *Delegation {
	del := new(Delegation)
	// Prevent null pointer initialization
	del.Released = new(big.Int).SetInt64(0)
	del.RestrictingPlan = new(big.Int).SetInt64(0)
	del.ReleasedHes = new(big.Int).SetInt64(0)
	del.RestrictingPlanHes = new(big.Int).SetInt64(0)
	del.CumulativeIncome = new(big.Int).SetInt64(0)
	return del
}

// the Delegate information
type Delegation struct {
	// The epoch number at delegate or edit
	DelegateEpoch uint32
	// The delegate von  is circulating for effective epoch (in effect)
	Released *big.Int
	// The delegate von  is circulating for hesitant epoch (in hesitation)
	ReleasedHes *big.Int
	// The delegate von  is RestrictingPlan for effective epoch (in effect)
	RestrictingPlan *big.Int
	// The delegate von  is RestrictingPlan for hesitant epoch (in hesitation)
	RestrictingPlanHes *big.Int
	// Cumulative delegate income (Waiting for withdrawal)
	CumulativeIncome *big.Int
}

func (del *Delegation) CleanCumulativeIncome(epoch uint32) {
	del.CumulativeIncome = new(big.Int)
	del.DelegateEpoch = epoch
}

func (del *Delegation) String() string {
	return fmt.Sprintf(`{"DelegateEpoch": "%d","Released": "%d","ReleasedHes": %d,"RestrictingPlan": %d,"RestrictingPlanHes": %d,"CumulativeIncome": %d}`,
		del.DelegateEpoch,
		del.Released,
		del.ReleasedHes,
		del.RestrictingPlan,
		del.RestrictingPlanHes,
		del.CumulativeIncome)
}

func (del *Delegation) IsEmpty() bool {
	return nil == del
}

type DelegationHex struct {
	// The epoch number at delegate or edit
	DelegateEpoch uint32
	// The delegate von  is circulating for effective epoch (in effect)
	Released *hexutil.Big
	// The delegate von  is circulating for hesitant epoch (in hesitation)
	ReleasedHes *hexutil.Big
	// The delegate von  is RestrictingPlan for effective epoch (in effect)
	RestrictingPlan *hexutil.Big
	// The delegate von  is RestrictingPlan for hesitant epoch (in hesitation)
	RestrictingPlanHes *hexutil.Big
	// Cumulative delegate income (Waiting for withdrawal)
	CumulativeIncome *hexutil.Big
}

func (delHex *DelegationHex) String() string {
	return fmt.Sprintf(`{"DelegateEpoch": "%d","Released": "%s","ReleasedHes": %s,"RestrictingPlan": %s,"RestrictingPlanHes": %s,"CumulativeIncome": %s}`,
		delHex.DelegateEpoch,
		delHex.Released,
		delHex.ReleasedHes,
		delHex.RestrictingPlan,
		delHex.RestrictingPlanHes,
		delHex.CumulativeIncome)
}

type DelegationEx struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
	DelegationHex
}

func (dex *DelegationEx) String() string {
	return fmt.Sprintf(`{"Addr": "%s","NodeId": "%s","StakingBlockNum": "%d","DelegateEpoch": "%d","Released": "%s","ReleasedHes": %s,"RestrictingPlan": %s,"RestrictingPlanHes": %s,"CumulativeIncome": %s}`,
		dex.Addr.String(),
		fmt.Sprintf("%x", dex.NodeId.Bytes()),
		dex.StakingBlockNum,
		dex.DelegateEpoch,
		dex.Released,
		dex.ReleasedHes,
		dex.RestrictingPlan,
		dex.RestrictingPlanHes,
		dex.CumulativeIncome)
}

func (dex *DelegationEx) IsEmpty() bool {
	return nil == dex
}

type DelegateRelated struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
}

func (dr *DelegateRelated) String() string {
	return fmt.Sprintf(`{"Addr": "%s","NodeId": "%s","StakingBlockNum": "%d"}`,
		dr.Addr.String(),
		fmt.Sprintf("%x", dr.NodeId.Bytes()),
		dr.StakingBlockNum)
}

type DelRelatedQueue []*DelegateRelated

func (queue DelRelatedQueue) IsEmpty() bool {
	return len(queue) == 0
}

func (queue DelRelatedQueue) String() string {
	arr := make([]string, len(queue))
	for i, r := range queue {
		arr[i] = r.String()
	}
	return "[" + strings.Join(arr, ",") + "]"
}
