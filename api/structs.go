package api

import (
	"encoding/json"
	"strconv"
)

// FlexTimestamp handles JSON timestamps that can be either int64 (Unix ms) or string
type FlexTimestamp struct {
	Value *int64
}

func (ft *FlexTimestamp) UnmarshalJSON(data []byte) error {
	// Try as number first
	var num int64
	if err := json.Unmarshal(data, &num); err == nil {
		ft.Value = &num
		return nil
	}

	// Try as string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == "" {
			ft.Value = nil
			return nil
		}
		// Try to parse string as number
		if num, err := strconv.ParseInt(str, 10, 64); err == nil {
			ft.Value = &num
			return nil
		}
		// Keep as nil if can't parse
		ft.Value = nil
		return nil
	}

	// Null or invalid - keep as nil
	ft.Value = nil
	return nil
}

func (ft FlexTimestamp) MarshalJSON() ([]byte, error) {
	if ft.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*ft.Value)
}

// IsSet returns true if the timestamp has a value
func (ft *FlexTimestamp) IsSet() bool {
	return ft != nil && ft.Value != nil && *ft.Value > 0
}

// Get returns the timestamp value or 0 if not set
func (ft *FlexTimestamp) Get() int64 {
	if ft == nil || ft.Value == nil {
		return 0
	}
	return *ft.Value
}

type Base struct {
	OID       string `json:"oid"`
	CreatedAt *int64 `json:"created_at,omitempty"`
	UpdatedAt *int64 `json:"updated_at,omitempty"`
	DeletedAt *int64 `json:"deleted_at,omitempty"`
}

type ServerDetail struct {
	Base
	Name                string             `json:"name"`
	Items               ServerItemsView    `json:"items"`
	Company             string             `json:"company" validate:"required,objectID"`
	CompanyName         string             `json:"company_name,omitempty"`
	Owner               string             `json:"owner" validate:"required,objectID"`
	ProjectOID          string             `json:"project_oid" validate:"objectID_e"`
	ForcedHypervisorOID string             `json:"forced_hypervisor_oid" validate:"objectID_e"`
	Hypervisor          string             `json:"hypervisor,omitempty"`
	ISOsOID             []string           `json:"isos_oid" validate:"omitempty,dive,objectID"`
	State               *string            `json:"state,omitempty"`
	StateOID            *string            `json:"state_oid,omitempty"`
	StatementOID        *string            `json:"statement_oid,omitempty"`
	Authentication      *Authentication    `json:"authentication,omitempty"`
	Material            *Material          `json:"material,omitempty"`
	KVM                 *KvmIPView         `json:"kvm_ip,omitempty"`
	Snapshots           *[]SnapshotDetail  `json:"snapshots,omitempty"`
	Demo                bool               `json:"demo"`
	Disable             bool               `json:"disable"`
	Notifications       *[]Notification    `json:"notifications,omitempty"`
	Terminations        *[]ScheduledAction `json:"terminations,omitempty"`
	UUID                string             `json:"uuid"`
	Site                string             `json:"site,omitempty"`
	Drp                 *Drp               `json:"drp,omitempty"`
	OutsourcingDate     int64              `json:"outsourcing_date,omitempty"`
	OutsourcingLevel    int                `json:"outsourcing_level,omitempty"`
	OutsourcingEnd      int64              `json:"outsourcing_end,omitempty"`
	TagInfo             []TagInfo          `json:"tag_info,omitempty"`
}

type Drp struct {
	Enabled             bool           `json:"enabled"`
	Status              int            `json:"status,omitempty"`
	ActiveSite          string         `json:"active_site,omitempty"`
	Site                string         `json:"site,omitempty"`
	DrpURL              string         `json:"drp_url,omitempty"`
	Interval            int            `json:"interval,omitempty"`
	StartTime           *FlexTimestamp `json:"start_time,omitempty"`
	SplitBrain          bool           `json:"split_brain,omitempty"`
	RequiresAttention   bool           `json:"requires_attention,omitempty"`
	LastError           string         `json:"last_error,omitempty"`
	LastFailoverAt      *FlexTimestamp `json:"last_failover_at,omitempty"`
	LastFailoverType    string         `json:"last_failover_type,omitempty"`
	LastResyncAt        *FlexTimestamp `json:"last_resync_at,omitempty"`
	LastOperationResult string         `json:"last_operation_result,omitempty"`
	PendingOperation    string         `json:"pending_operation,omitempty"`
	PendingOperationAt  *FlexTimestamp `json:"pending_operation_at,omitempty"`
	PendingOperationBy  string         `json:"pending_operation_by,omitempty"`
	MirroringState      map[string]int `json:"mirroring_state,omitempty"`
	MirroringLastSync   *FlexTimestamp `json:"mirroring_last_sync,omitempty"`
}

// DRP Status constants
const (
	DrpStatusOff        = 0 // DRP disabled or error
	DrpStatusOK         = 1 // DRP active and healthy
	DrpStatusPending    = 2 // Operation in progress
	DrpStatusSplitBrain = 3 // Split-brain detected, requires attention
)

// DRP Mirroring State constants
const (
	MirrorStateUnknown = 0
	MirrorStateError   = 1
	MirrorStateSyncing = 2
	MirrorStateStandby = 4
	MirrorStatePrimary = 6
)

// DrpStatus is the response from GET /server/{id}/drp/status
type DrpStatus struct {
	ServerOID           string         `json:"server_oid"`
	ServerUUID          string         `json:"server_uuid"`
	ServerName          string         `json:"server_name"`
	Enabled             bool           `json:"enabled"`
	Status              int            `json:"status"`
	Interval            int            `json:"interval,omitempty"`
	StartTime           *FlexTimestamp `json:"start_time,omitempty"`
	ActiveSite          string         `json:"active_site,omitempty"`
	SplitBrain          bool           `json:"split_brain,omitempty"`
	RequiresAttention   bool           `json:"requires_attention,omitempty"`
	LastError           string         `json:"last_error,omitempty"`
	LastFailoverAt      *FlexTimestamp `json:"last_failover_at,omitempty"`
	LastFailoverType    string         `json:"last_failover_type,omitempty"`
	LastResyncAt        *FlexTimestamp `json:"last_resync_at,omitempty"`
	LastOperationResult string         `json:"last_operation_result,omitempty"`
	PendingOperation    string         `json:"pending_operation,omitempty"`
	PendingOperationAt  *FlexTimestamp `json:"pending_operation_at,omitempty"`
	PendingOperationBy  string         `json:"pending_operation_by,omitempty"`
	IPs                 []DrpIPStatus  `json:"ips,omitempty"`
}

// DrpIPStatus represents IP DRP state
type DrpIPStatus struct {
	IPOID           string         `json:"ip_oid"`
	Address         string         `json:"address"`
	Version         int            `json:"version"`
	CurrentSite     string         `json:"current_site"`
	LastSwitchAt    *FlexTimestamp `json:"last_switch_at,omitempty"`
	LastSwitchError string         `json:"last_switch_error,omitempty"`
	MACAddress      string         `json:"mac_address,omitempty"`
}

// DrpOperationResult is the response from DRP operation endpoints
type DrpOperationResult struct {
	Success    bool   `json:"success"`
	Operation  string `json:"operation,omitempty"`
	ServerOID  string `json:"server_oid,omitempty"`
	Message    string `json:"message,omitempty"`
	SourceSite string `json:"source_site,omitempty"`
	TargetSite string `json:"target_site,omitempty"`
	Error      string `json:"error,omitempty"`
}

// DrpFailoverHardRequest is the request body for hard failover
type DrpFailoverHardRequest struct {
	TargetSite string `json:"target_site"`
}

// DrpResyncRequest is the request body for resync
type DrpResyncRequest struct {
	AuthoritativeSite string `json:"authoritative_site"`
}

// NetworkDrp represents DRP configuration for a network
type NetworkDrp struct {
	Enabled bool   `json:"enabled"`
	Site    string `json:"site,omitempty"`
}

type TagInfo struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Material struct {
	Hostname string `json:"hostname"`
}

type Notification struct {
	Base
	TargetOID       []string        `json:"target_oid"`
	Severity        int             `json:"severity"`
	Title           string          `json:"title"`
	Message         string          `json:"message"`
	MandatoryReboot bool            `json:"mandatory_reboot"`
	Servers         *[]ServerDetail `json:"servers,omitempty"`

	// Deprecated
	UUID string `json:"-"`
}

type ScheduledAction struct {
	OID        string              `json:"oid"`
	Type       string              `json:"type" validate:"required,scheduled_action_type"`
	Date       ScheduledActionDate `json:"date" validate:"required"`
	TargetOID  string              `json:"target_oid" validate:"required,objectID"`
	Terminated bool                `json:"terminated"`
}

type ScheduledActionDate struct {
	RequestedAt  int64  `json:"requested_at"`
	ScheduleDate int64  `json:"schedule_date"`
	State        string `json:"state"`
	ApplicantOID string `json:"applicant_oid"`
}

type ServerItemsView struct {
	OS      ItemLimited  `json:"os"`
	MAC     ItemLimited  `json:"mac"`
	RAM     ItemLimited  `json:"ram"`
	CPU     ItemLimited  `json:"cpu"`
	DISK    ItemLimited  `json:"disk"`
	License *ItemLimited `json:"license,omitempty"`
}

type ItemLimited struct {
	OID         string        `json:"oid,omitempty"`
	Name        string        `json:"name"`
	Code        string        `json:"code"`
	Description string        `json:"description"`
	Type        string        `json:"type"`
	Plan        string        `json:"plan,omitempty"`
	Priority    uint          `json:"priority,omitempty"`
	Primary     bool          `json:"primary,omitempty"`
	Package     bool          `json:"package"`
	Quantity    uint          `json:"quantity"`
	ItemUnit    ItemUnit      `json:"item_unit,omitempty"`
	SubItems    []ItemLimited `json:"sub_items,omitempty"`
	TargetOID   *string       `json:"target_oid,omitempty"`
	Template    *Template     `json:"template,omitempty"`
	IP          *IP           `json:"ip,omitempty"`
}

type IP struct {
	Base
	Address        string `json:"address"`
	Version        int    `json:"version"`
	CompanyOID     string `json:"company_oid"`
	Reverse        string `json:"reverse"`
	DefaultReverse string `json:"default_reverse"`
	ServerOID      string `json:"server_oid"`
	ServerName     string `json:"server_name"`

	// Internal fields - hidden from JSON output
	ICMP           bool   `json:"-"`
	Redirection    string `json:"-"`
	WanProviderOID string `json:"-"`
	Gateway        string `json:"-"`
	Network        string `json:"-"`
}

type TemplateOSItem struct {
	IsImage  bool       `json:"is_image"`
	OS       string     `json:"os"`
	Versions []Template `json:"versions"`
}

type Template struct {
	Base
	UUID       string     `json:"uuid"`
	Name       string     `json:"name,omitempty"`
	OS         string     `json:"os"`
	Version    string     `json:"version"`
	Type       string     `json:"type"`
	ImageInfo  *ImageInfo `json:"image_info,omitempty"`
	Enabled    bool       `json:"enabled"`
	HasLicense *bool      `json:"has_license,omitempty"`
	Licenses   *[]License `json:"licenses,omitempty"`
}

type ImageInfo struct {
	CompanyOID  string `json:"company_oid"`
	OwnerOID    string `json:"owner_oid"`
	TemplateOID string `json:"template_oid"`
	DiskSize    uint64 `json:"disk_size"`
}

type License struct {
	Base
	LicenseID string `json:"license_id"`
	Type      string `json:"type"`
	Number    string `json:"number"`
	Key       string `json:"key"`
	Code      string `json:"code"`
}

type Authentication struct {
	Base
	UserLogin string   `json:"user_login"`
	SSHKeys   []string `json:"-"`
	ServerOID string   `json:"server_oid"`
}

type KvmIPView struct {
	Base
	Deadline     int64   `json:"deadline"`
	URL          string  `json:"url"`
	ServerOID    string  `json:"server_oid"`
	State        string  `json:"state"`
	StatementOID *string `json:"-"`
}

type ItemUnit struct {
	Value uint64 `json:"value"`
	Unit  string `json:"unit"`
}

type ServerUpdateInfos struct {
	Name string `json:"name"`
}

type ServerScheduleTermination struct {
	Reason string `json:"reason"`
}

type NetworkList struct {
	Quota    uint            `json:"quota"`
	Networks []NetworkDetail `json:"networks"`
}

type Network struct {
	Base
	UUID       string       `json:"uuid"`
	CompanyOID string       `json:"company_oid"`
	Name       string       `json:"name"`
	Speed      NetworkSpeed `json:"speed"`
	Ports      uint         `json:"ports"`
	MaxMTU     uint         `json:"max_mtu"`
}

type NetworkDetail struct {
	Base       `bson:",inline"`
	UUID       string                    `json:"uuid"`
	CompanyOID string                    `json:"company_oid"`
	Name       string                    `json:"name"`
	Speed      NetworkSpeed              `json:"speed"`
	Ports      uint                      `json:"ports"`
	MaxMTU     uint                      `json:"max_mtu"`
	Interfaces []NetworkPrivateInterface `json:"interfaces"`
	State      string                    `json:"state,omitempty"`
	StateOID   string                    `json:"statement_oid,omitempty"`
	Site       string                    `json:"site,omitempty"`
	Drp        *NetworkDrp               `json:"drp,omitempty"`
}

type NetworkSpeed struct {
	Unit  string `json:"unit"`
	Value int32  `json:"value"`
}

type NetworkPrivateInterface struct {
	Base
	UUID       string       `json:"uuid"`
	MAC        string       `json:"mac"`
	NetworkOID string       `json:"network_oid"`
	Server     ServerDetail `json:"server"`
}

type Company struct {
	Base
	UUID                 string     `json:"uuid"`
	Name                 string     `json:"name"`
	Avatar               string     `json:"avatar,omitempty"`
	Email                string     `json:"email"`
	Website              string     `json:"website,omitempty"`
	Phone                string     `json:"phone,omitempty"`
	DefaultPaymentMethod *string    `json:"default_payment_method,omitempty"`
	AutoPaymentMethod    *string    `json:"auto_payment_method,omitempty"`
	Addresses            []Address  `json:"addresses,omitempty"`
	AddressShipping      Address    `json:"address_shipping"`
	AddressBilling       Address    `json:"address_billing"`
	VAT                  CompanyVAT `json:"vat"`
	RenewWithCredits     bool       `json:"renew_with_credits,omitempty"`
	Disable              bool       `json:"disable,omitempty"`

	// Internal fields - hidden from JSON output
	ClientNumber string `json:"-"`
	Siret        string `json:"-"`
	NAF          string `json:"-"`
	Discount     string `json:"-"`
	Commercial   string `json:"-"`
	QuotaOID     string `json:"-"`
}

type Address struct {
	Street2     string `json:"street2,omitempty"`
	Street      string `json:"street"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	PostalCode  string `json:"postal_code"`
}

type CompanyVAT struct {
	Number    string `json:"number,omitempty"`
	Valid     bool   `json:"valid"`
	Confirmed bool   `json:"confirmed"`
}

type ServerAction struct {
	Action  string `json:"action"`
	JobName string `json:"job_name,omitempty"` // Optional job name for tracking
}

type KvmIPRequest struct {
	ServerOID string `json:"server_oid"`
	Deadline  uint   `json:"deadline"`
}

type Snapshot struct {
	Base
	Name      string   `json:"name"`
	ServerOID string   `json:"server_oid"`
	Size      ItemUnit `json:"size"`
	UUID      string   `json:"uuid"`
}

type SnapshotDetail struct {
	Snapshot
	State        string `json:"state"`
	StatementOID string `json:"statement_oid"`
}

type Event struct {
	Type          string                  `json:"type"`
	SubType       string                  `json:"sub_type"`
	Action        string                  `json:"action"`
	Status        string                  `json:"status"`
	TargetOID     *string                 `json:"target_oid,omitempty"`
	TargetName    *string                 `json:"target_name,omitempty"`
	Timestamp     string                  `json:"timestamp"`
	UUID          *string                 `json:"uuid,omitempty"`
	FullMessage   *string                 `json:"full_message,omitempty"`
	ServerName    *string                 `json:"server_name,omitempty"`
	ServerOID     *string                 `json:"server_oid,omitempty"`
	CompanyOID    *string                 `json:"company_oid,omitempty"`
	UserOID       *string                 `json:"user_oid,omitempty"`
	UserUUID      *string                 `json:"user_uuid,omitempty"`
	UserEmail     *string                 `json:"user_email,omitempty"`
	UserIP        *string                 `json:"user_ip,omitempty"`
	UserAvatarOID *string                 `json:"avatar_oid,omitempty"`
	Fields        []EventAdditionalFields `json:"fields"`
}

type EventAdditionalFields struct {
	Name     *string `json:"name,omitempty"`
	OldValue *string `json:"old_data,omitempty"`
	NewValue *string `json:"new_data,omitempty"`
	Data     *string `json:"data,omitempty"`
}

type NetworkRename struct {
	Name string `json:"name"`
}

type NetworkOps struct {
	ServerOIDs []string `json:"servers_oid,omitempty"`
	ServerOID  string   `json:"server_oid,omitempty"`
}

type NetworkCreate struct {
	Name       string `json:"name"`
	CompanyOID string `json:"oid,omitempty"` // Company OID (optional, uses current user's company if not set)
}

type KvmIP struct {
	Base
	Deadline  int64  `json:"deadline"`
	URL       string `json:"url"`
	ServerOID string `json:"server_oid"`
}

type ServerMountISORequest struct {
	ISO      string `json:"iso_addr"`
	Protocol string `json:"protocol"`
}

type ServerAddonInfo struct {
	ReducibleItems  *ServerResources `json:"reducible_addons"`
	UpgradableItems []ItemWithPrice  `json:"upgradable_items"`
}

type ServerResources struct {
	CPU  uint
	RAM  uint
	Disk uint
	Plan int `json:"plan"`
}

type ItemWithPrice struct {
	OID               string   `json:"oid"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Enabled           bool     `json:"enabled"`
	Package           bool     `json:"package,omitempty"`
	Discountable      bool     `json:"discountable"`
	RecurringBills    bool     `json:"recurring_bills"`
	Type              string   `json:"type"`
	AccountingAccount string   `json:"accounting_account"`
	PriceUnitHT       int64    `json:"price_unit_ht"`
	Currency          string   `json:"currency"`
	StartQuantity     uint     `json:"start_quantity"`
	ItemUnit          ItemUnit `json:"item_unit"`
	Plan              string   `json:"plan,omitempty"`
}

type Return struct {
	Title string `json:"error,omitempty"`
	// Message is used Title equal ERROR_OTHER
	Message string `json:"message,omitempty"`
	// Data is used if Title equal ERROR_VALIDATION
	Data []ValidationError `json:"data,omitempty"`
}

type ValidationError struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

type IPAttachDetach struct {
	IPs []string `json:"ip"`
}

type APIVersion struct {
	Date    string `json:"release_date,omitempty"`
	Version string `json:"version"`
}

type User struct {
	Base
	Firstname         string           `json:"firstname"`
	Lastname          string           `json:"lastname"`
	Salutation        string           `json:"salutation"`
	Phone             string           `json:"phone"`
	Avatar            string           `json:"avatar"`
	DefaultCompanyOID string           `json:"company_oid"`
	Companies         []UserCompany    `json:"companies,omitempty"`
	Registration      UserRegistration `json:"registration"`
	Email             string           `json:"email"`
	Disabled          bool             `json:"disabled,omitempty"`
	Debug             bool             `json:"debug,omitempty"`
	LastLogin         int64            `json:"last_login"`
	Signature         string           `json:"signature"`
	LatestSignedCGV   *UserCGV         `json:"latest_signed_cgv"`
	Preference        *UserPreferences `json:"preference,omitempty"`

	// Deprecated
	UUID string `bson:"uuid" json:"uuid"`
}

type UserCompany struct {
	OID      string `json:"oid"`
	Role     string `json:"role"`
	RoleOID  string `json:"role_oid"`
	Position string `json:"position"`
}

type UserCGV struct {
	OID  string `json:"oid"`
	Date int64  `json:"date"`
	IP   string `json:"ip"`
}

type UserRegistration struct {
	TwoFA bool `json:"two_fa"`
}

type UserPreferences struct {
	OID                  string                       `json:"oid"`
	TargetOID            string                       `json:"target_oid"`
	ColorMode            string                       `json:"color_mode"`
	PreferredLanguage    string                       `json:"preferred_language"`
	ShowCartConfirmModal bool                         `json:"show_card_confirm_modal"`
	LastDashboardVersion string                       `json:"last_dashboard_version"`
	Tour                 UserTour                     `json:"tour"`
	SortingPreferences   map[string]SortingPreference `json:"sorting_preferences"`
	LastSeenMessage      string                       `json:"last_seen_msg"`
	Dashboard            []UserDashboardItem          `json:"dashboard"`
}

type UserTour struct {
	Global bool `json:"global"`
}

type SortingPreference struct {
	SortField string `json:"sort_field"`
	Sort      int8   `json:"sort"`
	Limit     uint8  `json:"limit"`
}

type UserDashboardItem struct {
	ServerOID  string `json:"server_oid"`
	CompanyOID string `json:"company_oid"`
	Name       string `json:"name"`
	Position   uint32 `json:"position"`
}

type SSHKey struct {
	OID   string `json:"oid"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type AddSSHKey struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeleteSSHKey struct {
	Name string `json:"name"`
}

type ResetServer struct {
	UserPassword string   `json:"user_password,omitempty"`
	UserSSHKeys  []string `json:"user_ssh_keys,omitempty"`
	TemplateOID  string   `json:"template_oid"`
}

type IPUpdateRequest struct {
	Reverse string `json:"reverse"`
}

type AddServerCart struct {
	CartOID  string              `json:"cart_oid"`
	Quantity int                 `json:"quantity,omitempty"`
	Items    []AddServerCartItem `json:"items"`
	Auth     ItemCartAuth        `json:"auth,omitempty"`
}

type ItemCartAuth struct {
	UserPassword string   `bson:"user_password" json:"user_password"`
	SSHKeys      []string `bson:"user_ssh_keys" json:"ssh_keys"`
}

type AddServerCartItem struct {
	OID       string `form:"oid" json:"oid" validate:"required,objectID"`
	Quantity  int    `form:"quantity" json:"quantity" validate:"required,gte=1" default:"1"`
	TargetOID string `form:"target_oid,omitempty" json:"target_oid,omitempty"`
}

type CreateServerCart struct {
	CartOID string `json:"cart_oid"`
}

// BuyCartRequest is the request body for POST /cart/buy
type BuyCartRequest struct {
	CartOID          string `json:"cart_oid"`
	PaymentMethodOID string `json:"payment_method_oid,omitempty"`
	SubscriptionOID  string `json:"subscription_oid,omitempty"` // Add to existing subscription instead of creating new one
	Discount         string `json:"discount,omitempty"`
}

// CartPrice represents the price preview for a cart
type CartPrice struct {
	Amount CartAmount `json:"amount"`
	Debug  bool       `json:"debug"`
}

// CartAmount contains pricing details
type CartAmount struct {
	HT        int64 `json:"ht"`        // Price excluding tax (in cents)
	TTC       int64 `json:"ttc"`       // Price including tax (in cents)
	TVA       int64 `json:"tva"`       // Tax amount (in cents)
	Initial   int64 `json:"initial"`   // Initial price before discount
	Remaining int64 `json:"remaining"` // Remaining to pay
}

// Subscription represents a billing subscription
type Subscription struct {
	OID              string              `json:"oid"`
	Name             string              `json:"name"`
	DocumentNumber   string              `json:"document_number"`
	CompanyOID       string              `json:"company_oid"`
	Company          SubscriptionCompany `json:"company,omitempty"`
	State            string              `json:"state"`
	Frequency        string              `json:"frequency"`
	NextFrequency    string              `json:"next_frequency,omitempty"`
	NextBillingDate  int64               `json:"next_billing_date,omitempty"`
	Amount           SubscriptionAmount  `json:"amount,omitempty"`
	PaymentMethodOID string              `json:"payment_method_oid,omitempty"`
	PaymentDisabled  bool                `json:"payment_disabled,omitempty"`
	CreatedAt        int64               `json:"created_at,omitempty"`
	UpdatedAt        int64               `json:"updated_at,omitempty"`
}

// SubscriptionCompany contains company info for subscription
type SubscriptionCompany struct {
	Name                  string `json:"name"`
	ClientNumber          string `json:"client_number,omitempty"`
	InvoicePaymentGrouped bool   `json:"invoice_payment_grouped,omitempty"`
}

// SubscriptionAmount contains the pricing for a subscription
type SubscriptionAmount struct {
	HT  int64 `json:"ht"`  // Price excluding tax (in cents)
	TTC int64 `json:"ttc"` // Price including tax (in cents)
}
