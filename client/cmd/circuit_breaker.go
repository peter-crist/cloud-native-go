package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/spf13/cobra"
)

type CircuitBreakerFlags struct {
	failureThreshold int
	attempts         int
	timeout          int
}

var cbFlags CircuitBreakerFlags

func NewCircuitBreakerCmd(c pb.ChatClient) *cobra.Command {
	circuitBreakerCmd := &cobra.Command{
		Use:   "cb",
		Short: "Demonstrates the Circuit Breaker pattern given the running gRPC server",
		Long: `The Circuit Breaker command calls the CircuitBreaker RPC which emulates connections
to an external dependency that either time out or fail. The command offers configuration for the
number of consecutive failures before the breaker flips, the raw number of connection attempts to
be made, and the "timeout" for simulating failures.

Default settings:
	-threshold: 2
	-attempts: 10
	-timeout: 3

Example with default settings: 
	bin/client cb
Example with a failure threshold of 3, 5 total attempts, and a 10 second timeout: 
	bin/client cb -f 3 -a 5 -t 10
		`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Calling CircuitBreaker with: %+v\n", cbFlags)
			resp, err := c.CircuitBreaker(
				context.Background(),
				&pb.CircuitBreakerRequest{
					FailureThreshold: int32(cbFlags.failureThreshold),
					Attempts:         int32(cbFlags.attempts),
					Timeout:          int32(cbFlags.timeout),
				},
			)
			if err != nil {
				log.Println(fmt.Errorf("Failed to send message to gRPC server: %w", err))
				return
			}

			result := "\n\t%s"
			log.Printf(result, resp.GetMessage())
		},
	}

	circuitBreakerCmd.Flags().IntVarP(&cbFlags.failureThreshold, "threshold", "f", 2, "Number of failed attempts to trigger breaker")
	circuitBreakerCmd.Flags().IntVarP(&cbFlags.attempts, "attempts", "a", 10, "Number of connection attempts")
	circuitBreakerCmd.Flags().IntVarP(&cbFlags.timeout, "timeout", "t", 3, "Duration in seconds for connection timeout")

	return circuitBreakerCmd
}
