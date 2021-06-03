package applicationServices

import (
	"crypto/rsa"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"jwt-auth/auth/core/authErrors"
	"jwt-auth/auth/core/domainEntities"
	"jwt-auth/auth/core/domainServices"
)

// LoginUser authenticates user with password, and returns a byteslice representing a signed JWT token for the user
func LoginUser(repository domainServices.IUserRepository, username string, password string, key *rsa.PrivateKey) ([]byte, error){
	// TODO: authN user
	user, err := repository.GetUserByName(username)
	if err != nil {
		return nil, authErrors.NewNoSuchUserError(username)
	}

	c := jwt.Claims{}
	c.Issuer = "auth server"
	c.Subject = "user"
	c.Set = make(map[string]interface{})
	c.Set["id"] = user.Id

	token, err := c.RSASign(jwt.RS256, key)
	if err != nil {
		return nil, fmt.Errorf("could not sign: %v", err)
	}

	return token, nil
}

func createUser(name string, password string) (domainEntities.User, error) {
	return domainEntities.User{}, fmt.Errorf("stub, not implemented")// TODO
}
