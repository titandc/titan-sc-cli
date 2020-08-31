package api

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
    Color         bool
    CLIVersion    string
    CLIos         string
    Token         string
    URI           string
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
    Notes        string `json:"notes"`
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
    ISOs           []APIISO `json:"isos"`
    PendingActions []string `json:"pending_actions"`
    ManagedNetwork string   `json:"managed_network"`
}

type APIServerUpdateInfos struct {
    Name    string `json:"name"`
    Notes   string `json:"notes"`
    Reverse string `json:"reverse"`
}

type APIISO struct {
    Protocol string `json:"protocol"`
    ISOPath  string `json:"iso_path"`
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
    Speed   APISize `json:"speed"`
    State   string  `json:"state"`
    UUID    string  `json:"uuid"`
    Managed bool    `json:"managed"`
    CIDR    string  `json:"cidr"`
    Gateway string  `json:"gateway"`
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
    CIDR   string  `json:"cidr,omitempty"`
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

type APIServerLOadISORequest struct {
    UUID     string `json:"uuid"`
    ISO      string `json:"iso"`
    Protocol string `json:"protocol"`
}

type APIReturn struct {
    Success string `json:"success,omitempty"`
    Error   string `json:"error,omitempty"`
}

type APIIP struct {
    IP      string `json:"ip"`
    Version int    `json:"type"`
}

type APIVersion struct {
    Date    string `json:"release_date"`
    Version string `json:"version"`
}

type ManagedServices struct {
    Company string `json:"company"`
}
