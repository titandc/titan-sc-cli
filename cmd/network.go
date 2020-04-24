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
    Use: "list",
    Aliases: []string{"ls"},
    Short: "List company network.",
    Long: "List company network details.",
    Run: API.NetworkList,
}

var networkCreate = &cobra.Command{
    Use: "create company-uuid",
    Short: "Create a new private network.",
    Long: "Create a new private network, need company UUID.",
    Args: cmdNeed1UUID,
    Run: API.NetworkCreate,
}

var networkRemove = &cobra.Command{
    Use: "remove network_uuid",
    Aliases: []string{"rm"},
    Short: "Remove one private network.",
    Long: "Remove a private network, need network UUID.",
    Args : cmdNeed1UUID,
    Run: API.NetworkRemove,
}

var networkAttachServer = &cobra.Command{
    Use: "attach [--network-uuid --server-uuid]",
    Aliases: []string{"add"},
    Short: "Attach server on private network.",
    Long: "Attach server on private network need companyUUID serverUUID.",
    Run: API.NetworkAttachServer,
}

var networkDetachServer = &cobra.Command{
    Use: "detach [--network-uuid --server-uuid]",
    Short: "Detach server on private network.",
    Long: "Detach server on private network need companyUUID serverUUID.",
    Run: API.NetworkDetachServer,
}

var networkRename = &cobra.Command{
    Use: "rename [--network-uuid --name].",
    Short: "Rename private network.",
    Long: "Update the name of an existing network. You must be an administrator of the company. No space accepted",
    Run: API.NetworkRename,
}

func networkCmdAdd() {
    rootCmd.AddCommand(network)
    network.AddCommand(networkList, networkCreate, networkRemove,
        networkAttachServer, networkDetachServer, networkRename)

    networkCreate.Flags().StringP("name", "n", "", "Set new network name.")
    networkList.Flags().StringP("company-uuid", "u", "", "Set company UUID.")

    networkAttachServer.Flags().StringP("network-uuid", "u", "", "Set network UUID.")
    networkAttachServer.Flags().StringP("server-uuid", "s", "", "Set server UUID.")

    networkDetachServer.Flags().StringP("network-uuid", "u", "", "Set network UUID.")
    networkDetachServer.Flags().StringP("server-uuid", "s", "", "Set server UUID.")

    networkRename.Flags().StringP("network-uuid", "u", "", "Set network UUID.")
    networkRename.Flags().StringP("name", "n", "", "Set network new name.")
}
