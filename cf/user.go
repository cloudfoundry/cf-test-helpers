package cf

type User struct {
	Username string
	Password string
	Org      string
	Space    string
}

var NewUser = func(username string, password string, org string, space string) User {
	u := User{}
	u.Username = username
	u.Password = password
	u.Org = org
	u.Space = space
	return u
}
