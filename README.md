![Build and Test odfe-cli](https://github.com/opendistro-for-elasticsearch/odfe-cli/workflows/Build%20and%20Test%20odfe-cli/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/opendistro-for-elasticsearch/odfe-cli/branch/main/graph/badge.svg?flag=odfe-cli)](https://codecov.io/gh/opendistro-for-elasticsearch/odfe-cli)
[![Documentation](https://img.shields.io/badge/doc-reference-blue)](https://opendistro.github.io/for-elasticsearch-docs/docs/cli/)
![PRs welcome!](https://img.shields.io/badge/PRs-welcome!-success)
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
or from the [releases](https://github.com/opendistro-for-elasticsearch/odfe-cli/releases) section.


## Supported versions and Interoperability

### Odfe cli and Opendistro plugins compatibility

| odfe plugins  | odfe versions |
| ------------- | --------------- |
| [Anomaly Detection](https://opendistro.github.io/for-elasticsearch-docs/docs/ad/)  | [1.12.0](https://github.com/opendistro-for-elasticsearch/opendistro-build/blob/main/release-notes/opendistro-for-elasticsearch-release-notes-1.12.0.md), [1.13.0](https://github.com/opendistro-for-elasticsearch/opendistro-build/blob/main/release-notes/opendistro-for-elasticsearch-release-notes-1.13.0.md), [1.13.1](https://github.com/opendistro-for-elasticsearch/opendistro-build/blob/main/release-notes/opendistro-for-elasticsearch-release-notes-1.13.1.md) |
| [k-NN](https://opendistro.github.io/for-elasticsearch-docs/docs/knn/)  | [1.12.0](https://github.com/opendistro-for-elasticsearch/opendistro-build/blob/main/release-notes/opendistro-for-elasticsearch-release-notes-1.12.0.md), [1.13.0](https://github.com/opendistro-for-elasticsearch/opendistro-build/blob/main/release-notes/opendistro-for-elasticsearch-release-notes-1.13.0.md), [1.13.1](https://github.com/opendistro-for-elasticsearch/opendistro-build/blob/main/release-notes/opendistro-for-elasticsearch-release-notes-1.13.1.md) |
            

## Development

### Minimum requirements

odfe-cli shares [minimum requirements](https://github.com/golang/go/wiki/MinimumRequirements#minimum-requirements) as Go
and [docker](https://docs.docker.com/get-docker/) to run integration tests.

### Build from source
1. Install [Go](https://golang.org/doc/install) > = 1.14
2. Clone the repository:
    ```
    cd $GOPATH/src
    git clone git@github.com:opendistro-for-elasticsearch/odfe-cli.git
    ```
3. Run build from source directory to generate binary:
   ```
   cd odfe-cli
   go build .
   ```
4. Make binary executable:
    ```
    chmod +x ./odfe-cli
    ```

### Unit Testing
Go has a simple tool for running tests. To run every unit test, use this command:
 ```
go test ./...
```
 
However, often when writing tests, you may want to run your new test as below
```
cd folder-path/to/test;
go test -v -run TestName; 
```

### Integration Testing
In order to test odfe-cli end-to-end, we need a running odfe cluster. We can use Docker to accomplish this. 
The [Docker Compose file](./docker-compose.yml) supports the ability to run integration tests for the project in local environments respectively.
If you have not installed docker-compose, you can install it from this [link](https://docs.docker.com/compose/install/)

Integration tests are often slower, so you may want to only run them after the unit test. In order to differentiate unit tests from integration tests, Go has a built-in mechanism for allowing you to logically separate your tests
with build tags. The build tag needs to be placed as close to the top of the file as possible, and must have a blank line beneath it.   
We recommend you to create all integration tests inside [this](./it) folder with build tag 'integration'.

#### Execute test integration command from your CLI
1. Run docker compose to start containers, by default it will launch latest odfe version.
    ```
    docker-compose up -d;
    ```
2. Run all integration tests with build tag 'integration'
    ```
    go test -tags=integration ./it/...
    ```

## Usage

```
$ odfe-cli --help

odfe-cli is a unified command line interface for managing ODFE clusters

Usage:
  odfe-cli [command]

Available Commands:
  ad          Manage the Anomaly Detection plugin
  completion  Generate completion script for your shell
  curl        Manage Elasticsearch core features
  help        Help about any command
  knn         Manage the k-NN plugin
  profile     Manage a collection of settings and credentials that you can apply to an odfe-cli command

Flags:
  -c, --config string    Configuration file for odfe-cli, default is /Users/balasvij/.odfe-cli/config.yaml
  -h, --help             Help for odfe-cli
  -p, --profile string   Use a specific profile from your configuration file
  -v, --version          Version for odfe-cli

```

### Create default profile
A profile is a collection of credentials that will be applied to the odfe-cli command. When a user specifies a profile, 
the settings and credentials of that profile will be used to execute the command.
Users can create one profile with the name "default", and is used when no profile is explicitly referenced. 

#### Examples:

1. Create default profile where the cluster's security uses HTTP basic authentication.
```
$ odfe-cli profile create --auth-type "basic" \
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
$ odfe-cli profile create --auth-type "aws-iam" \
                          --name "default" \
                          --endpoint "https://localhost:9200" 
AWS profile name (leave blank if you want to provide credentials using environment variables): readonly      
AWS service name where your cluster is deployed (for Amazon Elasticsearch Service, use 'es'. For EC2, use 'ec2'): es
Profile created successfully.
```
3. Create default profile where the cluster's security plugin is disabled.
```
$ odfe-cli profile create --auth-type "disabled" \
                          --name "default" \
                          --endpoint "https://cloud-service-endpoint:9200" 
Profile created successfully.
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
    $ odfe-cli ad stop-detector invalid-logins --profile prod
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

See [CONTRIBUTING](https://github.com/opendistro-for-elasticsearch/odfe-cli/blob/main/CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.

