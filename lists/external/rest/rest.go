package rest

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net/http"
	"github.com/pascaldekloe/jwt"
)

var authKey *rsa.PublicKey

// testing jwt authentication
func testHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := []byte(r.Header.Get("Authorization")[5:]) // TODO improve header parsing
	claims, err := jwt.RSACheck(authHeader, authKey)
	if err != nil {
		http.Error(w, "failed to parse auth header", http.StatusForbidden)
		return
	}

	_, _ = fmt.Fprintf(w, "Hello %v (uid: %v), you are authenticated!\n", claims.Subject, claims.Set["uid"])
}

func HandleHttp(port string, key *rsa.PublicKey) {
	authKey = key
	http.HandleFunc("/test/", testHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}
