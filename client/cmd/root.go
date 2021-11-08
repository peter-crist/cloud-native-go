package cmd

import (
	"github.com/spf13/cobra"
)

type RunFunc func(cmd *cobra.Command, args []string)

var rootCmd = &cobra.Command{
	Use:   "cloud-native-go",
	Short: "Cloud native patterns implemented in Go",
	Long: `
A gRPC server/client that can be used to showcase various distributed computing patterns.
For example, circuit breakers or debounce`,
}

func InitRootCmd(subcommands ...*cobra.Command) *cobra.Command {
	for _, cmd := range subcommands {
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}
