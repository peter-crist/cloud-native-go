package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/spf13/cobra"
)

type DebounceFlags struct {
	duration int
	attempts int
	delay    int
}

var debFlags DebounceFlags

func NewDebounceCmd(c pb.ChatClient) *cobra.Command {
	debounceCmd := &cobra.Command{
		Use:   "debounce",
		Short: "Demonstrates the Debounce pattern given the running gRPC server",
		Long: `The Debounce command calls the Debounce RPC which emulates connections
to an external dependency. The command offers configuration for the duration between
request clusters, the raw number of connection attempts to be made, and the delay
between subsequent requests.

Default settings:
	-duration(s): 3
	-attempts: 50
	-delay(ms): 200

Example with default settings: 
	bin/client debounce
Example with a duration of 3s, 50 total attempts, and a 200ms request delay: 
	bin/client debounce -t 3 -a 50 -d 200
`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Calling Debounce with: %+v\n", debFlags)
			resp, err := c.Debounce(
				context.Background(),
				&pb.DebounceRequest{
					Duration: int32(debFlags.duration),
					Attempts: int32(debFlags.attempts),
					Delay:    int32(debFlags.delay),
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

	debounceCmd.Flags().IntVarP(&debFlags.duration, "duration", "t", 3, "Duration of request clusters")
	debounceCmd.Flags().IntVarP(&debFlags.attempts, "attempts", "a", 50, "Number of connection attempts")
	debounceCmd.Flags().IntVarP(&debFlags.delay, "delay", "d", 200, "Delay in milliseconds before next connection attempt")

	return debounceCmd
}
