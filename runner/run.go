package runner

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const timeFormat = "2006-01-02 15:04:05.00 (MST)"

var CommandInterceptor = func(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}

func Run(executable string, args ...string) *gexec.Session {
	cmd := exec.Command(executable, args...)

	sayCommandWillRun(time.Now(), cmd)

	sess, err := gexec.Start(CommandInterceptor(cmd), ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return sess
}

func Curl(args ...string) *gexec.Session {
	args = append([]string{"-s"}, args...)
	return Run("curl", args...)
}

func sayCommandWillRun(startTime time.Time, cmd *exec.Cmd) {
	startColor := ""
	endColor := ""
	if !config.DefaultReporterConfig.NoColor {
		startColor = "\x1b[32m"
		endColor = "\x1b[0m"
	}
	fmt.Fprintf(ginkgo.GinkgoWriter, "\n%s[%s]> %s %s\n", startColor, startTime.UTC().Format(timeFormat), strings.Join(cmd.Args, " "), endColor)
}

func ExecWithTimeoutForExitCode(session *gexec.Session, timeout time.Duration, expectedExitCode int) *gexec.Session {
	cmdString := strings.Join(session.Command.Args, " ")
	select {
	case <-session.Exited:
		exitCode := session.ExitCode()
		Expect(exitCode).To(Equal(expectedExitCode), "Failed executing command (exit %d):\nCommand: %s\n\n[stdout]:\n%s\n\n[stderr]:\n%s", exitCode, cmdString, string(session.Out.Contents()), string(session.Err.Contents()))
	case <-time.After(timeout):
		exitCode := session.ExitCode()
		Expect(exitCode).To(Equal(expectedExitCode), "Timed out executing command (%v):\nCommand: %s\n\n[stdout]:\n%s\n\n[stderr]:\n%s", timeout.String(), cmdString, string(session.Out.Contents()), string(session.Err.Contents()))
	}
	return session
}

func ExecWithTimeout(session *gexec.Session, timeout time.Duration) *gexec.Session {
	return ExecWithTimeoutForExitCode(session, timeout, 0)
}
