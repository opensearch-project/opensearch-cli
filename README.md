![Build](https://github.com/opendistro-for-elasticsearch/odfe-cli/workflows/Build%20and%20Test%20odfe-cli/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/opendistro-for-elasticsearch/odfe-cli/branch/main/graph/badge.svg?flag=odfe-cli)](https://codecov.io/gh/opendistro-for-elasticsearch/odfe-cli)
[![Documentation](https://img.shields.io/badge/doc-reference-blue)](https://opendistro.github.io/for-elasticsearch-docs/docs/cli/)
[![Chat](https://img.shields.io/badge/chat-on%20forums-blue)](https://discuss.opendistrocommunity.dev/c/cli/55)
![PRs welcome!](https://img.shields.io/badge/PRs-welcome!-success)
# OpenSearch Command Line Interface

OpenSearch Command Line Interface (opensearch-cli) is an open source tool that lets you manage your OpenSearch cluster from the command line
and automate tasks. In addition to standard OpenSearch operations, you can configure,
manage, and use the plugins, such as Alerting, Anomaly Detection, and SQL

opensearch-cli is best suited for situations in which you want to quickly combine a few commands, possibly adding them to
a script for easy access or automation. This example moves a detector "ecommerce-count-qualtity" from staging
to prod cluster, provided both profiles are available in config file.

```
opensearch-cli ad get     ecommerce-count-qualtity      --profile stg > ecommerce-count-qualtity.json
opensearch-cli ad create  ecommerce-count-qualtity.json --profile prod
opensearch-cli ad start   ecommerce-count-qualtity.json --profile prod
opensearch-cli ad stop    ecommerce-count-qualtity      --profile stg
opensearch-cli ad delete  ecommerce-count-qualtity      --profile stg
```
## Installation:

You can download the binaries directly from the [downloads](https://opendistro.github.io/for-elasticsearch/downloads.html) page
or from the [releases](https://github.com/opendistro-for-elasticsearch/odfe-cli/releases) section.


## Development

### Minimum requirements

opensearch-cli shares [minimum requirements](https://github.com/golang/go/wiki/MinimumRequirements#minimum-requirements) as Go
and [docker](https://docs.docker.com/get-docker/) to run integration tests.

### Build from source
1. Install [Go](https://golang.org/doc/install) > = 1.16
2. Clone the repository:
    ```
    cd $GOPATH/src
    git clone git@github.com:opensearch-project/opensearch-cli.git
    ```
3. Run build from source directory to generate binary:
   ```
   cd opensearch-cli
   go build .
   ```
4. Make binary executable:
    ```
    chmod +x ./opensearch-cli
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
In order to test opensearch-cli end-to-end, we need a running OpenSearch cluster. We can use Docker to accomplish this. 
The [Docker Compose file](./docker-compose.yml) supports the ability to run integration tests for the project in local environments respectively.
If you have not installed docker-compose, you can install it from this [link](https://docs.docker.com/compose/install/)

Integration tests are often slower, so you may want to only run them after the unit test. In order to differentiate unit tests from integration tests, Go has a built-in mechanism for allowing you to logically separate your tests
with build tags. The build tag needs to be placed as close to the top of the file as possible, and must have a blank line beneath it.   
We recommend you to create all integration tests inside [this](./it) folder with build tag 'integration'.

#### Execute test integration command from your CLI
1. Run docker compose to start containers, by default it will launch latest OpenSearch cluster.
    ```
    docker-compose up -d;
    ```
2. Run all integration tests with build tag 'integration'
    ```
    go test -tags=integration ./it/...
    ```

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

#### Examples:

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
    
## Security

See [CONTRIBUTING](https://github.com/opendistro-for-elasticsearch/opensearch-cli/blob/main/CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.

