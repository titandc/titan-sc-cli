package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) PnatCmdAdd() {

	pnat := &cobra.Command{
		Use:     "port-nat",
		Aliases: []string{"pnat", "pat"},
		Short:   "Manage PNAT rules.",
		Long:    "Manage port address translation (PNAT) rules.",
	}

	ListServerPNATRules := &cobra.Command{
		Use:     "list --server-uuid SERVER_UUID",
		Aliases: []string{"lp"},
		Short:   "List all PNAT rules of a server.",
		Long:    "List all PNAT rules of a server.",
		Run:     cmd.runMiddleware.ListServerPNATRules,
	}

	AddIPPNATRule := &cobra.Command{
		Use:     "add --ip IP_ADDRESS --server-uuid SERVER_UUID [--transparent | --protocol --port-src --port-dst]",
		Aliases: []string{"new"},
		Short:   "Add a PNAT rule.",
		Long:    "Add a PNAT rule from a public IP to a managed server.",
		Run:     cmd.runMiddleware.IPPNATRuleAdd,
	}

	DelIPPNATRule := &cobra.Command{
		Use:     "del --ip IP_ADDRESS --server-uuid SERVER_UUID [--transparent | --protocol --port-src --port-dst]",
		Aliases: []string{"rm"},
		Short:   "Delete a PNAT rule.",
		Long:    "Delete a PNAT rule between a public IP and a managed server.",
		Run:     cmd.runMiddleware.IPPNATRuleDel,
	}

	cmd.RootCommand.AddCommand(pnat)
	pnat.AddCommand(ListServerPNATRules, AddIPPNATRule, DelIPPNATRule)

	ListServerPNATRules.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	_ = ListServerPNATRules.MarkFlagRequired("server-uuid")

	AddIPPNATRule.Flags().StringP("ip", "i", "", "Set public IP.")
	AddIPPNATRule.Flags().BoolP("transparent", "t", false, "Enable transparent redirection.")
	AddIPPNATRule.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
	AddIPPNATRule.Flags().StringP("protocol", "p", "", "Set protocol (not used for transparent rules).")
	AddIPPNATRule.Flags().Int64P("port-src", "s", 0, "Set source port from the public IP (not used for transparent rules).")
	AddIPPNATRule.Flags().Int64P("port-dst", "d", 0, "Set destination port to the server (not used for transparent rules).")
	_ = AddIPPNATRule.MarkFlagRequired("ip")
	_ = AddIPPNATRule.MarkFlagRequired("server-uuid")

	DelIPPNATRule.Flags().StringP("ip", "i", "", "Set public IP.")
	DelIPPNATRule.Flags().BoolP("transparent", "t", false, "Target a transparent redirection rule.")
	DelIPPNATRule.Flags().StringP("server-uuid", "u", "", "Set server UUID.")
	DelIPPNATRule.Flags().StringP("protocol", "p", "", "Set protocol (not used for transparent rules).")
	DelIPPNATRule.Flags().Int64P("port-src", "s", 0, "Set source port from the public IP (not used for transparent rules).")
	DelIPPNATRule.Flags().Int64P("port-dst", "d", 0, "Set destination port to the server (not used for transparent rules).")
	_ = DelIPPNATRule.MarkFlagRequired("ip")
	_ = DelIPPNATRule.MarkFlagRequired("server-uuid")
}
