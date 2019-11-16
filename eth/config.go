// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"math/big"
	"os"
	"os/user"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/eth/gasprice"
)

const (
	datadirCbftConfig = "cbft.json" // Path within the datadir to the cbft config
)

// DefaultConfig contains default settings for use on the Ethereum main net.
var DefaultConfig = Config{
	SyncMode: downloader.FullSync,
	CbftConfig: types.OptionsConfig{
		WalMode:           true,
		PeerMsgQueueSize:  1024,
		EvidenceDir:       "evidence",
		MaxPingLatency:    5000,
		MaxQueuesLimit:    4096,
		BlacklistDeadline: 60,
		Period:            20000,
		Amount:            10,
	},
	NetworkId:     1,
	LightPeers:    100,
	DatabaseCache: 768,
	TrieCache:     256,
	TrieTimeout:   60 * time.Minute,
	MinerGasFloor: params.GenesisGasLimit,
	//MinerGasCeil:  4000 * 21000 * 1.2,
	DBDisabledGC:  false,
	DBGCInterval:  86400,
	DBGCTimeout:   time.Minute,
	DBGCMpt:       true,
	DBGCBlock:     10,
	MinerGasPrice: big.NewInt(params.GVon),
	MinerRecommit: 3 * time.Second,

	MiningLogAtDepth:       7,
	TxChanSize:             4096,
	ChainHeadChanSize:      10,
	ChainSideChanSize:      10,
	ResultQueueSize:        10,
	ResubmitAdjustChanSize: 10,
	MinRecommitInterval:    1 * time.Second,
	MaxRecommitInterval:    15 * time.Second,
	IntervalAdjustRatio:    0.1,
	IntervalAdjustBias:     200 * 1000.0 * 1000.0,
	StaleThreshold:         7,
	DefaultCommitRatio:     0.95,

	BodyCacheLimit:    256,
	BlockCacheLimit:   256,
	MaxFutureBlocks:   256,
	BadBlockLimit:     10,
	TriesInMemory:     128,
	BlockChainVersion: 3,

	TxPool: core.DefaultTxPoolConfig,
	GPO: gasprice.Config{
		Blocks:     20,
		Percentile: 60,
	},

	//MPCPool: core.DefaultMPCPoolConfig,
	//VCPool:  core.DefaultVCPoolConfig,
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the Ethereum main net block is used.
	Genesis *core.Genesis `toml:",omitempty"`

	CbftConfig types.OptionsConfig `toml:",omitempty"`

	// Protocol options
	NetworkId uint64 // Network ID to use for selecting peers to connect to
	SyncMode  downloader.SyncMode
	NoPruning bool

	// Light client options
	LightServ  int `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
	LightPeers int `toml:",omitempty"` // Maximum number of LES client peers

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
	TrieCache          int
	TrieTimeout        time.Duration
	DBDisabledGC       bool
	DBGCInterval       uint64
	DBGCTimeout        time.Duration
	DBGCMpt            bool
	DBGCBlock          uint64

	// Mining-related options
	MinerExtraData []byte `toml:",omitempty"`
	MinerGasFloor  uint64
	//MinerGasCeil   uint64
	MinerGasPrice *big.Int
	MinerRecommit time.Duration
	MinerNoverify bool
	// minning conig
	MiningLogAtDepth       uint          // miningLogAtDepth is the number of confirmations before logging successful mining.
	TxChanSize             int           // txChanSize is the size of channel listening to NewTxsEvent.The number is referenced from the size of tx pool.
	ChainHeadChanSize      int           // chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	ChainSideChanSize      int           // chainSideChanSize is the size of channel listening to ChainSideEvent.
	ResultQueueSize        int           // resultQueueSize is the size of channel listening to sealing result.
	ResubmitAdjustChanSize int           // resubmitAdjustChanSize is the size of resubmitting interval adjustment channel.
	MinRecommitInterval    time.Duration // minRecommitInterval is the minimal time interval to recreate the mining block with any newly arrived transactions.
	MaxRecommitInterval    time.Duration // maxRecommitInterval is the maximum time interval to recreate the mining block with any newly arrived transactions.
	IntervalAdjustRatio    float64       // intervalAdjustRatio is the impact a single interval adjustment has on sealing work resubmitting interval.
	IntervalAdjustBias     float64       // intervalAdjustBias is applied during the new resubmit interval calculation in favor of increasing upper limit or decreasing lower limit so that the limit can be reachable.
	StaleThreshold         uint64        // staleThreshold is the maximum depth of the acceptable stale block.
	DefaultCommitRatio     float64

	// block config
	BodyCacheLimit           int
	BlockCacheLimit          int
	MaxFutureBlocks          int
	BadBlockLimit            int
	TriesInMemory            int
	BlockChainVersion        int // BlockChainVersion ensures that an incompatible database forces a resync from scratch.
	DefaultTxsCacheSize      int
	DefaultBroadcastInterval time.Duration

	// Transaction pool options
	TxPool core.TxPoolConfig

	// Gas Price Oracle options
	GPO gasprice.Config

	// Miscellaneous options
	DocRoot string `toml:"-"`

	// MPC pool options
	//MPCPool core.MPCPoolConfig
	//VCPool  core.VCPoolConfig
	Debug bool
}

type configMarshaling struct {
	MinerExtraData hexutil.Bytes
}
