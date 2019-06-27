package gov

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"testing"
)

func getProposal() TextProposal {
	return TextProposal{
		common.Hash{0x01},
		"p#01",
		Version,
		"up,up,up....",
		"哈哈哈哈哈哈",
		"em。。。。",
		1000,
		1000000,
		discover.NodeID{},
		TallyResult{},
	}
}

func TestMustDecoded(t *testing.T) {

	p := getProposal()

	println(p.GetUrl())

	encoded, _ := rlp.EncodeToBytes(p)
	println(encoded)

	var txt TextProposal

	rlp.DecodeBytes(encoded, &txt)

	println(txt.Url)
}

func Test_Nothing(t *testing.T) {
	proposal := getProposal()
	proposalBytes, _ := json.Marshal(proposal)

	fmt.Printf("%s \n", hex.EncodeToString(proposalBytes))
	proposalBytes = append(proposalBytes, byte(proposal.GetProposalType()))

	fmt.Printf("%s \n", hex.EncodeToString(proposalBytes))

	var txp TextProposal
	json.Unmarshal(proposalBytes, &txp)
	fmt.Println(txp.String())

}
