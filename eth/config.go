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
	"Platon-go/node"
	"fmt"
	"math/big"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/consensus/ethash"
	"Platon-go/core"
	"Platon-go/eth/downloader"
	"Platon-go/eth/gasprice"
	"Platon-go/log"
	"Platon-go/params"
)

const (
	datadirCbftConfig = "cbft.json" // Path within the datadir to the cbft config
)

// DefaultConfig contains default settings for use on the Ethereum main net.
var DefaultConfig = Config{
	SyncMode: downloader.FastSync,
	CbftConfig: CbftConfig{
		Period:           1,
		Epoch:            250000,
		MaxLatency:       600,
		LegalCoefficient: 1.0,
		Duration:         10,
	},
	Ethash: ethash.Config{
		CacheDir:       "ethash",
		CachesInMem:    2,
		CachesOnDisk:   3,
		DatasetsInMem:  1,
		DatasetsOnDisk: 2,
	},
	NetworkId:     1,
	LightPeers:    100,
	DatabaseCache: 768,
	TrieCache:     256,
	TrieTimeout:   60 * time.Minute,
	MinerGasFloor: 3150000000,
	MinerGasCeil:  3150000000,
	MinerGasPrice: big.NewInt(params.GWei),
	MinerRecommit: 3 * time.Second,

	TxPool: core.DefaultTxPoolConfig,
	GPO: gasprice.Config{
		Blocks:     20,
		Percentile: 60,
	},

	MPCPool: core.DefaultMPCPoolConfig ,

}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
	if runtime.GOOS == "windows" {
		DefaultConfig.Ethash.DatasetDir = filepath.Join(home, "AppData", "Ethash")
	} else {
		DefaultConfig.Ethash.DatasetDir = filepath.Join(home, ".ethash")
	}
}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the Ethereum main net block is used.
	Genesis *core.Genesis `toml:",omitempty"`

	CbftConfig	CbftConfig `toml:",omitempty"`

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

	// Mining-related options
	Etherbase      common.Address `toml:",omitempty"`
	MinerNotify    []string       `toml:",omitempty"`
	MinerExtraData []byte         `toml:",omitempty"`
	MinerGasFloor  uint64
	MinerGasCeil   uint64
	MinerGasPrice  *big.Int
	MinerRecommit  time.Duration
	MinerNoverify  bool

	// Ethash options
	Ethash ethash.Config

	// Transaction pool options
	TxPool core.TxPoolConfig

	// Gas Price Oracle options
	GPO gasprice.Config

	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool

	// Miscellaneous options
	DocRoot string `toml:"-"`

	// Type of the EWASM interpreter ("" for detault)
	EWASMInterpreter string
	// Type of the EVM interpreter ("" for default)
	EVMInterpreter string


	// MPC pool options
	MPCPool core.MPCPoolConfig
}

type CbftConfig struct {
	Period           uint64  `json:"period"`           // Number of seconds between blocks to enforce
	Epoch            uint64  `json:"epoch"`            // Epoch length to reset votes and checkpoint
	MaxLatency       int64   `json:"maxLatency"`
	LegalCoefficient float64 `json:"legalCoefficient"`
	Duration         int64   `json:"duration"`
	//mock
	//InitialNodes []discover.Node   `json:"initialNodes"`
	//NodeID       discover.NodeID   `json:"nodeID,omitempty"`
	//PrivateKey   *ecdsa.PrivateKey `json:"PrivateKey,omitempty"`
	Ppos 			*PposConfig 	`json:"ppos"`
}

// modify by platon
type PposConfig struct {
	Candidate 				*CandidateConfig 			`json:"candidate"`
}
// modify by platon
type CandidateConfig struct {
	// allow immediate elected max count
	MaxCount				uint64					`json:"maxCount"`
	// allow witness max count
	MaxChair				uint64					`json:"maxChair"`
	// allow block interval for refunds
	RefundBlockNumber 		uint64 					`json:"refundBlockNumber"`

}

type configMarshaling struct {
	MinerExtraData hexutil.Bytes
}

// StaticNodes returns a list of node enode URLs configured as static nodes.
func (c *Config) LoadCbftConfig(nodeConfig node.Config) *CbftConfig {
	return c.parsePersistentCbftConfig(filepath.Join(nodeConfig.DataDir, datadirCbftConfig))
}

// parsePersistentNodes parses a list of discovery node URLs loaded from a .json
// file from within the data directory.
func (c *Config) parsePersistentCbftConfig(path string) *CbftConfig {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	// Load the nodes from the config file.
	config := CbftConfig{}
	if err := common.LoadJSON(path, &config); err != nil {
		log.Error(fmt.Sprintf("Can't load cbft config file %s: %v", path, err))
		return nil
	}
	return &config
}
