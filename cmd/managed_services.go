package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) ManagedServicesCmdAdd() {

	managedServices := &cobra.Command{
		Use:   "managed-services --company-uuid COMPANY_UUID",
		Short: "Enable managed services.",
		Long:  "Enable managed services.",
		Run:   cmd.runMiddleware.ManagedServices,
	}

	cmd.RootCommand.AddCommand(managedServices)

	managedServices.Flags().StringP("company-uuid", "s", "", "Set company UUID.")
	_ = managedServices.MarkFlagRequired("company-uuid")
}
