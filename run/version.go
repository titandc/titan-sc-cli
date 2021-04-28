package run

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (run *RunMiddleware) VersionAPI(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	version, err := run.API.GetVersion()
	if err != nil {
		run.OutputError(err)
		return
	}
	if !run.HumanReadable {
		printAsJson(version)
	} else {
		fmt.Println("Titan API version:", version.Version)
	}
}
