package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	cb "github.com/peter-crist/cloud-native-go/circuitbreaker"
	"github.com/peter-crist/cloud-native-go/debounce"
	pb "github.com/peter-crist/cloud-native-go/proto"
	"github.com/peter-crist/cloud-native-go/retry"
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
func (s *server) DemoCircuitBreaker(
	ctx context.Context,
	req *pb.CircuitBreakerRequest,
) (
	*pb.CircuitBreakerResponse,
	error,
) {
	var (
		resp string
		err  error
	)
	conn := cb.Breaker(slowConnection, uint(req.GetFailureThreshold()))
	for i := 0; i < int(req.GetAttempts()); i++ {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(req.Timeout))
		defer cancel()
		resp, err = conn(ctxWithTimeout)
		log.Printf("⏱  Waiting 0.5s before trying to connect again.\n\n")
		time.Sleep(time.Millisecond * 500) //pause to simulate slower connection attempts to showcase resetting the breaker
	}

	log.Printf("🥳 %d connection attempts complete 🥳\n", req.GetAttempts())
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to dependency after %d attempts", req.GetAttempts())
	}

	return &pb.CircuitBreakerResponse{Message: resp}, nil
}

func (s *server) DemoDebounce(
	ctx context.Context,
	req *pb.DebounceRequest,
) (
	*pb.DebounceResponse,
	error,
) {
	var resp string
	attempts := req.GetAttempts()

	log.Println(req.GetDuration())
	conn := debounce.DebounceFirst(
		func(ctx context.Context) (string, error) {
			return "success", nil
		},
		time.Second*time.Duration(req.GetDuration()),
	)

	log.Printf("💻 Spamming %d connection attempts\n", attempts)
	for i := 0; i < int(attempts); i++ {
		resp, _ = conn(ctx)
		log.Printf("⏱ Waiting %dms before a new connection attempt...\n", req.GetDelay())
		time.Sleep(time.Millisecond * time.Duration(req.GetDelay()))
	}

	log.Printf("🥳 %d/%d connection attempts complete 🥳\n", attempts, attempts)
	return &pb.DebounceResponse{Message: resp}, nil
}

func (s *server) DemoRetry(
	ctx context.Context,
	req *pb.RetryRequest,
) (
	*pb.RetryResponse,
	error,
) {
	conn := retry.Retry(
		emulateTransientError,
		int(req.GetCount()),
		time.Millisecond*time.Duration(req.GetDelay()),
	)
	resp, err := conn(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.RetryResponse{Message: resp}, nil
}

func emulateTransientError(ctx context.Context) (string, error) {
	//randomly return error
	rand.Seed(time.Now().UnixNano())
	prob := 3

	isError := rand.Intn(10) > prob //roughly 1 in 3 calls will return an error
	if isError {
		return "", errors.New("❌ FAILED ❌")
	}

	return "✅ SUCCESS ✅", nil
}

func slowConnection(ctx context.Context) (string, error) {
	rand.Seed(time.Now().UnixNano())
	duration := rand.Intn(10)
	log.Printf("Simulating a long connection attempt for %d seconds", duration)
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
