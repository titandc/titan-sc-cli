package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) NetworkCmdAdd() {

	network := &cobra.Command{
		Use:     "network",
		Aliases: []string{"net"},
		Short:   "Manage private networks.",
		GroupID: "resources",
	}

	networkList := &cobra.Command{
		Use:     "list [--company-oid COMPANY_OID]",
		Aliases: []string{"ls"},
		Short:   "List all networks within your company.",
		Long: `List all private networks created within your company.

If --company-oid is not specified, your default company will be used.`,
		Run: cmd.runMiddleware.NetworkList,
	}

	networkDetail := &cobra.Command{
		Use:     "show --network-oid NETWORK_UUID",
		Aliases: []string{"get"},
		Short:   "Show network detail.",
		Long:    "Show detailed information about a network.",
		Run:     cmd.runMiddleware.NetworkDetail,
	}

	networkCreate := &cobra.Command{
		Use:     "create --name NETWORK_NAME",
		Aliases: []string{"add"},
		Short:   "Create a new network.",
		Long:    "Create a new private network.",
		Run:     cmd.runMiddleware.NetworkCreate,
	}

	networkDelete := &cobra.Command{
		Use:     "delete --network-oid NETWORK_UUID",
		Aliases: []string{"del"},
		Short:   "Delete a network.",
		Long:    "Completely delete a private network by OID.",
		Run:     cmd.runMiddleware.NetworkRemove,
	}

	networkAttachServer := &cobra.Command{
		Use:   "attach --server-oid SERVER_OID --network-oid NETWORK_UUID",
		Short: "Attach a server on private network.",
		Long:  "Attach a server on private network.",
		Run:   cmd.runMiddleware.NetworkAttachServer,
	}

	networkDetachServer := &cobra.Command{
		Use:   "detach --server-oid SERVER_OID --network-oid NETWORK_UUID",
		Short: "Detach a server from private network.",
		Long:  "Detach a server from private network.",
		Run:   cmd.runMiddleware.NetworkDetachServer,
	}

	networkRename := &cobra.Command{
		Use:   "rename --name NEW_NAME --network-oid NETWORK_UUID",
		Short: "Rename a network.",
		Long:  "Update the name of a private network, no space or special characters accepted.",
		Run:   cmd.runMiddleware.NetworkRename,
	}

	// DRP commands for network
	networkDrp := &cobra.Command{
		Use:   "drp",
		Short: "Manage network DRP (Disaster Recovery Plan).",
		Long:  "Manage network DRP (Disaster Recovery Plan).\nEnable or disable DRP replication for a private network.",
	}

	networkDrpEnable := &cobra.Command{
		Use:   "enable --network-oid NETWORK_OID",
		Short: "Enable DRP for a network.",
		Long:  "Enable DRP (Disaster Recovery Plan) replication for a private network.\nThis will replicate the network configuration to the target site.",
		Run:   cmd.runMiddleware.NetworkDrpEnable,
	}

	networkDrpDisable := &cobra.Command{
		Use:   "disable --network-oid NETWORK_OID --yes-i-understand-network-will-be-unavailable",
		Short: "Disable DRP for a network.",
		Long: `Disable DRP replication for a private network.

Disabling DRP will stop replication to the target site.
Note: If servers fail over to the secondary site, they will not have access
to this network until DRP is re-enabled or the network is manually recreated.

This operation requires the --yes-i-understand-network-will-be-unavailable flag to confirm.`,
		Run: cmd.runMiddleware.NetworkDrpDisable,
	}

	cmd.RootCommand.AddCommand(network)
	network.AddCommand(networkList, networkDetail, networkCreate, networkDelete,
		networkAttachServer, networkDetachServer, networkRename, networkDrp)

	// DRP subcommands
	networkDrp.AddCommand(networkDrpEnable, networkDrpDisable)

	networkCreate.Flags().StringP("name", "n", "", "Set new network name.")
	_ = networkCreate.MarkFlagRequired("name")

	networkList.Flags().StringP("company-oid", "c", "", "Company OID (uses your default company if not specified).")

	networkDelete.Flags().StringP("network-oid", "", "", "Set network OID.")
	_ = networkDelete.MarkFlagRequired("network-oid")

	networkAttachServer.Flags().StringP("network-oid", "", "", "Set network OID.")
	networkAttachServer.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = networkAttachServer.MarkFlagRequired("server-oid")
	_ = networkAttachServer.MarkFlagRequired("network-oid")

	networkDetachServer.Flags().StringP("network-oid", "", "", "Set network OID.")
	networkDetachServer.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = networkDetachServer.MarkFlagRequired("server-oid")
	_ = networkDetachServer.MarkFlagRequired("network-oid")

	networkRename.Flags().StringP("network-oid", "", "", "Set network OID.")
	networkRename.Flags().StringP("name", "n", "", "Set new network name.")
	_ = networkRename.MarkFlagRequired("name")
	_ = networkRename.MarkFlagRequired("network-oid")

	networkDetail.Flags().StringP("network-oid", "", "", "Set network OID.")
	_ = networkDetail.MarkFlagRequired("network-oid")

	// Network DRP flags
	networkDrpEnable.Flags().StringP("network-oid", "", "", "Network OID to enable DRP for.")
	_ = networkDrpEnable.MarkFlagRequired("network-oid")

	networkDrpDisable.Flags().StringP("network-oid", "", "", "Network OID to disable DRP for.")
	_ = networkDrpDisable.MarkFlagRequired("network-oid")
	networkDrpDisable.Flags().BoolP("yes-i-understand-network-will-be-unavailable", "", false, "Confirm that you understand this will break private network connectivity on the recovery site.")
}
