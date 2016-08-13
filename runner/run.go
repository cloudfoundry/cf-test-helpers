package runner

import (
	"os/exec"

	"github.com/onsi/gomega/gexec"
)


var SkipSSLValidation bool

func Run(executable string, args ...string) *gexec.Session {
	cmdStarter := NewCommandStarter()
	return cmdStarter.Start(executable, args...)
}
