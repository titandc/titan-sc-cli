package cmd

import (
    "github.com/spf13/cobra"
    . "titan-sc/api"
)

var network = &cobra.Command{
    Use: "network",
    Aliases: []string{"net"},
    Short: "Manage private networks.",
}

var networkList = &cobra.Command{
    Use: "list [--company-uuid COMPANY_UUID]",
    Aliases: []string{"ls"},
    Short: "List all networks within your company.",
    Long: "List all private networks created within your default company, another company UUID may be given.",
    Run: API.NetworkList,
}

var networkCreate = &cobra.Command{
    Use: "create NETWORK_UUID",
    Short: "Create a new network.",
    Long: "Create a new private network.",
    Args: cmdNeed1UUID,
    Run: API.NetworkCreate,
}

var networkDelete = &cobra.Command{
    Use: "delete NETWORK_UUID",
    Aliases: []string{"del"},
    Short: "Delete a network.",
    Long: "Completely delete a private network by UUID.",
    Args : cmdNeed1UUID,
    Run: API.NetworkRemove,
}

var networkAttachServer = &cobra.Command{
    Use: "attach --network-uuid NETWORK_UUID --server-uuid SERVER_UUID",
    Aliases: []string{"add"},
    Short: "Attach a server on private network.",
    Long: "Attach a server on private network.",
    Run: API.NetworkAttachServer,
}

var networkDetachServer = &cobra.Command{
    Use: "detach --network-uuid NETWORK_UUID --server-uuid SERVER_UUID",
    Short: "Detach a server from private network.",
    Long: "Detach a server from private network.",
    Run: API.NetworkDetachServer,
}

var networkRename = &cobra.Command{
    Use: "rename --network-uuid NETWORK_UUID --name NEW_NAME",
    Short: "Rename a network.",
    Long: "Update the name of a private network, no space or special characters accepted.",
    Run: API.NetworkRename,
}

func networkCmdAdd() {
    rootCmd.AddCommand(network)
    network.AddCommand(networkList, networkCreate, networkDelete,
        networkAttachServer, networkDetachServer, networkRename)

    networkCreate.Flags().StringP("name", "n", "", "Set new network name.")
    networkList.Flags().StringP("company-uuid", "c", "", "Set company UUID.")

    networkAttachServer.Flags().StringP("network-uuid", "n", "", "Set network UUID.")
    networkAttachServer.Flags().StringP("server-uuid", "s", "", "Set server UUID.")

    networkDetachServer.Flags().StringP("network-uuid", "n", "", "Set network UUID.")
    networkDetachServer.Flags().StringP("server-uuid", "s", "", "Set server UUID.")

    networkRename.Flags().StringP("network-uuid", "u", "", "Set network UUID.")
    networkRename.Flags().StringP("name", "n", "", "Set new network name.")
}
