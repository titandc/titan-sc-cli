package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

func (API *APITitan) SSHKeysList(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)

	if err := API.SendAndResponse(HTTPGet, "/auth/user/sshkeys", nil); err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		sshKeys := []APIUserSSHKey{}
		if err := json.Unmarshal(API.RespBody, &sshKeys); err != nil {
			fmt.Println(err.Error())
			return
		}

		sshkeySize := len(sshKeys)
		if len(sshKeys) == 0 {
			fmt.Println("SSH keys list is empty")
		} else {
			API.SSHKeysPrint(sshKeys, sshkeySize, "")
		}
	}
}

func (API *APITitan) SSHKeysGet() ([]APIUserSSHKey, error) {
	if err := API.SendAndResponse(HTTPGet, "/auth/user/sshkeys", nil); err != nil {
		return nil, err
	}
	sshKeys := []APIUserSSHKey{}
	if err := json.Unmarshal(API.RespBody, &sshKeys); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return sshKeys, nil
}

func (API *APITitan) SSHKeysPrint(sshKeys []APIUserSSHKey, sshkeySize int, ident string) {
	sshkeySize--
	fmt.Printf("%sSSH keys infos:\n", ident)
	for i, key := range sshKeys {
		fmt.Printf("%s  Title: %s\n"+
			"%s  Content: %s\n",
			ident, key.Title, ident, key.Content)
		if i != sshkeySize {
			fmt.Printf("\n")
		}
	}
}

func (API *APITitan) SSHKeyAdd(cmd *cobra.Command, args []string) {
	value := args[0]
	name, _ := cmd.Flags().GetString("name")
	API.ParseGlobalFlags(cmd)

	addSSHKey := APIAddUserSSHKey{
		Value: value,
		Name:  name,
	}
	API.SendAndPrintDefaultReply(HTTPPost, "/auth/user/sshkeys", addSSHKey)
}

func (API *APITitan) SSHKeyDel(cmd *cobra.Command, args []string) {
	_ = cmd
	name := args[0]
	API.ParseGlobalFlags(cmd)

	delSSHKey := APIDeleteUserSSHKey{
		Name: name,
	}
	API.SendAndPrintDefaultReply(HTTPDelete, "/auth/user/sshkeys", delSSHKey)
}
