package fakes

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/internal"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/gexec"
)

type FakeCmdStarter struct {
	CalledWith struct {
		Executable string
		Args       []string
		Reporter   internal.Reporter
	}
	ToReturn struct {
		Output    string
		Err       error
		ExitCode  int
		SleepTime int
	}
}

func (s *FakeCmdStarter) Start(reporter internal.Reporter, executable string, args ...string) (*gexec.Session, error) {
	s.CalledWith.Executable = executable
	s.CalledWith.Args = args
	s.CalledWith.Reporter = reporter

	// Default return values
	if s.ToReturn.Output == "" {
		s.ToReturn.Output = `\{\}`
	}

	reporter.Report(time.Now(), exec.Command(executable, args...))
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"echo %s; sleep %d; exit %d",
			s.ToReturn.Output,
			s.ToReturn.SleepTime,
			s.ToReturn.ExitCode,
		),
	)
	session, _ := gexec.Start(cmd, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	return session, s.ToReturn.Err
}
