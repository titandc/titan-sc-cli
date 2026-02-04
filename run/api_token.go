package run

import (
	"fmt"
	"time"

	"titan-sc/api"

	"github.com/spf13/cobra"
)

func (run *RunMiddleware) APITokenList(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)

	tokens, err := run.API.ListAPITokens()
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(tokens)
		return
	}

	if len(tokens) == 0 {
		fmt.Println("No API tokens found.")
		return
	}

	table := NewTable("NAME", "EXPIRES", "OID")
	table.SetNoColor(!run.Color)

	for _, token := range tokens {
		var expires string
		var expiresColorFn func(string) string

		if token.Expire != nil {
			t := time.Unix(*token.Expire, 0)
			if t.Before(time.Now()) {
				expires = t.Format("2006-01-02 15:04") + " (expired)"
				if run.Color {
					expiresColorFn = ColorFn("red")
				}
			} else {
				expires = t.Format("2006-01-02 15:04")
			}
		} else {
			expires = "never"
			if run.Color {
				expiresColorFn = ColorFn("green")
			}
		}

		table.AddRow(
			ColName(token.Name),
			ColColor(expires, expiresColorFn),
			ColOID(token.OID),
		)
	}
	table.Print()
}

func (run *RunMiddleware) APITokenShow(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	tokenOID, _ := cmd.Flags().GetString("token-oid")

	token, err := run.API.GetAPIToken(tokenOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(token)
		return
	}

	fmt.Printf("%s %s\n", run.Colorize("Name:", "cyan"), token.Name)
	fmt.Printf("%s %s\n", run.Colorize("OID:", "cyan"), token.OID)

	if token.Expire != nil {
		t := time.Unix(*token.Expire, 0)
		expireStr := t.Format("2006-01-02 15:04:05")
		if t.Before(time.Now()) {
			fmt.Printf("%s %s %s\n", run.Colorize("Expires:", "cyan"), run.Colorize(expireStr, "red"), "(expired)")
		} else {
			fmt.Printf("%s %s\n", run.Colorize("Expires:", "cyan"), expireStr)
		}
	} else {
		fmt.Printf("%s %s\n", run.Colorize("Expires:", "cyan"), run.Colorize("never", "green"))
	}

	fmt.Printf("%s %s\n", run.Colorize("Value:", "cyan"), token.Value)
}

func (run *RunMiddleware) APITokenCreate(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	name, _ := cmd.Flags().GetString("name")
	expireDays, _ := cmd.Flags().GetInt("expire-days")

	create := &api.APITokenCreate{
		Name: name,
	}

	// Handle expiration
	if expireDays > 0 {
		expireTime := time.Now().AddDate(0, 0, expireDays).Unix()
		create.Expire = &expireTime
	}

	token, err := run.API.CreateAPIToken(create)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(token)
		return
	}

	fmt.Printf("%s API token created successfully\n", run.Colorize("✓", "green"))
	fmt.Printf("%s %s\n", run.Colorize("Name:", "cyan"), token.Name)
	fmt.Printf("%s %s\n", run.Colorize("OID:", "cyan"), token.OID)

	if token.Expire != nil {
		t := time.Unix(*token.Expire, 0)
		fmt.Printf("%s %s\n", run.Colorize("Expires:", "cyan"), t.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("%s %s\n", run.Colorize("Expires:", "cyan"), run.Colorize("never", "green"))
	}

	fmt.Printf("%s %s\n", run.Colorize("Token:", "cyan"), token.Value)
}

func (run *RunMiddleware) APITokenUpdate(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	tokenOID, _ := cmd.Flags().GetString("token-oid")
	name, _ := cmd.Flags().GetString("name")
	expireDays, _ := cmd.Flags().GetInt("expire-days")
	noExpire, _ := cmd.Flags().GetBool("no-expire")

	update := &api.APITokenUpdate{}

	if name != "" {
		update.Name = name
	}

	// Handle expiration
	if noExpire {
		// Set expire to 0 to remove expiration (API interprets 0 as no expiration)
		var zero int64 = 0
		update.Expire = &zero
	} else if expireDays > 0 {
		expireTime := time.Now().AddDate(0, 0, expireDays).Unix()
		update.Expire = &expireTime
	}

	if name == "" && update.Expire == nil {
		fmt.Println("Nothing to update. Use --name, --expire-days, or --no-expire.")
		return
	}

	// Warn if updating expiration
	if update.Expire != nil {
		fmt.Printf("%s Updating expiration will regenerate the token value!\n", run.Colorize("⚠ WARNING:", "yellow"))
	}

	token, err := run.API.UpdateAPIToken(tokenOID, update)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(token)
		return
	}

	fmt.Printf("%s API token updated successfully\n", run.Colorize("✓", "green"))
	fmt.Printf("%s %s\n", run.Colorize("Name:", "cyan"), token.Name)
	fmt.Printf("%s %s\n", run.Colorize("OID:", "cyan"), token.OID)

	if token.Expire != nil {
		t := time.Unix(*token.Expire, 0)
		fmt.Printf("%s %s\n", run.Colorize("Expires:", "cyan"), t.Format("2006-01-02 15:04:05"))
	} else {
		fmt.Printf("%s %s\n", run.Colorize("Expires:", "cyan"), run.Colorize("never", "green"))
	}

	// Show new token value if expiration was updated
	if update.Expire != nil {
		fmt.Printf("%s %s\n", run.Colorize("Token:", "cyan"), token.Value)
	}
}

func (run *RunMiddleware) APITokenDelete(cmd *cobra.Command, args []string) {
	_ = args
	run.ParseGlobalFlags(cmd)
	tokenOID, _ := cmd.Flags().GetString("token-oid")

	err := run.API.DeleteAPIToken(tokenOID)
	if err != nil {
		run.OutputError(err)
		return
	}

	if run.JSONOutput {
		printAsJson(map[string]string{
			"status":  "success",
			"oid":     tokenOID,
			"message": "API token deleted successfully",
		})
	} else {
		fmt.Printf("%s API token %s deleted successfully\n", run.Colorize("✓", "green"), tokenOID)
	}
}
