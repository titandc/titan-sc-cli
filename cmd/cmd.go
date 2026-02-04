package cmd

import (
	"fmt"
	"os"
	"strings"

	"titan-sc/run"

	"github.com/spf13/cobra"
)

type CMD struct {
	runMiddleware  *run.RunMiddleware
	RootCommand    *cobra.Command
	tokenDefined   bool
	configFileName string
	VersionMajor   int
	VersionMinor   int
	VersionPatch   int
}

func NewCMD(programName, configFileName string, tokenDefined bool, runMiddleware *run.RunMiddleware, versionMajor,
	verionsMinor, versionPatch int) *CMD {
	// Preserve command order as defined (don't sort alphabetically)
	cobra.EnableCommandSorting = false

	cmd := &CMD{
		RootCommand: &cobra.Command{
			Use:           programName,
			Short:         "Titan SC CLI",
			Long:          "Titan Small Cloud - Command Line Interface",
			SilenceErrors: true, // We handle errors ourselves
			SilenceUsage:  true, // Don't show usage on error
		},
		tokenDefined:   tokenDefined,
		runMiddleware:  runMiddleware,
		configFileName: configFileName,
		VersionMajor:   versionMajor,
		VersionMinor:   verionsMinor,
		VersionPatch:   versionPatch,
	}
	cmd.RootCommand.PersistentPreRun = cmd.checkTokenRequirement
	return cmd
}

func (cmd *CMD) Execute() {
	if err := cmd.RootCommand.Execute(); err != nil {
		errMsg := err.Error()

		// For required flag errors, show the subcommand's help
		if strings.Contains(errMsg, "required flag") {
			fmt.Fprintf(os.Stderr, "Error: %s\n\n", errMsg)
			subCmd, _, _ := cmd.RootCommand.Find(os.Args[1:])
			if subCmd != nil {
				subCmd.Help()
			}
			os.Exit(1)
		}

		// For unknown command errors, show suggestion and root help
		if strings.Contains(errMsg, "unknown command") {
			fmt.Fprintln(os.Stderr, "Error:", errMsg)
			fmt.Fprintln(os.Stderr)
			cmd.RootCommand.Help()
			os.Exit(1)
		}

		// Check for json flag in args (flag parsing may have failed)
		cmd.runMiddleware.JSONOutput = hasJSONFlag(os.Args)
		cmd.runMiddleware.OutputError(err)
		os.Exit(1)
	}
}

// hasJSONFlag checks if -j or --json flag is present in args
func hasJSONFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-j" || arg == "--json" {
			return true
		}
	}
	return false
}

func (cmd *CMD) checkTokenRequirement(cobraCommand *cobra.Command, args []string) {
	_ = args
	arrCmd := strings.SplitN(cobraCommand.CommandPath(), " ", 3)
	if len(arrCmd) > 1 && arrCmd[1] == "setup" {
		return
	}
	if len(arrCmd) > 2 && arrCmd[1] == "version" && arrCmd[2] == "cli" {
		return
	}
	if len(arrCmd) == 2 && arrCmd[1] == "version" {
		return
	}
	if !cmd.tokenDefined {
		fmt.Fprintln(os.Stderr, "Error: Unable to retrieve token from configuration file.")
		fmt.Fprintln(os.Stderr, "Run 'titan-sc setup' to configure your API token.")
		os.Exit(1)
	}
}
