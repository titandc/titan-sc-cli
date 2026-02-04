package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RegisterFlagCompletions registers dynamic completion functions for common flags
func (cmd *CMD) RegisterFlagCompletions() {
	// Company-scoped completions (use --company-oid or default company)
	cmd.registerServerOIDCompletion()
	cmd.registerNetworkOIDCompletion()
	cmd.registerSubscriptionOIDCompletion()

	// Global/shared completions
	cmd.registerTemplateOIDCompletion()
	cmd.registerCompanyOIDCompletion()

	// User-scoped completions (belong to current user)
	cmd.registerSSHKeyOIDCompletion()
	cmd.registerTokenOIDCompletion()

	// Dependent completions (depend on other flags)
	cmd.registerSnapshotOIDCompletion()
	cmd.registerISOOIDCompletion()

	// Contextual completions (behavior depends on command context)
	cmd.registerIPCompletion()
	cmd.registerNetworkServerOIDCompletion()
}

// getCompanyOIDForCompletion returns the company OID to use for completion.
// It checks if --company-oid was specified, otherwise returns the user's default company.
func (cmd *CMD) getCompanyOIDForCompletion(c *cobra.Command) string {
	// Check if --company-oid was specified in the command
	companyOID, _ := c.Flags().GetString("company-oid")

	// If not specified, get user's default company
	if companyOID == "" {
		user, err := cmd.runMiddleware.API.GetUserInfos()
		if err != nil {
			return ""
		}
		companyOID = user.DefaultCompanyOID
	}

	return companyOID
}

// =============================================================================
// COMPANY-SCOPED COMPLETIONS
// These resources belong to a company, so we use --company-oid or default company
// =============================================================================

// registerServerOIDCompletion registers completion for --server-oid flag
func (cmd *CMD) registerServerOIDCompletion() {
	completionFunc := func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		companyOID := cmd.getCompanyOIDForCompletion(c)
		if companyOID == "" {
			return nil, cobra.ShellCompDirectiveError
		}

		// Fetch servers from API for the specific company
		servers, apiReturn, err := cmd.runMiddleware.API.ServerList(companyOID)
		if err != nil || apiReturn != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, server := range servers {
			state := "unknown"
			if server.State != nil {
				state = *server.State
			}
			desc := fmt.Sprintf("%s\t%s (%s)", server.OID, server.Name, state)
			completions = append(completions, desc)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	// Skip commands that have contextual server-oid completions
	networkDetachCmd := findCommand(cmd.RootCommand, "network", "detach")
	registerCompletionRecursive(cmd.RootCommand, "server-oid", completionFunc, networkDetachCmd)
}

// registerNetworkOIDCompletion registers completion for --network-oid flag
func (cmd *CMD) registerNetworkOIDCompletion() {
	completionFunc := func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		companyOID := cmd.getCompanyOIDForCompletion(c)
		if companyOID == "" {
			return nil, cobra.ShellCompDirectiveError
		}

		// Fetch networks from API for the specific company
		networkList, err := cmd.runMiddleware.API.GetNetworkList(companyOID)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, network := range networkList.Networks {
			// Format: OID with description (name and speed)
			speedInfo := fmt.Sprintf("%d%s", network.Speed.Value, network.Speed.Unit)
			desc := fmt.Sprintf("%s\t%s (%s)", network.OID, network.Name, speedInfo)
			completions = append(completions, desc)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	registerCompletionRecursive(cmd.RootCommand, "network-oid", completionFunc)
}

// registerSubscriptionOIDCompletion registers completion for --subscription-oid flag
func (cmd *CMD) registerSubscriptionOIDCompletion() {
	completionFunc := func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		companyOID := cmd.getCompanyOIDForCompletion(c)
		if companyOID == "" {
			return nil, cobra.ShellCompDirectiveError
		}

		// Fetch subscriptions from API for the specific company (active only)
		subscriptions, err := cmd.runMiddleware.API.GetSubscriptionList(companyOID, true)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, sub := range subscriptions {
			// Format: OID with description (name/document number and state)
			name := sub.Name
			if name == "" {
				name = sub.DocumentNumber
			}
			desc := fmt.Sprintf("%s\t%s (%s)", sub.OID, name, sub.State)
			completions = append(completions, desc)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	registerCompletionRecursive(cmd.RootCommand, "subscription-oid", completionFunc)
}

// =============================================================================
// GLOBAL/SHARED COMPLETIONS
// These resources are global or available to all users
// =============================================================================

// registerTemplateOIDCompletion registers completion for --template-oid flag
func (cmd *CMD) registerTemplateOIDCompletion() {
	completionFunc := func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Fetch templates from API (grouped by OS)
		templateGroups, _, err := cmd.runMiddleware.API.ListTemplates()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, group := range templateGroups {
			for _, tmpl := range group.Versions {
				// Format: OID with description (OS and version)
				desc := fmt.Sprintf("%s\t%s %s", tmpl.OID, group.OS, tmpl.Version)
				completions = append(completions, desc)
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	registerCompletionRecursive(cmd.RootCommand, "template-oid", completionFunc)
}

// registerCompanyOIDCompletion registers completion for --company-oid flag
func (cmd *CMD) registerCompanyOIDCompletion() {
	completionFunc := func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Fetch companies from API (user's accessible companies)
		companies, err := cmd.runMiddleware.API.GetListOfCompanies()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, company := range companies {
			// Format: OID with description (name)
			desc := fmt.Sprintf("%s\t%s", company.OID, company.Name)
			completions = append(completions, desc)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	registerCompletionRecursive(cmd.RootCommand, "company-oid", completionFunc)
}

// =============================================================================
// USER-SCOPED COMPLETIONS
// These resources belong to the current user
// =============================================================================

// registerSSHKeyOIDCompletion registers completion for --ssh-key-oid flag
func (cmd *CMD) registerSSHKeyOIDCompletion() {
	completionFunc := func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Get current user to fetch their SSH keys
		user, err := cmd.runMiddleware.API.GetUserInfos()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		sshKeys, err := cmd.runMiddleware.API.GetSSHKeyList(user.OID)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, key := range sshKeys {
			// Format: OID with description (name)
			desc := fmt.Sprintf("%s\t%s", key.OID, key.Name)
			completions = append(completions, desc)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	registerCompletionRecursive(cmd.RootCommand, "ssh-key-oid", completionFunc)
}

// registerTokenOIDCompletion registers completion for --token-oid flag
func (cmd *CMD) registerTokenOIDCompletion() {
	completionFunc := func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Fetch API tokens for the current user
		tokens, err := cmd.runMiddleware.API.ListAPITokens()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, token := range tokens {
			// Format: OID with description (name)
			desc := fmt.Sprintf("%s\t%s", token.OID, token.Name)
			completions = append(completions, desc)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	registerCompletionRecursive(cmd.RootCommand, "token-oid", completionFunc)
}

// =============================================================================
// DEPENDENT COMPLETIONS
// These resources depend on another flag being set first
// =============================================================================

// registerSnapshotOIDCompletion registers completion for --snapshot-oid flag
// Requires --server-oid to be specified first
func (cmd *CMD) registerSnapshotOIDCompletion() {
	completionFunc := func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Get server-oid from the command flags - snapshots belong to a server
		serverOID, _ := c.Flags().GetString("server-oid")

		if serverOID == "" {
			// No server specified, can't list snapshots
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}

		// Fetch snapshots for the specific server
		snapshots, _, err := cmd.runMiddleware.API.ListSnapshots(serverOID)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		var completions []string
		for _, snap := range snapshots {
			// Format: OID with description (name)
			desc := fmt.Sprintf("%s\t%s", snap.OID, snap.Name)
			completions = append(completions, desc)
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	registerCompletionRecursive(cmd.RootCommand, "snapshot-oid", completionFunc)
}

// registerISOOIDCompletion registers completion for --iso-oid flag
// Requires --server-oid to be specified first (ISOs are mounted on a server)
func (cmd *CMD) registerISOOIDCompletion() {
	completionFunc := func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Get server-oid from the command flags - ISOs belong to a server
		serverOID, _ := c.Flags().GetString("server-oid")

		if serverOID == "" {
			// No server specified, can't list ISOs
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		}

		// Fetch server detail to get mounted ISOs
		server, _, err := cmd.runMiddleware.API.GetServerOID(serverOID)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		// Return the mounted ISO OIDs (empty if none)
		return server.ISOsOID, cobra.ShellCompDirectiveNoFileComp
	}

	registerCompletionRecursive(cmd.RootCommand, "iso-oid", completionFunc)
}

// =============================================================================
// CONTEXTUAL COMPLETIONS
// These completions change behavior based on command context
// =============================================================================

// registerIPCompletion registers completion for --ip flag
// For "ip detach": shows IPs attached to the specified server (or all attached IPs if no server specified)
// For "ip attach": shows available (unattached) IPs from the company
func (cmd *CMD) registerIPCompletion() {
	// Register for ip detach command - show IPs attached to the server
	if ipDetachCmd := findCommand(cmd.RootCommand, "ip", "detach"); ipDetachCmd != nil {
		_ = ipDetachCmd.RegisterFlagCompletionFunc("ip", func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			serverOID, _ := c.Flags().GetString("server-oid")

			// If server is specified, show only IPs attached to that server
			if serverOID != "" {
				server, _, err := cmd.runMiddleware.API.GetServerOID(serverOID)
				if err == nil {
					ipList, err := cmd.runMiddleware.API.GetCompanyIPList(server.Company)
					if err == nil {
						var completions []string
						for _, ip := range ipList {
							if ip.ServerOID == serverOID {
								desc := fmt.Sprintf("%s\t%s", ip.Address, ip.Reverse)
								completions = append(completions, desc)
							}
						}
						return completions, cobra.ShellCompDirectiveNoFileComp
					}
				}
			}

			// Fallback: show all attached IPs from the default company
			companyOID := cmd.getCompanyOIDForCompletion(c)
			if companyOID == "" {
				return nil, cobra.ShellCompDirectiveError
			}
			ipList, err := cmd.runMiddleware.API.GetCompanyIPList(companyOID)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
			var completions []string
			for _, ip := range ipList {
				if ip.ServerOID != "" {
					desc := fmt.Sprintf("%s\t%s (%s)", ip.Address, ip.Reverse, ip.ServerName)
					completions = append(completions, desc)
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		})
	}

	// Register for ip attach command - show available (unattached) IPs
	if ipAttachCmd := findCommand(cmd.RootCommand, "ip", "attach"); ipAttachCmd != nil {
		_ = ipAttachCmd.RegisterFlagCompletionFunc("ip", func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var companyOID string

			// If server is specified, get company from server
			serverOID, _ := c.Flags().GetString("server-oid")
			if serverOID != "" {
				server, _, err := cmd.runMiddleware.API.GetServerOID(serverOID)
				if err == nil {
					companyOID = server.Company
				}
			}

			// Fallback to default company
			if companyOID == "" {
				companyOID = cmd.getCompanyOIDForCompletion(c)
			}
			if companyOID == "" {
				return nil, cobra.ShellCompDirectiveError
			}

			ipList, err := cmd.runMiddleware.API.GetCompanyIPList(companyOID)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			// Filter to only available (unattached) IPs
			var completions []string
			for _, ip := range ipList {
				if ip.ServerOID == "" {
					desc := fmt.Sprintf("%s\t%s", ip.Address, ip.Reverse)
					completions = append(completions, desc)
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		})
	}

	// Register for ip reverse command - show all company IPs
	if ipReverseCmd := findCommand(cmd.RootCommand, "ip", "reverse"); ipReverseCmd != nil {
		_ = ipReverseCmd.RegisterFlagCompletionFunc("ip", func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			companyOID := cmd.getCompanyOIDForCompletion(c)
			if companyOID == "" {
				return nil, cobra.ShellCompDirectiveError
			}

			ipList, err := cmd.runMiddleware.API.GetCompanyIPList(companyOID)
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			var completions []string
			for _, ip := range ipList {
				serverInfo := "available"
				if ip.ServerName != "" {
					serverInfo = ip.ServerName
				}
				desc := fmt.Sprintf("%s\t%s (%s)", ip.Address, ip.Reverse, serverInfo)
				completions = append(completions, desc)
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		})
	}
}

// registerNetworkServerOIDCompletion registers contextual --server-oid completion for network commands
// For "network detach": shows only servers attached to the specified network (if --network-oid is set)
func (cmd *CMD) registerNetworkServerOIDCompletion() {
	// Register for network detach command - show only servers attached to the network
	if netDetachCmd := findCommand(cmd.RootCommand, "network", "detach"); netDetachCmd != nil {
		_ = netDetachCmd.RegisterFlagCompletionFunc("server-oid", func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			networkOID, _ := c.Flags().GetString("network-oid")

			// If network-oid is specified, show only servers attached to that network
			if networkOID != "" {
				network, err := cmd.runMiddleware.API.GetNetworkDetail(networkOID)
				if err == nil && len(network.Interfaces) > 0 {
					var completions []string
					for _, iface := range network.Interfaces {
						state := "unknown"
						if iface.Server.State != nil {
							state = *iface.Server.State
						}
						desc := fmt.Sprintf("%s\t%s (%s)", iface.Server.OID, iface.Server.Name, state)
						completions = append(completions, desc)
					}
					return completions, cobra.ShellCompDirectiveNoFileComp
				}
			}

			// Fallback: show all servers from the company
			companyOID := cmd.getCompanyOIDForCompletion(c)
			if companyOID == "" {
				return nil, cobra.ShellCompDirectiveError
			}

			servers, apiReturn, err := cmd.runMiddleware.API.ServerList(companyOID)
			if err != nil || apiReturn != nil {
				return nil, cobra.ShellCompDirectiveError
			}

			var completions []string
			for _, server := range servers {
				state := "unknown"
				if server.State != nil {
					state = *server.State
				}
				desc := fmt.Sprintf("%s\t%s (%s)", server.OID, server.Name, state)
				completions = append(completions, desc)
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		})
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// findCommand finds a command by walking the path (e.g., "ip", "detach")
func findCommand(root *cobra.Command, path ...string) *cobra.Command {
	current := root
	for _, name := range path {
		found := false
		for _, sub := range current.Commands() {
			if sub.Name() == name {
				current = sub
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return current
}

// registerCompletionRecursive walks the command tree and registers completion
// for any command that has the specified flag, optionally skipping specified commands
func registerCompletionRecursive(c *cobra.Command, flagName string, completionFunc func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective), skip ...*cobra.Command) {
	// Check if this command should be skipped
	for _, skipCmd := range skip {
		if c == skipCmd {
			// Still recurse into subcommands even if skipping this one
			for _, subCmd := range c.Commands() {
				registerCompletionRecursive(subCmd, flagName, completionFunc, skip...)
			}
			return
		}
	}

	// Check if this command has the flag
	if flag := c.Flags().Lookup(flagName); flag != nil {
		_ = c.RegisterFlagCompletionFunc(flagName, completionFunc)
	}

	// Recurse into subcommands
	for _, subCmd := range c.Commands() {
		registerCompletionRecursive(subCmd, flagName, completionFunc, skip...)
	}
}
