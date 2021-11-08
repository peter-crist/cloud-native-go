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
		Short: "short",
		Long:  `long`,
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

	circuitBreakerCmd.Flags().IntVarP(&cbFlags.failureThreshold, "threshold", "f", 3, "Number of failed attempts to trigger breaker")
	circuitBreakerCmd.Flags().IntVarP(&cbFlags.attempts, "attempts", "a", 10, "Number of connection attempts")
	circuitBreakerCmd.Flags().IntVarP(&cbFlags.timeout, "timeout", "t", 100, "Duration in milliseconds for connection timeout")

	return circuitBreakerCmd
}
