package cmd

import (
    "github.com/spf13/cobra"
    . "titan-sc/api"
)

var snapshot = &cobra.Command{
    Use: "snapshot",
    Aliases: []string{"snap"},
    Short: "Manage servers' snapshots.",
}

var snapshotList = &cobra.Command{
    Use: "list SERVER_UUID",
    Aliases: []string{"ls"},
    Short: "List all snapshots of a server.",
    Long: "List all snapshots of a server.",
    Args: cmdNeed1UUID,
    Run: API.SnapshotList,
}

var snapshotCreate = &cobra.Command{
    Use: "create SERVER_UUID",
    Short: "Create a snapshot of a server.",
    Long: "Create a new snapshot of a server.",
    Args: cmdNeed1UUID,
    Run: API.SnapshotCreate,
}

var snapshotDelete = &cobra.Command{
    Use: "delete --server-uuid SERVER_UUID --snapshot-uuid SNAPSHOT_UUID",
    Aliases: []string{"del"},
    Short: "Delete a server's snapshot.",
    Long: "Delete a server's snapshot.",
    Run: API.SnapshotRemove,
}

func snapshotCmdAdd() {
    rootCmd.AddCommand(snapshot)
    snapshot.AddCommand(snapshotList, snapshotCreate, snapshotDelete)
    snapshotDelete.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
    snapshotDelete.Flags().StringP("snapshot-uuid", "s", "", "Set snapshot UUID.")
}
