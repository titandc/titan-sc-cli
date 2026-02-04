package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) ServerCmdAdd() {
	server := &cobra.Command{
		Use:     "server",
		Aliases: []string{"srv"},
		Short:   "Manage servers.",
		Long:    "Manage servers.",
	}

	serverList := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all servers.",
		Long:    "List all servers. Use --company-oid to filter by company.",
		Run:     cmd.runMiddleware.ServerList,
	}

	serverDetail := &cobra.Command{
		Use:     "show --server-oid SERVER_OID",
		Aliases: []string{"get"},
		Short:   "Show server detail.",
		Long:    "Show detailed information about a server.",
		Run:     cmd.runMiddleware.ServerDetail,
	}

	serverStart := &cobra.Command{
		Use:   "start --server-oid SERVER_OID",
		Short: "Start a server.",
		Long:  "Start a stopped server.",
		Run:   cmd.runMiddleware.ServerStart,
	}

	serverStop := &cobra.Command{
		Use:   "stop --server-oid SERVER_OID",
		Short: "Stop a server.",
		Long:  "Gracefully stop a running server (sends ACPI shutdown signal).",
		Run:   cmd.runMiddleware.ServerStop,
	}

	serverRestart := &cobra.Command{
		Use:     "restart --server-oid SERVER_OID",
		Aliases: []string{"reboot"},
		Short:   "Restart a server.",
		Long:    "Gracefully restart a server (sends ACPI reboot signal).",
		Run:     cmd.runMiddleware.ServerRestart,
	}

	serverHardstop := &cobra.Command{
		Use:   "hardstop --server-oid SERVER_OID",
		Short: "Force stop a server.",
		Long:  "Force stop a server immediately (equivalent to pulling the power cord).",
		Run:   cmd.runMiddleware.ServerHardstop,
	}

	serverChangeName := &cobra.Command{
		Use:   "rename --server-oid SERVER_OID --name NEW_NAME",
		Short: "Rename a server.",
		Long:  "Rename a server.",
		Run:   cmd.runMiddleware.ServerChangeName,
	}

	serverMountISO := &cobra.Command{
		Use:     "mount-iso --server-oid SERVER_OID --uri HTTPS_URI",
		Aliases: []string{"li"},
		Short:   "Mount an ISO image to a server.",
		Long:    "Mount a bootable ISO image from HTTPS URL to a server.",
		Run:     cmd.runMiddleware.ServerMountISO,
	}

	serverUmountISO := &cobra.Command{
		Use:     "umount-iso --server-oid SERVER_OID --iso-oid ISO_OID",
		Aliases: []string{"ui"},
		Short:   "Unmount an ISO image from a server.",
		Long:    "Unmount an ISO image from a server.",
		Run:     cmd.runMiddleware.ServerUmountISO,
	}

	ServerAddonsList := &cobra.Command{
		Use:   "addons --server-oid SERVER_OID",
		Short: "List available server addons.",
		Long:  "List available addons (CPU, RAM, Disk) that can be added to a server.",
		Run:   cmd.runMiddleware.ServerAddon,
	}

	serverGetTemplateList := &cobra.Command{
		Use:   "templates",
		Short: "List all available server templates.",
		Long:  "List all available server templates (operating systems).",
		Run:   cmd.runMiddleware.ServerListTemplates,
	}

	serverTermination := &cobra.Command{
		Use:   "termination --server-oid SERVER_OID",
		Short: "Schedule server termination.",
		Long:  "Schedule server termination (deletion).",
		Run:   cmd.runMiddleware.ServerScheduleTermination,
	}

	serverReset := &cobra.Command{
		Use:   "reset --server-oid SERVER_OID --template-oid TEMPLATE_OID",
		Short: "Reset a server to a new template.",
		Long:  "Reset a server to a new template (reinstall OS).",
		Run:   cmd.runMiddleware.ServerReset,
	}

	serverCreate := &cobra.Command{
		Use:   "create --plan PLAN --template-oid OID --payment-method OID --confirm-payment",
		Short: "Create a new server.",
		Long:  "Create a new server.\nGet OS and version list with: titan-sc server templates.\n\nPlans and default resources:\n  SC1: 1 CPU, 1 GB RAM, 10 GB disk\n  SC2: 4 CPU, 4 GB RAM, 80 GB disk\n  SC3: 6 CPU, 8 GB RAM, 100 GB disk",
		Run:   cmd.runMiddleware.ServerCreate,
	}

	// DRP (Disaster Recovery Plan) commands
	serverDrp := &cobra.Command{
		Use:   "drp",
		Short: "Manage server DRP (Disaster Recovery Plan).",
		Long:  "Manage server DRP (Disaster Recovery Plan).\nDRP provides disaster recovery between main and secondary sites.",
	}

	serverDrpStatus := &cobra.Command{
		Use:   "status --server-oid SERVER_OID",
		Short: "Get DRP status for a server.",
		Long:  "Get detailed DRP status for a server including mirroring states and pending operations.",
		Run:   cmd.runMiddleware.ServerDrpStatus,
	}

	serverDrpFailoverSoft := &cobra.Command{
		Use:   "failover-soft --server-oid SERVER_OID",
		Short: "Perform soft failover (server must be stopped).",
		Long:  "Perform soft failover to switch the server to the target site.\nThe server must be stopped before performing a soft failover.\nThis is the safest failover method as it ensures data consistency.",
		Run:   cmd.runMiddleware.ServerDrpFailoverSoft,
	}

	serverDrpFailoverHard := &cobra.Command{
		Use:   "failover-hard --server-oid SERVER_OID --target-site SITE --yes-i-understand-i-will-lose-data",
		Short: "⚠ DANGEROUS: Force failover (will cause data loss).",
		Long: `⚠ DANGEROUS OPERATION: Force failover to target site.

This operation can run while the VM is still running and WILL cause data loss.
Use only in emergencies when soft failover is not possible.

Target site must be specified:
  - main: Primary datacenter
  - secondary: Secondary datacenter

This operation requires the --yes-i-understand-i-will-lose-data flag to confirm
that you understand the risks involved.`,
		Run: cmd.runMiddleware.ServerDrpFailoverHard,
	}

	serverDrpResync := &cobra.Command{
		Use:   "resync --server-oid SERVER_OID --authoritative-site SITE --yes-i-understand-i-will-lose-data",
		Short: "⚠ DANGEROUS: Resync DRP (overrides target data).",
		Long: `⚠ DANGEROUS OPERATION: Resync DRP after split-brain or failure.

This operation will resynchronize DRP from the authoritative site.
ALL DATA ON THE NON-AUTHORITATIVE SITE WILL BE OVERWRITTEN.

Authoritative site must be specified:
  - main: Use main site as the source of truth
  - secondary: Use secondary site as the source of truth

This operation requires the --yes-i-understand-i-will-lose-data flag to confirm
that you understand the risks involved.`,
		Run: cmd.runMiddleware.ServerDrpResync,
	}

	cmd.RootCommand.AddCommand(server)
	server.AddCommand(serverList,
		serverDetail,
		serverStart,
		serverStop,
		serverRestart,
		serverHardstop,
		serverMountISO,
		serverUmountISO,
		serverChangeName,
		ServerAddonsList,
		serverGetTemplateList,
		serverCreate,
		serverTermination,
		serverReset,
		serverDrp)

	// DRP subcommands
	serverDrp.AddCommand(serverDrpStatus, serverDrpFailoverSoft, serverDrpFailoverHard, serverDrpResync)

	// Command arguments
	serverList.Flags().StringP("company-oid", "c", "", "Filter servers by company OID.")

	serverDetail.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = serverDetail.MarkFlagRequired("server-oid")

	serverStart.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = serverStart.MarkFlagRequired("server-oid")

	serverStop.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = serverStop.MarkFlagRequired("server-oid")

	serverRestart.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = serverRestart.MarkFlagRequired("server-oid")

	serverHardstop.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = serverHardstop.MarkFlagRequired("server-oid")

	serverMountISO.Flags().StringP("server-oid", "s", "", "Set server OID.")
	serverMountISO.Flags().StringP("uri", "u", "", "Set remote ISO URI (HTTPS only).")
	_ = serverMountISO.MarkFlagRequired("server-oid")
	_ = serverMountISO.MarkFlagRequired("uri")

	serverUmountISO.Flags().StringP("server-oid", "s", "", "Set server OID.")
	serverUmountISO.Flags().StringP("iso-oid", "i", "", "Set ISO OID.")
	_ = serverUmountISO.MarkFlagRequired("server-oid")
	_ = serverUmountISO.MarkFlagRequired("iso-oid")

	serverChangeName.Flags().StringP("server-oid", "s", "", "Set server OID.")
	serverChangeName.Flags().StringP("name", "n", "", "Set new server's name.")
	_ = serverChangeName.MarkFlagRequired("server-oid")
	_ = serverChangeName.MarkFlagRequired("name")

	// Addon info
	ServerAddonsList.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = ServerAddonsList.MarkFlagRequired("server-oid")

	// server reset
	serverReset.Flags().StringP("server-oid", "s", "", "Set server OID.")
	serverReset.Flags().StringP("template-oid", "", "", "Set template used for create server.")
	serverReset.Flags().StringP("password", "", "", "Set user password.")
	serverReset.Flags().StringP("ssh-keys-name", "", "", "Set ssh keys: keyname1,keyname2,...,keynameN.")
	_ = serverReset.MarkFlagRequired("server-oid")
	_ = serverReset.MarkFlagRequired("template-oid")

	serverTermination.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = serverTermination.MarkFlagRequired("server-oid")

	// Server create
	serverCreate.Flags().StringP("plan", "p", "", "Server plan (SC1, SC2, SC3).")
	serverCreate.Flags().StringP("template-oid", "t", "", "Template OID for the OS.")
	serverCreate.Flags().StringP("password", "", "", "Password for login.")
	serverCreate.Flags().StringP("ssh-keys-name", "", "", "SSH keys: keyname1,keyname2,...,keynameN.")
	serverCreate.Flags().IntP("quantity", "", 1, "Number of servers to create.")
	serverCreate.Flags().IntP("cpu", "c", 0, "Total CPU cores (0 = plan default).")
	serverCreate.Flags().IntP("ram", "r", 0, "Total RAM in GB (0 = plan default).")
	serverCreate.Flags().IntP("disk", "d", 0, "Total disk in GB, must be multiple of 10 (0 = plan default).")
	serverCreate.Flags().StringP("payment-method", "", "", "Payment method OID (required).")
	serverCreate.Flags().StringP("subscription-oid", "", "", "Add server to existing subscription OID (optional, creates new subscription if not set).")
	serverCreate.Flags().BoolP("confirm-payment", "", false, "Confirm the payment (required to proceed).")
	_ = serverCreate.MarkFlagRequired("plan")
	_ = serverCreate.MarkFlagRequired("template-oid")
	_ = serverCreate.MarkFlagRequired("payment-method")

	// DRP status
	serverDrpStatus.Flags().StringP("server-oid", "s", "", "Server OID.")
	_ = serverDrpStatus.MarkFlagRequired("server-oid")

	// DRP failover soft
	serverDrpFailoverSoft.Flags().StringP("server-oid", "s", "", "Server OID.")
	_ = serverDrpFailoverSoft.MarkFlagRequired("server-oid")

	// DRP failover hard
	serverDrpFailoverHard.Flags().StringP("server-oid", "s", "", "Server OID.")
	serverDrpFailoverHard.Flags().StringP("target-site", "", "", "Target site for failover (main or secondary).")
	serverDrpFailoverHard.Flags().BoolP("yes-i-understand-i-will-lose-data", "", false, "Confirm that you understand this operation will cause data loss.")
	_ = serverDrpFailoverHard.MarkFlagRequired("server-oid")
	_ = serverDrpFailoverHard.MarkFlagRequired("target-site")
	_ = serverDrpFailoverHard.MarkFlagRequired("yes-i-understand-i-will-lose-data")

	// DRP resync
	serverDrpResync.Flags().StringP("server-oid", "s", "", "Server OID.")
	serverDrpResync.Flags().StringP("authoritative-site", "", "", "Authoritative site (main or secondary) - this site's data will be preserved.")
	serverDrpResync.Flags().BoolP("yes-i-understand-i-will-lose-data", "", false, "Confirm that you understand this operation will cause data loss.")
	_ = serverDrpResync.MarkFlagRequired("server-oid")
	_ = serverDrpResync.MarkFlagRequired("authoritative-site")
	_ = serverDrpResync.MarkFlagRequired("yes-i-understand-i-will-lose-data")
}
