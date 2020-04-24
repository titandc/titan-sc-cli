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
    VersionMajor = 1
    VersionMinor = 0
    VersionPatch = 0
    LocalApp = "./titan-sc"
    UsrLocalBinApp = "/usr/local/bin/titan-sc"
)

var weatherMap = &cobra.Command{
    Use: "weathermap",
    Short: "Show weather map.",
    Long: "Show weathermap.",
    Run: API.WeatherMap,
}

var versionCmd = &cobra.Command{
    Use: "version",
    Short: "Print version and exit.",
    Long: "Print current CLI version and exit.",
    Run: version,
}

var AddTokenCmd = &cobra.Command{
    Use: "setup \"token-string\"",
    Short: "Automated config/install.",
    Long: "Automated config/install.",
    Args: cmdNeed1Args,
    Run: InitApp,
}

func utilsCmdAdd() {
    rootCmd.AddCommand(weatherMap, versionCmd, AddTokenCmd)
}

func version(cmd *cobra.Command, args []string) {
    _ = cmd
    _ = args

    fmt.Printf("Titan cloud CLI version %d.%d.%d\n",
        VersionMajor, VersionMinor, VersionPatch)
    os.Exit(0)
}

func InitApp(cmd *cobra.Command, args []string) {
    _ = cmd
    var err error

    if runtime.GOOS == "windows" {
        path := ConfigFileName
        err = InitCreateFile(path, args[0])
    } else {
        err = InitAppUnixLike(args[0])
    }
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println("Init success.")
    }
}

func InitAppUnixLike(token string) error {
    // install bin in /usr/local/bin
    _, err := os.Stat(UsrLocalBinApp)
    if !os.IsNotExist(err) {
        if err := os.Remove(UsrLocalBinApp); err != nil {
            return err
        }
    }
    if err := os.Link(LocalApp, UsrLocalBinApp); err != nil {
        return err
    }

    path := os.Getenv("HOME") + "/" + ".titan"
    if err := os.Mkdir(path, os.ModePerm); err != nil {
        if !strings.Contains(err.Error(), ": file exists") {
            return fmt.Errorf("Init: create dir path %s error.", path)
        }
    }
    path = path + "/" + ConfigFileName
    return InitCreateFile(path, token)
}

func InitCreateFile(path, token string) error {
    data := map[string]interface{}{
        "token": token,
    }
    viper.Set("default", data)
    viper.SetConfigType("toml")
    viper.SetConfigName("config")
    err := viper.WriteConfigAs(path)
    return err
}
