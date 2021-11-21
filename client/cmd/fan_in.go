package cmd

import (
	"context"
	"fmt"
	"log"

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/spf13/cobra"
)

type FanInFlags struct {
	SourceCount int
}

var fiFlags FanInFlags

func NewFanInCmd(c pb.ChatClient) *cobra.Command {
	fanInCmd := &cobra.Command{
		Use:   "fanin",
		Short: "Demonstrates the FanIn pattern given the running gRPC server",
		Long: `The FanIn command calls the FanIn RPC which will demonstrate funneling
a configurable amount of source channels to one destination channel. The individual channels
simply count to a random value up to 10.

Default settings:
	-source_count: 5

Example with default settings: 
	bin/client fanin
Example with 10 source channels to fan in: 
	bin/client fanin -s 10
`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Calling FanIn with: %+v\n", fiFlags)
			resp, err := c.DemoFanIn(
				context.Background(),
				&pb.FanInRequest{
					SourceCount: int32(fiFlags.SourceCount),
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

	fanInCmd.Flags().IntVarP(&fiFlags.SourceCount, "source_count", "s", 5, "Number of source channels to funnel")

	return fanInCmd
}
