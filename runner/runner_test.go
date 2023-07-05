package runner_test

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	runner "github.com/cprates/gotest/runner"
)

type RunnerTests struct {
	*runner.Runner
	t            *testing.T
	testsCounter int
}

func (s *RunnerTests) Test1(t *runner.T) {
	require.Equal(s.t, "Test1", t.Name)
	s.testsCounter++
}

func (s *RunnerTests) Test2(t *runner.T) {
	require.Equal(s.t, "Test2", t.Name)
	s.testsCounter++

	t.Run("Tes2.SubTest1", func(t *runner.T) {
		s.testsCounter++
	})
}

func (s *RunnerTests) Test3(t *runner.T) {
	require.Equal(s.t, "Test3", t.Name)
	s.testsCounter++

	panic(errors.New("panic 3"))
}

func (s *RunnerTests) Test4(t *runner.T) {
	require.Equal(s.t, "Test4", t.Name)
	s.testsCounter++

	t.Run("Tes4.SubTest1", func(t *runner.T) {
		panic(errors.New("panic 4"))
	})
}

func TestRunner(t *testing.T) {
	r := &RunnerTests{t: t, Runner: runner.New(logrus.WithField("loggers", "smoketests"))}
	r.Run(r)

	require.Equal(t, 5, r.testsCounter)
	results := r.Results()
	require.Equal(t, 4, len(results))

	// Test1
	require.Empty(t, results[0].Failures)
	require.Empty(t, results[0].SubTests)

	// Test2
	require.Empty(t, results[1].Failures)
	require.Equal(t, 1, len(results[1].SubTests))

	// Test3
	require.Equal(t, 1, len(results[2].Failures))
	require.Equal(t, "Test paniked: panic 3", results[2].Failures[0])

	// Test4
	require.Equal(t, 1, len(results[3].SubTests[0].Failures))
	require.Equal(t, "Test paniked: panic 4", results[3].SubTests[0].Failures[0])
}

func TestRunnerTCollectsErrors(t *testing.T) {
	st := &runner.T{}
	assert.Equal(st, 1, 1)
	assert.Equal(st, 1, 2, "never true")
	assert.NoError(st, errors.New("dummy error"))

	require.Equal(t, 2, len(st.Failures))
	require.Contains(t, st.Failures[0], "never true")
	require.Contains(t, st.Failures[1], "dummy error")
}

func TestRunnerSubTests(t *testing.T) {
	st := &runner.T{Name: "TestRunnerSubTests"}
	st.Run("SubTest1", func(subT *runner.T) {
		assert.Equal(subT, 1, 1)
		assert.NoError(subT, errors.New("dummy error"))
	})

	require.Equal(t, 1, len(st.SubTests))
	require.Equal(t, 1, len(st.SubTests[0].Failures))
	require.Equal(t, "TestRunnerSubTests.SubTest1", st.SubTests[0].Name)
	require.Contains(t, st.SubTests[0].Failures[0], "dummy error")
	require.Equal(t, 0, len(st.Failures))
}

func TestRunnerRunsAllSubTestsAfterPanics(t *testing.T) {
	st := &runner.T{Name: "TestRunnerSubTests"}
	st.Run("SubTest1", func(subT *runner.T) {
		panic(errors.New("dummy error"))
	})
	st.Run("SubTest2", func(subT *runner.T) {
		assert.Equal(st, 1, 1)
	})

	require.Equal(t, 2, len(st.SubTests))
	require.Equal(t, 1, len(st.SubTests[0].Failures))
	require.Equal(t, "TestRunnerSubTests.SubTest1", st.SubTests[0].Name)
	require.Contains(t, st.SubTests[0].Failures[0], "dummy error")
	require.Equal(t, 0, len(st.Failures))
}
