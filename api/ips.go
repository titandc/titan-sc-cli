package api

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

func (API *APITitan) IPAttach(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	ip, _ := cmd.Flags().GetString("ip")

	ipOpt := APIIP{
		IP:      ip,
		Version: 4,
	}
	API.SendAndPrintDefaultReply(HTTPPost, "/compute/servers/"+serverUUID+"/ips", ipOpt)
}

func (API *APITitan) IPDetach(cmd *cobra.Command, args []string) {
	_ = args
	API.ParseGlobalFlags(cmd)
	serverUUID, _ := cmd.Flags().GetString("server-uuid")
	ip, _ := cmd.Flags().GetString("ip")

	ipOpt := APIIP{
		IP:      ip,
		Version: 4,
	}
	API.SendAndPrintDefaultReply(HTTPDelete, "/compute/servers/"+serverUUID+"/ips", ipOpt)
}

func (API *APITitan) IPsList(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)

	err := API.SendAndResponse(HTTPGet, "/compute/ips", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !API.HumanReadable {
		API.PrintJson()
	} else {
		APIIP := make([]APIIP, 0)
		if err := json.Unmarshal(API.RespBody, &APIIP); err != nil {
			fmt.Println(err.Error())
			return
		}
		API.IPsPrint(&APIIP)
	}
}

func (API *APITitan) IPsCompanyList(cmd *cobra.Command, args []string) {
	API.ParseGlobalFlags(cmd)

	companyUUID := args[0]
	err := API.SendAndResponse(HTTPGet, "/companies/"+companyUUID+"/ips", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	if !API.HumanReadable {
		API.PrintJson()
	} else {
		APIIP := make([]APIIP, 0)
		if err := json.Unmarshal(API.RespBody, &APIIP); err != nil {
			fmt.Println(err.Error())
			return
		}
		API.IPsPrint(&APIIP)
	}
}

func (API *APITitan) IPsPrint(ipArray *[]APIIP) {

	if len(*ipArray) == 0 {
		fmt.Println("Empty IPs list")
		return
	}

	var w *tabwriter.Writer
	if API.HumanReadable {
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	}

	_, _ = fmt.Fprintf(w, "IP\tVERSION\t\n")
	for _, ip := range *ipArray {
		_, _ = fmt.Fprintf(w, "%s\t%d\t\n", ip.IP, ip.Version)
	}
	_ = w.Flush()
}
