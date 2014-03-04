package cf

import (
	. "github.com/vito/cmdtest/matchers"
	. "github.com/onsi/gomega"
)

func CfAsUser(username string, password string, actions func() error) error {
	defer func() {
		Expect(Cf("logout")).To(ExitWith(0))
	}()

	if Expect(Cf("login", username, password)).To(ExitWith(0)) {
		return actions()
	}

	return nil
}
