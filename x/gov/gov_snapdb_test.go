package gov

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"testing"
)

func getTProposal() TextProposal {
	return TextProposal{
		common.Hash{0x01},
		"p#01",
		Version,
		"up,up,up....",
		"哈哈哈哈哈哈",
		"em。。。。",
		big.NewInt(1000),
		big.NewInt(1000000),
		discover.NodeID{},
		TallyResult{},
	}
}

func getVProposal() VersionProposal {
	return VersionProposal{
		getTProposal(),
		100,
		big.NewInt(1000),
	}
}

func TestGovSnapshotDB_AddVotingProposal(t *testing.T) {
	//var  govSnapdb GovSnapshotDB
	//pId := common.Hash{0x01}
	//govSnapdb.AddVotingProposal(pId)

	var arr []common.Hash
	arr = append(arr, common.Hash{0x01})

	bytes, err := rlp.EncodeToBytes(arr)
	if err != nil {
		t.Error(err)
	}
	var arr2 []common.Hash
	if err = rlp.DecodeBytes(bytes, &arr2); err != nil {
		t.Error(err)
	}

	fmt.Println(arr2)

}
