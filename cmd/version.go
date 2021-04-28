package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
	_ = cmd
	_ = args

	fmt.Printf("Titan cloud CLI version %d.%d.%d\n",
		cmd.VersionMajor, cmd.VersionMinor, cmd.VersionPatch)
	os.Exit(0)
}
