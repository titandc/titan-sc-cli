package cmd

import (
    "github.com/spf13/cobra"
    . "titan-sc/api"
)

var server = &cobra.Command{
    Use: "server",
    Aliases: []string{"srv"},
    Short: "Administrate servers.",
    Long: "Administrate servers.",
}

var serverList = &cobra.Command{
    Use: "list",
    Aliases: []string{"ls"},
    Short: "Show servers detail on all your compagnies.",
    Long: "Show servers detail on all your compagnies.",
    Run: API.ServerList,
}

var serverDetail = &cobra.Command{
    Use: "show server_uuid",
    Aliases: []string{"get"},
    Short: "Show servers detail on one server.",
    Long: "Show servers detail on one server.",
    Args: cmdNeed1UUID,
    Run: API.ServerDetail,
}

var serverStart = &cobra.Command{
    Use: "start server_uuid",
    Short: "Show detail on one server.",
    Long: "Show detail on one server.",
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

func serverCmdAdd() {
    rootCmd.AddCommand(server)
    server.AddCommand(serverList, serverDetail, serverStart,
        serverStop, serverRestart, serverHardstop)
    serverList.PersistentFlags().StringP("company-uuid", "u", "", "Json ident output.")
}
