package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) SSHKeysCmdAdd() {

	sshKeys := &cobra.Command{
		Use:   "ssh-key",
		Short: "Manage your user ssh keys.",
		Long:  "Manage your user ssh keys.",
	}

	sshKeysList := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Get all user informations.",
		Long:    "Get all user informations.",
		Run:     cmd.runMiddleware.SSHKeysList,
	}

	sshKeyAdd := &cobra.Command{
		Use:   "add --name \"NAME\" --value \"SSH_KEYS_VALUE\"",
		Short: "Add one ssh key.",
		Long:  "Add one ssh key\nNeed name and ssh key value.",
		Run:   cmd.runMiddleware.SSHKeyAdd,
	}

	sshKeyDel := &cobra.Command{
		Use:     "delete --name \"SSH_KEY_NAME\"",
		Aliases: []string{"del"},
		Short:   "Delete one ssh key by name.",
		Long:    "Delete one ssh key by name.",
		Run:     cmd.runMiddleware.SSHKeyDel,
	}

	cmd.RootCommand.AddCommand(sshKeys)

	sshKeys.AddCommand(sshKeysList, sshKeyAdd, sshKeyDel)

	sshKeyAdd.Flags().StringP("name", "n", "", "Name of ssh KEY.")
	sshKeyAdd.Flags().StringP("value", "v", "", "value of ssh KEY.")
	_ = sshKeyAdd.MarkFlagRequired("name")
	_ = sshKeyAdd.MarkFlagRequired("value")

	sshKeyDel.Flags().StringP("name", "n", "", "Name of ssh KEY.")
	_ = sshKeyDel.MarkFlagRequired("name")
}
