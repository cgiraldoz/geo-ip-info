package cli

import (
	"context"
	"github.com/cgiraldoz/geo-ip-info/internal/interfaces"
	"github.com/cgiraldoz/geo-ip-info/internal/services"
	"github.com/spf13/cobra"
)

func NewStatsCmd(redisCache interfaces.Cache) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "View usage distance statistics of the service",
		Long:  `Display distance statistics for service usage, including farthest, closest, and average distances from Buenos Aires.`,
		Run: func(cmd *cobra.Command, args []string) {

			stats, err := services.GetDistanceStatsFromCache(context.Background(), redisCache)
			if err != nil {
				cmd.PrintErrln("Error retrieving stats:", err)
				return
			}

			cmd.Println("Distance Statistics:")
			cmd.Printf("  Farthest Distance: %.2f km (Country: %s)\n", stats.FarthestDistance, stats.FarthestCountryName)
			cmd.Printf("  Closest Distance: %.2f km (Country: %s)\n", stats.ClosestDistance, stats.ClosestCountryName)
			cmd.Printf("  Total Distance: %.2f km\n", stats.TotalDistance)
			cmd.Printf("  Total Requests: %d\n", stats.TotalRequests)

			averageDistance := services.CalculateWeightedAverageDistance(stats)
			cmd.Printf("  Average Distance: %.2f km\n", averageDistance)

			cmd.Println("\nDistance by Country:")
			for country, data := range stats.CountryDistances {
				cmd.Printf("  - %s: %.2f km (Requests: %d)\n", country, data.TotalDistance/float64(data.Requests), data.Requests)
			}
		},
	}
}
