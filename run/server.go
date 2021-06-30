package run

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
	"titan-sc/api"
)

const (
	deleteServerReasonCLI = "cli"
)

func (run *RunMiddleware) ServerChangeName(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	newServerName, _ := cmd.Flags().GetString("name")
	apiReturn, err := run.API.ServerChangeName(newServerName, serverUUID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerChangeReverse(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	newServerReverse, _ := cmd.Flags().GetString("reverse")
	apiReturn, err := run.API.ServerChangeReverse(newServerReverse, serverUUID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	companyUUID, _ := cmd.Flags().GetString("company-uuid")

	servers, apiReturn, err := run.API.ServerList(companyUUID)
	if err != nil || apiReturn != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}

	if !run.HumanReadable {
		printAsJson(servers)
	} else {
		var w *tabwriter.Writer
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "COMPANY\tUUID\tPLAN\tSTATE\tOS\tNAME\tMANAGED\t\n")

		for _, server := range servers {
			state := serverStateSetColor(run.Color, server.State)
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%t\t\n",
				server.Company.Name, server.UUID, server.Plan, state, server.Template,
				server.Name, server.Managed)
		}
		_ = w.Flush()
	}
}

func (run *RunMiddleware) ServerDetail(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	server, apiReturn, err := run.API.GetServerUUID(serverUUID)
	if err != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}
	if !run.HumanReadable {
		printAsJson(server)
	} else {
		run.printServerDetail(server)
	}
}

func (run *RunMiddleware) printServerDetail(server *api.APIServer) {
	date := run.DateFormat(server.CreationDate)
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
		server.Company.Name, server.Hypervisor.Hostname, server.KvmIp.Status)

	if server.KvmIp.Status == "started" && server.KvmIp.URI != "" {
		fmt.Println("IP Kvm URI:", server.KvmIp.URI)
	}

	// Collect server IPs in a single string
	ips := ""
	for _, ip := range server.IPs {
		ips += ip.IP + " "
	}

	fmt.Println("Managed network uuid:", server.ManagedNetwork)
	fmt.Printf("Network:\n"+
		"  - IPv4: %s\n"+
		"  - IPv6: %s\n"+
		"  - Mac: %s\n"+
		"  - Gateway: %s\n"+
		"  - Bandwidth in/out: %d/%d %s\n"+
		"  - Reverse: %s\n",
		ips, server.IPv6, server.Mac, server.Gateway,
		server.Bandwidth.Input, server.Bandwidth.Output,
		server.Bandwidth.Unit, server.Reverse)

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

func (run *RunMiddleware) ServerStart(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	apiReturn, err := run.API.ServerStateAction("start", serverUUID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerStop(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	apiReturn, err := run.API.ServerStateAction("stop", serverUUID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerRestart(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	apiReturn, err := run.API.ServerStateAction("reboot", serverUUID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerHardstop(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	apiReturn, err := run.API.ServerStateAction("hardstop", serverUUID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerLoadISO(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	uriISO, _ := cmd.Flags().GetString("uri")
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	apiReturn, err := run.API.ServerLoadISO(uriISO, serverUUID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerUnloadISO(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	apiReturn, err := run.API.ServerUnloadISO(serverUUID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerListTemplates(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	templates, apiReturn, err := run.API.ServerListTemplates()
	// Render error output
	if err != nil || apiReturn != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}

	if !run.HumanReadable {
		printAsJson(templates)
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

func (run *RunMiddleware) ServerDelete(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	apiReturn, err := run.API.ServerDelete(serverUUID, deleteServerReasonCLI)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerReset(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")

	serverReset := &api.APIResetServer{}
	sshKeys, _ := cmd.Flags().GetString("ssh-keys-name")
	serverReset.UserPassword, _ = cmd.Flags().GetString("password")
	serverReset.TemplateOS, _ = cmd.Flags().GetString("os")
	serverReset.TemplateVersion, _ = cmd.Flags().GetString("os-version")

	if sshKeys != "" {
		var err error
		serverReset.UserSSHKeys, err = run.serverSearchSSHKeys(sshKeys)
		if err != nil {
			return
		}
	}
	apiReturn, err := run.API.ServerReset(serverUUID, serverReset)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) serverSearchSSHKeys(sshKeysName string) ([]string, error) {
	sshKeysList, err := run.API.GetSSHKeyList()
	if err != nil {
		return []string{}, err
	}
	sshKeys := make([]string, 0)
	for _, keyRequest := range strings.Split(strings.TrimSpace(sshKeysName), ",") {
		found := false
		for _, keyExist := range sshKeysList {
			if keyExist.Title == keyRequest {
				sshKeys = append(sshKeys, keyExist.Title)
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

func (run *RunMiddleware) ServerCreate(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	server := &api.APICreateServers{
		CreateServersDetail: make([]api.CreateServersDetail, 1),
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

	var allAddons []api.APIInstallAddonsAddon
	if DiskAddonsNumber > 0 || RamAddonsNumber > 0 || CpuAddonsNumber > 0 {
		addons, err := run.GetAllAddons()
		if err != nil {
			return
		}

		if DiskAddonsNumber > 0 {
			if err := run.ServerCreateAddAddonInArray(addons, &allAddons, DiskAddonsNumber, "disk"); err != nil {
				run.OutputError(err)
				return
			}
		}
		if RamAddonsNumber > 0 {
			if err := run.ServerCreateAddAddonInArray(addons, &allAddons, RamAddonsNumber, "ram"); err != nil {
				run.OutputError(err)
				return
			}
		}
		if CpuAddonsNumber > 0 {
			if err := run.ServerCreateAddAddonInArray(addons, &allAddons, CpuAddonsNumber, "cpu"); err != nil {
				run.OutputError(err)
				return
			}
		}
	}
	server.CreateServersDetail[0].Addons = allAddons

	if sshKeys != "" {
		var err error
		server.CreateServersDetail[0].UserSSHKeys, err = run.serverSearchSSHKeys(sshKeys)
		if err != nil {
			run.OutputError(errors.New("Fail to get ssh keys."))
			return
		}
	}

	apiReturn, err := run.API.ServerCreate(server)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerCreateAddAddonInArray(addonsList []api.APIAddonsItem,
	allAddons *[]api.APIInstallAddonsAddon, addonNumber int, addonName string) error {
	addonUUID, err := run.AddonGetUUIDByName(addonsList, addonName)
	if err != nil {
		run.OutputError(err)
		return err
	}
	addonItem := &api.APIInstallAddonsAddon{
		Item:     addonUUID,
		Quantity: int64(addonNumber),
	}
	*allAddons = append(*allAddons, *addonItem)
	return nil
}
