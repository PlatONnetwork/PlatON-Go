package cbft

import "Platon-go/consensus"

type API struct {
	chain consensus.ChainReader
	cbft  *Cbft
}
