# Developer Guide

So you want to contribute code to OpenSearch CLI?  Excellent!  We're glad you're here.  Here's what you need to do:

- [Prerequisites](#prerequisites)
- [Unit Testing](#unit-testing)
- [Integration Testing](#integration-testing)
- [Submitting your changes](#submitting-your-changes)


## Prerequisites

* see [minimum requirements for Go](https://github.com/golang/go/wiki/MinimumRequirements#minimum-requirements) as OpenSearch shares the same requirements
* [docker](https://docs.docker.com/get-docker/) to run integration tests.

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

## Unit Testing
Go has a simple tool for running tests. To run every unit test, use this command:
 ```
go test ./...
```
 
However, often when writing tests, you may want to run your new test as below
```
cd folder-path/to/test;
go test -v -run TestName; 
```

## Integration Testing
In order to test opensearch-cli end-to-end, we need a running OpenSearch cluster. We can use Docker to accomplish this. 
The [Docker Compose file](./docker-compose.yml) supports the ability to run integration tests for the project in local environments respectively.
If you have not installed docker-compose, you can install it from this [link](https://docs.docker.com/compose/install/)

Integration tests are often slower, so you may want to only run them after the unit test. In order to differentiate unit tests from integration tests, Go has a built-in mechanism for allowing you to logically separate your tests
with build tags. The build tag needs to be placed as close to the top of the file as possible, and must have a blank line beneath it.   
We recommend you to create all integration tests inside [this](./it) folder with build tag 'integration'.

### Execute test integration command from your CLI
1. Run docker compose to start containers, by default it will launch latest OpenSearch cluster.
    ```
    docker-compose up -d;
    ```
2. Run all integration tests with build tag 'integration'
    ```
    go test -tags=integration ./it/...
    ```

## Submitting your changes

Once your changes and tests are ready to submit for review:

1. Test your changes

    Run the test suite to make sure local tests passes, add unit tests for all the new code paths introduced by your change,
    make sure nothing your change does not cause any regression.

2. Rebase your changes

Update your local repository with the most recent code from the main opensearch-cli repository, and rebase your branch on top of the latest master branch. We prefer your initial changes to be squashed into a single commit. Later, if we ask you to make changes, add them as separate commits.  This makes them easier to review.  As a final step before merging we will either ask you to squash all commits yourself, or we'll do it for you.

3. Please provide detailed description on your change.

   What is your change?
   What features are enabled in your change if any?
   What is fixed in your change if any?

4. Submit a pull request

    Push your local changes to your forked copy of the repository and [submit a pull request](https://help.github.com/articles/using-pull-requests). In the pull request, choose a title which sums up the changes that you have made, and in the body provide more details about your changes. Also, mention the number of the issue where discussion has taken place, eg "Closes #123".

Unless your change is trivial, there will probably be discussion about the pull request and, if any changes are needed, we would love to work with you to get your pull request merged into opensearch-cli.

5. Code review and approval

   Invite domain expert to review your code.

Please adhere to the general guideline that you should never force push
to a publicly shared branch. Once you have opened your pull request, you
should consider your branch publicly shared. Instead of force pushing
you can just add incremental commits; this is generally easier on your
reviewers. If you need to pick up changes from master, you can merge
master into your branch. A reviewer might ask you to rebase a
long-running pull request in which case force pushing is okay for that
request. Note that squashing at the end of the review process should
also not be done, that can be done when the pull request is [integrated
via GitHub](https://github.com/blog/2141-squash-your-commits).

