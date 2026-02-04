package run

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (run *RunMiddleware) UserInfo(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	user, err := run.API.GetUserInfos()
	if err != nil {
		run.OutputError(err)
		return
	}
	if run.JSONOutput {
		printAsJson(user)
	} else {
		twoFAStatus := run.Colorize("false", "yellow")
		if user.Registration.TwoFA {
			twoFAStatus = run.Colorize("true", "green")
		}
		fmt.Printf("%s\n"+
			"  OID: %s\n"+
			"  Name: %s\n"+
			"  Email: %s\n"+
			"  Phone: %s\n"+
			"  Salutation: %s\n"+
			"  Created at: %s\n"+
			"  Last login: %s\n"+
			"  Preferred language: %s\n"+
			"  2FA enabled: %s\n"+
			"  Default company OID: %s\n",
			run.Colorize("User:", "cyan"), user.OID,
			run.Colorize(fmt.Sprintf("%s %s", user.Firstname, user.Lastname), "cyan"),
			run.Colorize(user.Email, "cyan"), user.Phone,
			user.Salutation, run.Colorize(DatePtrFormat(user.CreatedAt), "dim"),
			run.Colorize(DateFormat(user.LastLogin), "dim"),
			user.Preference.PreferredLanguage, twoFAStatus, user.DefaultCompanyOID)

		if user.LatestSignedCGV != nil {
			fmt.Printf("%s\n"+
				"  OID: %s\n"+
				"  Signed at: %s\n"+
				"  IP: %s\n",
				run.Colorize("CGV:", "cyan"), user.LatestSignedCGV.OID,
				run.Colorize(DateFormat(user.LatestSignedCGV.Date), "dim"),
				user.LatestSignedCGV.IP)
		}

		if len(user.Companies) > 0 {
			fmt.Printf("%s %s\n", run.Colorize("Companies:", "cyan"), run.Colorize(fmt.Sprintf("%d", len(user.Companies)), "green"))
			for _, company := range user.Companies {
				roleColor := getRoleColor(company.Role)
				fmt.Printf("  - OID: %s\n"+
					"    Role: %s\n",
					company.OID, run.Colorize(company.Role, roleColor))
				if company.RoleOID != "" {
					fmt.Printf("    Role OID: %s\n", company.RoleOID)
				}
				if company.Position != "" {
					fmt.Printf("    Position: %s\n", company.Position)
				}
			}
		}
	}
}
