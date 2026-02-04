package run

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// NormalizeVersion strips any "v" prefix from a version string for consistent output
func NormalizeVersion(version string) string {
	return strings.TrimPrefix(version, "v")
}

// PrintVersion outputs a version in a consistent format (human or JSON)
func (run *RunMiddleware) PrintVersion(label, version string) {
	normalized := NormalizeVersion(version)

	if run.JSONOutput {
		data, _ := json.MarshalIndent(struct {
			Version string `json:"version"`
		}{Version: normalized}, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("%s: %s\n", label, run.Colorize(normalized, "green"))
	}
}

func (run *RunMiddleware) VersionAPI(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	version, err := run.API.GetVersion()
	if err != nil {
		run.OutputError(err)
		return
	}

	run.PrintVersion("Titan API version", version.Version)
}
