package run

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

// displayWidth calculates the visual display width of a string.
// Uses go-runewidth for proper Unicode width handling including:
// - Multi-byte UTF-8 characters
// - Wide characters (CJK, emojis)
// - Combining characters and emoji modifiers
// - ZWJ sequences
func displayWidth(s string) int {
	return runewidth.StringWidth(s)
}

// padRight pads a string to the specified display width with spaces
func padRight(s string, targetWidth int) string {
	return runewidth.FillRight(s, targetWidth)
}

/*
Color System Guidelines for Human Output:

SEMANTIC COLORS (use consistently across all commands):
  - cyan:    Names, identifiers, section headers, URLs, emails, usernames
  - blue:    OIDs, unique identifiers, references
  - green:   Success messages, counts/totals, prices/amounts, positive booleans (true/yes/enabled)
  - yellow:  Warnings, deadlines, negative booleans (false/no), paused states, IP addresses
  - orange:  Transitional states (pending, updating, starting, stopping, creating, deleting)
  - red:     Errors, failed states, disabled features, deleted items
  - dim:     Timestamps, secondary info, comments
  - magenta: High privilege roles (SUPER_ADMINISTRATOR)

STATE COLORS (via StateColorFn/ColState):
  - GREEN:  started, created, ongoing, active, enabled, connected, attached, up
  - ORANGE: creating, starting, stopping, pending, processing, updating, deleting (transitional)
  - YELLOW: stopped, paused, suspended, unmanaged (warning)
  - RED:    deleted, cancelled, failed, error, disabled, down (negative)

ROLE COLORS (via getRoleColorFn/getRoleColor):
  - magenta: SUPER_ADMINISTRATOR
  - yellow:  ADMINISTRATOR
  - green:   USER
  - dim:     unknown/other

SEMANTIC COLUMN HELPERS (use for consistent styling):
  - ColName(s)      -> cyan    : Names, labels, identifiers
  - ColOID(s)       -> blue    : Object IDs, unique references
  - ColIP(s)        -> yellow  : IP addresses (v4/v6)
  - ColState(s)     -> varies  : State values (auto-colored based on state)
  - ColCount(s)     -> green   : Counts, numbers, amounts
  - ColTimestamp(s) -> dim     : Dates, times, timestamps
*/

// TableColumn defines a column with optional colorization
type TableColumn struct {
	Value   string
	ColorFn func(string) string // Optional: colorize function for this cell
}

// Table handles dynamic column width calculation and ANSI-safe printing
type Table struct {
	headers []string
	rows    [][]TableColumn
	widths  []int
	noColor bool // When true, ignore all ColorFn
}

// NewTable creates a new table with the given headers
func NewTable(headers ...string) *Table {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = displayWidth(h)
	}
	return &Table{
		headers: headers,
		rows:    make([][]TableColumn, 0),
		widths:  widths,
		noColor: false,
	}
}

// SetNoColor disables color output for the table
func (t *Table) SetNoColor(noColor bool) {
	t.noColor = noColor
}

// AddRow adds a row of columns to the table
func (t *Table) AddRow(cols ...TableColumn) {
	if len(cols) != len(t.headers) {
		return // Silently ignore mismatched rows
	}

	// Update max widths based on display width (handles emojis correctly)
	for i, col := range cols {
		colWidth := displayWidth(col.Value)
		if colWidth > t.widths[i] {
			t.widths[i] = colWidth
		}
	}

	t.rows = append(t.rows, cols)
}

// Print outputs the table with proper alignment
func (t *Table) Print() {
	// Print header with proper padding
	headerArgs := make([]interface{}, len(t.headers))
	for i, h := range t.headers {
		headerArgs[i] = padRight(h, t.widths[i])
	}
	headerFormat := strings.Repeat("%s  ", len(t.headers)-1) + "%s\n"
	fmt.Printf(headerFormat, headerArgs...)

	// Print rows
	for _, row := range t.rows {
		args := make([]interface{}, len(row))
		for i, col := range row {
			// Pad value using display width (handles emojis correctly)
			padded := padRight(col.Value, t.widths[i])
			if col.ColorFn != nil && !t.noColor {
				args[i] = col.ColorFn(padded)
			} else {
				args[i] = padded
			}
		}
		// Use simple format since padding is pre-applied
		simpleFormat := strings.Repeat("%s  ", len(row)-1) + "%s\n"
		fmt.Printf(simpleFormat, args...)
	}
}

// Col is a helper to create a simple column without color
func Col(value string) TableColumn {
	return TableColumn{Value: value}
}

// ColColor creates a column with a color function
func ColColor(value string, colorFn func(string) string) TableColumn {
	return TableColumn{Value: value, ColorFn: colorFn}
}

// StateColorFn returns a color function for the given state
func StateColorFn(state string) func(string) string {
	color := getStateColor(state)
	if color == "" {
		return nil
	}
	return func(s string) string {
		return color + s + "\033[0m"
	}
}

// ColorFn returns a generic color function for the given color name
func ColorFn(colorName string) func(string) string {
	code := colorCode(colorName)
	if code == "" {
		return nil
	}
	return func(s string) string {
		return code + s + "\033[0m"
	}
}

// colorCode returns the ANSI escape code for a color name
func colorCode(colorName string) string {
	switch colorName {
	case "red":
		return "\033[1;31m"
	case "green":
		return "\033[1;32m"
	case "yellow":
		return "\033[1;33m"
	case "blue":
		return "\033[1;34m"
	case "magenta":
		return "\033[1;35m"
	case "cyan":
		return "\033[1;36m"
	case "white":
		return "\033[1;37m"
	case "dim":
		return "\033[2m"
	default:
		return ""
	}
}

// --- Semantic Color Functions ---
// Use these for consistent coloring across all commands

// ColIP creates a column for IP addresses (yellow - stands out, network related)
func ColIP(ip string) TableColumn {
	return TableColumn{Value: ip, ColorFn: ColorFn("yellow")}
}

// ColOID creates a column for OIDs (blue - distinct identifier)
func ColOID(oid string) TableColumn {
	return TableColumn{Value: oid, ColorFn: ColorFn("blue")}
}

// ColName creates a column for names/identifiers (cyan - primary identifier)
func ColName(name string) TableColumn {
	return TableColumn{Value: name, ColorFn: ColorFn("cyan")}
}

// ColState creates a column with state-based coloring
func ColState(state string) TableColumn {
	return TableColumn{Value: state, ColorFn: StateColorFn(state)}
}

// ColCount creates a column for counts/numbers (green - positive values)
func ColCount(count string) TableColumn {
	return TableColumn{Value: count, ColorFn: ColorFn("green")}
}

// ColTimestamp creates a column for timestamps (dim - secondary info)
func ColTimestamp(ts string) TableColumn {
	return TableColumn{Value: ts, ColorFn: ColorFn("dim")}
}
