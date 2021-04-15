# Developer Guide

So you want to contribute code to OpenSearch CLI?  Excellent!  We're glad you're here.  Here's what you need to do:

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
If you have not installed docker-compose, you can install it from this [link](https://docs.docker.com/compose/install/).

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

### Submitting your changes

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


### Reviewing and accepting your contribution

We deeply appreciate everyone who takes the time to make a contribution.  We will review all contributions as quickly as possible, but there are a few things you can do to help us with the process:

First and foremost, opening an issue and discussing your change before you make it is the best way to smooth the PR process.  This will prevent a rejection because someone else is already working on the problem, or because the solution is incompatible with our architectural direction.

Additionally:
1) Please make sure you've run `make format` before submitting to format according to our guideline.
2) Please make sure your change includes the tests that correspond with your changes.
3) Please make sure local tests pass, add unit tests for all the new code paths introduced by your change.
4) Please write both unit tests and integration test for your change
5) Smaller changes are easier to digest than large ones.
6) Given the limits of the team, we will reject PRs that are simple refactorings or "tidying up".  So make sure you're clear about what problem your PR is solving.

During the PR process, expect that they'll be some back and forth.  Please try to respond to comments in a timely fashion, and if you don't wish to continue with the PR, let us know.  If a PR takes too many iterations for its complexity or size, we may reject it.  Additionally, if you stop responding, we may close the PR as abandoned.  In either case, if you feel this was done in error, please add a comment on the PR.

If we accept the PR, we will merge your change and usually take care of backporting it to appropriate branches ourselves.

If we reject the PR, we will close the pull request with a comment explaining why. This decision isn't always final: if you feel we have
misunderstood your intended change or otherwise think that we should reconsider then please continue the conversation with a comment on the pull request, and we'll do our best to address any further points you raise.
