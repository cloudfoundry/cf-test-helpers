package commandstarter_test

import (
	"bytes"
	"io"
	"os/exec"
	"time"

	"github.com/cloudfoundry-incubator/cf-test-helpers/v2/commandstarter"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

type fakeReporter struct {
	calledWith struct {
		time time.Time
		cmd  *exec.Cmd
	}
}

func (f *fakeReporter) Report(b bool, t time.Time, cmd *exec.Cmd) {
	f.calledWith.time = t
	f.calledWith.cmd = cmd
}

var _ = Describe("CommandStarter", func() {
	var cmdStarter *commandstarter.CommandStarter
	var reporter *fakeReporter

	BeforeEach(func() {
		cmdStarter = commandstarter.NewCommandStarter()
		reporter = &fakeReporter{}
	})

	It("reports the command that it's running", func() {
		session, err := cmdStarter.Start(reporter, "bash", "-c", "echo \"hello world\"")
		Expect(err).To(Succeed())
		Expect(reporter.calledWith.cmd.Args).To(Equal([]string{"bash", "-c", "echo \"hello world\""}))
		Eventually(session).Should(Say("hello world"))
	})

	When("created with stdin", func() {
		var stdin io.Reader

		BeforeEach(func() {
			stdin = bytes.NewBufferString("name from input")
			cmdStarter = commandstarter.NewCommandStarterWithStdin(stdin)
		})

		It("reports what command is running and sends the input to the command", func() {
			session, err := cmdStarter.Start(reporter, "bash", "-c", `echo "hello $(cat -)"`)
			Expect(err).To(Succeed())
			Expect(reporter.calledWith.cmd.Args).To(Equal([]string{"bash", "-c", `echo "hello $(cat -)"`}))
			Eventually(session).Should(Say("hello name from input"))
		})
	})
})
