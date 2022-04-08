package helpers

import (
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/onsi/ginkgo/v2/reporters"
)

type artifactsDirectoryConfig interface {
	GetArtifactsDirectory() string
}

func EnableCFTrace(config artifactsDirectoryConfig, componentName string) {
	_ = os.Setenv("CF_TRACE", traceLogFilePath(config, componentName))
}

func NewJUnitReporter(config artifactsDirectoryConfig, componentName string) *reporters.JUnitReporter {
	return reporters.NewJUnitReporter(jUnitReportFilePath(config, componentName))
}

func traceLogFilePath(config artifactsDirectoryConfig, componentName string) string {
	return filepath.Join(config.GetArtifactsDirectory(), fmt.Sprintf("CATS-TRACE-%s-%d.txt", sanitizeComponentName(componentName), ginkgoNode()))
}

func jUnitReportFilePath(config artifactsDirectoryConfig, componentName string) string {
	return filepath.Join(config.GetArtifactsDirectory(), fmt.Sprintf("junit-%s-%d.xml", sanitizeComponentName(componentName), ginkgoNode()))
}

func ginkgoNode() int {
	return ginkgo.GinkgoParallelProcess()
}

func sanitizeComponentName(componentName string) string {
	return strings.Replace(componentName, " ", "_", -1)
}
