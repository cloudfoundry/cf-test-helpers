package helpers

import (
	"github.com/onsi/gomega/gexec"
)

var SkipSSLValidation bool

func Curl(args ...string) *gexec.Session {
	return CurlSkipSSL(SkipSSLValidation, args...)
}

func CurlSkipSSL(skip bool, args ...string) *gexec.Session {
	curlArgs := append([]string{"-s"}, args...)
	if skip {
		curlArgs = append([]string{"-k"}, curlArgs...)
	}
	return Run("curl", curlArgs...)
}
