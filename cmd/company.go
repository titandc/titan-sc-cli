package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) CompanyCmdAdd() {

	company := &cobra.Command{
		Use:     "company",
		Aliases: []string{"co"},
		Short:   "Retrieve information about your companies.",
	}

	companyList := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all your companies.",
		Long:    "List all your companies.",
		Run:     cmd.runMiddleware.CompaniesList,
	}

	companyShow := &cobra.Command{
		Use:     "show --company-uuid COMPANY_UUID",
		Aliases: []string{"get"},
		Short:   "Show company detail.",
		Long:    "Show detailed information about a company.",
		Run:     cmd.runMiddleware.CompanyDetail,
	}

	cmd.RootCommand.AddCommand(company)
	company.AddCommand(companyList, companyShow)

	companyShow.Flags().StringP("company-uuid", "c", "", "Set company UUID.")
	_ = companyShow.MarkFlagRequired("company-uuid")
}
