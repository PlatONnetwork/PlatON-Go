package main

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/lru"
	"github.com/PlatONnetwork/PlatON-Go/life/exec"
	"bytes"
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"time"
)

// 单元测试使用的runner, 主要负责测试Platonlib c++ 库，根据dir 扫描目录下wasm文件执行 wasm从main入口进入，根据--outdir创建db, 测试工具根据日志信息判断测试结果
var (
	wasmFileFlag = cli.StringFlag{
		Name:  "file",
		Usage: "wasm file",
	}

	loopFlag = cli.IntFlag{
		Name:  "loop",
		Usage: "execute count",
		Value: 1,
	}

	profFlag = cli.StringFlag{
		Name:  "prof",
		Usage: "write cpuprofie to file",
	}
)

var benchmarkCommand = cli.Command{
	Action:    benchmarkCmd,
	Name:      "benchmark",
	Usage:     "benchmark wasm vm",
	ArgsUsage: "<dir>",
	Flags: []cli.Flag{
		wasmFileFlag,
		outDirFlag,
		loopFlag,
		profFlag,
	},
	HideHelp: false,
}

func benchmarkCmd(ctx *cli.Context) error {
	wasmFile := ctx.String(wasmFileFlag.Name)
	outDir := ctx.String(outDirFlag.Name)
	loop := ctx.Int(loopFlag.Name)
	profFile := ctx.String(profFlag.Name)

	if profFile != "" {
		f, err := os.Create(profFile)
		if err != nil {
			return err
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	start := time.Now()

	err := benchmark(wasmFile, outDir, loop)

	if err != nil {
		return err
	}

	end := time.Now()

	fmt.Println("execute time:", end.Sub(start).String(), "loop:", loop)

	return nil
}
func benchmark(wasmFile string, outDir string, loop int) error {
	dbPath := outDir + testDBName
	logStream := bytes.NewBuffer(make([]byte, 65535))
	os.RemoveAll(dbPath)
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("open leveldb %s failed :%v", dbPath, err))
	}
	code, err := ioutil.ReadFile(wasmFile)
	if err != nil {
		return err
	}

	m, functionCode, err := exec.ParseModuleAndFunc(code, nil)

	if err != nil {
		return err
	}

	if err := lru.SetWasmDB(outDir); err != nil {
		return err
	}

	addr := common.HexToAddress("0x43355c787c50b647c425f594b441d4bd751951c1")
	lru.WasmCache().Add(addr, &lru.WasmModule{m, functionCode})

	for i := 0; i < loop; i++ {
		m, ok := lru.WasmCache().Get(addr)
		if !ok {
			return errors.New("get wasm cache error")
		}
		if err := runModule(m.Module, m.FunctionCode, db, logStream); err != nil {
			return err
		}
		//lru.WasmCache().Purge()
	}
	db.Close()
	return nil
}
