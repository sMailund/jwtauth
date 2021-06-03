package domainServices

import "jwt-auth/auth/core/domainEntities"

type IUserRepository interface {
	GetUserByName(name string) (domainEntities.User, error)
	Save(user domainEntities.User) error
}
