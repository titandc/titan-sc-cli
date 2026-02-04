package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func (cmd *CMD) CompletionCmdAdd() {

	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `Generate shell completion script for auto-completing commands and flags.

Bash:
  # Add to current session:
  source <(titan-sc completion bash)

  # Install permanently:
  # Linux:
  titan-sc completion bash > /etc/bash_completion.d/titan-sc
  # macOS:
  titan-sc completion bash > /usr/local/etc/bash_completion.d/titan-sc

Zsh:
  # Add to current session:
  source <(titan-sc completion zsh)

  # Install permanently (requires shell restart):
  titan-sc completion zsh > ~/.titan/_titan-sc
  echo 'fpath=(~/.titan $fpath)' >> ~/.zshrc
  echo 'autoload -Uz compinit && compinit' >> ~/.zshrc

Fish:
  # Add to current session:
  titan-sc completion fish | source

  # Install permanently:
  titan-sc completion fish > ~/.config/fish/completions/titan-sc.fish

PowerShell:
  # Add to current session:
  titan-sc completion powershell | Out-String | Invoke-Expression

  # Install permanently (add to $PROFILE):
  titan-sc completion powershell >> $PROFILE

After installation, restart your shell or source the config file.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
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
