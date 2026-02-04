package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (cmd *CMD) VersionCmdAdd() {

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show API or CLI version.",
		Long:  "Show API or CLI version.",
	}

	versionCliCmd := &cobra.Command{
		Use:   "cli",
		Short: "Show CLI version.",
		Long:  "Show CLI version.",
		Run:   cmd.versionCli,
	}

	versionAPICmd := &cobra.Command{
		Use:   "api",
		Short: "Show API version.",
		Long:  "Show API version.",
		Run:   cmd.runMiddleware.VersionAPI,
	}

	cmd.RootCommand.AddCommand(versionCmd)
	versionCmd.AddCommand(versionAPICmd, versionCliCmd)
}

func (cmd *CMD) versionCli(cobraCommand *cobra.Command, args []string) {
	_ = args
	cmd.runMiddleware.ParseGlobalFlags(cobraCommand)

	versionStr := fmt.Sprintf("%d.%d.%d", cmd.VersionMajor, cmd.VersionMinor, cmd.VersionPatch)
	cmd.runMiddleware.PrintVersion("Titan CLI version", versionStr)
}
