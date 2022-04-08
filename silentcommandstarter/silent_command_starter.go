package silentcommandstarter

import (
	"os/exec"
	"time"

	"github.com/KevinJCross/cf-test-helpers/v2/internal"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/gexec"
)

type CommandStarter struct {
}

func NewCommandStarter() *CommandStarter {
	return &CommandStarter{}
}

func (r *CommandStarter) Start(reporter internal.Reporter, executable string, args ...string) (*gexec.Session, error) {
	cmd := exec.Command(executable, args...)
	reporter.Report(false, time.Now(), cmd)

	writer := ginkgo.GinkgoWriter
	_, _ = writer.Write([]byte("SILENCING COMMAND OUTPUT"))

	return gexec.Start(cmd, nil, nil)
}
