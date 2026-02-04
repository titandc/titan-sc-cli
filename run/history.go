package run

import (
	"fmt"
	"regexp"
	"strings"

	"titan-sc/api"

	"github.com/spf13/cobra"
)

const (
	HistoryNumberMin     = 5
	HistoryNumberMax     = 50
	HistoryNumberDefault = 25
)

func (run *RunMiddleware) EventHistory(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	companyOID, _ := cmd.Flags().GetString("company-oid")
	number, _ := cmd.Flags().GetInt("number")
	offset, _ := cmd.Flags().GetInt("offset")

	// If server specified, use it; otherwise resolve company
	if serverOID != "" {
		setLimits(&number, HistoryNumberMax, HistoryNumberMin)
		strNumber := fmt.Sprintf("%d", number)
		strOffset := fmt.Sprintf("%d", offset)
		run.historyByServer(strNumber, strOffset, serverOID)
		return
	}

	// No server specified, resolve company
	var err error
	if companyOID == "" {
		companyOID, err = run.ResolveCompanyOID(cmd)
		if err != nil {
			run.OutputError(err)
			return
		}
	}

	setLimits(&number, HistoryNumberMax, HistoryNumberMin)
	strNumber := fmt.Sprintf("%d", number)
	strOffset := fmt.Sprintf("%d", offset)
	run.historyByCompany(strNumber, strOffset, companyOID)
}

func setLimits(n *int, max, min int) {
	if *n > max {
		*n = max
	} else if *n < min {
		*n = min
	}
}

func (run *RunMiddleware) historyByCompany(number, offset, companyOID string) {
	events, apiReturn, err := run.API.GetEvents(number, offset, companyOID, api.EventTypeCompany)
	if err != nil {
		run.OutputError(err)
		return
	}
	if apiReturn != nil && apiReturn.Error() {
		err = api.ConcatAPIValidationError(apiReturn)
		run.OutputError(err)
		return
	}
	run.printEvents(events)
}

func (run *RunMiddleware) historyByServer(number, offset, serverOID string) {
	events, apiReturn, err := run.API.GetEvents(number, offset, serverOID, api.EventTypeServer)
	if err != nil {
		run.OutputError(err)
		return
	}
	if apiReturn != nil && apiReturn.Error() {
		err := api.ConcatAPIValidationError(apiReturn)
		run.OutputError(err)
		return
	}
	run.printEvents(events)
}

func (run *RunMiddleware) printEvents(events []api.Event) {
	// Sanitize field data before output (fix backend Go fmt errors)
	for i := range events {
		for j := range events[i].Fields {
			if events[i].Fields[j].Data != nil {
				sanitized := sanitizeFieldData(*events[i].Fields[j].Data)
				events[i].Fields[j].Data = &sanitized
			}
		}
	}

	if run.JSONOutput {
		printAsJson(events)
	} else {
		run.printEventsTable(events)
	}
}

func (run *RunMiddleware) printEventsTable(events []api.Event) {
	table := NewTable("TIME", "TYPE", "ACTION", "STATUS", "TARGET", "SERVER", "USER")
	table.SetNoColor(!run.Color)

	for _, event := range events {
		// Build event type string (e.g., "network/vswitch")
		eventType := event.Type
		if event.SubType != "" {
			eventType = fmt.Sprintf("%s/%s", event.Type, event.SubType)
		}

		// Action (e.g., "detach", "created")
		action := event.Action

		// Status with icon
		status := run.getStatusWithIcon(event.Status)

		// Target info
		target := ""
		if event.TargetName != nil {
			target = *event.TargetName
		}

		// Server info
		server := ""
		if event.ServerName != nil {
			server = *event.ServerName
		}

		// User info
		user := ""
		if event.UserEmail != nil {
			user = *event.UserEmail
		}

		table.AddRow(
			ColColor(event.Timestamp, ColorFn("dim")),
			ColColor(eventType, ColorFn("cyan")),
			Col(action),
			ColColor(status, StateColorFn(event.Status)),
			ColColor(target, ColorFn("cyan")),
			ColColor(server, ColorFn("cyan")),
			Col(user),
		)
	}

	table.Print()
}

func (run *RunMiddleware) collectEventDetails(fields []api.EventAdditionalFields) string {
	var parts []string
	for _, field := range fields {
		if field.Data != nil && *field.Data != "" {
			// Skip values that look like Unix timestamps (large numbers that result in 1970 dates)
			if !looksLikeUnixTimestamp(*field.Data) {
				parts = append(parts, *field.Data)
			}
		}
		if field.OldValue != nil && field.NewValue != nil {
			old := formatFieldValue(*field.OldValue)
			new := formatFieldValue(*field.NewValue)
			if old != "" || new != "" {
				parts = append(parts, fmt.Sprintf("%s→%s", old, new))
			}
		} else if field.NewValue != nil {
			val := formatFieldValue(*field.NewValue)
			if val != "" {
				name := ""
				if field.Name != nil {
					name = *field.Name + ": "
				}
				parts = append(parts, name+val)
			}
		}
	}
	return strings.Join(parts, ", ")
}

// looksLikeUnixTimestamp checks if a string looks like a Unix timestamp or 1970-era date
func looksLikeUnixTimestamp(s string) bool {
	// Check for 1970 dates (common for unset/zero timestamps)
	if strings.HasPrefix(s, "1970-") {
		return true
	}
	return false
}

// formatFieldValue formats a field value, filtering out bogus timestamps
func formatFieldValue(s string) string {
	// Filter out 1970 dates and zero values
	if s == "0" || strings.HasPrefix(s, "1970-") {
		return ""
	}
	return s
}

func (run *RunMiddleware) getStatusWithIcon(status string) string {
	switch status {
	case "success":
		return "✓ success"
	case "request":
		return "→ request"
	case "error", "failed":
		return "✗ " + status
	case "alert":
		return "⚠ alert"
	default:
		return status
	}
}

// sanitizeFieldData fixes corrupted strings from backend Go fmt errors
// e.g., ">85%!D(MISSING)uration: 10mn" -> ">85% - Duration: 10mn"
func sanitizeFieldData(s string) string {
	// Pattern: %!X(MISSING) where X is a letter
	re := regexp.MustCompile(`%!([A-Za-z])\(MISSING\)`)
	return re.ReplaceAllString(s, "% - $1")
}

func (run *RunMiddleware) printEventHuman(event api.Event) {
	// Build event title based on type
	title := run.buildEventTitle(event)
	fmt.Printf("%s\n", run.Colorize(title, "cyan"))

	// Timestamp
	fmt.Printf("  Time:    %s\n", run.Colorize(event.Timestamp, "dim"))

	// Server info (if present)
	if event.ServerName != nil {
		fmt.Printf("  Server:  %s", run.Colorize(*event.ServerName, "cyan"))
		if event.ServerOID != nil {
			fmt.Printf(" %s", run.Colorize(fmt.Sprintf("(%s)", *event.ServerOID), "dim"))
		}
		fmt.Println()
	}

	// Target info (for snapshots, etc)
	if event.TargetName != nil {
		fmt.Printf("  Target:  %s", run.Colorize(*event.TargetName, "cyan"))
		if event.TargetOID != nil {
			fmt.Printf(" %s", run.Colorize(fmt.Sprintf("(%s)", *event.TargetOID), "dim"))
		}
		fmt.Println()
	}

	// User info (if present)
	if event.UserEmail != nil {
		fmt.Printf("  User:    %s", run.Colorize(*event.UserEmail, "cyan"))
		if event.UserIP != nil {
			fmt.Printf(" %s", run.Colorize(fmt.Sprintf("(IP: %s)", *event.UserIP), "dim"))
		}
		fmt.Println()
	}

	// Fields - display data/changes
	for _, field := range event.Fields {
		run.printEventField(field)
	}
}

func (run *RunMiddleware) buildEventTitle(event api.Event) string {
	// For metric alerts, show: "⚠ CPU Alert" or "⚠ RAM Alert"
	if event.Type == "metric" {
		metricName := capitalizeFirst(event.SubType)
		return fmt.Sprintf("%s %s Alert", run.Colorize("⚠", "yellow"), metricName)
	}

	// For storage operations: "Snapshot Created", "Snapshot Deleted", etc.
	if event.Type == "storage" {
		action := capitalizeFirst(event.Action)
		subType := capitalizeFirst(event.SubType)
		statusIcon := run.getStatusIcon(event.Status)
		return fmt.Sprintf("%s %s %s", statusIcon, subType, action)
	}

	// Generic fallback
	return fmt.Sprintf("%s/%s: %s (%s)", event.Type, event.SubType, event.Action, event.Status)
}

func (run *RunMiddleware) printEventField(field api.EventAdditionalFields) {
	if field.Data != nil && *field.Data != "" {
		// For metric data like ">95%!D(MISSING)uration: 60mn" or ">95% - Duration: 60mn"
		fmt.Printf("  Details: %s\n", *field.Data)
	}
	if field.OldValue != nil && field.NewValue != nil {
		fmt.Printf("  Changed: %s -> %s\n", *field.OldValue, *field.NewValue)
	} else if field.NewValue != nil {
		name := "Value"
		if field.Name != nil {
			name = capitalizeFirst(*field.Name)
		}
		fmt.Printf("  %s:   %s\n", name, *field.NewValue)
	}
}

func (run *RunMiddleware) getStatusIcon(status string) string {
	switch status {
	case "success":
		return run.Colorize("✓", "green")
	case "request":
		return run.Colorize("→", "cyan")
	case "error", "failed":
		return run.Colorize("✗", "red")
	case "alert":
		return run.Colorize("⚠", "yellow")
	default:
		return "•"
	}
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	// Handle common abbreviations
	switch s {
	case "cpu":
		return "CPU"
	case "ram":
		return "RAM"
	case "disk":
		return "Disk"
	case "snapshot":
		return "Snapshot"
	}
	// Default: capitalize first letter
	if len(s) == 1 {
		return string(s[0] - 32)
	}
	return string(s[0]-32) + s[1:]
}
