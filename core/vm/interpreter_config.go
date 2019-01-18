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
	// Enable recording of SHA3/keccak preimages
	EnablePreimageRecording bool
	// JumpTable contains the EVM instruction table. This
	// may be left uninitialised and will be set to the default table.
	JumpTable [256]operation

	// Type of the EWASM interpreter
	EWASMInterpreter string
	// Type of the EVM interpreter
	EVMInterpreter string

	ConsoleOutput bool
}