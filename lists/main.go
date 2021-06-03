package main

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/sMailund/jwtauth/auth/external/grpc/auth"
	"github.com/sMailund/jwtauth/lists/external/rest"
	"google.golang.org/grpc"
	"log"
)

const authAddress = "localhost:8888"
const restPort = ":8081"

func getAuthPublicKey(client auth.AuthClient) (*rsa.PublicKey, error) {
	response, err := client.GetPublicKey(context.Background(), &auth.PublicKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("could not get public key of auth server: %v\n", err)
	}

	var key rsa.PublicKey
	err = json.Unmarshal(response.PublicKey, &key)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key from auth server: %v\n", err)
	}

	return &key, nil
}

func main() {
	fmt.Println("Waiting for connection to auth server...")
	conn, err := grpc.Dial(authAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("Getting public key from auth server...")
	authClient := auth.NewAuthClient(conn)
	authKey, err := getAuthPublicKey(authClient)
	if err != nil {
		panic(err)
	}

	_ = conn.Close()
	fmt.Println("Done.")

	rest.HandleHttp(restPort, authKey)
}
