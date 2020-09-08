package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

func (API *APITitan) UserShowAllInfos(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)

	if err := API.SendAndResponse(HTTPGet, "/auth/user", nil); err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		userInfos := &APIUserInfos{}
		if err := json.Unmarshal(API.RespBody, userInfos); err != nil {
			fmt.Println(err.Error())
			return
		}
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
			API.DateFormat(userInfos.CreatedAt), API.DateFormat(userInfos.LastLogin),
			userInfos.Mobile, userInfos.PreferredLanguage, userInfos.Salutation,
			userInfos.LatestCGVSigned, API.DateFormat(userInfos.CGVDate),
			userInfos.CGVLink, userInfos.CGVVersion, userInfos.TwoFA)

		sshkeySize := len(userInfos.SSHKeys)
		if sshkeySize == 0 {
			fmt.Println("  SSH keys: -")
		} else {
			API.SSHKeysPrint(userInfos.SSHKeys, sshkeySize, "  ")
		}
	}
}
