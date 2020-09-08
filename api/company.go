package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
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
func (API *APITitan) CompaniesList(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)

	if compagnies, err := API.GetCompagnies(); err != nil {
		fmt.Println(err.Error())
	} else {
		if !API.HumanReadable {
			API.PrintJson()
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			_, _ = fmt.Fprintf(w, "COMPANY UUID\tNAME\tROLE\t\n")
			for _, company := range compagnies {
				_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t\n",
					company.Company.UUID, company.Company.Name, company.Role.Name)
				_ = w.Flush()
			}
		}
	}
	return
}

func (API *APITitan) CompanyDetail(cmd *cobra.Command, args []string) {
	companyUUID := args[0]
	API.ParseGlobalFlags(cmd)

	err := API.SendAndResponse(HTTPGet, "/companies/"+companyUUID, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		company := &APICompanyDetail{}
		if err := json.Unmarshal(API.RespBody, &company); err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Company %s informations:\n", company.Name)
		fmt.Printf("  UUID: %s\n"+
			"  Phone: %s\n"+
			"  Description: %s\n"+
			"  Email: %s\n"+
			"  CA: %d\n"+
			"  NAF: %s\n"+
			"  Siret: %s\n"+
			"  TVA number: %s\n"+
			"  TVA rate: %d\n"+
			"  Website: %s\n"+
			"  Note: %s\n"+
			"  Qutas:\n"+
			"    CPUs: %d\n"+
			"    Networks: %d\n"+
			"    Servers: %d\n",
			company.UUID, company.Phone,
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

func (API *APITitan) GetCompagnies() ([]APICompany, error) {
	if err := API.SendAndResponse(HTTPGet, "/companies", nil); err != nil {
		return nil, err
	}

	compagnies := make([]APICompany, 0)
	if err := json.Unmarshal(API.RespBody, &compagnies); err != nil {
		return nil, err
	}
	return compagnies, nil
}

func (API *APITitan) GetCompanyServers(companyUUID string) ([]APIServer, error) {
	err := API.SendAndResponse(HTTPGet, "/compute/servers?company_uuid="+companyUUID, nil)
	if err != nil {
		return nil, err
	}

	servers := make([]APIServer, 0)
	if err := json.Unmarshal(API.RespBody, &servers); err != nil {
		return nil, err
	}
	return servers, nil
}
