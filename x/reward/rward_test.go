// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package reward

import (
	"log"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
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
		&DelegateRewardPer{big.NewInt(300), 1, nil, nil},
		&DelegateRewardPer{big.NewInt(500), 2, nil, nil},
		&DelegateRewardPer{big.NewInt(550), 3, nil, nil},
		&DelegateRewardPer{big.NewInt(800), 4, nil, nil},
		&DelegateRewardPer{big.NewInt(550), 5, nil, nil},
	}
	index := list.DecreaseTotalAmount(receives)
	if index != 4 {
		t.Errorf("receives index is wrong,%v", index)
	}
	if list.Pers[1].Left.Cmp(big.NewInt(300)) != 0 {
		t.Errorf("first Left  is wrong,%v", list.Pers[1].Left)
	}

	if list.Pers[len(list.Pers)-1].Left.Cmp(big.NewInt(150)) != 0 {
		t.Errorf("latest Left  is wrong,%v", list.Pers[1].Left)
	}
}

func TestSize(t *testing.T) {
	delegate := new(big.Int).Mul(new(big.Int).SetInt64(10000000), big.NewInt(params.LAT))
	reward, _ := new(big.Int).SetString("135840374364973262032076", 10)
	per := new(big.Int).Div(reward, delegate)
	key := DelegateRewardPerKey(discover.MustHexID("0aa9805681d8f77c05f317efc141c97d5adb511ffb51f5a251d2d7a4a3a96d9a12adf39f06b702f0ccdff9eddc1790eb272dca31b0c47751d49b5931c58701e7"), 100, 10)

	list := NewDelegateRewardPerList()
	for i := 0; i < DelegateRewardPerLength; i++ {
		list.AppendDelegateRewardPer(NewDelegateRewardPer(uint64(i), per, delegate))
	}
	val, err := rlp.EncodeToBytes(list)
	if err != nil {
		t.Error(err)
		return
	}
	length := len(key) + len(val)

	log.Print("size of per", length*101/(1024*1024))

}
