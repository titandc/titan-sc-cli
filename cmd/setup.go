package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
	"titan-sc/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	AppName  = "titan-sc"
	LocalApp = "./titan-sc"
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
	token, err := cobraCommand.Flags().GetString("token")
	if err != nil {
		fmt.Println("Fail to parse --token option.")
		return
	}

	if runtime.GOOS == "windows" {
		err = installOnWindows(token, cmd.configFileName)
	} else {
		err = installOnUnixLike(token, cmd.configFileName)
	}

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Init success.")
	}
}

func installOnWindows(token, path string) error {
	return initCreateFile(token, path)
}

func installOnUnixLike(token, configFileName string) error {
	if err := installIsRoot(); err != nil {
		return err
	}

	if err := installBin(LocalApp); err != nil {
		return err
	}

	return installConfigFile(token, configFileName)
}

func installIsRoot() error {
	// check user UID
	usr, err := user.Current()
	if err != nil {
		return err
	}
	if usr.Uid != "0" {
		return errors.New("you must have root privileges to use the 'setup' command")
	}
	return nil
}

func installBin(src string) error {
	dirs := []string{"/usr/local/bin/", "/usr/bin"}

	dstPath, err := deleteBinIfExists(dirs)
	if err != nil {
		return err
	}

	if dstPath == "" {
		dir, err := getBinDestination(dirs)
		if err != nil {
			return err
		}
		dstPath = fmt.Sprintf("%s%s", dir, AppName)
	}

	if err = os.Link(src, dstPath); err != nil {
		return err
	}
	return nil
}

func deleteBinIfExists(dirs []string) (string, error) {
	var dstPath string

	for _, dir := range dirs {
		dstPath = fmt.Sprintf("%s%s", dir, AppName)
		if pathExists(dstPath) {
			if err := os.Remove(dstPath); err != nil {
				return "", err
			}
			return dstPath, nil
		}
	}
	return "", nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getBinDestination(dirs []string) (string, error) {
	for _, dir := range dirs {
		if pathExists(dir) {
			return dir, nil
		}
	}
	return "", errors.New("no available path found for install CLI")
}

func installConfigFile(token, configFileName string) error {
	path := os.Getenv("HOME") + "/" + ".titan"
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		if !strings.Contains(err.Error(), ": file exists") {
			return fmt.Errorf("init: create dir path <%s> error", path)
		}
	}
	path = path + "/" + configFileName
	return initCreateFile(token, path)
}

func initCreateFile(token, path string) error {
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
