package workflowhelpers

import (
	"time"
)

type cfUser interface {
	SetCfHomeDir() (string, string)
	UnsetCfHomeDir(string, string)
	Login(time.Duration)
	Logout(time.Duration)
	TargetSpace(time.Duration)
}

var AsUser = func(user cfUser, timeout time.Duration, actions func()) {
	originalCfHomeDir, currentCfHomeDir := user.SetCfHomeDir()
	user.Login(timeout)
	defer func() {
		user.Logout(timeout)
		user.UnsetCfHomeDir(originalCfHomeDir, currentCfHomeDir)
	}()

	user.TargetSpace(timeout)
	actions()
}
