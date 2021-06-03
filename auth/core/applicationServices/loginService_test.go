package applicationServices

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/pascaldekloe/jwt"
	"jwt-auth/auth/core/domainServices"
	"jwt-auth/auth/external/database/simpledb"
	"testing"
)

func setup() (*rsa.PrivateKey, domainServices.IUserRepository) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	repo := simpledb.NewDb()
	return key, repo
}

func TestLoginUserSuccessful(t *testing.T) {
	key, repo := setup()
	password := "testpassword"
	user1, _ := CreateUser(repo, "test1", password)

	got, err := LoginUser(repo, user1.Name, password, key)

	if err != nil {
		t.Fatalf("unexpected error: %v\n", err)
	}

	claims, err := jwt.RSACheck(got, &key.PublicKey)
	if err != nil {
		t.Fatalf("jwt signature error: %v\n", err)
	}

	if user1.Name != claims.Subject {
		t.Fatalf("expected subject with username %v, got %v\n", user1.Name, claims.Subject)
	}

	if user1.Id != claims.Set["uid"] {
		t.Fatalf("expected subject with uid %v, got %v\n", user1.Id, claims.Set["uid"])
	}
}

func TestLoginUserFailed(t *testing.T) {
	key, repo := setup()
	password := "testpassword"
	notpassword := "incorrect"
	user1, _ := CreateUser(repo, "test1", password)
	_, err := LoginUser(repo, user1.Name, notpassword, key)

	if err == nil {
		t.Fatalf("expected error on password mismatch %v != %v\n", password, notpassword)
	}
}

func TestLoginUserFailedIncorrectUid(t *testing.T) {
	key, repo := setup()
	password1 := "testpassword"
	password2 := "otherpassword"
	user1, _ := CreateUser(repo, "test1", password1)
	user2, _ := CreateUser(repo, "test2", password2)

	_, err := LoginUser(repo, user1.Name, password2, key)
	if err == nil {
		t.Fatalf("expected error on password1 mismatch %v != %v\n", password1, password2)
	}

	_, err = LoginUser(repo, user2.Name, password1, key)
	if err == nil {
		t.Fatalf("expected error on password1 mismatch %v != %v\n", password1, password2)
	}
}

