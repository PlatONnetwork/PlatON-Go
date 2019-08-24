package fetcher

import (
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/stretchr/testify/assert"
)

func TestLimitFetcher_AddTask(t *testing.T) {
	fetcher := NewLimitFetcher()

	msg1 := common.BytesToHash(utils.Rand32Bytes(32))
	msg2 := common.BytesToHash(utils.Rand32Bytes(32))
	msg3 := common.BytesToHash(utils.Rand32Bytes(32))
	assert.True(t, fetcher.AddTask(msg1))
	assert.True(t, fetcher.AddTask(msg2))
	assert.False(t, fetcher.AddTask(msg1))
	assert.True(t, fetcher.AddTask(msg3))
	assert.False(t, fetcher.AddTask(msg2))

	assert.Equal(t, 3, len(fetcher.fetching))

	time.Sleep(keepTimeout + 100*time.Millisecond)

	assert.Equal(t, 0, len(fetcher.fetching))

	fetcher.Stop()
}
