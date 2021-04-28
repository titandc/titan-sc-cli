package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"runtime"
	"strings"
	"titan-sc/api"
)

const (
	LocalApp       = "./titan-sc"
	UsrLocalBinApp = "/usr/local/bin/titan-sc"
)

func (cmd *CMD) SetupCmdAdd() {

	addTokenCmd := &cobra.Command{
		Use:   "setup --token \"token-string\"",
		Short: "Automated config/install.",
		Long:  "Automated config/install.",
		Run:   cmd.setupApp,
	}

	cmd.RootCommand.AddCommand(addTokenCmd)

	addTokenCmd.Flags().StringP("token", "t", "", "Set user API token.")
	_ = addTokenCmd.MarkFlagRequired("token")
}

func (cmd *CMD) setupApp(cobraCommand *cobra.Command, args []string) {
	_ = args
	cmd.runMiddleware.ParseGlobalFlags(cobraCommand)
	token, _ := cobraCommand.Flags().GetString("token")
	var err error

	if runtime.GOOS == "windows" {
		path := cmd.configFileName
		err = initCreateFile(path, token)
	} else {
		err = initAppUnixLike(token, cmd.configFileName)
	}
	if err != nil {
		cmd.runMiddleware.OutputError(err)
	} else {
		fmt.Println("Init success.")
	}
}

func initAppUnixLike(token, configFileName string) error {
	// check user UID
	usr, err := user.Current()
	if err != nil {
		return err
	}
	if usr.Uid != "0" {
		return fmt.Errorf("You must have root privileges to use the 'setup' command.")
	}

	// Install bin in /usr/local/bin
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
	path = path + "/" + configFileName
	return initCreateFile(path, token)
}

func initCreateFile(path, token string) error {
	data := map[string]interface{}{
		"token": token,
		"uri":   api.DefaultURI,
	}
	viper.Set("default", data)
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	err := viper.WriteConfigAs(path)
	return err
}
