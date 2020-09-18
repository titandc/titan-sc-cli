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

var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(titan-sc completion bash)

# To load completions for each session, execute once:
Linux:
  $ titan-sc completion bash > /etc/bash_completion.d/titan-sc
MacOS:
  $ titan-sc completion bash > /usr/local/etc/bash_completion.d/titan-sc

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ titan-sc completion zsh > "${fpath[1]}/_titan-sc"

# You will need to start a new shell for this setup to take effect.

Fish:

$ titan-sc completion fish | source

# To load completions for each session, execute once:
$ titan-sc completion fish > ~/.config/fish/completions/titan-sc.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	companyCmdAdd()
	serverCmdAdd()
	snapshotCmdAdd()
	historyCmdAdd()
	utilsCmdAdd()
	networkCmdAdd()
	firewallCmdAdd()
	kvmIpCmdAdd()
	ipCmdAdd()
	userCmdAdd()
	sshKeysCmdAdd()
	versionCmdAdd()
	pnatCmdAdd()
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
	_ = args
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
