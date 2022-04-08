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

type CommandReporter struct {
	Writer io.Writer
}

func NewCommandReporter(writers ...io.Writer) *CommandReporter {
	var writer io.Writer
	switch len(writers) {
	case 0:
		writer = ginkgo.GinkgoWriter
	case 1:
		writer = writers[0]
	default:
		panic("newCommandReporter can only take one writer")
	}

	return &CommandReporter{
		Writer: writer,
	}
}

func (r *CommandReporter) Report(withColour bool, startTime time.Time, cmd *exec.Cmd) {
	PrintCommand(withColour, startTime, r.Writer, strings.Join(cmd.Args, " "))
}

func PrintCommand(withColour bool, startTime time.Time, writer io.Writer, command string) {
	startColor := ""
	endColor := ""
	startBold := ""
	if withColour {
		startColor = "\x1b[32m"
		startBold = "\x1b[32;1m"
		endColor = "\x1b[0m"
	}
	_, _ = fmt.Fprintf(
		writer,
		"\n%s[%s]> %s%s %s\n",
		startColor,
		startTime.UTC().Format(timeFormat),
		startBold,
		command,
		endColor,
	)
}
