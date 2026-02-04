package api

/*
 *
 *
 ******************
 * API structures
 ******************
 *
 *
 */

// APISize represents a size with unit (e.g., RAM, disk, bandwidth)
type APISize struct {
	Unit  string `json:"unit"`
	Value int32  `json:"value"`
}

// APIServer represents a server resource from the API
type APIServer struct {
	UUID       string                 `json:"uuid"`
	Name       string                 `json:"Name"`
	Plan       string                 `json:"plan"`
	State      string                 `json:"state"`
	Template   APIServerTemplateInfos `json:"template_full"`
	Image      *APIServerImageInfos   `json:"image_full"`
	Disksource APIServerDiskSource    `json:"disk_source"`
	SSHKeys    []string               `json:"ssh_keys"`
	Login      string                 `json:"user_login"`
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
	IPv6    string `json:"ipv6"`
	Company struct {
		UUID        string `json:"uuid"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Disabled    bool   `json:"disabled"`
	} `json:"company"`
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

// APIServerDiskSource represents disk source info for a server
type APIServerDiskSource struct {
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

// APIServerTemplateInfos represents OS template information
type APIServerTemplateInfos struct {
	UUID    string `form:"uuid" json:"uuid"`
	Type    string `form:"type" json:"type"`
	OS      string `form:"os" json:"os"`
	Version string `form:"version" json:"version"`
}

// APIServerImageInfos represents image information for a server
type APIServerImageInfos struct {
	UUID         string              `form:"uuid" json:"uuid"`
	Name         string              `form:"name" json:"name"`
	State        string              `form:"state" json:"state"`
	CompanyUUID  string              `form:"company_uuid" json:"company_uuid"`
	Owner        string              `form:"owner" json:"owner"`
	CreatedAt    int64               `form:"created_at" json:"created_at"`
	DiskSize     APISize             `form:"disk_size" json:"disk_size"`
	TemplateUUID string              `form:"template_uuid" json:"template_uuid"`
	OS           APIImageOSFullInfos `form:"os" json:"os"`
}

// APIImageOSFullInfos represents OS info within an image
type APIImageOSFullInfos struct {
	Type    string `form:"type" json:"type"`
	Name    string `form:"name" json:"name"`
	Version string `form:"version" json:"version"`
}

// APIISO represents an ISO attached to a server
type APIISO struct {
	Protocol string `json:"protocol"`
	ISOPath  string `json:"iso_path"`
}

// APIKvmIP represents KVM IP status
type APIKvmIP struct {
	Status string `json:"status"`
	URI    string `json:"uri"`
}

// APISnapshot represents a server snapshot
type APISnapshot struct {
	UUID      string  `json:"uuid"`
	CreatedAt int64   `json:"created_at"`
	Name      string  `json:"name"`
	Size      APISize `json:"size"`
	State     string  `json:"state"`
}

// APIReturn represents a standard API response
type APIReturn struct {
	Code    string `json:"code" binding:"required"`
	Success string `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}
