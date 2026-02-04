package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) SnapshotCmdAdd() {

	snapshot := &cobra.Command{
		Use:     "snapshot",
		Aliases: []string{"snap"},
		Short:   "Manage servers' snapshots.",
		GroupID: "resources",
	}

	snapshotList := &cobra.Command{
		Use:     "list --server-oid SERVER_OID",
		Aliases: []string{"ls"},
		Short:   "List all snapshots of a server.",
		Long: `List all snapshots of a server.

Examples:
  # Using API v2 (recommended)
  titan-sc snapshot list --server-oid sc-abc123

  # Using API v1 (legacy, deprecated - will be removed in future versions)
  titan-sc snapshot list --server-uuid 12345678-1234-1234-1234-123456789abc`,
		Run: cmd.runMiddleware.SnapshotList,
	}

	snapshotCreate := &cobra.Command{
		Use:   "create --server-oid SERVER_OID",
		Short: "Create a snapshot of a server.",
		Long: `Create a new snapshot of a server.

Examples:
  # Using API v2 (recommended)
  titan-sc snapshot create --server-oid sc-abc123

  # Force create when quota is reached
  titan-sc snapshot create --server-oid sc-abc123 --yes-i-agree-to-erase-oldest-snapshot

  # Using API v1 (legacy, deprecated - will be removed in future versions)
  titan-sc snapshot create --server-uuid 12345678-1234-1234-1234-123456789abc`,
		Run: cmd.runMiddleware.SnapshotCreate,
	}

	snapshotDelete := &cobra.Command{
		Use:     "delete --snapshot-oid SNAPSHOT_OID",
		Aliases: []string{"del"},
		Short:   "Delete a server's snapshot.",
		Long: `Delete a server's snapshot.

Examples:
  # Using API v2 (recommended)
  titan-sc snapshot delete --snapshot-oid snap-xyz789

  # Using API v1 (legacy, deprecated - will be removed in future versions)
  titan-sc snapshot delete --server-uuid 12345678-... --snapshot-uuid 87654321-...`,
		Run: cmd.runMiddleware.SnapshotDelete,
	}

	snapshotRotate := &cobra.Command{
		Use:   "rotate --server-oid SERVER_OID",
		Short: "Rotate the server's snapshots.",
		Long: `Create a new snapshot and delete the oldest one if necessary.

Examples:
  # Using API v2 (recommended)
  titan-sc snapshot rotate --server-oid sc-abc123

  # Force rotation without confirmation
  titan-sc snapshot rotate --server-oid sc-abc123 --force

  # Using API v1 (legacy, deprecated - will be removed in future versions)
  titan-sc snapshot rotate --server-uuid 12345678-1234-1234-1234-123456789abc --force`,
		Run: cmd.runMiddleware.SnapshotRotate,
	}

	snapshotRestore := &cobra.Command{
		Use:   "restore --snapshot-oid SNAPSHOT_OID",
		Short: "Restore a server snapshot.",
		Long: `Restore a server snapshot.

WARNING: The server must be stopped before restoring. This operation erases all
data on the server's disk and replaces it with the snapshot content.

Examples:
  # Using API v2 (recommended)
  titan-sc snapshot restore --snapshot-oid snap-xyz789

  # Using API v1 (legacy, deprecated - will be removed in future versions)
  titan-sc snapshot restore --snapshot-uuid 87654321-1234-1234-1234-123456789abc`,
		Run: cmd.runMiddleware.SnapshotRestore,
	}

	cmd.RootCommand.AddCommand(snapshot)
	snapshot.AddCommand(snapshotList, snapshotCreate, snapshotDelete, snapshotRotate, snapshotRestore)

	// Delete: OID or legacy UUID (legacy requires both server and snapshot UUID)
	snapshotDelete.Flags().StringP("snapshot-oid", "o", "", "Set snapshot OID (API v2).")
	snapshotDelete.Flags().StringP("server-uuid", "u", "", "Legacy: Set server UUID (API v1, requires --snapshot-uuid).")
	snapshotDelete.Flags().StringP("snapshot-uuid", "s", "", "Legacy: Set snapshot UUID (API v1, requires --server-uuid).")
	snapshotDelete.MarkFlagsMutuallyExclusive("snapshot-oid", "server-uuid")
	snapshotDelete.MarkFlagsMutuallyExclusive("snapshot-oid", "snapshot-uuid")

	// Create: OID or legacy UUID
	snapshotCreate.Flags().StringP("server-oid", "s", "", "Set server OID (API v2).")
	snapshotCreate.Flags().StringP("server-uuid", "u", "", "Legacy: Set server UUID (API v1).")
	snapshotCreate.Flags().BoolP("yes-i-agree-to-erase-oldest-snapshot", "", false,
		"Automatically erase oldest snapshot if quota has been reached.")
	snapshotCreate.MarkFlagsMutuallyExclusive("server-oid", "server-uuid")

	// List: OID or legacy UUID
	snapshotList.Flags().StringP("server-oid", "s", "", "Set server OID (API v2).")
	snapshotList.Flags().StringP("server-uuid", "u", "", "Legacy: Set server UUID (API v1).")
	snapshotList.MarkFlagsMutuallyExclusive("server-oid", "server-uuid")

	// Rotate: OID or legacy UUID
	snapshotRotate.Flags().StringP("server-oid", "s", "", "Set server OID (API v2).")
	snapshotRotate.Flags().StringP("server-uuid", "u", "", "Legacy: Set server UUID (API v1).")
	snapshotRotate.Flags().BoolP("force", "f", false, "Force the rotation. "+
		"The oldest snapshot will be automatically deleted without prompting.")
	snapshotRotate.MarkFlagsMutuallyExclusive("server-oid", "server-uuid")

	// Restore: OID or legacy UUID
	snapshotRestore.Flags().StringP("snapshot-oid", "o", "", "Set snapshot OID (API v2).")
	snapshotRestore.Flags().StringP("snapshot-uuid", "s", "", "Legacy: Set snapshot UUID (API v1).")
	snapshotRestore.MarkFlagsMutuallyExclusive("snapshot-oid", "snapshot-uuid")
}
