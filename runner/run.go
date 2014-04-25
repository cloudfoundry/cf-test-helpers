package runner

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var CommandInterceptor = func(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}

func Run(executable string, args ...string) *gexec.Session {
	cmd := exec.Command(executable, args...)

	if config.DefaultReporterConfig.Verbose {
		fmt.Println("\n", "> ", strings.Join(cmd.Args, " "))
	}

	sess, err := gexec.Start(CommandInterceptor(cmd), ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return sess
}

func Curl(args ...string) *gexec.Session {
	args = append([]string{"-s"}, args...)
	return Run("curl", args...)
}
