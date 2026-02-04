package run

import (
	"fmt"
	"time"
	"titan-sc/api"

	"github.com/spf13/cobra"
)

// Simplified structs for clean JSON output
type networkDetailOutput struct {
	OID        string                   `json:"oid"`
	Name       string                   `json:"name"`
	State      string                   `json:"state"`
	CreatedAt  string                   `json:"created_at,omitempty"`
	Ports      uint                     `json:"ports"`
	MaxMTU     uint                     `json:"max_mtu"`
	Speed      api.NetworkSpeed         `json:"speed"`
	Interfaces []networkInterfaceOutput `json:"interfaces"`
}

type networkInterfaceOutput struct {
	OID        string `json:"oid"`
	MAC        string `json:"mac"`
	ServerOID  string `json:"server_oid"`
	ServerName string `json:"server_name"`
	ServerIP   string `json:"server_ip,omitempty"`
}

func (run *RunMiddleware) NetworkList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	companyOID, err := run.ResolveCompanyOID(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}

	networks, err := run.API.GetNetworkList(companyOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		// Build clean output - array at root like other list commands
		output := make([]networkDetailOutput, 0, len(networks.Networks))
		for _, net := range networks.Networks {
			output = append(output, toNetworkDetailOutput(&net))
		}
		printAsJson(output)
	} else {
		if len(networks.Networks) == 0 {
			fmt.Println("No networks found.")
			return
		}

		table := NewTable("NAME", "STATE", "DRP", "SPEED", "PORTS", "SERVERS", "OID")
		table.SetNoColor(!run.Color)
		for _, net := range networks.Networks {
			speed := fmt.Sprintf("%d %s", net.Speed.Value, net.Speed.Unit)
			servers := fmt.Sprintf("%d", len(net.Interfaces))

			// DRP indicator
			drpIndicator := "-"
			if net.Drp != nil && net.Drp.Enabled {
				drpIndicator = "✓"
				if run.Color {
					drpIndicator = run.Colorize(drpIndicator, "green")
				}
			}

			var stateColorFn func(string) string
			if run.Color {
				stateColorFn = StateColorFn(net.State)
			}

			table.AddRow(
				ColName(net.Name),
				ColColor(net.State, stateColorFn),
				Col(drpIndicator),
				Col(speed),
				Col(fmt.Sprintf("%d", net.Ports)),
				Col(servers),
				ColOID(net.OID),
			)
		}
		table.Print()
	}
}

func (run *RunMiddleware) NetworkDetail(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkOID, _ := cmd.Flags().GetString("network-oid")

	network, err := run.API.GetNetworkDetail(networkOID)
	if err != nil {
		run.OutputError(err)
		return
	}
	if run.JSONOutput {
		// Use clean output struct for consistency
		printAsJson(toNetworkDetailOutput(network))
	} else {
		run.printNetworkDetail(network)
		fmt.Printf("\n")
	}
}

func (run *RunMiddleware) NetworkAttachServer(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkOID, _ := cmd.Flags().GetString("network-oid")
	serverOID, _ := cmd.Flags().GetString("server-oid")

	act := api.NetworkOps{
		ServerOIDs: []string{serverOID},
	}
	_, apiRequest, err := run.API.SendRequestToAPI(api.HTTPPut, "/network/switch/"+networkOID+"/attach", act)
	run.handleErrorAndGenericOutput(apiRequest, err)
}

func (run *RunMiddleware) NetworkDetachServer(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkOID, _ := cmd.Flags().GetString("network-oid")
	serverOID, _ := cmd.Flags().GetString("server-oid")

	act := api.NetworkOps{
		ServerOID: serverOID,
	}
	_, apiReturn, err := run.API.SendRequestToAPI(api.HTTPPut, "/network/switch/"+networkOID+"/detach", act)
	// Render error or success output
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) NetworkCreate(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkName, _ := cmd.Flags().GetString("name")

	network, err := run.API.CreateNetwork(&api.NetworkCreate{
		Name: networkName,
	})
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(network)
	} else {
		run.printNetwork(network)
	}
}

func (run *RunMiddleware) NetworkRemove(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkOID, _ := cmd.Flags().GetString("network-oid")
	// Handle error, print reply
	if err := run.API.RemoveNetwork(networkOID); err != nil {

	}
}

func DatePtrFormat(timestamp *int64) string {
	if timestamp == nil {
		return ""
	}
	return DateFormat(*timestamp)
}

// FlexTimestampFormat formats a FlexTimestamp for display
func FlexTimestampFormat(ft *api.FlexTimestamp) string {
	if ft == nil || ft.Value == nil || *ft.Value == 0 {
		return ""
	}
	return DateFormat(*ft.Value)
}

func DateFormat(timestamp int64) string {
	dateMls := time.Unix(0, timestamp*int64(time.Millisecond))
	date := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		dateMls.Year(), dateMls.Month(), dateMls.Day(),
		dateMls.Hour(), dateMls.Minute(), dateMls.Second())
	return date
}

func (run *RunMiddleware) printNetwork(net *api.Network) {
	date := DatePtrFormat(net.CreatedAt)
	fmt.Printf("%s %s\n", run.Colorize("Network:", "cyan"), run.Colorize(net.Name, "cyan"))
	fmt.Printf("  OID:      %s\n", run.Colorize(net.OID, "blue"))
	fmt.Printf("  Created:  %s\n", run.Colorize(date, "dim"))
	fmt.Printf("  Speed:    %d %s\n", net.Speed.Value, net.Speed.Unit)
	fmt.Printf("  Ports:    %d\n", net.Ports)
	fmt.Printf("  Max MTU:  %d\n", net.MaxMTU)
}

func (run *RunMiddleware) printNetworkDetail(net *api.NetworkDetail) {
	date := DatePtrFormat(net.CreatedAt)
	state := GetStateColorized(run.Color, net.State)
	fmt.Printf("%s %s\n", run.Colorize("Network:", "cyan"), run.Colorize(net.Name, "cyan"))
	fmt.Printf("  OID:        %s\n", run.Colorize(net.OID, "blue"))
	fmt.Printf("  State:      %s\n", state)
	fmt.Printf("  Created:    %s\n", run.Colorize(date, "dim"))
	if net.Site != "" {
		fmt.Printf("  Site:       %s\n", run.Colorize(mapSiteToPublic(net.Site), "cyan"))
	}
	fmt.Printf("  Speed:      %d %s\n", net.Speed.Value, net.Speed.Unit)
	fmt.Printf("  Ports:      %d\n", net.Ports)
	fmt.Printf("  Max MTU:    %d\n", net.MaxMTU)
	connectedStr := fmt.Sprintf("%d servers", len(net.Interfaces))
	if run.Color && len(net.Interfaces) > 0 {
		connectedStr = run.Colorize(connectedStr, "green")
	}
	fmt.Printf("  Connected:  %s\n", connectedStr)

	// DRP Section
	if net.Drp != nil {
		fmt.Printf("%s\n", run.Colorize("Disaster Recovery Plan (DRP):", "cyan"))
		if net.Drp.Enabled {
			fmt.Printf("  Enabled:    %s\n", run.Colorize("Yes", "green"))
			if net.Drp.Site != "" {
				fmt.Printf("  Target Site: %s\n", mapSiteToPublic(net.Drp.Site))
			}
		} else {
			fmt.Printf("  Enabled:    %s\n", run.Colorize("No", "yellow"))
		}
	}

	if len(net.Interfaces) > 0 {
		fmt.Printf("  Servers:\n")
		for _, iface := range net.Interfaces {
			// Try to get primary IP
			serverIP := ""
			for _, subItem := range iface.Server.Items.MAC.SubItems {
				if subItem.Primary && subItem.IP != nil && subItem.IP.Version == 4 {
					serverIP = subItem.IP.Address
					break
				}
			}
			serverName := iface.Server.Name
			if run.Color {
				serverName = run.Colorize(serverName, "cyan")
				if serverIP != "" {
					serverIP = run.Colorize(serverIP, "yellow")
				}
			}
			if serverIP != "" {
				fmt.Printf("    - %s (%s) [%s]\n", serverName, serverIP, iface.MAC)
			} else {
				fmt.Printf("    - %s [%s]\n", serverName, iface.MAC)
			}
		}
	}
}

func (run *RunMiddleware) NetworkRename(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	networkOID, _ := cmd.Flags().GetString("network-oid")
	name, _ := cmd.Flags().GetString("name")

	netRename := &api.NetworkRename{Name: name}
	_, apiReturn, err := run.API.SendRequestToAPI(api.HTTPPut, "/network/switch/"+networkOID, netRename)
	// Handle error, print output
	run.handleErrorAndGenericOutput(apiReturn, err)
}

// toNetworkDetailOutput converts API NetworkDetail to clean output struct
func toNetworkDetailOutput(net *api.NetworkDetail) networkDetailOutput {
	interfaces := make([]networkInterfaceOutput, 0, len(net.Interfaces))
	for _, iface := range net.Interfaces {
		// Extract primary IP if available
		serverIP := ""
		for _, subItem := range iface.Server.Items.MAC.SubItems {
			if subItem.Primary && subItem.IP != nil {
				serverIP = subItem.IP.Address
				break
			}
		}
		interfaces = append(interfaces, networkInterfaceOutput{
			OID:        iface.OID,
			MAC:        iface.MAC,
			ServerOID:  iface.Server.OID,
			ServerName: iface.Server.Name,
			ServerIP:   serverIP,
		})
	}

	createdAt := ""
	if net.CreatedAt != nil {
		createdAt = DatePtrFormat(net.CreatedAt)
	}

	return networkDetailOutput{
		OID:        net.OID,
		Name:       net.Name,
		State:      net.State,
		CreatedAt:  createdAt,
		Ports:      net.Ports,
		MaxMTU:     net.MaxMTU,
		Speed:      net.Speed,
		Interfaces: interfaces,
	}
}
