package cfinternal

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func Cf(cmdStarter starter, args ...string) *gexec.Session {
	reporter := runner.NewDefaultReporter()
	request, err := cmdStarter.Start(reporter, "cf", args...)
	Expect(err).NotTo(HaveOccurred())

	return request
}
