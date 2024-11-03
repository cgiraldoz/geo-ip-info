package cli

import (
	"github.com/cgiraldoz/geo-ip-info/internal/interfaces"
	"github.com/cgiraldoz/geo-ip-info/internal/services"
	"github.com/spf13/cobra"
)

func NewIPCmd(redisCache interfaces.Cache, httpClient interfaces.Client) *cobra.Command {
	return &cobra.Command{
		Use:     "ip [ip address]",
		Short:   "Get information about an IP address",
		Long:    `Get information about an IP address, such as the country, currency, and timezone.`,
		Args:    cobra.ExactArgs(1),
		Example: "gip ip 8.8.8.8",
		Run: func(cmd *cobra.Command, args []string) {
			ip := args[0]
			ipDetails, err := services.GetIPLocationDetails(redisCache, httpClient, ip)

			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			cmd.Println("IP Location Details:")
			cmd.Printf("\nCountry: %s\n", ipDetails.CountryName)
			cmd.Printf("\nCountry Code: %s\n", ipDetails.Cca2)

			cmd.Println("\nCurrencies:")
			for code, currency := range ipDetails.Currencies {
				cmd.Printf("  - %s: %s (%s)\n", code, currency.Name, currency.Symbol)
			}

			cmd.Println("\nRelative Exchange Rates (compared to USD):")
			for code, rate := range ipDetails.RelativeRates {
				cmd.Printf("  - %s: %.4f\n", code, rate)
			}

			cmd.Println("\nCurrent Time by Timezone:")
			for timezone, currentTime := range ipDetails.CurrentTimeByTimezone {
				cmd.Printf("  - %s: %s\n", timezone, currentTime)
			}
		},
	}
}
