package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"jwt-auth/auth/core/domainEntities"
	"jwt-auth/auth/external/database/simpledb"
	"jwt-auth/auth/external/grpc/authgrpc"
	"jwt-auth/auth/external/rest"
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

	go authgrpc.HandleGrpc(grpcPort, &privateKey.PublicKey)
	rest.HandleHttp(restPort, privateKey, db)
}
