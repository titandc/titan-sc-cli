package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"titan-sc/api"
)

/*
 *
 *
 ******************
 * Company function
 ******************
 *
 *
 */
func (run *RunMiddleware) CompaniesList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	isAdmin, _, err := run.IsAdmin()
	if err != nil {
		fmt.Println("Get Isadmin:", err.Error())
		return
	}
	listOfCompanies, err := run.API.GetListOfCompanies()
	if err != nil {
		run.OutputError(err)
		return
	}

	if !run.HumanReadable {
		printAsJson(listOfCompanies)
	} else {
		if err := run.PrintCompanies(isAdmin, listOfCompanies); err != nil {
			run.OutputError(err)
		}
	}
}

func (run *RunMiddleware) PrintCompanies(isAdmin bool, companies []api.APICompany) error {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	if !isAdmin {
		_, _ = fmt.Fprintf(w, "COMPANY UUID\tNAME\tROLE\t\n")
	} else {
		_, _ = fmt.Fprintf(w, "COMPANY UUID\tNAME\t\n")
	}
	for _, company := range companies {
		if !isAdmin {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t\n",
				company.Company.UUID, company.Company.Name, company.Role.Name)
		} else {
			_, _ = fmt.Fprintf(w, "%s\t%s\t\n",
				company.Company.UUID, company.Company.Name)
		}
		_ = w.Flush()
	}
	return nil
}

func (run *RunMiddleware) CompanyDetail(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	companyUUID, _ := cmd.Flags().GetString("company-uuid")

	company, err := run.API.GetCompanyDetails(companyUUID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if !run.HumanReadable {
		printAsJson(company)
	} else {
		fmt.Printf("Company %s informations:\n", company.Name)
		fmt.Printf("  UUID: %s\n"+
			"  Phone: %s\n"+
			"  Managed: %t\n"+
			"  Description: %s\n"+
			"  Email: %s\n"+
			"  CA: %d\n"+
			"  NAF: %s\n"+
			"  Siret: %s\n"+
			"  TVA number: %s\n"+
			"  TVA rate: %d\n"+
			"  Website: %s\n"+
			"  Note: %s\n"+
			"  Quotas:\n"+
			"    CPUs: %d\n"+
			"    Networks: %d\n"+
			"    Servers: %d\n",
			company.UUID, company.Phone, company.Managed,
			company.Description, company.Email, company.CA,
			company.Naf, company.Siret, company.TvaNumber,
			company.TvaRate, company.Note, company.Website,
			company.Quotas.CPUs, company.Quotas.Networks,
			company.Quotas.Servers)
		fmt.Println("  Address:")
		for _, addr := range company.Addresses {
			fmt.Printf("    - Name: %s\n"+
				"      City: %s\n"+
				"      Country: %s\n"+
				"      Postal code: %s\n"+
				"      Type: %s\n"+
				"      Street: %s\n",
				addr.Name, addr.City, addr.Country,
				addr.PostalCode, addr.Type, addr.Street)
			if addr.Street2 != "" {
				fmt.Printf("      Street2: %s\n", addr.Street2)
			}
		}
		fmt.Println("  Members:")
		for _, member := range company.Members {
			fmt.Printf("    - Name: %s %s (%s)\n"+
				"      Phone: %s\n"+
				"      UUID: %s\n",
				member.Firstname, member.Lastname,
				member.Email, member.Phone, member.UUID)
		}
	}
}

func (run *RunMiddleware) CompanyServers(companyUUID string) ([]api.APIServer, error) {
	return run.API.GetServersOfCompany(companyUUID)
}
