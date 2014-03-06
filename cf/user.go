package cf

type UserContext struct {
	Username string
	Password string
	Org      string
	Space    string
}

var NewUserContext = func(username string, password string, org string, space string) UserContext {
	uc := UserContext{}
	uc.Username = username
	uc.Password = password
	uc.Org = org
	uc.Space = space
	return uc
}
