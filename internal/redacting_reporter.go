package internal

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/v2/commandreporter"
	"io"
	"os/exec"
	"strings"
	"time"
)

type RedactingReporter struct {
	writer   io.Writer
	redactor Redactor
}

var _ Reporter = new(RedactingReporter)

func NewRedactingReporter(writer io.Writer, redactor Redactor) *RedactingReporter {
	return &RedactingReporter{
		writer:   writer,
		redactor: redactor,
	}
}

func (r *RedactingReporter) Report(withColour bool, startTime time.Time, cmd *exec.Cmd) {
	commandreporter.PrintCommand(withColour, startTime, r.writer, r.redactor.Redact(strings.Join(cmd.Args, " ")))
}
