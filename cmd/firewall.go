package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) FirewallCmdAdd() {

	firewall := &cobra.Command{
		Use:   "firewall",
		Short: "Manage your networks firewall rules.",
		Long:  "Manage your networks firewall rules.",
	}

	firewallRulesList := &cobra.Command{
		Use:     "list --network-uuid NETWORK_UUID",
		Aliases: []string{"ls"},
		Short:   "Get all firewall rules.",
		Long:    "Get all firewall rules of a managed network.",
		Run:     cmd.runMiddleware.FirewallListRules,
	}

	firewallRulesAdd := &cobra.Command{
		Use:   "add --network-uuid NETWORK_UUID",
		Short: "Add a firewall rule.",
		Long:  "Add a firewall rule to a managed network.",
		Run:   cmd.runMiddleware.FirewallAddRule,
	}

	firewallRulesDel := &cobra.Command{
		Use:     "delete --network-uuid NETWORK_UUID",
		Aliases: []string{"del"},
		Short:   "Delete a firewall rule.",
		Long:    "Delete a firewall rule from a managed network.",
		Run:     cmd.runMiddleware.FirewallDelRule,
	}

	cmd.RootCommand.AddCommand(firewall)
	firewall.AddCommand(firewallRulesList, firewallRulesAdd, firewallRulesDel)

	// Firewall list rules
	firewallRulesList.Flags().StringP("network-uuid", "n", "", "Set network UUID.")
	_ = firewallRulesList.MarkFlagRequired("network-uuid")

	// firewall rules add
	firewallRulesAdd.Flags().StringP("network-uuid", "n", "", "Set network UUID.")
	firewallRulesAdd.Flags().StringP("server-uuid", "u", "", "Set targeted server UUID or '*' for the whole subnet.")
	firewallRulesAdd.Flags().StringP("protocol", "p", "", "Set targeted protocol, it could be 'tcp', 'udp' or '*' for both protocols.")
	firewallRulesAdd.Flags().StringP("port", "o", "", "Set targeted port(s), use '*' to open all ports, use 'PORT_START:PORT_END' to open a range of ports (eg. '2000:3000').")
	firewallRulesAdd.Flags().StringP("source", "s", "", "Set a specific source IP or CIDR, use '*' to allow any source.")
	_ = firewallRulesAdd.MarkFlagRequired("network-uuid")
	_ = firewallRulesAdd.MarkFlagRequired("server-uuid")
	_ = firewallRulesAdd.MarkFlagRequired("protocol")
	_ = firewallRulesAdd.MarkFlagRequired("port")
	_ = firewallRulesAdd.MarkFlagRequired("source")

	// Firewall rules del
	firewallRulesDel.Flags().StringP("network-uuid", "n", "", "Set network UUID.")
	firewallRulesDel.Flags().StringP("server-uuid", "u", "", "Set targeted server UUID or '*' for the whole subnet.")
	firewallRulesDel.Flags().StringP("protocol", "p", "", "Set targeted protocol, it could be 'tcp', 'udp' or '*' for both protocols.")
	firewallRulesDel.Flags().StringP("port", "o", "", "Set targeted port(s), use '*' to open all ports, use 'PORT_START:PORT_END' to open a range of ports (eg. '2000:3000').")
	firewallRulesDel.Flags().StringP("source", "s", "", "Set a specific source IP or CIDR, use '*' to allow any source.")
	_ = firewallRulesDel.MarkFlagRequired("network-uuid")
	_ = firewallRulesDel.MarkFlagRequired("server-uuid")
	_ = firewallRulesDel.MarkFlagRequired("protocol")
	_ = firewallRulesDel.MarkFlagRequired("port")
	_ = firewallRulesDel.MarkFlagRequired("source")
}
