package rest

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/sMailund/jwtauth/auth/core/applicationServices"
	"github.com/sMailund/jwtauth/auth/core/authErrors"
	"github.com/sMailund/jwtauth/auth/core/domainEntities"
	"github.com/sMailund/jwtauth/auth/core/domainServices"
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

type UserDao struct {
	Username string
	Uid      string
}

func mapToUserDao(user domainEntities.User) ([]byte, error) {
	dao := UserDao{
		Username: user.Name,
		Uid:      user.Id,
	}

	return json.Marshal(dao)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // TODO fix duplicate code in create user
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "incorrect post body format", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" {
		http.Error(w, "missing username", http.StatusBadRequest)
		return
	}
	if password == "" {
		http.Error(w, "missing password", http.StatusBadRequest)
		return
	}

	token, err := applicationServices.LoginUser(repo, username, password, privateKey)
	if err != nil { // TODO improve error handling
		switch err.(type) {
		case authErrors.NoSuchUser:
			http.Error(w, "could not find user", 404)
		case authErrors.IncorrectPassword:
			http.Error(w, "incorrect password", 401)
		default:
			http.Error(w, "unknown server error", 500)
		}

		return
	}

	expire := time.Now().Add(20 * time.Minute)

	authCookie := http.Cookie{
		Name:     "Authorization",
		Value:    "JWT: " + string(token),
		Expires:  expire,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &authCookie)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "incorrect post body format", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" {
		http.Error(w, "missing username", http.StatusBadRequest)
		return
	}
	if password == "" {
		http.Error(w, "missing password", http.StatusBadRequest)
		return
	}

	created, err := applicationServices.CreateUser(repo, username, password)
	dao, err := mapToUserDao(created)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(dao)
}

func HandleHttp(port string, key *rsa.PrivateKey, userRepo domainServices.IUserRepository) {
	privateKey = key
	repo = userRepo
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/create/", createHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}
