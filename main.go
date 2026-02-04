package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"titan-sc/api"
	"titan-sc/cmd"
	"titan-sc/run"

	"github.com/spf13/viper"
)

const (
	EnvApiToken    = "TITAN_API_TOKEN"
	EnvApiUri      = "TITAN_URI"
	ConfigFileName = "config"
	VersionMajor   = 4
	VersionMinor   = 0
	VersionPatch   = 0
)

var (
	cmdInstance *cmd.CMD
	runInstance *run.RunMiddleware
	apiInstance *api.API
)

func init() {
	loadConfigurationFile()
	// Retrieve optional API URI or use default one
	uri := getApiUriFromEnv()
	if uri == "" {
		uri = getApiUriFromFile()
		if uri == "" {
			//run.OutputError(fmt.Errorf("Unable to retrieve URI from configuration file."))
			uri = api.DefaultURI
		}
	}

	// Retrieve mandatory API token
	tokenDefined := true
	token := getApiTokenFromEnv()
	if token == "" {
		token = getApiTokenFromFile()
		if token == "" {
			tokenDefined = false
		}
	}

	operatingsystem := runtime.GOOS

	apiInstance = api.NewAPI(token, uri, operatingsystem, fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch))
	runInstance = run.NewRunMiddleware(apiInstance)
	cmdInstance = cmd.NewCMD(getProgramName(), ConfigFileName, tokenDefined, runInstance, VersionMajor, VersionMinor,
		VersionPatch)

	cmdInstance.CompletionCmdAdd()
	cmdInstance.CompanyCmdAdd()
	cmdInstance.ServerCmdAdd()
	cmdInstance.SnapshotCmdAdd()
	cmdInstance.HistoryCmdAdd()
	cmdInstance.NetworkCmdAdd()
	cmdInstance.KvmIpCmdAdd()
	cmdInstance.IpCmdAdd()
	cmdInstance.UserCmdAdd()
	cmdInstance.SetupCmdAdd()
	cmdInstance.SSHKeysCmdAdd()
	cmdInstance.SubscriptionCmdAdd()
	cmdInstance.VersionCmdAdd()
	cmdInstance.RootCommand.PersistentFlags().BoolP("json", "j", false,
		"Output in JSON format (disables colors).")
	cmdInstance.RootCommand.PersistentFlags().Bool("no-color", false,
		"Disable colorized output.")

}

func main() {
	cmdInstance.Execute()
}

func getApiTokenFromEnv() string {
	return os.Getenv(EnvApiToken)
}

func getApiTokenFromFile() string {
	return viper.GetString("default.token")
}

func getApiUriFromFile() string {
	return viper.GetString("default.uri")
}

func getApiUriFromEnv() string {
	return os.Getenv(EnvApiUri)
}

func getProgramName() string {
	var progname string
	idx := strings.LastIndex(os.Args[0], "/")
	if idx != 0 && len(os.Args[0][idx+1:]) > 0 {
		progname = os.Args[0][idx+1:]
	} else {
		progname = os.Args[0]
	}
	return progname
}

func loadConfigurationFile() {
	viper.SetConfigType("toml")
	viper.SetConfigName(ConfigFileName)

	// Add config paths in order of priority (first found wins)
	if runtime.GOOS == "windows" {
		// Windows: check APPDATA first, then current directory
		if appData := os.Getenv("APPDATA"); appData != "" {
			viper.AddConfigPath(appData + "\\titan")
		}
		viper.AddConfigPath(".")
	} else {
		// Unix/macOS: check ~/.titan first, then current directory
		if home := os.Getenv("HOME"); home != "" {
			viper.AddConfigPath(home + "/.titan")
		}
		viper.AddConfigPath(".")
	}

	_ = viper.ReadInConfig()
}
