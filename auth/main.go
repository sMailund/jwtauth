package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/sMailund/jwtauth/auth/core/domainEntities"
	"github.com/sMailund/jwtauth/auth/external/database/simpledb"
	"github.com/sMailund/jwtauth/auth/external/grpc/grpcserver"
	"github.com/sMailund/jwtauth/auth/external/rest"
)

const grpcPort = ":8888"
const restPort = ":8080"

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

	go grpcserver.HandleGrpc(grpcPort, &privateKey.PublicKey)
	rest.HandleHttp(restPort, privateKey, db)
}
