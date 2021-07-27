<img src="https://opensearch.org/assets/brand/SVG/Logo/opensearch_logo_default.svg" height="64px"/>

- [OpenSearch Command Line Interface](#opensearch-command-line-interface)
- [Usage](#usage)
- [Contributing](#contributing)
- [Getting Help](#getting-help)
- [Code of Conduct](#code-of-conduct)
- [Security](#security)
- [License](#license)
- [Copyright](#copyright)

[![tests](https://github.com/opensearch-project/opensearch-cli/actions/workflows/test-build-workflow.yml/badge.svg?branch=main)](https://github.com/opensearch-project/opensearch-cli/actions/workflows/test-build-workflow.yml)
[![codecov](https://codecov.io/gh/opensearch-project/opensearch-cli/branch/main/graph/badge.svg?token=BTU8GI0LJP)](https://codecov.io/gh/opensearch-project/opensearch-cli)
[![Documentation](https://img.shields.io/badge/doc-reference-blue)](https://opendistro.github.io/for-elasticsearch-docs/docs/cli/)
[![Chat](https://img.shields.io/badge/chat-on%20forums-blue)](https://discuss.opendistrocommunity.dev/c/cli/55)
![PRs welcome!](https://img.shields.io/badge/PRs-welcome!-success)

## OpenSearch Command Line Interface

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

## Installing

You can download the binaries directly from the [downloads](https://opendistro.github.io/for-elasticsearch/downloads.html) page
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

#### Example: Create default profile where the cluster's security uses HTTP basic authentication.

```
$ opensearch-cli profile create --auth-type "basic" \
                          --name "default" \
                          --endpoint "https://localhost:9200" 
Username: admin
Password: *******
Profile created successfully.
```
#### Example: Create default profile where the cluster's security uses AWS IAM ARNs as users.

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
#### Example: Create default profile where the cluster's security plugin is disabled.

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
    
## Contributing

See [developer guide](DEVELOPER_GUIDE.md) and [how to contribute to this project](CONTRIBUTING.md). 

## Getting Help

If you find a bug, or have a feature request, please don't hesitate to open an issue in this repository.

For more information, see [project website](https://opensearch.org/) and [documentation](https://docs-beta.opensearch.org/). If you need help and are unsure where to open an issue, try [forums](https://discuss.opendistrocommunity.dev/).

## Code of Conduct

This project has adopted the [Amazon Open Source Code of Conduct](CODE_OF_CONDUCT.md). For more information see the [Code of Conduct FAQ](https://aws.github.io/code-of-conduct-faq), or contact [opensource-codeofconduct@amazon.com](mailto:opensource-codeofconduct@amazon.com) with any additional questions or comments.

## Security

If you discover a potential security issue in this project we ask that you notify AWS/Amazon Security via our [vulnerability reporting page](http://aws.amazon.com/security/vulnerability-reporting/). Please do **not** create a public GitHub issue.

## License

This project is licensed under the [Apache v2.0 License](LICENSE.txt).

## Copyright

Copyright 2020-2021 Amazon.com, Inc. or its affiliates. All Rights Reserved.
