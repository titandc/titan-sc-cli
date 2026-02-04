package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) SnapshotCmdAdd() {

	snapshot := &cobra.Command{
		Use:     "snapshot",
		Aliases: []string{"snap"},
		Short:   "Manage servers' snapshots.",
	}

	snapshotList := &cobra.Command{
		Use:     "list --server-uuid SERVER_UUID",
		Aliases: []string{"ls"},
		Short:   "List all snapshots of a server.",
		Long:    "List all snapshots of a server.",
		Run:     cmd.runMiddleware.SnapshotList,
	}

	snapshotCreate := &cobra.Command{
		Use:   "create --server-uuid SERVER_UUID",
		Short: "Create a snapshot of a server.",
		Long:  "Create a new snapshot of a server.",
		Run:   cmd.runMiddleware.SnapshotCreate,
	}

	snapshotDelete := &cobra.Command{
		Use:     "delete --server-uuid SERVER_UUID --snapshot-uuid SNAPSHOT_UUID",
		Aliases: []string{"del"},
		Short:   "Delete a server's snapshot.",
		Long:    "Delete a server's snapshot.",
		Run:     cmd.runMiddleware.SnapshotDelete,
	}

	snapshotRotate := &cobra.Command{
		Use:   "rotate --server-uuid SERVER_UUID",
		Short: "Rotate the server's snapshots.",
		Long:  "Create a new snapshot and delete the oldest one if necessary.",
		Run:   cmd.runMiddleware.SnapshotRotate,
	}

	cmd.RootCommand.AddCommand(snapshot)
	snapshot.AddCommand(snapshotList, snapshotCreate, snapshotDelete, snapshotRotate)

	snapshotDelete.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
	snapshotDelete.Flags().StringP("snapshot-uuid", "s", "", "Set snapshot UUID.")

	_ = snapshotDelete.MarkFlagRequired("server-uuid")
	_ = snapshotDelete.MarkFlagRequired("snapshot-uuid")

	snapshotCreate.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
	snapshotCreate.Flags().BoolP("yes-i-agree-to-erase-oldest-snapshot", "", false,
		"Automatically erase oldest snapshot if quota has been reached.")
	_ = snapshotCreate.MarkFlagRequired("server-uuid")

	snapshotList.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
	_ = snapshotList.MarkFlagRequired("server-uuid")

	snapshotRotate.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
	snapshotRotate.Flags().BoolP("force", "f", false, "Force the rotation. "+
		"The oldest snapshot will be automatically deleted without prompting.")
	_ = snapshotRotate.MarkFlagRequired("server-uuid")
}
