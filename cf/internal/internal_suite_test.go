package cfinternal_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

type fakeStarter struct {
	calledWith struct {
		executable string
		args       []string
	}
	toReturn struct {
		output    string
		err       error
		exitCode  int
		sleepTime int
	}
}

func (s *fakeStarter) Start(executable string, args ...string) (*gexec.Session, error) {
	s.calledWith.executable = executable
	s.calledWith.args = args

	// Default return values
	if s.toReturn.output == "" {
		s.toReturn.output = `\{\}`
	}
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"echo %s; sleep %d; exit %d",
			s.toReturn.output,
			s.toReturn.sleepTime,
			s.toReturn.exitCode,
		),
	)
	session, _ := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	return session, s.toReturn.err
}

func TestInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CF Internal Suite")
}
