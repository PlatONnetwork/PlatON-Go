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
	"io"
	"os"
	"reflect"
	"unicode"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/cmd/utils"
	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/naoina/toml"
)

var (
	dumpConfigCommand = cli.Command{
		Action:    utils.MigrateFlags(dumpConfig),
		Name:      "dumpconfig",
		Usage:     "Show configuration values",
		ArgsUsage: "",
		//Flags:       append(append(nodeFlags, rpcFlags...), whisperFlags...),
		Flags:       append(append(nodeFlags, rpcFlags...)),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}

	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
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
		link := ""
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
	Eth      eth.Config
	Node     node.Config
	Ethstats ethstatsConfig
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
	cfg.Version = params.VersionWithCommit(gitCommit)
	cfg.HTTPModules = append(cfg.HTTPModules, "platon")
	cfg.WSModules = append(cfg.WSModules, "platon")
	cfg.IPCPath = "platon.ipc"
	return cfg
}

func makeConfigNode(ctx *cli.Context) (*node.Node, platonConfig) {

	// Load defaults.
	cfg := platonConfig{
		Eth:  eth.DefaultConfig,
		Node: defaultNodeConfig(),
	}

	//

	// Load config file.
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		/*	if err := loadConfig(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}*/
		if err := loadConfigFile(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// Current version only supports full syncmode
	// ctx.GlobalSet(utils.SyncModeFlag.Name, cfg.Eth.SyncMode.String())

	// Apply flags.
	utils.SetNodeConfig(ctx, &cfg.Node)
	utils.SetCbft(ctx, &cfg.Eth.CbftConfig, &cfg.Node)
	stack, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
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

	//if ctx.GlobalIsSet(utils.EthStatsURLFlag.Name) {
	//	cfg.Ethstats.URL = ctx.GlobalString(utils.EthStatsURLFlag.Name)
	//}

	//utils.SetShhConfig(ctx, stack, &cfg.Shh)

	return stack, cfg
}

func makeFullNode(ctx *cli.Context) *node.Node {

	stack, cfg := makeConfigNode(ctx)

	snapshotdb.SetDBPathWithNode(stack.ResolvePath(snapshotdb.DBPath))

	utils.RegisterEthService(stack, &cfg.Eth)

	// Add the Ethereum Stats daemon if requested.
	if cfg.Ethstats.URL != "" {
		utils.RegisterEthStatsService(stack, cfg.Ethstats.URL)
	}
	return stack
}

func makeFullNodeForCBFT(ctx *cli.Context) (*node.Node, platonConfig) {
	stack, cfg := makeConfigNode(ctx)
	snapshotdb.SetDBPathWithNode(stack.ResolvePath(snapshotdb.DBPath))

	utils.RegisterEthService(stack, &cfg.Eth)

	// Add the Ethereum Stats daemon if requested.
	if cfg.Ethstats.URL != "" {
		utils.RegisterEthStatsService(stack, cfg.Ethstats.URL)
	}
	return stack, cfg
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
	io.WriteString(os.Stdout, comment)
	os.Stdout.Write(out)
	return nil
}
