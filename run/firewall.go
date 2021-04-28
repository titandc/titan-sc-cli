package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"titan-sc/api"
)

func (run *RunMiddleware) FirewallAddRule(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	protocol, _ := cmd.Flags().GetString("protocol")
	port, _ := cmd.Flags().GetString("port")
	source, _ := cmd.Flags().GetString("source")

	run.handleErrorAndGenericOutput(run.API.PostFirewallAdd(networkUUID, serverUUID, protocol,
		port, source))
}

func (run *RunMiddleware) FirewallDelRule(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	protocol, _ := cmd.Flags().GetString("protocol")
	port, _ := cmd.Flags().GetString("port")
	source, _ := cmd.Flags().GetString("source")

	run.handleErrorAndGenericOutput(run.API.DeleteFirewall(networkUUID, serverUUID, protocol,
		port, source))
}

func (run *RunMiddleware) FirewallListRules(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")

	firewallFullInfos, err := run.API.GetFireWallFullInfos(networkUUID)
	if err != nil {
		run.OutputError(err)
		return
	}
	if !run.HumanReadable {
		printAsJson(firewallFullInfos)
	} else {
		run.FirewallPrint(firewallFullInfos)
	}
}

func (run *RunMiddleware) FirewallPrint(firewallInfos *api.APINetworkFullInfosFirewall) {

	fmt.Println("Policy: " + firewallInfos.Policy)

	var w *tabwriter.Writer
	if run.HumanReadable {
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	}

	_, _ = fmt.Fprintf(w, "Server UUID\tProtocol\tPort\tSource\t\n")
	for _, rule := range firewallInfos.Rules {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", rule.Server, rule.Protocol, rule.Port, rule.Source)
	}
	_ = w.Flush()
}
