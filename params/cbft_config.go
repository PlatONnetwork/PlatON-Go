package params

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
)

type CbftNode struct {
	Node      *enode.Node   `json:"node"`
	BlsPubKey bls.PublicKey `json:"blsPubKey"`
}

type initNode struct {
	Enode     string
	BlsPubkey string
}

type CbftConfig struct {
	Period        uint64     `json:"period,omitempty"`        // Number of seconds between blocks to enforce
	Amount        uint32     `json:"amount,omitempty"`        //The maximum number of blocks generated per cycle
	InitialNodes  []CbftNode `json:"initialNodes,omitempty"`  //Genesis consensus node
	ValidatorMode string     `json:"validatorMode,omitempty"` //Validator mode for easy testing
}

func ConvertNodeUrl(initialNodes []initNode) []CbftNode {
	bls.Init(bls.BLS12_381)
	NodeList := make([]CbftNode, 0, len(initialNodes))
	for _, n := range initialNodes {

		cbftNode := new(CbftNode)

		if node, err := enode.Parse(enode.ValidSchemes, n.Enode); nil == err {
			cbftNode.Node = node
		}

		if n.BlsPubkey != "" {
			var blsPk bls.PublicKey
			if err := blsPk.UnmarshalText([]byte(n.BlsPubkey)); nil == err {
				cbftNode.BlsPubKey = blsPk
			}
		}

		NodeList = append(NodeList, *cbftNode)
	}
	return NodeList
}

// String implements the fmt.Stringer interface.
func (c *CbftConfig) String() string {
	initialNodes := make([]initNode, 0)
	for _, node := range c.InitialNodes {
		initialNodes = append(initialNodes, initNode{
			Enode: node.Node.String(),
			//	BlsPubkey: node.BlsPubKey.GetHexString(),
		})
	}

	return fmt.Sprintf("{period: %v  amount: %v initialNodes: %v validatorMode: %v}",
		c.Period,
		c.Amount,
		initialNodes,
		c.ValidatorMode,
	)
}
