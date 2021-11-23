package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/spf13/cobra"
)

type FanOutFlags struct {
	SourceCount      int
	DestinationCount int
}

var foFlags FanOutFlags

func NewFanOutCmd(c pb.ChatClient) *cobra.Command {
	fanOutCmd := &cobra.Command{
		Use:   "fanout",
		Short: "Demonstrates the FanOut pattern given the running gRPC server",
		Long: `The FanOut command calls the FanOut RPC which will demonstrate splitting
values from a source channel to a configurable amount of destination channels. The source
channel counts up to a configurable value.

Default settings:
	-source_count: 10
	-destination_count: 5

Example with default settings: 
	bin/client fanout
Example with 10 destination channels splitting a 20 count source channel: 
	bin/client fanout -s 20 -d 10
`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Calling FanOut with: %+v\n", foFlags)
			resp, err := c.DemoFanOut(
				context.Background(),
				&pb.FanOutRequest{
					SourceCount:      int32(foFlags.SourceCount),
					DestinationCount: int32(foFlags.DestinationCount),
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

	fanOutCmd.Flags().IntVarP(&foFlags.SourceCount, "source_count", "s", 10, "Source channel count limit")
	fanOutCmd.Flags().IntVarP(&foFlags.DestinationCount, "destination_count", "d", 5, "Number of destination channels")

	return fanOutCmd
}
