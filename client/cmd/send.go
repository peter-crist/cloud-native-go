package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	pb "github.com/peter-crist/cloud-native-go/proto"

	"github.com/spf13/cobra"
)

func NewSendCmd(c pb.ChatClient) *cobra.Command {
	return &cobra.Command{
		Use:   "send",
		Short: "Simple request/response command to validate communication with the gRPC server is established",
		Long: `Sends an input message to the gRPC server which simulates a variable processing time and returns
SHA of the received message.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("send command did not receive the correct amount of arguments (1)")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := c.Send(context.Background(), &pb.SendRequest{Message: args[0]})
			if err != nil {
				log.Println(fmt.Errorf("Failed to send message to gRPC server: %w", err))
				return
			}

			result := "\n\t%s\n\tSha: %s"
			log.Printf(result, resp.GetMessage(), resp.GetSha())
		},
	}
}
