package grpcserver

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"github.com/sMailund/jwtauth/auth/external/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type server struct {
	*auth.UnimplementedAuthServer
}

var key *rsa.PublicKey

func (s server) GetPublicKey(ctx context.Context, request *auth.PublicKeyRequest) (*auth.PublicKeyResponse, error) {
	payload, err := json.Marshal(key)
	if err != nil {
		return &auth.PublicKeyResponse{}, status.Errorf(13, "internal server error: %v", err)
	}
	return &auth.PublicKeyResponse{PublicKey: payload}, nil
}

func HandleGrpc(port string, publicKey *rsa.PublicKey) {
	key = publicKey
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	auth.RegisterAuthServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
