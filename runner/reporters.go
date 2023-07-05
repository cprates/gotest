package runner

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

var _ Reporter = (*LogReport)(nil)

// LogReport is a reporter to log results to a given logger.
type LogReport struct {
	Logger *logrus.Entry
}

func NewLogReport(l *logrus.Entry) *LogReport {
	return &LogReport{Logger: l}
}

func (l *LogReport) Report(results []*T) {
	for _, result := range results {
		if len(result.Failures) == 0 && len(result.SubTests) == 0 {
			if result.skipped {
				l.Logger.WithField("skip_reason", result.skippedReason).Infof("Test %s: Skipped", result.Name)
			} else {
				l.Logger.Infof("Test %s: Pass", result.Name)
			}
			continue
		}

		if len(result.Failures) > 0 {
			l.Logger.WithField("errors", result.Failures).Errorf("Test %s: Failed", result.Name)
		}

		if len(result.SubTests) > 0 {
			l.Report(result.SubTests)
		}
	}
}

var _ Reporter = (*JSONReporter)(nil)

// LogReport is a reporter to log results as JSON.
type JSONReporter struct {
	out  io.Writer
	name string
}

type JSONReportEntry struct {
	TestName   string
	Failed     bool
	Skipped    bool
	SkipReason string
	Errors     []string
}

type JSONReport struct {
	Name    string
	Results []JSONReportEntry
}

func NewJSONReport(out io.Writer, name string) *JSONReporter {
	return &JSONReporter{
		out:  out,
		name: name,
	}
}

func (j *JSONReporter) Report(results []*T) {
	report := JSONReport{
		Name:    j.name,
		Results: make([]JSONReportEntry, 0),
	}
	for _, result := range results {
		report.Results = append(report.Results, JSONReportEntry{
			TestName:   result.Name,
			Failed:     len(result.Failures) > 0,
			Skipped:    result.skipped,
			SkipReason: result.skippedReason,
			Errors:     append(result.Failures[:0:0], result.Failures...),
		})

		if len(result.SubTests) > 0 {
			report.Results = append(report.Results, jsonReportAux(result.Name, result.SubTests)...)
		}
	}

	b, err := json.Marshal(&report)
	if err != nil {
		fmt.Fprintf(j.out, "Failed to marshal report: %s", err)
		return
	}

	fmt.Fprintln(j.out, string(b))
}

func jsonReportAux(parentName string, results []*T) []JSONReportEntry {
	var report []JSONReportEntry
	for _, result := range results {
		report = append(report, JSONReportEntry{
			TestName:   parentName + "." + result.Name,
			Failed:     len(result.Failures) > 0,
			Skipped:    result.skipped,
			SkipReason: result.skippedReason,
			Errors:     append(result.Failures[:0:0], result.Failures...),
		})

		if len(result.SubTests) > 0 {
			report = append(report, jsonReportAux(result.Name, result.SubTests)...)
		}
	}

	return report
}
