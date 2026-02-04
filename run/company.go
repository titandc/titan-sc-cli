package run

import (
	"fmt"
	"titan-sc/api"

	"github.com/spf13/cobra"
)

func (run *RunMiddleware) CompaniesList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	listOfCompanies, err := run.API.GetListOfCompanies()
	if err != nil {
		run.OutputError(err)
		return
	}

	user, err := run.API.GetUserInfos()
	if err != nil {
		user = &api.User{
			Companies: []api.UserCompany{},
		}
	}

	if run.JSONOutput {
		printAsJson(listOfCompanies)
	} else {
		if err := run.PrintCompanies(listOfCompanies, user); err != nil {
			run.OutputError(err)
		}
	}
}

func (run *RunMiddleware) PrintCompanies(companies []api.Company, user *api.User) error {
	table := NewTable("NAME", "ROLE", "OID")
	table.SetNoColor(!run.Color)

	for _, company := range companies {
		role := getRoleByCompanyOID(company.OID, user.Companies)

		var roleColorFn func(string) string
		if run.Color {
			roleColorFn = getRoleColorFn(role)
		}

		table.AddRow(
			ColName(company.Name),
			ColColor(role, roleColorFn),
			ColOID(company.OID),
		)
	}

	table.Print()
	return nil
}

// getRoleColorFn returns a color function based on role privilege level
func getRoleColorFn(role string) func(string) string {
	return ColorFn(getRoleColor(role))
}

// getRoleColor returns the color name for a given role (shared across files)
func getRoleColor(role string) string {
	switch role {
	case "SUPER_ADMINISTRATOR":
		return "magenta"
	case "ADMINISTRATOR":
		return "yellow"
	case "USER":
		return "green"
	default:
		return "dim"
	}
}

func getRoleByCompanyOID(companyOID string, companies []api.UserCompany) string {
	var role string

	for _, co := range companies {
		if companyOID == co.OID {
			if co.Role != "" {
				role = co.Role
			} else {
				role = co.RoleOID
			}
			break
		}
	}
	return role
}

func (run *RunMiddleware) CompanyDetail(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	companyOID, err := run.ResolveCompanyOID(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}

	company, err := run.API.GetCompanyDetails(companyOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(company)
	} else {
		run.printCompanyHuman(company)
	}
}

func (run *RunMiddleware) printCompanyHuman(company *api.Company) {
	fmt.Printf("%s %s\n", run.Colorize("Company:", "cyan"), run.Colorize(company.Name, "cyan"))
	fmt.Printf("  OID:      %s\n", company.OID)
	fmt.Printf("  UUID:     %s\n", company.UUID)
	fmt.Printf("  Email:    %s\n", run.Colorize(company.Email, "cyan"))
	if company.Phone != "" {
		fmt.Printf("  Phone:    %s\n", company.Phone)
	}
	if company.Website != "" {
		fmt.Printf("  Website:  %s\n", run.Colorize(company.Website, "cyan"))
	}

	// VAT information
	fmt.Println(run.Colorize("  VAT:", "cyan"))
	if company.VAT.Number != "" {
		fmt.Printf("    Number:    %s\n", company.VAT.Number)
		validStatus := run.Colorize("No", "yellow")
		if company.VAT.Valid {
			validStatus = run.Colorize("Yes", "green")
		}
		fmt.Printf("    Valid:     %s\n", validStatus)
	} else {
		fmt.Println(run.Colorize("    (not set)", "dim"))
	}

	// Addresses
	run.printAddress("  Shipping Address:", &company.AddressShipping)
	run.printAddress("  Billing Address:", &company.AddressBilling)

	// Account settings
	fmt.Println(run.Colorize("  Settings:", "cyan"))
	renewStatus := run.Colorize("No", "yellow")
	if company.RenewWithCredits {
		renewStatus = run.Colorize("Yes", "green")
	}
	fmt.Printf("    Renew with credits: %s\n", renewStatus)
	if company.DefaultPaymentMethod != nil && *company.DefaultPaymentMethod != "" {
		fmt.Printf("    Payment method:    %s\n", *company.DefaultPaymentMethod)
	}
	if company.Disable {
		fmt.Printf("    Status:            %s\n", run.Colorize("DISABLED", "red"))
	}
}

func (run *RunMiddleware) printAddress(addrType string, addr *api.Address) {
	fmt.Println(run.Colorize(addrType, "cyan"))
	fmt.Printf("    Street:      %s\n", addr.Street)
	if addr.Street2 != "" {
		fmt.Printf("    Street2:     %s\n", addr.Street2)
	}
	fmt.Printf("    City:        %s\n", addr.City)
	fmt.Printf("    Postal code: %s\n", addr.PostalCode)
	fmt.Printf("    Country:     %s (%s)\n", addr.Country, addr.CountryCode)
}

// Colorize adds ANSI color codes if color is enabled
func (run *RunMiddleware) Colorize(text, color string) string {
	if !run.Color {
		return text
	}
	colors := map[string]string{
		"red":     "\033[1;31m",
		"green":   "\033[1;32m",
		"yellow":  "\033[1;33m",
		"blue":    "\033[1;34m",
		"magenta": "\033[1;35m",
		"cyan":    "\033[1;36m",
		"white":   "\033[1;37m",
		"dim":     "\033[2m",
		"reset":   "\033[0m",
	}
	return colors[color] + text + colors["reset"]
}
