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

	if compagnies, err := API.GetCompagnies(); err != nil {
		fmt.Println(err.Error())
	} else {
		var w *tabwriter.Writer
		if API.HumanReadable {
			w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		}

		for _, company := range compagnies {
			if companyUUID == "" || (companyUUID != "" && companyUUID == company.Company.UUID) {
				servers, err := API.GetCompanyServers(company.Company.UUID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				if !API.HumanReadable {
					API.PrintJson()
				} else {
					_, _ = fmt.Fprintf(w, "UUID\tPLAN\tSTATE\tOS\tNAME\tMANAGED\t\n")
					for _, server := range servers {
						state := API.ServerStateSetColor(server.State)
						_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%t\t\n",
							server.UUID, server.Plan, state, server.OS, server.Name, server.Managed)
					}
					_ = w.Flush()
				}
			}
		}
	}
}

func (API *APITitan) ServerDetail(cmd *cobra.Command, args []string) {

	serverUUID := args[0]
	API.ParseGlobalFlags(cmd)

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

	serverUUID := args[0]
	API.ParseGlobalFlags(cmd)

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
	API.ParseGlobalFlags(cmd)
	API.ServerStateAction("start", args[0])
}

func (API *APITitan) ServerStop(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	API.ServerStateAction("stop", args[0])
}

func (API *APITitan) ServerRestart(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	API.ServerStateAction("reboot", args[0])
}

func (API *APITitan) ServerHardstop(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	API.ServerStateAction("hardstop", args[0])
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
	API.ParseGlobalFlags(cmd)
	serverUUID := args[0]
	path := "/compute/servers/" + serverUUID + "/iso"
	API.SendAndPrintDefaultReply(HTTPDelete, path, nil)
}

func (API *APITitan) ServerCreate(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)

	server := &APICreateServers{}
	sshKeys, _ := cmd.Flags().GetString("ssh-keys-name")
	server.Quantity, _ = cmd.Flags().GetInt64("quantity")
	server.UserPassword, _ = cmd.Flags().GetString("password")
	server.UserLogin, _ = cmd.Flags().GetString("login")
	server.TemplateOS, _ = cmd.Flags().GetString("os")
	server.TemplateVersion, _ = cmd.Flags().GetString("os-version")
	server.ManagedNetwork, _ = cmd.Flags().GetString("network-uuid")
	server.Plan, _ = cmd.Flags().GetString("plan")
	CpuAddonsNumber, _ := cmd.Flags().GetInt("cpu-addon")
	RamAddonsNumber, _ := cmd.Flags().GetInt("ram-addon")
	DiskAddonsNumber, _ := cmd.Flags().GetInt("disk-addon")

	server.Plan = strings.ToUpper(server.Plan)
	switch server.Plan {
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
		fmt.Println("Invalid --plan value !\nHelper: argument --plan only value accepted: SC1 SC2 SC3")
		return
	}

	allAddons := []APIInstallAddonsAddon{}
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
	server.Addons = allAddons

	if sshKeys != "" {
		var err error
		server.UserSSHKey, err = API.ServerSearchAndConcatSSHKeys(sshKeys)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	servers := make([]*APICreateServers, 0)
	servers = append(servers, server)
	API.SendAndPrintDefaultReply(HTTPPost, "/compute/servers", servers)
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
		return fmt.Errorf("Error: %s max CPU addons is %d", plan, cpuMax)
	}

	if ramAddonsNumber > ramMax {
		return fmt.Errorf("Error: %s max RAM addons is %d", plan, ramMax)
	}

	if diskAddonsNumber > diskMax {
		return fmt.Errorf("Error: %s max disk size addons is %d", plan, diskMax)
	}
	return nil
}

func (API *APITitan) ServerDelete(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)
	serverUUID := args[0]
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

	templates := []APITemplateFullInfos{}
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
	API.ParseGlobalFlags(cmd)
	serverUUID := args[0]

	serverReset := &APIResetServer{}
	sshKeys, _ := cmd.Flags().GetString("ssh-keys-name")
	serverReset.UserPassword, _ = cmd.Flags().GetString("password")
	serverReset.TemplateOS, _ = cmd.Flags().GetString("os")
	serverReset.TemplateVersion, _ = cmd.Flags().GetString("os-version")

	if sshKeys != "" {
		var err error
		serverReset.UserSSHKey, err = API.ServerSearchAndConcatSSHKeys(sshKeys)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	API.SendAndPrintDefaultReply(HTTPPut, "/compute/servers/"+serverUUID+"/reset", serverReset)
}

func (API *APITitan) ServerSearchAndConcatSSHKeys(sshKeys string) (string, error) {
	sshKeysList, err := API.SSHKeysGet()
	if err != nil {
		return "", err
	}

	sshKey := ""
	for _, keyRequest := range strings.Split(sshKeys, ",") {
		find := false
		for _, keyExist := range sshKeysList {
			if keyExist.Title == keyRequest {
				sshKey += keyExist.Content + "\n"
				find = true
				break
			}
		}

		// Check if ssh key find
		if !find {
			return "", fmt.Errorf("SSH keys name <%s> not found\n", keyRequest)
		}
	}
	return sshKeys, nil
}
