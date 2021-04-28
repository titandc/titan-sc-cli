package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func (cmd *CMD) CompletionCmdAdd() {

	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
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

	cmd.RootCommand.AddCommand(completionCmd)
}
