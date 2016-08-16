package helpers

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func Run(executable string, args ...string) *gexec.Session {
	cmdStarter := runner.NewCommandStarter()
	reporter := runner.NewDefaultReporter()
	session, err := cmdStarter.Start(reporter, executable, args...)
	Expect(err).NotTo(HaveOccurred())
	return session
}
