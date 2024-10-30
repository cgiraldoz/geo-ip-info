package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:     "gip [ip]",
		Short:   "A CLI for querying IP address geolocation data.",
		Long:    "Geo IP Info is a CLI for querying IP address geolocation data.",
		Example: "gip 0.0.0.0\ngip stats",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				fmt.Printf("Consulting geolocation for IP: %s\n", args[0])
			} else {
				fmt.Println("Please provide a valid IP address or command.")
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(NewStatsCmd())
	rootCmd.AddCommand(NewApiCmd())
}

func Execute() error {
	return rootCmd.Execute()
}
