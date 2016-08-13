package runner

import (
	"os/exec"

	"github.com/onsi/gomega/gexec"
)

const timeFormat = "2006-01-02 15:04:05.00 (MST)"

var SkipSSLValidation bool

func Run(executable string, args ...string) *gexec.Session {
	cmdStarter := NewCommandStarter()
	return cmdStarter.Start(executable, args...)
}
