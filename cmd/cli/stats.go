package cli

import "github.com/spf13/cobra"

func NewStatsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "View usage distance statistic of the service",
		Long:  `Display distance statistics for service usage.`,
		Run: func(cmd *cobra.Command, args []string) {
			println("Stats command executed")
		},
	}
}
