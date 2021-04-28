package run

import (
	"github.com/spf13/cobra"
)

func (run *RunMiddleware) ManagedServices(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	companyUUID, _ := cmd.Flags().GetString("company-uuid")

	run.handleErrorAndGenericOutput(run.API.PostManagedServices(companyUUID))
}
