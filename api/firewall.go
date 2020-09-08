package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

func (API *APITitan) FirewallAddRule(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	networkUUID := args[0]
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	protocol, _ := cmd.Flags().GetString("protocol")
	port, _ := cmd.Flags().GetString("port")
	source, _ := cmd.Flags().GetString("source")

	firewallOpt := APINetworkFirewallRule{
		ServerUUID: serverUUID,
		Protocol:   protocol,
		Port:       port,
		Source:     source,
	}

	API.SendAndPrintDefaultReply(HTTPPost, "/compute/networks/"+networkUUID+"/firewall", firewallOpt)
}

func (API *APITitan) FirewallDelRule(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	networkUUID := args[0]
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	protocol, _ := cmd.Flags().GetString("protocol")
	port, _ := cmd.Flags().GetString("port")
	source, _ := cmd.Flags().GetString("source")

	firewallOpt := APINetworkFirewallRule{
		ServerUUID: serverUUID,
		Protocol:   protocol,
		Port:       port,
		Source:     source,
	}

	API.SendAndPrintDefaultReply(HTTPDelete, "/compute/networks/"+networkUUID+"/firewall", firewallOpt)
}

func (API *APITitan) FirewallListRules(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	networkUUID := args[0]

	err := API.SendAndResponse(HTTPGet, "/compute/networks/"+networkUUID+"/firewall", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !API.HumanReadable {
		API.PrintJson()
	} else {
		var apiReplyFirewall APINetworkFullInfosFirewall
		if err := json.Unmarshal(API.RespBody, &apiReplyFirewall); err != nil {
			fmt.Println(err.Error())
			return
		}
		API.FirewallPrint(&apiReplyFirewall)
	}
}

func (API *APITitan) FirewallPrint(firewallInfos *APINetworkFullInfosFirewall) {

	fmt.Println("Policy: " + firewallInfos.Policy)

	var w *tabwriter.Writer
	if API.HumanReadable {
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	}

	_, _ = fmt.Fprintf(w, "Server UUID\tProtocol\tPort\tSource\t\n")
	for _, rule := range firewallInfos.Rules {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", rule.Server, rule.Protocol, rule.Port, rule.Source)
	}
	_ = w.Flush()
}
