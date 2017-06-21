package internal

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"io"

	"github.com/onsi/ginkgo/config"
)

const timeFormat string = "2006-01-02 15:04:05.00 (MST)"

type RedactingReporter struct {
	writer    io.Writer
	redactees []string
}

func NewRedactingReporter(writer io.Writer, redactees ...string) Reporter {
	return &RedactingReporter{
		writer:    writer,
		redactees: redactees,
	}
}

func (r *RedactingReporter) Report(startTime time.Time, cmd *exec.Cmd) {
	startColor := ""
	endColor := ""
	if !config.DefaultReporterConfig.NoColor {
		startColor = "\x1b[32m"
		endColor = "\x1b[0m"
	}

	fmt.Fprintf(
		r.writer,
		"\n%s[%s]> %s %s\n",
		startColor,
		startTime.UTC().Format(timeFormat),
		r.redactCommandArgs(cmd.Args),
		endColor,
	)
}

func (r *RedactingReporter) redactCommandArgs(args []string) string {
	if len(r.redactees) == 0 {
		return strings.Join(args, " ")
	}

	var out []string
	for _, arg := range args {
		if r.shouldBeRedacted(arg) {
			out = append(out, "[REDACTED]")
		} else {
			out = append(out, arg)
		}
	}

	return strings.Join(out, " ")
}

func (r *RedactingReporter) shouldBeRedacted(val string) bool {
	for _, v := range r.redactees {
		if v == val {
			return true
		}
	}

	return false
}
