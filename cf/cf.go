package cf

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/runner"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var Cf = func(args ...string) *Session {
	cmdStarter := runner.NewCommandStarter()
	request, err := cmdStarter.Start("cf", args...)
	Expect(err).NotTo(HaveOccurred())

	return request
}
