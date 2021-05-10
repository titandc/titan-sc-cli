package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"runtime"
	"strings"
	"titan-sc/api"
	"titan-sc/cmd"
	"titan-sc/run"
)

const (
	EnvApiToken    = "TITAN_API_TOKEN"
	EnvApiUri      = "TITAN_URI"
	ConfigFileName = "config"
	VersionMajor   = 3
	VersionMinor   = 0
	VersionPatch   = 1
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
	cmdInstance.FirewallCmdAdd()
	cmdInstance.KvmIpCmdAdd()
	cmdInstance.IpCmdAdd()
	cmdInstance.UserCmdAdd()
	cmdInstance.SetupCmdAdd()
	cmdInstance.SSHKeysCmdAdd()
	cmdInstance.VersionCmdAdd()
	cmdInstance.PnatCmdAdd()
	cmdInstance.ManagedServicesCmdAdd()
	cmdInstance.WeatherMapCmdAdd()
	cmdInstance.RootCommand.PersistentFlags().BoolP("human", "H", false,
		"Format output for human.")
	cmdInstance.RootCommand.PersistentFlags().BoolP("color", "C", false,
		"Enable colorized output.")

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
	var path string
	if runtime.GOOS == "windows" {
		path = "./"
	} else {
		path = os.Getenv("HOME") + "/.titan"
	}
	viper.SetConfigType("toml")
	viper.SetConfigName(ConfigFileName)
	viper.AddConfigPath(path)
	_ = viper.ReadInConfig()
}
