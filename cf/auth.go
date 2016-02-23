package cf

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const timeFormat = "2006-01-02 15:04:05.00 (MST)"

var CfAuth = func(user, password string) *gexec.Session {
	cfCmd := exec.Command("cf", "auth", user, password) //Cf("auth", user, password)
	// session := runner.NewCmdRunner(cfCmd, 1*time.Second).Run()

	sayCommandWillRun(time.Now(), cfCmd)

	sess, err := gexec.Start(runner.CommandInterceptor(cfCmd), ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return sess
}

func sayCommandWillRun(startTime time.Time, cmd *exec.Cmd) {
	cfCmd := cmd.Args[0]
	authCmd := cmd.Args[1]
	user := cmd.Args[2]

	startColor := ""
	endColor := ""
	if !config.DefaultReporterConfig.NoColor {
		startColor = "\x1b[32m"
		endColor = "\x1b[0m"
	}
	fmt.Fprintf(
		ginkgo.GinkgoWriter,
		"\n%s[%s]> %s %s %s %s %s\n",
		startColor,
		startTime.UTC().Format(timeFormat),
		cfCmd,
		authCmd,
		user,
		"[REDACTED]",
		endColor)
}
