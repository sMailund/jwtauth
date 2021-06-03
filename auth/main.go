package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"jwt-auth/auth/core/domainEntities"
	"jwt-auth/auth/external/database/simpledb"
	"jwt-auth/auth/external/rest"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	db := simpledb.NewDb()

	testUser := domainEntities.User{
		Name:     "testuser",
		Id:       "123",
		Password: "dføasldfjaøsl",
	}

	err = db.Save(testUser)
	if err != nil {
		panic(err)
	}

	fmt.Printf("public key: %v\n", privateKey.Public())

	rest.HandleHttp(privateKey, db)
}


