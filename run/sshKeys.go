package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"titan-sc/api"
)

func (run *RunMiddleware) SSHKeysList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	sshKeyList, err := run.API.GetSSHKeyList()
	if err != nil {
		run.OutputError(err)
		return
	}

	if !run.HumanReadable {
		printAsJson(sshKeyList)
	} else {

		sshkeySize := len(sshKeyList)
		if sshkeySize == 0 {
			fmt.Println("SSH keys list is empty")
		} else {
			run.SSHKeysPrint(sshKeyList, sshkeySize, "")
		}
	}
}

func (run *RunMiddleware) SSHKeysPrint(sshKeys []api.APIUserSSHKey, sshkeySize int, ident string) {
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

func (run *RunMiddleware) SSHKeyAdd(cmd *cobra.Command, args []string) {
	_ = args
	name, _ := cmd.Flags().GetString("name")
	value, _ := cmd.Flags().GetString("value")
	run.ParseGlobalFlags(cmd)

	run.handleErrorAndGenericOutput(run.API.PostSSHKeyAdd(name, value))
}

func (run *RunMiddleware) SSHKeyDel(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	name, _ := cmd.Flags().GetString("name")

	run.handleErrorAndGenericOutput(run.API.DeleteSSHKey(name))
}
