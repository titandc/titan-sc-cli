package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) IpCmdAdd() {
	ip := &cobra.Command{
		Use:     "ip",
		Short:   "Manage IP addresses.",
		Long:    "Manage IP addresses.",
		GroupID: "resources",
	}

	listCompanyAvailableIPs := &cobra.Command{
		Use:     "list [--company-oid COMPANY_OID]",
		Aliases: []string{"ls"},
		Short:   "List available IPs on a company.",
		Long: `List all available IPs (not attached to a server) on a company.

If --company-oid is not specified, your default company will be used.`,
		Run: cmd.runMiddleware.IPsCompanyList,
	}

	ipUpdateReverse := &cobra.Command{
		Use:   "reverse --ip IP --reverse REVERSE",
		Short: "Change IP reverse.",
		Long:  "Change IP reverse.",
		Run:   cmd.runMiddleware.IPUpdateReverse,
	}

	ipDetach := &cobra.Command{
		Use:     "detach --server-oid SERVER_OID --ip IP_ADDRESS",
		Aliases: []string{"unset"},
		Short:   "Detach an IP from a server.",
		Long:    "Detach an IP from a server.",
		Run:     cmd.runMiddleware.IPDetach,
	}

	ipAttach := &cobra.Command{
		Use:     "attach --server-oid SERVER_OID --ip IP_ADDRESS",
		Aliases: []string{"set"},
		Short:   "Attach an IP to a server.",
		Long:    "Attach an IP to a server.",
		Run:     cmd.runMiddleware.IPAttach,
	}

	cmd.RootCommand.AddCommand(ip)
	ip.AddCommand(listCompanyAvailableIPs, ipDetach, ipAttach, ipUpdateReverse)

	listCompanyAvailableIPs.Flags().StringP("company-oid", "c", "", "Company OID (uses your default company if not specified).")

	ipDetach.Flags().StringP("server-oid", "s", "", "Set server OID.")
	ipDetach.Flags().StringP("ip", "i", "", "Set IP to detach.")
	_ = ipDetach.MarkFlagRequired("server-oid")
	_ = ipDetach.MarkFlagRequired("ip")

	ipAttach.Flags().StringP("server-oid", "s", "", "Set server OID.")
	ipAttach.Flags().StringP("ip", "i", "", "Set IP to attach.")
	_ = ipAttach.MarkFlagRequired("server-oid")
	_ = ipAttach.MarkFlagRequired("ip")

	ipUpdateReverse.Flags().StringP("ip", "i", "", "Set IP address.")
	ipUpdateReverse.Flags().StringP("reverse", "r", "", "Set new IP reverse.")
	_ = ipUpdateReverse.MarkFlagRequired("ip")
	_ = ipUpdateReverse.MarkFlagRequired("reverse")
}
