package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"runtime"
	"strings"
	. "titan-sc/api"
)

const (
	VersionMajor   = 1
	VersionMinor   = 2
	VersionPatch   = 0
	LocalApp       = "./titan-sc"
	UsrLocalBinApp = "/usr/local/bin/titan-sc"
)

var weatherMap = &cobra.Command{
	Use:   "weathermap",
	Short: "Show weather map.",
	Long:  "Show weathermap.",
	Run:   API.WeatherMap,
}

var addTokenCmd = &cobra.Command{
	Use:   "setup \"token-string\"",
	Short: "Automated config/install.",
	Long:  "Automated config/install.",
	Args:  cmdNeed1Args,
	Run:   InitApp,
}

var managedServices = &cobra.Command{
	Use:   "managed-services COMPANY_UUID",
	Short: "Enable managed services.",
	Long:  "Enable managed services.",
	Args:  cmdNeed1UUID,
	Run:   API.ManagedServices,
}

func utilsCmdAdd() {
	rootCmd.AddCommand(weatherMap, addTokenCmd, managedServices)
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
	// check user UID
	usr, err := user.Current()
	if err != nil {
		return err
	}
	if usr.Uid != "0" {
		return fmt.Errorf("You must have root privileges to use the 'setup' command.")
	}

	// install bin in /usr/local/bin
	_, err = os.Stat(UsrLocalBinApp)
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
			return fmt.Errorf("Init: create dir path <%s> error.", path)
		}
	}
	path = path + "/" + ConfigFileName
	return InitCreateFile(path, token)
}

func InitCreateFile(path, token string) error {
	data := map[string]interface{}{
		"token": token,
		"uri":   DefaultURI,
	}
	viper.Set("default", data)
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	err := viper.WriteConfigAs(path)
	return err
}
