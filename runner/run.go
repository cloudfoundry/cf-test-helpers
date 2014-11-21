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

const timeFormat = "2006-01-02 15:04:00 (MST)"

var CommandInterceptor = func(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}

func Run(executable string, args ...string) *gexec.Session {
	cmd := exec.Command(executable, args...)

	startTime := time.Now()
	sayCommandWillRun(startTime, cmd)

	sess, err := gexec.Start(CommandInterceptor(cmd), ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	sayCommandDidRun(startTime, time.Now(), cmd)

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

func sayCommandDidRun(startTime, finishedTime time.Time, cmd *exec.Cmd) {
	startColor := ""
	endColor := ""
	if !config.DefaultReporterConfig.NoColor {
		startColor = "\x1b[33m"
		endColor = "\x1b[0m"
	}
	elapsedTime := finishedTime.Sub(startTime)
	fmt.Fprintf(ginkgo.GinkgoWriter, "\n%s[Completed: %s]> %s %s\n", startColor, elapsedTime, strings.Join(cmd.Args, " "), endColor)
}
