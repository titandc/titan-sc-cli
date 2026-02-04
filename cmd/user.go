package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) UserCmdAdd() {

	apiUser := &cobra.Command{
		Use:     "user",
		Short:   "Manage your user information.",
		Long:    "Manage your user informations.",
		GroupID: "resources",
	}

	userInfo := &cobra.Command{
		Use:   "info",
		Short: "Get user information.",
		Long:  "Get user information.",
		Run:   cmd.runMiddleware.UserInfo,
	}

	cmd.RootCommand.AddCommand(apiUser)
	apiUser.AddCommand(userInfo)
}
