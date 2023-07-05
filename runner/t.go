package runner

import (
	"fmt"
	"runtime"

	"github.com/stretchr/testify/require"
)

type T struct {
	// Name of the test. Method name for tests or Run name argument for subtests
	Name string
	// List of any failures
	Failures []string
	SubTests []*T
	// Set by t.Skip()
	skipped       bool
	skippedReason string
}

var _ require.TestingT = (*T)(nil)

// Run runs f as a subtest of t with name. The subtest runs on the same goroutine, always returning true (has to
// satisfy the interface). T.parallel() won't work here.
func (t *T) Run(name string, f func(t *T)) bool {
	subT := &T{
		Name: t.Name + "." + name,
	}
	t.SubTests = append(t.SubTests, subT)

	doneC := make(chan struct{})
	go func() {
		finished := false
		defer func() {
			err := recover()
			if !finished && err == nil {
				// test called panic(nil) (which is weird) or Goexit (expected from FailNow or Skip)
				close(doneC)
			}
			if err != nil {
				// record when tests panic
				subT.Failures = append(t.Failures, fmt.Sprintf("Test paniked: %s", err))
				close(doneC)
			}
		}()

		f(subT)
		finished = true
		close(doneC)
	}()

	<-doneC

	return true
}

func (t *T) Errorf(format string, args ...interface{}) {
	for _, arg := range args {
		t.Failures = append(t.Failures, fmt.Sprintf(format, arg))
	}
}

func (t *T) FailNow() {
	runtime.Goexit()
}

func (t *T) Skip(msg string) {
	t.skipped = true
	t.skippedReason = msg
	runtime.Goexit()
}
