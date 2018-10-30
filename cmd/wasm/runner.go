package main

import (
	"Platon-go/cmd/utils"
	"Platon-go/common"
	"Platon-go/core"
	"Platon-go/core/state"
	"Platon-go/core/vm"
	"Platon-go/ethdb"
	"Platon-go/life/runtime"
	"Platon-go/log"
	"Platon-go/params"
	"Platon-go/rlp"
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	covert "Platon-go/life/utils"
	goruntime "runtime"
)

var runCommond = cli.Command{
	Action:      runCmd,
	Name:        "run",
	Usage:       "run arbitrary wasm binary",
	ArgsUsage:   "<code>",
	Description: "The run command runs arbitrary WASM code",
}

// readGenesis will read the given JSON format genesis file and
// return the initialized Genesis structure
func readGenesis(genesisPath string) *core.Genesis {
	if len(genesisPath) == 0 {
		utils.Fatalf("Must supply path to genesis JSON file")
	}
	file, err := os.Open(genesisPath)
	if err != nil {
		utils.Fatalf("Failed to read genesis file : %v", err)
	}
	defer file.Close()

	genesis := new(core.Genesis)
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		utils.Fatalf("invalid genesis file : %v", err)
	}
	return genesis
}

func runCmd(ctx *cli.Context) error {
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(ctx.GlobalInt(VerbosityFlag.Name)))
	log.Root().SetHandler(glogger)
	logconfig := &vm.LogConfig{
		DisableMemory: false,
		DisableStack:  false,
		Debug:         ctx.GlobalBool(DebugFlag.Name),
	}

	var (
		tracer			vm.Tracer
		debugLogger		*vm.StructLogger
		statedb			*state.StateDB
		chainConfig		*params.ChainConfig
		sender			= common.BytesToAddress([]byte("sender"))
		receiver		= common.BytesToAddress([]byte("receiver"))
		genesisConfig	*core.Genesis
	)
	if ctx.GlobalBool(MachineFlag.Name) {
		tracer = NewJSONLogger(logconfig, os.Stdout)
	} else if ctx.GlobalBool(DebugFlag.Name) {
		debugLogger = vm.NewStructLogger(logconfig)
		tracer = debugLogger
	} else {
		debugLogger = vm.NewStructLogger(logconfig)
	}
	if ctx.GlobalString(GenesisFlag.Name) != "" {
		gen := readGenesis(ctx.GlobalString(GenesisFlag.Name))
		genesisConfig = gen
		db := ethdb.NewMemDatabase()
		genesis := gen.ToBlock(db)
		statedb, _ = state.New(genesis.Root(), state.NewDatabase(db))
		chainConfig = gen.Config
	} else {
		statedb, _ = state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
		genesisConfig = new(core.Genesis)
	}
	if ctx.GlobalString(SenderFlag.Name) != "" {
		sender = common.HexToAddress(ctx.GlobalString(SenderFlag.Name))
	}
	statedb.CreateAccount(sender)
	if ctx.GlobalString(ReceiverFlag.Name) != "" {
		receiver = common.HexToAddress(ctx.GlobalString(ReceiverFlag.Name))
	}

	var (
		code 	[]byte
		abi		[]byte
		ret		[]byte
		err 	error
	)

	// The '--code' or '--codefile' flag overrides code in state
	if ctx.GlobalString(CodeFileFlag.Name) != "" {
		var hexcode []byte
		var err error
		if ctx.GlobalString(CodeFileFlag.Name) == "-" {
			// try reading from stdin
			if hexcode, err = ioutil.ReadAll(os.Stdin); err != nil {
				utils.Fatalf("Could not load code from stdin : %v", err)
			}
		} else {
			// codefile with hex assembly
			if hexcode, err = ioutil.ReadFile(ctx.GlobalString(CodeFileFlag.Name)); err != nil {
				utils.Fatalf("Could not load code from file: %v", err)
			}
		}
		// 剔除换行符
		code = common.Hex2Bytes(string(bytes.TrimRight(hexcode,"\n")))
	} else if ctx.GlobalString(CodeFlag.Name) != "" {
		code = common.Hex2Bytes(ctx.GlobalString(CodeFlag.Name))
	}

	// The '--abi' or '--abifile' flag overrides abi in state
	if ctx.GlobalString(AbiFileFlag.Name) != "" {
		var strabi []byte
		var err error
		if ctx.GlobalString(AbiFileFlag.Name) == "-" {
			// try reading from stdin
			if strabi, err = ioutil.ReadAll(os.Stdin); err != nil {
				utils.Fatalf("Could not load abi from stdin : %v", err)
			}
		} else {
			// codefile with hex assembly
			if strabi, err = ioutil.ReadFile(ctx.GlobalString(AbiFileFlag.Name)); err != nil {
				utils.Fatalf("Could not load abi from file: %v", err)
			}
		}
		hexabi := common.Bytes2Hex(bytes.TrimRight(strabi,"\n"))
		abi = common.Hex2Bytes(hexabi)
	} else if ctx.GlobalString(AbiFlag.Name) != "" {
		abi = []byte(ctx.GlobalString(AbiFlag.Name))
	}

	initialGas := ctx.GlobalUint64(GasFlag.Name)
	if genesisConfig.GasLimit != 0 {
		initialGas = genesisConfig.GasLimit
	}
	runtimeConfig := runtime.Config{
		Origin:      sender,
		State:       statedb,
		GasLimit:    initialGas,
		GasPrice:    utils.GlobalBig(ctx, GasPriceFlag.Name),
		Value:       utils.GlobalBig(ctx, ValueFlag.Name),
		Difficulty:  genesisConfig.Difficulty,
		Time:        new(big.Int).SetUint64(genesisConfig.Timestamp),
		Coinbase:    genesisConfig.Coinbase,
		BlockNumber: new(big.Int).SetUint64(genesisConfig.Number),
		EVMConfig: vm.Config{
			Tracer: tracer,
			Debug:  ctx.GlobalBool(DebugFlag.Name) || ctx.GlobalBool(MachineFlag.Name),
		},
	}

	if chainConfig != nil {
		runtimeConfig.ChainConfig = chainConfig
	}

	txType := ctx.GlobalInt64(TxTypeFlag.Name)

	tstart := time.Now()
	var leftOverGas uint64
	if ctx.GlobalBool(CreateFlag.Name) {
		// 合约创建逻辑，input为外部输入, 可能是参数，在以太坊中。在wasm中需要进行编码才能完成
		rlpData := make([][]byte,0)
		rlpData = append(rlpData, covert.Int64ToBytes(txType), abi, code)

		buffer := new(bytes.Buffer)
		err := rlp.Encode(buffer, rlpData)
		if err != nil {
			utils.Fatalf("rlp parse fail: %v", err)
		}
		ret, _, leftOverGas, err = runtime.Create(buffer.Bytes(), &runtimeConfig)
	} else {
		if len(code) > 0 {
			statedb.SetCode(receiver, code)
		}
		if len(abi) > 0 {
			statedb.SetAbi(receiver, abi)
		}
		// input : rlp.encoded format.
		input := common.Hex2Bytes(ctx.GlobalString(InputFlag.Name))
		ret, leftOverGas, err = runtime.Call(receiver, input, &runtimeConfig)
	}
	execTime := time.Since(tstart)

	statedb.IntermediateRoot(true)
	fmt.Println(string(statedb.Dump()))

	if ctx.GlobalBool(StatDumpFlag.Name) {
		var mem goruntime.MemStats
		goruntime.ReadMemStats(&mem)
		fmt.Fprintf(os.Stderr, `evm execution time: %v
heap objects:       %d
allocations:        %d
total allocations:  %d
GC calls:           %d
Gas used:           %d

`, execTime, mem.HeapObjects, mem.Alloc, mem.TotalAlloc, mem.NumGC, initialGas-leftOverGas)
	}

	if tracer == nil {
		fmt.Printf("0x%x\n", ret)
		if err != nil {
			fmt.Printf(" error: %v\n", err)
		}
	}

	return nil
}
