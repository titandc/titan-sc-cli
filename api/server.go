package api

import (
    "encoding/json"
    "fmt"
    "github.com/spf13/cobra"
    "log"
    "os"
    "text/tabwriter"
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

    if serverUUID == "" {
        log.Println("Server change name error: missing --server-uuid argument.")
        return
    }
    if newName == "" {
        log.Println("Server change name error: missing --name argument.")
        return
    }

    updateInfos := &APIServerUpdateInfos{
        Name: newName,
    }
    API.ServerUpdateInfos(serverUUID, updateInfos)
}

func (API *APITitan) ServerChangeReverse(cmd *cobra.Command, args []string) {
    _ = args
    API.ParseGlobalFlags(cmd)
    serverUUID, _ := cmd.Flags().GetString("server-uuid")
    newReverse, _ := cmd.Flags().GetString("reverse")

    if serverUUID == "" {
        log.Println("Server change reverse error: missing --server-uuid argument.")
        return
    }
    if newReverse == "" {
        log.Println("Server change reverse error: missing --reverse argument.")
        return
    }

    updateInfos := &APIServerUpdateInfos{
        Reverse: newReverse,
    }
    API.ServerUpdateInfos(serverUUID, updateInfos)
}

func (API *APITitan) ServerUpdateInfos(serverUUID string, updateInfos *APIServerUpdateInfos) {
    reqData, err := json.Marshal(updateInfos)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    err = API.SendAndResponse(HTTPPut, "/compute/servers/"+serverUUID, reqData)
    if err != nil {
        fmt.Println(err.Error())
    }
    if !API.HumanReadable {
        API.PrintJson()
    } else {
        API.DefaultPrintReturn()
    }
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
                    _, _ = fmt.Fprintf(w, "UUID\tPLAN\tSTATE\tOS\tNAME\t\n")
                    for _, server := range servers {
                        state := API.ServerStateSetColor(server.State)
                        _, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n",
                            server.UUID, server.Plan, state, server.OS, server.Name)
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
    reqData, e := json.Marshal(act)
    if e != nil {
        fmt.Println(e.Error())
        return
    }
    err = API.SendAndResponse(HTTPPut, "/compute/servers/"+server.UUID+"/action", reqData)
    if err != nil {
        fmt.Println(err.Error())
    }
    if !API.HumanReadable {
        API.PrintJson()
    } else {
        API.DefaultPrintReturn()
    }
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
        colorState = "\033[1;31m"+state+"\033[0m"
    case "started":
        colorState = "\033[1;32m"+state+"\033[0m"
    case "stopped":
        colorState = "\033[1;33m"+state+"\033[0m"
    }
    return colorState
}