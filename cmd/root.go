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
    ConfigFileName = "config"
)

var rootCmd = &cobra.Command{
    Use:              GetProgramName(),
    Short:            "Titan SC CLI",
    PersistentPreRun: GetAPIToken,
    Long:             "Titan Small Cloud command line interface.",
}

func init() {
    companyCmdAdd()
    serverCmdAdd()
    snapshotCmdAdd()
    historyCmdAdd()
    utilsCmdAdd()
    networkCmdAdd()
    ipkvmCmdAdd()
    rootCmd.PersistentFlags().BoolP("human", "H", false, "Format output for human.")
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

/* juste get token API */
func GetAPIToken(cmd *cobra.Command, args []string) {

    arrCmd := strings.SplitN(cmd.CommandPath(), " ", 3)
    if len(arrCmd) > 1 && (arrCmd[1] == "version" || arrCmd[1] == "setup") {
        return
    }

    GetApiTokenByEnv()
    if API.Token == "" {
        if err := GetApiTokenByFile(); err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
    }

    if API.Token == "" {
        fmt.Println("Token api missing")
        os.Exit(1)
    }
}

func GetApiTokenByEnv() {
    API.Token = os.Getenv(EnvApiToken)
}

func GetApiTokenByFile() error {
    var path string

    if runtime.GOOS == "windows" {
        path = "./"
    } else {
        path = os.Getenv("HOME") + "/.titan"
    }
    viper.SetConfigType("toml")
    viper.SetConfigName(ConfigFileName)
    viper.AddConfigPath(path)
    if err := viper.ReadInConfig(); err != nil {
        return fmt.Errorf("Error to read configuration file: %s", err.Error())
    }
    API.Token = viper.GetString("default.token")
    return nil
}
