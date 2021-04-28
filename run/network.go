package run

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
	"titan-sc/api"
)

/*
 *
 *
 ******************
 * Network function
 ******************
 *
 *
 */

func (run *RunMiddleware) NetworkList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	companyUUID, _ := cmd.Flags().GetString("company-uuid")

	networks, err := run.API.GetNetworkList(companyUUID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if !run.HumanReadable {
		printAsJson(networks)
	} else {
		fmt.Println("Quota:", networks.Quota)
		for _, net := range networks.NetInfos {
			run.NetworkPrintBase(&net)
			fmt.Println("  Servers list:")
			for _, server := range net.Servers {
				fmt.Printf("    - Name: %s\n"+
					"      OS: %s\n"+
					"      Plan: %s\n"+
					"      State: %s\n"+
					"      UUID: %s\n",
					server.Name, server.OS, server.Plan, server.State, server.UUID)
			}
			fmt.Printf("\n")
		}
	}
}

func (run *RunMiddleware) NetworkDetail(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")

	network, err := run.API.GetNetworkDetail(networkUUID)
	if err != nil {
		run.OutputError(err)
		return
	}
	if !run.HumanReadable {
		printAsJson(network)
	} else {
		run.NetworkPrintBase(network)
		fmt.Println("  Servers list:")
		for _, server := range network.Servers {
			fmt.Printf("    - Name: %s\n"+
				"      OS: %s\n"+
				"      Plan: %s\n"+
				"      State: %s\n"+
				"      UUID: %s\n",
				server.Name, server.OS, server.Plan, server.State, server.UUID)
		}
		fmt.Printf("\n")
	}
}

type APINetworkOps struct {
	ServerUUID string `json:"server_uuid"`
}

func (run *RunMiddleware) NetworkAttachServer(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	act := APINetworkOps{
		ServerUUID: serverUUID,
	}
	_, apiRequest, err := run.API.SendRequestToAPI(api.HTTPPut, "/compute/networks/"+networkUUID+"/attach", act)
	run.handleErrorAndGenericOutput(apiRequest, err)
}

func (run *RunMiddleware) NetworkDetachServer(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	act := APINetworkOps{
		ServerUUID: serverUUID,
	}
	_, apiReturn, err := run.API.SendRequestToAPI(api.HTTPPut,
		"/compute/networks/"+networkUUID+"/detach", act)
	// Render error or success output
	run.handleErrorAndGenericOutput(apiReturn, err)

}

func (run *RunMiddleware) NetworkCreate(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	companyUUID, _ := cmd.Flags().GetString("company-uuid")
	networkName, _ := cmd.Flags().GetString("name")
	cidr, _ := cmd.Flags().GetString("cidr")

	net := api.APINetworkCreate{
		MaxMTU: 8948,
		Name:   networkName,
		Ports:  6,
		CIDR:   cidr,
	}
	net.Speed.Value = 1
	net.Speed.Unit = "Gbps"

	network, err := run.API.CreateNetwork(companyUUID, net)
	if err != nil {
		run.OutputError(err)
		return
	}

	if !run.HumanReadable {
		printAsJson(network)
	} else {
		run.NetworkPrintBase(network)
	}
}

func (run *RunMiddleware) NetworkRemove(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")
	// Handle error, print reply
	if err := run.API.RemoveNetwork(networkUUID); err != nil {

	}
}

func (run *RunMiddleware) DateFormat(timestamp int64) string {
	dateMls := time.Unix(0, int64(timestamp)*int64(time.Millisecond))
	date := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		dateMls.Year(), dateMls.Month(), dateMls.Day(),
		dateMls.Hour(), dateMls.Minute(), dateMls.Second())
	return date
}

func (run *RunMiddleware) NetworkPrintBase(net *api.APINetwork) {
	date := run.DateFormat(net.CreatedAt)
	fmt.Printf("Network information:\n"+
		"  Name: %s\n"+
		"  Created at: %s\n"+
		"  Ports: %d\n"+
		"  Speed: %d %s\n"+
		"  State: %s\n"+
		"  UUID: %s\n"+
		"  Company: %s\n"+
		"  Max MTU: %d\n",
		net.Name, date, net.Ports, net.Speed.Value, net.Speed.Unit,
		net.State, net.UUID, net.Company, net.MaxMtu)

	if net.Managed {
		fmt.Printf("  Managed: %t\n"+
			"  CIDR: %s\n",
			net.Managed, net.CIDR)
		if net.Gateway != "" {
			fmt.Printf("  Gateway: %s\n", net.Gateway)
		}
		fmt.Printf("  Firewall:\n"+
			"    Policy: %s\n",
			net.Firewall.Policy)
		if len(net.Firewall.Rules) > 0 {
			fmt.Println("    Rules:")
			for _, rule := range net.Firewall.Rules {
				fmt.Printf("      - Server: %s\n"+
					"        Protocol: %s\n"+
					"        Port: %s\n"+
					"        Source: %s\n",
					rule.Server, rule.Protocol, rule.Port, rule.Source)
			}
		}
	}

	fmt.Printf("  Owner informations:\n"+
		"    Name: %s %s (%s)\n"+
		"    UUID: %s\n",
		net.Owner.Firstname, net.Owner.Lastname, net.Owner.Email,
		net.Owner.UUID)
}

type APINetworkRename struct {
	Name string `json:"name"`
}

func (run *RunMiddleware) NetworkRename(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")
	name, _ := cmd.Flags().GetString("name")

	netRename := &APINetworkRename{Name: name}
	_, apiReturn, err := run.API.SendRequestToAPI(api.HTTPPut, "/compute/networks/"+networkUUID, netRename)
	// Handle error, print output
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) NetworkSetGateway(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")
	ip, _ := cmd.Flags().GetString("ip")

	ipData := api.APIIPAttachDetach{
		IP:      ip,
		Version: 4,
	}
	_, apiReturn, err := run.API.SendRequestToAPI(api.HTTPPut, "/compute/networks/"+networkUUID+"/gateway", ipData)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) NetworkUnsetGateway(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkUUID, _ := cmd.Flags().GetString("network-uuid")
	_, apiReturn, err := run.API.SendRequestToAPI(api.HTTPDelete, "/compute/networks/"+networkUUID+"/gateway",
		nil)
	// Handle error, print output
	run.handleErrorAndGenericOutput(apiReturn, err)
}
