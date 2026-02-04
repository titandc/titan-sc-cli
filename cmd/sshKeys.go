package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) SSHKeysCmdAdd() {

	sshKeys := &cobra.Command{
		Use:   "ssh-key",
		Short: "Manage your ssh keys.",
		Long:  "Manage your ssh keys.",
	}

	sshKeysList := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all SSH keys.",
		Long:    "List all SSH keys for the current user.",
		Run:     cmd.runMiddleware.SSHKeysList,
	}

	sshKeyShow := &cobra.Command{
		Use:   "show --oid OID",
		Short: "Show SSH key details.",
		Long:  "Show full details of an SSH key by OID.",
		Run:   cmd.runMiddleware.SSHKeyShow,
	}

	sshKeyAdd := &cobra.Command{
		Use:   "add --name \"NAME\" --value \"SSH_KEYS_VALUE\"",
		Short: "Add one ssh key.",
		Long:  "Add one ssh key\nNeed name and ssh key value.",
		Run:   cmd.runMiddleware.SSHKeyAdd,
	}

	sshKeyDel := &cobra.Command{
		Use:     "delete --oid \"SSH_KEY_NAME\"",
		Aliases: []string{"del"},
		Short:   "Delete one ssh key by oid.",
		Long:    "Delete one ssh key by oid.",
		Run:     cmd.runMiddleware.SSHKeyDel,
	}

	cmd.RootCommand.AddCommand(sshKeys)

	sshKeys.AddCommand(sshKeysList, sshKeyShow, sshKeyAdd, sshKeyDel)

	sshKeyShow.Flags().StringP("oid", "o", "", "OID of the SSH key.")
	_ = sshKeyShow.MarkFlagRequired("oid")

	sshKeyAdd.Flags().StringP("name", "n", "", "Name of SSH key.")
	sshKeyAdd.Flags().StringP("value", "v", "", "Value of SSH key.")
	_ = sshKeyAdd.MarkFlagRequired("name")
	_ = sshKeyAdd.MarkFlagRequired("value")

	sshKeyDel.Flags().StringP("oid", "o", "", "OID of SSH key.")
	_ = sshKeyDel.MarkFlagRequired("oid")
}
