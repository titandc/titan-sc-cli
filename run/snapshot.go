package run

import (
	"errors"
	"fmt"

	"titan-sc/api"

	"github.com/spf13/cobra"
)

const snapshotTimeFormat = "2006-01-02T15:04:05-07:00"

// getServerIdentifier extracts server OID or UUID from flags.
// Returns (identifier, useLegacy, error).
// If --server-uuid is provided, useLegacy=true and API v1 should be used.
func getServerIdentifier(cmd *cobra.Command) (string, bool, error) {
	serverOID, _ := cmd.Flags().GetString("server-oid")
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	// Mutual exclusivity is enforced by Cobra's MarkFlagsMutuallyExclusive
	if serverOID != "" {
		return serverOID, false, nil
	}
	if serverUUID != "" {
		return serverUUID, true, nil
	}
	return "", false, errors.New("either --server-oid or --server-uuid is required")
}

// getSnapshotIdentifier extracts snapshot OID or UUID from flags.
// Returns (identifier, useLegacy, error).
// If --snapshot-uuid is provided, useLegacy=true and API v1 should be used.
func getSnapshotIdentifier(cmd *cobra.Command) (string, bool, error) {
	snapOID, _ := cmd.Flags().GetString("snapshot-oid")
	snapUUID, _ := cmd.Flags().GetString("snapshot-uuid")

	// Mutual exclusivity is enforced by Cobra's MarkFlagsMutuallyExclusive
	if snapOID != "" {
		return snapOID, false, nil
	}
	if snapUUID != "" {
		return snapUUID, true, nil
	}
	return "", false, errors.New("either --snapshot-oid or --snapshot-uuid is required")
}

func (run *RunMiddleware) SnapshotList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	serverID, useLegacy, err := getServerIdentifier(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}

	var snapshots []api.Snapshot
	var apiReturn *api.Return

	if useLegacy {
		snapshots, apiReturn, err = run.API.ListSnapshotsLegacy(serverID)
	} else {
		snapshots, apiReturn, err = run.API.ListSnapshots(serverID)
	}

	// Render error output
	if err != nil || apiReturn != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}

	// Render success output
	if !run.JSONOutput {
		if len(snapshots) == 0 {
			fmt.Println("Snapshot list is empty.")
			return
		}
		// Show UUID column for legacy API, OID column for v2 API
		var table *Table
		if useLegacy {
			table = NewTable("NAME", "TIMESTAMP", "SIZE", "UUID")
		} else {
			table = NewTable("NAME", "TIMESTAMP", "SIZE", "OID")
		}
		table.SetNoColor(!run.Color)
		for _, snap := range snapshots {
			run.addSnapshotRow(table, &snap, useLegacy)
		}
		table.Print()
		return
	}
	printAsJson(snapshots)
}

func (run *RunMiddleware) SnapshotCreate(cmd *cobra.Command, args []string) {
	// Parse flags
	_ = args
	run.ParseGlobalFlags(cmd)

	serverID, useLegacy, err := getServerIdentifier(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}

	forceErase, _ := cmd.Flags().GetBool("yes-i-agree-to-erase-oldest-snapshot")

	// Execute query - use appropriate API version
	var snapshot *api.SnapshotDetail
	var apiReturn *api.Return

	if useLegacy {
		snapshot, apiReturn, err = run.API.CreateSnapshotLegacy(serverID)
	} else {
		snapshot, apiReturn, err = run.API.CreateSnapshot(serverID)
	}

	if err != nil {
		run.OutputError(err)
		return
	}

	// Check API error
	if apiReturn != nil {
		// Check if it's a limit exceeded error
		// v2 uses error field (Title), v1 uses code field (Code)
		isLimitExceeded := apiReturn.Title == api.SnapshotCreateErrorLimitExceeded ||
			apiReturn.Code == api.SnapshotCreateErrorLimitExceeded
		// API error is fatal unless it's limit exceeded and forceErase is true
		if !(isLimitExceeded && forceErase) {
			run.printAPIReturn(apiReturn)
			return
		}

		// Get list of existing snapshots
		var snapshots []api.Snapshot
		if useLegacy {
			snapshots, apiReturn, err = run.API.ListSnapshotsLegacy(serverID)
		} else {
			snapshots, apiReturn, err = run.API.ListSnapshots(serverID)
		}
		if err != nil || apiReturn != nil {
			run.handleErrorAndGenericOutput(apiReturn, err)
			return
		}

		// Find the oldest one
		oldestSnapshot, err := getOldestSnapshotFromList(snapshots)
		if err != nil {
			run.OutputError(err)
			return
		}

		// Delete oldest snapshot
		if useLegacy {
			// In legacy mode, use UUID field for the snapshot identifier
			apiReturn, err = run.API.DeleteSnapshotLegacy(serverID, oldestSnapshot.UUID)
		} else {
			apiReturn, err = run.API.DeleteSnapshot(oldestSnapshot.OID)
		}
		if err != nil {
			run.OutputError(err)
			return
		}
		if apiReturn != nil && apiReturn.Error() {
			run.printAPIReturn(apiReturn)
			return
		}

		// Create new snapshot
		if useLegacy {
			snapshot, apiReturn, err = run.API.CreateSnapshotLegacy(serverID)
		} else {
			snapshot, apiReturn, err = run.API.CreateSnapshot(serverID)
		}
		if err != nil || apiReturn != nil {
			run.handleErrorAndGenericOutput(apiReturn, err)
			return
		}
	}

	// Render success output
	if !run.JSONOutput {
		fmt.Printf("%s\n", run.Colorize("Creating new snapshot:", "green"))
		// Show UUID column for legacy API, OID column for v2 API
		var table *Table
		if useLegacy {
			table = NewTable("NAME", "TIMESTAMP", "SIZE", "UUID")
		} else {
			table = NewTable("NAME", "TIMESTAMP", "SIZE", "OID")
		}
		table.SetNoColor(!run.Color)
		run.addSnapshotRow(table, &snapshot.Snapshot, useLegacy)
		table.Print()
		return
	}
	printAsJson(snapshot)
}

func (run *RunMiddleware) SnapshotDelete(cmd *cobra.Command, args []string) {
	// Parse flags
	_ = args
	run.ParseGlobalFlags(cmd)

	// Check for API v2 (OID) or API v1 (UUID) mode
	snapOID, _ := cmd.Flags().GetString("snapshot-oid")
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	snapUUID, _ := cmd.Flags().GetString("snapshot-uuid")

	// Validate flags
	if snapOID != "" && (serverUUID != "" || snapUUID != "") {
		run.OutputError(errors.New("cannot mix --snapshot-oid with legacy --server-uuid/--snapshot-uuid flags"))
		return
	}

	var apiReturn *api.Return
	var err error

	if snapOID != "" {
		// API v2 mode: only need snapshot OID
		apiReturn, err = run.API.DeleteSnapshot(snapOID)
	} else if serverUUID != "" && snapUUID != "" {
		// API v1 legacy mode: need both server UUID and snapshot UUID
		apiReturn, err = run.API.DeleteSnapshotLegacy(serverUUID, snapUUID)
	} else if serverUUID != "" || snapUUID != "" {
		run.OutputError(errors.New("legacy mode requires both --server-uuid and --snapshot-uuid"))
		return
	} else {
		run.OutputError(errors.New("--snapshot-oid is required (or --server-uuid and --snapshot-uuid for legacy mode)"))
		return
	}

	// Format output
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) SnapshotRotate(cmd *cobra.Command, args []string) {
	// Parse flags
	_ = args
	run.ParseGlobalFlags(cmd)

	serverID, useLegacy, err := getServerIdentifier(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}

	forceRotation, _ := cmd.Flags().GetBool("force")

	// Try to create a snapshot - use appropriate API version
	var snapshot *api.SnapshotDetail
	var apiReturn *api.Return

	if useLegacy {
		snapshot, apiReturn, err = run.API.CreateSnapshotLegacy(serverID)
	} else {
		snapshot, apiReturn, err = run.API.CreateSnapshot(serverID)
	}

	if err != nil {
		run.OutputError(err)
		return
	}

	// Check API error
	if apiReturn != nil {
		// Check if it's a limit exceeded error (need to rotate)
		// v2 uses error field (Title), v1 uses code field (Code)
		isLimitExceeded := apiReturn.Title == api.SnapshotCreateErrorLimitExceeded ||
			apiReturn.Code == api.SnapshotCreateErrorLimitExceeded
		if !isLimitExceeded {
			// We had a fatal error (not limit exceeded)
			run.printAPIReturn(apiReturn)
			return
		}

		// Get list of existing snapshots
		var snapshots []api.Snapshot
		if useLegacy {
			snapshots, apiReturn, err = run.API.ListSnapshotsLegacy(serverID)
		} else {
			snapshots, apiReturn, err = run.API.ListSnapshots(serverID)
		}
		if err != nil || apiReturn != nil {
			run.handleErrorAndGenericOutput(apiReturn, err)
			return
		}

		// Find the oldest one
		oldestSnapshot, err := getOldestSnapshotFromList(snapshots)
		if err != nil {
			run.OutputError(err)
			return
		}
		// Prompt user if needed
		if !forceRotation {
			var idLabel, idValue string
			if useLegacy {
				idLabel = "UUID"
				idValue = oldestSnapshot.UUID
			} else {
				idLabel = "OID"
				idValue = oldestSnapshot.OID
			}
			promptString := fmt.Sprint("This action will immediately delete snapshot '", oldestSnapshot.Name,
				"' (", idLabel, ": ", idValue, ") created at ", DatePtrFormat(oldestSnapshot.CreatedAt), ".",
				" \nAre you sure you want to continue? (y/N): ")
			lowerText := keyboardPromptToLower(promptString)
			// If the response is anything other than "yes" or "y"
			if lowerText != "y" && lowerText != "yes" {
				fmt.Println("The response was something other than 'yes' or 'y', so no further actions will " +
					"be taken. No snapshots have been deleted, and no new snapshots will be created.")
				return
			}
		}

		// Delete oldest snapshot
		if useLegacy {
			// In legacy mode, use UUID field for the snapshot identifier
			apiReturn, err = run.API.DeleteSnapshotLegacy(serverID, oldestSnapshot.UUID)
		} else {
			apiReturn, err = run.API.DeleteSnapshot(oldestSnapshot.OID)
		}
		if err != nil {
			run.OutputError(err)
			return
		}
		if apiReturn != nil && apiReturn.Error() {
			run.printAPIReturn(apiReturn)
			return
		}

		// Create new snapshot
		if useLegacy {
			snapshot, apiReturn, err = run.API.CreateSnapshotLegacy(serverID)
		} else {
			snapshot, apiReturn, err = run.API.CreateSnapshot(serverID)
		}
		if err != nil || apiReturn != nil {
			run.handleErrorAndGenericOutput(apiReturn, err)
			return
		}
	}

	// Render success output
	if !run.JSONOutput {
		fmt.Printf("%s\n", run.Colorize("Creating new snapshot:", "green"))
		// Show UUID column for legacy API, OID column for v2 API
		var table *Table
		if useLegacy {
			table = NewTable("NAME", "TIMESTAMP", "SIZE", "UUID")
		} else {
			table = NewTable("NAME", "TIMESTAMP", "SIZE", "OID")
		}
		table.SetNoColor(!run.Color)
		run.addSnapshotRow(table, &snapshot.Snapshot, useLegacy)
		table.Print()
		return
	}
	printAsJson(snapshot)
}

func (run *RunMiddleware) SnapshotRestore(cmd *cobra.Command, args []string) {
	// Parse flags
	_ = args
	run.ParseGlobalFlags(cmd)

	snapID, useLegacy, err := getSnapshotIdentifier(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}

	if useLegacy {
		// Restore was not available in API v1 - only v2 supports it
		run.OutputError(errors.New("snapshot restore is not supported in legacy mode (API v1); use --snapshot-oid instead"))
		return
	}

	// Execute query
	apiReturn, err := run.API.RestoreSnapshot(snapID)

	// Format output
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) addSnapshotRow(table *Table, snap *api.Snapshot, useLegacy bool) {
	size := fmt.Sprintf("%d %s", snap.Size.Value, snap.Size.Unit)

	var timestampColorFn func(string) string
	if run.Color {
		timestampColorFn = ColorFn("dim")
	}

	// Show UUID for legacy API, OID for v2 API
	var idCol TableColumn
	if useLegacy {
		idCol = TableColumn{Value: snap.UUID, ColorFn: ColorFn("blue")}
	} else {
		idCol = ColOID(snap.OID)
	}

	table.AddRow(
		ColName(snap.Name),
		ColColor(DatePtrFormat(snap.CreatedAt), timestampColorFn),
		Col(size),
		idCol,
	)
}

func getOldestSnapshotFromList(snapshots []api.Snapshot) (api.Snapshot, error) {
	if len(snapshots) == 0 {
		return api.Snapshot{}, errors.New("empty snapshot list")
	}
	// Because the program cannot unmarshal the time from the snapshot returned by the go-api
	// into an int64, it must be passed as a string and parsed with time.Parse
	//oldestDate := millisecondsToTime(snapshots[0].CreatedAt)
	oldestDate := *snapshots[0].CreatedAt
	oldestSnapshot := snapshots[0]
	for _, snap := range snapshots {
		currentDate := *snap.CreatedAt
		if currentDate < oldestDate {
			oldestSnapshot = snap
			oldestDate = currentDate
		}
	}
	return oldestSnapshot, nil
}
