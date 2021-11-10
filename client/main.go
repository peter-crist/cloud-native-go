package main

import (
	"fmt"
	"log"

	"github.com/peter-crist/cloud-native-go/client/cmd"
	pb "github.com/peter-crist/cloud-native-go/proto"
	"google.golang.org/grpc"
)

var (
	port string = "5300"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%s", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := pb.NewChatClient(conn)
	sendCmd := cmd.NewSendCmd(c)
	cbCmd := cmd.NewCircuitBreakerCmd(c)
	debounceCmd := cmd.NewDebounceCmd(c)
	rootCmd := cmd.InitRootCmd(
		sendCmd,
		cbCmd,
		debounceCmd,
	)

	rootCmd.Execute()
}
