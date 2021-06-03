package applicationServices

import (
	"fmt"
	"github.com/google/uuid"
	"jwt-auth/auth/core/domainEntities"
	"jwt-auth/auth/core/domainServices"
)

// CreateUser creates a new user with a hashed password,
// returns error if username is already taken
func CreateUser(repository domainServices.IUserRepository, name string, password string) (domainEntities.User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return domainEntities.User{}, fmt.Errorf("could not hash password: %v\n", err)
	}

	user := domainEntities.User{
		Name:     name,
		Id:       uuid.NewString(),
		Password: hashedPassword,
	}

	err = repository.Save(user)
	if err != nil {
		return domainEntities.User{}, fmt.Errorf("could not create user: %v\n", err)
	}

	return user, nil
}
