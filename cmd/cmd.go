package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"titan-sc/run"
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
	cmd := &CMD{
		RootCommand: &cobra.Command{
			Use:   programName,
			Short: "Titan SC CLI",
			Long:  "Titan Small Cloud - Command Line Interface",
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
		cmd.runMiddleware.OutputError(err)
		os.Exit(1)
	}
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
		// TODO: format better
		log.Println("Unable to retrieve token from configuration file.")
		os.Exit(1)
	}
}
