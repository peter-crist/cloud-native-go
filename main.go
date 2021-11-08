package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	cb "github.com/peter-crist/cloud-native-go/circuitbreaker"
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

//Send takes a message, performs a pseudo-slow operation, and returns the SHA of the message input
func (s *server) Send(
	ctx context.Context,
	req *pb.SendRequest,
) (
	resp *pb.SendResponse,
	err error,
) {
	slowConnection(ctx)

	msg := req.GetMessage()
	log.Printf("Received: %v", msg)

	bv := []byte(msg)
	hasher := sha1.New()
	hasher.Write(bv)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return &pb.SendResponse{
		Message: fmt.Sprintf("Message sent: %s", msg),
		Sha:     sha,
	}, nil
}

//CircuitBreaker takes in a request which specifies the desired failureThreshold and demonstrates
//the CircuitBreaker pattern with a pseudorandomally failing dependency
func (s *server) CircuitBreaker(
	ctx context.Context,
	req *pb.CircuitBreakerRequest,
) (
	resp *pb.CircuitBreakerResponse,
	err error,
) {
	var result string
	conn := cb.Breaker(slowConnection, uint(req.GetFailureThreshold()))
	for i := 0; i < int(req.GetAttempts()); i++ {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(req.Timeout))
		defer cancel()
		result, err = conn(ctxWithTimeout)
		log.Println("⏱ Waiting 0.5s before trying to connect again.")
		time.Sleep(time.Millisecond * 500) //pause to simulate slower connection attempts to showcase resetting the breaker
	}

	log.Printf("🥳 %d connection attempts complete 🥳", req.GetAttempts())
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to dependency after %d attempts", req.GetAttempts())
	}

	return &pb.CircuitBreakerResponse{Message: result}, nil
}

func slowConnection(ctx context.Context) (string, error) {
	duration := rand.Intn(10)
	log.Printf("Simulating a long connection attempt for %d milliseconds", duration)
	for i := 0; i < duration; i++ {
		select {
		case <-ctx.Done():
			log.Println("Failed to connect in time...")
			return "", ctx.Err()
		default:
			log.Printf("%ds elapsed...", i+1)
			time.Sleep(time.Second * 1)
		}
	}
	success := "Connection complete!"
	log.Println(success)
	return success, nil
}

//Heres what going to happen:
// user enters cb -f 3 -a 100 -t 10
