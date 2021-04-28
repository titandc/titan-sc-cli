package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"titan-sc/api"
)

func (run *RunMiddleware) IPPNATRuleAdd(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	ip, _ := cmd.Flags().GetString("ip")
	transparent, _ := cmd.Flags().GetBool("transparent")
	protocol, _ := cmd.Flags().GetString("protocol")
	portSrc, _ := cmd.Flags().GetInt64("port-src")
	portDst, _ := cmd.Flags().GetInt64("port-dst")
	run.handleErrorAndGenericOutput(run.API.PostIPPNATRuleAdd(serverUUID, ip, protocol,
		transparent, portSrc, portDst))
}

func (run *RunMiddleware) IPPNATRuleDel(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	ip, _ := cmd.Flags().GetString("ip")
	transparent, _ := cmd.Flags().GetBool("transparent")
	protocol, _ := cmd.Flags().GetString("protocol")
	portSrc, _ := cmd.Flags().GetInt64("port-src")
	portDst, _ := cmd.Flags().GetInt64("port-dst")
	run.handleErrorAndGenericOutput(run.API.DeleteIPPNATRule(serverUUID, ip, protocol, transparent,
		portSrc, portDst))
}

func (run *RunMiddleware) ListServerPNATRules(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	PNATRulesList, err := run.API.GetServerPNATRulesList(serverUUID)
	if err != nil {
		run.OutputError(err)
		return
	}
	if !run.HumanReadable {
		printAsJson(PNATRulesList)
	} else {
		run.PNATRulesPrint(PNATRulesList)
	}
}

func (run *RunMiddleware) PNATRulesPrint(pnatRulesArray []api.APIPNATRuleInfos) {
	if len(pnatRulesArray) == 0 {
		fmt.Println("Empty PNAT rules list")
		return
	}

	var w *tabwriter.Writer
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	_, _ = fmt.Fprintf(w, "IP\tTRANSPARENT\tPROTOCOL\tPORT_SRC\tPORT_DST\t\n")
	for _, pnatRule := range pnatRulesArray {
		if pnatRule.Transparent {
			_, _ = fmt.Fprintf(w, "%s\t%t\t%s\t%s\t%s\t\n",
				pnatRule.IP, pnatRule.Transparent, "-", "-", "-")

		} else {
			_, _ = fmt.Fprintf(w, "%s\t%t\t%s\t%d\t%d\t\n",
				pnatRule.IP, pnatRule.Transparent, pnatRule.Protocol, pnatRule.PortSrc, pnatRule.PortDst)

		}
	}
	_ = w.Flush()
}
