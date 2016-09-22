package config_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	cfg "github.com/cloudfoundry-incubator/cf-test-helpers/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type RequiredConfig struct {
	ApiEndpoint       string `json:"api"`
	AdminUser         string `json:"admin_user"`
	AdminPassword     string `json:"admin_password"`
	SkipSSLValidation bool   `json:"skip_ssl_validation"`
	AppsDomain        string `json:"apps_domain"`
	UseHttp           bool   `json:"use_http"`
}

var config cfg.Config
var tmpFile *os.File
var err error
var _ = Describe("Config", func() {
	BeforeEach(func() {
		requiredConfig := RequiredConfig{
			ApiEndpoint:       "somewhere.over.the.rainbow",
			AdminUser:         "admin",
			AdminPassword:     "admin",
			SkipSSLValidation: true,
			AppsDomain:        "cf-app.over.the.rainbow",
			UseHttp:           true,
		}

		tmpFile, err = ioutil.TempFile("", "cf-test-helpers-config")
		Expect(err).NotTo(HaveOccurred())

		encoder := json.NewEncoder(tmpFile)
		err = encoder.Encode(requiredConfig)
		Expect(err).NotTo(HaveOccurred())

		err = tmpFile.Close()
		Expect(err).NotTo(HaveOccurred())

		os.Setenv("CONFIG", tmpFile.Name())
		fmt.Println(tmpFile.Name())
		config = cfg.LoadConfig()
	})

	AfterEach(func() {
		err := os.Remove(tmpFile.Name())
		Expect(err).NotTo(HaveOccurred())
	})

	It("should have the right defaults", func() {
		Expect(config.IncludeApps).To(BeTrue())
		Expect(config.DefaultTimeout).To(Equal(30 * time.Second))
		Expect(config.CfPushTimeout).To(Equal(2 * time.Minute))
		Expect(config.LongCurlTimeout).To(Equal(2 * time.Minute))
		Expect(config.BrokerStartTimeout).To(Equal(5 * time.Minute))
		Expect(config.AsyncServiceOperationTimeout).To(Equal(2 * time.Minute))

		// undocumented
		Expect(config.DetectTimeout).To(Equal(5 * time.Minute))
		Expect(config.SleepTimeout).To(Equal(30 * time.Second))
	})
})
