package cmd

import (
    "github.com/spf13/cobra"
    . "titan-sc/api"
)

var company = &cobra.Command{
    Use: "company",
    Aliases: []string{"co"},
    Short: "Retrieve information about your companies.",
}

var companyList = &cobra.Command{
    Use: "list",
    Aliases: []string{"ls"},
    Short: "List all your companies.",
    Long: "List all your compagnies.",
    Run: API.CompaniesList,
}

var companyShow = &cobra.Command{
    Use: "show [company_uuid]",
    Aliases: []string{"get"},
    Short: "Show company detail.",
    Long: "Show company detail (need company UUID).",
    Args: cmdNeed1UUID,
    Run: API.CompanyDetail,
}

func companyCmdAdd() {
    rootCmd.AddCommand(company)
    company.AddCommand(companyList, companyShow)
}
