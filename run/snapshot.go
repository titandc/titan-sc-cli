package run

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"titan-sc/api"
)

const snapshotTimeFormat = "2006-01-02T15:04:05-07:00"

func (run *RunMiddleware) SnapshotList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	snapshots, apiReturn, err := run.API.ListSnapshots(serverUUID)

	// Render error output
	if err != nil || apiReturn != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}

	// Render success output
	if run.HumanReadable {
		if len(snapshots) == 0 {
			fmt.Println("Snapshot list is empty.")
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "SNAPSHOT UUID\tTIMESTAMP\tSIZE\tNAME\t\n")
		for _, snap := range snapshots {
			printSnapshotInfos(w, &snap)
		}
		return
	}
	printAsJson(snapshots)
}

func (run *RunMiddleware) SnapshotCreate(cmd *cobra.Command, args []string) {
	// Parse flags
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	forceErase, _ := cmd.Flags().GetBool("yes-i-agree-to-erase-oldest-snapshot")

	// Execute query
	snapshot, apiReturn, err := run.API.CreateSnapshot(serverUUID)
	if err != nil {
		run.OutputError(err)
		return
	}

	// Check API error
	if apiReturn != nil {
		// API error is fatal
		if !(apiReturn.Code == api.SnapshotCreateErrorLimitExceeded && forceErase) {
			run.printAPIReturn(apiReturn)
			return
		}

		// Get list of existing snapshots
		snapshots, apiReturn, err := run.API.ListSnapshots(serverUUID)
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
		apiReturn, err = run.API.DeleteSnapshot(serverUUID, oldestSnapshot.UUID)
		if err != nil {
			run.OutputError(err)
			return
		}
		if apiReturn != nil && apiReturn.Error != "" {
			run.printAPIReturn(apiReturn)
			return
		}

		// Create new snapshot
		snapshot, apiReturn, err = run.API.CreateSnapshot(serverUUID)
		if err != nil || apiReturn != nil {
			run.handleErrorAndGenericOutput(apiReturn, err)
			return
		}
	}

	// Render success output
	if run.HumanReadable {
		fmt.Println("Creating new snapshot:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "SNAPSHOT UUID\tTIMESTAMP\tSIZE\t\tNAME\t\n")
		printSnapshotInfos(w, snapshot)
		return
	}
	printAsJson(snapshot)
}

func (run *RunMiddleware) SnapshotDelete(cmd *cobra.Command, args []string) {
	// Parse flags
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	snapUUID, _ := cmd.Flags().GetString("snapshot-uuid")

	// Execute query
	apiReturn, err := run.API.DeleteSnapshot(serverUUID, snapUUID)

	// Format output
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) SnapshotRotate(cmd *cobra.Command, args []string) {
	// Parse flags
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	forceRotation, _ := cmd.Flags().GetBool("force")

	// Try to create a snapshot
	snapshot, apiReturn, err := run.API.CreateSnapshot(serverUUID)
	if err != nil {
		run.OutputError(err)
		return
	}

	// Check API error
	if apiReturn != nil {
		// We had a fatal error
		if apiReturn.Code != api.SnapshotCreateErrorLimitExceeded {
			run.printAPIReturn(apiReturn)
			return
		}

		// Get list of existing snapshots
		snapshots, apiReturn, err := run.API.ListSnapshots(serverUUID)
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
			promptString := fmt.Sprint("This action will immediately delete snapshot '", oldestSnapshot.Name,
				"' (UUID: ", oldestSnapshot.UUID, ") created at ", FormatTimestampToString(oldestSnapshot.CreatedAt), ".",
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
		apiReturn, err = run.API.DeleteSnapshot(serverUUID, oldestSnapshot.UUID)
		if err != nil {
			run.OutputError(err)
			return
		}
		if apiReturn != nil && apiReturn.Error != "" {
			run.printAPIReturn(apiReturn)
			return
		}

		// Create new snapshot
		snapshot, apiReturn, err = run.API.CreateSnapshot(serverUUID)
		if err != nil || apiReturn != nil {
			run.handleErrorAndGenericOutput(apiReturn, err)
			return
		}
	}

	// Render success output
	if run.HumanReadable {
		fmt.Println("Creating new snapshot:")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "SNAPSHOT UUID\tTIMESTAMP\tSIZE\t\tNAME\t\n")
		printSnapshotInfos(w, snapshot)
		return
	}
	printAsJson(snapshot)
}

func (run *RunMiddleware) SnapshotRestore(cmd *cobra.Command, args []string) {
	// Parse falgs
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	snapUUID, _ := cmd.Flags().GetString("snapshot-uuid")

	// Execute query
	apiReturn, err := run.API.RestoreSnapshot(serverUUID, snapUUID)

	// Format output
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func printSnapshotInfos(w *tabwriter.Writer, snap *api.APISnapshot) {
	_, _ = fmt.Fprintf(w, "%s\t%s\t%d %s\t%s\t\n",
		snap.UUID, FormatTimestampToString(snap.CreatedAt), snap.Size.Value,
		snap.Size.Unit, snap.Name)
	_ = w.Flush()
}

func getOldestSnapshotFromList(snapshots []api.APISnapshot) (api.APISnapshot, error) {
	if len(snapshots) == 0 {
		return api.APISnapshot{}, errors.New("empty snapshot list")
	}
	// Because the program cannot unmarshal the time from the snapshot returned by the go-api
	// into an int64, it must be passed as a string and parsed with time.Parse
	//oldestDate := millisecondsToTime(snapshots[0].CreatedAt)
	oldestDate := snapshots[0].CreatedAt
	oldestSnapshot := snapshots[0]
	for _, snap := range snapshots {
		currentDate := snap.CreatedAt
		if currentDate < oldestDate {
			oldestSnapshot = snap
			oldestDate = currentDate
		}
	}
	return oldestSnapshot, nil
}

func FormatTimestampToString(createdAt int64) string {
	return millisecondsToTime(createdAt).Format(snapshotTimeFormat)
}
