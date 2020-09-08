package cmd

import (
	"github.com/spf13/cobra"
	. "titan-sc/api"
)

var apiUser = &cobra.Command{
	Use:   "user",
	Short: "Manage your user information.",
	Long:  "Manage your user informations.",
}

var userAllInfos = &cobra.Command{
	Use:   "infos",
	Short: "Get all user informations.",
	Long:  "Get all user informations.",
	Run:   API.UserShowAllInfos,
}

func userCmdAdd() {
	rootCmd.AddCommand(apiUser)
	apiUser.AddCommand(userAllInfos)
}
