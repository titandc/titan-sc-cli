package cmd

import (
	"github.com/spf13/cobra"
)

// APITokenCmdAdd adds the api-token command and subcommands
func (cmd *CMD) APITokenCmdAdd() {
	apiToken := &cobra.Command{
		Use:     "api-token",
		Aliases: []string{"token"},
		Short:   "Manage API tokens.",
		Long:    "Manage API tokens for authentication.",
		GroupID: "resources",
	}

	apiTokenList := &cobra.Command{
		Use:   "list",
		Short: "List all API tokens.",
		Long:  "List all API tokens for the authenticated user.",
		Run:   cmd.runMiddleware.APITokenList,
	}

	apiTokenShow := &cobra.Command{
		Use:   "show --token-oid TOKEN_OID",
		Short: "Show API token details.",
		Long:  "Show details of a specific API token.",
		Run:   cmd.runMiddleware.APITokenShow,
	}

	apiTokenCreate := &cobra.Command{
		Use:   "create --name NAME",
		Short: "Create a new API token.",
		Long: `Create a new API token.

The token name must be alphanumeric with hyphens and underscores only.
Use --expire-days to set expiration relative to now (e.g., --expire-days 30).

If no expiration is set, the token will never expire.`,
		Run: cmd.runMiddleware.APITokenCreate,
	}

	apiTokenUpdate := &cobra.Command{
		Use:   "update --token-oid TOKEN_OID",
		Short: "Update an API token.",
		Long: `Update an existing API token.

You can update the name and/or expiration date.
Use --no-expire to remove expiration (make token permanent).
WARNING: Updating the expiration will regenerate the token value!`,
		Run: cmd.runMiddleware.APITokenUpdate,
	}

	apiTokenDelete := &cobra.Command{
		Use:   "delete --token-oid TOKEN_OID",
		Short: "Delete an API token.",
		Long:  "Delete an API token by its OID.",
		Run:   cmd.runMiddleware.APITokenDelete,
	}

	cmd.RootCommand.AddCommand(apiToken)
	apiToken.AddCommand(apiTokenList, apiTokenShow, apiTokenCreate, apiTokenUpdate, apiTokenDelete)

	// Show flags
	apiTokenShow.Flags().StringP("token-oid", "", "", "API token OID.")
	_ = apiTokenShow.MarkFlagRequired("token-oid")

	// Create flags
	apiTokenCreate.Flags().StringP("name", "n", "", "Token name (alphanumeric, hyphens, underscores).")
	apiTokenCreate.Flags().IntP("expire-days", "", 0, "Expiration in days from now (0 = never expires).")
	_ = apiTokenCreate.MarkFlagRequired("name")

	// Update flags
	apiTokenUpdate.Flags().StringP("token-oid", "", "", "API token OID.")
	apiTokenUpdate.Flags().StringP("name", "n", "", "New token name.")
	apiTokenUpdate.Flags().IntP("expire-days", "", 0, "New expiration in days from now (0 = never expires, regenerates token!).")
	apiTokenUpdate.Flags().BoolP("no-expire", "", false, "Remove expiration, make token permanent (regenerates token!).")
	_ = apiTokenUpdate.MarkFlagRequired("token-oid")

	// Delete flags
	apiTokenDelete.Flags().StringP("token-oid", "", "", "API token OID.")
	_ = apiTokenDelete.MarkFlagRequired("token-oid")
}
