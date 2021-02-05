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
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	err := API.SendAndResponse(HTTPGet, "/compute/servers/"+serverUUID+"/snapshots", nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var snapshots []APISnapshot
	if err := json.Unmarshal(API.RespBody, &snapshots); err != nil {
		if !API.HumanReadable {
			API.PrintJson()
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		if len(snapshots) == 0 {
			fmt.Println("Snapshot list is empty.")
			return
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "SNAPSHOT UUID\tTIMESTAMP\tSIZE\tNAME\t\n")
		for _, snap := range snapshots {
			API.PrintSnapshotInfos(w, &snap)
		}
	}
}

func (API *APITitan) SnapshotCreate(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

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

func (API *APITitan) SnapshotDelete(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	snapUUID, _ := cmd.Flags().GetString("snapshot-uuid")

	err := API.SendAndResponse(HTTPDelete, "/compute/servers/"+
		serverUUID+"/snapshots/"+snapUUID, nil)
	if err != nil {
		if !API.HumanReadable {
			API.PrintJson()
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	if !API.HumanReadable {
		API.PrintJson()
	} else {
		fmt.Println("Snapshot deleting request successfully sent.")
	}
}

func (API *APITitan) PrintSnapshotInfos(w *tabwriter.Writer, snap *APISnapshot) {

	_, _ = fmt.Fprintf(w, "%s\t%s\t%d %s\t%s\t\n",
		snap.UUID, snap.CreatedAt, snap.Size.Value,
		snap.Size.Unit, snap.Name)
	_ = w.Flush()
}
