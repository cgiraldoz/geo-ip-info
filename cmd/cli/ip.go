package cli

import (
	"context"
	"fmt"
	"github.com/cgiraldoz/geo-ip-info/internal/interfaces"
	"github.com/cgiraldoz/geo-ip-info/internal/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewIPCmd(redisCache interfaces.Cache, httpClient interfaces.Client) *cobra.Command {
	return &cobra.Command{
		Use:     "ip [ip address]",
		Short:   "Get information about an IP address",
		Long:    `Get information about an IP address, such as the country, currency, and timezone.`,
		Args:    cobra.ExactArgs(1),
		Example: "gip ip 8.8.8.8",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(redisCache)
			IpLocation, err := services.NewIPLocation(httpClient)
			if err != nil {
				fmt.Println("Error creating IP location service:", err)
				return
			}

			contextTimeout := viper.GetDuration("context.timeout")
			ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
			defer cancel()

			ip := args[0]
			info, err := IpLocation.GetIPLocation(ctx, ip)
			if err != nil {
				fmt.Println("Error getting IP location:", err)
				return
			}

			fmt.Printf("Country: %s\n", info.Name)
		},
	}
}
