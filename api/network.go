package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"time"
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

func (API *APITitan) NetworkList(cmd *cobra.Command, args []string) {
	companyUUID, _ := cmd.Flags().GetString("company-uuid")
	API.ParseGlobalFlags(cmd)

	err := API.SendAndResponse(HTTPGet, "/compute/networks?company_uuid="+companyUUID, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		networks := &APINetworkList{}
		if err := json.Unmarshal(API.RespBody, &networks); err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("Quota:", networks.Quota)
		for _, net := range networks.NetInfos {
			API.NetworkPrintBase(&net)
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

type APINetworkOps struct {
	ServerUUID string `json:"server_uuid"`
}

func (API *APITitan) NetworkAttachServer(cmd *cobra.Command, args []string) {
	networkUUID := args[0]
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	act := APINetworkOps{
		ServerUUID: serverUUID,
	}
	API.SendAndPrintDefaultReply(HTTPPut, "/compute/networks/"+networkUUID+"/attach", act)
}

func (API *APITitan) NetworkDetachServer(cmd *cobra.Command, args []string) {
	networkUUID := args[0]
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	act := APINetworkOps{
		ServerUUID: serverUUID,
	}
	API.SendAndPrintDefaultReply(HTTPPut, "/compute/networks/"+networkUUID+"/detach", act)
}

func (API *APITitan) NetworkCreate(cmd *cobra.Command, args []string) {
	companyUUID := args[0]
	API.ParseGlobalFlags(cmd)
	networkName, _ := cmd.Flags().GetString("name")
	cidr, _ := cmd.Flags().GetString("cidr")

	net := APINetworkCreate{
		MaxMTU: 8948,
		Name:   networkName,
		Ports:  6,
		CIDR:   cidr,
	}
	net.Speed.Value = 1
	net.Speed.Unit = "Gbps"

	reqData, e := json.Marshal(net)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	err := API.SendAndResponse(HTTPPost, "/compute/networks/?company_uuid="+companyUUID, reqData)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		net := &APINetwork{}
		if err := json.Unmarshal(API.RespBody, net); err != nil {
			fmt.Println(err.Error())
			return
		}
		API.NetworkPrintBase(net)
	}
}

func (API *APITitan) NetworkRemove(cmd *cobra.Command, args []string) {
	_ = cmd
	networkUUID := args[0]
	API.SendAndPrintDefaultReply(HTTPDelete, "/compute/networks/"+networkUUID, nil)
}

func (API *APITitan) DateFormat(timestamp int64) string {
	dateMls := time.Unix(0, int64(timestamp)*int64(time.Millisecond))
	date := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		dateMls.Year(), dateMls.Month(), dateMls.Day(),
		dateMls.Hour(), dateMls.Minute(), dateMls.Second())
	return date
}

func (API *APITitan) NetworkPrintBase(net *APINetwork) {
	date := API.DateFormat(net.CreatedAt)
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
			fmt.Printf("  %s\n", net.Gateway)
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

func (API *APITitan) NetworkRename(cmd *cobra.Command, args []string) {
	networkUUID := args[0]
	name, _ := cmd.Flags().GetString("name")

	netRename := &APINetworkRename{Name: name}
	API.SendAndPrintDefaultReply(HTTPPut, "/compute/networks/"+networkUUID, netRename)
}

func (API *APITitan) NetworkSetGateway(cmd *cobra.Command, args []string) {
	networkUUID := args[0]
	ip, _ := cmd.Flags().GetString("ip")

	ipData := APIIP{
		IP:      ip,
		Version: 4,
	}
	API.SendAndPrintDefaultReply(HTTPPut, "/compute/networks/"+networkUUID+"/gateway", ipData)
}

func (API *APITitan) NetworkUnsetGateway(cmd *cobra.Command, args []string) {
	_ = cmd
	networkUUID := args[0]
	API.SendAndPrintDefaultReply(HTTPDelete, "/compute/networks/"+networkUUID+"/gateway", nil)
}
