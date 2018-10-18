package main

import (
	"Platon-go/life/exec"
	"Platon-go/life/resolver"
	"Platon-go/lifer"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	//entryFunctionFlag := flag.String("entry", "app_main", "entry function id")
	//dynamicPages := flag.Int("dynamicPages", 1, "dynamic memory pages")

	//jitFlag := flag.Bool("jit", false, "enable jit")
	//flag.Parse()

	// mocking test
	flag := false
	pages := 1
	functionFlag := "transfer"
	jitFlag := &flag
	dynamicPages := &pages
	entryFunctionFlag := &functionFlag

	rl := resolver.NewResolver(0x01)
	// Read WebAssembly *.wasm file.
	//input, err := ioutil.ReadFile(flag.Arg(0))
	input, err := ioutil.ReadFile("C:\\sunzone\\MyDocument\\liteide\\src\\Platon-go\\core\\vm\\life\\contract\\hello.wasm")
	//fmt.Println(common.ToHex(input))
	if err != nil {
		panic(err)
	}

	// Instantiate a new WebAssembly VM with a few resolved imports.
	vm, err := exec.NewVirtualMachine(input, exec.VMContext{
		Config: exec.VMConfig{
			EnableJIT:          *jitFlag,
			DefaultMemoryPages: 128,
			DefaultTableSize:   65536,
			DynamicMemoryPages: *dynamicPages,
		},
		Addr:     [20]byte{},
		GasUsed:  0,
		GasLimit: 20000000,
	}, rl, nil)

	if err != nil {
		panic(err)
	}


	*entryFunctionFlag = "transfer"
	// Get the function ID of the entry function to be executed.
	entryID, ok := vm.GetFunctionExport(*entryFunctionFlag)
	if !ok {
		fmt.Printf("Entry function %s not found; starting from 0.\n", *entryFunctionFlag)
		entryID = 0
	}

	start := time.Now()

	// If any function prior to the entry function was declared to be
	// called by the module, run it first.
	if vm.Module.Base.Start != nil {
		startID := int(vm.Module.Base.Start.Index)
		_, err := vm.Run(startID)
		if err != nil {
			vm.PrintStackTrace()
			panic(err)
		}
	}

	// 合约方法结构：int transfer(address from, address to, int asset)
	// todo: 此处确定了ID，传入需要的参数，然后获取返回值？返回值？？？
	// Run the WebAssembly module's entry function.
	ret, err := vm.Run(entryID, resolver.MallocString(vm, "helloh汉子来了"), resolver.MallocString(vm, "world"), 45)
	if err != nil {
		vm.PrintStackTrace()
		panic(err)
	}
	end := time.Now()

	fmt.Printf("return value = %d, duration = %v  gasUsed:%d \n", ret, end.Sub(start), vm.Context.GasUsed)
}
