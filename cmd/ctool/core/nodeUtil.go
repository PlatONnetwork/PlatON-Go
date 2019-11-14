package core

import (
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/internal/cmdtest"
)

type testPlatON struct {
	*cmdtest.TestCmd

	// template variables for expect
	Datadir string
}

func runPlatON(t *testing.T, args ...string) *testPlatON {
	tt := &testPlatON{}
	tt.TestCmd = cmdtest.NewTestCmd(t, tt)
	for i, arg := range args {
		switch {
		case arg == "-datadir" || arg == "--datadir":
			if i < len(args)-1 {
				tt.Datadir = args[i+1]
			}
		}
	}
	if tt.Datadir == "" {
		tt.Datadir = tmpdir(t)
		tt.Cleanup = func() { os.RemoveAll(tt.Datadir) }
		args = append([]string{"-datadir", tt.Datadir}, args...)
		// Remove the temporary datadir if something fails below.
		defer func() {
			if t.Failed() {
				tt.Cleanup()
			}
		}()
	}

	// Boot "platon". This actually runs the test binary but the TestMain
	// function will prevent any tests from running.
	tt.Run("platon-test", args...)

	return tt
}
