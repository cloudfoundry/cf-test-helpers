package fakes

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/internal"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega/gexec"
)

type callToStartMethod struct {
	Executable string
	Args       []string
	Reporter   internal.Reporter
}

type startMethodStub struct {
	Output    string
	Stderr    string
	Err       error
	ExitCode  int
	SleepTime time.Duration
}

type FakeCmdStarter struct {
	CalledWith        []callToStartMethod
	ToReturn          []startMethodStub
	TotalCallsToStart int
}

func NewFakeCmdStarter() *FakeCmdStarter {
	return &FakeCmdStarter{
		CalledWith: []callToStartMethod{},
		ToReturn:   make([]startMethodStub, 10),
	}
}

func (s *FakeCmdStarter) Start(reporter internal.Reporter, executable string, args ...string) (*gexec.Session, error) {
	output := s.ToReturn[s.TotalCallsToStart].Output
	if output == "" {
		output = `\{\}`
	}
	stderr := s.ToReturn[s.TotalCallsToStart].Stderr
	if stderr == "" {
		stderr = `\{\}`
	}
	sleepTime := s.ToReturn[s.TotalCallsToStart].SleepTime
	exitCode := s.ToReturn[s.TotalCallsToStart].ExitCode
	err := s.ToReturn[s.TotalCallsToStart].Err

	s.TotalCallsToStart += 1

	callToStart := callToStartMethod{
		Executable: executable,
		Args:       args,
		Reporter:   reporter,
	}
	s.CalledWith = append(s.CalledWith, callToStart)

	reporter.Report(false, time.Now(), exec.Command(executable, args...))
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"echo %[1]s; echo %[2]s > /dev/stderr; [ %[4]d -eq 0 ] && sleep %[3]f; exit %[4]d",
			output,
			stderr,
			sleepTime.Seconds(),
			exitCode,
		),
	)
	session, _ := gexec.Start(cmd, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	return session, err
}
