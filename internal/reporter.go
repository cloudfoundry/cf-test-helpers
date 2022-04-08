package internal

import (
	"os/exec"
	"time"
)

type Reporter interface {
	Report(bool, time.Time, *exec.Cmd)
}
