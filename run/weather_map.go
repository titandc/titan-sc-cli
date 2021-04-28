package run

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (run *RunMiddleware) WeatherMap(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	weatherMap, err := run.API.GetWeatherMap()
	if err != nil {
		run.OutputError(err)
		return
	}
	if !run.HumanReadable {
		printAsJson(weatherMap)
	} else {
		fmt.Printf("Titan Weather Map:\n"+
			"  Compute: %s\n"+
			"  Storage: %s\n"+
			"  Public network: %s\n"+
			"  Private network: %s\n",
			weatherMap.Compute, weatherMap.Storage,
			weatherMap.PublicNetwork, weatherMap.PrivateNetwork)
	}
}
