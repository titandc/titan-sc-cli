package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"titan-sc/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cmd *CMD) SetupCmdAdd() {
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Configure CLI with your API credentials.",
		Long: `Configure the CLI with your API token and optional custom API endpoint.

The token will be validated before saving.

Configuration is stored in:
  - Linux/macOS: ~/.titan/config
  - Windows:     %APPDATA%\titan\config

The CLI also checks the current directory for a config or config.toml file as a fallback.

Examples:
  titan-sc setup --token "your-api-token"
  titan-sc setup --token "your-api-token" --uri "https://custom-api.example.com/v2"`,
		Run:     cmd.setupApp,
		GroupID: "config",
	}

	cmd.RootCommand.AddCommand(setupCmd)

	setupCmd.Flags().StringP("token", "t", "", "API authentication token.")
	setupCmd.Flags().String("uri", api.DefaultURI, "Custom API endpoint URL.")
	_ = setupCmd.MarkFlagRequired("token")
}

func (cmd *CMD) setupApp(cobraCommand *cobra.Command, args []string) {
	_ = args
	cmd.runMiddleware.ParseGlobalFlags(cobraCommand)

	token, _ := cobraCommand.Flags().GetString("token")
	uri, _ := cobraCommand.Flags().GetString("uri")

	// Validate token
	fmt.Print("Validating token... ")
	if err := validateToken(token, uri); err != nil {
		fmt.Println(cmd.runMiddleware.Colorize("FAILED", "red"))
		fmt.Printf("Error: %s\n", err.Error())
		fmt.Println("\nPlease check your token and try again.")
		fmt.Println("You can generate a new token from the Titan dashboard.")
		return
	}
	fmt.Println(cmd.runMiddleware.Colorize("OK", "green"))

	// Save configuration
	fmt.Print("Saving configuration... ")
	if err := saveConfig(token, uri); err != nil {
		fmt.Println(cmd.runMiddleware.Colorize("FAILED", "red"))
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	fmt.Println(cmd.runMiddleware.Colorize("OK", "green"))

	// Get config path for display
	configPath := getConfigPath()
	fmt.Printf("\nConfiguration saved to: %s\n", cmd.runMiddleware.Colorize(configPath, "cyan"))
	fmt.Println(cmd.runMiddleware.Colorize("\nSetup complete!", "green"))
}

// validateToken checks if the token is valid by making an API call
func validateToken(token, uri string) error {
	tempAPI := api.NewAPI(token, uri, runtime.GOOS, "setup")
	_, err := tempAPI.GetUserInfos()
	if err != nil {
		return fmt.Errorf("invalid token or unable to connect to API")
	}
	return nil
}

// getConfigDir returns the directory where config should be stored
func getConfigDir() string {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return appData + "\\titan"
		}
		return "."
	}
	if home := os.Getenv("HOME"); home != "" {
		return home + "/.titan"
	}
	return "."
}

// getConfigPath returns the full path to the config file
func getConfigPath() string {
	dir := getConfigDir()
	if runtime.GOOS == "windows" {
		return dir + "\\config"
	}
	return dir + "/config"
}

// saveConfig writes the token and URI to the config file
func saveConfig(token, uri string) error {
	configDir := getConfigDir()
	configPath := getConfigPath()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Build config data
	data := map[string]interface{}{
		"token": token,
	}
	// Only save URI if it's not the default
	if uri != "" && uri != api.DefaultURI {
		data["uri"] = uri
	}

	viper.Set("default", data)
	viper.SetConfigType("toml")

	// Write config file
	if err := viper.WriteConfigAs(configPath); err != nil {
		// If file exists, try to overwrite
		if strings.Contains(err.Error(), "already exists") {
			return viper.WriteConfigAs(configPath)
		}
		return err
	}

	return nil
}
