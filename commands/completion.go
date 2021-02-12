/*
 * Copyright 2021 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
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

$ source <(odfe-cli completion bash)

# To enable auto completion for commands for each session, execute once:
Linux:
  $ odfe-cli completion bash > /etc/bash_completion.d/odfe-cli
MacOS:
  $ odfe-cli completion bash > /usr/local/etc/bash_completion.d/odfe-cli

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To enable auto completion for commands for each session, execute once:
$ odfe-cli completion zsh > "${fpath[1]}/_odfe-cli"

# You will need to start a new shell for this setup to take effect.

Fish:

$ odfe-cli completion fish | source

# To enable auto completion for commands for each session, execute once:
$ odfe-cli completion fish > ~/.config/fish/completions/odfe-cli.fish

Powershell:

PS> odfe-cli completion powershell | Out-String | Invoke-Expression

# To enable auto completion for commands for each session, execute once:
PS> odfe-cli completion powershell > odfe-cli.ps1
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
