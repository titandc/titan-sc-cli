package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "os"
    "runtime"
    "strings"
    . "titan-sc/api"
)

const (
    EnvApiToken    = "TITAN_API_TOKEN"
    EnvApiUri      = "TITAN_URI"
    ConfigFileName = "config"
)

var rootCmd = &cobra.Command{
    Use:              GetProgramName(),
    Short:            "Titan SC CLI",
    PersistentPreRun: GetAPIToken,
    Long:             "Titan Small Cloud - Command Line Interface",
}

func init() {
    companyCmdAdd()
    serverCmdAdd()
    snapshotCmdAdd()
    historyCmdAdd()
    utilsCmdAdd()
    networkCmdAdd()
    kvmIpCmdAdd()
    ipCmdAdd()
    versionCmdAdd()
    rootCmd.PersistentFlags().BoolP("human", "H", false, "Format output for human.")
    rootCmd.PersistentFlags().BoolP("color", "C", false, "Enable colorized output.")
    loadConfigurationFile()
}

func Execute() {
    API.CLIos = runtime.GOOS
    API.CLIVersion = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func GetProgramName() string {
    var progname string

    idx := strings.LastIndex(os.Args[0], "/")
    if idx != 0 && len(os.Args[0][idx+1:]) > 0 {
        progname = os.Args[0][idx+1:]
    } else {
        progname = os.Args[0]
    }
    return progname
}

func GetAPIToken(cmd *cobra.Command, args []string) {

    arrCmd := strings.SplitN(cmd.CommandPath(), " ", 3)
    if len(arrCmd) > 1 && arrCmd[1] == "setup" {
        return
    }
    if len(arrCmd) > 2 && arrCmd[1] == "version" && arrCmd[2] == "cli" {
        return
    }
    if len(arrCmd) == 2 && arrCmd[1] == "version" {
        return
    }

    // Retrieve mandatory API token
    GetApiTokenFromEnv()
    if API.Token == "" {
        if err := GetApiTokenFromFile(); err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
    }

    // Retrieve optional API URI
    GetApiUriFromEnv()
    if API.URI == "" {
        if err := GetApiUriFromFile(); err != nil {
            API.URI = DefaultURI
        }
    }
}

func GetApiTokenFromEnv() {
    API.Token = os.Getenv(EnvApiToken)
}

func GetApiTokenFromFile() error {
    token := viper.GetString("default.token")
    if token == "" {
        return fmt.Errorf("Unable to retrieve token from configuration file.")
    }
    API.Token = token
    return nil
}

func GetApiUriFromEnv() {
    API.URI = os.Getenv(EnvApiUri)
}

func GetApiUriFromFile() error {
    uri := viper.GetString("default.uri")
    if uri == "" {
        return fmt.Errorf("Unable to retrieve token from configuration file.")
    }
    API.URI = uri
    return nil
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
