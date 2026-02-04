package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) ServerCmdAdd() {

	server := &cobra.Command{
		Use:     "server",
		Aliases: []string{"srv"},
		Short:   "Manage servers.",
		Long:    "Manage servers.",
	}

	serverList := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Show detail of all servers in your companies.",
		Long:    "Show detail of all servers in your companies.",
		Run:     cmd.runMiddleware.ServerList,
	}

	serverDetail := &cobra.Command{
		Use:     "show --server-uuid SERVER_UUID",
		Aliases: []string{"get"},
		Short:   "Show server detail.",
		Long:    "Show detailed information about a server.",
		Run:     cmd.runMiddleware.ServerDetail,
	}

	cmd.RootCommand.AddCommand(server)
	server.AddCommand(serverList, serverDetail)

	// Command arguments
	serverList.Flags().StringP("company-uuid", "c", "", "Set company UUID.")

	serverDetail.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = serverDetail.MarkFlagRequired("server-uuid")
}
