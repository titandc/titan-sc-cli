package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func (run *RunMiddleware) UserShowAllInfos(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	userInfos, err := run.API.GetUserInfos()
	if err != nil {
		run.OutputError(err)
		return
	}
	if !run.HumanReadable {
		printAsJson(userInfos)
	} else {
		log.Printf("User informations:\n"+
			"  UUID: %s\n"+
			"  Firstname: %s\n"+
			"  Lastname: %s\n"+
			"  Email: %s\n"+
			"  Phone: %s\n"+
			"  Created at: %s"+
			"  Last Login: %s\n"+
			"  Mobile: %s\n"+
			"  Preferred language: %s\n"+
			"  Salutation: %s\n"+
			"  Latest CGV signed: %t\n"+
			"  CGV signed date: %s\n"+
			"  CGV link: %s\n"+
			"  CGV version: %s\n"+
			"  2 FA enabled: %t\n",
			userInfos.UUID, userInfos.Firstname, userInfos.Lastname,
			userInfos.Email, userInfos.Phone,
			run.DateFormat(userInfos.CreatedAt), run.DateFormat(userInfos.LastLogin),
			userInfos.Mobile, userInfos.PreferredLanguage, userInfos.Salutation,
			userInfos.LatestCGVSigned, run.DateFormat(userInfos.CGVDate),
			userInfos.CGVLink, userInfos.CGVVersion, userInfos.TwoFA)

		sshkeySize := len(userInfos.SSHKeys)
		if sshkeySize == 0 {
			fmt.Println("  SSH keys: -")
		} else {
			run.SSHKeysPrint(userInfos.SSHKeys, sshkeySize, "  ")
		}
	}
}
