package domainServices

import "github.com/sMailund/jwtauth/auth/core/domainEntities"

type IUserRepository interface {
	GetUserByName(name string) (domainEntities.User, error)
	Save(user domainEntities.User) error
}
