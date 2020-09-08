package cmd

import (
	"github.com/spf13/cobra"
	. "titan-sc/api"
)

var network = &cobra.Command{
	Use:     "network",
	Aliases: []string{"net"},
	Short:   "Manage private networks.",
}

var networkList = &cobra.Command{
	Use:     "list [--company-uuid COMPANY_UUID]",
	Aliases: []string{"ls"},
	Short:   "List all networks within your company.",
	Long:    "List all private networks created within your default company, another company UUID may be given.",
	Run:     API.NetworkList,
}

var networkCreate = &cobra.Command{
	Use:     "create COMPANY_UUID",
	Aliases: []string{"add"},
	Short:   "Create a new network.",
	Long:    "Create a new private network.",
	Args:    cmdNeed1UUID,
	Run:     API.NetworkCreate,
}

var networkDelete = &cobra.Command{
	Use:     "delete NETWORK_UUID",
	Aliases: []string{"del"},
	Short:   "Delete a network.",
	Long:    "Completely delete a private network by UUID.",
	Args:    cmdNeed1UUID,
	Run:     API.NetworkRemove,
}

var networkAttachServer = &cobra.Command{
	Use:   "attach --server-uuid SERVER_UUID NETWORK_UUID",
	Short: "Attach a server on private network.",
	Long:  "Attach a server on private network.",
	Args:  cmdNeed1UUID,
	Run:   API.NetworkAttachServer,
}

var networkDetachServer = &cobra.Command{
	Use:   "detach --server-uuid SERVER_UUID NETWORK_UUID",
	Short: "Detach a server from private network.",
	Long:  "Detach a server from private network.",
	Args:  cmdNeed1UUID,
	Run:   API.NetworkDetachServer,
}

var networkRename = &cobra.Command{
	Use:   "rename --name NEW_NAME NETWORK_UUID",
	Short: "Rename a network.",
	Long:  "Update the name of a private network, no space or special characters accepted.",
	Args:  cmdNeed1UUID,
	Run:   API.NetworkRename,
}

var networkSetGW = &cobra.Command{
	Use:   "set-gw --ip IP_ADDRESS NETWORK_UUID",
	Short: "Set the gateway for a managed network.",
	Long:  "Set the gateway for a managed network.",
	Args:  cmdNeed1UUID,
	Run:   API.NetworkSetGateway,
}

var networkUnsetGW = &cobra.Command{
	Use:   "unset-gw NETWORK_UUID",
	Short: "Unset the gateway of a managed network.",
	Long:  "Unset the gateway of a managed network.",
	Args:  cmdNeed1UUID,
	Run:   API.NetworkUnsetGateway,
}

func networkCmdAdd() {
	rootCmd.AddCommand(network)
	network.AddCommand(networkList, networkCreate, networkDelete,
		networkAttachServer, networkDetachServer, networkRename,
		networkSetGW, networkUnsetGW)

	networkCreate.Flags().StringP("name", "n", "", "Set new network name.")
	networkCreate.Flags().StringP("cidr", "c", "", "Provide a CIDR to enable managed network.")

	networkList.Flags().StringP("company-uuid", "c", "", "Set company UUID.")

	networkAttachServer.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = networkAttachServer.MarkFlagRequired("server-uuid")

	networkDetachServer.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = networkDetachServer.MarkFlagRequired("server-uuid")

	networkRename.Flags().StringP("name", "n", "", "Set new network name.")
	_ = networkRename.MarkFlagRequired("name")

	networkSetGW.Flags().StringP("ip", "i", "", "Set gateway IP.")
	_ = networkSetGW.MarkFlagRequired("ip")
}
