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
/*
 *
 *
 ******************
 * API structure
 ******************
 *
 *
 */

type APISize struct {
	Unit  string `json:"unit"`
	Value int32  `json:"value"`
}

type APIServer struct {
	UUID       string `json:"uuid"`
	Login      string `json:"user_login"`
	Name       string `json:"Name"`
	Plan       string `json:"plan"`
	State      string `json:"state"`
	Template   string `json:"template"`
	Hypervisor struct {
		UUID     string `json:"uuid"`
		Hostname string `json:"hostname"`
		State    string `json:"state"`
	} `json:"hypervisor"`
	Reverse string `json:"reverse"`
	Gateway string `json:"gateway"`
	Mac     string `json:"mac"`
	IPs     []struct {
		IP   string `json:"ip"`
		Type int64  `json:"type"`
	} `json:"ips"`
	IPv6         string `json:"ipv6"`
	CompanyName  string `json:"company"`
	CreationDate int64  `json:"creation_date"`
	Notes        string `json:"notes"`
	Bandwidth    struct {
		Unit   string `json:"unit"`
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
	Managed        bool     `json:"managed"`
	Notifications  []string `json:"notifications"`
}

type APIServerUpdateInfos struct {
	Name    string `json:"name"`
	Notes   string `json:"notes"`
	Reverse string `json:"reverse"`
}

type APIDeleteServer struct {
	Reason string `json:"reason"`
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
	CreatedAt int64    `json:"created_at"`
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
	Speed    APISize `json:"speed"`
	State    string  `json:"state"`
	UUID     string  `json:"uuid"`
	Managed  bool    `json:"managed"`
	CIDR     string  `json:"cidr"`
	Gateway  string  `json:"gateway"`
	Firewall struct {
		Policy string `json:"policy"`
		Rules  []struct {
			Server   string `json:"server"`
			Protocol string `json:"protocol"`
			Port     string `json:"port"`
			Source   string `json:"source"`
		} `json:"rules"`
	} `json:"firewall"`
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
		Description string `json:"description"`
		Name        string `json:"name"`
		UUID        string `json:"uuid"`
		Role        struct {
			Name string `json:"name"`
		} `json:"role"`
	}
}

type APIServerAction struct {
	Action string `json:"action"`
}

type APISnapshot struct {
	UUID      string  `json:"uuid"`
	CreatedAt int64   `json:"created_at"`
	Name      string  `json:"name"`
	Size      APISize `json:"size"`
	State     string  `json:"state"`
}

type APIHistoryEvent struct {
	Server    APIServer `json:"server"`
	State     string    `json:"state"`
	Status    string    `json:"status"`
	Timestamp int64     `json:"timestamp"`
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
	Name   string  `json:"name"`
	MaxMTU int32   `json:"max_mtu"`
	Ports  int32   `json:"ports"`
	Speed  APISize `json:"speed"`
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
	Managed   bool   `json:"managed_services"`
}

type APIServerLOadISORequest struct {
	ISO      string `json:"iso"`
	Protocol string `json:"protocol"`
}

type APIReturn struct {
	Code    string `json:"code" binding:"required"`
	Success string `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

type APIIPAttachDetach struct {
	IP      string `json:"ip"`
	Version int    `json:"type"`
}

type APIPNATRuleInfos struct {
	IP          string `json:"ip" binding:"required"`
	Transparent bool   `json:"transparent,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	PortSrc     int64  `json:"port_src,omitempty"`
	PortDst     int64  `json:"port_dst,omitempty"`
}

type APIPNATRuleAddDel struct {
	IP          string `json:"ip" binding:"required"`
	Transparent bool   `json:"transparent"`
	Protocol    string `json:"protocol"`
	PortSrc     int64  `json:"port_src"`
	PortDst     int64  `json:"port_dst"`
}

type APIVersion struct {
	Date    string `json:"release_date"`
	Version string `json:"version"`
}

type APIManagedServices struct {
	Company string `json:"company"`
}

type APINetworkFirewallRule struct {
	ServerUUID string `json:"server_uuid" binding:"required"`
	Protocol   string `json:"protocol" binding:"required"`
	Port       string `json:"port" binding:"required"`
	Source     string `json:"source" binding:"required"`
}

type APINetworkFullInfosFirewall struct {
	Policy string                            `form:"policy,omitempty" json:"policy,omitempty"`
	Rules  []APINetworkFullInfosFirewallRule `form:"rules,omitempty" json:"rules,omitempty"`
}

type APINetworkFullInfosFirewallRule struct {
	Server   string `form:"server,omitempty" json:"server,omitempty"`
	Protocol string `form:"protocol,omitempty" json:"protocol,omitempty"`
	Port     string `form:"port,omitempty" json:"port,omitempty"`
	Source   string `form:"source,omitempty" json:"source,omitempty"`
}

type APIUserInfos struct {
	UUID              string          `json:"uuid" binding:"required"`
	Firstname         string          `json:"firstname"`
	Lastname          string          `json:"lastname"`
	Email             string          `json:"email"`
	Phone             string          `json:"phone"`
	CreatedAt         int64           `json:"created_at"`
	LastLogin         int64           `json:"last_login"`
	Mobile            string          `json:"mobile"`
	SSHKeys           []APIUserSSHKey `json:"ssh_keys"`
	PreferredLanguage string          `json:"preferred_language"`
	Salutation        string          `json:"salutation"`
	LatestCGVSigned   bool            `json:"latest_cgv_signed"`
	CGVDate           int64           `json:"cgv_signature_date"`
	CGVLink           string          `json:"cgv_link"`
	CGVVersion        string          `json:"cgv_version"`
	Signature         string          `json:"signature"`
	TwoFA             bool            `json:"twofa"`
}

type APIUserSSHKeys struct {
	SSHKeys []APIUserSSHKey `json:"ssh_keys"`
}

type APIUserSSHKey struct {
	Title   string `bson:"name,omitempty" json:"name,omitempty"`
	Content string `bson:"value,omitempty" json:"value,omitempty"`
}

type APIAddUserSSHKey struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type APIDeleteUserSSHKey struct {
	Name string `json:"name"`
}

type CreateServersDetail struct {
	Username        string                  `json:"username"`
	Quantity        int64                   `json:"quantity"`
	UserPassword    string                  `json:"user_password,omitempty"`
	UserLogin       string                  `json:"user_login,omitempty"`
	UserSSHKeys     []string                `json:"user_ssh_keys,omitempty"`
	TemplateOS      string                  `json:"template_os"`
	TemplateVersion string                  `json:"template_version"`
	Plan            string                  `json:"plan"`
	Addons          []APIInstallAddonsAddon `json:"addons,omitempty"`
	ManagedNetwork  string                  `json:"managed_network,omitempty"`
}

type APICreateServers struct {
	CreateServersDetail []CreateServersDetail `json:"servers"`
}

type APIInstallAddonsAddon struct {
	Item     string `json:"item"`
	Quantity int64  `json:"quantity"`
}

type APIAddonsItem struct {
	UUID       string               `bson:"uuid" json:"uuid"`
	Name       string               `bson:"name" json:"name"`
	Amount     APIAddonsItemAmount  `bson:"amount" json:"amount"`
	ZohoID     string               `bson:"zoho_id" json:"zoho_id"`
	PricingSC1 APIAddonsItemPricing `bson:"pricing_SC1" json:"pricing_SC1"`
	PricingSC2 APIAddonsItemPricing `bson:"pricing_SC2" json:"pricing_SC2"`
	PricingSC3 APIAddonsItemPricing `bson:"pricing_SC3" json:"pricing_SC3"`
}

type APIAddonsItemAmount struct {
	Unit  string `bson:"unit" json:"unit"`
	Value int64  `bson:"value" json:"value"`
}

type APIAddonsItemPricing struct {
	Value    float64 `bson:"value" json:"value"`
	Currency string  `bson:"currency" json:"currency"`
}

type APITemplateFullInfos struct {
	OS       string                        `json:"os"`
	Versions []APITemplateFullInfosVersion `json:"version"`
}

type APITemplateFullInfosVersion struct {
	UUID       string `json:"uuid"`
	Version    string `json:"version"`
	LastUpdate int64  `json:"last_update"`
}

type APIResetServer struct {
	UserPassword    string   `json:"user_password,omitempty"`
	UserSSHKeys     []string `json:"user_ssh_keys,omitempty"`
	TemplateOS      string   `json:"template_os"`
	TemplateVersion string   `json:"template_version"`
}

type IsAdminStruct struct {
	Admin bool `json:"admin"`
}
