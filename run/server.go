package run

import (
	"fmt"
	"os"
	"text/tabwriter"
	"titan-sc/api"

	"github.com/spf13/cobra"
)

const (
	// Disk source type
	DiskSourceTemplate = "template"
	DiskSourceImage    = "image"
)

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
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		_, _ = fmt.Fprintf(w, "COMPANY\tUUID\tPLAN\tSTATE\tOS\tNAME\tMANAGED\t\n")

		for _, server := range servers {
			var osInfos string
			state := GetStateColorized(run.Color, server.State)

			if server.Disksource.Type == DiskSourceTemplate {
				osInfos = fmt.Sprintf("%s %s", server.Template.OS, server.Template.Version)
			} else {
				osInfos = fmt.Sprintf("%s %s", server.Image.OS.Name, server.Image.OS.Version)
			}
			_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%t\t\n",
				server.Company.Name, server.UUID, server.Plan, state, osInfos,
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
	var osInfos string
	date := millisecondsToTime(server.CreationDate).Format("2006-01-02 15:04:05")

	if server.Disksource.Type == DiskSourceTemplate {
		osInfos = fmt.Sprintf("%s %s", server.Template.OS, server.Template.Version)
	} else {
		osInfos = fmt.Sprintf("%s %s", server.Image.OS.Name, server.Image.OS.Version)
	}

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
		server.Plan, osInfos, server.Company.Name,
		server.Hypervisor.Hostname, server.KvmIp.Status)

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

	if server.Notes == "" {
		fmt.Println("Notes: -")
	} else {
		fmt.Println("Notes:", server.Notes)
	}
}
