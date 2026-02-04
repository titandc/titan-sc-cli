package run

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func (run *RunMiddleware) SSHKeysList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	// Get current user to use as target_oid
	user, err := run.API.GetUserInfos()
	if err != nil {
		run.OutputError(err)
		return
	}

	sshKeyList, err := run.API.GetSSHKeyList(user.OID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(sshKeyList)
	} else {
		if len(sshKeyList) == 0 {
			fmt.Println("No SSH keys found.")
		} else {
			table := NewTable("NAME", "TYPE", "COMMENT", "OID")
			table.SetNoColor(!run.Color)
			for _, key := range sshKeyList {
				keyType, comment := parseSSHKey(key.Value)
				var commentColorFn func(string) string
				if run.Color {
					commentColorFn = ColorFn("dim")
				}
				table.AddRow(
					ColName(key.Name),
					Col(keyType),
					ColColor(comment, commentColorFn),
					ColOID(key.OID),
				)
			}
			table.Print()
		}
	}
}

// parseSSHKey extracts the key type and comment from an SSH public key
func parseSSHKey(value string) (keyType, comment string) {
	parts := strings.Fields(value)
	if len(parts) >= 1 {
		keyType = parts[0]
	}
	if len(parts) >= 3 {
		comment = parts[2]
	}
	return
}

func (run *RunMiddleware) SSHKeyShow(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	oid, _ := cmd.Flags().GetString("oid")

	sshKey, err := run.API.GetSSHKey(oid)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(sshKey)
	} else {
		keyType, comment := parseSSHKey(sshKey.Value)
		fmt.Printf("%s\n"+
			"  OID: %s\n"+
			"  Name: %s\n"+
			"  Type: %s\n"+
			"  Comment: %s\n"+
			"  Value: %s\n",
			run.Colorize("SSH Key:", "cyan"), sshKey.OID, run.Colorize(sshKey.Name, "cyan"), keyType, run.Colorize(comment, "dim"), sshKey.Value)
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
	oid, _ := cmd.Flags().GetString("oid")
	run.handleErrorAndGenericOutput(run.API.DeleteSSHKey(oid))
}
