package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
)

type API struct {
	chain consensus.ChainReader
	cbft  *Cbft
}

// Get the block address of producer.
func (api *API) GetProducer(number *rpc.BlockNumber) (common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the signers from its snapshot
	if header == nil {
		return common.Address{}, errUnknownBlock
	}
	nodeID, _, err := ecrecover(header)

	if err != nil {
		return common.Address{}, err
	}

	var signer common.Address
	copy(signer[:], crypto.Keccak256(nodeID.Bytes()[1:])[12:])

	return signer, nil
}
