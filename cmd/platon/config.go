// Copyright 2017 The go-ethereum Authors
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

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/accounts/external"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/internal/flags"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"os"
	"reflect"
	"unicode"

	"github.com/PlatONnetwork/PlatON-Go/eth/ethconfig"
	"github.com/PlatONnetwork/PlatON-Go/internal/ethapi"
	"github.com/PlatONnetwork/PlatON-Go/metrics"

	"github.com/PlatONnetwork/PlatON-Go/accounts/keystore"
	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/naoina/toml"
	cli "github.com/urfave/cli/v2"
)

var (
	dumpConfigCommand = &cli.Command{
		Action:      dumpConfig,
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       flags.Merge(nodeFlags, rpcFlags),
		Description: `The dumpconfig command shows configuration values.`,
	}

	configFileFlag = &cli.StringFlag{
		Name:     "config",
		Usage:    "TOML configuration file",
		Category: flags.EthCategory,
	}
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		id := fmt.Sprintf("%s.%s", rt.String(), field)
		if deprecated(id) {
			log.Warn("Config field is deprecated and won't have an effect", "name", id)
			return nil
		}
		var link string
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}

type platonConfig struct {
	Eth      ethconfig.Config
	Node     node.Config
	Ethstats ethstatsConfig
	Metrics  metrics.Config
}

func loadConfig(file string, cfg *platonConfig) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

func loadConfigFile(filePath string, cfg *platonConfig) error {
	file, err := os.Open(filePath)
	if err != nil {
		utils.Fatalf("Failed to read config file: %v", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(cfg)
	if err != nil {
		utils.Fatalf("invalid config file: %v", err)
	}
	return err
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit, gitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "platon")
	cfg.WSModules = append(cfg.WSModules, "platon")
	cfg.IPCPath = "platon.ipc"
	return cfg
}

// makeConfigNode loads geth configuration and creates a blank node instance.
func makeConfigNode(ctx *cli.Context) (*node.Node, platonConfig) {
	// Load defaults.
	cfg := platonConfig{
		Eth:     ethconfig.Defaults,
		Node:    defaultNodeConfig(),
		Metrics: metrics.DefaultConfig,
	}

	// Load config file.
	if file := ctx.String(configFileFlag.Name); file != "" {
		/*	if err := loadConfig(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}*/
		if err := loadConfigFile(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}
	// Apply flags.
	utils.SetNodeConfig(ctx, &cfg.Node)
	utils.SetCbft(ctx, &cfg.Eth.CbftConfig, &cfg.Node)
	stack, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}
	// Node doesn't by default populate account manager backends
	if err := setAccountManagerBackends(stack); err != nil {
		utils.Fatalf("Failed to set account manager backends: %v", err)
	}

	utils.SetEthConfig(ctx, stack, &cfg.Eth)

	// pass on the rpc port to mpc pool conf.
	//cfg.Eth.MPCPool.LocalRpcPort = cfg.Node.HTTPPort

	// pass on the rpc port to vc pool conf.
	//cfg.Eth.VCPool.LocalRpcPort = cfg.Node.HTTPPort

	// load cbft config file.
	//if cbftConfig := cfg.Eth.LoadCbftConfig(cfg.Node); cbftConfig != nil {
	//	cfg.Eth.CbftConfig = *cbftConfig
	//}

	//if ctx.IsSet(utils.EthStatsURLFlag.Name) {
	//	cfg.Ethstats.URL = ctx.String(utils.EthStatsURLFlag.Name)
	//}

	//utils.SetShhConfig(ctx, stack, &cfg.Shh)

	applyMetricConfig(ctx, &cfg)

	return stack, cfg
}

// makeFullNode loads platon configuration and creates the backend.
func makeFullNode(ctx *cli.Context) (*node.Node, ethapi.Backend) {

	stack, cfg := makeConfigNode(ctx)
	snapshotdb.SetDBPathWithNode(stack.ResolvePath(snapshotdb.DBPath))

	backend, eth := utils.RegisterEthService(stack, &cfg.Eth)
	// Warn users to migrate if they have a legacy freezer format.
	if eth != nil && !ctx.IsSet(utils.IgnoreLegacyReceiptsFlag.Name) {
		firstIdx := uint64(0)
		// Hack to speed up check for mainnet because we know
		// the first non-empty block.
		ghash := rawdb.ReadCanonicalHash(eth.ChainDb(), 0)
		if cfg.Eth.NetworkId == 1 && ghash == params.MainnetGenesisHash {
			firstIdx = 1
		}
		isLegacy, _, err := dbHasLegacyReceipts(eth.ChainDb(), firstIdx)
		if err != nil {
			log.Error("Failed to check db for legacy receipts", "err", err)
		} else if isLegacy {
			stack.Close()
			utils.Fatalf("Database has receipts with a legacy format. Please run `geth db freezer-migrate`.")
		}
	}

	// Configure log filter RPC API.
	filterSystem := utils.RegisterFilterAPI(stack, backend, &cfg.Eth)

	// Configure GraphQL if requested.
	if ctx.IsSet(utils.GraphQLEnabledFlag.Name) {
		utils.RegisterGraphQLService(stack, backend, filterSystem, &cfg.Node)
	}

	// Add the Ethereum Stats daemon if requested.
	if cfg.Ethstats.URL != "" {
		utils.RegisterEthStatsService(stack, backend, cfg.Ethstats.URL)
	}
	return stack, backend
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	_, cfg := makeConfigNode(ctx)
	comment := ""

	if cfg.Eth.Genesis != nil {
		cfg.Eth.Genesis = nil
		comment += "# Note: this config doesn't contain the genesis block.\n\n"
	}

	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}

	dump := os.Stdout
	if ctx.NArg() > 0 {
		dump, err = os.OpenFile(ctx.Args().Get(0), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer dump.Close()
	}
	dump.WriteString(comment)
	dump.Write(out)

	return nil
}

func applyMetricConfig(ctx *cli.Context, cfg *platonConfig) {
	if ctx.IsSet(utils.MetricsEnabledFlag.Name) {
		cfg.Metrics.Enabled = ctx.Bool(utils.MetricsEnabledFlag.Name)
	}
	if ctx.IsSet(utils.MetricsEnabledExpensiveFlag.Name) {
		cfg.Metrics.EnabledExpensive = ctx.Bool(utils.MetricsEnabledExpensiveFlag.Name)
	}
	if ctx.IsSet(utils.MetricsHTTPFlag.Name) {
		cfg.Metrics.HTTP = ctx.String(utils.MetricsHTTPFlag.Name)
	}
	if ctx.IsSet(utils.MetricsPortFlag.Name) {
		cfg.Metrics.Port = ctx.Int(utils.MetricsPortFlag.Name)
	}
	if ctx.IsSet(utils.MetricsEnableInfluxDBFlag.Name) {
		cfg.Metrics.EnableInfluxDB = ctx.Bool(utils.MetricsEnableInfluxDBFlag.Name)
	}
	if ctx.IsSet(utils.MetricsInfluxDBEndpointFlag.Name) {
		cfg.Metrics.InfluxDBEndpoint = ctx.String(utils.MetricsInfluxDBEndpointFlag.Name)
	}
	if ctx.IsSet(utils.MetricsInfluxDBDatabaseFlag.Name) {
		cfg.Metrics.InfluxDBDatabase = ctx.String(utils.MetricsInfluxDBDatabaseFlag.Name)
	}
	if ctx.IsSet(utils.MetricsInfluxDBUsernameFlag.Name) {
		cfg.Metrics.InfluxDBUsername = ctx.String(utils.MetricsInfluxDBUsernameFlag.Name)
	}
	if ctx.IsSet(utils.MetricsInfluxDBPasswordFlag.Name) {
		cfg.Metrics.InfluxDBPassword = ctx.String(utils.MetricsInfluxDBPasswordFlag.Name)
	}
	if ctx.IsSet(utils.MetricsInfluxDBTagsFlag.Name) {
		cfg.Metrics.InfluxDBTags = ctx.String(utils.MetricsInfluxDBTagsFlag.Name)
	}
	if ctx.IsSet(utils.MetricsEnableInfluxDBV2Flag.Name) {
		cfg.Metrics.EnableInfluxDBV2 = ctx.Bool(utils.MetricsEnableInfluxDBV2Flag.Name)
	}
	if ctx.IsSet(utils.MetricsInfluxDBTokenFlag.Name) {
		cfg.Metrics.InfluxDBToken = ctx.String(utils.MetricsInfluxDBTokenFlag.Name)
	}
	if ctx.IsSet(utils.MetricsInfluxDBBucketFlag.Name) {
		cfg.Metrics.InfluxDBBucket = ctx.String(utils.MetricsInfluxDBBucketFlag.Name)
	}
	if ctx.IsSet(utils.MetricsInfluxDBOrganizationFlag.Name) {
		cfg.Metrics.InfluxDBOrganization = ctx.String(utils.MetricsInfluxDBOrganizationFlag.Name)
	}
}

func deprecated(field string) bool {
	return false
}

func setAccountManagerBackends(stack *node.Node) error {
	conf := stack.Config()
	am := stack.AccountManager()
	keydir := stack.KeyStoreDir()
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	if conf.UseLightweightKDF {
		scryptN = keystore.LightScryptN
		scryptP = keystore.LightScryptP
	}

	// Assemble the supported backends
	if len(conf.ExternalSigner) > 0 {
		log.Info("Using external signer", "url", conf.ExternalSigner)
		if extapi, err := external.NewExternalBackend(conf.ExternalSigner); err == nil {
			am.AddBackend(extapi)
			return nil
		} else {
			return fmt.Errorf("error connecting to external signer: %v", err)
		}
	}

	// For now, we're using EITHER external signer OR local signers.
	// If/when we implement some form of lockfile for USB and keystore wallets,
	// we can have both, but it's very confusing for the user to see the same
	// accounts in both externally and locally, plus very racey.
	am.AddBackend(keystore.NewKeyStore(keydir, scryptN, scryptP))

	return nil
}
