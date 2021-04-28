package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) WeatherMapCmdAdd() {

	weatherMap := &cobra.Command{
		Use:   "weathermap",
		Short: "Show weather map.",
		Long:  "Show weathermap.",
		Run:   cmd.runMiddleware.WeatherMap,
	}

	cmd.RootCommand.AddCommand(weatherMap)
}
