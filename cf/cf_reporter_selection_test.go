package cf

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/cf-test-helpers/v2/commandreporter"
)

var _ = Describe("defaultReporter", func() {
	It("returns a plain command reporter (observer teeing now lives in internal)", func() {
		Expect(defaultReporter()).To(BeAssignableToTypeOf(commandreporter.NewCommandReporter()))
	})
})
