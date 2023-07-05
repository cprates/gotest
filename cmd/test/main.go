package main

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/cprates/gotest/runner"
)

func main() {
	l := logrus.New().WithField("loggers", "PostgresRunner")
	l.Logger.SetOutput(io.Discard) // just as an example
	postgresTests := &PostgresTests{Runner: runner.New(l)}
	// report as JSON
	postgresTests.Run(postgresTests, runner.NewJSONReport(os.Stdout, "SQLiteTests"))

	l = logrus.WithField("loggers", "SQLiteRunner")
	sqliteTests := &SQLiteTests{Runner: runner.New(l)}
	l = logrus.WithField("loggers", "SQLiteTests")
	sqliteTests.Run(sqliteTests, runner.NewLogReport(l))
}
