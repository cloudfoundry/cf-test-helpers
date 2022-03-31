package internal

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/commandreporter"
	"github.com/onsi/ginkgo/v2"
	"os/exec"
	"strings"
	"time"
)

const timeFormat string = "2006-01-02 15:04:05.00 (MST)"

type RedactingReporter struct {
	redactor Redactor
}
type RedactorStruct struct {
	commandreporter.StringerStruct
	redactor Redactor
}

// ColorableString for ReportEntry to use
func (s RedactorStruct) ColorableString() string {
	return s.redactor.Redact(s.StringerStruct.ColorableString())
}

// non-colorable String() is used by go's string formatting support but ignored by ReportEntry
func (s RedactorStruct) String() string {
	return  s.redactor.Redact(s.StringerStruct.String())
}


var _ Reporter = new(RedactingReporter)

func NewRedactingReporter(redactor Redactor) *RedactingReporter {
	return &RedactingReporter{
		redactor: redactor,
	}
}

func (r *RedactingReporter) Report(startTime time.Time, cmd *exec.Cmd) {
	ginkgo.AddReportEntry("CF", RedactorStruct{StringerStruct: commandreporter.StringerStruct{Time: startTime.UTC().Format(timeFormat), Command: cmd.Args[0], Args: strings.Join(cmd.Args[1:], " ")}, redactor: r.redactor}, ginkgo.ReportEntryVisibilityFailureOrVerbose)
}
