package cf_test

import (
	"testing"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/commandstarter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var originalCf = cf.Cf
var originalCommandInterceptor = commandstarter.CommandInterceptor

var _ = AfterEach(func() {
	cf.Cf = originalCf
	commandstarter.CommandInterceptor = originalCommandInterceptor
})

func TestCf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cf Suite")
}
