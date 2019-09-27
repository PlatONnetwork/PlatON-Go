package main

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/life/compiler"
	"github.com/PlatONnetwork/PlatON-Go/life/exec"
	"github.com/PlatONnetwork/PlatON-Go/life/resolver"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

// The runner used in the unit test is mainly responsible for testing the Platonlib c++ library.
// The wasm file is executed according to the dir scan directory.
// The wasm is entered from the main entry.
// The db is created according to --outdir.
// The test tool judges the test result based on the log information.
var (
	resultReg = regexp.MustCompile(`([\d]+)\s+tests,\s+([\d]+)\s+assertions,\s+([\d]+)\s+failures`)

	testDBName = "/testdb"
	//dbPathFlag = cli.StringFlag{
	//	Name: "dbpath",
	//	Usage: "unittest leveldb path",
	//}
	//
	//logFileFlag = cli.StringFlag{
	//	Name: "logfile",
	//	Usage: "unittest log path",
	//}

	testDirFlag = cli.StringFlag{
		Name:  "dir",
		Usage: "unittest directory",
	}

	outDirFlag = cli.StringFlag{
		Name:  "outdir",
		Usage: "unittest output directory",
	}
)

var unittestCommand = cli.Command{
	Action:    unittestCmd,
	Name:      "unittest",
	Usage:     "executes the given unit tests",
	ArgsUsage: "<dir>",
	Flags: []cli.Flag{
		testDirFlag,
		outDirFlag,
	},
}

func unittestCmd(ctx *cli.Context) error {
	testDir := ctx.String(testDirFlag.Name)
	outDir := ctx.String(outDirFlag.Name)

	dbPath := outDir + testDBName

	logStream := bytes.NewBuffer(make([]byte, 65535))

	err := runTestDir(testDir, dbPath, logStream)

	if err != nil {
		return err
	}

	fmt.Println("all test pass")

	return nil
}

func runTestDir(testDir, dbPath string, logStream *bytes.Buffer) (retErr error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(logStream.String())
			retErr = err.(error)
		}
	}()
	files, err := ioutil.ReadDir(testDir)

	if err != nil {
		return errors.New(fmt.Sprintf("read unittest dir failed :%v", err))
	}

	for _, fi := range files {
		logStream.Reset()

		if fi.IsDir() {
			runTestDir(testDir+"/"+fi.Name(), dbPath, logStream)
		} else if path.Ext(fi.Name()) == ".wasm" {
			fmt.Println("exec unittest file:" + fi.Name())
			os.RemoveAll(dbPath)
			db, err := leveldb.OpenFile(dbPath, nil)
			if err != nil {
				return errors.New(fmt.Sprintf("open leveldb %s failed :%v", dbPath, err))
			}

			code, err := ioutil.ReadFile(testDir + "/" + fi.Name())

			if err != nil {
				return errors.New(fmt.Sprintf("open test file %s error :%v", fi.Name(), err))
			}

			runTest(code, db, logStream)

			states := resultReg.FindStringSubmatch(logStream.String())
			//fmt.Println("[",logStream.String(), "]")
			if len(states) != 4 {
				fmt.Println(logStream.String())
				return errors.New(fmt.Sprintf("%s unittest output result error : need 3 states such as []tests []assertions []failures", fi.Name()))
			}

			if states[3] != "0" {
				return errors.New(fmt.Sprintf("unittest :%s error \n, %s", fi.Name(), logStream.String()))
			}

			db.Close()

		}
	}
	return nil
}

func newContext(logStream *bytes.Buffer) *exec.VMContext {
	logger := log.New("wasm")
	logger.SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StreamHandler(logStream, log.FormatFunc(func(r *log.Record) []byte {
		return []byte(r.Msg)
	}))))

	wasmLog := vm.NewWasmLogger(vm.Config{Debug: true}, logger)

	return &exec.VMContext{
		Config: exec.VMConfig{
			EnableJIT:          false,
			DynamicMemoryPages: 16,
			MaxMemoryPages:     256,
			MaxTableSize:       65536,
			MaxValueSlots:      10000,
			MaxCallStackDepth:  512,
			DefaultMemoryPages: 128,
			DefaultTableSize:   65536,
			GasLimit:           1000000000000,
		},
		Addr:     [20]byte{},
		GasLimit: 1000000000000,
		StateDB:  nil,
		Log:      wasmLog,
	}
}

func runMain(wasm *exec.VirtualMachine) error {
	entryID, ok := wasm.GetFunctionExport("_Z4mainiPPc")
	if !ok {
		return errors.New("find main error")
	}

	_, err := wasm.Run(entryID, 0, 0)

	if logger, ok := wasm.Context.Log.(*vm.WasmLogger); ok {
		logger.Flush()
	}

	wasm.Stop()
	if err != nil {
		return err
	}

	return nil
}

func runModule(m *compiler.Module, functionCode []compiler.InterpreterCode, db *leveldb.DB, logStream *bytes.Buffer) error {
	context := newContext(logStream)

	wasm, err := exec.NewVirtualMachineWithModule(m, functionCode, context, newUnitTestResolver(db, logStream), nil)

	if err != nil {
		return err
	}
	return runMain(wasm)
}

func runTest(code []byte, db *leveldb.DB, logStream *bytes.Buffer) error {
	context := newContext(logStream)


	wasm, err := exec.NewVirtualMachine(code, context, newUnitTestResolver(db, logStream), nil)

	if err != nil {
		return err
	}
	return runMain(wasm)
}

type UnitTestResolver struct {
	db        *leveldb.DB
	resolver  exec.ImportResolver
	funcs     map[string]map[string]*exec.FunctionImport
	logStream *bytes.Buffer
}

func newUnitTestResolver(ldb *leveldb.DB, logStream *bytes.Buffer) *UnitTestResolver {

	constGasFunc := func(vm *exec.VirtualMachine) (uint64, error) {
		return 1, nil
	}

	resolver := &UnitTestResolver{
		db:        ldb,
		resolver:  resolver.NewResolver(0x01),
		logStream: logStream,
	}
	resolver.funcs = map[string]map[string]*exec.FunctionImport{
		"env": {
			"setState":                 &exec.FunctionImport{Execute: resolver.envSetState, GasCost: constGasFunc},
			"getState":                 &exec.FunctionImport{Execute: resolver.envGetState, GasCost: constGasFunc},
			"getStateSize":             &exec.FunctionImport{Execute: resolver.envGetStateSize, GasCost: constGasFunc},
			"getTestLog":               &exec.FunctionImport{Execute: resolver.envGetTestLog, GasCost: constGasFunc},
			"getTestLogSize":           &exec.FunctionImport{Execute: resolver.envGetTestLogSize, GasCost: constGasFunc},
			"clearLog":                 &exec.FunctionImport{Execute: resolver.envClearLog, GasCost: constGasFunc},
			"setStateDB":               &exec.FunctionImport{Execute: resolver.envSetStateDB, GasCost: constGasFunc},
			"platonCallString":         &exec.FunctionImport{Execute: resolver.envPlatonCall, GasCost: constGasFunc},
			"platonCallInt64":          &exec.FunctionImport{Execute: resolver.envPlatonCall, GasCost: constGasFunc},
			"platonDelegateCallString": &exec.FunctionImport{Execute: resolver.envPlatonCall, GasCost: constGasFunc},
			"platonDelegateCallInt64":  &exec.FunctionImport{Execute: resolver.envPlatonCall, GasCost: constGasFunc},
			"platonCall":               &exec.FunctionImport{Execute: resolver.envPlatonCall, GasCost: constGasFunc},
			"platonDelegateCall":       &exec.FunctionImport{Execute: resolver.envPlatonCall, GasCost: constGasFunc},
			"emitEvent":                &exec.FunctionImport{Execute: resolver.envEmitEvent, GasCost: constGasFunc},
			"bigintAdd":                &exec.FunctionImport{Execute: resolver.envBigintAdd, GasCost: constGasFunc},
			"envMalloc":                &exec.FunctionImport{Execute: resolver.envMalloc, GasCost: constGasFunc},
			"envFree":                  &exec.FunctionImport{Execute: resolver.envFree, GasCost: constGasFunc},
		},
	}
	return resolver
}

func (r *UnitTestResolver) envMalloc(vm *exec.VirtualMachine) int64 {

	mem := vm.Memory
	size := int(uint32(vm.GetCurrentFrame().Locals[0]))
	pos := mem.Malloc(size)

	return int64(pos)
}

func (r *UnitTestResolver) envFree(vm *exec.VirtualMachine) int64 {
	mem := vm.Memory
	offset := int(uint32(vm.GetCurrentFrame().Locals[0]))
	err := mem.Free(offset)
	if err != nil {
		return -1
	}
	return 0
}

func (r *UnitTestResolver) ResolveFunc(module, field string) *exec.FunctionImport {
	if m, exist := r.funcs[module]; exist == true {
		if f, exist := m[field]; exist == true {
			return f
		}
	}
	return r.resolver.ResolveFunc(module, field)
}

func (r *UnitTestResolver) ResolveGlobal(module, field string) int64 {
	return r.resolver.ResolveGlobal(module, field)
}

func (r *UnitTestResolver) envSetState(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(int32(vm.GetCurrentFrame().Locals[2]))
	valueLen := int(int32(vm.GetCurrentFrame().Locals[3]))
	r.db.Put(vm.Memory.Memory[key:key+keyLen], vm.Memory.Memory[value:value+valueLen], nil)
	return 0
}

func (r *UnitTestResolver) envGetState(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	value := int(int32(vm.GetCurrentFrame().Locals[2]))
	valueLen := int(int32(vm.GetCurrentFrame().Locals[3]))

	val, err := r.db.Get(vm.Memory.Memory[key:key+keyLen], nil)
	if len(val) > valueLen || err != nil {
		return 0
	}

	copy(vm.Memory.Memory[value:value+valueLen], val)
	return 0
}

func (r *UnitTestResolver) envGetStateSize(vm *exec.VirtualMachine) int64 {
	key := int(int32(vm.GetCurrentFrame().Locals[0]))
	keyLen := int(int32(vm.GetCurrentFrame().Locals[1]))

	val, err := r.db.Get(vm.Memory.Memory[key:key+keyLen], nil)
	if err != nil {
		return 0
	}
	return int64(len(val))
}

func (r *UnitTestResolver) envGetTestLog(v *exec.VirtualMachine) int64 {
	if logger, ok := v.Context.Log.(*vm.WasmLogger); ok {
		logger.Flush()
	}
	data := int(int32(v.GetCurrentFrame().Locals[0]))
	dataLen := int(int32(v.GetCurrentFrame().Locals[1]))
	size := len(r.logStream.Bytes())
	if dataLen < size {
		panic("out of buffer")
	}

	copy(v.Memory.Memory[data:data+dataLen], r.logStream.Bytes())

	return int64(size)
}

func (r *UnitTestResolver) envGetTestLogSize(v *exec.VirtualMachine) int64 {
	if logger, ok := v.Context.Log.(*vm.WasmLogger); ok {
		logger.Flush()
	}
	return int64(len(r.logStream.Bytes()))
}

func (r *UnitTestResolver) envClearLog(vm *exec.VirtualMachine) int64 {
	r.logStream.Reset()
	return 0
}

func (r *UnitTestResolver) envSetStateDB(vm *exec.VirtualMachine) int64 {
	data := int(int32(vm.GetCurrentFrame().Locals[0]))
	dataLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	t := testState{}
	err := json.Unmarshal([]byte(vm.Memory.Memory[data:data+dataLen]), &t)
	if err != nil {
		return -1
	}

	vm.Context.StateDB = &stateDB{
		state: &t,
	}
	return 0
}

func (r *UnitTestResolver) envPlatonCall(vm *exec.VirtualMachine) int64 {
	addr := int(int32(vm.GetCurrentFrame().Locals[0]))
	params := int(int32(vm.GetCurrentFrame().Locals[1]))
	paramsLen := int(int32(vm.GetCurrentFrame().Locals[2]))
	vm.Context.Log.Debug(hex.EncodeToString(vm.Memory.Memory[addr : addr+20]))
	vm.Context.Log.Debug(" ")
	vm.Context.Log.Debug(hex.EncodeToString(vm.Memory.Memory[params : params+paramsLen]))
	return 0
}

func (r *UnitTestResolver) envEmitEvent(vm *exec.VirtualMachine) int64 {
	topic := int(int32(vm.GetCurrentFrame().Locals[0]))
	topicLen := int(int32(vm.GetCurrentFrame().Locals[1]))
	data := int(int32(vm.GetCurrentFrame().Locals[2]))
	dataLen := int(int32(vm.GetCurrentFrame().Locals[3]))
	vm.Context.Log.Debug(string(vm.Memory.Memory[topic : topic+topicLen]))
	vm.Context.Log.Debug(" ")
	vm.Context.Log.Debug(hex.EncodeToString(vm.Memory.Memory[data : data+dataLen]))

	return 0
}

func (r *UnitTestResolver) envBigintAdd(vm *exec.VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	src := int(int32(frame.Locals[0]))
	srcLen := int(int32(frame.Locals[1]))
	dst := int(int32(frame.Locals[2]))
	dstLen := int(int32(frame.Locals[3]))

	i := new(big.Int)
	i.SetBytes(vm.Memory.Memory[src : src+srcLen])

	ii := new(big.Int)
	ii.SetUint64(1)
	i = i.Add(i, ii)

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(i.Bytes())))

	copy(vm.Memory.Memory[dst:], i.Bytes())
	copy(vm.Memory.Memory[dstLen:], buf)
	return 0
}

type testState struct {
	GasPrice    int64             `json:"gasPrice,omitempty"`
	BlockHash   map[string]string `json:"blockHash,omitempty"`
	Number      int64             `json:"number,omitempty"`
	GasLimit    uint64            `json:"gasLimit,omitempty"`
	Timestamp   int64             `json:"timestamp,omitempty"`
	CoinBase    string            `json:"coinbase,omitempty"`
	Balance     int64             `json:"balance,omitempty"`
	Origin      string            `json:"origin,omitempty"`
	Caller      string            `json:"caller,omitempty"`
	Value       int64             `json:"value,omitempty"`
	Address     string            `json:"address,omitempty"`
	CallerNonce int64             `json:"nonce,omitempty"`
	Account     map[string]int64  `json:"account,omitempty"`
}

type stateDB struct {
	vm.StateDB
	state *testState
}

func (s *stateDB) GasPrice() int64 {
	return s.state.GasPrice
}
func (s *stateDB) BlockHash(num uint64) common.Hash {
	hash, e := s.state.BlockHash[strconv.FormatUint(num, 10)]
	if !e {
		return common.Hash{}
	}
	return common.HexToHash(hash)
}
func (s *stateDB) BlockNumber() *big.Int {
	return big.NewInt(s.state.Number)
}
func (s *stateDB) GasLimimt() uint64 {
	return s.state.GasLimit
}
func (s *stateDB) Time() *big.Int {
	return big.NewInt(s.state.Timestamp)
}
func (s *stateDB) Coinbase() common.Address {
	return common.HexToAddress(s.state.CoinBase)
}
func (s *stateDB) GetBalance(addr common.Address) *big.Int {
	//fmt.Println("addr:", addr.Hex())
	balance, ok := s.state.Account[strings.ToLower(addr.Hex())]
	if !ok {
		return big.NewInt(0)
	}
	return big.NewInt(balance)
}

func (s *stateDB) Origin() common.Address {
	return common.HexToAddress(s.state.Origin)
}
func (s *stateDB) Caller() common.Address {
	return common.HexToAddress(s.state.Caller)
}
func (s *stateDB) Address() common.Address {
	return common.HexToAddress(s.state.Address)
}
func (s *stateDB) CallValue() *big.Int {
	return big.NewInt(s.state.Value)
}
func (s *stateDB) AddLog(address common.Address, topics []common.Hash, data []byte, bn uint64) {

}
func (s *stateDB) SetState(key []byte, value []byte) {

}
func (s *stateDB) GetState(key []byte) []byte {
	return nil
}

func (s *stateDB) GetCallerNonce() int64 {
	return s.state.CallerNonce
}
func (s *stateDB) Transfer(addr common.Address, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	fromBalance, ok := s.state.Account[s.state.Address]
	if !ok {
		return nil, 0, fmt.Errorf("not found account %s", s.state.Address)
	}
	fBigBalance := new(big.Int)
	fBigBalance.SetInt64(fromBalance)
	if fBigBalance.Cmp(value) == -1 {
		return nil, 0, errors.New("amount not enough")
	}


	tBalance, ok := s.state.Account[strings.ToLower(addr.Hex())]
	if !ok {
		return nil, 0, fmt.Errorf("not found account %s", strings.ToLower(addr.Hex()))
	}
	tBigBalance := new(big.Int)
	tBigBalance.SetInt64(tBalance)

	tBigBalance= tBigBalance.Add(tBigBalance, value)
	s.state.Account[strings.ToLower(addr.Hex())] = tBigBalance.Int64()

	fBigBalance = fBigBalance.Sub(fBigBalance, value)
	s.state.Balance = fBigBalance.Int64()
	s.state.Account[s.state.Address] = fBigBalance.Int64()

	return nil, 0, nil
}
func (s *stateDB) Call(addr, params []byte) ([]byte, error) {
	return nil, nil
}
func (s *stateDB) DelegateCall(addr, params []byte) ([]byte, error) {
	return nil, nil
}


//func (s *stateDB) CreateAccount(common.Address){}
//
//func (s *stateDB) SubBalance(common.Address, *big.Int){}
//func (s *stateDB) AddBalance(common.Address, *big.Int){}
////func (s *stateDB) GetBalance(common.Address) *big.Int{return nil}
//
//func (s *stateDB) GetNonce(common.Address) uint64{return 0}
//func (s *stateDB) SetNonce(common.Address, uint64){}
//
//func (s *stateDB) GetCodeHash(common.Address) common.Hash{return common.Hash{}}
//func (s *stateDB) GetCode(common.Address) []byte{return nil}
//func (s *stateDB) SetCode(common.Address, []byte){}
//func (s *stateDB) GetCodeSize(common.Address) int{return 0}

//// todo: new func for abi of contract.
//func (stateDB) GetAbiHash(common.Address) common.Hash{return common.Hash{}}
//func (stateDB) GetAbi(common.Address) []byte{return nil}
//func (stateDB) SetAbi(common.Address, []byte){}
//
//func (stateDB) AddRefund(uint64){}
//func (stateDB) SubRefund(uint64){}
//func (stateDB) GetRefund() uint64{return 0}
//
//func (stateDB) GetCommittedState(common.Address, []byte) []byte{return nil}
//func (stateDB) GetState(common.Address, []byte) []byte{return []byte("world+++++++**")}
//func (stateDB) SetState(common.Address, []byte, []byte){}
//func (stateDB) Suicide(common.Address) bool{return true}
//func (stateDB) HasSuicided(common.Address) bool{return true}
//
//// Exist reports whether the given account exists in state.
//// Notably this should also return true for suicided accounts.
//func (stateDB) Exist(common.Address) bool {return true}
//// Empty returns whether the given account is empty. Empty
//// is defined according to EIP161 (balance = nonce = code = 0).
//func (stateDB) Empty(common.Address) bool {return true}
//
//func (stateDB) RevertToSnapshot(int){}
//func (stateDB) Snapshot() int {return 0}
//
//func (stateDB) AddPreimage(common.Hash, []byte){}
//
//func (stateDB) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool){}
//func (stateDB) Address() common.Address {
//	return common.Address{}
//}
//
//func (stateDB)  BlockHash(num uint64) common.Hash {
//	return common.Hash{}
//}
//
//func (stateDB) BlockNumber() *big.Int {
//	return big.NewInt(0)
//}
//func (stateDB) AddLog(*types.Log) {
//	fmt.Println("add log")
//}
