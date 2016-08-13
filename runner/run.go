package runner

import "github.com/onsi/gomega/gexec"

func Run(executable string, args ...string) *gexec.Session {
	cmdStarter := NewCommandStarter()
	return cmdStarter.Start(executable, args...)
}
