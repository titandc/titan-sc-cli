package cmd

import (
    "github.com/spf13/cobra"
    . "titan-sc/api"
)

var snapshot = &cobra.Command{
    Use: "snapshot",
    Aliases: []string{"snap"},
    Short: "Manage server's snapshots.",
}

var snapshotList = &cobra.Command{
    Use: "list server_uuid",
    Aliases: []string{"ls"},
    Short: "Remove one server snapshot.",
    Long: "Remove one server snapshot (need server UUID).",
    Args: cmdNeed1UUID,
    Run: API.SnapshotList,
}

var snapshotCreate = &cobra.Command{
    Use: "create server-uuid",
    Short: "Create a new server snapshopt.",
    Long: "Create a new snapshot for a server (need server UUID).",
    Args: cmdNeed1UUID,
    Run: API.SnapshotCreate,
}

var snapshotRemove = &cobra.Command{
    Use: "remove [--server-uuid --snap-uuid]",
    Aliases: []string{"rm"},
    Short: "Remove one server snapshot.",
    Long: "Remove one server snapshot.",
    Run: API.SnapshotRemove,
}

func snapshotCmdAdd() {
    rootCmd.AddCommand(snapshot)
    snapshot.AddCommand(snapshotList, snapshotCreate, snapshotRemove)
    snapshotRemove.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
    snapshotRemove.Flags().StringP("snapshot-uuid", "s", "", "Set snapshot UUID.")
}
