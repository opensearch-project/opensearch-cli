[![tests](https://github.com/opensearch-project/opensearch-cli/actions/workflows/test-build-workflow.yml/badge.svg?branch=main)](https://github.com/opensearch-project/opensearch-cli/actions/workflows/test-build-workflow.yml)
[![codecov](https://codecov.io/gh/opensearch-project/opensearch-cli/branch/main/graph/badge.svg?token=BTU8GI0LJP)](https://codecov.io/gh/opensearch-project/opensearch-cli)
[![Documentation](https://img.shields.io/badge/doc-reference-blue)](https://opensearch.org/docs/latest/tools/cli/)
[![Chat](https://img.shields.io/badge/chat-on%20forums-blue)](https://forum.opensearch.org/)
![PRs welcome!](https://img.shields.io/badge/PRs-welcome!-success)

- [OpenSearch Command Line Interface](#opensearch-command-line-interface)
  - [User Guide](#user-guide)
  - [Compatibility with OpenSearch](#compatibility-with-opensearch)
  - [Security](#security)
  - [License](#license)
  - [Copyright](#copyright)

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

## User Guide

See [User Guide](USER_GUIDE.md).

## Compatibility with OpenSearch

See [Compatibility](COMPATIBILITY.md).

## Security

See [SECURITY](SECURITY.md) for more information.

## License

This project is licensed under the [Apache v2.0 License](LICENSE.txt).

## Copyright

Copyright OpenSearch Contributors. See [NOTICE](NOTICE.txt) for details.
