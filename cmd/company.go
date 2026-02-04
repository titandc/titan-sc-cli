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
		Use:     "show [--company-oid COMPANY_OID]",
		Aliases: []string{"get"},
		Short:   "Show company detail.",
		Long: `Show detailed information about a company.

If you have only one company, --company-oid is optional and will be auto-detected.
If you have multiple companies, you must specify --company-oid.`,
		Run: cmd.runMiddleware.CompanyDetail,
	}

	cmd.RootCommand.AddCommand(company)
	company.AddCommand(companyList, companyShow)

	companyShow.Flags().StringP("company-oid", "c", "", "Company OID (auto-detected if you have only one company).")
}
