package run

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"sort"

	"titan-sc/api"

	"github.com/spf13/cobra"
)

func (run *RunMiddleware) IPAttach(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	ip, _ := cmd.Flags().GetString("ip")

	parsedIP := net.ParseIP(ip)
	apiReturn, err := run.API.IPAttach(serverOID, []string{parsedIP.String()})
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) IPDetach(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	ip, _ := cmd.Flags().GetString("ip")

	parsedIP := net.ParseIP(ip)
	apiReturn, err := run.API.IPDetach(serverOID, []string{parsedIP.String()})
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) IPsCompanyList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	companyOID, err := run.ResolveCompanyOID(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}

	ipList, err := run.API.GetCompanyIPList(companyOID)
	if err != nil {
		run.OutputError(err)
		return
	}
	if run.JSONOutput {
		printAsJson(ipList)
	} else {
		run.IPsPrint(ipList)
	}
}

func (run *RunMiddleware) IPUpdateReverse(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	argIP, _ := cmd.Flags().GetString("ip")
	newIPReverse, _ := cmd.Flags().GetString("reverse")

	user, err := run.API.GetUserInfos()
	if err != nil {
		run.OutputError(err)
	}

	ips, err := run.API.GetCompanyIPList(user.DefaultCompanyOID)
	if err != nil {
		run.OutputError(err)
	}

	for _, ip := range ips {
		if ip.Address == argIP {
			apiReturn, err := run.API.IPUpdateReverse(ip.OID, newIPReverse)
			run.handleErrorAndGenericOutput(apiReturn, err)
			return
		}
	}

	run.OutputError(errors.New("IP not found"))
}

func (run *RunMiddleware) IPsPrint(ipArray []api.IP) {
	if len(ipArray) == 0 {
		fmt.Println("Empty IPs list")
		return
	}

	// Sort IPs: IPv4 first (sorted), then IPv6 (sorted)
	sort.Slice(ipArray, func(i, j int) bool {
		ipI := net.ParseIP(ipArray[i].Address)
		ipJ := net.ParseIP(ipArray[j].Address)

		// IPv4 before IPv6
		isIPv4i := ipI.To4() != nil
		isIPv4j := ipJ.To4() != nil
		if isIPv4i != isIPv4j {
			return isIPv4i // IPv4 comes first
		}

		// Same version, sort by IP bytes
		return bytes.Compare(ipI, ipJ) < 0
	})

	w := NewTable("IP", "REVERSE", "SERVER")
	for _, ip := range ipArray {
		serverName := ip.ServerName
		if serverName == "" {
			serverName = "(unassigned)"
		}

		var serverColorFn func(string) string
		if run.Color {
			if ip.ServerName != "" {
				serverColorFn = ColorFn("cyan") // Names are always cyan
			} else {
				serverColorFn = ColorFn("dim") // Unassigned is dim
			}
		}

		w.AddRow(
			ColIP(ip.Address),
			Col(ip.Reverse),
			ColColor(serverName, serverColorFn),
		)
	}
	w.Print()
}
