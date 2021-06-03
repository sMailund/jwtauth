package rest

import (
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"jwt-auth/auth/core/applicationServices"
	"jwt-auth/auth/core/authErrors"
	"jwt-auth/auth/core/domainServices"
	"log"
	"net/http"
	"time"
)

var privateKey *rsa.PrivateKey
var repo domainServices.IUserRepository

type loginForm struct {
	username string
	password string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not parse POST body", http.StatusBadRequest)
		return
	}

	var login loginForm
	err = json.Unmarshal(body, &login)
	if err != nil {
		http.Error(w, "could not parse POST body", http.StatusBadRequest)
		return
	}

	token, err := applicationServices.LoginUser(repo, login.username, login.password, privateKey)
	if err != nil {
		if authErrors.IsNotFoundError(err) {
			http.Error(w, "could not find user", 404)
		} else {
			http.Error(w, "unknown server error", 500)
		}

		return
	}

	authCookie := http.Cookie{
		Name:       "Authorization",
		Value:      "JWT: " + string(token),
		Expires:    time.Time{},
		Secure:     true,
		HttpOnly:   true,
	}

	http.SetCookie(w, &authCookie)
}

func HandleHttp(key *rsa.PrivateKey, userRepo domainServices.IUserRepository) {
	privateKey = key
	repo = userRepo
	http.HandleFunc("/login/", loginHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
