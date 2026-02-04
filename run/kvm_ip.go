package run

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"titan-sc/api"

	"github.com/spf13/cobra"
)

const noVNCBaseURL = "https://novnc.titandc.io/vnc.html?uri="

// kvmJSONOutput is a wrapper for JSON output that includes the console URL
type kvmJSONOutput struct {
	OID       string `json:"oid"`
	ServerOID string `json:"server_oid"`
	State     string `json:"state,omitempty"`
	URL       string `json:"url,omitempty"`
	WebURL    string `json:"web_url,omitempty"`
	CreatedAt *int64 `json:"created_at,omitempty"`
	UpdatedAt *int64 `json:"updated_at,omitempty"`
	Deadline  int64  `json:"deadline,omitempty"`
}

func buildKvmJSONOutput(kvm *api.KvmIPView) kvmJSONOutput {
	out := kvmJSONOutput{
		OID:       kvm.OID,
		ServerOID: kvm.ServerOID,
		State:     kvm.State,
		URL:       kvm.URL,
		Deadline:  kvm.Deadline,
		CreatedAt: kvm.CreatedAt,
		UpdatedAt: kvm.UpdatedAt,
	}
	if kvm.URL != "" {
		out.WebURL = noVNCBaseURL + url.QueryEscape(kvm.URL)
	}
	return out
}

func (run *RunMiddleware) KVMIPGetInfos(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	kvmOID, _ := cmd.Flags().GetString("kvm-oid")

	var kvm *api.KvmIPView
	var serverName string

	if serverOID != "" {
		// Get KVM info via server detail endpoint
		server, apiReturn, err := run.API.GetServerOID(serverOID)
		if err != nil {
			run.OutputError(err)
			return
		}
		if apiReturn != nil && apiReturn.Error() {
			run.OutputError(api.ConcatAPIValidationError(apiReturn))
			return
		}
		if server.KVM == nil {
			run.OutputError(fmt.Errorf("no active KVM session for this server"))
			return
		}
		kvm = server.KVM
		serverName = server.Name
	} else {
		// Get KVM info directly via KVM OID
		path := fmt.Sprintf("/console/kvmip/%s", kvmOID)
		body, apiReturn, err := run.API.SendRequestToAPI(api.HTTPGet, path, nil)
		if err != nil {
			run.OutputError(err)
			return
		}
		if apiReturn != nil && apiReturn.Error() {
			run.OutputError(api.ConcatAPIValidationError(apiReturn))
			return
		}
		var kvmInfo api.KvmIPView
		if err := json.Unmarshal(body, &kvmInfo); err != nil {
			run.OutputError(err)
			return
		}
		kvm = &kvmInfo
		serverOID = kvm.ServerOID

		// Fetch server info to get the state (KVM direct endpoint doesn't return state)
		if serverOID != "" && kvm.State == "" {
			server, _, _ := run.API.GetServerOID(serverOID)
			if server != nil {
				serverName = server.Name
				if server.KVM != nil && server.KVM.State != "" {
					kvm.State = server.KVM.State
				}
			}
		}
	}

	if run.JSONOutput {
		printAsJson(buildKvmJSONOutput(kvm))
	} else {
		run.printKvmInfo(serverOID, serverName, kvm)
	}
}

func formatTimestamp(ts int64) string {
	if ts == 0 {
		return "-"
	}
	// API returns milliseconds, convert to seconds
	t := time.UnixMilli(ts)
	return t.Format("2006-01-02 15:04:05")
}

func (run *RunMiddleware) printKvmInfo(serverOID, serverName string, kvm *api.KvmIPView) {
	fmt.Printf("%s\n", run.Colorize("KVM Session:", "cyan"))
	if serverName != "" {
		fmt.Printf("  Server:     %s %s\n", run.Colorize(serverName, "cyan"), run.Colorize(fmt.Sprintf("(%s)", serverOID), "dim"))
	} else {
		fmt.Printf("  Server:     %s\n", serverOID)
	}
	fmt.Printf("  KVM OID:    %s\n", kvm.OID)
	fmt.Printf("  State:      %s\n", GetStateColorized(run.Color, kvm.State))
	if kvm.URL != "" {
		noVNCURL := noVNCBaseURL + url.QueryEscape(kvm.URL)
		fmt.Printf("  Web URL:    %s\n", run.Colorize(noVNCURL, "cyan"))
	}
	if kvm.CreatedAt != nil {
		fmt.Printf("  Created:    %s\n", run.Colorize(formatTimestamp(*kvm.CreatedAt), "dim"))
	}
	if kvm.UpdatedAt != nil {
		fmt.Printf("  Updated:    %s\n", run.Colorize(formatTimestamp(*kvm.UpdatedAt), "dim"))
	}
	if kvm.Deadline > 0 {
		fmt.Printf("  Deadline:   %s\n", run.Colorize(formatTimestamp(kvm.Deadline), "yellow"))
	}
}

func (run *RunMiddleware) KVMIPStart(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")

	req := &api.KvmIPRequest{
		ServerOID: serverOID,
	}

	body, apiReturn, err := run.API.SendRequestToAPI(api.HTTPPost, "/console/kvmip/", req)
	if err != nil {
		run.OutputError(err)
		return
	}
	if apiReturn != nil && apiReturn.Error() {
		run.OutputError(api.ConcatAPIValidationError(apiReturn))
		return
	}

	// Parse the KVM response
	var kvmip api.KvmIP
	if err := json.Unmarshal(body, &kvmip); err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(kvmip)
	} else {
		fmt.Printf("%s\n", run.Colorize("KVM Session Started:", "green"))
		fmt.Printf("  Server OID: %s\n", serverOID)
		fmt.Printf("  KVM OID:    %s\n", kvmip.OID)
		if kvmip.URL != "" {
			noVNCURL := noVNCBaseURL + url.QueryEscape(kvmip.URL)
			fmt.Printf("  Web URL:    %s\n", run.Colorize(noVNCURL, "cyan"))
		}
		if kvmip.Deadline > 0 {
			fmt.Printf("  Deadline:   %s\n", run.Colorize(formatTimestamp(kvmip.Deadline), "yellow"))
		}
	}
}

func (run *RunMiddleware) KVMIPStop(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")

	path := fmt.Sprintf("/console/kvmip/%s", serverOID)
	_, apiReturn, err := run.API.SendRequestToAPI(api.HTTPDelete, path, nil)
	if err != nil {
		run.OutputError(err)
		return
	}
	if apiReturn != nil && apiReturn.Error() {
		run.OutputError(api.ConcatAPIValidationError(apiReturn))
		return
	}

	if run.JSONOutput {
		fmt.Println(`{"status": "success", "message": "KVM session stopped"}`)
	} else {
		fmt.Printf("%s for server %s\n", run.Colorize("KVM session stopped", "green"), run.Colorize(serverOID, "cyan"))
	}
}
