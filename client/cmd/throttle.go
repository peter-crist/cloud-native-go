package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/spf13/cobra"
)

type throttleFlags struct {
	attempts   int
	max        int
	refillRate int
	duration   int
}

var tFlags throttleFlags

func NewThrottleCmd(c pb.ChatClient) *cobra.Command {
	throttleCmd := &cobra.Command{
		Use:   "throttle",
		Short: "Demonstrates the Throttle pattern given the running gRPC server",
		Long: `The Throttle command calls the Throttle RPC which emulates a high rate of requests.
The command offers configuration for the max amount of requests, the token refill count,
and the duration to wait before refilling.

Default settings:
	-attempts: 100
	-max: 5
	-refill: 5
	-duration(s): 5

Example with default settings: 
	bin/client throttle
Example with 200 requests, a 3 tokens, and a 3 token per 10second refill rate: 
	bin/client timeout -a 200 -m 3 -r 3 -d 10
`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Calling Throttle with: %+v\n", tFlags)
			resp, err := c.DemoThrottle(
				context.Background(),
				&pb.ThrottleRequest{
					Attempts: int32(tFlags.attempts),
					Max:      int32(tFlags.max),
					Refill:   int32(tFlags.refillRate),
					Duration: int32(tFlags.duration),
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

	throttleCmd.Flags().IntVarP(&tFlags.attempts, "attempts", "a", 100, "Number of request attempts to make")
	throttleCmd.Flags().IntVarP(&tFlags.max, "max", "m", 5, "Number of request tokens in bucket")
	throttleCmd.Flags().IntVarP(&tFlags.refillRate, "refill", "r", 5, "Number of tokens to refill at a time")
	throttleCmd.Flags().IntVarP(&tFlags.duration, "duration", "d", 5, "Seconds to wait before refilling")

	return throttleCmd
}
