package cmd

import (
	"github.com/spf13/cobra"
)

func (cmd *CMD) IpCmdAdd() {

	ip := &cobra.Command{
		Use:   "ip",
		Short: "Manage IP addresses.",
		Long:  "Manage IP addresses.",
	}

	/*var ListAvailableIPs = &cobra.Command{
	    Use: "new-list",
	    Aliases: []string{"nl"},
	    Short: "List a few available IPs that can be ordered.",
	    Long: "List a few available IPs that can be ordered.",
	    Run: RunMiddleware.IPsList,
	}*/

	listCompanyAvailableIPs := &cobra.Command{
		Use:     "company-list --company-uuid COMPANY_UUID",
		Aliases: []string{"cl"},
		Short:   "List available IPs on a company.",
		Long:    "List all available IPs (not attached to a server) on a company.",
		Run:     cmd.runMiddleware.IPsCompanyList,
	}

	ipDetach := &cobra.Command{
		Use:     "detach --server-uuid SERVER_UUID --ip IP_ADDRESS",
		Aliases: []string{"unset"},
		Short:   "Detach an IP from a server.",
		Long:    "Detach an IP from a server.",
		Run:     cmd.runMiddleware.IPDetach,
	}

	ipAttach := &cobra.Command{
		Use:     "attach --server-uuid SERVER_UUID --ip IP_ADDRESS",
		Aliases: []string{"set"},
		Short:   "Attach an IP to a server.",
		Long:    "Attach an IP to a server.",
		Run:     cmd.runMiddleware.IPAttach,
	}

	cmd.RootCommand.AddCommand(ip)
	ip.AddCommand(listCompanyAvailableIPs, ipDetach, ipAttach)

	listCompanyAvailableIPs.Flags().StringP("company-uuid", "c", "", "Set company UUID.")
	_ = listCompanyAvailableIPs.MarkFlagRequired("company-uuid")

	ipDetach.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	ipDetach.Flags().StringP("ip", "i", "", "Set IP to detach.")
	_ = ipDetach.MarkFlagRequired("server-uuid")
	_ = ipDetach.MarkFlagRequired("ip")

	ipAttach.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	ipAttach.Flags().StringP("ip", "i", "", "Set IP to attach.")
	_ = ipAttach.MarkFlagRequired("server-uuid")
	_ = ipAttach.MarkFlagRequired("ip")
}
