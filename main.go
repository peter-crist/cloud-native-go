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

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/peter-crist/cloud-native-go/circuitbreaker"
	"github.com/peter-crist/cloud-native-go/debounce"
	"github.com/peter-crist/cloud-native-go/fanin"
	"github.com/peter-crist/cloud-native-go/retry"
	"github.com/peter-crist/cloud-native-go/throttle"
	"github.com/peter-crist/cloud-native-go/timeout"

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
	conn := circuitbreaker.Breaker(slowConnection, uint(req.GetFailureThreshold()))
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

func (s *server) DemoThrottle(
	ctx context.Context,
	req *pb.ThrottleRequest,
) (
	*pb.ThrottleResponse,
	error,
) {
	var resp string
	conn := throttle.Throttle(
		func(ctx context.Context) (string, error) {
			return "success", nil
		},
		uint(req.GetMax()),
		uint(req.GetRefill()),
		time.Second*time.Duration(req.GetDuration()),
	)

	attempts := req.GetAttempts()
	log.Printf("💻 Spamming %d attempts\n", attempts)
	for i := 0; i < int(attempts); i++ {
		resp, _ = conn(ctx)
		log.Printf("⏱ Waiting %dms before a new attempt...\n", 200)
		time.Sleep(time.Millisecond * 200)
	}

	log.Printf("🥳 %d/%d connection attempts complete 🥳\n", attempts, attempts)
	return &pb.ThrottleResponse{Message: resp}, nil
}

func (s *server) DemoTimeout(
	ctx context.Context,
	req *pb.TimeoutRequest,
) (
	*pb.TimeoutResponse,
	error,
) {
	var resp string
	conn := timeout.Timeout(
		slowFunction(time.Second * time.Duration(req.Duration)),
	)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Second*time.Duration(req.Timeout))
	defer cancel()
	resp, err := conn(ctxWithTimeout, "input1")
	if err != nil {
		return nil, err
	}
	return &pb.TimeoutResponse{Message: resp}, nil
}

func (s *server) DemoFanIn(
	ctx context.Context,
	req *pb.FanInRequest,
) (
	*pb.FanInResponse,
	error,
) {
	// Race condition here. It is finishing somehow before all values are added to the sources.
	// Need to inspect fanin function.

	var resp string
	sources := make([]<-chan int, 0) // Create an empty channel slice

	for i := 1; i <= int(req.SourceCount); i++ {
		ch := make(chan int)
		sources = append(sources, ch) // Create a channel; add to sources

		go func(channel_id int) { // Run a toy goroutine for each
			defer close(ch) // Close ch when the routine ends

			rand.Seed(time.Now().UnixNano())
			count := rand.Intn(9) + 1
			log.Printf("Source #%d set to count up to %d", channel_id, count)
			for j := 1; j <= count; j++ {
				// Each sent value will correspond to its specific channel indicated by the initial digits
				// Ex. 34 is the 3rd source with and 4th value
				val := channel_id*10 + j
				log.Printf("➕ Adding %d to source channel #%d ➕", val, channel_id)
				ch <- val
				time.Sleep(time.Second)
			}
			log.Printf("✅ All %d values pushed to channel #%d. Closing channel... ✅", count, channel_id)
		}(i)
	}

	dest := fanin.Funnel(sources...)
	for d := range dest {
		log.Printf("📖 Reading %d off destination channel 📖", d)
	}

	log.Printf("🥳 All values successfully fanned-in 🥳")
	return &pb.FanInResponse{Message: resp}, nil
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

func slowFunction(d time.Duration) timeout.SlowFunction {
	return func(s string) (string, error) {
		log.Printf("Received argument: %v", s)
		log.Printf("Emulating slow function for %s", d)
		time.Sleep(d)
		return "✅ Slow Function completed ✅", nil
	}
}
