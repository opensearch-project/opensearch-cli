- [User Guide](#user-guide)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Create default profile](#create-default-profile)
      - [Examples](#examples)
    - [List existing profile](#list-existing-profile)
    - [Using profile with opensearch-cli command](#using-profile-with-opensearch-cli-command)

# User Guide

## Installation

You can download the binaries directly from the [downloads](https://opensearch.org/downloads.html) page
or from the [releases](https://github.com/opensearch-project/opensearch-cli/releases) section.

## Usage

```
$ opensearch-cli --help

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
  -c, --config string    Configuration file for opensearch-cli, default is /Users/balasvij/.opensearch-cli/config.yaml
  -h, --help             Help for opensearch-cli
  -p, --profile string   Use a specific profile from your configuration file
  -v, --version          Version for opensearch-cli

```

### Create default profile
A profile is a collection of credentials that will be applied to the opensearch-cli command. When a user specifies a profile, 
the settings and credentials of that profile will be used to execute the command.
Users can create one profile with the name "default", and is used when no profile is explicitly referenced. 

#### Examples

1. Create default profile where the cluster's security uses HTTP basic authentication.
```
$ opensearch-cli profile create --auth-type "basic" \
                          --name "default" \
                          --endpoint "https://localhost:9200" 
Username: admin
Password: *******
Profile created successfully.
```
2. Create default profile where the cluster's security uses AWS IAM ARNs as users.
AWS credentials can be provided either by specifying aws profile name or using [environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).
You can find details about creating aws profiles [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html).
```
$ opensearch-cli profile create --auth-type "aws-iam" \
                          --name "default" \
                          --endpoint "https://localhost:9200" 
AWS profile name (leave blank if you want to provide credentials using environment variables): readonly      
AWS service name where your cluster is deployed (for Amazon Elasticsearch Service, use 'es'. For EC2, use 'ec2'): es
Profile created successfully.
```
3. Create default profile where the cluster's security plugin is disabled.
```
$ opensearch-cli profile create --auth-type "disabled" \
                          --name "default" \
                          --endpoint "https://cloud-service-endpoint:9200" 
Profile created successfully.
```

### List existing profile

```
$ opensearch-cli profile list -l
Name         UserName            Endpoint-url             
----         --------            ------------              
default      admin               https://localhost:9200   
prod         admin               https://node1:9200
                 
```

### Using profile with opensearch-cli command

You can specify profiles in two ways.

1. The first way is to add the --profile <name> option:    
    ```
    $ opensearch-cli ad stop-detector invalid-logins --profile prod
    ```
    This example stops the invalid-logins detector using the credentials and settings in the prod profile.
    
2. The second way is to use an environment variable.

    On Linux or macOS :
    ```
    $ export OPENSEARCH_PROFILE=prod
    ```
    Windows
    ```
    C:\> setx OPENSEARCH_PROFILE prod
    ```
   These variables last for the duration of your shell session, but you can add them to .zshenv or .bash_profile
   for a more permanent option.
