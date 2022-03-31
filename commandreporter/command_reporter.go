package commandreporter

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
)

const timeFormat = "2006-01-02 15:04:05.00 (MST)"

type StringerStruct struct {
	Time    string
	Command string
	Args    string
}

// ColorableString for ReportEntry to use
func (s StringerStruct) ColorableString() string {
	return fmt.Sprintf("{{green}}[%s]> {{bold}}%s %s{{/}}", s.Time, s.Command, s.Args)
}

// non-colorable String() is used by go's string formatting support but ignored by ReportEntry
func (s StringerStruct) String() string {
	return fmt.Sprintf("[%s]> %s %s", s.Time, s.Command, s.Args)
}

type CommandReporter struct {
	Writer io.Writer
}

func NewCommandReporter() *CommandReporter {
	return &CommandReporter{}
}

func (r *CommandReporter) Report(startTime time.Time, cmd *exec.Cmd) {
	ginkgo.AddReportEntry("CF", StringerStruct{Time: startTime.UTC().Format(timeFormat), Command: cmd.Args[0], Args: strings.Join(cmd.Args[1:], " ")}, ginkgo.ReportEntryVisibilityFailureOrVerbose)
}
