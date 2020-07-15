package cmd

import (
    "github.com/spf13/cobra"
    . "titan-sc/api"
)

var server = &cobra.Command{
    Use: "server",
    Aliases: []string{"srv"},
    Short: "Manage servers.",
    Long: "Manage servers.",
}

var serverList = &cobra.Command{
    Use: "list",
    Aliases: []string{"ls"},
    Short: "Show detail of all servers in your companies.",
    Long: "Show detail of all servers in your companies.",
    Run: API.ServerList,
}

var serverDetail = &cobra.Command{
    Use: "show SERVER_UUID",
    Aliases: []string{"get"},
    Short: "Show server detail.",
    Long: "Show detailed information about a server.",
    Args: cmdNeed1UUID,
    Run: API.ServerDetail,
}

var serverStart = &cobra.Command{
    Use: "start SERVER_UUID",
    Short: "Send an action request to start a server.",
    Long: "Send an action request to start a server." +
        "\nList of available actions:" +
        "\n  start\n  stop\n  hardstop\n  reboot\n",
    Args: cmdNeed1UUID,
    Run: API.ServerStart,
}

var serverStop = &cobra.Command{
    Use: "stop SERVER_UUID",
    Short: "Send an action request to stop a server.",
    Long: "Send an action request to stop a server." +
        "\nList of available actions:" +
        "\n  start\n  stop\n  hardstop\n  reboot\n",
    Args: cmdNeed1UUID,
    Run: API.ServerStop,
}

var serverRestart = &cobra.Command{
    Use: "restart SERVER_UUID",
    Aliases: []string{"reboot"},
    Short: "Send an action request to restart a server.",
    Long: "Send an action request to restart a server." +
        "\nList of available actions:" +
        "\n  start\n  stop\n  hardstop\n  reboot\n",
    Args : cmdNeed1UUID,
    Run: API.ServerRestart,
}

var serverHardstop = &cobra.Command{
    Use: "hardstop SERVER_UUID",
    Short: "Send an action request to hardstop a server.",
    Long: "Send an action request to hardstop a server." +
        "\nList of available actions:" +
        "\n  start\n  stop\n  hardstop\n  reboot\n",
    Args : cmdNeed1UUID,
    Run: API.ServerHardstop,
}

var serverChangeName = &cobra.Command{
    Use: "rename --server-uuid SERVER_UUID --new-name NEW_NAME",
    Short: "Send a request to change server's name.",
    Long: "Send a request to change server's name.",
    Run: API.ServerChangeName,
}

var serverChangeReverse = &cobra.Command{
    Use: "reverse NEW_REVERSE",
    Short: "Send a request to change server's reverse.",
    Long: "Send a request to change server's reverse.",
    Run: API.ServerChangeReverse,
}

var serverLoadISO = &cobra.Command{
    Use: "load-iso --uri HTTPS_URI --server-uuid SERVER_UUID",
    Aliases: []string{"li"},
    Short: "Send a request to load an ISO from HTTPS.",
    Long: "Send a request to load a bootable ISO from HTTPS.",
    Run: API.ServerLoadISO,
}

var serverUnloadISO = &cobra.Command{
    Use: "unload-iso SERVER_UUID",
    Aliases: []string{"ui"},
    Short: "Send a request to unload previously loaded ISO(s).",
    Long: "Send a request to unload all previously loaded ISO(s).",
    Args : cmdNeed1UUID,
    Run: API.ServerUnloadISO,
}

func serverCmdAdd() {
    rootCmd.AddCommand(server)
    server.AddCommand(serverList, serverDetail, serverStart,
        serverStop, serverRestart, serverHardstop, serverLoadISO,
        serverUnloadISO, serverChangeName, serverChangeReverse)
    serverList.Flags().StringP("company-uuid", "c", "", "Set company UUID.")

    serverLoadISO.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
    serverLoadISO.Flags().StringP("uri", "u", "", "Set remote ISO URI (HTTPS only).")

    serverChangeName.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
    serverChangeName.Flags().StringP("name", "n", "", "Set new server's name.")

    serverChangeReverse.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
    serverChangeReverse.Flags().StringP("reverse", "r", "", "Set new server's reverse.")
}
