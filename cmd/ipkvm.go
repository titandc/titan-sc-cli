package cmd

import (
    "github.com/spf13/cobra"
    . "titan-sc/api"
)

var KVMIP = &cobra.Command{
    Use: "kvmip",
    Aliases: []string{"kvm"},
    Short: "Manage servers' KVM IP.",
    Long: "Manage servers' KVM IP.",
}

var KVMIPStart = &cobra.Command{
    Use: "start SERVER_UUID",
    Short: "Start a KVM IP.",
    Long: "Start KVM IP on a server.",
    Args: cmdNeed1UUID,
    Run: API.KVMIPStart,
}

var KVMIPStop = &cobra.Command{
    Use: "stop SERVER_UUID",
    Short: "Stop a KVM IP.",
    Long: "Stop KVM IP on a server.",
    Args: cmdNeed1UUID,
    Run: API.KVMIPStop,
}

var KVMIPShow = &cobra.Command{
    Use: "show SERVER_UUID",
    Aliases: []string{"get"},
    Short: "Show KVM IP information.",
    Long: "Show KVM IP information of a server.",
    Args: cmdNeed1UUID,
    Run: API.KVMIPGetInfos,
}

func kvmIpCmdAdd() {
    rootCmd.AddCommand(KVMIP)
    KVMIP.AddCommand(KVMIPStart)
    KVMIP.AddCommand(KVMIPStop)
    KVMIP.AddCommand(KVMIPShow)
}
