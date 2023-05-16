# Using the opensearch-cli

This section introduces you to many of the common features and options available in the 
OpenSearch Command Line Interface (opensearch-cli).


**Topics**
+ [Getting help](./usage.md#getting-help)
+ [Command structure](./usage.md#command-structure)
+ [Specifying parameter values](./usage.md#specifying-parameter-values)
+ [Auto complete](./usage.md#auto-complete)
+ [Environment variables](./usage.md#environment-variables)


## Getting help

You can get help with any command when using the OpenSearch Command Line Interface.
To do so, simply type `--help` or `help` at the end of a command name. 

For example, the following command displays help for the general opensearch-cli options and the available top-level commands. 

```
$ opensearch-cli help

opensearch-cli is a unified command line interface for managing OpenSearch clusters

Usage:
  opensearch-cli [command]

Available Commands:
  ad          Manage the Anomaly Detection plugin
  completion  Generate completion script for your shell
  curl        Manage OpenSearch core features
  help        Help about any command
  knn         Manage the k-NN plugin
  profile     Manage a collection of settings and credentials that you can apply to an opensearch-cli command

Flags:
  -c, --config string    Configuration file for opensearch-cli, default is /Users//.opensearch-cli/config.yaml
  -h, --help             Help for opensearch-cli
  -p, --profile string   Use a specific profile from your configuration file
  -v, --version          Version for opensearch-cli

Use "opensearch-cli [command] --help" for more information about a command.

```

For example, the following command displays Anomaly Detection plugin specific commands. 

```
$ opensearch-cli ad help

Use the Anomaly Detection commands to create, configure, and manage detectors.

Usage:
  opensearch-cli ad [command]

Available Commands:
  create      Create detectors based on JSON files
  delete      Delete detectors based on a list of IDs, names, or name regex patterns
  get         Get detectors based on a list of IDs, names, or name regex patterns
  start       Start detectors based on a list of IDs, names, or name regex patterns
  stop        Stop detectors based on a list of IDs, names, or name regex patterns
  update      Update detectors based on JSON files

Flags:
  -h, --help   Help for Anomaly Detection

Global Flags:
  -c, --config string    Configuration file for opensearch-cli, default is /Users/opensearch/.opensearch-cli/config.yaml
  -p, --profile string   Use a specific profile from your configuration file

Use "opensearch-cli ad [command] --help" for more information about a command.

```


## Command structure

The opensearch-cli uses a multipart structure on the command line that must be specified in this order:

1. The base call to the `opensearch-cli` program.

1. The top-level command, which typically corresponds to OpenSearch plugins supported by the opensearch-cli.

1. The subcommand that specifies which operation to perform.

1. General CLI flags or parameters required by the operation.

    ```
    $ opensearch-cli <command> <subcommand> [flags and parameters]
    ```

## Specifying parameter values

Many flags used in the opensearch-cli are simple string, numeric values or boolean, such as the path and pretty in the following example. 
Here, path accepts string value while pretty is a boolean flag.

**Note:** you must use quotation marks around strings that include one or more space characters
```
$ opensearch-cli curl get --path _cluster/health --pretty
```

## Auto complete
opensearch-cli includes a command-completion feature that enables you to use the Tab key to complete a partially entered command.
This feature isn't automatically installed, you need to configure it manually.

### Configure auto complete

To enable the auto complete feature, you can use the following command and perform operation as mentioned in the help,
based on your operating system
```
opensearch-cli completion --help

To enable shell autocompletion:

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

Usage:
  opensearch-cli completion [ bash | zsh | fish | powershell ]

Flags:
  -h, --help   help for completion

Global Flags:
  -c, --config string    Configuration file for opensearch-cli, default is /Users/opensearch/.opensearch-cli/config.yaml
  -p, --profile string   Use a specific profile from your configuration file

```
### How to use Auto complete
When you partially enter a command, parameter, or flag, the command-completion feature either automatically
completes your command or displays a suggested list of commands. To prompt command completion, you partially enter a command and press Tab.

The following examples show different ways that you can use command completion:
* Partially enter a command and press Tab to display a suggested list of commands.                                                                                      
```
opensearch-cli curl p[tab key]
patch    post    put
```
* Partially enter a parameter and press Tab to display a suggested list of parameters.

```
opensearch-cli curl get --[tab key]
--config          --data            --filter-path     --headers         --help            --output-format=  --path=           --profile         --query-params    
--config=         --data=           --filter-path=    --headers=        --output-format   --path            --pretty          --profile=        --query-params= 
```
* Enter a command and press Tab to display sub commands.
```
opensearch-cli curl [tab key]
delete  get     patch     post    put     
```

## Environment variables

The opensearch-cli supports the following environment variables.

`OPENSEARCH_CONFIG_FILE`  
Specifies the location of the file that the opensearch-cli saves configuration profiles.
The default file location is `~/.opensearch-cli/config.yaml`.

`OPENSEARCH_MAX_RETRY`  
Specifies a value of maximum retry attempts the opensearch-cli performs, excluding initial call.
If defined, `OPENSEARCH_MAX_RETRY` overrides the value for the individual profiles setting `max_retry`.

`OPENSEARCH_PROFILE`  
Specifies the name of the ofe-cli profile to use.
If defined, this environment variable overrides the behavior of using the profile named `[default]` in the configuration file.
You can override this environment variable by using the `--profile` command line parameter.

`OPENSEARCH_TIMEOUT`  
Specifies maximum time  in  seconds  that you allow the connection to the server to take.
If defined, `OPENSEARCH_TIMEOUT` overrides the value for the individual profiles setting `timeout`.
This only limits  the  connection  phase, once timeout happens, client will only exit, it doesn't terminate the
request that already reached the server.
