package helpers_test

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/cf-test-helpers/config"
	. "github.com/cloudfoundry/cf-test-helpers/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Artifacts", func() {
	Describe("EnableCFTrace", func() {
		var componentName string
		var config config.Config
		var expectedGinkgoNode int

		BeforeEach(func() {
			componentName = "fakeComponentName"
			config.ArtifactsDirectory = "/some/dir"
			expectedGinkgoNode = GinkgoParallelProcess()
		})

		It("Sets the CF_TRACE variable correctly", func() {
			EnableCFTrace(&config, componentName)
			Expect(os.Getenv("CF_TRACE")).To(Equal(fmt.Sprintf("%s/CATS-TRACE-%s-%d.txt", config.ArtifactsDirectory, componentName, expectedGinkgoNode)))
		})

		Context("when the component name has spaces", func() {
			BeforeEach(func() {
				componentName = "fake component name"
			})
			It("replaces them with underscores", func() {
				EnableCFTrace(&config, componentName)
				Expect(os.Getenv("CF_TRACE")).To(Equal(fmt.Sprintf("%s/CATS-TRACE-fake_component_name-%d.txt", config.ArtifactsDirectory, expectedGinkgoNode)))
			})
		})

		Context("when the ArtifactsDirectory is not set", func() {
			BeforeEach(func() {
				config.ArtifactsDirectory = ""
			})

			It("uses the current directory", func() {
				EnableCFTrace(&config, componentName)
				Expect(os.Getenv("CF_TRACE")).To(Equal(fmt.Sprintf("CATS-TRACE-%s-%d.txt", componentName, expectedGinkgoNode)))
			})
		})
	})
})
