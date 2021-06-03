package applicationServices

import (
	"crypto/rsa"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"golang.org/x/crypto/bcrypt"
	"jwt-auth/auth/core/authErrors"
	"jwt-auth/auth/core/domainServices"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// LoginUser authenticates user with password, and returns a byteslice representing a signed JWT token for the user
func LoginUser(repository domainServices.IUserRepository, username string, password string, key *rsa.PrivateKey) ([]byte, error) {
	// TODO: authN user, check password
	user, err := repository.GetUserByName(username)
	if err != nil {
		return nil, authErrors.NewNoSuchUserError(username)
	}

	hash := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return nil, authErrors.IncorrectPassword{Name: username}
	}

	c := jwt.Claims{}
	c.Issuer = "auth server"
	c.Subject = username
	c.Set = make(map[string]interface{})
	c.Set["id"] = user.Id

	token, err := c.RSASign(jwt.RS256, key)
	if err != nil {
		return nil, fmt.Errorf("could not sign: %v", err)
	}

	return token, nil
}
