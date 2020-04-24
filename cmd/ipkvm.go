package cmd

import (
    "github.com/spf13/cobra"
    . "titan-sc/api"
)

var IPKvm = &cobra.Command{
    Use: "ipkvm",
    Short: "Manage server's IP KVM.",
    Long: "Manage server's IP KVM.",
}

var IPKvmStart = &cobra.Command{
    Use: "start server_uuid",
    Short: "Start ip-kvm.",
    Long: "Start ip-kvm on server UUID.",
    Args: cmdNeed1UUID,
    Run: API.IPKvmStart,
}

var IPKvmStop = &cobra.Command{
    Use: "stop server_uuid",
    Short: "Stop ip-kvm.",
    Long: "Stop ip-kvm on server UUID.",
    Args: cmdNeed1UUID,
    Run: API.IPKvmStop,
}

var IPKvmShow = &cobra.Command{
    Use: "show server_uuid",
    Aliases: []string{"get"},
    Short: "Show IP Kvm informations.",
    Long: "Show IP Kvm information (status,URI).",
    Args: cmdNeed1UUID,
    Run: API.IPKvmGetInfos,
}


func ipkvmCmdAdd() {
    rootCmd.AddCommand(IPKvm)
    IPKvm.AddCommand(IPKvmStart)
    IPKvm.AddCommand(IPKvmStop)
    IPKvm.AddCommand(IPKvmShow)
}
