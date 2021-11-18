package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/spf13/cobra"
)

type timeoutFlags struct {
	duration int
	timeout  int
}

var toFlags timeoutFlags

func NewTimeoutCmd(c pb.ChatClient) *cobra.Command {
	timeoutCmd := &cobra.Command{
		Use:   "timeout",
		Short: "Demonstrates the Timeout pattern given the running gRPC server",
		Long: `The Timeout command calls the Timeout RPC which emulates a potentially
long-running function that not only doesn’t accept a Context value, but hypothetically
comes from a package you don’t control. This timeout pattern allows us to implement a proper timeout.

Default settings:
	-duration(s): 5
	-timeout(s): 2

Example with default settings: 
	bin/client timeout
Example with a slow function that takes about 10s and there is a 5s timeout: 
	bin/client timeout -d 5 -t 2
`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Calling Timeout with: %+v\n", toFlags)
			resp, err := c.DemoTimeout(
				context.Background(),
				&pb.TimeoutRequest{
					Duration: int32(toFlags.duration),
					Timeout:  int32(toFlags.timeout),
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

	timeoutCmd.Flags().IntVarP(&toFlags.duration, "duration", "d", 5, "Duration of the emulated 'slow' function in seconds")
	timeoutCmd.Flags().IntVarP(&toFlags.timeout, "timeout", "t", 2, "Timeout in seconds")

	return timeoutCmd
}
