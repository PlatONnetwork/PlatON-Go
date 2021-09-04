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
	"math"
	"os"
	"runtime"
	godebug "runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/panjf2000/ants/v2"
	"gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/accounts"
	"github.com/PlatONnetwork/PlatON-Go/accounts/keystore"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/console"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/ethclient"
	"github.com/PlatONnetwork/PlatON-Go/internal/debug"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
	"github.com/PlatONnetwork/PlatON-Go/node"

	"github.com/elastic/gosigar"
)

const (
	clientIdentifier = "platon" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	gitDate   = ""
	// The app that holds all commands and flags.
	app = utils.NewApp(gitCommit, gitDate, "the platon-go command line interface")
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
		utils.KeyStoreDirFlag,
		utils.NoUSBFlag,
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
		utils.LightServFlag,
		utils.LightPeersFlag,
		utils.LightKDFFlag,
		utils.CacheFlag,
		utils.CacheDatabaseFlag,
		utils.CacheGCFlag,
		utils.CacheTrieDBFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxConsensusPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.MinerGasTargetFlag,
		//utils.MinerGasLimitFlag,
		utils.MinerGasPriceFlag,
		//	utils.MinerExtraDataFlag,
		//utils.MinerLegacyExtraDataFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		//	utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
		utils.DeveloperPeriodFlag,
		utils.MainFlag,
		utils.TestnetFlag,
		utils.NetworkIdFlag,
		//utils.EthStatsURLFlag,
		utils.NoCompactionFlag,
		utils.GpoBlocksFlag,
		utils.GpoPercentileFlag,
		configFileFlag,
	}

	rpcFlags = []cli.Flag{
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RPCApiFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
		utils.RPCGlobalGasCap,
	}

	//whisperFlags = []cli.Flag{
	//	utils.WhisperEnabledFlag,
	//	utils.WhisperMaxMessageSizeFlag,
	//	utils.WhisperRestrictConnectionBetweenLightClientsFlag,
	//}

	metricsFlags = []cli.Flag{
		utils.MetricsEnabledFlag,
		utils.MetricsEnabledExpensiveFlag,
		utils.MetricsEnableInfluxDBFlag,
		utils.MetricsInfluxDBEndpointFlag,
		utils.MetricsInfluxDBDatabaseFlag,
		utils.MetricsInfluxDBUsernameFlag,
		utils.MetricsInfluxDBPasswordFlag,
		utils.MetricsInfluxDBTagsFlag,
	}

	//mpcFlags = []cli.Flag{
	//	//utils.MPCEnabledFlag,
	//	utils.MPCIceFileFlag,
	//	utils.MPCActorFlag,
	//}
	//vcFlags = []cli.Flag{
	//	utils.VCEnabledFlag,
	//	utils.VCActorFlag,
	//	utils.VCPasswordFlag,
	//}

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
	app.Copyright = "Copyright 2019 The PlatON-Go Authors"
	app.Commands = []cli.Command{
		// See chaincmd.go:
		initCommand,
		//importCommand,
		//exportCommand,
		importPreimagesCommand,
		exportPreimagesCommand,
		copydbCommand,
		removedbCommand,
		dumpCommand,
		inspectCommand,
		// See accountcmd.go:
		accountCommand,
		// See consolecmd.go:
		consoleCommand,
		attachCommand,
		javascriptCommand,
		versionCommand,
		bugCommand,
		licenseCommand,
		// See config.go
		dumpConfigCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, consoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	//app.Flags = append(app.Flags, whisperFlags...)
	app.Flags = append(app.Flags, metricsFlags...)

	// for mpc
	//app.Flags = append(app.Flags, mpcFlags...)
	//// for vc
	//app.Flags = append(app.Flags, vcFlags...)
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

		logdir := ""
		if err := debug.Setup(ctx, logdir); err != nil {
			return err
		}

		//init wasm logfile
		if err := debug.SetupWasmLog(ctx); err != nil {
			return err
		}

		// Cap the cache allowance and tune the garbage collector
		var mem gosigar.Mem
		// Workaround until OpenBSD support lands into gosigar
		// Check https://github.com/elastic/gosigar#supported-platforms
		if runtime.GOOS != "openbsd" {
			if err := mem.Get(); err == nil {
				allowance := int(mem.Total / 1024 / 1024 / 3)
				if cache := ctx.GlobalInt(utils.CacheFlag.Name); cache > allowance {
					log.Warn("Sanitizing cache to Go's GC limits", "provided", cache, "updated", allowance)
					ctx.GlobalSet(utils.CacheFlag.Name, strconv.Itoa(allowance))
				}
			}
		}
		// Ensure Go's GC ignores the database cache for trigger percentage
		cache := ctx.GlobalInt(utils.CacheFlag.Name)
		gogc := math.Max(20, math.Min(100, 100/(float64(cache)/1024)))

		log.Debug("Sanitizing Go's GC trigger", "percent", int(gogc))
		godebug.SetGCPercent(int(gogc))

		// Start metrics export if enabled
		utils.SetupMetrics(ctx)

		// Start system runtime metrics collection
		go metrics.CollectProcessMetrics(3 * time.Second)

		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
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
	node := makeFullNode(ctx)
	startNode(ctx, node)
	node.Wait()
	return nil
}

// startNode boots up the system node and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces and the
// miner.
func startNode(ctx *cli.Context, stack *node.Node) {
	debug.Memsize.Add("node", stack)

	// Start up the node itself
	utils.StartNode(stack)

	// Unlock any account specifically requested
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	passwords := utils.MakePasswordList(ctx)
	unlocks := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for i, account := range unlocks {
		if trimmed := strings.TrimSpace(account); trimmed != "" {
			unlockAccount(ctx, ks, trimmed, i, passwords)
		}
	}
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

				derivationPath := accounts.DefaultBaseDerivationPath
				if event.Wallet.URL().Scheme == "ledger" {
					derivationPath = accounts.DefaultLedgerBaseDerivationPath
				}
				event.Wallet.SelfDerive(derivationPath, stateReader)

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
	var ethereum *eth.Ethereum
	if err := stack.Service(&ethereum); err != nil {
		utils.Fatalf("PlatON service not running: %v", err)
	}
	// Set the gas price to the limits from the CLI and start mining
	gasprice := utils.GlobalBig(ctx, utils.MinerGasPriceFlag.Name)

	ethereum.TxPool().SetGasPrice(gasprice)

	if err := ethereum.StartMining(); err != nil {
		utils.Fatalf("Failed to start mining: %v", err)
	}
}
