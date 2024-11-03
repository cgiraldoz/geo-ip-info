package cli

import (
	"github.com/cgiraldoz/geo-ip-info/internal/interfaces"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gip",
	Short: "A CLI for querying IP address geolocation data.",
	Long:  "Geo IP Info is a CLI for querying IP address geolocation data.",
}

func InitializeCommands(redisCache interfaces.Cache, httpClient interfaces.Client) {
	rootCmd.AddCommand(NewStatsCmd(redisCache))
	rootCmd.AddCommand(NewApiCmd(redisCache, httpClient))
	rootCmd.AddCommand(NewIPCmd(redisCache, httpClient))
}

func Execute(redisCache interfaces.Cache, httpClient interfaces.Client) error {
	InitializeCommands(redisCache, httpClient)
	return rootCmd.Execute()
}
