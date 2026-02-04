package run

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"titan-sc/api"

	"github.com/spf13/cobra"
)

// State List
const (
	StateCreating = "creating"
	StateCreated  = "created"
	StateDeleted  = "deleted"
	StateStarted  = "started"
	StateStopped  = "stopped"
)

type RunMiddleware struct {
	JSONOutput bool
	Color      bool
	CLIVersion string
	CLIos      string
	API        *api.API
}

func NewRunMiddleware(api *api.API) *RunMiddleware {
	return &RunMiddleware{
		API:        api,
		JSONOutput: false,
		Color:      true,
	}
}

func (run *RunMiddleware) ParseGlobalFlags(cmd *cobra.Command) {
	var err error

	run.JSONOutput, err = cmd.Flags().GetBool("json")
	if err != nil {
		run.JSONOutput = false
	}

	// Color is enabled by default for human output, disabled for JSON
	run.Color = true
	if run.JSONOutput {
		run.Color = false
	} else {
		// Only check --no-color when not in JSON mode
		noColor, _ := cmd.Flags().GetBool("no-color")
		if noColor {
			run.Color = false
		}
	}
}

// ResolveCompanyOID auto-resolves the company OID when only one company exists.
// If company-oid flag is provided, use it. Otherwise, fetch companies and auto-select if only one.
// Returns the company OID and any error encountered.
func (run *RunMiddleware) ResolveCompanyOID(cmd *cobra.Command) (string, error) {
	companyOID, _ := cmd.Flags().GetString("company-oid")

	// If explicitly provided, use it
	if companyOID != "" {
		return companyOID, nil
	}

	// Fetch user's companies
	companies, err := run.API.GetListOfCompanies()
	if err != nil {
		return "", fmt.Errorf("failed to fetch companies: %w", err)
	}

	if len(companies) == 0 {
		return "", fmt.Errorf("no companies found for your account")
	}

	if len(companies) == 1 {
		// Auto-select the only company
		return companies[0].OID, nil
	}

	// Multiple companies - list them and ask user to specify
	var companyList string
	for _, c := range companies {
		companyList += fmt.Sprintf("\n  - %s (%s)", c.Name, c.OID)
	}
	return "", fmt.Errorf("multiple companies found, please specify --company-oid (-c):%s", companyList)
}

// GetDefaultCompanyOID returns the user's default (main) company OID.
// If company-oid flag is provided, use it. Otherwise, fetch user info and return their default company.
// This is useful for super admins who have access to all companies but want to see their own by default.
func (run *RunMiddleware) GetDefaultCompanyOID(cmd *cobra.Command) (string, error) {
	companyOID, _ := cmd.Flags().GetString("company-oid")

	// If explicitly provided, use it
	if companyOID != "" {
		return companyOID, nil
	}

	// Fetch user's info to get their default company
	user, err := run.API.GetUserInfos()
	if err != nil {
		return "", fmt.Errorf("failed to fetch user info: %w", err)
	}

	if user.DefaultCompanyOID == "" {
		return "", fmt.Errorf("no default company found for your account")
	}

	return user.DefaultCompanyOID, nil
}

func (run *RunMiddleware) handleErrorAndGenericOutput(apiReturn *api.Return, err error) {
	// Communication or marshalling error
	if err != nil {
		run.OutputError(err)
		return
	}

	// API parsed error (automatically handle JSON vs string)
	if apiReturn != nil {
		run.printAPIReturn(apiReturn)
		return
	}
}

func (run *RunMiddleware) OutputError(err error) {
	if run.JSONOutput {
		// Output as clean JSON error object
		errObj := map[string]string{"error": err.Error()}
		printAsJson(errObj)
	} else {
		fmt.Printf("%s %s\n", run.Colorize("Error:", "red"), err.Error())
	}
}

func (run *RunMiddleware) printAPIReturn(apiReturn *api.Return) {
	if run.JSONOutput {
		printAsJson(apiReturn)
	} else {
		run.printAPIReturnAsString(apiReturn)
	}
	return
}

func printAsJson(data interface{}) {
	switch data.(type) {
	case []byte:
		// do nothing
		fmt.Println(string(data.([]byte)))
	default:
		dataToPrint, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Printf("{'error': '%s'}\n", err.Error())
			return
		}
		fmt.Println(string(dataToPrint))
	}

}

func (run *RunMiddleware) printAPIReturnAsString(apiReturn *api.Return) {
	if apiReturn.Error() {
		fmt.Printf("%s %s\n", run.Colorize("Error:", "red"), api.ConcatAPIValidationError(apiReturn).Error())
	} else if apiReturn.Success != "" {
		// API v1 success response with message
		fmt.Printf("%s %s\n", run.Colorize("Success:", "green"), apiReturn.Success)
	} else {
		fmt.Printf("%s\n", run.Colorize("Success", "green"))
	}
}

func millisecondsToTime(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}

func keyboardPromptToLower(promptString string) string {
	// Read user input
	fmt.Print(promptString)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSuffix(text, "\n")
	lowerText := strings.ToLower(text)
	return lowerText
}

// GetStateColorized returns a colorized state string if color is enabled
// Works for server, network, subscription, and any other state fields
func GetStateColorized(color bool, state string) string {
	if !color {
		return state
	}
	return colorizeState(state, state)
}

// getStateColor returns the ANSI color code for a given state
// Unified color scheme for all state types across the application
func getStateColor(state string) string {
	switch strings.ToLower(state) {
	// Active/positive states - GREEN
	case StateStarted, StateCreated, "ongoing", "active", "enabled", "connected", "attached", "up":
		return "\033[1;32m"
	// In-progress/transitional states - ORANGE
	case StateCreating, "starting", "stopping", "pending", "processing", "updating", "deleting":
		return "\033[38;5;208m"
	// Warning/paused states - YELLOW
	case StateStopped, "paused", "suspended", "unmanaged":
		return "\033[1;33m"
	// Error/negative states - RED
	case StateDeleted, "cancelled", "canceled", "failed", "error", "disabled", "down":
		return "\033[1;31m"
	default:
		return ""
	}
}

// colorizeState wraps text with state-appropriate color codes
func colorizeState(state, text string) string {
	color := getStateColor(state)
	if color == "" {
		return text
	}
	return color + text + "\033[0m"
}

// getDrpStatusIndicator returns a concise DRP status indicator for list views
// Returns the raw value and a color function (for use with ColColor)
// Shows: - (no DRP), ✓/OK (good), ⏳/PENDING, ⚠/SPLIT, ✗/ERROR
func (run *RunMiddleware) getDrpStatusIndicator(drp *api.Drp) (string, func(string) string) {
	if drp == nil || !drp.Enabled {
		return "-", nil
	}

	// Check status-based indicators
	switch drp.Status {
	case api.DrpStatusOK:
		if run.Color {
			return "OK", ColorFn("green")
		}
		return "OK", nil
	case api.DrpStatusPending:
		if run.Color {
			return "...", ColorFn("orange")
		}
		return "...", nil
	case api.DrpStatusSplitBrain:
		if run.Color {
			return "ERR", ColorFn("red")
		}
		return "ERR", nil
	case api.DrpStatusOff:
		if drp.RequiresAttention || drp.LastError != "" {
			if run.Color {
				return "ERR", ColorFn("red")
			}
			return "ERR", nil
		}
		return "-", nil
	default:
		// Fallback to checking boolean flags for older API responses
		if drp.SplitBrain {
			if run.Color {
				return "ERR", ColorFn("red")
			}
			return "ERR", nil
		}
		if drp.RequiresAttention {
			if run.Color {
				return "ATT", ColorFn("yellow")
			}
			return "ATT", nil
		}
		if drp.PendingOperation != "" {
			if run.Color {
				return "...", ColorFn("orange")
			}
			return "...", nil
		}
		if run.Color {
			return "OK", ColorFn("green")
		}
		return "OK", nil
	}
}

// getDrpStatusText returns human-readable DRP status text
func getDrpStatusText(status int) string {
	switch status {
	case api.DrpStatusOff:
		return "Disabled/Error"
	case api.DrpStatusOK:
		return "Healthy"
	case api.DrpStatusPending:
		return "Pending"
	case api.DrpStatusSplitBrain:
		return "Split-Brain"
	default:
		return "Unknown"
	}
}

// getMirrorStateText returns human-readable mirroring state text
func getMirrorStateText(state int) string {
	switch state {
	case api.MirrorStateUnknown:
		return "Unknown"
	case api.MirrorStateError:
		return "Error"
	case api.MirrorStateSyncing:
		return "Syncing"
	case api.MirrorStateStandby:
		return "Standby"
	case api.MirrorStatePrimary:
		return "Primary"
	default:
		return "Unknown"
	}
}

// mapSiteToPublic converts internal site names (tas/lms) to user-friendly names (main/secondary)
// This is needed because server/network objects still return internal names,
// while DRP-specific endpoints return already-transformed names.
func mapSiteToPublic(internalSite string) string {
	switch strings.ToLower(internalSite) {
	case "tas":
		return "main"
	case "lms":
		return "secondary"
	case "main", "secondary":
		// Already public names (from DRP endpoints)
		return internalSite
	case "":
		return ""
	default:
		// Return as-is if unknown
		return internalSite
	}
}
