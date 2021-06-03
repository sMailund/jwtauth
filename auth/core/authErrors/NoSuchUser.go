package authErrors

import "fmt"

type NoSuchUser struct {
	name string
}

func NewNoSuchUserError(name string) NoSuchUser {
	e := new(NoSuchUser)
	e.name = name
	return *e
}

func (u NoSuchUser) Error() string {
	return fmt.Sprintf("no user with username %v", u.name)
}
