package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (cmd *CMD) KvmIpCmdAdd() {

	KVMIP := &cobra.Command{
		Use:     "kvmip",
		Aliases: []string{"kvm"},
		Short:   "Manage servers' KVM IP.",
		Long:    "Manage servers' KVM IP.",
		GroupID: "resources",
	}

	KVMIPStart := &cobra.Command{
		Use:   "start --server-oid SERVER_OID",
		Short: "Start a KVM IP.",
		Long:  "Start KVM IP on a server.",
		Run:   cmd.runMiddleware.KVMIPStart,
	}

	KVMIPStop := &cobra.Command{
		Use:   "stop --server-oid SERVER_OID",
		Short: "Stop a KVM IP.",
		Long:  "Stop KVM IP on a server.",
		Run:   cmd.runMiddleware.KVMIPStop,
	}

	KVMIPShow := &cobra.Command{
		Use:     "show [--server-oid SERVER_OID | --kvm-oid KVM_OID]",
		Aliases: []string{"get"},
		Short:   "Show KVM IP information.",
		Long: `Show KVM IP information.

You can query by either:
  - Server OID (-s): Get KVM info for a server (recommended)
  - KVM OID (-k): Get KVM info directly by its OID

Exactly one of --server-oid or --kvm-oid must be provided.`,
		Example: `  titan-sc kvmip show -s 604a19c439430d34d52028be
  titan-sc kvmip show --kvm-oid 698348a9bf4df844b0fc85e8`,
		Run: cmd.runMiddleware.KVMIPGetInfos,
		PreRunE: func(c *cobra.Command, args []string) error {
			serverOID, _ := c.Flags().GetString("server-oid")
			kvmOID, _ := c.Flags().GetString("kvm-oid")
			if serverOID == "" && kvmOID == "" {
				return fmt.Errorf("one of --server-oid (-s) or --kvm-oid (-k) is required")
			}
			return nil
		},
	}

	cmd.RootCommand.AddCommand(KVMIP)
	KVMIP.AddCommand(KVMIPStart, KVMIPStop, KVMIPShow)

	KVMIPStart.Flags().StringP("server-oid", "s", "", "Server OID to start KVM on.")
	_ = KVMIPStart.MarkFlagRequired("server-oid")

	KVMIPStop.Flags().StringP("server-oid", "s", "", "Server OID to stop KVM on.")
	_ = KVMIPStop.MarkFlagRequired("server-oid")

	KVMIPShow.Flags().StringP("server-oid", "s", "", "Server OID to get KVM info for.")
	KVMIPShow.Flags().StringP("kvm-oid", "k", "", "KVM session OID to get info for.")
	KVMIPShow.MarkFlagsMutuallyExclusive("server-oid", "kvm-oid")
}
