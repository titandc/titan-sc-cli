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
    Use: "show server_uuid",
    Aliases: []string{"get"},
    Short: "Show server detail.",
    Long: "Show server detail.",
    Args: cmdNeed1UUID,
    Run: API.ServerDetail,
}

var serverStart = &cobra.Command{
    Use: "start server_uuid",
    Short: "Send a request for change state server action.",
    Long: "Send a request for change server state.\nAction list:" +
        "\n  start\n  stop\n  hardstop\n  reboot\n",
    Args: cmdNeed1UUID,
    Run: API.ServerStart,
}

var serverStop = &cobra.Command{
    Use: "stop server_uuid",
    Short: "Send a request for change state server action.",
    Long: "Send a request for change server state.\nAction list:" +
        "\n  start\n  stop\n  hardstop\n  reboot\n",
    Args: cmdNeed1UUID,
    Run: API.ServerStop,
}

var serverRestart = &cobra.Command{
    Use: "restart server_uuid",
    Aliases: []string{"reboot"},
    Short: "Send a request for change state server action.",
    Long: "Send a request for change server state.\nAction list:" +
        "\n  start\n  stop\n  hardstop\n  reboot\n",
    Args : cmdNeed1UUID,
    Run: API.ServerRestart,
}

var serverHardstop = &cobra.Command{
    Use: "hardstop server_uuid",
    Short: "Send a request for change state server action.",
    Long: "Send a request for change server state.\nAction list:" +
        "\n  start\n  stop\n  hardstop\n  reboot\n",
    Args : cmdNeed1UUID,
    Run: API.ServerHardstop,
}

var serverLoadISO = &cobra.Command{
    Use: "load-iso",
    Short: "Send a request to load an ISO from HTTPS.",
    Long: "Send a request to load a bootable ISO from HTTPS.",
    Run: API.ServerLoadISO,
}

var serverUnloadISO = &cobra.Command{
    Use: "unload-iso server_uuid",
    Short: "Send a request to unload previously loaded ISO(s).",
    Long: "Send a request to unload all previously loaded ISO(s)",
    Args : cmdNeed1UUID,
    Run: API.ServerUnloadISO,
}

func serverCmdAdd() {
    rootCmd.AddCommand(server)
    server.AddCommand(serverList, serverDetail, serverStart,
        serverStop, serverRestart, serverHardstop, serverLoadISO,
        serverUnloadISO)
    serverList.Flags().StringP("company-uuid", "u", "", "JSON ident output.")
    serverLoadISO.Flags().StringP("url", "u", "", "Load ISO by HTTPS.")
    serverLoadISO.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
}
