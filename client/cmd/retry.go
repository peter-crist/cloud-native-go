package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/spf13/cobra"
)

type RetryFlags struct {
	retries int
	delay   int
}

var rFlags RetryFlags

func NewRetryCmd(c pb.ChatClient) *cobra.Command {
	retryCmd := &cobra.Command{
		Use:   "retry",
		Short: "Demonstrates the Retry pattern given the running gRPC server",
		Long: `The Retry command calls the Retry RPC which emulates transient failures.
The command offers configuration for the number of retries to attempt, and the delay
between retrying.

Default settings:
	-retries: 5
	-delay(ms): 500

Example with default settings: 
	bin/client retry
Example with a 3 retries, and a 200ms request delay: 
	bin/client retry -r 3 -d 200
`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Calling Retry with: %+v\n", rFlags)
			resp, err := c.DemoRetry(
				context.Background(),
				&pb.RetryRequest{
					Count: int32(rFlags.retries),
					Delay: int32(rFlags.delay),
				},
			)
			if err != nil {
				log.Println(fmt.Errorf("Failed to send message to gRPC server: %w", err))
				return
			}

			result := "%s"
			log.Printf(result, resp.GetMessage())
		},
	}

	retryCmd.Flags().IntVarP(&rFlags.retries, "retries", "r", 5, "Number of retry attempts")
	retryCmd.Flags().IntVarP(&rFlags.delay, "delay", "d", 500, "Delay in milliseconds before next attempt")

	return retryCmd
}
