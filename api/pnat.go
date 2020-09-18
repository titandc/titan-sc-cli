package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

func (API *APITitan) IPPNATRuleAdd(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	ip, _ := cmd.Flags().GetString("ip")
	transparent, _ := cmd.Flags().GetBool("transparent")
	protocol, _ := cmd.Flags().GetString("protocol")
	portSrc, _ := cmd.Flags().GetInt64("port-src")
	portDst, _ := cmd.Flags().GetInt64("port-dst")

	pnatOpt := APIPNATRuleAddDel{
		IP:          ip,
		Transparent: transparent,
		Protocol:    protocol,
		PortSrc:     portSrc,
		PortDst:     portDst,
	}
	API.SendAndPrintDefaultReply(HTTPPost, "/compute/servers/"+serverUUID+"/pnat", pnatOpt)
}

func (API *APITitan) IPPNATRuleDel(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	ip, _ := cmd.Flags().GetString("ip")
	transparent, _ := cmd.Flags().GetBool("transparent")
	protocol, _ := cmd.Flags().GetString("protocol")
	portSrc, _ := cmd.Flags().GetInt64("port-src")
	portDst, _ := cmd.Flags().GetInt64("port-dst")

	pnatOpt := APIPNATRuleAddDel{
		IP:          ip,
		Transparent: transparent,
		Protocol:    protocol,
		PortSrc:     portSrc,
		PortDst:     portDst,
	}
	API.SendAndPrintDefaultReply(HTTPDelete, "/compute/servers/"+serverUUID+"/pnat", pnatOpt)
}

func (API *APITitan) ListServerPNATRules(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	err := API.SendAndResponse(HTTPGet, "/compute/servers/"+serverUUID+"/pnat", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !API.HumanReadable {
		API.PrintJson()
	} else {
		APIPNATRules := make([]APIPNATRuleInfos, 0)
		if err := json.Unmarshal(API.RespBody, &APIPNATRules); err != nil {
			fmt.Println(err.Error())
			return
		}
		API.PNATRulesPrint(&APIPNATRules)
	}
}

func (API *APITitan) PNATRulesPrint(pnatRulesArray *[]APIPNATRuleInfos) {
	if len(*pnatRulesArray) == 0 {
		fmt.Println("Empty PNAT rules list")
		return
	}

	var w *tabwriter.Writer
	w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	_, _ = fmt.Fprintf(w, "IP\tTRANSPARENT\tPROTOCOL\tPORT_SRC\tPORT_DST\t\n")
	for _, pnatRule := range *pnatRulesArray {
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
