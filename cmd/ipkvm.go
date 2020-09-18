package cmd

import (
	"github.com/spf13/cobra"
	. "titan-sc/api"
)

var KVMIP = &cobra.Command{
	Use:     "kvmip",
	Aliases: []string{"kvm"},
	Short:   "Manage servers' KVM IP.",
	Long:    "Manage servers' KVM IP.",
}

var KVMIPStart = &cobra.Command{
	Use:   "start --server-uuid SERVER_UUID",
	Short: "Start a KVM IP.",
	Long:  "Start KVM IP on a server.",
	Run:   API.KVMIPStart,
}

var KVMIPStop = &cobra.Command{
	Use:   "stop --server-uuid SERVER_UUID",
	Short: "Stop a KVM IP.",
	Long:  "Stop KVM IP on a server.",
	Run:   API.KVMIPStop,
}

var KVMIPShow = &cobra.Command{
	Use:     "show --server-uuid SERVER_UUID",
	Aliases: []string{"get"},
	Short:   "Show KVM IP information.",
	Long:    "Show KVM IP information of a server.",
	Run:     API.KVMIPGetInfos,
}

func kvmIpCmdAdd() {
	rootCmd.AddCommand(KVMIP)
	KVMIP.AddCommand(KVMIPStart, KVMIPStop, KVMIPShow)

	KVMIPStart.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
	_ = KVMIPStart.MarkFlagRequired("server-uuid")

	KVMIPStop.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
	_ = KVMIPStop.MarkFlagRequired("server-uuid")

	KVMIPShow.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
	_ = KVMIPShow.MarkFlagRequired("server-uuid")

}
