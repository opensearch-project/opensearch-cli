# ODFE Command Line Interface

ODFE Command Line Interface (odfe-cli) is an open source tool that lets you manage your Open Distro for Elasticsearch
cluster from the command line and automate tasks. In addition to standard Elasticsearch operations, you can configure,
manage, and use the ODFE plugins, such as Alerting, Anomaly Detection, and SQL

odfe-cli is best suited for situations in which you want to quickly combine a few commands, possibly adding them to
a script for easy access or automation. This example moves a detector "ecommerce-count-qualtity" from staging
to prod cluster, provided both profiles are available in config file.

```
odfe-cli ad get ecommerce-count-qualtity --profile stg > ecommerce-count-qualtity.json
odfe-cli ad create ecommerce-count-qualtity.json --profile prod
odfe-cli ad start ecommerce-count-qualtity.json --profile prod
odfe-cli ad stop ecommerce-count-qualtity --profile stg
odfe-cli ad delete ecommerce-count-qualtity --profile stg
```
## Installation:

You can download the binaries directly from the [downloads](https://opendistro.github.io/for-elasticsearch/downloads.html) page
or from the [releases](https://github.com/opendistro-for-elasticsearch/es-cli/releases) section.


## Development

### Minimum requirements

odfe-cli shares the same [minimum requirements](https://github.com/golang/go/wiki/MinimumRequirements#minimum-requirements) as Go.

### Build from source
1. Install [Go](https://golang.org/doc/install) > = 1.14
2. Clone the repository:
    ```
    cd $GOPATH/src
    git clone git@github.com:opendistro-for-elasticsearch/es-cli.git
    ```
3. Run build from source directory to generate binary:
   ```
   cd es-cli/odfe-cli
   go build .
   ```
4. Make binary executable:
    ```
    chmod +x ./odfe-cli
    ```

## Usage

```
odfe-cli --help
```

### Create default profile
A profile is a collection of credentials that will be applied to the odfe-cli command. When a user specifies a profile, 
the settings and credentials of that profile will be used to execute the command.
Users can create one profile with the name "default", and is used when no profile is explicitly referenced. 

```
$ odfe-cli profile create
Enter profile's name: default
Elasticsearch Endpoint: https://localhost:9200  
User Name: admin
Password: 
```

### List existing profile

```
$ odfe-cli profile list -l
Name         UserName            Endpoint-url             
----         --------            ------------              
default      admin               https://localhost:9200   
prod         admin               https://odfe-node1:9200
                 
```

### Using profile with odfe-cli command

You can specify profiles in two ways.

1. The first way is to add the --profile <name> option:    
    ```
    $ odfe ad stop-detector invalid-logins --profile prod
    ```
    This example stops the invalid-logins detector using the credentials and settings in the prod profile.
    
2. The second way is to use an environment variable.

    On Linux or macOS :
    ```
    $ export ODFE_PROFILE=prod
    ```
    Windows
    ```
    C:\> setx ODFE_PROFILE prod
    ```
   These variables last for the duration of your shell session, but you can add them to .zshenv or .bash_profile
   for a more permanent option.
    
## Security

See [CONTRIBUTING](https://github.com/opendistro-for-elasticsearch/es-cli/blob/main/CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.

