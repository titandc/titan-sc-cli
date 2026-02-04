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
		GroupID: "resources",
	}

	serverList := &cobra.Command{
		Use:     "list [--company-oid COMPANY_OID]",
		Aliases: []string{"ls"},
		Short:   "List all servers.",
		Long: `List all servers within your company.

If --company-oid is not specified, your default company will be used.`,
		Run: cmd.runMiddleware.ServerList,
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

	// ISO subcommand group
	serverISO := &cobra.Command{
		Use:   "iso",
		Short: "Manage server ISO images.",
		Long:  "Manage server ISO images (mount, unmount, show).",
	}

	serverISOMount := &cobra.Command{
		Use:   "mount --server-oid SERVER_OID --uri HTTPS_URI",
		Short: "Mount an ISO image to a server.",
		Long:  "Mount a bootable ISO image from HTTPS URL to a server.",
		Run:   cmd.runMiddleware.ServerISOMount,
	}

	serverISOUmount := &cobra.Command{
		Use:     "umount --server-oid SERVER_OID [--iso-oid ISO_OID]",
		Aliases: []string{"unmount"},
		Short:   "Unmount an ISO image from a server.",
		Long:    "Unmount an ISO image from a server. If --iso-oid is not specified and only one ISO is mounted, it will be unmounted automatically.",
		Run:     cmd.runMiddleware.ServerISOUmount,
	}

	serverISOShow := &cobra.Command{
		Use:     "show --server-oid SERVER_OID",
		Aliases: []string{"list", "ls"},
		Short:   "Show mounted ISOs on a server.",
		Long:    "Show all currently mounted ISO images on a server.",
		Run:     cmd.runMiddleware.ServerISOShow,
	}

	serverAddons := &cobra.Command{
		Use:   "addons",
		Short: "Manage server addons (CPU, RAM, Disk).",
		Long:  "Manage server addons (CPU, RAM, Disk).",
	}

	serverAddonsList := &cobra.Command{
		Use:   "list --server-oid SERVER_OID",
		Short: "List available server addons.",
		Long:  "List available addons (CPU, RAM, Disk) that can be added to a server.",
		Run:   cmd.runMiddleware.ServerAddon,
	}

	serverTermination := &cobra.Command{
		Use:    "delete --server-oid SERVER_OID",
		Short:  "Schedule server deletion.",
		Long:   "Schedule server deletion (termination).",
		Run:    cmd.runMiddleware.ServerScheduleTermination,
		Hidden: true, // Hidden until properly tested (requires payment)
	}

	serverReset := &cobra.Command{
		Use:   "reset --server-oid SERVER_OID --template-oid TEMPLATE_OID",
		Short: "Reset a server to a new template.",
		Long:  "Reset a server to a new template (reinstall OS).",
		Run:   cmd.runMiddleware.ServerReset,
	}

	serverCreate := &cobra.Command{
		Use:    "create --plan PLAN --template-oid OID --confirm-payment",
		Short:  "Create a new server.",
		Long:   "Create a new server.\nGet OS and version list with: titan-sc template list.\n\nPlans and default resources:\n  SC1: 1 CPU, 1 GB RAM, 10 GB disk\n  SC2: 4 CPU, 4 GB RAM, 80 GB disk\n  SC3: 6 CPU, 8 GB RAM, 100 GB disk",
		Run:    cmd.runMiddleware.ServerCreate,
		Hidden: true, // Hidden until properly tested (requires payment)
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
		serverISO,
		serverChangeName,
		serverAddons,
		serverCreate,
		serverTermination,
		serverReset,
		serverDrp)

	// ISO subcommands
	serverISO.AddCommand(serverISOMount, serverISOUmount, serverISOShow)

	// Addons subcommands
	serverAddons.AddCommand(serverAddonsList)

	// DRP subcommands
	serverDrp.AddCommand(serverDrpStatus, serverDrpFailoverSoft, serverDrpFailoverHard, serverDrpResync)

	// Command arguments
	serverList.Flags().StringP("company-oid", "c", "", "Company OID (uses your default company if not specified).")

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

	// ISO mount
	serverISOMount.Flags().StringP("server-oid", "s", "", "Set server OID.")
	serverISOMount.Flags().StringP("uri", "u", "", "Set remote ISO URI (HTTPS only).")
	_ = serverISOMount.MarkFlagRequired("server-oid")
	_ = serverISOMount.MarkFlagRequired("uri")

	// ISO umount
	serverISOUmount.Flags().StringP("server-oid", "s", "", "Set server OID.")
	serverISOUmount.Flags().StringP("iso-oid", "i", "", "Set ISO OID (auto-detected if only one ISO is mounted).")
	_ = serverISOUmount.MarkFlagRequired("server-oid")

	// ISO show
	serverISOShow.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = serverISOShow.MarkFlagRequired("server-oid")

	serverChangeName.Flags().StringP("server-oid", "s", "", "Set server OID.")
	serverChangeName.Flags().StringP("name", "n", "", "Set new server's name.")
	_ = serverChangeName.MarkFlagRequired("server-oid")
	_ = serverChangeName.MarkFlagRequired("name")

	// Addon info
	serverAddonsList.Flags().StringP("server-oid", "s", "", "Set server OID.")
	_ = serverAddonsList.MarkFlagRequired("server-oid")

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
	serverCreate.Flags().StringP("payment-method", "", "", "Payment method OID (uses default if not specified).")
	serverCreate.Flags().StringP("subscription-oid", "", "", "Add server to existing subscription OID (optional, creates new subscription if not set).")
	serverCreate.Flags().BoolP("confirm-payment", "", false, "Confirm the payment (required to proceed).")
	_ = serverCreate.MarkFlagRequired("plan")
	_ = serverCreate.MarkFlagRequired("template-oid")

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

	// DRP resync
	serverDrpResync.Flags().StringP("server-oid", "s", "", "Server OID.")
	serverDrpResync.Flags().StringP("authoritative-site", "", "", "Authoritative site (main or secondary) - this site's data will be preserved.")
	serverDrpResync.Flags().BoolP("yes-i-understand-i-will-lose-data", "", false, "Confirm that you understand this operation will cause data loss.")
	_ = serverDrpResync.MarkFlagRequired("server-oid")
}
