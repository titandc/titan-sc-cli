package run

import (
	"errors"
	"fmt"
	"strings"
	"titan-sc/api"

	"github.com/spf13/cobra"
)

const (
	deleteServerReasonCLI = "cli"
)

// Plan
const (
	SC1 = "SC1"
	SC2 = "SC2"
	SC3 = "SC3"
)

// Items type
const (
	ItemTypeCPU  = "CPU"
	ItemTypeRAM  = "RAM"
	ItemTypeDisk = "DISK"
	ItemTypeMac  = "MAC"
	ItemTypeOS   = "OS"
)

const (
	OSTypeWindows = "windows"
)

type CreateServerInfo struct {
	plan        string
	templateOID string
	password    string
	sshKeysName string
	quantity    int
	cpu         int
	ram         int
	disk        int
	template    *api.Template
}

var (
	ErrCreateServerPlanInvalid = errors.New("create server fail, plan unknown")
	ErrCreateServerEmptyAuth   = errors.New("create server fail, server authentication is empty")
)

func (run *RunMiddleware) ServerChangeName(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	newServerName, _ := cmd.Flags().GetString("name")
	apiReturn, err := run.API.ServerChangeName(newServerName, serverOID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	companyOID, err := run.GetDefaultCompanyOID(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}

	servers, apiReturn, err := run.API.ServerList(companyOID)
	if err != nil || apiReturn != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}

	if run.JSONOutput {
		// Transform internal site names to public names in JSON output
		for i := range servers {
			if servers[i].Drp != nil {
				servers[i].Drp.ActiveSite = mapSiteToPublic(servers[i].Drp.ActiveSite)
				servers[i].Drp.Site = mapSiteToPublic(servers[i].Drp.Site)
			}
			if servers[i].Site != "" {
				servers[i].Site = mapSiteToPublic(servers[i].Site)
			}
		}
		printAsJson(servers)
	} else {
		table := NewTable("NAME", "PLAN", "STATE", "DRP", "OS", "UUID", "OID")
		table.SetNoColor(!run.Color)

		for _, server := range servers {
			var osInfos string
			if server.Items.OS.Template != nil {
				osInfos = fmt.Sprintf("%s %s", server.Items.OS.Template.OS, server.Items.OS.Template.Version)
			}
			stateRaw := ""
			if server.State != nil {
				stateRaw = *server.State
			}

			var stateColorFn func(string) string
			if run.Color {
				stateColorFn = StateColorFn(stateRaw)
			}

			// DRP status indicator
			drpStatus, drpColorFn := run.getDrpStatusIndicator(server.Drp)

			table.AddRow(
				ColName(server.Name),
				Col(server.Items.CPU.Plan),
				ColColor(stateRaw, stateColorFn),
				ColColor(drpStatus, drpColorFn),
				Col(osInfos),
				Col(server.UUID),
				ColOID(server.OID),
			)
		}

		table.Print()
	}
}

func (run *RunMiddleware) ServerDetail(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")

	server, apiReturn, err := run.API.GetServerOID(serverOID)
	if err != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}
	if run.JSONOutput {
		// Transform internal site names to public names in JSON output
		if server.Drp != nil {
			server.Drp.ActiveSite = mapSiteToPublic(server.Drp.ActiveSite)
			server.Drp.Site = mapSiteToPublic(server.Drp.Site)
		}
		if server.Site != "" {
			server.Site = mapSiteToPublic(server.Site)
		}
		printAsJson(server)
	} else {
		run.printServerDetail(server)
	}
}

func (run *RunMiddleware) printServerDetail(server *api.ServerDetail) {
	state := ""
	if server.State != nil {
		state = GetStateColorized(run.Color, *server.State)
	}

	fmt.Printf("%s\n"+
		"  OID: %s\n"+
		"  UUID: %s\n"+
		"  Name: %s\n"+
		"  State: %s\n"+
		"  Created: %s\n",
		run.Colorize("Server:", "cyan"), server.OID, server.UUID, run.Colorize(server.Name, "cyan"), state, run.Colorize(DatePtrFormat(server.CreatedAt), "dim"))
	if server.UpdatedAt != nil && *server.UpdatedAt > 0 {
		fmt.Printf("  Updated: %s\n", run.Colorize(DatePtrFormat(server.UpdatedAt), "dim"))
	}
	if server.Site != "" {
		fmt.Printf("  Site: %s\n", mapSiteToPublic(server.Site))
	}

	fmt.Printf("%s\n"+
		"  Name: %s\n"+
		"  OID: %s\n"+
		"  Project OID: %s\n",
		run.Colorize("Company:", "cyan"), run.Colorize(server.CompanyName, "cyan"), server.Company, server.ProjectOID)

	cpuTotal := server.Items.CPU.ItemUnit.Value * uint64(server.Items.CPU.Quantity)
	ramTotal := server.Items.RAM.ItemUnit.Value * uint64(server.Items.RAM.Quantity)
	diskTotal := server.Items.DISK.ItemUnit.Value * uint64(server.Items.DISK.Quantity)
	fmt.Printf("%s\n"+
		"  Plan: %s\n"+
		"  CPU: %s\n"+
		"  RAM: %s\n"+
		"  Disk: %s\n",
		run.Colorize("Resources:", "cyan"), server.Items.CPU.Plan,
		run.Colorize(fmt.Sprintf("%d %s", cpuTotal, server.Items.CPU.ItemUnit.Unit), "green"),
		run.Colorize(fmt.Sprintf("%d %s", ramTotal, server.Items.RAM.ItemUnit.Unit), "green"),
		run.Colorize(fmt.Sprintf("%d %s", diskTotal, server.Items.DISK.ItemUnit.Unit), "green"))
	if server.Items.DISK.Description != "" {
		fmt.Printf("  Disk I/O: %s\n", server.Items.DISK.Description)
	}

	if server.Items.OS.Template != nil {
		template := server.Items.OS.Template
		fmt.Printf("%s\n"+
			"  OS: %s %s\n"+
			"  Type: %s\n"+
			"  Template OID: %s\n",
			run.Colorize("Operating System:", "cyan"), template.OS, template.Version, template.Type, template.OID)
	}

	fmt.Printf("%s\n"+
		"  MAC Address: %s\n"+
		"  Bandwidth: %s\n",
		run.Colorize("Network:", "cyan"), server.Items.MAC.Name, server.Items.MAC.Description)
	if len(server.Items.MAC.SubItems) > 0 {
		for _, item := range server.Items.MAC.SubItems {
			if item.IP == nil {
				continue
			}
			primary := ""
			if item.Primary {
				primary = run.Colorize(" (primary)", "green")
			}
			fmt.Printf("  IP: %s%s\n", run.Colorize(item.IP.Address, "cyan"), primary)
			if item.IP.Reverse != "" {
				fmt.Printf("    Reverse: %s\n", run.Colorize(item.IP.Reverse, "dim"))
			}
		}
	}

	if server.Authentication != nil {
		fmt.Printf("%s\n"+
			"  Login: %s\n",
			run.Colorize("Authentication:", "cyan"), run.Colorize(server.Authentication.UserLogin, "cyan"))
	}

	fmt.Printf("%s\n", run.Colorize("Infrastructure:", "cyan"))
	if server.Material != nil && server.Material.Hostname != "" {
		fmt.Printf("  Hypervisor: %s\n", server.Material.Hostname)
	}
	fmt.Printf("  Hypervisor OID: %s\n", server.Hypervisor)

	if server.ISOsOID != nil && len(server.ISOsOID) > 0 {
		fmt.Printf("%s\n", run.Colorize("Mounted ISOs:", "cyan"))
		for _, iso := range server.ISOsOID {
			fmt.Printf("  - %s\n", iso)
		}
	}

	if server.Snapshots != nil && len(*server.Snapshots) > 0 {
		fmt.Printf("%s %s\n", run.Colorize("Snapshots:", "cyan"), run.Colorize(fmt.Sprintf("%d", len(*server.Snapshots)), "green"))
		for _, snap := range *server.Snapshots {
			fmt.Printf("  - Name: %s\n"+
				"    OID: %s\n"+
				"    Size: %d %s\n"+
				"    State: %s\n",
				run.Colorize(snap.Name, "cyan"), snap.OID, snap.Size.Value, snap.Size.Unit, GetStateColorized(run.Color, snap.State))
		}
	}

	if server.Notifications != nil && len(*server.Notifications) > 0 {
		fmt.Printf("%s %s\n", run.Colorize("Notifications:", "yellow"), run.Colorize(fmt.Sprintf("%d", len(*server.Notifications)), "yellow"))
		for _, notification := range *server.Notifications {
			fmt.Printf("  - Title: %s\n"+
				"    Message: %s\n",
				run.Colorize(notification.Title, "yellow"), notification.Message)
		}
	}

	if server.Terminations != nil && len(*server.Terminations) > 0 {
		fmt.Printf("Scheduled Terminations: %d\n", len(*server.Terminations))
		for _, term := range *server.Terminations {
			fmt.Printf("  - Type: %s\n"+
				"    Scheduled: %s\n",
				term.Type, DateFormat(term.Date.ScheduleDate))
		}
	}

	// DRP Section
	if server.Drp != nil && server.Drp.Enabled {
		fmt.Printf("%s\n", run.Colorize("Disaster Recovery Plan (DRP):", "cyan"))
		fmt.Printf("  Enabled: %s\n", run.Colorize("Yes", "green"))
		fmt.Printf("  Status: %s\n", getDrpStatusText(server.Drp.Status))
		fmt.Printf("  Active Site: %s\n", run.Colorize(mapSiteToPublic(server.Drp.ActiveSite), "cyan"))
		fmt.Printf("  Interval: %d minutes\n", server.Drp.Interval)

		// Mirroring state per site
		if server.Drp.MirroringState != nil {
			fmt.Printf("  Mirroring States:\n")
			for site, state := range server.Drp.MirroringState {
				stateText := getMirrorStateText(state)
				fmt.Printf("    %s: %s\n", mapSiteToPublic(site), stateText)
			}
		}

		// Last sync time
		if server.Drp.MirroringLastSync != nil && server.Drp.MirroringLastSync.IsSet() {
			fmt.Printf("  Last Sync: %s\n", FlexTimestampFormat(server.Drp.MirroringLastSync))
		}

		// Pending operation
		if server.Drp.PendingOperation != "" {
			fmt.Printf("  Pending Operation: %s\n", run.Colorize(server.Drp.PendingOperation, "yellow"))
			if server.Drp.PendingOperationAt != nil && server.Drp.PendingOperationAt.IsSet() {
				fmt.Printf("  Operation Started: %s\n", FlexTimestampFormat(server.Drp.PendingOperationAt))
			}
			if server.Drp.PendingOperationBy != "" {
				fmt.Printf("  Operation By: %s\n", server.Drp.PendingOperationBy)
			}
		}

		// Last failover
		if server.Drp.LastFailoverAt != nil && server.Drp.LastFailoverAt.IsSet() {
			fmt.Printf("  Last Failover: %s", FlexTimestampFormat(server.Drp.LastFailoverAt))
			if server.Drp.LastFailoverType != "" {
				fmt.Printf(" (%s)", server.Drp.LastFailoverType)
			}
			fmt.Printf("\n")
		}

		// Last resync
		if server.Drp.LastResyncAt != nil && server.Drp.LastResyncAt.IsSet() {
			fmt.Printf("  Last Resync: %s\n", FlexTimestampFormat(server.Drp.LastResyncAt))
		}

		// Last operation result
		if server.Drp.LastOperationResult != "" {
			fmt.Printf("  Last Result: %s\n", server.Drp.LastOperationResult)
		}

		// Last error
		if server.Drp.LastError != "" {
			fmt.Printf("  Last Error: %s\n", run.Colorize(server.Drp.LastError, "red"))
		}

		// Split-Brain warning
		if server.Drp.SplitBrain || server.Drp.Status == api.DrpStatusSplitBrain {
			fmt.Printf("  %s\n", run.Colorize("⚠ WARNING: Split-Brain detected! Manual intervention required.", "red"))
		}

		// Requires attention
		if server.Drp.RequiresAttention && !server.Drp.SplitBrain {
			fmt.Printf("  %s\n", run.Colorize("⚠ DRP requires attention", "yellow"))
		}
	} else if server.Drp != nil && !server.Drp.Enabled {
		fmt.Printf("%s\n", run.Colorize("Disaster Recovery Plan (DRP):", "cyan"))
		fmt.Printf("  Enabled: %s\n", run.Colorize("No", "yellow"))
	}
}

func (run *RunMiddleware) ServerStart(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	apiReturn, err := run.API.ServerStateAction("start", serverOID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerStop(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	apiReturn, err := run.API.ServerStateAction("stop", serverOID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerRestart(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	apiReturn, err := run.API.ServerStateAction("reboot", serverOID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerHardstop(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	apiReturn, err := run.API.ServerStateAction("hardstop", serverOID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerMountISO(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	uriISO, _ := cmd.Flags().GetString("uri")
	serverOID, _ := cmd.Flags().GetString("server-oid")

	apiReturn, err := run.API.ServerMountISO(uriISO, serverOID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerUmountISO(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	isoOID, _ := cmd.Flags().GetString("iso-oid")
	apiReturn, err := run.API.ServerUmountISO(serverOID, isoOID)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerListTemplates(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	templates, apiReturn, err := run.API.ListTemplates()
	// Render error output
	if err != nil || apiReturn != nil {
		run.handleErrorAndGenericOutput(apiReturn, err)
		return
	}

	if run.JSONOutput {
		// Create output without is_image field (redundant - can be inferred from os: "image")
		type TemplateOutput struct {
			OS       string         `json:"os"`
			Versions []api.Template `json:"versions"`
		}
		output := make([]TemplateOutput, 0, len(templates))
		for _, t := range templates {
			if len(t.Versions) == 0 {
				continue
			}
			output = append(output, TemplateOutput{
				OS:       t.OS,
				Versions: t.Versions,
			})
		}
		printAsJson(output)
	} else {
		// Separate system templates from user images
		var systemTemplates []api.TemplateOSItem
		var userImages []api.TemplateOSItem

		for _, template := range templates {
			if len(template.Versions) == 0 {
				continue
			}
			if template.IsImage {
				userImages = append(userImages, template)
			} else {
				systemTemplates = append(systemTemplates, template)
			}
		}

		// Print system templates
		if len(systemTemplates) > 0 {
			fmt.Println("System Templates:")
			fmt.Println("─────────────────")
			for _, template := range systemTemplates {
				fmt.Printf("  %s:\n", template.OS)
				for _, version := range template.Versions {
					fmt.Printf("    %s (%s)\n", version.Version, version.OID)
				}
			}
		}

		// Print user images
		if len(userImages) > 0 {
			if len(systemTemplates) > 0 {
				fmt.Println()
			}
			fmt.Println("User Images:")
			fmt.Println("────────────")
			for _, template := range userImages {
				for _, version := range template.Versions {
					name := version.Name
					if name == "" {
						name = fmt.Sprintf("%s %s", version.OS, version.Version)
					}
					fmt.Printf("  %s\n", name)
					fmt.Printf("    OID: %s\n", version.OID)
					fmt.Printf("    Base: %s %s\n", version.OS, version.Version)
					if version.ImageInfo != nil {
						fmt.Printf("    Disk Size: %d GB\n", version.ImageInfo.DiskSize)
					}
				}
			}
		}
	}
}

func (run *RunMiddleware) ServerScheduleTermination(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	serverOID, _ := cmd.Flags().GetString("server-oid")
	apiReturn, err := run.API.ServerScheduleTermination(serverOID, deleteServerReasonCLI)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) ServerReset(cmd *cobra.Command, args []string) {
	_ = args
	var err error
	run.ParseGlobalFlags(cmd)
	reset := new(api.ResetServer)

	serverOID, _ := cmd.Flags().GetString("server-oid")
	sshKeys, _ := cmd.Flags().GetString("ssh-keys-name")
	reset.TemplateOID, _ = cmd.Flags().GetString("template-oid")
	reset.UserPassword, _ = cmd.Flags().GetString("password")

	reset.UserSSHKeys, err = run.serverSearchSSHKeys(sshKeys)
	if err != nil {
		run.OutputError(err)
		return
	}

	apiReturn, err := run.API.ServerReset(serverOID, reset)
	run.handleErrorAndGenericOutput(apiReturn, err)
}

func (run *RunMiddleware) serverSearchSSHKeys(sshKeysName string) ([]string, error) {
	if sshKeysName == "" {
		return []string{}, nil
	}

	// Get current user to use as target_oid
	user, err := run.API.GetUserInfos()
	if err != nil {
		return []string{}, err
	}

	sshKeysList, err := run.API.GetSSHKeyList(user.OID)
	if err != nil {
		return []string{}, err
	}

	sshKeys := make([]string, 0)
	for _, keyRequest := range strings.Split(strings.TrimSpace(sshKeysName), ",") {
		found := false
		for _, keyExist := range sshKeysList {
			if keyExist.Name == keyRequest {
				sshKeys = append(sshKeys, keyExist.Value)
				found = true
				break
			}
		}
		if !found {
			return []string{}, fmt.Errorf("SSH keys name <%s> not found.\n", keyRequest)
		}
	}
	return sshKeys, nil
}

// ServerCreateResult is returned as JSON for non-human output
type ServerCreateResult struct {
	CartOID  string  `json:"cart_oid"`
	PriceHT  float64 `json:"price_ht"`
	PriceTTC float64 `json:"price_ttc"`
	Status   string  `json:"status"`
}

func (run *RunMiddleware) ServerCreate(cmd *cobra.Command, _ []string) {
	info := CreateServerInfo{}

	err := info.parse(cmd)
	if err != nil {
		run.OutputError(err)
		return
	}

	paymentMethodOID, _ := cmd.Flags().GetString("payment-method")
	confirmPayment, _ := cmd.Flags().GetBool("confirm-payment")

	if !confirmPayment {
		run.OutputError(errors.New("payment not confirmed: add --confirm-payment flag to proceed"))
		return
	}

	cart := &api.AddServerCart{
		Quantity: info.quantity,
	}

	user, err := run.API.GetUserInfos()
	if err != nil {
		run.OutputError(err)
		return
	}
	cart.CartOID = user.OID

	if err = run.setServerItems(cart, &info); err != nil {
		run.OutputError(err)
		return
	}

	if info.template.Type != OSTypeWindows {
		if err = run.setServerAuth(cart, info.password, info.sshKeysName); err != nil {
			run.OutputError(err)
			return
		}
	}

	cartOID, err := run.API.CreateServerCart(cart)
	if err != nil {
		run.OutputError(err)
		return
	}

	// Get price preview
	price, err := run.API.GetCartPrice(cartOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	priceHT := float64(price.Amount.HT) / 100
	priceTTC := float64(price.Amount.TTC) / 100

	if !run.JSONOutput {
		fmt.Printf("Cart created: %s\n", cartOID)
		fmt.Printf("Price: %.2f€ TTC (%.2f€ HT)\n", priceTTC, priceHT)
		fmt.Println("Processing payment...")
	}

	subscriptionOID, _ := cmd.Flags().GetString("subscription-oid")

	if err = run.API.BuyCart(cartOID, paymentMethodOID, subscriptionOID); err != nil {
		run.OutputError(err)
		return
	}

	if !run.JSONOutput {
		fmt.Println("Server order completed successfully.")
	} else {
		result := ServerCreateResult{
			CartOID:  cartOID,
			PriceHT:  priceHT,
			PriceTTC: priceTTC,
			Status:   "completed",
		}
		printAsJson(result)
	}
}

// Plan minimum resources (disk is in GB, internally divided by 10 for API)
var planMinResources = map[string]struct{ cpu, ram, diskGB int }{
	SC1: {1, 1, 10},
	SC2: {4, 4, 80},
	SC3: {6, 8, 100},
}

func (a *CreateServerInfo) parse(cmd *cobra.Command) error {
	a.plan, _ = cmd.Flags().GetString("plan")
	a.templateOID, _ = cmd.Flags().GetString("template-oid")
	a.password, _ = cmd.Flags().GetString("password")
	a.sshKeysName, _ = cmd.Flags().GetString("ssh-keys-name")
	a.quantity, _ = cmd.Flags().GetInt("quantity")
	a.cpu, _ = cmd.Flags().GetInt("cpu")
	a.ram, _ = cmd.Flags().GetInt("ram")
	a.disk, _ = cmd.Flags().GetInt("disk")

	a.plan = strings.ToUpper(a.plan)
	if a.plan != SC1 && a.plan != SC2 && a.plan != SC3 {
		return ErrCreateServerPlanInvalid
	}

	// Validate minimum resources for plan
	min := planMinResources[a.plan]
	if a.cpu > 0 && a.cpu < min.cpu {
		return fmt.Errorf("plan %s requires minimum %d CPU(s), got %d", a.plan, min.cpu, a.cpu)
	}
	if a.ram > 0 && a.ram < min.ram {
		return fmt.Errorf("plan %s requires minimum %d GB RAM, got %d", a.plan, min.ram, a.ram)
	}
	if a.disk > 0 {
		if a.disk%10 != 0 {
			return fmt.Errorf("disk must be a multiple of 10 GB, got %d", a.disk)
		}
		if a.disk < min.diskGB {
			return fmt.Errorf("plan %s requires minimum %d GB disk, got %d", a.plan, min.diskGB, a.disk)
		}
	}

	return nil
}

func (run *RunMiddleware) setServerAuth(cart *api.AddServerCart, password, sshKeysName string) error {
	var err error

	cart.Auth.UserPassword = password
	if sshKeysName != "" {
		cart.Auth.SSHKeys, err = run.getSSHKeysValue(sshKeysName)
		if err != nil {
			return err
		}
	}
	if cart.Auth.UserPassword == "" && len(cart.Auth.SSHKeys) == 0 {
		return ErrCreateServerEmptyAuth
	}

	return nil
}

func (run *RunMiddleware) setServerItems(cart *api.AddServerCart, info *CreateServerInfo) error {
	items, err := run.API.ListItems()
	if err != nil {
		return err
	}

	for _, itemType := range []string{ItemTypeCPU, ItemTypeRAM, ItemTypeDisk} {
		if err = getPackageItems(cart, info.plan, itemType, items); err != nil {
			return err
		}
	}

	if err = getItems(cart, info.plan, ItemTypeMac, items, false); err != nil {
		return err
	}

	info.template, err = run.getOSItem(cart, info.templateOID, items)
	if err != nil {
		return err
	}

	// Calculate addons needed (total - plan minimum)
	min := planMinResources[info.plan]
	cpuAddons := 0
	ramAddons := 0
	diskAddons := 0
	if info.cpu > min.cpu {
		cpuAddons = info.cpu - min.cpu
	}
	if info.ram > min.ram {
		ramAddons = info.ram - min.ram
	}
	// Convert disk from GB to units (10GB per unit)
	minDiskUnits := min.diskGB / 10
	userDiskUnits := info.disk / 10
	if userDiskUnits > minDiskUnits {
		diskAddons = userDiskUnits - minDiskUnits
	}

	if err = getAddonItem(cart, info.plan, ItemTypeCPU, items, cpuAddons); err != nil {
		return err
	}

	if err = getAddonItem(cart, info.plan, ItemTypeRAM, items, ramAddons); err != nil {
		return err
	}

	if err = getAddonItem(cart, info.plan, ItemTypeDisk, items, diskAddons); err != nil {
		return err
	}

	return nil
}

func getAddonItem(cart *api.AddServerCart, plan, itemType string, items []api.ItemLimited, addonQuantity int) error {
	if addonQuantity == 0 {
		return nil
	}

	if err := getItems(cart, plan, itemType, items, false); err != nil {
		return err
	}
	i := len(cart.Items) - 1
	cart.Items[i].Quantity = addonQuantity
	return nil
}

func getPackageItems(cart *api.AddServerCart, plan, itemType string, items []api.ItemLimited) error {
	if err := getItems(cart, plan, itemType, items, true); err != nil {
		return err
	}
	if itemType == ItemTypeCPU || itemType == ItemTypeRAM || itemType == ItemTypeDisk {
		i := len(cart.Items) - 1
		setQuantityByPlan(&cart.Items[i], plan, itemType)
	}
	return nil
}

func setQuantityByPlan(cartItem *api.AddServerCartItem, plan, itemType string) {
	min := planMinResources[plan]
	switch itemType {
	case ItemTypeCPU:
		cartItem.Quantity = min.cpu
	case ItemTypeRAM:
		cartItem.Quantity = min.ram
	case ItemTypeDisk:
		cartItem.Quantity = min.diskGB / 10 // Convert GB to units
	}
}

func getItems(cart *api.AddServerCart, plan, itemType string, items []api.ItemLimited, packageFlag bool) error {
	item := searchItems(packageFlag, plan, itemType, items)
	if item == nil {
		return fmt.Errorf("%s item not found", itemType)
	}
	cart.Items = append(cart.Items, api.AddServerCartItem{
		OID:      item.OID,
		Quantity: 1,
	})
	return nil
}

func (run *RunMiddleware) getOSItem(cart *api.AddServerCart, templateOID string, items []api.ItemLimited) (*api.Template, error) {
	var err error
	template, err := run.API.GetTemplateByOID(templateOID)
	if err != nil {
		return nil, err
	}

	item := searchItems(false, "", ItemTypeOS, items)
	if item == nil {
		return nil, fmt.Errorf("%s item not found", ItemTypeOS)
	}

	cart.Items = append(cart.Items, api.AddServerCartItem{
		OID:       item.OID,
		Quantity:  1,
		TargetOID: template.OID,
	})
	return template, nil
}

func searchItems(itemPackage bool, plan, itemType string, items []api.ItemLimited) *api.ItemLimited {
	for i := range items {
		if plan == items[i].Plan && itemPackage == items[i].Package && itemType == items[i].Type {
			return &items[i]
		}
	}
	return nil
}

func (run *RunMiddleware) getSSHKeysValue(sshKeysName string) ([]string, error) {
	var values []string

	// Get current user to use as target_oid
	user, err := run.API.GetUserInfos()
	if err != nil {
		return nil, err
	}

	sshKeysList, err := run.API.GetSSHKeyList(user.OID)
	if err != nil {
		return nil, err
	}

	for _, name := range strings.Split(sshKeysName, ",") {
		find := false
		for _, sshKey := range sshKeysList {
			if sshKey.Name == name {
				find = true
				values = append(values, sshKey.Value)
			}
		}
		if !find {
			return nil, fmt.Errorf("SSHKey %s not found", name)
		}
	}

	return values, nil
}
