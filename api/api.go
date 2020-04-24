package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/spf13/cobra"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "text/tabwriter"
    "time"
)

const (
    HTTPHost   = "https://sc.titandc.net/api/v1"
    HTTPGet    = http.MethodGet
    HTTPPut    = http.MethodPut
    HTTPPost   = http.MethodPost
    HTTPDelete = http.MethodDelete
    // Get N history entry
    NHistory = 25
)

/*
 *
 *
 ******************
 * API structure
 ******************
 *
 *
 */

type APITitan struct {
    HumanReadable bool
    CLIVersion    string
    CLIos         string
    Token         string
    RespBody      []byte
}

type APISize struct {
    Unit  string `json:"unit"`
    Value int    `json:"value"`
}

type APIServer struct {
    UUID         string `json:"uuid"`
    Login        string `json:"user_login"`
    OS           string `json:"os"`
    Name         string `json:"Name"`
    Plan         string `json:"Plan"`
    State        string `json:"State"`
    Template     string `json:"template"`
    Hypervisor   string `json:"hypervisor"`
    Reverse      string `json:"reverse"`
    Gateway      string `json:"gateway"`
    Mac          string `json:"mac"`
    IP           string `json:"ip"`
    IPv6         string `json:"ipv6"`
    CompanyName  string `json:"company"`
    Creationdate int    `json:"creation_date"`
    Bandwidth    struct {
        Uint   string `json:"unit"`
        Input  int    `json:"input"`
        Output int    `json:"output"`
    } `json:"bandwidth"`
    KvmIp APIKvmIP `json:"kvm_ip"`
    CPU   struct {
        NbCores int `json:"nb_cores"`
    } `json:"cpu"`
    RAM  APISize `json:"ram"`
    Disk struct {
        Size APISize `json:"size"`
        QoS  struct {
            Unit  string `json:"unit"`
            Read  int    `json:"read"`
            Write int    `json:"write"`
        } `json:"qos"`
        IOPS struct {
            Unit      string `json:"unit"`
            Read      int64  `json:"read"`
            Write     int64  `json:"write"`
            BlockSize string `json:"block_size"`
        } `json:"iops"`
    } `json:"disk"`
}

type APIOwner struct {
    Email      string `json:"email"`
    Firstname  string `json:"firstname"`
    Lastname   string `json:"lastname"`
    LastLogin  int    `json:"last_login"`
    Salutation string `json:"salutation"`
    UUID       string `json:"uuid"`
}

type APINetwork struct {
    Company   string   `json:"company"`
    CreatedAt int      `json:"created_at"`
    MaxMtu    int      `json:"max_mtu"`
    Name      string   `json:"name"`
    Owner     APIOwner `json:"owner"`
    Ports     int      `json:"ports"`
    Servers   []struct {
        Name  string `json:"name"`
        OS    string `json:"os"`
        Plan  string `json:"plan"`
        State string `json:"state"`
        UUID  string `json:"uuid"`
    } `json:"servers"`
    Speed APISize `json:"speed"`
    State string  `json:"state"`
    UUID  string  `json:"uuid"`
}

type APINetworkList struct {
    Quota    int64        `json:"NetworksQuota"`
    NetInfos []APINetwork `json:"NetworkFullInfos"`
}

type APICompany struct {
    Company struct {
        Description string `json:"description"`
        Name        string `json:"name"`
        UUID        string `json:"uuid"`
    } `json:"company"`
    Role struct {
        Name string `json:"name"`
    } `json:"role"`
}

type APIServerAction struct {
    Action string `json:"action"`
}

type APISnapshot struct {
    CreatedAt string  `json:"created_at"`
    Name      string  `json:"name"`
    Size      APISize `json:"size"`
    State     string  `json:"state"`
    UUID      string  `json:"uuid"`
}

type APIHistoryEvent struct {
    Server    APIServer `json:"server"`
    State     string    `json:"state"`
    Status    string    `json:"status"`
    Timestamp int       `json:"timestamp"`
    Type      string    `json:"type"`
}

type APIServerResourceHistory struct {
    UUID             string  `json:"uuid"`
    Timestamp        int     `json:"timestamp"`
    CurrentSize      int     `json:"current_size,omitempty"`
    OSPercentageUsed float64 `json:"os_percentage_used,omitempty"`
    OSMaxSize        int     `json:"os_max_size,omitempty"`
    CPUTime          int     `json:"cpu_time,omitempty"`
    RamTotal         int     `json:"ram_total,omitempty"`
    RamUsed          int     `json:"ram_used,omitempty"`
    RxBytes          int     `json:"rx_bytes,omitempty"`
    TxBytes          int     `json:"tx_bytes,omitempty"`
    CpuPercent       int     `json:"cpu_percent,omitempty"`
    RamPercent       int     `json:"ram_percent,omitempty"`
}

type APIWeatherMap struct {
    Compute        string `json:"compute"`
    PrivateNetwork string `json:"private_network"`
    PublicNetwork  string `json:"public_network"`
    Storage        string `json:"storage"`
}

type APINetworkCreate struct {
    MaxMTU int     `json:"max_mtu"`
    Name   string  `json:"name"`
    Ports  int     `json:"ports"`
    Speed  APISize `json:"value"`
}

type APIKvmIP struct {
    Status string `json:"status"`
    URI    string `json:"uri"`
}

type APICompanyDetail struct {
    Addresses []struct {
        Name       string `json:"name"`
        City       string `json:"city"`
        Country    string `json:"country"`
        PostalCode string `json:"postal_code"`
        Street     string `json:"street"`
        Street2    string `json:"street2"`
        Type       string `json:"type"`
    } `json:"addresses"`
    CA          int    `json:"ca"`
    Description string `json:"description"`
    Email       string `json:"email"`
    Members     []struct {
        Email     string `json:"email"`
        Firstname string `json:"firstname"`
        Lastname  string `json:"lastname"`
        Phone     string `json:"phone"`
        LastLogin string `json:"last_login"`
        Roles     []struct {
            Accountant    bool   `json:"accountant"`
            Administrator bool   `json:"administrator"`
            Company       string `json:"company"`
            Manager       bool   `json:"manager"`
            Name          string `json:"name"`
            Position      string `json:"position"`
            UUID          string `json:"uuid"`
        } `json:"roles"`
        Salutation string `json:"salutation"`
        UUID       string `json:"uuid"`
    } `json:"members"`
    Naf    string   `json:"naf"`
    Name   string   `json:"name"`
    Note   string   `json:"note"`
    Owner  APIOwner `json:"owner"`
    Phone  string   `json:"phone"`
    Quotas struct {
        CPUs     int `json:"cpus"`
        Networks int `json:"networks"`
        Servers  int `json:"servers"`
    } `json:"quotas"`
    Siret     string `json:"siret"`
    TvaNumber string `json:"tva_number"`
    TvaRate   int    `json:"tva_rate"`
    UUID      string `json:"uuid"`
    Website   string `json:"website"`
}

/*
 *
 *
 ******************
 * Global variable
 ******************
 *
 *
 */
var API = &APITitan{
    Token:         "",
    HumanReadable: false,
}

/*
 *
 *
 ******************
 * Company function
 ******************
 *
 *
 */
func (API *APITitan) CompaniesList(cmd *cobra.Command, args []string) {
    _ = args
    API.ParseGlobalFlags(cmd)

    if compagnies, err := API.GetCompagnies()
        err != nil {
        fmt.Println(err.Error())
    } else {
        if !API.HumanReadable {
            API.PrintJson()
        } else {
            w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
            _, _ = fmt.Fprintf(w, "COMPANY UUID\tNAME\tROLE\t\n")
            for _, company := range compagnies {
                _, _ = fmt.Fprintf(w, "%s\t%s\t%s\t\n",
                    company.Company.UUID, company.Company.Name, company.Role.Name)
                _ = w.Flush()
            }
        }
    }
    return
}

func (API *APITitan) CompanyDetail(cmd *cobra.Command, args []string) {
    companyUUID := args[0]
    API.ParseGlobalFlags(cmd)

    err := API.SendAndResponse(HTTPGet, "/companies/"+companyUUID, nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    if !API.HumanReadable {
        API.PrintJson()
    } else {
        company := &APICompanyDetail{}
        if err := json.Unmarshal(API.RespBody, &company)
            err != nil {
            fmt.Println(err.Error())
            return
        }
        fmt.Printf("Company %s informations:\n", company.Name)
        fmt.Printf("  UUID: %s\n"+
            "  Phone: %s\n"+
            "  Description: %s\n"+
            "  Email: %s\n"+
            "  CA: %d\n"+
            "  NAF: %s\n"+
            "  Siret: %s\n"+
            "  TVA number: %s\n"+
            "  TVA rate: %d\n"+
            "  Website: %s\n"+
            "  Note: %s\n"+
            "  Qutas:\n"+
            "    CPUs: %d\n"+
            "    Networks: %d\n"+
            "    Servers: %d\n",
            company.UUID, company.Phone,
            company.Description, company.Email, company.CA,
            company.Naf, company.Siret, company.TvaNumber,
            company.TvaRate, company.Note, company.Website,
            company.Quotas.CPUs, company.Quotas.Networks,
            company.Quotas.Servers)
        fmt.Println("  Address:")
        for _, addr := range company.Addresses {
            fmt.Printf("    - Name: %s\n"+
                "      City: %s\n"+
                "      Country: %s\n"+
                "      Postal code: %s\n"+
                "      Type: %s\n"+
                "      Street: %s\n",
                addr.Name, addr.City, addr.Country,
                addr.PostalCode, addr.Type, addr.Street)
            if addr.Street2 != "" {
                fmt.Printf("      Street2: %s\n", addr.Street2)
            }
        }
        fmt.Println("  Members:")
        for _, member := range company.Members {
            fmt.Printf("    - Name: %s %s (%s)\n"+
                "      Phone: %s\n"+
                "      UUID: %s\n",
                member.Firstname, member.Lastname,
                member.Email, member.Phone, member.UUID)
        }
    }
}

func (API *APITitan) GetCompagnies() ([]APICompany, error) {
    if err := API.SendAndResponse(HTTPGet, "/companies", nil); err != nil {
        return nil, err
    }

    compagnies := make([]APICompany, 0)
    if err := json.Unmarshal(API.RespBody, &compagnies); err != nil {
        return nil, err
    }
    return compagnies, nil
}

func (API *APITitan) GetCompanyServers(companyUUID string) ([]APIServer, error) {
    err := API.SendAndResponse(HTTPGet, "/compute/servers?company_uuid="+companyUUID, nil)
    if err != nil {
        return nil, err
    }

    servers := make([]APIServer, 0)
    if err := json.Unmarshal(API.RespBody, &servers); err != nil {
        return nil, err
    }
    return servers, nil
}

/*
 *
 *
 ******************
 * Servers function
 ******************
 *
 *
 */
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
                        _, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n",
                            server.UUID, server.Plan, server.State, server.OS, server.Name)
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
    fmt.Println("Sending server action success.")
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
    fmt.Printf("Name: %s\n"+
        "UUID: %s\n"+
        "VM Login: %s\n"+
        "State: %s\n"+
        "OS and version: %s\n"+
        "company: %s\n"+
        "Hypervisor: %s\n"+
        "Reverse: %s\n"+
        "IPv4: %s\n"+
        "IPv6: %s\n"+
        "Mac: %s\n"+
        "Gateway: %s\n"+
        "Bandwidth in/out: %d/%d %s\n"+
        "IP Kvm: %s\n"+
        "Cpu(s): %d\n"+
        "RAM: %d %s\n"+
        "Disk: %d %s\n"+
        "Disk QoS Read/Write: %d/%d %s\n"+
        "Disk IOPS Read/Write/BlockSize: %d/%d/%s %s\n",
        server.Name, server.UUID, server.Login, server.State, server.Template,
        server.CompanyName, server.Hypervisor, server.Reverse,
        server.IP, server.IPv6, server.Mac, server.Gateway,
        server.Bandwidth.Input, server.Bandwidth.Output,
        server.Bandwidth.Uint, server.KvmIp.Status,
        server.CPU.NbCores, server.RAM.Value, server.RAM.Unit,
        server.Disk.Size.Value, server.Disk.Size.Unit,
        server.Disk.QoS.Read, server.Disk.QoS.Write, server.Disk.QoS.Unit,
        server.Disk.IOPS.Read, server.Disk.IOPS.Write, server.Disk.IOPS.BlockSize,
        server.Disk.IOPS.Unit)
}

/*
 *
 *
 **************************
 * Snapshot server function
 **************************
 *
 *
 */
func (API *APITitan) SnapshotList(cmd *cobra.Command, args []string) {

    serverUUID := args[0]
    API.ParseGlobalFlags(cmd)

    snapshots, err := API.SnapshotServerUUIDList(serverUUID)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    if len(snapshots) == 0 {
        fmt.Println("0 Snapshot")
        return
    }

    for _, snap := range snapshots {
        if !API.HumanReadable {
            API.PrintJson()
        } else {
            w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
            _, _ = fmt.Fprintf(w, "SNAPSHOT UUID\tTIMESTAMP\tSIZE\tNAME\t\n")
            API.PrintSnapshotInfos(w, &snap)
        }
    }
}

func (API *APITitan) SnapshotCreate(cmd *cobra.Command, args []string) {

    serverUUID := args[0]
    API.ParseGlobalFlags(cmd)

    err := API.SendAndResponse(HTTPPost, "/compute/servers/"+serverUUID+"/snapshots", nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    if !API.HumanReadable {
        API.PrintJson()
    } else {
        snap := &APISnapshot{}
        if err := json.Unmarshal(API.RespBody, &snap); err != nil {
            log.Println("Human readable format error: ", err.Error())
            return
        }
        w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
        _, _ = fmt.Fprintf(w, "SNAPSHOT UUID\tTIMESTAMP\tSIZE\t\tNAME\t\n")
        API.PrintSnapshotInfos(w, snap)
    }
}

func (API *APITitan) SnapshotRemove(cmd *cobra.Command, args []string) {

    _ = args
    API.ParseGlobalFlags(cmd)

    serverUUID, _ := cmd.Flags().GetString("server-uuid")
    snapUUID, _ := cmd.Flags().GetString("snapshot-uuid")

    if snapUUID == "" {
        fmt.Println("error: Snapshot UUID missing.")
        return
    }
    if serverUUID == "" {
        fmt.Println("error: Server UUID missing.")
        return
    }

    snapshots, err := API.SnapshotServerUUIDList(serverUUID)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    for _, snap := range snapshots {
        if snapUUID != snap.UUID {
            continue
        }

        err = API.SendAndResponse(HTTPDelete, "/compute/servers/"+
            serverUUID+"/snapshots/"+snapUUID, nil)
        if err != nil {
            fmt.Println(err.Error())
        }

        fmt.Println("Snapshot deleting request successfully sent.")
        return
    }
    fmt.Println("Snapshot UUID", snapUUID, "not found")
}

func (API *APITitan) SnapshotServerUUIDList(serverUUID string) ([]APISnapshot, error) {

    err := API.SendAndResponse(HTTPGet, "/compute/servers/"+serverUUID+"/snapshots", nil)
    if err != nil {
        return nil, err
    }

    snap := make([]APISnapshot, 0)
    if err := json.Unmarshal(API.RespBody, &snap); err != nil {
        return nil, err
    }
    return snap, nil
}

func (API *APITitan) PrintSnapshotInfos(w *tabwriter.Writer, snap *APISnapshot) {

    _, _ = fmt.Fprintf(w, "%s\t%s\t%d %s\t%s\t\n",
        snap.UUID, snap.CreatedAt, snap.Size.Value,
        snap.Size.Unit, snap.Name)
    _ = w.Flush()
}

/*
 *
 *
 ******************
 * history function
 ******************
 *
 *
 */

func (API *APITitan) HistoryCompanyEvent(cmd *cobra.Command, args []string) {

    _ = args
    API.ParseGlobalFlags(cmd)

    serverUUID, _ := cmd.Flags().GetString("server-uuid")
    companyUUID, _ := cmd.Flags().GetString("company-uuid")
    number, _ := cmd.Flags().GetInt("number")

    if serverUUID != "" && companyUUID != "" ||
        (companyUUID == "" && serverUUID == "") {
        _ = cmd.Help()
        return
    }

    if number < 1 {
        number = NHistory
    }
    strNumber := fmt.Sprintf("%d", number)

    if companyUUID != "" {
        API.HistoryByCompany(strNumber, companyUUID)
    } else {
        API.HistoryByServer(strNumber, serverUUID)
    }
}

func (API *APITitan) HistoryByCompany(number, companyUUID string) {

    err := API.SendAndResponse(HTTPGet, "/compute/servers/events?nb="+
        number+"&company_uuid="+companyUUID, nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    if !API.HumanReadable {
        API.PrintJson()
    } else {
        history := make([]APIHistoryEvent, 0)
        if err := json.Unmarshal(API.RespBody, &history); err != nil {
            fmt.Println(err.Error())
            return
        }
        for _, event := range history {
            date := API.DateFormat(event.Timestamp)
            fmt.Printf("Server UUID:\r\t\t%s\n"+
                "Server Name:\r\t\t%s\n"+
                "Event type:\r\t\t%s (%s)\n"+
                "Event status:\r\t\t%s\n"+
                "Date:\r\t\t%s\n\n",
                event.Server.UUID, event.Server.Name, event.Type,
                event.State, event.Status, date)
        }
    }
}

func (API *APITitan) HistoryByServer(number, serverUUID string) {

    err := API.SendAndResponse(HTTPGet, "/compute/servers/"+serverUUID+
        "/events?nb="+number, nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    if !API.HumanReadable {
        API.PrintJson()
    } else {
        history := make([]APIHistoryEvent, 0)
        if err := json.Unmarshal(API.RespBody, &history); err != nil {
            fmt.Println(err.Error())
            return
        }
        for _, event := range history {
            date := API.DateFormat(event.Timestamp)
            fmt.Printf("Event type:\r\t\t%s (%s)\n"+
                "Event status:\r\t\t%s\n"+
                "Date:\r\t\t%s\n\n",
                event.Type, event.State, event.Status, date)
        }
    }
}

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
    companyUUID, _ := cmd.Flags().GetString("server-uuid")
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

func (API *APITitan) NetworkAttachServer(cmd *cobra.Command, args []string) {
    _ = args
    networkUUID, _ := cmd.Flags().GetString("network-uuid")
    serverUUID, _ := cmd.Flags().GetString("server-uuid")
    API.NetworkServerOps(cmd, networkUUID, serverUUID, "attach")
}

func (API *APITitan) NetworkDetachServer(cmd *cobra.Command, args []string) {
    _ = args
    networkUUID, _ := cmd.Flags().GetString("network-uuid")
    serverUUID, _ := cmd.Flags().GetString("server-uuid")
    API.NetworkServerOps(cmd, networkUUID, serverUUID, "detach")
}

type APINetworkOps struct {
    ServerUUID string `json:"server_uuid"`
}

func (API *APITitan) NetworkServerOps(cmd *cobra.Command, networkUUID, serverUUID, ops string) {

    if networkUUID == "" {
        fmt.Printf("Error: --network-uuid missing.\n\n")
        fmt.Println(cmd.Help())
        return
    }
    if serverUUID == "" {
        fmt.Printf("Error: --server-uuid missing.\n\n")
        fmt.Println(cmd.Help())
        return
    }

    // send request
    act := APINetworkOps{
        ServerUUID: serverUUID,
    }
    reqData, err := json.Marshal(act)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    err = API.SendAndResponse(HTTPPut, "/compute/networks/"+networkUUID+"/"+ops, reqData)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println("Sending request server", ops, " to network success.")
}

func (API *APITitan) NetworkCreate(cmd *cobra.Command, args []string) {

    companyUUID := args[0]
    API.ParseGlobalFlags(cmd)
    networkName, _ := cmd.Flags().GetString("name")

    net := APINetworkCreate{
        MaxMTU: 8948,
        Name:   networkName,
        Ports:  6,
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

    err := API.SendAndResponse(HTTPDelete, "/compute/networks/"+networkUUID, nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println("Remove network sending success.")
    return
}

func (API *APITitan) DateFormat(timestamp int) string {
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
        "  Max MTU: %d\n"+
        "  Owner informations:\n"+
        "    Name: %s %s (%s)\n"+
        "    UUID: %s\n",
        net.Name, date, net.Ports, net.Speed.Value, net.Speed.Unit,
        net.State, net.UUID, net.Company, net.MaxMtu,
        net.Owner.Firstname, net.Owner.Lastname, net.Owner.Email,
        net.Owner.UUID)
}

type APINetworkRename struct {
    Name string `json:"name"`
}

func (API *APITitan) NetworkRename(cmd *cobra.Command, args []string) {

    networkUUID, _ := cmd.Flags().GetString("network-uuid")
    name, _ := cmd.Flags().GetString("name")

    if networkUUID == "" {
        fmt.Println("Error: --network-uuid missing.")
        return
    }
    if name == "" {
        fmt.Println("Error: --name missing.")
        return
    }

    netRename := &APINetworkRename{Name: name}
    reqData, err := json.Marshal(netRename)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    err = API.SendAndResponse(HTTPPut, "/compute/networks/"+networkUUID, reqData)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Println("Sending rename network name success.")
}

/*
 *
 *
 *****************
 * IP Kvm function
 *****************
 *
 *
 */
func (API *APITitan) IPKvmGetInfos(cmd *cobra.Command, args []string) {

    serverUUID := args[0]
    API.ParseGlobalFlags(cmd)

    err := API.SendAndResponse(HTTPGet, "/compute/servers/"+serverUUID+"/ipkvm", nil)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    if !API.HumanReadable {
        API.PrintJson()
    } else {
        API.IPKvmPrint(serverUUID)
    }
}

func (API *APITitan) IPKvmStart(cmd *cobra.Command, args []string) {
    API.ParseGlobalFlags(cmd)
    API.IPKvmAction("start", args[0])
}

func (API *APITitan) IPKvmStop(cmd *cobra.Command, args []string) {
    API.ParseGlobalFlags(cmd)
    API.IPKvmAction("stop", args[0])
}

func (API *APITitan) IPKvmAction(action, serverUUID string) {
    act := APIServerAction{
        Action: action,
    }
    reqData, e := json.Marshal(act)
    if e != nil {
        fmt.Println(e.Error())
        return
    }
    err := API.SendAndResponse(HTTPPut, "/compute/servers/"+serverUUID+"/ipkvm", reqData)
    if err != nil {
        fmt.Println(err.Error())
    }

    if !API.HumanReadable {
        API.PrintJson()
    } else {
        fmt.Println("Send IP Kvm request succes.")
    }
}

func (API *APITitan) IPKvmPrint(serverUUID string) {
    ipkvm := &APIKvmIP{}
    if err := json.Unmarshal(API.RespBody, ipkvm); err != nil {
        fmt.Println(err.Error())
        return
    }
    fmt.Printf("IP KVM informations:\n  Server UUID: %s\n  Status: %s\n", serverUUID, ipkvm.Status)
    if ipkvm.URI != "" {
        fmt.Println("  URI:", ipkvm.URI)
    }
}

/*
 *
 *
 *********************
 * Other API function
 *********************
 *
 *
 */
func (API *APITitan) WeatherMap(cmd *cobra.Command, args []string) {

    _ = args
    API.ParseGlobalFlags(cmd)

    if err := API.SendAndResponse(HTTPGet, "/weather", nil); err != nil {
        fmt.Println(err.Error())
        return
    }

    if !API.HumanReadable {
        API.PrintJson()
    } else {
        weatherMap := &APIWeatherMap{}
        if err := json.Unmarshal(API.RespBody, weatherMap); err != nil {
            log.Println(err.Error())
            return
        }
        fmt.Printf("Titan Weather Map:\n"+
            "  Compute: %s\n"+
            "  Storage: %s\n"+
            "  Public network: %s\n"+
            "  Private network: %s\n",
            weatherMap.Compute, weatherMap.Storage,
            weatherMap.PublicNetwork, weatherMap.PrivateNetwork)
    }
}

/*
 *
 *
 ******************
 * Utils function
 ******************
 *
 *
 */

func (API *APITitan) PrintJson() {
    dst := &bytes.Buffer{}
    if err := json.Indent(dst, API.RespBody, "", "  "); err != nil {
        fmt.Printf(string(API.RespBody))
        return
    }
    fmt.Printf(dst.String())
}

func (API *APITitan) ParseGlobalFlags(cmd *cobra.Command) {
    var err error

    API.HumanReadable, err = cmd.Flags().GetBool("human")
    if err != nil {
        API.HumanReadable = false
    }
}

func (API *APITitan) SendAndResponse(method, path string, req []byte) error {
    request, err := http.NewRequest(method, HTTPHost+path, bytes.NewBuffer(req))
    if err != nil {
        return err
    }
    request.Header.Add("X-API-KEY", API.Token)
    request.Header.Add("Titan-Cli-Os", API.CLIos)
    request.Header.Add("Titan-Cli-Version", API.CLIVersion)
    request.Header.Set("Content-Type", "application/json; charset=utf-8")

    client := &http.Client{}
    resp, err := client.Do(request)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    API.RespBody, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    return nil
}
