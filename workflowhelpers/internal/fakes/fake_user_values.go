package fakes

type FakeUserValues struct {
	username string
	password string
	origin   string
}

func NewFakeUserValues(username, password, origin string) *FakeUserValues {
	return &FakeUserValues{
		username: username,
		password: password,
		origin: origin,
	}
}

func (user *FakeUserValues) Username() string {
	return user.username
}

func (user *FakeUserValues) Password() string {
	return user.password
}

func (user *FakeUserValues) Origin() string {
	return user.origin
}
