package cli

import (
	"github.com/cgiraldoz/geo-ip-info/cmd/api"
	"github.com/spf13/cobra"
)

func NewApiCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "api",
		Short:   "Start the Geo IP Info API server.",
		Long:    "Start the Geo IP Info API server to query IP address geolocation data.",
		Example: "gip api",
		Run: func(cmd *cobra.Command, args []string) {
			api.StartAPI()
		},
	}
}
