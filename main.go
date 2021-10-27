package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net"

	pb "github.com/peter-crist/cloud-native-go/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	port string = "5300"
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	pb.RegisterChatServer(s, &server{})

	log.Printf("gRPC server listening on port %s...\n", port)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct{}

func (s *server) Send(
	ctx context.Context,
	req *pb.Request,
) (
	resp *pb.Response,
	err error,
) {
	msg := req.GetMessage()
	log.Printf("Received: %v", msg)

	bv := []byte(msg)
	hasher := sha1.New()
	hasher.Write(bv)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return &pb.Response{
		Message: fmt.Sprintf("Message sent: %s", msg),
		Sha:     sha,
	}, nil
}
