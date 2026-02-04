package run

import (
	"fmt"
	"strings"
	"titan-sc/api"

	"github.com/spf13/cobra"
)

// ServerDrpStatus shows the DRP status for a server
func (run *RunMiddleware) ServerDrpStatus(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")

	status, err := run.API.GetDrpStatus(serverOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(status)
	} else {
		run.printDrpStatus(status)
	}
}

func (run *RunMiddleware) printDrpStatus(status *api.DrpStatus) {
	fmt.Printf("%s\n", run.Colorize("DRP Status:", "cyan"))
	fmt.Printf("  Server: %s (%s)\n", run.Colorize(status.ServerName, "cyan"), status.ServerOID)

	// Enabled status
	if status.Enabled {
		fmt.Printf("  Enabled: %s\n", run.Colorize("Yes", "green"))
	} else {
		fmt.Printf("  Enabled: %s\n", run.Colorize("No", "yellow"))
		return
	}

	// Overall status
	statusText := getDrpStatusText(status.Status)
	if status.Status == api.DrpStatusOK {
		statusText = run.Colorize(statusText, "green")
	} else if status.Status == api.DrpStatusSplitBrain {
		statusText = run.Colorize(statusText, "red")
	} else if status.Status == api.DrpStatusPending {
		statusText = run.Colorize(statusText, "yellow")
	} else if status.Status == api.DrpStatusOff {
		statusText = run.Colorize(statusText, "red")
	}
	fmt.Printf("  Status: %s\n", statusText)

	// Active site
	if status.ActiveSite != "" {
		fmt.Printf("  Active Site: %s\n", run.Colorize(mapSiteToPublic(status.ActiveSite), "cyan"))
	}

	// Interval
	if status.Interval > 0 {
		fmt.Printf("  Sync Interval: %d hours\n", status.Interval)
	}

	// Start time
	if status.StartTime != nil && status.StartTime.IsSet() {
		fmt.Printf("  DRP Start Time: %s\n", FlexTimestampFormat(status.StartTime))
	}

	// Last failover info
	if status.LastFailoverAt != nil && status.LastFailoverAt.IsSet() {
		fmt.Printf("  Last Failover: %s", FlexTimestampFormat(status.LastFailoverAt))
		if status.LastFailoverType != "" {
			fmt.Printf(" (%s)", status.LastFailoverType)
		}
		fmt.Printf("\n")
	}

	// Last resync
	if status.LastResyncAt != nil && status.LastResyncAt.IsSet() {
		fmt.Printf("  Last Resync: %s\n", FlexTimestampFormat(status.LastResyncAt))
	}

	// Last operation result
	if status.LastOperationResult != "" {
		fmt.Printf("  Last Operation Result: %s\n", status.LastOperationResult)
	}

	// Pending operation
	if status.PendingOperation != "" {
		fmt.Printf("  %s %s\n", run.Colorize("Pending Operation:", "yellow"), status.PendingOperation)
		if status.PendingOperationAt != nil && status.PendingOperationAt.IsSet() {
			fmt.Printf("    Started: %s\n", FlexTimestampFormat(status.PendingOperationAt))
		}
		if status.PendingOperationBy != "" {
			fmt.Printf("    By: %s\n", status.PendingOperationBy)
		}
	}

	// IPs status
	if status.IPs != nil && len(status.IPs) > 0 {
		fmt.Printf("  IPs:\n")
		for _, ip := range status.IPs {
			fmt.Printf("    %s (v%d): Site %s", ip.Address, ip.Version, mapSiteToPublic(ip.CurrentSite))
			if ip.MACAddress != "" {
				fmt.Printf(" [MAC: %s]", ip.MACAddress)
			}
			fmt.Printf("\n")
			if ip.LastSwitchAt != nil && ip.LastSwitchAt.IsSet() {
				fmt.Printf("      Last switch: %s\n", FlexTimestampFormat(ip.LastSwitchAt))
			}
			if ip.LastSwitchError != "" {
				// Map site names in error message for consistency
				errMsg := ip.LastSwitchError
				errMsg = strings.ReplaceAll(errMsg, " tas", " "+mapSiteToPublic("tas"))
				errMsg = strings.ReplaceAll(errMsg, " lms", " "+mapSiteToPublic("lms"))
				fmt.Printf("      %s\n", run.Colorize("Last switch error: "+errMsg, "red"))
			}
		}
	}

	// Split-brain warning
	if status.SplitBrain {
		fmt.Printf("\n  %s\n", run.Colorize("⚠ WARNING: Split-Brain detected!", "red"))
		fmt.Printf("  %s\n", run.Colorize("Manual intervention required. Use 'drp resync' to resolve.", "red"))
	}

	// Requires attention
	if status.RequiresAttention && !status.SplitBrain {
		fmt.Printf("\n  %s\n", run.Colorize("⚠ DRP requires attention", "yellow"))
	}

	// Last error
	if status.LastError != "" {
		// Map site names in error message for consistency
		errMsg := status.LastError
		errMsg = strings.ReplaceAll(errMsg, "'tas'", "'"+mapSiteToPublic("tas")+"'")
		errMsg = strings.ReplaceAll(errMsg, "'lms'", "'"+mapSiteToPublic("lms")+"'")
		fmt.Printf("  Last Error: %s\n", run.Colorize(errMsg, "red"))
	}
}

// ServerDrpFailoverSoft performs a soft failover
func (run *RunMiddleware) ServerDrpFailoverSoft(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")

	result, err := run.API.DrpFailoverSoft(serverOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(result)
	} else {
		run.printDrpOperationResult("Soft Failover", result)
	}
}

// ServerDrpFailoverHard performs a hard failover (DANGEROUS)
func (run *RunMiddleware) ServerDrpFailoverHard(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	targetSite, _ := cmd.Flags().GetString("target-site")
	confirmed, _ := cmd.Flags().GetBool("yes-i-understand-i-will-lose-data")

	// Check for missing required flags
	targetSite = strings.ToLower(strings.TrimSpace(targetSite))
	if targetSite == "" || !confirmed {
		fmt.Printf("%s Hard failover will cause data loss.\n\n", run.Colorize("⚠ DANGEROUS:", "red"))
		fmt.Printf("%s\n", run.Colorize("Required flags:", "yellow"))
		if targetSite == "" {
			fmt.Printf("  --target-site <main|secondary>\n")
		}
		if !confirmed {
			fmt.Printf("  --yes-i-understand-i-will-lose-data\n")
		}
		fmt.Printf("\nUse -h for more information.\n")
		return
	}

	// Validate target site value
	if targetSite != "main" && targetSite != "secondary" {
		run.OutputError(fmt.Errorf("invalid target site '%s': must be 'main' or 'secondary'", targetSite))
		return
	}

	// Print warning
	fmt.Printf("%s\n", run.Colorize("⚠ Performing HARD FAILOVER to "+targetSite, "yellow"))
	fmt.Printf("This operation may cause data loss...\n\n")

	result, err := run.API.DrpFailoverHard(serverOID, targetSite)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(result)
	} else {
		run.printDrpOperationResult("Hard Failover", result)
	}
}

// ServerDrpResync resynchronizes DRP after split-brain (DANGEROUS)
func (run *RunMiddleware) ServerDrpResync(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	authoritativeSite, _ := cmd.Flags().GetString("authoritative-site")
	confirmed, _ := cmd.Flags().GetBool("yes-i-understand-i-will-lose-data")

	// Check for missing required flags
	authoritativeSite = strings.ToLower(strings.TrimSpace(authoritativeSite))
	if authoritativeSite == "" || !confirmed {
		fmt.Printf("%s Resync will OVERWRITE all data on the non-authoritative site.\n\n", run.Colorize("⚠ DANGEROUS:", "red"))
		fmt.Printf("%s\n", run.Colorize("Required flags:", "yellow"))
		if authoritativeSite == "" {
			fmt.Printf("  --authoritative-site <main|secondary>\n")
		}
		if !confirmed {
			fmt.Printf("  --yes-i-understand-i-will-lose-data\n")
		}
		fmt.Printf("\nUse -h for more information.\n")
		return
	}

	// Validate authoritative site value
	if authoritativeSite != "main" && authoritativeSite != "secondary" {
		run.OutputError(fmt.Errorf("invalid authoritative site '%s': must be 'main' or 'secondary'", authoritativeSite))
		return
	}

	// Determine which site loses data
	losingDataSite := "secondary"
	if authoritativeSite == "secondary" {
		losingDataSite = "main"
	}

	// Print warning
	fmt.Printf("%s\n", run.Colorize("⚠ Performing DRP RESYNC", "yellow"))
	fmt.Printf("Authoritative site: %s (data will be preserved)\n", run.Colorize(authoritativeSite, "green"))
	fmt.Printf("Non-authoritative site: %s (data will be OVERWRITTEN)\n\n", run.Colorize(losingDataSite, "red"))

	result, err := run.API.DrpResync(serverOID, authoritativeSite)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(result)
	} else {
		run.printDrpOperationResult("Resync", result)
	}
}

func (run *RunMiddleware) printDrpOperationResult(operation string, result *api.DrpOperationResult) {
	if result.Success {
		fmt.Printf("%s %s\n", run.Colorize("✓", "green"), run.Colorize(operation+" initiated successfully", "green"))
	} else {
		fmt.Printf("%s %s\n", run.Colorize("✗", "red"), run.Colorize(operation+" failed", "red"))
	}

	if result.Operation != "" {
		fmt.Printf("  Operation: %s\n", result.Operation)
	}

	if result.Message != "" {
		fmt.Printf("  Message: %s\n", result.Message)
	}

	if result.SourceSite != "" || result.TargetSite != "" {
		if result.SourceSite != "" {
			fmt.Printf("  Source Site: %s\n", result.SourceSite)
		}
		if result.TargetSite != "" {
			fmt.Printf("  Target Site: %s\n", result.TargetSite)
		}
	}

	if !result.Success && result.Error != "" {
		fmt.Printf("  Error: %s\n", run.Colorize(result.Error, "red"))
	}
}

// NetworkDrpEnable enables DRP for a network
func (run *RunMiddleware) NetworkDrpEnable(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkOID, _ := cmd.Flags().GetString("network-oid")

	network, err := run.API.DrpNetworkEnable(networkOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(network)
	} else {
		fmt.Printf("%s Network DRP enabled successfully\n", run.Colorize("✓", "green"))
		fmt.Printf("  Network: %s (%s)\n", run.Colorize(network.Name, "cyan"), network.OID)
		if network.Drp != nil {
			fmt.Printf("  DRP Enabled: %s\n", run.Colorize("Yes", "green"))
			if network.Drp.Site != "" {
				fmt.Printf("  Target Site: %s\n", mapSiteToPublic(network.Drp.Site))
			}
		}
	}
}

// NetworkDrpDisable disables DRP for a network
func (run *RunMiddleware) NetworkDrpDisable(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkOID, _ := cmd.Flags().GetString("network-oid")
	confirmed, _ := cmd.Flags().GetBool("yes-i-understand-network-will-be-unavailable")

	// Double-check confirmation flag
	if !confirmed {
		fmt.Printf("%s\n", run.Colorize("⚠ WARNING", "yellow"))
		fmt.Printf("Disabling DRP will stop network replication to the target site.\n")
		fmt.Printf("If servers fail over to the secondary site, they will lose\n")
		fmt.Printf("private network connectivity until DRP is re-enabled.\n")
		fmt.Printf("Use --yes-i-understand-network-will-be-unavailable to confirm.\n")
		return
	}

	fmt.Printf("%s\n", run.Colorize("⚠ Disabling network DRP...", "yellow"))

	network, err := run.API.DrpNetworkDisable(networkOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(network)
	} else {
		fmt.Printf("%s Network DRP disabled successfully\n", run.Colorize("✓", "green"))
		fmt.Printf("  Network: %s (%s)\n", run.Colorize(network.Name, "cyan"), network.OID)
		if network.Drp != nil {
			fmt.Printf("  DRP Enabled: %s\n", run.Colorize("No", "yellow"))
		}
	}
}
