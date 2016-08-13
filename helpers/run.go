package helpers

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func Run(executable string, args ...string) *gexec.Session {
	cmdStarter := runner.NewCommandStarter()
	session, err := cmdStarter.Start(executable, args...)
	Expect(err).NotTo(HaveOccurred())
	return session
}
