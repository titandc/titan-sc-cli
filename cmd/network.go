package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) NetworkCmdAdd() {

	network := &cobra.Command{
		Use:     "network",
		Aliases: []string{"net"},
		Short:   "Manage private networks.",
	}

	networkList := &cobra.Command{
		Use:     "list [--company-uuid COMPANY_UUID]",
		Aliases: []string{"ls"},
		Short:   "List all networks within your company.",
		Long:    "List all private networks created within your default company, another company UUID may be given.",
		Run:     cmd.runMiddleware.NetworkList,
	}

	networkDetail := &cobra.Command{
		Use:     "show --network-uuid NETWORK_UUID",
		Aliases: []string{"get"},
		Short:   "Show network detail.",
		Long:    "Show detailed information about a network.",
		Run:     cmd.runMiddleware.NetworkDetail,
	}

	networkCreate := &cobra.Command{
		Use:     "create --company-uuid COMPANY_UUID --name NETWORK_NAME --cidr CIDR_VALUE",
		Aliases: []string{"add"},
		Short:   "Create a new network.",
		Long:    "Create a new private network.",
		Run:     cmd.runMiddleware.NetworkCreate,
	}

	networkDelete := &cobra.Command{
		Use:     "delete --network-uuid NETWORK_UUID",
		Aliases: []string{"del"},
		Short:   "Delete a network.",
		Long:    "Completely delete a private network by UUID.",
		Run:     cmd.runMiddleware.NetworkRemove,
	}

	networkAttachServer := &cobra.Command{
		Use:   "attach --server-uuid SERVER_UUID --network-uuid NETWORK_UUID",
		Short: "Attach a server on private network.",
		Long:  "Attach a server on private network.",
		Run:   cmd.runMiddleware.NetworkAttachServer,
	}

	networkDetachServer := &cobra.Command{
		Use:   "detach --server-uuid SERVER_UUID --network-uuid NETWORK_UUID",
		Short: "Detach a server from private network.",
		Long:  "Detach a server from private network.",
		Run:   cmd.runMiddleware.NetworkDetachServer,
	}

	networkRename := &cobra.Command{
		Use:   "rename --name NEW_NAME --network-uuid NETWORK_UUID",
		Short: "Rename a network.",
		Long:  "Update the name of a private network, no space or special characters accepted.",
		Run:   cmd.runMiddleware.NetworkRename,
	}

	networkSetGW := &cobra.Command{
		Use:   "set-gw --ip IP_ADDRESS --network-uuid NETWORK_UUID",
		Short: "Set the gateway for a managed network.",
		Long:  "Set the gateway for a managed network.",
		Run:   cmd.runMiddleware.NetworkSetGateway,
	}

	networkUnsetGW := &cobra.Command{
		Use:   "unset-gw --network-uuid NETWORK_UUID",
		Short: "Unset the gateway of a managed network.",
		Long:  "Unset the gateway of a managed network.",
		Run:   cmd.runMiddleware.NetworkUnsetGateway,
	}

	cmd.RootCommand.AddCommand(network)
	network.AddCommand(networkList, networkDetail, networkCreate, networkDelete,
		networkAttachServer, networkDetachServer, networkRename,
		networkSetGW, networkUnsetGW)

	networkCreate.Flags().StringP("company-uuid", "", "", "Set company uuid.")
	networkCreate.Flags().StringP("name", "n", "", "Set new network name.")
	networkCreate.Flags().StringP("cidr", "c", "", "Provide a CIDR to enable managed network.")
	_ = networkCreate.MarkFlagRequired("network-uuid")
	_ = networkCreate.MarkFlagRequired("name")

	networkList.Flags().StringP("company-uuid", "c", "", "Set company UUID.")

	networkDelete.Flags().StringP("network-uuid", "", "", "Set network uuid.")
	_ = networkDelete.MarkFlagRequired("network-uuid")

	networkAttachServer.Flags().StringP("network-uuid", "", "", "Set network uuid.")
	networkAttachServer.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = networkAttachServer.MarkFlagRequired("server-uuid")
	_ = networkAttachServer.MarkFlagRequired("network-uuid")

	networkDetachServer.Flags().StringP("network-uuid", "", "", "Set network uuid.")
	networkDetachServer.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = networkDetachServer.MarkFlagRequired("server-uuid")
	_ = networkDetachServer.MarkFlagRequired("network-uuid")

	networkRename.Flags().StringP("network-uuid", "", "", "Set network uuid.")
	networkRename.Flags().StringP("name", "n", "", "Set new network name.")
	_ = networkRename.MarkFlagRequired("name")
	_ = networkRename.MarkFlagRequired("network-uuid")

	networkSetGW.Flags().StringP("network-uuid", "", "", "Set network uuid.")
	networkSetGW.Flags().StringP("ip", "i", "", "Set gateway IP.")
	_ = networkSetGW.MarkFlagRequired("ip")
	_ = networkSetGW.MarkFlagRequired("network-uuid")

	networkUnsetGW.Flags().StringP("network-uuid", "u", "", "Set network uuid.")
	_ = networkUnsetGW.MarkFlagRequired("network-uuid")

	networkDetail.Flags().StringP("network-uuid", "u", "", "Set network uuid.")
	_ = networkDetail.MarkFlagRequired("network-uuid")
}
