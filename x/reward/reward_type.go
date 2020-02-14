package reward

import (
	"encoding/json"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

func NewDelegateRewardPer(epoch uint64, totalReward, totalDelegate *big.Int) *DelegateRewardPer {
	return &DelegateRewardPer{
		Left:     totalDelegate,
		Epoch:    epoch,
		Delegate: totalDelegate,
		Reward:   totalReward,
	}
}

type DelegateRewardPer struct {
	//this is the node total effective delegate  amount at this epoch,will Decrease when receive delegate reward
	Left  *big.Int
	Epoch uint64
	//this is the node total delegate reward per amount at this epoch
	Delegate *big.Int
	//this is the node total effective delegate  amount at this epoch
	Reward *big.Int
}

func (d *DelegateRewardPer) CalDelegateReward(delegate *big.Int) *big.Int {
	tmp := new(big.Int).Mul(delegate, d.Reward)
	return new(big.Int).Div(tmp, d.Delegate)
}

func NewDelegateRewardPerList() *DelegateRewardPerList {
	del := new(DelegateRewardPerList)
	del.Pers = make([]*DelegateRewardPer, 0)
	return del
}

type DelegateRewardPerList struct {
	Pers    []*DelegateRewardPer
	changed bool
}

func (d DelegateRewardPerList) String() string {
	v, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(v)
}

func (d *DelegateRewardPerList) AppendDelegateRewardPer(per *DelegateRewardPer) {
	d.Pers = append(d.Pers, per)
}

func (d *DelegateRewardPerList) DecreaseTotalAmount(receipt []DelegateRewardReceipt) int {
	var indexOfReceipt int
	for indexOfList := 0; indexOfList < len(d.Pers) && indexOfReceipt < len(receipt); {
		if d.Pers[indexOfList].Epoch < receipt[indexOfReceipt].Epoch {
			indexOfList++
		} else if d.Pers[indexOfList].Epoch > receipt[indexOfReceipt].Epoch {
			indexOfReceipt++
		} else {
			d.Pers[indexOfList].Left.Sub(d.Pers[indexOfList].Left, receipt[indexOfReceipt].Delegate)
			if d.Pers[indexOfList].Left.Cmp(common.Big0) <= 0 {
				d.Pers = append(d.Pers[:indexOfList], d.Pers[indexOfList+1:]...)
			} else {
				indexOfList++
			}
			indexOfReceipt++
			d.changed = true
		}
	}
	return indexOfReceipt
}

func (d *DelegateRewardPerList) ShouldDel() bool {
	if len(d.Pers) == 0 {
		return true
	}
	return false
}

func (d *DelegateRewardPerList) IsChange() bool {
	return d.changed
}

type NodeDelegateReward struct {
	NodeID     discover.NodeID `json:"nodeID"`
	StakingNum uint64          `json:"stakingNum"`
	Reward     *big.Int        `json:"reward" rlp:"nil"`
}

type NodeDelegateRewardPresenter struct {
	NodeID     discover.NodeID `json:"nodeID" `
	Reward     *hexutil.Big    `json:"reward" `
	StakingNum uint64          `json:"stakingNum"`
}

type DelegateRewardReceipt struct {
	//this is the account  total effective delegate amount with the node  on this epoch
	Delegate *big.Int
	Epoch    uint64
}
