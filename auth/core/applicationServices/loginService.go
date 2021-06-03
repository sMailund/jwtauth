package applicationServices

import (
	"crypto/rsa"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"github.com/sMailund/jwtauth/auth/core/authErrors"
	"github.com/sMailund/jwtauth/auth/core/domainServices"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// LoginUser authenticates user with password, and returns a byteslice representing a signed JWT token for the user
func LoginUser(repository domainServices.IUserRepository, username string, password string, key *rsa.PrivateKey) ([]byte, error) {
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
	c.Set["uid"] = user.Id

	token, err := c.RSASign(jwt.RS256, key)
	if err != nil {
		return nil, fmt.Errorf("could not sign: %v", err)
	}

	return token, nil
}
