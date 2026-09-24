package commandstarter_test

import (
	"bytes"
	"io"
	"os/exec"
	"sync"
	"time"

	"github.com/cloudfoundry/cf-test-helpers/v2/commandstarter"

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

func (f *fakeReporter) Report(t time.Time, cmd *exec.Cmd) {
	f.calledWith.time = t
	f.calledWith.cmd = cmd
}

type completionReporter struct {
	mu           sync.Mutex
	completed    bool
	completedCmd []string
	elapsed      time.Duration
	exitCode     int
}

func (f *completionReporter) Report(t time.Time, cmd *exec.Cmd) {}

func (f *completionReporter) ReportCompletion(cmd *exec.Cmd, elapsed time.Duration, exitCode int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.completed = true
	f.completedCmd = cmd.Args
	f.elapsed = elapsed
	f.exitCode = exitCode
}

func (f *completionReporter) snapshot() (bool, []string, time.Duration, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.completed, f.completedCmd, f.elapsed, f.exitCode
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

	When("the reporter implements CompletionReporter", func() {
		It("reports completion with exit code and non-negative elapsed, without the caller calling Wait()", func() {
			cr := &completionReporter{}
			session, err := cmdStarter.Start(cr, "bash", "-c", "exit 3")
			Expect(err).To(Succeed())

			// Intentionally do NOT call session.Wait(); the library's own goroutine
			// must observe completion via session.Exited.
			Eventually(func() bool {
				done, _, _, _ := cr.snapshot()
				return done
			}, 5*time.Second).Should(BeTrue())

			_, args, elapsed, exit := cr.snapshot()
			Expect(args).To(Equal([]string{"bash", "-c", "exit 3"}))
			Expect(exit).To(Equal(3))
			Expect(elapsed).To(BeNumerically(">=", 0))

			_ = session
		})
	})
})
