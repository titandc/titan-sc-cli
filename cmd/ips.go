package cmd

import (
	"github.com/spf13/cobra"
	. "titan-sc/api"
)

var ip = &cobra.Command{
	Use:   "ip",
	Short: "Manage IP addresses.",
	Long:  "Manage IP addresses.",
}

/*var ListAvailableIPs = &cobra.Command{
    Use: "new-list",
    Aliases: []string{"nl"},
    Short: "List a few available IPs that can be ordered.",
    Long: "List a few available IPs that can be ordered.",
    Run: API.IPsList,
}*/

var ListCompanyAvailableIPs = &cobra.Command{
	Use:     "company-list --company-uuid COMPANY_UUID",
	Aliases: []string{"cl"},
	Short:   "List available IPs on a company.",
	Long:    "List all available IPs (not attached to a server) on a company.",
	Run:     API.IPsCompanyList,
}

var IPDetach = &cobra.Command{
	Use:     "detach --server-uuid SERVER_UUID --ip IP_ADDRESS",
	Aliases: []string{"unset"},
	Short:   "Detach an IP from a server.",
	Long:    "Detach an IP from a server.",
	Run:     API.IPDetach,
}

var IPAttach = &cobra.Command{
	Use:     "attach --server-uuid SERVER_UUID --ip IP_ADDRESS",
	Aliases: []string{"set"},
	Short:   "Attach an IP to a server.",
	Long:    "Attach an IP to a server.",
	Run:     API.IPAttach,
}

func ipCmdAdd() {
	rootCmd.AddCommand(ip)
	ip.AddCommand(ListCompanyAvailableIPs, IPDetach, IPAttach)

	ListCompanyAvailableIPs.Flags().StringP("company-uuid", "c", "", "Set company UUID.")
	_ = ListCompanyAvailableIPs.MarkFlagRequired("company-uuid")

	IPDetach.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	IPDetach.Flags().StringP("ip", "i", "", "Set IP to detach.")
	_ = IPDetach.MarkFlagRequired("server-uuid")
	_ = IPDetach.MarkFlagRequired("ip")

	IPAttach.Flags().StringP("server-uuid", "s", "", "Set server UUID.")
	IPAttach.Flags().StringP("ip", "i", "", "Set IP to attach.")
	_ = IPAttach.MarkFlagRequired("server-uuid")
	_ = IPAttach.MarkFlagRequired("ip")
}
