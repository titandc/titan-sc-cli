package api

import (
    "encoding/json"
    "fmt"
    "github.com/spf13/cobra"
    "log"
    "os"
    "text/tabwriter"
)

func (API *APITitan) IPAttach(cmd *cobra.Command, args []string) {
    API.ParseGlobalFlags(cmd)
    serverUUID := args[0]
    IP, _ := cmd.Flags().GetString("ip")
    IPVersion := 4

    if IP == "" {
        log.Println("Attach server IP: missing --ip argument.")
        return
    }
    API.AttachDeatchIPServer(HTTPPost, serverUUID, IP, IPVersion)
}

func (API *APITitan) IPDetach(cmd *cobra.Command, args []string) {
    API.ParseGlobalFlags(cmd)
    serverUUID := args[0]
    IP, _ := cmd.Flags().GetString("ip")
    IPVersion := 4

    if IP == "" {
        log.Println("Deatch server IP: missing --ip argument.")
        return
    }
    API.AttachDeatchIPServer(HTTPDelete, serverUUID, IP, IPVersion)
}

func (API *APITitan) AttachDeatchIPServer(HttpMethod, serverUUID, ip string, version int) {
    ipOpt := APIIP{
        IP:      ip,
        Version: version,
    }
    reqData, err := json.Marshal(ipOpt)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    err = API.SendAndResponse(HttpMethod, "/compute/servers/"+serverUUID+"/ips", reqData)
    if err != nil {
        fmt.Println(err.Error())
    }
    if !API.HumanReadable {
        API.PrintJson()
    } else {
        API.DefaultPrintReturn()
    }
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
