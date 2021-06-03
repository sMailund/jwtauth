package authErrors

import "fmt"

type IncorrectPassword struct {
	Name string
}

func (i IncorrectPassword) Error() string {
	return fmt.Sprintf("incorrect password for user %v\n", i.Name)
}
