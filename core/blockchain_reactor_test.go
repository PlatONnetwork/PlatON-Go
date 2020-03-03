package core

import (
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
)

func TestBlockChainReactorClose(t *testing.T) {
	t.Run("close after commit", func(t *testing.T) {
		eventmux := new(event.TypeMux)
		reacter := NewBlockChainReactor(eventmux, big.NewInt(101))
		reacter.Start(common.PPOS_VALIDATOR_MODE)
		var parenthash common.Hash
		cbftress := make(chan cbfttypes.CbftResult, 5)
		go func() {
			for i := 1; i < 11; i++ {
				header := new(types.Header)
				header.Number = big.NewInt(int64(i))
				header.Time = big.NewInt(int64(i))
				header.ParentHash = parenthash
				block := types.NewBlock(header, nil, nil)
				snapshotdb.Instance().NewBlock(header.Number, header.ParentHash, block.Hash())
				parenthash = block.Hash()
				cbftress <- cbfttypes.CbftResult{Block: block}
			}
			close(cbftress)
		}()

		for value := range cbftress {
			eventmux.Post(value)
		}

		reacter.Close()

		time.Sleep(time.Second)
		snapshotdb.Instance().Clear()
	})
}
