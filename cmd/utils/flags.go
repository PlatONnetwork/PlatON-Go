// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

// Package utils contains internal helper functions for go-ethereum commands.
package utils

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/eth/tracers"
	"github.com/PlatONnetwork/PlatON-Go/ethdb/remotedb"
	"io"
	"math"
	"math/big"
	"os"
	"path/filepath"
	godebug "runtime/debug"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/eth/ethconfig"
	"github.com/PlatONnetwork/PlatON-Go/internal/flags"

	"github.com/PlatONnetwork/PlatON-Go/metrics/exp"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"

	"github.com/PlatONnetwork/PlatON-Go/graphql"
	"github.com/PlatONnetwork/PlatON-Go/internal/ethapi"

	"github.com/PlatONnetwork/PlatON-Go/miner"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/accounts"
	"github.com/PlatONnetwork/PlatON-Go/accounts/keystore"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/fdlimit"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core"
	types2 "github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/eth/gasprice"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/ethstats"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
	"github.com/PlatONnetwork/PlatON-Go/metrics/influxdb"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/nat"
	"github.com/PlatONnetwork/PlatON-Go/p2p/netutil"
	"github.com/PlatONnetwork/PlatON-Go/params"
	gopsutil "github.com/shirou/gopsutil/mem"
)

func init() {
	cli.AppHelpTemplate = `{{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
	cli.CommandHelpTemplate = flags.CommandHelpTemplate
	cli.HelpPrinter = printHelp
}

func printHelp(out io.Writer, templ string, data interface{}) {
	funcMap := template.FuncMap{"join": strings.Join}
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	w := tabwriter.NewWriter(out, 38, 8, 2, ' ', 0)
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

// These are all the command line flags we support.
// If you add to this list, please remember to include the
// flag in the appropriate command definition.
//
// The flags are defined here so their names and help texts
// are the same for all commands.

var (
	// General settings
	DataDirFlag = DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: DirectoryString(node.DefaultDataDir()),
	}
	RemoteDBFlag = cli.StringFlag{
		Name:  "remotedb",
		Usage: "URL for remote database",
	}
	AncientFlag = DirectoryFlag{
		Name:  "datadir.ancient",
		Usage: "Data directory for ancient chain segments (default = inside chaindata)",
	}
	MinFreeDiskSpaceFlag = DirectoryFlag{
		Name:  "datadir.minfreedisk",
		Usage: "Minimum free disk space in MB, once reached triggers auto shut down (default = --cache.gc converted to MB, 0 = disabled)",
	}
	KeyStoreDirFlag = DirectoryFlag{
		Name:  "keystore",
		Usage: "Directory for the keystore (default = inside the datadir)",
	}
	NetworkIdFlag = cli.Uint64Flag{
		Name:  "networkid",
		Usage: "Explicitly set network id (integer)",
		Value: ethconfig.Defaults.NetworkId,
	}
	MainFlag = cli.BoolFlag{
		Name:  "main",
		Usage: "Mainnet network: pre-configured main network (default network)",
	}
	TestnetFlag = cli.BoolFlag{
		Name:  "testnet",
		Usage: "Testnet network: pre-configured test network",
	}
	AddressHRPFlag = cli.StringFlag{
		Name:  "addressHRP",
		Usage: "set the address hrp,if not set,use default address hrp",
	}
	IdentityFlag = cli.StringFlag{
		Name:  "identity",
		Usage: "Custom node name",
	}
	DocRootFlag = DirectoryFlag{
		Name:  "docroot",
		Usage: "Document Root for HTTPClient file scheme",
		Value: DirectoryString(HomeDir()),
	}
	IterativeOutputFlag = cli.BoolTFlag{
		Name:  "iterative",
		Usage: "Print streaming JSON iteratively, delimited by newlines",
	}
	ExcludeStorageFlag = cli.BoolFlag{
		Name:  "nostorage",
		Usage: "Exclude storage entries (save db lookups)",
	}
	IncludeIncompletesFlag = cli.BoolFlag{
		Name:  "incompletes",
		Usage: "Include accounts for which we don't have the address (missing preimage)",
	}
	ExcludeCodeFlag = cli.BoolFlag{
		Name:  "nocode",
		Usage: "Exclude contract code (save db lookups)",
	}
	StartKeyFlag = cli.StringFlag{
		Name:  "start",
		Usage: "Start position. Either a hash or address",
		Value: "0x0000000000000000000000000000000000000000000000000000000000000000",
	}
	DumpLimitFlag = cli.Uint64Flag{
		Name:  "limit",
		Usage: "Max number of elements (0 = no limit)",
		Value: 0,
	}
	defaultSyncMode = ethconfig.Defaults.SyncMode
	SyncModeFlag    = TextMarshalerFlag{
		Name:  "syncmode",
		Usage: `Blockchain sync mode ("snap", "full")`,
		Value: &defaultSyncMode,
	}
	TxLookupLimitFlag = cli.Uint64Flag{
		Name:  "txlookuplimit",
		Usage: "Number of recent blocks to maintain transactions index by-hash for (default = index all blocks)",
		Value: 0,
	}
	LightKDFFlag = cli.BoolFlag{
		Name:  "lightkdf",
		Usage: "Reduce key-derivation RAM & CPU usage at some expense of KDF strength",
	}
	// Transaction pool settings
	TxPoolLocalsFlag = cli.StringFlag{
		Name:  "txpool.locals",
		Usage: "Comma separated accounts to treat as locals (no flush, priority inclusion)",
	}
	TxPoolNoLocalsFlag = cli.BoolFlag{
		Name:  "txpool.nolocals",
		Usage: "Disables price exemptions for locally submitted transactions",
	}
	TxPoolJournalFlag = cli.StringFlag{
		Name:  "txpool.journal",
		Usage: "Disk journal for local transaction to survive node restarts",
		Value: core.DefaultTxPoolConfig.Journal,
	}
	TxPoolRejournalFlag = cli.DurationFlag{
		Name:  "txpool.rejournal",
		Usage: "Time interval to regenerate the local transaction journal",
		Value: core.DefaultTxPoolConfig.Rejournal,
	}
	TxPoolPriceBumpFlag = cli.Uint64Flag{
		Name:  "txpool.pricebump",
		Usage: "Price bump percentage to replace an already existing transaction",
		Value: ethconfig.Defaults.TxPool.PriceBump,
	}
	TxPoolAccountSlotsFlag = cli.Uint64Flag{
		Name:  "txpool.accountslots",
		Usage: "Minimum number of executable transaction slots guaranteed per account",
		Value: ethconfig.Defaults.TxPool.AccountSlots,
	}
	TxPoolGlobalSlotsFlag = cli.Uint64Flag{
		Name:  "txpool.globalslots",
		Usage: "Maximum number of executable transaction slots for all accounts",
		Value: ethconfig.Defaults.TxPool.GlobalSlots,
	}
	TxPoolAccountQueueFlag = cli.Uint64Flag{
		Name:  "txpool.accountqueue",
		Usage: "Maximum number of non-executable transaction slots permitted per account",
		Value: ethconfig.Defaults.TxPool.AccountQueue,
	}
	TxPoolGlobalQueueFlag = cli.Uint64Flag{
		Name:  "txpool.globalqueue",
		Usage: "Maximum number of non-executable transaction slots for all accounts",
		Value: ethconfig.Defaults.TxPool.GlobalQueue,
	}
	TxPoolGlobalTxCountFlag = cli.Uint64Flag{
		Name:  "txpool.globaltxcount",
		Usage: "Maximum number of transactions for package",
		Value: ethconfig.Defaults.TxPool.GlobalTxCount,
	}
	TxPoolLifetimeFlag = cli.DurationFlag{
		Name:  "txpool.lifetime",
		Usage: "Maximum amount of time non-executable transaction are queued",
		Value: ethconfig.Defaults.TxPool.Lifetime,
	}
	TxPoolCacheSizeFlag = cli.Uint64Flag{
		Name:  "txpool.cacheSize",
		Usage: "After receiving the specified number of transactions from the remote, move the transactions in the queen to pending",
		Value: ethconfig.Defaults.TxPool.TxCacheSize,
	}
	// Performance tuning settings
	CacheFlag = cli.IntFlag{
		Name:  "cache",
		Usage: "Megabytes of memory allocated to internal caching",
		Value: 1024,
	}
	CacheDatabaseFlag = cli.IntFlag{
		Name:  "cache.database",
		Usage: "Percentage of cache memory allowance to use for database io",
		Value: 50,
	}
	CacheTrieFlag = cli.IntFlag{
		Name:  "cache.trie",
		Usage: "Percentage of cache memory allowance to use for trie caching",
		Value: 25,
	}
	CacheTrieJournalFlag = cli.StringFlag{
		Name:  "cache.trie.journal",
		Usage: "Disk journal directory for trie cache to survive node restarts",
		Value: ethconfig.Defaults.TrieCleanCacheJournal,
	}
	CacheTrieRejournalFlag = cli.DurationFlag{
		Name:  "cache.trie.rejournal",
		Usage: "Time interval to regenerate the trie cache journal",
		Value: ethconfig.Defaults.TrieCleanCacheRejournal,
	}
	SnapshotFlag = cli.BoolTFlag{
		Name:  "snapshot",
		Usage: `Enables snapshot-database mode (default = enable)`,
	}
	BloomFilterSizeFlag = cli.Uint64Flag{
		Name:  "bloomfilter.size",
		Usage: "Megabytes of memory allocated to bloom-filter for pruning",
		Value: 2048,
	}
	CacheGCFlag = cli.IntFlag{
		Name:  "cache.gc",
		Usage: "Percentage of cache memory allowance to use for trie pruning (default = 25% full mode, 0% archive mode)",
		Value: 25,
	}
	CacheSnapshotFlag = cli.IntFlag{
		Name:  "cache.snapshot",
		Usage: "Percentage of cache memory allowance to use for snapshot caching (default = 10% full mode, 20% archive mode)",
		Value: 10,
	}
	CacheTrieDBFlag = cli.IntFlag{
		Name:  "cache.triedb",
		Usage: "Megabytes of memory allocated to triedb internal caching",
		Value: ethconfig.Defaults.TrieDBCache,
	}
	CachePreimagesFlag = cli.BoolFlag{
		Name:  "cache.preimages",
		Usage: "Enable recording the SHA3/keccak preimages of trie keys",
	}
	FDLimitFlag = cli.IntFlag{
		Name:  "fdlimit",
		Usage: "Raise the open file descriptor resource limit (default = system fd limit)",
	}
	MinerGasPriceFlag = BigFlag{
		Name:  "miner.gasprice",
		Usage: "Minimum gas price for mining a transaction",
		Value: ethconfig.Defaults.Miner.GasPrice,
	}
	// Account settings
	UnlockedAccountFlag = cli.StringFlag{
		Name:  "unlock",
		Usage: "Comma separated list of accounts to unlock",
		Value: "",
	}
	PasswordFileFlag = cli.StringFlag{
		Name:  "password",
		Usage: "Password file to use for non-interactive password input",
		Value: "",
	}
	InsecureUnlockAllowedFlag = cli.BoolFlag{
		Name:  "allow-insecure-unlock",
		Usage: "Allow insecure account unlocking when account-related RPCs are exposed by http",
	}
	RPCGlobalGasCapFlag = cli.Uint64Flag{
		Name:  "rpc.gascap",
		Usage: "Sets a cap on gas that can be used in platon_call/estimateGas (0=infinite)",
		Value: ethconfig.Defaults.RPCGasCap,
	}
	RPCGlobalEVMTimeoutFlag = cli.DurationFlag{
		Name:  "rpc.evmtimeout",
		Usage: "Sets a timeout used for eth_call (0=infinite)",
		Value: ethconfig.Defaults.RPCEVMTimeout,
	}
	RPCGlobalTxFeeCapFlag = cli.Float64Flag{
		Name:  "rpc.txfeecap",
		Usage: "Sets a cap on transaction fee (in ether) that can be sent via the RPC APIs (0 = no cap)",
		Value: ethconfig.Defaults.RPCTxFeeCap,
	}
	// Authenticated RPC HTTP settings
	AuthHostFlag = cli.StringFlag{
		Name:  "authrpc.host",
		Usage: "Listening address for authenticated APIs",
		Value: node.DefaultConfig.AuthHost,
	}
	AuthPortFlag = cli.IntFlag{
		Name:  "authrpc.port",
		Usage: "Listening port for authenticated APIs",
		Value: node.DefaultConfig.AuthPort,
	}
	JWTSecretFlag = cli.StringFlag{
		Name:  "authrpc.jwtsecret",
		Usage: "Path to a JWT secret to use for authenticated RPC endpoints",
	}
	// Logging and debug settings
	EthStatsURLFlag = cli.StringFlag{
		Name:  "ethstats",
		Usage: "Reporting URL of a ethstats service (nodename:secret@host:port)",
	}
	NoCompactionFlag = cli.BoolFlag{
		Name:  "nocompaction",
		Usage: "Disables db compaction after import",
	}
	IgnoreLegacyReceiptsFlag = cli.BoolFlag{
		Name:  "ignore-legacy-receipts",
		Usage: "Platon will start up even if there are legacy receipts in freezer",
	}
	// RPC settings
	IPCDisabledFlag = cli.BoolFlag{
		Name:  "ipcdisable",
		Usage: "Disable the IPC-RPC server",
	}
	IPCPathFlag = DirectoryFlag{
		Name:  "ipcpath",
		Usage: "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
	}
	HTTPEnabledFlag = cli.BoolFlag{
		Name:  "http",
		Usage: "Enable the HTTP-RPC server",
	}
	HTTPListenAddrFlag = cli.StringFlag{
		Name:  "http.addr",
		Usage: "HTTP-RPC server listening interface",
		Value: node.DefaultHTTPHost,
	}
	HTTPPortFlag = cli.IntFlag{
		Name:  "http.port",
		Usage: "HTTP-RPC server listening port",
		Value: node.DefaultHTTPPort,
	}
	HTTPCORSDomainFlag = cli.StringFlag{
		Name:  "http.corsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests (browser enforced)",
		Value: "",
	}
	HTTPVirtualHostsFlag = cli.StringFlag{
		Name:  "http.vhosts",
		Usage: "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard.",
		Value: strings.Join(node.DefaultConfig.HTTPVirtualHosts, ","),
	}
	HTTPApiFlag = cli.StringFlag{
		Name:  "http.api",
		Usage: "API's offered over the HTTP-RPC interface",
		Value: "",
	}
	HTTPPathPrefixFlag = cli.StringFlag{
		Name:  "http.rpcprefix",
		Usage: "HTTP path path prefix on which JSON-RPC is served. Use '/' to serve on all paths.",
		Value: "",
	}
	HTTPEnabledEthCompatibleFlag = cli.BoolFlag{
		Name:  "http.ethcompatible",
		Usage: "Enable eth compatible",
	}
	GraphQLEnabledFlag = cli.BoolFlag{
		Name:  "graphql",
		Usage: "Enable GraphQL on the HTTP-RPC server. Note that GraphQL can only be started if an HTTP server is started as well.",
	}
	GraphQLCORSDomainFlag = cli.StringFlag{
		Name:  "graphql.corsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests (browser enforced)",
		Value: "",
	}
	GraphQLVirtualHostsFlag = cli.StringFlag{
		Name:  "graphql.vhosts",
		Usage: "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard.",
		Value: strings.Join(node.DefaultConfig.GraphQLVirtualHosts, ","),
	}
	WSEnabledFlag = cli.BoolFlag{
		Name:  "ws",
		Usage: "Enable the WS-RPC server",
	}
	WSListenAddrFlag = cli.StringFlag{
		Name:  "ws.addr",
		Usage: "WS-RPC server listening interface",
		Value: node.DefaultWSHost,
	}
	WSPortFlag = cli.IntFlag{
		Name:  "ws.port",
		Usage: "WS-RPC server listening port",
		Value: node.DefaultWSPort,
	}
	WSApiFlag = cli.StringFlag{
		Name:  "ws.api",
		Usage: "API's offered over the WS-RPC interface",
		Value: "",
	}
	WSAllowedOriginsFlag = cli.StringFlag{
		Name:  "ws.origins",
		Usage: "Origins from which to accept websockets requests",
		Value: "",
	}
	WSPathPrefixFlag = cli.StringFlag{
		Name:  "ws.rpcprefix",
		Usage: "HTTP path prefix on which JSON-RPC is served. Use '/' to serve on all paths.",
		Value: "",
	}
	ExecFlag = cli.StringFlag{
		Name:  "exec",
		Usage: "Execute JavaScript statement",
	}
	PreloadJSFlag = cli.StringFlag{
		Name:  "preload",
		Usage: "Comma separated list of JavaScript files to preload into the console",
	}
	AllowUnprotectedTxs = &cli.BoolFlag{
		Name:  "rpc.allow-unprotected-txs",
		Usage: "Allow for unprotected (non EIP155 signed) transactions to be submitted via RPC",
	}

	// Network Settings
	MaxPeersFlag = cli.IntFlag{
		Name:  "maxpeers",
		Usage: "Maximum number of network peers (network disabled if set to 0)",
		Value: node.DefaultConfig.P2P.MaxPeers,
	}
	MaxConsensusPeersFlag = cli.IntFlag{
		Name:  "maxconsensuspeers",
		Usage: "Maximum number of network consensus peers (network disabled if set to 0)",
		Value: 40,
	}
	MaxPendingPeersFlag = cli.IntFlag{
		Name:  "maxpendpeers",
		Usage: "Maximum number of pending connection attempts (defaults used if set to 0)",
		Value: node.DefaultConfig.P2P.MaxPendingPeers,
	}
	ListenPortFlag = cli.IntFlag{
		Name:  "port",
		Usage: "Network listening port",
		Value: 16789,
	}
	BootnodesFlag = cli.StringFlag{
		Name:  "bootnodes",
		Usage: "Comma separated enode URLs for P2P discovery bootstrap (set v4+v5 instead for light servers)",
		Value: "",
	}
	BootnodesV4Flag = cli.StringFlag{
		Name:  "bootnodesv4",
		Usage: "Comma separated enode URLs for P2P v4 discovery bootstrap (light server, full nodes)",
		Value: "",
	}
	/*
		BootnodesV5Flag = cli.StringFlag{
			Name:  "bootnodesv5",
			Usage: "Comma separated enode URLs for P2P v5 discovery bootstrap (light server, light nodes)",
			Value: "",
		}*/
	NodeKeyFileFlag = cli.StringFlag{
		Name:  "nodekey",
		Usage: "P2P node key file",
	}
	NodeKeyHexFlag = cli.StringFlag{
		Name:  "nodekeyhex",
		Usage: "P2P node key as hex (for testing)",
	}
	NATFlag = cli.StringFlag{
		Name:  "nat",
		Usage: "NAT port mapping mechanism (any|none|upnp|pmp|extip:<IP>)",
		Value: "any",
	}
	NoDiscoverFlag = cli.BoolFlag{
		Name:  "nodiscover",
		Usage: "Disables the peer discovery mechanism (manual peer addition)",
	}
	/*
		DiscoveryV5Flag = cli.BoolFlag{
			Name:  "v5disc",
			Usage: "Enables the experimental RLPx V5 (Topic Discovery) mechanism",
		}*/
	NetrestrictFlag = cli.StringFlag{
		Name:  "netrestrict",
		Usage: "Restricts network communication to the given IP networks (CIDR masks)",
	}
	DNSDiscoveryFlag = cli.StringFlag{
		Name:  "discovery.dns",
		Usage: "Sets DNS discovery entry points (use \"\" to disable DNS)",
	}

	// ATM the url is left to the user and deployment to
	JSpathFlag = DirectoryFlag{
		Name:  "jspath",
		Usage: "JavaScript root path for `loadScript`",
		Value: DirectoryString("."),
	}

	// Gas price oracle settings
	GpoBlocksFlag = cli.IntFlag{
		Name:  "gpo.blocks",
		Usage: "Number of recent blocks to check for gas prices",
		Value: ethconfig.Defaults.GPO.Blocks,
	}
	GpoPercentileFlag = cli.IntFlag{
		Name:  "gpo.percentile",
		Usage: "Suggested gas price is the given percentile of a set of recent transaction gas prices",
		Value: ethconfig.Defaults.GPO.Percentile,
	}
	GpoMaxGasPriceFlag = cli.Int64Flag{
		Name:  "gpo.maxprice",
		Usage: "Maximum transaction priority fee (or gasprice before London fork) to be recommended by gpo",
		Value: ethconfig.Defaults.GPO.MaxPrice.Int64(),
	}
	GpoIgnoreGasPriceFlag = cli.Int64Flag{
		Name:  "gpo.ignoreprice",
		Usage: "Gas price below which gpo will ignore transactions",
		Value: ethconfig.Defaults.GPO.IgnorePrice.Int64(),
	}

	// Metrics flags
	MetricsEnabledFlag = cli.BoolFlag{
		Name:  "metrics",
		Usage: "Enable metrics collection and reporting",
	}
	MetricsEnabledExpensiveFlag = cli.BoolFlag{
		Name:  "metrics.expensive",
		Usage: "Enable expensive metrics collection and reporting",
	}

	// MetricsHTTPFlag defines the endpoint for a stand-alone metrics HTTP endpoint.
	// Since the pprof service enables sensitive/vulnerable behavior, this allows a user
	// to enable a public-OK metrics endpoint without having to worry about ALSO exposing
	// other profiling behavior or information.
	MetricsHTTPFlag = cli.StringFlag{
		Name:  "metrics.addr",
		Usage: "Enable stand-alone metrics HTTP server listening interface",
		Value: metrics.DefaultConfig.HTTP,
	}
	MetricsPortFlag = cli.IntFlag{
		Name:  "metrics.port",
		Usage: "Metrics HTTP server listening port",
		Value: metrics.DefaultConfig.Port,
	}
	MetricsEnableInfluxDBFlag = cli.BoolFlag{
		Name:  "metrics.influxdb",
		Usage: "Enable metrics export/push to an external InfluxDB database",
	}
	MetricsInfluxDBEndpointFlag = cli.StringFlag{
		Name:  "metrics.influxdb.endpoint",
		Usage: "InfluxDB API endpoint to report metrics to",
		Value: metrics.DefaultConfig.InfluxDBEndpoint,
	}
	MetricsInfluxDBDatabaseFlag = cli.StringFlag{
		Name:  "metrics.influxdb.database",
		Usage: "InfluxDB database name to push reported metrics to",
		Value: metrics.DefaultConfig.InfluxDBDatabase,
	}
	MetricsInfluxDBUsernameFlag = cli.StringFlag{
		Name:  "metrics.influxdb.username",
		Usage: "Username to authorize access to the database",
		Value: metrics.DefaultConfig.InfluxDBUsername,
	}
	MetricsInfluxDBPasswordFlag = cli.StringFlag{
		Name:  "metrics.influxdb.password",
		Usage: "Password to authorize access to the database",
		Value: metrics.DefaultConfig.InfluxDBPassword,
	}
	// Tags are part of every measurement sent to InfluxDB. Queries on tags are faster in InfluxDB.
	// For example `host` tag could be used so that we can group all nodes and average a measurement
	// across all of them, but also so that we can select a specific node and inspect its measurements.
	// https://docs.influxdata.com/influxdb/v1.4/concepts/key_concepts/#tag-key
	MetricsInfluxDBTagsFlag = cli.StringFlag{
		Name:  "metrics.influxdb.tags",
		Usage: "Comma-separated InfluxDB tags (key/values) attached to all measurements",
		Value: metrics.DefaultConfig.InfluxDBTags,
	}

	MetricsEnableInfluxDBV2Flag = cli.BoolFlag{
		Name:  "metrics.influxdbv2",
		Usage: "Enable metrics export/push to an external InfluxDB v2 database",
	}

	MetricsInfluxDBTokenFlag = cli.StringFlag{
		Name:  "metrics.influxdb.token",
		Usage: "Token to authorize access to the database (v2 only)",
		Value: metrics.DefaultConfig.InfluxDBToken,
	}

	MetricsInfluxDBBucketFlag = cli.StringFlag{
		Name:  "metrics.influxdb.bucket",
		Usage: "InfluxDB bucket name to push reported metrics to (v2 only)",
		Value: metrics.DefaultConfig.InfluxDBBucket,
	}

	MetricsInfluxDBOrganizationFlag = cli.StringFlag{
		Name:  "metrics.influxdb.organization",
		Usage: "InfluxDB organization name (v2 only)",
		Value: metrics.DefaultConfig.InfluxDBOrganization,
	}

	CbftPeerMsgQueueSize = cli.Uint64Flag{
		Name:  "cbft.msg_queue_size",
		Usage: "Message queue size",
		Value: 1024,
	}

	CbftWalDisabledFlag = cli.BoolFlag{
		Name:  "cbft.wal.disabled",
		Usage: "Disable the Wal server",
	}

	CbftMaxPingLatency = cli.Int64Flag{
		Name:  "cbft.max_ping_latency",
		Usage: "Maximum latency of ping",
		Value: 2000,
	}

	CbftBlsPriKeyFileFlag = cli.StringFlag{
		Name:  "cbft.blskey",
		Usage: "BLS key file",
	}

	CbftBlacklistDeadlineFlag = cli.StringFlag{
		Name:  "cbft.blacklist_deadline",
		Usage: "Blacklist effective time. uint:minute",
		Value: "60",
	}

	DBNoGCFlag = cli.BoolFlag{
		Name:  "db.nogc",
		Usage: "Disables database garbage collection",
	}
	DBGCIntervalFlag = cli.Uint64Flag{
		Name:  "db.gc_interval",
		Usage: "Block interval for garbage collection",
		Value: ethconfig.Defaults.DBGCInterval,
	}
	DBGCTimeoutFlag = cli.DurationFlag{
		Name:  "db.gc_timeout",
		Usage: "Maximum time for database garbage collection",
		Value: ethconfig.Defaults.DBGCTimeout,
	}
	DBGCMptFlag = cli.BoolFlag{
		Name:  "db.gc_mpt",
		Usage: "Enables database garbage collection MPT",
	}
	DBGCBlockFlag = cli.IntFlag{
		Name:  "db.gc_block",
		Usage: "Number of cache block states, default 10",
		Value: ethconfig.Defaults.DBGCBlock,
	}
	DBValidatorsHistoryFlag = cli.BoolFlag{
		Name:  "db.validators_history",
		Usage: "Store the list of validators for each consensus round",
	}

	VMWasmType = cli.StringFlag{
		Name:   "vm.wasm_type",
		Usage:  "The actual implementation type of the wasm instance",
		EnvVar: "",
		Value:  ethconfig.Defaults.VMWasmType,
	}

	VmTimeoutDuration = cli.Uint64Flag{
		Name:   "vm.timeout_duration",
		Usage:  "The VM execution timeout duration (uint: ms)",
		EnvVar: "",
		Value:  ethconfig.Defaults.VmTimeoutDuration,
	}
)

var (
	// DatabasePathFlags is the flag group of all database path flags.
	DatabasePathFlags = []cli.Flag{
		DataDirFlag,
		AncientFlag,
		RemoteDBFlag,
	}
)

// GroupFlags combines the given flag slices together and returns the merged one.
func GroupFlags(groups ...[]cli.Flag) []cli.Flag {
	var ret []cli.Flag
	for _, group := range groups {
		ret = append(ret, group...)
	}
	return ret
}

// MakeDataDir retrieves the currently requested data directory, terminating
// if none (or the empty string) is specified. If the node is starting a testnet,
// then a subdirectory of the specified datadir will be used.
func MakeDataDir(ctx *cli.Context) string {
	if path := ctx.GlobalString(DataDirFlag.Name); path != "" {

		if ctx.GlobalBool(TestnetFlag.Name) {
			return filepath.Join(path, "testnet")
		}
		return path
	}
	Fatalf("Cannot determine default data directory, please set manually (--datadir)")
	return ""
}

// setNodeKey creates a node key from set command line flags, either loading it
// from a file or as a specified hex value. If neither flags were provided, this
// method returns nil and an emphemeral key is to be generated.
func setNodeKey(ctx *cli.Context, cfg *p2p.Config) {
	var (
		hex  = ctx.GlobalString(NodeKeyHexFlag.Name)
		file = ctx.GlobalString(NodeKeyFileFlag.Name)
		key  *ecdsa.PrivateKey
		err  error
	)
	switch {
	case file != "" && hex != "":
		Fatalf("Options %q and %q are mutually exclusive", NodeKeyFileFlag.Name, NodeKeyHexFlag.Name)
	case file != "":
		if key, err = crypto.LoadECDSA(file); err != nil {
			Fatalf("Option %q: %v", NodeKeyFileFlag.Name, err)
		}
		cfg.PrivateKey = key
	case hex != "":
		if key, err = crypto.HexToECDSA(hex); err != nil {
			Fatalf("Option %q: %v", NodeKeyHexFlag.Name, err)
		}
		cfg.PrivateKey = key
	}
}

// setNodeUserIdent creates the user identifier from CLI flags.
func setNodeUserIdent(ctx *cli.Context, cfg *node.Config) {
	if identity := ctx.GlobalString(IdentityFlag.Name); len(identity) > 0 {
		cfg.UserIdent = identity
	}
}

// setBootstrapNodes creates a list of bootstrap nodes from the command line
// flags, reverting to pre-configured ones if none have been specified.
func setBootstrapNodes(ctx *cli.Context, cfg *p2p.Config) {
	urls := params.MainnetBootnodes

	switch {
	case ctx.GlobalIsSet(BootnodesFlag.Name) || ctx.GlobalIsSet(BootnodesV4Flag.Name):
		if ctx.GlobalIsSet(BootnodesV4Flag.Name) {
			urls = SplitAndTrim(ctx.GlobalString(BootnodesV4Flag.Name))
		} else {
			urls = SplitAndTrim(ctx.GlobalString(BootnodesFlag.Name))
		}
	case ctx.GlobalBool(TestnetFlag.Name):
		urls = params.TestnetBootnodes
	case cfg.BootstrapNodes != nil:
		return // already set, don't apply defaults.
	}

	cfg.BootstrapNodes = make([]*enode.Node, 0, len(urls))
	for _, url := range urls {
		if url != "" {
			node, err := enode.Parse(enode.ValidSchemes, url)
			if err != nil {
				log.Crit("Bootstrap URL invalid", "enode", url, "err", err)
			}
			cfg.BootstrapNodes = append(cfg.BootstrapNodes, node)
		}
	}
}

// setBootstrapNodesV5 creates a list of bootstrap nodes from the command line
// flags, reverting to pre-configured ones if none have been specified.
/*
func setBootstrapNodesV5(ctx *cli.Context, cfg *p2p.Config) {
	urls := params.DiscoveryV5Bootnodes
	switch {
	case ctx.GlobalIsSet(BootnodesFlag.Name) || ctx.GlobalIsSet(BootnodesV5Flag.Name):
		if ctx.GlobalIsSet(BootnodesV5Flag.Name) {
			urls = splitAndTrim(ctx.GlobalString(BootnodesV5Flag.Name))
		} else {
			urls = splitAndTrim(ctx.GlobalString(BootnodesV5Flag.Name))
		}
	case cfg.BootstrapNodesV5 != nil:
		return // already set, don't apply defaults.
	}

	cfg.BootstrapNodesV5 = make([]*discv5.Node, 0, len(urls))
	for _, url := range urls {
		node, err := discv5.ParseNode(url)
		if err != nil {
			log.Error("Bootstrap URL invalid", "enode", url, "err", err)
			continue
		}
		cfg.BootstrapNodesV5 = append(cfg.BootstrapNodesV5, node)
	}
}*/

// setListenAddress creates a TCP listening address string from set command
// line flags.
func setListenAddress(ctx *cli.Context, cfg *p2p.Config) {
	if ctx.GlobalIsSet(ListenPortFlag.Name) {
		cfg.ListenAddr = fmt.Sprintf(":%d", ctx.GlobalInt(ListenPortFlag.Name))
	}
}

// setNAT creates a port mapper from command line flags.
func setNAT(ctx *cli.Context, cfg *p2p.Config) {
	if ctx.GlobalIsSet(NATFlag.Name) {
		natif, err := nat.Parse(ctx.GlobalString(NATFlag.Name))
		if err != nil {
			Fatalf("Option %s: %v", NATFlag.Name, err)
		}
		cfg.NAT = natif
	}
}

// SplitAndTrim splits input separated by a comma
// and trims excessive white space from the substrings.
func SplitAndTrim(input string) (ret []string) {
	l := strings.Split(input, ",")
	for _, r := range l {
		r = strings.TrimSpace(r)
		if len(r) > 0 {
			ret = append(ret, r)
		}
	}
	return ret
}

// setHTTP creates the HTTP RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setHTTP(ctx *cli.Context, cfg *node.Config) {
	if ctx.GlobalBool(HTTPEnabledFlag.Name) && cfg.HTTPHost == "" {
		cfg.HTTPHost = "127.0.0.1"
		if ctx.GlobalIsSet(HTTPListenAddrFlag.Name) {
			cfg.HTTPHost = ctx.GlobalString(HTTPListenAddrFlag.Name)
		}
	}
	if ctx.GlobalIsSet(HTTPPortFlag.Name) {
		cfg.HTTPPort = ctx.GlobalInt(HTTPPortFlag.Name)
	}
	if ctx.GlobalIsSet(AuthHostFlag.Name) {
		cfg.AuthHost = ctx.GlobalString(AuthHostFlag.Name)
	}
	if ctx.GlobalIsSet(AuthPortFlag.Name) {
		cfg.AuthPort = ctx.GlobalInt(AuthPortFlag.Name)
	}

	if ctx.GlobalIsSet(HTTPCORSDomainFlag.Name) {
		cfg.HTTPCors = SplitAndTrim(ctx.GlobalString(HTTPCORSDomainFlag.Name))
	}
	if ctx.GlobalIsSet(HTTPApiFlag.Name) {
		cfg.HTTPModules = SplitAndTrim(ctx.GlobalString(HTTPApiFlag.Name))
	}
	if ctx.GlobalBool(HTTPEnabledEthCompatibleFlag.Name) {
		types2.HttpEthCompatible = true
	}
	if ctx.GlobalIsSet(HTTPVirtualHostsFlag.Name) {
		cfg.HTTPVirtualHosts = SplitAndTrim(ctx.GlobalString(HTTPVirtualHostsFlag.Name))
	}
	if ctx.IsSet(AllowUnprotectedTxs.Name) {
		cfg.AllowUnprotectedTxs = ctx.Bool(AllowUnprotectedTxs.Name)
	}

	if ctx.GlobalIsSet(HTTPPathPrefixFlag.Name) {
		cfg.HTTPPathPrefix = ctx.GlobalString(HTTPPathPrefixFlag.Name)
	}
}

// setGraphQL creates the GraphQL listener interface string from the set
// command line flags, returning empty if the GraphQL endpoint is disabled.
func setGraphQL(ctx *cli.Context, cfg *node.Config) {
	if ctx.GlobalIsSet(GraphQLCORSDomainFlag.Name) {
		cfg.GraphQLCors = SplitAndTrim(ctx.GlobalString(GraphQLCORSDomainFlag.Name))
	}
	if ctx.GlobalIsSet(GraphQLVirtualHostsFlag.Name) {
		cfg.GraphQLVirtualHosts = SplitAndTrim(ctx.GlobalString(GraphQLVirtualHostsFlag.Name))
	}
}

// setWS creates the WebSocket RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setWS(ctx *cli.Context, cfg *node.Config) {
	if ctx.GlobalBool(WSEnabledFlag.Name) && cfg.WSHost == "" {
		cfg.WSHost = "127.0.0.1"
		if ctx.GlobalIsSet(LegacyWSListenAddrFlag.Name) {
			cfg.WSHost = ctx.GlobalString(LegacyWSListenAddrFlag.Name)
			log.Warn("The flag --wsaddr is deprecated and will be removed in the future, please use --ws.addr")
		}
		if ctx.GlobalIsSet(WSListenAddrFlag.Name) {
			cfg.WSHost = ctx.GlobalString(WSListenAddrFlag.Name)
		}
	}
	if ctx.GlobalIsSet(LegacyWSPortFlag.Name) {
		cfg.WSPort = ctx.GlobalInt(LegacyWSPortFlag.Name)
		log.Warn("The flag --wsport is deprecated and will be removed in the future, please use --ws.port")
	}
	if ctx.GlobalIsSet(WSPortFlag.Name) {
		cfg.WSPort = ctx.GlobalInt(WSPortFlag.Name)
	}

	if ctx.GlobalIsSet(LegacyWSAllowedOriginsFlag.Name) {
		cfg.WSOrigins = SplitAndTrim(ctx.GlobalString(LegacyWSAllowedOriginsFlag.Name))
		log.Warn("The flag --wsorigins is deprecated and will be removed in the future, please use --ws.origins")
	}
	if ctx.GlobalIsSet(WSAllowedOriginsFlag.Name) {
		cfg.WSOrigins = SplitAndTrim(ctx.GlobalString(WSAllowedOriginsFlag.Name))
	}

	if ctx.GlobalIsSet(LegacyWSApiFlag.Name) {
		cfg.WSModules = SplitAndTrim(ctx.GlobalString(LegacyWSApiFlag.Name))
		log.Warn("The flag --wsapi is deprecated and will be removed in the future, please use --ws.api")
	}
	if ctx.GlobalIsSet(WSApiFlag.Name) {
		cfg.WSModules = SplitAndTrim(ctx.GlobalString(WSApiFlag.Name))
	}

	if ctx.GlobalIsSet(WSPathPrefixFlag.Name) {
		cfg.WSPathPrefix = ctx.GlobalString(WSPathPrefixFlag.Name)
	}
}

// setIPC creates an IPC path configuration from the set command line flags,
// returning an empty string if IPC was explicitly disabled, or the set path.
func setIPC(ctx *cli.Context, cfg *node.Config) {
	CheckExclusive(ctx, IPCDisabledFlag, IPCPathFlag)
	switch {
	case ctx.GlobalBool(IPCDisabledFlag.Name):
		cfg.IPCPath = ""
	case ctx.GlobalIsSet(IPCPathFlag.Name):
		cfg.IPCPath = ctx.GlobalString(IPCPathFlag.Name)
	}
}

// MakeDatabaseHandles raises out the number of allowed file handles per process
// for Geth and returns half of the allowance to assign to the database.
func MakeDatabaseHandles(max int) int {
	limit, err := fdlimit.Maximum()
	if err != nil {
		Fatalf("Failed to retrieve file descriptor allowance: %v", err)
	}
	switch {
	case max == 0:
		// User didn't specify a meaningful value, use system limits
	case max < 128:
		// User specified something unhealthy, just use system defaults
		log.Error("File descriptor limit invalid (<128)", "had", max, "updated", limit)
	case max > limit:
		// User requested more than the OS allows, notify that we can't allocate it
		log.Warn("Requested file descriptors denied by OS", "req", max, "limit", limit)
	default:
		// User limit is meaningful and within allowed range, use that
		limit = max
	}
	raised, err := fdlimit.Raise(uint64(limit))
	if err != nil {
		Fatalf("Failed to raise file descriptor allowance: %v", err)
	}
	return int(raised / 2) // Leave half for networking and other stuff
}

// MakeAddress converts an account specified directly as a hex encoded string or
// a key index in the key store to an internal account representation.
func MakeAddress(ks *keystore.KeyStore, account string) (accounts.Account, error) {
	// If the specified account is a valid address, return it
	if common.IsBech32Address(account) {
		return accounts.Account{Address: common.MustBech32ToAddress(account)}, nil
	}
	if common.IsHexAddress(account) {
		return accounts.Account{Address: common.HexToAddress(account)}, nil
	}
	// Otherwise try to interpret the account as a keystore index
	index, err := strconv.Atoi(account)
	if err != nil || index < 0 {
		return accounts.Account{}, fmt.Errorf("invalid account address or index %q", account)
	}
	log.Warn("-------------------------------------------------------------------")
	log.Warn("Referring to accounts by order in the keystore folder is dangerous!")
	log.Warn("This functionality is deprecated and will be removed in the future!")
	log.Warn("Please use explicit addresses! (can search via `platon account list`)")
	log.Warn("-------------------------------------------------------------------")

	accs := ks.Accounts()
	if len(accs) <= index {
		return accounts.Account{}, fmt.Errorf("index %d higher than number of accounts %d", index, len(accs))
	}
	return accs[index], nil
}

// MakePasswordList reads password lines from the file specified by the global --password flag.
func MakePasswordList(ctx *cli.Context) []string {
	path := ctx.GlobalString(PasswordFileFlag.Name)
	if path == "" {
		return nil
	}
	text, err := os.ReadFile(path)
	if err != nil {
		Fatalf("Failed to read password file: %v", err)
	}
	lines := strings.Split(string(text), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines
}

func SetP2PConfig(ctx *cli.Context, cfg *p2p.Config) {
	setNodeKey(ctx, cfg)
	setNAT(ctx, cfg)
	setListenAddress(ctx, cfg)
	setBootstrapNodes(ctx, cfg)
	// setBootstrapNodesV5(ctx, cfg)

	if ctx.GlobalIsSet(MaxConsensusPeersFlag.Name) {
		cfg.MaxConsensusPeers = ctx.GlobalInt(MaxConsensusPeersFlag.Name)
	}

	if ctx.GlobalIsSet(MaxPeersFlag.Name) {
		cfg.MaxPeers = ctx.GlobalInt(MaxPeersFlag.Name)
	}

	ethPeers := cfg.MaxPeers

	if cfg.MaxPeers <= cfg.MaxConsensusPeers {
		log.Error("MaxPeers is less than MaxConsensusPeers", "MaxPeers", cfg.MaxPeers, "MaxConsensusPeers", cfg.MaxConsensusPeers)
		Fatalf("MaxPeers is less than MaxConsensusPeers, MaxPeers: %d, MaxConsensusPeers: %d", cfg.MaxPeers, cfg.MaxConsensusPeers)
	}
	log.Info("Maximum peer count", "ETH", ethPeers, "consensusTotal", cfg.MaxConsensusPeers, "total", cfg.MaxPeers)

	if ctx.GlobalIsSet(MaxPendingPeersFlag.Name) {
		cfg.MaxPendingPeers = ctx.GlobalInt(MaxPendingPeersFlag.Name)
	}
	if ctx.GlobalIsSet(NoDiscoverFlag.Name) {
		cfg.NoDiscovery = true
	}

	// if we're running a light client or server, force enable the v5 peer discovery
	// unless it is explicitly disabled with --nodiscover note that explicitly specifying
	// --v5disc overrides --nodiscover, in which case the later only disables v4 discovery
	//forceV5Discovery := /*(lightClient || lightServer) && */ !ctx.GlobalBool(NoDiscoverFlag.Name)
	/*if ctx.GlobalIsSet(DiscoveryV5Flag.Name) {
		cfg.DiscoveryV5 = ctx.GlobalBool(DiscoveryV5Flag.Name)
	} else if forceV5Discovery {
		cfg.DiscoveryV5 = true
	}*/
	// 没用discoverv5

	if netrestrict := ctx.GlobalString(NetrestrictFlag.Name); netrestrict != "" {
		list, err := netutil.ParseNetlist(netrestrict)
		if err != nil {
			Fatalf("Option %q: %v", NetrestrictFlag.Name, err)
		}
		cfg.NetRestrict = list
	}
}

// SetNodeConfig applies node-related command line flags to the config.
func SetNodeConfig(ctx *cli.Context, cfg *node.Config) {
	SetP2PConfig(ctx, &cfg.P2P)
	setIPC(ctx, cfg)
	setHTTP(ctx, cfg)
	setGraphQL(ctx, cfg)
	setWS(ctx, cfg)
	setNodeUserIdent(ctx, cfg)

	switch {
	case ctx.GlobalIsSet(DataDirFlag.Name):
		cfg.DataDir = ctx.GlobalString(DataDirFlag.Name)
	case ctx.GlobalBool(TestnetFlag.Name):
		cfg.DataDir = filepath.Join(node.DefaultDataDir(), "testnet")
	}

	if ctx.GlobalIsSet(JWTSecretFlag.Name) {
		cfg.JWTSecret = ctx.GlobalString(JWTSecretFlag.Name)
	}

	if ctx.GlobalIsSet(KeyStoreDirFlag.Name) {
		cfg.KeyStoreDir = ctx.GlobalString(KeyStoreDirFlag.Name)
	}
	if ctx.GlobalIsSet(LightKDFFlag.Name) {
		cfg.UseLightweightKDF = ctx.GlobalBool(LightKDFFlag.Name)
	}
	if ctx.GlobalIsSet(InsecureUnlockAllowedFlag.Name) {
		cfg.InsecureUnlockAllowed = ctx.GlobalBool(InsecureUnlockAllowedFlag.Name)
	}
}

func setGPO(ctx *cli.Context, cfg *gasprice.Config) {
	if ctx.GlobalIsSet(GpoBlocksFlag.Name) {
		cfg.Blocks = ctx.GlobalInt(GpoBlocksFlag.Name)
	}
	if ctx.GlobalIsSet(GpoPercentileFlag.Name) {
		cfg.Percentile = ctx.GlobalInt(GpoPercentileFlag.Name)
	}
	if ctx.GlobalIsSet(GpoMaxGasPriceFlag.Name) {
		cfg.MaxPrice = big.NewInt(ctx.GlobalInt64(GpoMaxGasPriceFlag.Name))
	}
	if ctx.GlobalIsSet(GpoIgnoreGasPriceFlag.Name) {
		cfg.IgnorePrice = big.NewInt(ctx.GlobalInt64(GpoIgnoreGasPriceFlag.Name))
	}
}

func setTxPool(ctx *cli.Context, cfg *core.TxPoolConfig) {
	if ctx.GlobalIsSet(TxPoolLocalsFlag.Name) {
		locals := strings.Split(ctx.GlobalString(TxPoolLocalsFlag.Name), ",")
		for _, account := range locals {
			if trimmed := strings.TrimSpace(account); !common.IsBech32Address(trimmed) {
				Fatalf("Invalid account in --txpool.locals: %s", trimmed)
			} else {
				cfg.Locals = append(cfg.Locals, common.MustBech32ToAddress(account))
			}
		}
	}
	if ctx.GlobalIsSet(TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.GlobalBool(TxPoolNoLocalsFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolJournalFlag.Name) {
		cfg.Journal = ctx.GlobalString(TxPoolJournalFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolRejournalFlag.Name) {
		cfg.Rejournal = ctx.GlobalDuration(TxPoolRejournalFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolPriceBumpFlag.Name) {
		cfg.PriceBump = ctx.GlobalUint64(TxPoolPriceBumpFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolAccountSlotsFlag.Name) {
		cfg.AccountSlots = ctx.GlobalUint64(TxPoolAccountSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolGlobalSlotsFlag.Name) {
		cfg.GlobalSlots = ctx.GlobalUint64(TxPoolGlobalSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolAccountQueueFlag.Name) {
		cfg.AccountQueue = ctx.GlobalUint64(TxPoolAccountQueueFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolGlobalQueueFlag.Name) {
		cfg.GlobalQueue = ctx.GlobalUint64(TxPoolGlobalQueueFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolGlobalTxCountFlag.Name) {
		cfg.GlobalTxCount = ctx.GlobalUint64(TxPoolGlobalTxCountFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.GlobalDuration(TxPoolLifetimeFlag.Name)
	}
	if ctx.GlobalIsSet(TxPoolCacheSizeFlag.Name) {
		cfg.TxCacheSize = ctx.GlobalUint64(TxPoolCacheSizeFlag.Name)
	}
}

func setMiner(ctx *cli.Context, cfg *miner.Config) {
	if ctx.GlobalIsSet(MinerGasPriceFlag.Name) {
		cfg.GasPrice = GlobalBig(ctx, MinerGasPriceFlag.Name)
	}
}

// SetEthConfig applies eth-related command line flags to the config.
func SetEthConfig(ctx *cli.Context, stack *node.Node, cfg *ethconfig.Config) {
	// Avoid conflicting network flags
	CheckExclusive(ctx, TestnetFlag)

	setGPO(ctx, &cfg.GPO)
	setTxPool(ctx, &cfg.TxPool)

	// Cap the cache allowance and tune the garbage collector
	mem, err := gopsutil.VirtualMemory()
	if err == nil {
		if 32<<(^uintptr(0)>>63) == 32 && mem.Total > 2*1024*1024*1024 {
			log.Warn("Lowering memory allowance on 32bit arch", "available", mem.Total/1024/1024, "addressable", 2*1024)
			mem.Total = 2 * 1024 * 1024 * 1024
		}
		allowance := int(mem.Total / 1024 / 1024 / 3)
		if cache := ctx.GlobalInt(CacheFlag.Name); cache > allowance {
			log.Warn("Sanitizing cache to Go's GC limits", "provided", cache, "updated", allowance)
			ctx.GlobalSet(CacheFlag.Name, strconv.Itoa(allowance))
		}
	}
	// Ensure Go's GC ignores the database cache for trigger percentage
	cache := ctx.GlobalInt(CacheFlag.Name)
	gogc := math.Max(20, math.Min(100, 100/(float64(cache)/1024)))

	log.Debug("Sanitizing Go's GC trigger", "percent", int(gogc))
	godebug.SetGCPercent(int(gogc))

	if ctx.GlobalIsSet(CbftWalDisabledFlag.Name) {
		cfg.CbftConfig.WalMode = false
	}
	if ctx.GlobalIsSet(SyncModeFlag.Name) {
		cfg.SyncMode = *GlobalTextMarshaler(ctx, SyncModeFlag.Name).(*downloader.SyncMode)
	}
	if ctx.GlobalIsSet(NetworkIdFlag.Name) {
		cfg.NetworkId = ctx.GlobalUint64(NetworkIdFlag.Name)
	}
	if ctx.GlobalIsSet(CacheFlag.Name) || ctx.GlobalIsSet(CacheDatabaseFlag.Name) {
		cfg.DatabaseCache = ctx.GlobalInt(CacheFlag.Name) * ctx.GlobalInt(CacheDatabaseFlag.Name) / 100
	}
	if ctx.GlobalIsSet(CacheTrieJournalFlag.Name) {
		cfg.TrieCleanCacheJournal = ctx.GlobalString(CacheTrieJournalFlag.Name)
	}
	if ctx.GlobalIsSet(CacheTrieRejournalFlag.Name) {
		cfg.TrieCleanCacheRejournal = ctx.GlobalDuration(CacheTrieRejournalFlag.Name)
	}
	cfg.DatabaseHandles = MakeDatabaseHandles(ctx.GlobalInt(FDLimitFlag.Name))
	if ctx.GlobalIsSet(AncientFlag.Name) {
		cfg.DatabaseFreezer = ctx.GlobalString(AncientFlag.Name)
	}

	if ctx.GlobalIsSet(TxLookupLimitFlag.Name) {
		cfg.TxLookupLimit = ctx.GlobalUint64(TxLookupLimitFlag.Name)
	}

	cfg.NoPruning = true

	if ctx.GlobalIsSet(CacheFlag.Name) || ctx.GlobalIsSet(CacheTrieFlag.Name) {
		cfg.TrieCleanCache = ctx.GlobalInt(CacheFlag.Name) * ctx.GlobalInt(CacheTrieFlag.Name) / 100
	}
	if ctx.GlobalIsSet(CacheFlag.Name) || ctx.GlobalIsSet(CacheGCFlag.Name) {
		cfg.TrieDirtyCache = ctx.GlobalInt(CacheFlag.Name) * ctx.GlobalInt(CacheGCFlag.Name) / 100
	}
	if ctx.GlobalIsSet(CacheFlag.Name) || ctx.GlobalIsSet(CacheSnapshotFlag.Name) {
		cfg.SnapshotCache = ctx.GlobalInt(CacheFlag.Name) * ctx.GlobalInt(CacheSnapshotFlag.Name) / 100
	}
	if !ctx.GlobalBool(SnapshotFlag.Name) {
		// If snap-sync is requested, this flag is also required
		if cfg.SyncMode == downloader.SnapSync {
			log.Info("Snap sync requested, enabling --snapshot")
		} else {
			cfg.TrieCleanCache += cfg.SnapshotCache
			cfg.SnapshotCache = 0 // Disabled
		}
	}
	if ctx.GlobalIsSet(CacheTrieDBFlag.Name) {
		cfg.TrieDBCache = ctx.GlobalInt(CacheTrieDBFlag.Name)
	}
	if ctx.GlobalIsSet(DocRootFlag.Name) {
		cfg.DocRoot = ctx.GlobalString(DocRootFlag.Name)
	}

	if ctx.GlobalIsSet(RPCGlobalGasCapFlag.Name) {
		cfg.RPCGasCap = ctx.GlobalUint64(RPCGlobalGasCapFlag.Name)
	}
	if cfg.RPCGasCap != 0 {
		log.Info("Set global gas cap", "cap", cfg.RPCGasCap)
	} else {
		log.Info("Global gas cap disabled")
	}
	if ctx.GlobalIsSet(RPCGlobalEVMTimeoutFlag.Name) {
		cfg.RPCEVMTimeout = ctx.GlobalDuration(RPCGlobalEVMTimeoutFlag.Name)
	}
	if ctx.GlobalIsSet(RPCGlobalTxFeeCapFlag.Name) {
		cfg.RPCTxFeeCap = ctx.GlobalFloat64(RPCGlobalTxFeeCapFlag.Name)
	}

	// Override any default configs for hard coded networks.
	switch {

	// Test NetWork
	case ctx.GlobalBool(TestnetFlag.Name):
		if !ctx.GlobalIsSet(NetworkIdFlag.Name) {
			cfg.NetworkId = 2000
		}
		cfg.Genesis = core.DefaultTestnetGenesisBlock()
	default:
		SetDNSDiscoveryDefaults(cfg, params.MainnetGenesisHash)
	}

	if ctx.GlobalIsSet(DBNoGCFlag.Name) {
		cfg.DBDisabledGC = ctx.GlobalBool(DBNoGCFlag.Name)
	}
	if ctx.GlobalIsSet(DBGCIntervalFlag.Name) {
		cfg.DBGCInterval = ctx.GlobalUint64(DBGCIntervalFlag.Name)
	}
	if ctx.GlobalIsSet(DBGCTimeoutFlag.Name) {
		cfg.DBGCTimeout = ctx.GlobalDuration(DBGCTimeoutFlag.Name)
	}
	if ctx.GlobalIsSet(DBGCMptFlag.Name) {
		cfg.DBGCMpt = ctx.GlobalBool(DBGCMptFlag.Name)
	}
	if ctx.GlobalIsSet(DBGCBlockFlag.Name) {
		b := ctx.GlobalInt(DBGCBlockFlag.Name)
		if b > 0 {
			cfg.DBGCBlock = b
		}
	}
	if ctx.GlobalIsSet(DBValidatorsHistoryFlag.Name) {
		cfg.DBValidatorsHistory = ctx.GlobalBool(DBValidatorsHistoryFlag.Name)
	}

	// Read the value from the flag no matter if it's set or not.
	cfg.Preimages = ctx.GlobalBool(CachePreimagesFlag.Name)
	if cfg.DBDisabledGC && !cfg.Preimages {
		cfg.Preimages = true
		log.Info("Enabling recording of key preimages since archive mode is used")
	}

	// vm options
	if ctx.GlobalIsSet(VMWasmType.Name) {
		cfg.VMWasmType = ctx.GlobalString(VMWasmType.Name)
	}
	if ctx.GlobalIsSet(VmTimeoutDuration.Name) {
		cfg.VmTimeoutDuration = ctx.GlobalUint64(VmTimeoutDuration.Name)
	}

}

// SetDNSDiscoveryDefaults configures DNS discovery with the given URL if
// no URLs are set.
func SetDNSDiscoveryDefaults(cfg *ethconfig.Config, genesis common.Hash) {
	if cfg.EthDiscoveryURLs != nil {
		return // already set through flags/config
	}
	protocol := "all"
	if cfg.SyncMode == downloader.LightSync {
		protocol = "les"
	}
	if url := params.KnownDNSNetwork(genesis, protocol); url != "" {
		cfg.EthDiscoveryURLs = []string{url}
		cfg.SnapDiscoveryURLs = cfg.EthDiscoveryURLs
	}
}

func SetCbft(ctx *cli.Context, cfg *types.OptionsConfig, nodeCfg *node.Config) {
	if nodeCfg.P2P.PrivateKey != nil {
		cfg.NodePriKey = nodeCfg.P2P.PrivateKey
		cfg.Node = enode.NewV4(&cfg.NodePriKey.PublicKey, nil, 0, 0)
		cfg.NodeID = cfg.Node.IDv0()
	}

	if ctx.GlobalIsSet(CbftBlsPriKeyFileFlag.Name) {
		priKey, err := bls.LoadBLS(ctx.GlobalString(CbftBlsPriKeyFileFlag.Name))
		if err != nil {
			Fatalf("Failed to load bls key from file: %v", err)
		}
		cfg.BlsPriKey = priKey
	} else {
		cfg.BlsPriKey = nodeCfg.BlsKey()
	}
	nodeCfg.P2P.BlsPublicKey = *(cfg.BlsPriKey.GetPublicKey())

	if ctx.GlobalIsSet(CbftWalDisabledFlag.Name) {
		cfg.WalMode = !ctx.GlobalBool(CbftWalDisabledFlag.Name)
	}

	if ctx.GlobalIsSet(CbftPeerMsgQueueSize.Name) {
		cfg.PeerMsgQueueSize = ctx.GlobalUint64(CbftPeerMsgQueueSize.Name)
	}

	if ctx.GlobalIsSet(CbftMaxPingLatency.Name) {
		cfg.MaxPingLatency = ctx.GlobalInt64(CbftMaxPingLatency.Name)
	}
	if ctx.GlobalIsSet(CbftBlacklistDeadlineFlag.Name) {
		cfg.BlacklistDeadline = ctx.GlobalInt64(CbftBlacklistDeadlineFlag.Name)
	}

}

// RegisterEthService adds an Ethereum client to the stack.
func RegisterEthService(stack *node.Node, cfg *ethconfig.Config) (ethapi.Backend, *eth.Ethereum) {
	if cfg.SyncMode == downloader.LightSync {
		Fatalf("Failed to register the Platon service: not les")
		return nil, nil
	} else {
		backend, err := eth.New(stack, cfg)
		if err != nil {
			Fatalf("Failed to register the PlatON service: %v", err)
		}
		stack.RegisterAPIs(tracers.APIs(backend.APIBackend))
		return backend.APIBackend, backend
	}
}

// RegisterShhService configures Whisper and adds it to the given node.
//func RegisterShhService(stack *node.Node, cfg *whisper.Config) {
//	if err := stack.Register(func(n *node.ServiceContext) (node.Service, error) {
//		return whisper.New(cfg), nil
//	}); err != nil {
//		Fatalf("Failed to register the Whisper service: %v", err)
//	}
//}

// RegisterEthStatsService configures the Ethereum Stats daemon and adds it to
// the given node.
func RegisterEthStatsService(stack *node.Node, backend ethapi.Backend, url string) {
	if err := ethstats.New(stack, backend, backend.Engine(), url); err != nil {
		Fatalf("Failed to register the Ethereum Stats service: %v", err)
	}
}

// RegisterGraphQLService is a utility function to construct a new service and register it against a node.
func RegisterGraphQLService(stack *node.Node, backend ethapi.Backend, cfg node.Config) {
	if err := graphql.New(stack, backend, cfg.GraphQLCors, cfg.GraphQLVirtualHosts); err != nil {
		Fatalf("Failed to register the GraphQL service: %v", err)
	}
}

func SetupMetrics(ctx *cli.Context) {
	if metrics.Enabled {
		log.Info("Enabling metrics collection")

		var (
			enableExport   = ctx.GlobalBool(MetricsEnableInfluxDBFlag.Name)
			enableExportV2 = ctx.GlobalBool(MetricsEnableInfluxDBV2Flag.Name)
		)

		if enableExport || enableExportV2 {
			CheckExclusive(ctx, MetricsEnableInfluxDBFlag, MetricsEnableInfluxDBV2Flag)

			v1FlagIsSet := ctx.GlobalIsSet(MetricsInfluxDBUsernameFlag.Name) ||
				ctx.GlobalIsSet(MetricsInfluxDBPasswordFlag.Name)

			v2FlagIsSet := ctx.GlobalIsSet(MetricsInfluxDBTokenFlag.Name) ||
				ctx.GlobalIsSet(MetricsInfluxDBOrganizationFlag.Name) ||
				ctx.GlobalIsSet(MetricsInfluxDBBucketFlag.Name)

			if enableExport && v2FlagIsSet {
				Fatalf("Flags --influxdb.metrics.organization, --influxdb.metrics.token, --influxdb.metrics.bucket are only available for influxdb-v2")
			} else if enableExportV2 && v1FlagIsSet {
				Fatalf("Flags --influxdb.metrics.username, --influxdb.metrics.password are only available for influxdb-v1")
			}
		}

		var (
			endpoint = ctx.GlobalString(MetricsInfluxDBEndpointFlag.Name)
			database = ctx.GlobalString(MetricsInfluxDBDatabaseFlag.Name)
			username = ctx.GlobalString(MetricsInfluxDBUsernameFlag.Name)
			password = ctx.GlobalString(MetricsInfluxDBPasswordFlag.Name)

			token        = ctx.GlobalString(MetricsInfluxDBTokenFlag.Name)
			bucket       = ctx.GlobalString(MetricsInfluxDBBucketFlag.Name)
			organization = ctx.GlobalString(MetricsInfluxDBOrganizationFlag.Name)
		)

		if enableExport {
			tagsMap := SplitTagsFlag(ctx.GlobalString(MetricsInfluxDBTagsFlag.Name))

			log.Info("Enabling metrics export to InfluxDB")

			go influxdb.InfluxDBWithTags(metrics.DefaultRegistry, 10*time.Second, endpoint, database, username, password, "platon.", tagsMap)
		} else if enableExportV2 {
			tagsMap := SplitTagsFlag(ctx.GlobalString(MetricsInfluxDBTagsFlag.Name))

			log.Info("Enabling metrics export to InfluxDB (v2)")

			go influxdb.InfluxDBV2WithTags(metrics.DefaultRegistry, 10*time.Second, endpoint, token, bucket, organization, "platon.", tagsMap)
		}

		if ctx.GlobalIsSet(MetricsHTTPFlag.Name) {
			address := fmt.Sprintf("%s:%d", ctx.GlobalString(MetricsHTTPFlag.Name), ctx.GlobalInt(MetricsPortFlag.Name))
			log.Info("Enabling stand-alone metrics HTTP endpoint", "address", address)
			exp.Setup(address)
		}
	}
}

func SplitTagsFlag(tagsFlag string) map[string]string {
	tags := strings.Split(tagsFlag, ",")
	tagsMap := map[string]string{}

	for _, t := range tags {
		if t != "" {
			kv := strings.Split(t, "=")

			if len(kv) == 2 {
				tagsMap[kv[0]] = kv[1]
			}
		}
	}

	return tagsMap
}

// MakeChainDatabase open an LevelDB using the flags passed to the client and will hard crash if it fails.
func MakeChainDatabase(ctx *cli.Context, stack *node.Node, readonly bool) ethdb.Database {
	var (
		cache   = ctx.GlobalInt(CacheFlag.Name) * ctx.GlobalInt(CacheDatabaseFlag.Name) / 100
		handles = MakeDatabaseHandles(ctx.GlobalInt(FDLimitFlag.Name))

		err     error
		chainDb ethdb.Database
	)
	switch {
	case ctx.GlobalIsSet(RemoteDBFlag.Name):
		log.Info("Using remote db", "url", ctx.GlobalString(RemoteDBFlag.Name))
		chainDb, err = remotedb.New(ctx.GlobalString(RemoteDBFlag.Name))
	case ctx.GlobalString(SyncModeFlag.Name) == "light":
		chainDb, err = stack.OpenDatabase("lightchaindata", cache, handles, "", readonly)
	default:
		chainDb, err = stack.OpenDatabaseWithFreezer("chaindata", cache, handles, ctx.GlobalString(AncientFlag.Name), "", readonly)
	}
	if err != nil {
		Fatalf("Could not open database: %v", err)
	}
	return chainDb
}

func MakeGenesis(ctx *cli.Context) *core.Genesis {
	var genesis *core.Genesis
	switch {
	case ctx.GlobalBool(TestnetFlag.Name):
		genesis = core.DefaultTestnetGenesisBlock()
	}
	return genesis
}

// MakeChain creates a chain manager from set command line flags.
func MakeChain(ctx *cli.Context, stack *node.Node) (chain *core.BlockChain, chainDb ethdb.Database) {
	var err error
	chainDb = MakeChainDatabase(ctx, stack, false)
	basedb, err := snapshotdb.Open(stack.ResolvePath(snapshotdb.DBPath), 0, 0, true)
	if err != nil {
		Fatalf("%v", err)
	}
	config, _, err := core.SetupGenesisBlock(chainDb, basedb, MakeGenesis(ctx))
	if err != nil {
		Fatalf("%v", err)
	}
	if err := basedb.Close(); err != nil {
		Fatalf("%v", err)
	}
	var engine consensus.Engine
	//todo: Merge confirmation.
	engine = consensus.NewFaker()

	cache := &core.CacheConfig{
		Disabled:        true,
		TrieCleanLimit:  ethconfig.Defaults.TrieCleanCache,
		TrieDirtyLimit:  ethconfig.Defaults.TrieCache,
		TrieTimeLimit:   ethconfig.Defaults.TrieTimeout,
		BodyCacheLimit:  ethconfig.Defaults.BodyCacheLimit,
		BlockCacheLimit: ethconfig.Defaults.BlockCacheLimit,
		MaxFutureBlocks: ethconfig.Defaults.MaxFutureBlocks,
		TriesInMemory:   ethconfig.Defaults.TriesInMemory,
		Preimages:       ctx.GlobalBool(CachePreimagesFlag.Name),
		SnapshotLimit:   ethconfig.Defaults.SnapshotCache,
	}
	if ctx.GlobalIsSet(CacheFlag.Name) || ctx.GlobalIsSet(CacheTrieFlag.Name) {
		cache.TrieCleanLimit = ctx.GlobalInt(CacheFlag.Name) * ctx.GlobalInt(CacheTrieFlag.Name) / 100
	}
	if !ctx.GlobalBool(SnapshotFlag.Name) {
		cache.SnapshotLimit = 0 // Disabled
	}
	if ethconfig.Defaults.DBDisabledGC && !cache.Preimages {
		cache.Preimages = true
		log.Info("Enabling recording of key preimages since archive mode is used")
	}
	if ctx.GlobalIsSet(CacheFlag.Name) || ctx.GlobalIsSet(CacheGCFlag.Name) {
		cache.TrieDirtyLimit = ctx.GlobalInt(CacheFlag.Name) * ctx.GlobalInt(CacheGCFlag.Name) / 100
	}
	vmcfg := vm.Config{}

	// TODO(rjl493456442) disable snapshot generation/wiping if the chain is read only.
	// Disable transaction indexing/unindexing by default.
	chain, err = core.NewBlockChain(chainDb, cache, config, engine, vmcfg, nil, nil)
	if err != nil {
		Fatalf("Can't create BlockChain: %v", err)
	}

	return chain, chainDb
}

// MakeConsolePreloads retrieves the absolute paths for the console JavaScript
// scripts to preload before starting.
func MakeConsolePreloads(ctx *cli.Context) []string {
	// Skip preloading if there's nothing to preload
	if ctx.GlobalString(PreloadJSFlag.Name) == "" {
		return nil
	}
	// Otherwise resolve absolute paths and return them
	preloads := []string{}

	for _, file := range strings.Split(ctx.GlobalString(PreloadJSFlag.Name), ",") {
		preloads = append(preloads, strings.TrimSpace(file))
	}
	return preloads
}

// MigrateFlags sets the global flag from a local flag when it's set.
// This is a temporary function used for migrating old command/flags to the
// new format.
//
// e.g. platon account new --keystore /tmp/mykeystore --lightkdf
//
// is equivalent after calling this method with:
//
// platon --keystore /tmp/mykeystore --lightkdf account new
//
// This allows the use of the existing configuration functionality.
// When all flags are migrated this function can be removed and the existing
// configuration functionality must be changed that is uses local flags
func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		for _, name := range ctx.FlagNames() {
			if ctx.IsSet(name) {
				ctx.GlobalSet(name, ctx.String(name))
			}
		}
		return action(ctx)
	}
}

// CheckExclusive verifies that only a single instance of the provided flags was
// set by the user. Each flag might optionally be followed by a string type to
// specialize it further.
func CheckExclusive(ctx *cli.Context, args ...interface{}) {
	set := make([]string, 0, 1)
	for i := 0; i < len(args); i++ {
		// Make sure the next argument is a flag and skip if not set
		flag, ok := args[i].(cli.Flag)
		if !ok {
			panic(fmt.Sprintf("invalid argument, not cli.Flag type: %T", args[i]))
		}
		// Check if next arg extends current and expand its name if so
		name := flag.GetName()

		if i+1 < len(args) {
			switch option := args[i+1].(type) {
			case string:
				// Extended flag check, make sure value set doesn't conflict with passed in option
				if ctx.GlobalString(flag.GetName()) == option {
					name += "=" + option
					set = append(set, "--"+name)
				}
				// shift arguments and continue
				i++
				continue

			case cli.Flag:
			default:
				panic(fmt.Sprintf("invalid argument, not cli.Flag or string extension: %T", args[i+1]))
			}
		}
		// Mark the flag if it's set
		if ctx.GlobalIsSet(flag.GetName()) {
			set = append(set, "--"+name)
		}
	}
	if len(set) > 1 {
		Fatalf("Flags %v can't be used at the same time", strings.Join(set, ", "))
	}
}
