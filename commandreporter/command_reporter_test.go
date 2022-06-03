package commandreporter_test

import (
	"bytes"
	"github.com/cloudfoundry/cf-test-helpers/commandreporter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("CommandReporter", func() {
	Describe("NewCommandReporter", func() {
		var writers []io.Writer

		Context("when no writers are provided", func() {
			BeforeEach(func() {
				writers = []io.Writer{}
			})

			It("uses the GinkgoWriter", func() {
				reporter := commandreporter.NewCommandReporter(writers...)
				writer := reporter.Writer
				Expect(writer).To(BeAssignableToTypeOf(GinkgoWriter))
			})
		})

		Context("when a single writer is provided", func() {
			BeforeEach(func() {
				writers = []io.Writer{
					&bytes.Buffer{},
				}
			})

			It("uses the provided writer", func() {
				reporter := commandreporter.NewCommandReporter(writers...)
				Expect(reporter.Writer).To(Equal(writers[0]))
			})
		})

		Context("when there is more than one writer provided", func() {
			BeforeEach(func() {
				writers = []io.Writer{
					&bytes.Buffer{},
					&bytes.Buffer{},
				}
			})

			It("panics", func() {
				Expect(func() { commandreporter.NewCommandReporter(writers...) }).To(Panic())
			})
		})
	})
})
