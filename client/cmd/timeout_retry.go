package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/spf13/cobra"
)

type timeoutAndRetryFlags struct {
	duration int
	timeout  int
}

var trFlags timeoutAndRetryFlags

func NewTimeoutAndRetryCmd(c pb.ChatClient) *cobra.Command {
	timeoutCmd := &cobra.Command{
		Use:   "tr",
		Short: "Demonstrates the Timeout in combination with a Retry pattern given the running gRPC server",
		Long: `The Timeout command calls the Timeout RPC which emulates a potentially
long-running function that not only doesn’t accept a Context value, but hypothetically
comes from a package you don’t control. Retry allows us to retry a timed out connection.

Default settings:
	-retries: 3
	-delay(ms): 200
	-duration(s): 5
	-timeout(s): 2

Example with default settings: 
	bin/client tr
Example with a slow function that takes about 10s and there is a 5s timeout: 
	bin/client tr -d 10 -t 2
`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Calling TimeoutAndRetry with: %+v\n", trFlags)
			resp, err := c.DemoTimeout(
				context.Background(),
				&pb.TimeoutRequest{
					Duration: int32(trFlags.duration),
					Timeout:  int32(trFlags.timeout),
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

	timeoutCmd.Flags().IntVarP(&trFlags.duration, "retries", "r", 3, "Number of retries")
	timeoutCmd.Flags().IntVarP(&trFlags.timeout, "delay", "s", 2, "Delay between retries in milliseconds")
	timeoutCmd.Flags().IntVarP(&toFlags.duration, "duration", "d", 5, "Duration of the emulated 'slow' function in seconds")
	timeoutCmd.Flags().IntVarP(&toFlags.timeout, "timeout", "t", 2, "Timeout in seconds")

	return timeoutCmd
}
