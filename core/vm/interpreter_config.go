package vm

// Config are the configuration options for the interpreter
type Config struct {
	// Debug enable debugging Interpreter options
	Debug bool
	// Tracer is the op code logger
	Tracer Tracer
	// NoRecursion disabled interpreter call, callcode,
	// delegate call and create
	NoRecursion bool

	// JumpTable contains the EVM instruction table. This
	// may be left uninitialised and will be set to the default table.
	JumpTable [256]operation

	ConsoleOutput bool

	// The actual implementation type of the wasm instance
	// This option is used in the configuration or command line
	WasmType WasmInsType

	// VM execution timeout duration (unit: ms)
	VmTimeoutDuration uint64
}