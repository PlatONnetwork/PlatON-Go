package reward

import (
	"math/big"
	"testing"
)

func TestDecreaseDelegateReward(t *testing.T) {
	var receives []DelegateRewardReceipt
	var list DelegateRewardPerList

	receives = []DelegateRewardReceipt{

		{big.NewInt(200), 2},
		{big.NewInt(550), 3},
		{big.NewInt(400), 4},
		{big.NewInt(400), 5},
		{big.NewInt(800), 6},

		{big.NewInt(600), 7},
	}

	list.Pers = []*DelegateRewardPer{
		&DelegateRewardPer{big.NewInt(300), 1, nil},
		&DelegateRewardPer{big.NewInt(500), 2, nil},
		&DelegateRewardPer{big.NewInt(550), 3, nil},
		&DelegateRewardPer{big.NewInt(800), 4, nil},
		&DelegateRewardPer{big.NewInt(550), 5, nil},
	}
	index := list.DecreaseTotalAmount(receives)
	if index != 4 {
		t.Errorf("receives index is wrong,%v", index)
	}
	if list.Pers[1].DelegateAmount.Cmp(big.NewInt(300)) != 0 {
		t.Errorf("first DelegateAmount  is wrong,%v", list.Pers[1].DelegateAmount)
	}

	if list.Pers[len(list.Pers)-1].DelegateAmount.Cmp(big.NewInt(150)) != 0 {
		t.Errorf("latest DelegateAmount  is wrong,%v", list.Pers[1].DelegateAmount)
	}
}
