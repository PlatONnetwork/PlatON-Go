// Copyright 2014 The go-ethereum Authors
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

// platon is the official command-line client for Ethereum.
package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/internal/flags"

	"github.com/PlatONnetwork/PlatON-Go/console/prompt"

	"github.com/panjf2000/ants/v2"
	"gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/accounts"
	"github.com/PlatONnetwork/PlatON-Go/accounts/keystore"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/ethclient"
	"github.com/PlatONnetwork/PlatON-Go/internal/debug"
	"github.com/PlatONnetwork/PlatON-Go/internal/ethapi"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
	"github.com/PlatONnetwork/PlatON-Go/node"

	// Force-load the tracer engines to trigger registration
	_ "github.com/PlatONnetwork/PlatON-Go/eth/tracers/js"
	_ "github.com/PlatONnetwork/PlatON-Go/eth/tracers/native"
)

const (
	clientIdentifier = "platon" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	gitDate   = ""
	// The app that holds all commands and flags.
	app = flags.NewApp(gitCommit, gitDate, "the platon-go command line interface")
	// flags that configure the node
	nodeFlags = []cli.Flag{
		utils.IdentityFlag,
		utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.BootnodesFlag,
		utils.BootnodesV4Flag,
		//	utils.BootnodesV5Flag,
		utils.DataDirFlag,
		utils.AncientFlag,
		utils.MinFreeDiskSpaceFlag,
		utils.KeyStoreDirFlag,
		utils.TxPoolLocalsFlag,
		utils.TxPoolNoLocalsFlag,
		utils.TxPoolJournalFlag,
		utils.TxPoolRejournalFlag,
		utils.TxPoolPriceBumpFlag,
		utils.TxPoolAccountSlotsFlag,
		utils.TxPoolGlobalSlotsFlag,
		utils.TxPoolAccountQueueFlag,
		utils.TxPoolGlobalQueueFlag,
		utils.TxPoolGlobalTxCountFlag,
		utils.TxPoolLifetimeFlag,
		utils.TxPoolCacheSizeFlag,
		utils.SyncModeFlag,
		utils.SnapshotFlag,
		utils.TxLookupLimitFlag,
		utils.LightKDFFlag,
		utils.BloomFilterSizeFlag,
		utils.CacheFlag,
		utils.CacheDatabaseFlag,
		utils.CacheTrieFlag,
		utils.CacheTrieJournalFlag,
		utils.CacheTrieRejournalFlag,
		utils.CacheGCFlag,
		utils.CacheSnapshotFlag,
		utils.CacheTrieDBFlag,
		utils.CachePreimagesFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxConsensusPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.MinerGasPriceFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		//	utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
		//	utils.DNSDiscoveryFlag,
		utils.MainFlag,
		utils.TestnetFlag,
		utils.NetworkIdFlag,
		//utils.EthStatsURLFlag,
		//utils.NoCompactionFlag,
		utils.GpoBlocksFlag,
		utils.LegacyGpoBlocksFlag,
		utils.GpoPercentileFlag,
		utils.LegacyGpoPercentileFlag,
		utils.GpoMaxGasPriceFlag,
		utils.GpoIgnoreGasPriceFlag,
		configFileFlag,
	}

	rpcFlags = []cli.Flag{
		utils.HTTPEnabledFlag,
		utils.HTTPListenAddrFlag,
		utils.HTTPPortFlag,
		utils.HTTPCORSDomainFlag,
		utils.HTTPVirtualHostsFlag,
		utils.GraphQLEnabledFlag,
		utils.GraphQLCORSDomainFlag,
		utils.GraphQLVirtualHostsFlag,
		utils.HTTPApiFlag,
		utils.HTTPPathPrefixFlag,
		utils.HTTPEnabledEthCompatibleFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.LegacyWSListenAddrFlag,
		utils.WSPortFlag,
		utils.LegacyWSPortFlag,
		utils.WSApiFlag,
		utils.LegacyWSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.WSPathPrefixFlag,
		utils.LegacyWSAllowedOriginsFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
		utils.InsecureUnlockAllowedFlag,
		utils.RPCGlobalGasCapFlag,
		utils.RPCGlobalEVMTimeoutFlag,
		utils.RPCGlobalTxFeeCapFlag,
		utils.AllowUnprotectedTxs,
	}

	metricsFlags = []cli.Flag{
		utils.MetricsEnabledFlag,
		utils.MetricsEnabledExpensiveFlag,
		utils.MetricsHTTPFlag,
		utils.MetricsPortFlag,
		utils.MetricsEnableInfluxDBFlag,
		utils.MetricsInfluxDBEndpointFlag,
		utils.MetricsInfluxDBDatabaseFlag,
		utils.MetricsInfluxDBUsernameFlag,
		utils.MetricsInfluxDBPasswordFlag,
		utils.MetricsInfluxDBTagsFlag,
		utils.MetricsEnableInfluxDBV2Flag,
		utils.MetricsInfluxDBTokenFlag,
		utils.MetricsInfluxDBBucketFlag,
		utils.MetricsInfluxDBOrganizationFlag,
	}

	cbftFlags = []cli.Flag{
		utils.CbftPeerMsgQueueSize,
		utils.CbftWalDisabledFlag,
		utils.CbftMaxPingLatency,
		utils.CbftBlsPriKeyFileFlag,
		utils.CbftBlacklistDeadlineFlag,
	}

	dbFlags = []cli.Flag{
		utils.DBNoGCFlag,
		utils.DBGCIntervalFlag,
		utils.DBGCTimeoutFlag,
		utils.DBGCMptFlag,
		utils.DBGCBlockFlag,
		utils.DBValidatorsHistoryFlag,
		utils.DBNoCacheFlag,
		utils.DBCacheEpochFlag,
	}

	vmFlags = []cli.Flag{
		utils.VMWasmType,
		utils.VmTimeoutDuration,
	}
)

func init() {
	// Initialize the CLI app and start PlatON
	app.Action = platon
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2023 The PlatON-Go Authors"
	app.Commands = []cli.Command{
		// See chaincmd.go:
		initCommand,
		//importCommand,
		//exportCommand,
		importPreimagesCommand,
		exportPreimagesCommand,
		removedbCommand,
		dumpCommand,
		dumpGenesisCommand,
		// See accountcmd.go:
		accountCommand,
		// See consolecmd.go:
		consoleCommand,
		attachCommand,
		javascriptCommand,
		versionCommand,
		licenseCommand,
		// See config.go
		dumpConfigCommand,
		// see dbcmd.go
		dbCommand,
		// See cmd/utils/flags_legacy.go
		utils.ShowDeprecated,
		// See snapshot.go
		snapshotCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, consoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	//app.Flags = append(app.Flags, whisperFlags...)
	app.Flags = append(app.Flags, metricsFlags...)

	// for cbft
	app.Flags = append(app.Flags, cbftFlags...)
	app.Flags = append(app.Flags, dbFlags...)
	app.Flags = append(app.Flags, vmFlags...)

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		err := bls.Init(int(bls.BLS12_381))
		if err != nil {
			return err
		}

		if err := debug.Setup(ctx); err != nil {
			return err
		}

		//init wasm logfile
		if err := debug.SetupWasmLog(ctx); err != nil {
			return err
		}

		// Start metrics export if enabled
		utils.SetupMetrics(ctx)

		// Start system runtime metrics collection
		go metrics.CollectProcessMetrics(3 * time.Second)

		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		prompt.Stdin.Close() // Resets terminal mode.
		ants.Release()
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// platon is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func platon(ctx *cli.Context) error {
	if args := ctx.Args(); len(args) > 0 {
		return fmt.Errorf("invalid command: %q", args[0])
	}
	stack, backend := makeFullNode(ctx)
	defer stack.Close()
	startNode(ctx, stack, backend, false)
	stack.Wait()
	return nil
}

// startNode boots up the system node and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces and the
// miner.
func startNode(ctx *cli.Context, stack *node.Node, backend ethapi.Backend, isConsole bool) {
	debug.Memsize.Add("node", stack)

	// Start up the node itself
	utils.StartNode(ctx, stack, isConsole)

	// Unlock any account specifically requested
	unlockAccounts(ctx, stack)

	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	go func() {
		// Create a chain state reader for self-derivation
		rpcClient, err := stack.Attach()
		if err != nil {
			utils.Fatalf("Failed to attach to self: %v", err)
		}
		stateReader := ethclient.NewClient(rpcClient)

		// Open any wallets already attached
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				var derivationPaths []accounts.DerivationPath
				if event.Wallet.URL().Scheme == "ledger" {
					derivationPaths = append(derivationPaths, accounts.LegacyLedgerBaseDerivationPath)
				}
				derivationPaths = append(derivationPaths, accounts.DefaultBaseDerivationPath)

				event.Wallet.SelfDerive(derivationPaths, stateReader)

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()
	// Start auxiliary services if enabled
	// Mining only makes sense if a full Ethereum node is running
	if ctx.GlobalString(utils.SyncModeFlag.Name) == "light" {
		utils.Fatalf("Light clients do not support mining")
	}
	ethBackend, ok := backend.(*eth.EthAPIBackend)
	if !ok {
		utils.Fatalf("Ethereum service not running")
	}
	// Set the gas price to the limits from the CLI and start mining
	gasprice := utils.GlobalBig(ctx, utils.MinerGasPriceFlag.Name)

	ethBackend.TxPool().SetGasPrice(gasprice)

	if err := ethBackend.StartMining(); err != nil {
		utils.Fatalf("Failed to start mining: %v", err)
	}
}

// unlockAccounts unlocks any account specifically requested.
func unlockAccounts(ctx *cli.Context, stack *node.Node) {
	var unlocks []string
	inputs := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for _, input := range inputs {
		if trimmed := strings.TrimSpace(input); trimmed != "" {
			unlocks = append(unlocks, trimmed)
		}
	}
	// Short circuit if there is no account to unlock.
	if len(unlocks) == 0 {
		return
	}
	// If insecure account unlocking is not allowed if node's APIs are exposed to external.
	// Print warning log to user and skip unlocking.
	if !stack.Config().InsecureUnlockAllowed && stack.Config().ExtRPCEnabled() {
		utils.Fatalf("Account unlock with HTTP access is forbidden!")
	}
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	passwords := utils.MakePasswordList(ctx)
	for i, account := range unlocks {
		unlockAccount(ks, account, i, passwords)
	}
}
