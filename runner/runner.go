package runner

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

type Reporter interface {
	Report(tests []*T)
}

type Runner struct {
	logger *logrus.Entry

	tests []*T
}

func New(logger *logrus.Entry) *Runner {
	return &Runner{logger: logger}
}

func (r *Runner) Results() []*T {
	return r.tests
}

func (r *Runner) Run(this interface{}, reporters ...Reporter) {
	t := reflect.TypeOf(this)

	tests := []reflect.Method{}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if !validMethod(m) {
			r.logger.Debugf("Discarding method %q", m.Name)
			continue
		}

		tests = append(tests, m)
	}

	for _, test := range tests {
		t := &T{Name: test.Name}
		r.tests = append(r.tests, t)

		in := []reflect.Value{reflect.ValueOf(this), reflect.ValueOf(t)}
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
					t.Failures = append(t.Failures, fmt.Sprintf("Test paniked: %s", err))
					close(doneC)
				}
			}()

			_ = test.Func.Call(in)
			finished = true
			close(doneC)
		}()

		<-doneC
	}

	r.logger.Info("Reporting smoke tests results")
	for _, reporter := range reporters {
		reporter.Report(r.tests)
	}
}

func validMethod(m reflect.Method) bool {
	t := m.Type

	prefix := "Test"
	switch {
	case !strings.HasPrefix(m.Name, prefix) || len(m.Name) == len(prefix):
		return false
	case t.NumIn() != 2:
		return false
	case t.NumOut() != 0:
		return false
	case t.In(1) != reflect.TypeOf(&T{}):
		return false
	}

	return true
}
