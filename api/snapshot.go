package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"text/tabwriter"
)

/*
 *
 *
 **************************
 * Snapshot server function
 **************************
 *
 *
 */
func (API *APITitan) SnapshotList(cmd *cobra.Command, args []string) {

	serverUUID := args[0]
	API.ParseGlobalFlags(cmd)

	snapshots, err := API.SnapshotServerUUIDList(serverUUID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(snapshots) == 0 {
		fmt.Println("0 Snapshot")
		return
	}

	for _, snap := range snapshots {
		if !API.HumanReadable {
			API.PrintJson()
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			_, _ = fmt.Fprintf(w, "SNAPSHOT UUID\tTIMESTAMP\tSIZE\tNAME\t\n")
			API.PrintSnapshotInfos(w, &snap)
		}
	}
}

func (API *APITitan) SnapshotCreate(cmd *cobra.Command, args []string) {

	serverUUID := args[0]
	API.ParseGlobalFlags(cmd)

	err := API.SendAndResponse(HTTPPost, "/compute/servers/"+serverUUID+"/snapshots", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		snap := &APISnapshot{}
		if err := json.Unmarshal(API.RespBody, &snap); err != nil {
			log.Println("Human readable format error: ", err.Error())
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "SNAPSHOT UUID\tTIMESTAMP\tSIZE\t\tNAME\t\n")
		API.PrintSnapshotInfos(w, snap)
	}
}

func (API *APITitan) SnapshotRemove(cmd *cobra.Command, args []string) {

	_ = args
	API.ParseGlobalFlags(cmd)

	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	snapUUID, _ := cmd.Flags().GetString("snapshot-uuid")

	if snapUUID == "" {
		fmt.Println("Error: Snapshot UUID missing.\n" +
			"Example: ./titan-sc snap rm --server-uuid SERVER_UUID --snapshot-uuid SNAP_UUID")
		return
	}
	if serverUUID == "" {
		fmt.Println("Error: Server UUID missing.\n" +
			"Example: ./titan-sc snap rm --server-uuid SERVER_UUID --snapshot-uuid SNAP_UUID")
		return
	}

	snapshots, err := API.SnapshotServerUUIDList(serverUUID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, snap := range snapshots {
		if snapUUID != snap.UUID {
			continue
		}

		err = API.SendAndResponse(HTTPDelete, "/compute/servers/"+
			serverUUID+"/snapshots/"+snapUUID, nil)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Snapshot deleting request successfully sent.")
		return
	}
	fmt.Println("Snapshot UUID", snapUUID, "not found")
}

func (API *APITitan) SnapshotServerUUIDList(serverUUID string) ([]APISnapshot, error) {

	err := API.SendAndResponse(HTTPGet, "/compute/servers/"+serverUUID+"/snapshots", nil)
	if err != nil {
		return nil, err
	}

	snap := make([]APISnapshot, 0)
	if err := json.Unmarshal(API.RespBody, &snap); err != nil {
		return nil, err
	}
	return snap, nil
}

func (API *APITitan) PrintSnapshotInfos(w *tabwriter.Writer, snap *APISnapshot) {

	_, _ = fmt.Fprintf(w, "%s\t%s\t%d %s\t%s\t\n",
		snap.UUID, snap.CreatedAt, snap.Size.Value,
		snap.Size.Unit, snap.Name)
	_ = w.Flush()
}
