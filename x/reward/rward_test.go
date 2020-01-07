package reward

import (
	"math/big"
	"testing"
)

func TestDecreaseDelegateReward(t *testing.T) {
	var receives []DelegateRewardReceipt
	var list DelegateRewardPerList

	receives = []DelegateRewardReceipt{

		DelegateRewardReceipt{big.NewInt(300), 1},
		{big.NewInt(200), 2},
		{big.NewInt(300), 3},
	}

	list.Pers = []*DelegateRewardPer{
		&DelegateRewardPer{big.NewInt(300), 1, nil},
		&DelegateRewardPer{big.NewInt(500), 2, nil},
		&DelegateRewardPer{big.NewInt(550), 3, nil},
		&DelegateRewardPer{big.NewInt(550), 5, nil},
	}
	list.DecreaseTotalAmount(receives)
	if len(list.Pers) != 3 {
		t.Error("list must decrease")
	}
	if list.Pers[0].DelegateAmount.Cmp(big.NewInt(300)) != 0 {
		t.Error("epoch 2 must same")
	}
	if list.Pers[0].Epoch != 2 {
		t.Error("epoch must 2")
	}
	if list.Pers[1].DelegateAmount.Cmp(big.NewInt(250)) != 0 {
		t.Error("epoch 3 must same")
	}
	if list.Pers[2].DelegateAmount.Cmp(big.NewInt(550)) != 0 {
		t.Error("epoch 5 must same")
	}
	if list.changed != true {
		t.Error("must changed")
	}

}
