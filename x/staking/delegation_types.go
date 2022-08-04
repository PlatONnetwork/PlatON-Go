package staking

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
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
	del.LockReleasedHes = new(big.Int).SetInt64(0)
	del.LockRestrictingHes = new(big.Int).SetInt64(0)
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

	// 处于犹豫期的委托金,源自锁定期
	LockReleasedHes    *big.Int
	LockRestrictingHes *big.Int
}

type v1StoredDelegationRlp struct {
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
	return fmt.Sprintf(`{DelegateEpoch: %d,Released: %d,ReleasedHes: %d,RestrictingPlan: %d,RestrictingPlanHes: %d,CumulativeIncome: %d,LockReleasedHes:%d,LockRestrictingHes,%d}`,
		del.DelegateEpoch,
		del.Released,
		del.ReleasedHes,
		del.RestrictingPlan,
		del.RestrictingPlanHes,
		del.CumulativeIncome,
		del.LockReleasedHes,
		del.LockRestrictingHes,
	)
}

func (del *Delegation) IsEmpty() bool {
	return nil == del
}

type DelegationForStorage Delegation

// DecodeRLP implements rlp.Decoder, and loads both consensus and implementation
// fields of a Delegation from an RLP stream.
func (r *DelegationForStorage) DecodeRLP(s *rlp.Stream) error {
	// Retrieve the entire delegation blob as we need to try multiple decoders
	blob, err := s.Raw()
	if err != nil {
		return err
	}
	// Try decoding from the newest format for future proofness, then the older one
	if err := decodeStoredDelegateRLP(r, blob); err == nil {
		return nil
	}
	return decodeV1StoredDelegateRLP(r, blob)
}

func decodeStoredDelegateRLP(r *DelegationForStorage, blob []byte) error {
	var stored Delegation
	if err := rlp.DecodeBytes(blob, &stored); err != nil {
		return err
	}
	r.DelegateEpoch = stored.DelegateEpoch
	r.Released = stored.Released
	r.ReleasedHes = stored.ReleasedHes
	r.RestrictingPlan = stored.RestrictingPlan
	r.RestrictingPlanHes = stored.RestrictingPlanHes
	r.CumulativeIncome = stored.CumulativeIncome
	r.LockReleasedHes = stored.LockReleasedHes
	r.LockRestrictingHes = stored.LockRestrictingHes
	return nil
}

func decodeV1StoredDelegateRLP(r *DelegationForStorage, blob []byte) error {
	var stored v1StoredDelegationRlp
	if err := rlp.DecodeBytes(blob, &stored); err != nil {
		return err
	}
	r.DelegateEpoch = stored.DelegateEpoch
	r.Released = stored.Released
	r.ReleasedHes = stored.ReleasedHes
	r.RestrictingPlan = stored.RestrictingPlan
	r.RestrictingPlanHes = stored.RestrictingPlanHes
	r.CumulativeIncome = stored.CumulativeIncome
	r.LockReleasedHes = new(big.Int)
	r.LockRestrictingHes = new(big.Int)
	return nil
}

func encodeStoredDelegateRLP(d *Delegation) ([]byte, error) {
	return rlp.EncodeToBytes(d)
}

func encodeV1StoredDelegateRLP(d *Delegation) ([]byte, error) {
	stored := new(v1StoredDelegationRlp)
	stored.DelegateEpoch = d.DelegateEpoch
	stored.Released = d.Released
	stored.ReleasedHes = d.ReleasedHes
	stored.RestrictingPlan = d.RestrictingPlan
	stored.RestrictingPlanHes = d.RestrictingPlanHes
	stored.CumulativeIncome = d.CumulativeIncome
	return rlp.EncodeToBytes(stored)
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

	LockReleasedHes    *hexutil.Big
	LockRestrictingHes *hexutil.Big
}

type DelegationHexV1 struct {
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
	return fmt.Sprintf(`{"DelegateEpoch": "%d","Released": "%s","ReleasedHes": %s,"RestrictingPlan": %s,"RestrictingPlanHes": %s,"CumulativeIncome": %s,"LockReleasedHes":%s,"LockRestrictingHes",%s}`,
		delHex.DelegateEpoch,
		delHex.Released,
		delHex.ReleasedHes,
		delHex.RestrictingPlan,
		delHex.RestrictingPlanHes,
		delHex.CumulativeIncome,
		delHex.LockReleasedHes,
		delHex.LockRestrictingHes,
	)
}

type DelegationEx struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
	DelegationHex
}

func (dex *DelegationEx) String() string {
	return fmt.Sprintf(`{"Addr": "%s","NodeId": "%s","StakingBlockNum": "%d","DelegateEpoch": "%d","Released": "%s","ReleasedHes": %s,"RestrictingPlan": %s,"RestrictingPlanHes": %s,"CumulativeIncome": %s,"LockReleasedHes":%s,"LockRestrictingHes",%s}`,
		dex.Addr.String(),
		fmt.Sprintf("%x", dex.NodeId.Bytes()),
		dex.StakingBlockNum,
		dex.DelegateEpoch,
		dex.Released,
		dex.ReleasedHes,
		dex.RestrictingPlan,
		dex.RestrictingPlanHes,
		dex.CumulativeIncome,
		dex.LockReleasedHes,
		dex.LockRestrictingHes,
	)
}

func (dex *DelegationEx) IsEmpty() bool {
	return nil == dex
}

type DelegationExV1 struct {
	Addr            common.Address
	NodeId          discover.NodeID
	StakingBlockNum uint64
	DelegationHexV1
}

func (dex *DelegationEx) V1() *DelegationExV1 {
	return &DelegationExV1{
		Addr:            dex.Addr,
		NodeId:          dex.NodeId,
		StakingBlockNum: dex.StakingBlockNum,
		DelegationHexV1: DelegationHexV1{
			DelegateEpoch:      dex.DelegateEpoch,
			Released:           dex.Released,
			ReleasedHes:        dex.ReleasedHes,
			RestrictingPlan:    dex.RestrictingPlan,
			RestrictingPlanHes: dex.RestrictingPlanHes,
			CumulativeIncome:   dex.CumulativeIncome,
		},
	}
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

type DelegationInfo struct {
	NodeID           discover.NodeID
	StakeBlockNumber uint64
	Delegation       *Delegation
}

type DelByDelegateEpoch []*DelegationInfo

func (d DelByDelegateEpoch) Len() int { return len(d) }
func (d DelByDelegateEpoch) Less(i, j int) bool {
	return d[i].Delegation.DelegateEpoch < d[j].Delegation.DelegateEpoch
}
func (d DelByDelegateEpoch) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

type DelegationLock struct {
	// 锁定期
	Locks []*DelegationLockPeriod

	// 解锁期
	Released        *big.Int
	RestrictingPlan *big.Int

	change bool
}

func NewDelegationLock() *DelegationLock {
	store := make([]*DelegationLockPeriod, 0)
	return &DelegationLock{
		Locks:           store,
		Released:        new(big.Int),
		RestrictingPlan: new(big.Int),
	}
}

func (d *DelegationLock) AddLock(lockEpoch uint32, released, restrictingPlan *big.Int) {
	lastStore := len(d.Locks) - 1
	if len(d.Locks) == 0 || d.Locks[lastStore].Epoch != lockEpoch {
		d.Locks = append(d.Locks, newDelegationLockPeriod(lockEpoch, released, restrictingPlan))
	} else {
		if released != nil {
			d.Locks[lastStore].Released.Add(d.Locks[lastStore].Released, released)
		}
		if restrictingPlan != nil {
			d.Locks[lastStore].RestrictingPlan.Add(d.Locks[lastStore].RestrictingPlan, restrictingPlan)
		}
	}
	d.change = true
}

func (d *DelegationLock) Change() bool {
	return d.change
}

// update 根据当前结算周期更新锁定期与解锁期的委托金
func (d *DelegationLock) update(currentEpoch uint32) {
	end := 0
	change := false
	for i := 0; i < len(d.Locks); i++ {
		if d.Locks[i].Epoch < currentEpoch {
			d.Released.Add(d.Released, d.Locks[i].Released)
			d.RestrictingPlan.Add(d.RestrictingPlan, d.Locks[i].RestrictingPlan)
			end = i
			change = true
		} else {
			break
		}
	}
	if change {
		if end+1 == len(d.Locks) {
			d.Locks = d.Locks[:0]
		} else {
			d.Locks = d.Locks[end+1:]
		}
		d.change = true
	}
}

// AdvanceLockedFunds 使用锁定期的委托金
func (d *DelegationLock) AdvanceLockedFunds(amount *big.Int) (*big.Int, *big.Int, error) {
	if len(d.Locks) == 0 {
		return nil, nil, ErrDelegateLockBalanceNotEnough
	}
	restricting, released, left := new(big.Int), new(big.Int), new(big.Int).Set(amount)

	for i := len(d.Locks) - 1; i >= 0; i-- {
		// 使用来自锁仓账户部分
		if left.Cmp(d.Locks[i].RestrictingPlan) > 0 {
			left.Sub(left, d.Locks[i].RestrictingPlan)
			restricting.Add(restricting, d.Locks[i].RestrictingPlan)
		} else {
			d.Locks[i].RestrictingPlan.Sub(d.Locks[i].RestrictingPlan, left)
			restricting.Add(restricting, left)
			d.change = true
			d.Locks = d.Locks[:i+1]
			return released, restricting, nil
		}

		// 使用来自余额部分
		if left.Cmp(d.Locks[i].Released) > 0 {
			left.Sub(left, d.Locks[i].Released)
			released.Add(released, d.Locks[i].Released)
		} else {
			d.Locks[i].RestrictingPlan.SetInt64(0)
			d.Locks[i].Released.Sub(d.Locks[i].Released, left)
			released.Add(released, left)
			d.change = true
			if d.Locks[i].Released.Cmp(common.Big0) == 0 {
				d.Locks = d.Locks[:i]
			} else {
				d.Locks = d.Locks[:i+1]
			}
			return released, restricting, nil
		}
	}
	return nil, nil, ErrDelegateLockBalanceNotEnough
}

func (d *DelegationLock) shouldDel() bool {
	if len(d.Locks) == 0 {
		total := new(big.Int).Add(d.RestrictingPlan, d.Released)
		if total.Cmp(common.Big0) == 0 {
			return true
		}
	}
	return false
}

func (d *DelegationLock) ToHex() *DelegationLockHex {
	hex := new(DelegationLockHex)
	hex.Released = (*hexutil.Big)(d.Released)
	hex.RestrictingPlan = (*hexutil.Big)(d.RestrictingPlan)
	hex.Locks = make([]*DelegationLockPeriodHex, 0)
	for _, lock := range d.Locks {
		hex.Locks = append(hex.Locks, lock.ToHex())
	}
	return hex
}

type DelegationLockPeriod struct {
	// 锁定截止周期
	Epoch uint32
	//处于锁定期的委托金,解锁后释放到用户余额
	Released *big.Int
	//处于锁定期的委托金,解锁后释放到用户锁仓账户
	RestrictingPlan *big.Int
}

func (d *DelegationLockPeriod) ToHex() *DelegationLockPeriodHex {
	hex := new(DelegationLockPeriodHex)
	hex.Released = (*hexutil.Big)(d.Released)
	hex.Epoch = d.Epoch
	hex.RestrictingPlan = (*hexutil.Big)(d.RestrictingPlan)
	return hex
}

func (d *DelegationLockPeriod) String() string {
	return fmt.Sprintf(`{Epoch: %d,Released: %d,RestrictingPlan: %d}`,
		d.Epoch,
		d.Released,
		d.RestrictingPlan,
	)
}

func newDelegationLockPeriod(epoch uint32, released, restrictingPlan *big.Int) *DelegationLockPeriod {
	info := new(DelegationLockPeriod)
	info.Epoch = epoch
	if released != nil {
		info.Released = new(big.Int).Set(released)
	} else {
		info.Released = new(big.Int)
	}
	if restrictingPlan != nil {
		info.RestrictingPlan = new(big.Int).Set(restrictingPlan)
	} else {
		info.RestrictingPlan = new(big.Int)
	}
	return info
}

type DelegationLockHex struct {
	// 锁定期
	Locks []*DelegationLockPeriodHex
	//处于解锁期的委托金
	Released *hexutil.Big
	//处于解锁期的委托金
	RestrictingPlan *hexutil.Big
}

type DelegationLockPeriodHex struct {
	// 锁定截止周期
	Epoch uint32
	//处于锁定期的委托金,解锁后释放到用户余额
	Released *hexutil.Big
	//处于锁定期的委托金,解锁后释放到用户锁仓账户
	RestrictingPlan *hexutil.Big
}
