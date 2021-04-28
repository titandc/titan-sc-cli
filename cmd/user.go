package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) UserCmdAdd() {

	apiUser := &cobra.Command{
		Use:   "user",
		Short: "Manage your user information.",
		Long:  "Manage your user informations.",
	}

	userAllInfos := &cobra.Command{
		Use:   "infos",
		Short: "Get all user informations.",
		Long:  "Get all user informations.",
		Run:   cmd.runMiddleware.UserShowAllInfos,
	}

	cmd.RootCommand.AddCommand(apiUser)
	apiUser.AddCommand(userAllInfos)
}
