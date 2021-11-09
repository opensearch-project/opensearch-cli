/*
 * SPDX-License-Identifier: Apache-2.0
 *
 * The OpenSearch Contributors require contributions made to
 * this file be licensed under the Apache-2.0 license or a
 * compatible open source license.
 *
 * Modifications Copyright OpenSearch Contributors. See
 * GitHub history for details.
 */

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	CompletionCommandName = "completion"
	BashShell             = "bash"
	ZshShell              = "zsh"
	FishShell             = "fish"
	PowerShell            = "powershell"
)

var longText = `To enable shell autocompletion:

Bash:

$ source <(opensearch-cli completion bash)

# To enable auto completion for commands for each session, execute once:
Linux:
  $ opensearch-cli completion bash > /etc/bash_completion.d/opensearch-cli
MacOS:
  $ opensearch-cli completion bash > /usr/local/etc/bash_completion.d/opensearch-cli

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To enable auto completion for commands for each session, execute once:
$ opensearch-cli completion zsh > "${fpath[1]}/_opensearch-cli"

# You will need to start a new shell for this setup to take effect.

Fish:

$ opensearch-cli completion fish | source

# To enable auto completion for commands for each session, execute once:
$ opensearch-cli completion fish > ~/.config/fish/completions/opensearch-cli.fish

Powershell:

PS> opensearch-cli completion powershell | Out-String | Invoke-Expression

# To enable auto completion for commands for each session, execute once:
PS> opensearch-cli completion powershell > opensearch-cli.ps1
# and source this file from your powershell profile.
`

var completionCmd = &cobra.Command{
	Use:                   fmt.Sprintf("%s [ %s | %s | %s | %s ]", CompletionCommandName, BashShell, ZshShell, FishShell, PowerShell),
	Short:                 "Generate completion script for your shell",
	Long:                  longText,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{BashShell, ZshShell, FishShell, PowerShell},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case BashShell:
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case ZshShell:
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case FishShell:
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		case PowerShell:
			err = cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
		DisplayError(err, CompletionCommandName)
	},
}

func init() {
	GetRoot().AddCommand(completionCmd)
}
