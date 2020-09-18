package cmd

import (
	"github.com/spf13/cobra"
	. "titan-sc/api"
)

var company = &cobra.Command{
	Use:     "company",
	Aliases: []string{"co"},
	Short:   "Retrieve information about your companies.",
}

var companyList = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all your companies.",
	Long:    "List all your companies.",
	Run:     API.CompaniesList,
}

var companyShow = &cobra.Command{
	Use:     "show --company-uuid COMPANY_UUID",
	Aliases: []string{"get"},
	Short:   "Show company detail.",
	Long:    "Show detailed information about a company.",
	Run:     API.CompanyDetail,
}

func companyCmdAdd() {
	rootCmd.AddCommand(company)
	company.AddCommand(companyList, companyShow)

	companyShow.Flags().StringP("company-uuid", "c", "", "Set company UUID.")
	_ = companyShow.MarkFlagRequired("company-uuid")
}
