package command

// Executor represents a command builder and executor.
// It allows creating a command with an initial set of arguments 
// via Execute(), and then modifying it (e.g., appending arguments, 
// changing directory) before execution.
type Executor interface {
	// Execute creates a new command with the specified name and arguments.
	// It returns itself to allow method chaining.
	Execute(name string, args ...string) Executor
	Runner
}

// Runner represents an executable command instance.
// It provides methods to configure the command and execute it.
type Runner interface {
	// Run executes the command and returns an error if it fails.
	Run() error

	// Path returns the absolute path of the executable command.
	Path() string

	// CombinedOutput runs the command and returns its combined standard
	// output and standard error.
	CombinedOutput() ([]byte, error)

	// WithPath sets the absolute path to the executable binary.
	WithPath(path string) Runner

	// WithDir sets the working directory for the command execution.
	WithDir(dirPath string) Runner

	// AppendArgs adds additional arguments to the existing command.
	AppendArgs(args ...string) Runner
}
