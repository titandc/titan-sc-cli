package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	// Max cpu addons by plan
	SC1MaxCpuAddons = 6
	SC2MaxCpuAddons = 12
	SC3MaxCpuAddons = 32

	// Max ram addons by plan
	SC1MaxRamAddons = 8
	SC2MaxRamAddons = 32
	SC3MaxRamAddons = 256

	// Max disk addons by plan
	SC1MaxDiskAddons = 2000
	SC2MaxDiskAddons = 2000
	SC3MaxDiskAddons = 2000
)

/*
 *
 *
 ******************
 * Servers function
 ******************
 *
 *
 */
func (API *APITitan) ServerChangeName(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	newName, _ := cmd.Flags().GetString("name")

	updateInfos := &APIServerUpdateInfos{
		Name: newName,
	}
	path := "/compute/servers/" + serverUUID
	API.SendAndPrintDefaultReply(HTTPPut, path, updateInfos)
}

func (API *APITitan) ServerChangeReverse(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	newReverse, _ := cmd.Flags().GetString("reverse")

	updateInfos := &APIServerUpdateInfos{
		Reverse: newReverse,
	}
	path := "/compute/servers/" + serverUUID
	API.SendAndPrintDefaultReply(HTTPPut, path, updateInfos)
}

func (API *APITitan) ServerList(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	companyUUID, _ := cmd.Flags().GetString("company-uuid")

	err := API.SendAndResponse(HTTPGet, "/compute/servers/detail?company_uuid="+companyUUID, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		var servers []APIServer

		if err := json.Unmarshal(API.RespBody, &servers); err != nil {
			fmt.Println(err.Error())
			return
		}

		var w *tabwriter.Writer
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "COMPANY\tUUID\tPLAN\tSTATE\tOS\tNAME\tMANAGED\t\n")

		for _, server := range servers {
			state := API.ServerStateSetColor(server.State)
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%t\t\n",
				server.CompanyName, server.UUID, server.Plan, state, server.Template,
				server.Name, server.Managed)
		}
		_ = w.Flush()
	}
}

func (API *APITitan) ServerDetail(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	server, err := API.GetServerUUID(serverUUID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !API.HumanReadable {
		API.PrintJson()
	} else {
		API.PrintServerDetail(server)
	}
}

func (API *APITitan) ShowServerDetail(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	if server, err := API.GetServerUUID(serverUUID); err != nil {
		fmt.Println(err)
	} else {
		if !API.HumanReadable {
			API.PrintJson()
		} else {
			API.PrintServerDetail(server)
		}
	}
}

func (API *APITitan) ServerStart(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	API.ServerStateAction("start", serverUUID)
}

func (API *APITitan) ServerStop(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	API.ServerStateAction("stop", serverUUID)
}

func (API *APITitan) ServerRestart(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	API.ServerStateAction("reboot", serverUUID)
}

func (API *APITitan) ServerHardstop(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	API.ServerStateAction("hardstop", serverUUID)
}

func (API *APITitan) ServerStateAction(state, serverUUID string) {
	// check server exist
	server, err := API.GetServerUUID(serverUUID)
	if err != nil {
		fmt.Println(err)
		return
	}

	// send request
	act := APIServerAction{
		Action: state,
	}
	path := "/compute/servers/" + server.UUID + "/action"
	API.SendAndPrintDefaultReply(HTTPPut, path, act)
}

func (API *APITitan) GetServerUUID(serverUUID string) (*APIServer, error) {
	err := API.SendAndResponse(HTTPGet, "/compute/servers/"+serverUUID, nil)
	if err != nil {
		return nil, err
	}

	server := &APIServer{}
	if err := json.Unmarshal(API.RespBody, &server); err != nil {
		return nil, err
	}
	return server, nil
}

func (API *APITitan) PrintServerDetail(server *APIServer) {
	date := API.DateFormat(server.Creationdate)
	fmt.Printf("Name: %s\n"+
		"UUID: %s\n"+
		"Created at: %s\n"+
		"VM Login: %s\n"+
		"State: %s\n"+
		"Plan: %s\n"+
		"OS version: %s\n"+
		"Company: %s\n"+
		"Hypervisor: %s\n"+
		"IP Kvm: %s\n",
		server.Name, server.UUID, date, server.Login, server.State,
		server.Plan, server.Template,
		server.CompanyName, server.Hypervisor, server.KvmIp.Status)

	if server.KvmIp.Status == "started" && server.KvmIp.URI != "" {
		fmt.Println("IP Kvm URI:", server.KvmIp.URI)
	}

	fmt.Println("Managed network uuid:", server.ManagedNetwork)
	fmt.Printf("Network:\n"+
		"  - IPv4: %s\n"+
		"  - IPv6: %s\n"+
		"  - Mac: %s\n"+
		"  - Gateway: %s\n"+
		"  - Bandwidth in/out: %d/%d %s\n"+
		"  - Reverse: %s\n",
		server.IP, server.IPv6, server.Mac, server.Gateway,
		server.Bandwidth.Input, server.Bandwidth.Output,
		server.Bandwidth.Uint, server.Reverse)

	fmt.Printf("Resources:\n"+
		"  - Cpu(s): %d\n"+
		"  - RAM: %d %s\n"+
		"  - Disk: %d %s\n"+
		"  - Disk QoS Read/Write: %d/%d %s\n"+
		"  - Disk IOPS Read/Write/BlockSize: %d/%d/%s %s\n",
		server.CPU.NbCores, server.RAM.Value, server.RAM.Unit,
		server.Disk.Size.Value, server.Disk.Size.Unit,
		server.Disk.QoS.Read, server.Disk.QoS.Write, server.Disk.QoS.Unit,
		server.Disk.IOPS.Read, server.Disk.IOPS.Write, server.Disk.IOPS.BlockSize,
		server.Disk.IOPS.Unit)

	if len(server.PendingActions) > 0 {
		fmt.Println("Pending actions:")
		for _, action := range server.PendingActions {
			fmt.Printf("  - %s\n", action)
		}
	} else {
		fmt.Println("Pending action(s): -")
	}

	if len(server.PendingActions) > 0 {
		fmt.Println("Pending actions:")
		for _, action := range server.PendingActions {
			fmt.Printf("  - %s\n", action)
		}
	} else {
		fmt.Println("Pending action(s): -")
	}

	if server.Notes == "" {
		fmt.Println("Notes: -")
	} else {
		fmt.Println("Notes:", server.Notes)
	}
}

func (API *APITitan) ServerStateSetColor(state string) string {
	if !API.Color {
		return state
	}

	colorState := state
	switch state {
	case "deleted":
		colorState = "\033[1;31m" + state + "\033[0m"
	case "started":
		colorState = "\033[1;32m" + state + "\033[0m"
	case "stopped":
		colorState = "\033[1;33m" + state + "\033[0m"
	}
	return colorState
}

func (API *APITitan) ServerLoadISO(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	uriISO, _ := cmd.Flags().GetString("uri")
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	reqStruct := &APIServerLOadISORequest{
		Protocol: "https",
		ISO:      uriISO,
	}
	path := "/compute/servers/" + serverUUID + "/iso"
	API.SendAndPrintDefaultReply(HTTPPost, path, reqStruct)
}

func (API *APITitan) ServerUnloadISO(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	path := "/compute/servers/" + serverUUID + "/iso"
	API.SendAndPrintDefaultReply(HTTPDelete, path, nil)
}

func (API *APITitan) ServerCreate(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)

	server := &APICreateServers{
		CreateServersDetail: make([]CreateServersDetail, 1),
	}

	sshKeys, _ := cmd.Flags().GetString("ssh-keys-name")
	server.CreateServersDetail[0].Quantity, _ = cmd.Flags().GetInt64("quantity")
	server.CreateServersDetail[0].UserPassword, _ = cmd.Flags().GetString("password")
	server.CreateServersDetail[0].UserLogin, _ = cmd.Flags().GetString("login")
	server.CreateServersDetail[0].TemplateOS, _ = cmd.Flags().GetString("os")
	server.CreateServersDetail[0].TemplateVersion, _ = cmd.Flags().GetString("os-version")
	server.CreateServersDetail[0].ManagedNetwork, _ = cmd.Flags().GetString("network-uuid")
	server.CreateServersDetail[0].Plan, _ = cmd.Flags().GetString("plan")
	CpuAddonsNumber, _ := cmd.Flags().GetInt("cpu-addon")
	RamAddonsNumber, _ := cmd.Flags().GetInt("ram-addon")
	DiskAddonsNumber, _ := cmd.Flags().GetInt("disk-addon")

	server.CreateServersDetail[0].Plan = strings.ToUpper(server.CreateServersDetail[0].Plan)
	server.CreateServersDetail[0].TemplateOS = strings.ToTitle(server.CreateServersDetail[0].TemplateOS)
	
	switch server.CreateServersDetail[0].Plan {
	case "SC1":
		if err := API.ServerCheckAddonsNumber(CpuAddonsNumber, RamAddonsNumber, DiskAddonsNumber,
			SC1MaxCpuAddons, SC1MaxRamAddons, SC1MaxDiskAddons, "SC1"); err != nil {
			fmt.Println(err.Error())
			return
		}
	case "SC2":
		if err := API.ServerCheckAddonsNumber(CpuAddonsNumber, RamAddonsNumber, DiskAddonsNumber,
			SC2MaxCpuAddons, SC2MaxRamAddons, SC2MaxDiskAddons, "SC2"); err != nil {
			fmt.Println(err.Error())
			return
		}
	case "SC3":
		if err := API.ServerCheckAddonsNumber(CpuAddonsNumber, RamAddonsNumber, DiskAddonsNumber,
			SC3MaxCpuAddons, SC3MaxRamAddons, SC3MaxDiskAddons, "SC3"); err != nil {
			fmt.Println(err.Error())
			return
		}
	default:
		fmt.Println("Invalid --plan value !\nHelper: argument --plan only value accepted: SC1 SC2 SC3.")
		return
	}

	var allAddons []APIInstallAddonsAddon
	if DiskAddonsNumber > 0 || RamAddonsNumber > 0 || CpuAddonsNumber > 0 {
		addons, err := API.GetAllAddons()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if DiskAddonsNumber > 0 {
			if err := API.ServerCreateAddAddonInArray(addons, &allAddons, DiskAddonsNumber, "disk"); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		if RamAddonsNumber > 0 {
			if err := API.ServerCreateAddAddonInArray(addons, &allAddons, RamAddonsNumber, "ram"); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		if CpuAddonsNumber > 0 {
			if err := API.ServerCreateAddAddonInArray(addons, &allAddons, CpuAddonsNumber, "cpu"); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
	server.CreateServersDetail[0].Addons = allAddons

	if sshKeys != "" {
		var err error
		server.CreateServersDetail[0].UserSSHKeys, err = API.ServerSearchSSHKeys(sshKeys)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	API.SendAndPrintDefaultReply(HTTPPost, "/compute/servers", server)
}

func (API *APITitan) ServerCreateAddAddonInArray(addonsList []APIAddonsItem,
	allAddons *[]APIInstallAddonsAddon, addonNumber int, addonName string) error {
	addonUUID, err := API.AddonGetUUIDByName(addonsList, addonName)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	addonItem := &APIInstallAddonsAddon{
		Item:     addonUUID,
		Quantity: int64(addonNumber),
	}
	*allAddons = append(*allAddons, *addonItem)
	return nil
}

func (API *APITitan) ServerCheckAddonsNumber(cpuAddonsNumber, ramAddonsNumber, diskAddonsNumber,
	cpuMax, ramMax, diskMax int, plan string) error {
	if cpuAddonsNumber > cpuMax {
		return fmt.Errorf("Error: %s max CPU addons is %d.", plan, cpuMax)
	}

	if ramAddonsNumber > ramMax {
		return fmt.Errorf("Error: %s max RAM addons is %d.", plan, ramMax)
	}

	if diskAddonsNumber > diskMax {
		return fmt.Errorf("Error: %s max disk size addons is %d.", plan, diskMax)
	}
	return nil
}

func (API *APITitan) ServerDelete(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	API.SendAndPrintDefaultReply(HTTPDelete, "/compute/servers/"+serverUUID, nil)
}

func (API *APITitan) ServerGetTemplateList(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)

	templates, err := API.ServerGetTemplates()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !API.HumanReadable {
		API.PrintJson()
	} else {
		for _, template := range templates {
			str := "Name: " + template.OS + " version: "
			for _, version := range template.Versions {
				str += version.Version + " - "
			}
			str = strings.TrimSuffix(str, " - ")
			fmt.Println(str)
		}
	}
}

func (API *APITitan) ServerGetTemplates() ([]APITemplateFullInfos, error) {
	err := API.SendAndResponse(HTTPGet, "/compute/templates", nil)
	if err != nil {
		return nil, err
	}

	var templates []APITemplateFullInfos
	if API.HumanReadable {
		err := json.Unmarshal(API.RespBody, &templates)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
	}
	return templates, nil
}

func (API *APITitan) ServerReset(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	serverReset := &APIResetServer{}
	sshKeys, _ := cmd.Flags().GetString("ssh-keys-name")
	serverReset.UserPassword, _ = cmd.Flags().GetString("password")
	serverReset.TemplateOS, _ = cmd.Flags().GetString("os")
	serverReset.TemplateVersion, _ = cmd.Flags().GetString("os-version")

	if sshKeys != "" {
		var err error
		serverReset.UserSSHKeys, err = API.ServerSearchSSHKeys(sshKeys)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	API.SendAndPrintDefaultReply(HTTPPut, "/compute/servers/"+serverUUID+"/reset", serverReset)
}

func (API *APITitan) ServerSearchSSHKeys(sshKeysName string) ([]string, error) {
	sshKeysList, err := API.SSHKeysGet()
	if err != nil {
		return []string{}, err
	}
	sshKeys := make([]string, 0)
	for _, keyRequest := range strings.Split(strings.TrimSpace(sshKeysName), ",") {
		found := false
		for _, keyExist := range sshKeysList {
			if keyExist.Title == keyRequest {
				sshKeys = append(sshKeys, keyExist.Content)
				found = true
				break
			}
		}
		if !found {
			return []string{}, fmt.Errorf("SSH keys name <%s> not found.\n", keyRequest)
		}
	}
	return sshKeys, nil
}
