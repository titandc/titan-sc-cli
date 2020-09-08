package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	. "titan-sc/api"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show API or CLI version.",
	Long:  "Show API or CLI version.",
}

var versionCliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Show CLI version.",
	Long:  "Show CLI version.",
	Run:   versionCli,
}

var versionAPICmd = &cobra.Command{
	Use:   "api",
	Short: "Show API version.",
	Long:  "Show API version.",
	Run:   API.VersionAPI,
}

func versionCmdAdd() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.AddCommand(versionAPICmd, versionCliCmd)
}

func versionCli(cmd *cobra.Command, args []string) {
	_ = cmd
	_ = args

	fmt.Printf("Titan cloud CLI version %d.%d.%d\n",
		VersionMajor, VersionMinor, VersionPatch)
	os.Exit(0)
}
