package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os"
	"text/tabwriter"
	"titan-sc/api"
)

func (run *RunMiddleware) IPAttach(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	ip, _ := cmd.Flags().GetString("ip")

	parsedIP := net.ParseIP(ip)
	ipOpt := api.APIIPAttachDetach{
		IP: parsedIP.String(),
	}
	if parsedIP.To4() == nil {
		ipOpt.Version = 6
	} else {
		ipOpt.Version = 4
	}
	apiReturn, err := run.API.PostIPAttach(serverUUID, []api.APIIPAttachDetach{ipOpt})
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) IPDetach(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	ip, _ := cmd.Flags().GetString("ip")

	parsedIP := net.ParseIP(ip)
	ipOpt := api.APIIPAttachDetach{
		IP: parsedIP.String(),
	}
	if parsedIP.To4() == nil {
		ipOpt.Version = 6
	} else {
		ipOpt.Version = 4
	}
	apiReturn, err := run.API.DeleteIPDetach(serverUUID, ipOpt)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) IPsList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	ipList, err := run.API.GetIPList()
	if err != nil {
		run.OutputError(err)
		return
	}
	if !run.HumanReadable {
		printAsJson(ipList)
	} else {
		run.IPsPrint(ipList)
	}
}

func (run *RunMiddleware) IPsCompanyList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	companyUUID, _ := cmd.Flags().GetString("company-uuid")

	ipList, err := run.API.GetCompanyIPList(companyUUID)
	if err != nil {
		run.OutputError(err)
	}
	if !run.HumanReadable {
		printAsJson(ipList)
	} else {
		run.IPsPrint(ipList)
	}
}

func (run *RunMiddleware) IPsPrint(ipArray []api.APIIPAttachDetach) {
	if len(ipArray) == 0 {
		fmt.Println("Empty IPs list")
		return
	}

	var w *tabwriter.Writer
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	_, _ = fmt.Fprintf(w, "IP\tVERSION\t\n")
	for _, ip := range ipArray {
		_, _ = fmt.Fprintf(w, "%s\t%d\t\n", ip.IP, ip.Version)
	}
	_ = w.Flush()
}
