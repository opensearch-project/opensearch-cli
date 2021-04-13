# we will put our integration testing in this path
INTEGRATION_TEST_PATH=./it/...

# set of env variables that you need for testing
ENV_LOCAL_TEST=\
  OPENSEARCH_ENDPOINT="https://localhost:9200" \
  OPENSEARCH_USER="admin" \
  OPENSEARCH_PASSWORD="admin"

# this command will start a docker components that we set in docker-compose.yml
docker.start.components:
	docker-compose up -d;

# shutting down docker components
docker.stop:
	docker-compose down;

# clean up docker
docker.clean:
	docker system prune --volumes --force;

# this command will trigger integration test
# INTEGRATION_TEST_SUITE_PATH is used for run specific test in Golang
test.integration:
	go clean -testcache;
	$(ENV_LOCAL_TEST) go test -tags=integration $(INTEGRATION_TEST_PATH);

test.unit:
	go clean -testcache;
	go test ./...;

# format project using goimports tool
format:
	goimports -w .;

# generate opensearch-cli file in current directory
# update GOOS / GOARCH if you want to build for other operating systems and architecture
build:
	go build .
