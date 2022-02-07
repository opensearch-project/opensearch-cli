// +build integration

/*
 * SPDX-License-Identifier: Apache-2.0
 *
 * The OpenSearch Contributors require contributions made to
 * this file be licensed under the Apache-2.0 license or a
 * compatible open source license.
 *
 * Modifications Copyright OpenSearch Contributors. See
 * GitHub history for details.
 */

package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"opensearch-cli/client"
	ctrl "opensearch-cli/controller/platform"
	"opensearch-cli/entity"
	"opensearch-cli/entity/platform"
	"opensearch-cli/environment"
	gateway "opensearch-cli/gateway/platform"
	"opensearch-cli/it"
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const GetBulkIndexName = "bulk-user-request"

//OpenSearchTestSuite suite tests OpenSearch REST API REQUESTS
type OpenSearchTestSuite struct {
	it.CLISuite
	Gateway    gateway.Gateway
	Controller ctrl.Controller
}

type result struct {
	Source map[string]interface{} `json:"_source"`
}

//SetupSuite runs once for every test suite
func (a *OpenSearchTestSuite) SetupSuite() {
	var err error
	a.Client, err = client.New(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.Profile = &entity.Profile{
		Name:     "test",
		Endpoint: os.Getenv(environment.OPENSEARCH_ENDPOINT),
		UserName: os.Getenv(environment.OPENSEARCH_USER),
		Password: os.Getenv(environment.OPENSEARCH_PASSWORD),
	}
	if err = a.ValidateProfile(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.Gateway, _ = gateway.New(a.Client, a.Profile)
	a.Controller = ctrl.New(a.Gateway)
	a.CreateIndex(GetBulkIndexName, "")
}
func (a *OpenSearchTestSuite) TearDownSuite() {
	a.DeleteIndex(GetBulkIndexName)
}

func (a *OpenSearchTestSuite) TestCurlGet() {
	request := platform.CurlCommandRequest{
		Action: "Get",
		Pretty: true,
	}
	a.T().Run("get document count for an index", func(t *testing.T) {
		ctx := context.Background()
		request.Path = fmt.Sprintf("_cat/count/%s", GetBulkIndexName)
		request.QueryParams = "v=true"
		response, err := a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)
		assert.True(t, strings.Contains(string(response), "5"))
		assert.True(t, strings.Contains(string(response), "epoch"))
		assert.True(t, strings.Contains(string(response), "timestamp"))
		assert.True(t, strings.Contains(string(response), "count"))
	})

	a.T().Run("health status of a cluster", func(t *testing.T) {
		ctx := context.Background()
		request.QueryParams = ""
		request.Path = "_cluster/health"
		response, err := a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)
		var health map[string]interface{}
		assert.NoError(t, json.Unmarshal(response, &health))
		assert.True(t, len(health) > 0)
		assert.NotNil(t, health["status"])
		assert.NotNil(t, health["number_of_nodes"])
	})
	a.T().Run("health status of a cluster in yaml", func(t *testing.T) {
		ctx := context.Background()
		request.QueryParams = ""
		request.Path = "_cluster/health"
		request.OutputFormat = "yaml"
		response, err := a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)
		var health struct {
			Nodes  string `yaml:"number_of_nodes"`
			Status string `yaml:"status"`
		}
		assert.NoError(t, yaml.Unmarshal(response, &health))
		assert.True(t, len(health.Status) > 0)
		assert.True(t, len(health.Nodes) > 0)
	})
}

func (a *OpenSearchTestSuite) TestCurlPatch() {
	plugins := []string{"opensearch-security"}
	if !a.IsPluginFromInputInstalled(plugins) {
		a.T().Skipf("plugin %s is not installed", plugins)
	}
	request := platform.CurlCommandRequest{
		Action: "PATCH",
		Pretty: true,
	}
	a.T().Run("add user to a role mapping", func(t *testing.T) {
		ctx := context.Background()
		request.Path = "_plugins/_security/api/rolesmapping"
		request.QueryParams = ""
		request.Data = `[{ "op": "add", "path": "/all_access/users", "value": ["testuser"] }]`
		response, err := a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)

		var result map[string]interface{}
		assert.NoError(t, json.Unmarshal(response, &result))
		assert.True(t, len(result) > 0)
		assert.EqualValues(t, "OK", result["status"])
		assert.EqualValues(t, "Resource updated.", result["message"])

		//get rolemappings
		request.Action = "GET"
		request.Path = "_plugins/_security/api/rolesmapping"
		request.QueryParams = ""
		request.Data = ""
		response, err = a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)

		var rolemappings map[string]struct {
			Hosts           []string `json:"hosts"`
			Users           []string `json:"users"`
			Reserved        bool     `json:"reserved"`
			Hidden          bool     `json:"hidden"`
			BackendRoles    []string `json:"backend_roles"`
			AndBackendRoles []string `json:"and_backend_roles"`
		}
		assert.NoError(t, json.Unmarshal(response, &rolemappings))
		assert.True(t, len(rolemappings["all_access"].Users) > 0)
		assert.EqualValues(t, "testuser", rolemappings["all_access"].Users[0])
	})
}

func (a *OpenSearchTestSuite) TestCurlPost() {
	request := platform.CurlCommandRequest{
		Action: "Post",
		Pretty: true,
	}
	expectedDocument := "Taming Text: How to Find, Organize, and Manipulate It"
	a.T().Run("bulk request", func(t *testing.T) {
		ctx := context.Background()
		request.Path = "test-index-3/_bulk"
		request.QueryParams = "refresh"
		request.Data = `@testdata/sample-index`
		response, err := a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)

		//get document
		request.Action = "GET"
		request.Path = "test-index-3/_doc/2"
		request.QueryParams = ""
		request.Data = ""
		response, err = a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)

		var index result
		assert.NoError(t, json.Unmarshal(response, &index))
		assert.True(t, len(index.Source) > 0)
		assert.EqualValues(t, expectedDocument, index.Source["title"])
		assert.EqualValues(t, 12, index.Source["num_reviews"])

		a.DeleteIndex("test-index-3")
	})

	a.T().Run("bulk request compressed", func(t *testing.T) {
		ctx := context.Background()
		request.Action = "PUT"
		request.Path = "test-index-4/_bulk"
		request.QueryParams = "refresh"
		request.Data = `@testdata/sample-index-compressed.gz`
		request.Headers = "content-encoding: gzip"
		response, err := a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)

		//get document
		request.Action = "GET"
		request.Path = "test-index-4/_doc/2"
		request.QueryParams = ""
		request.Data = ""
		response, err = a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)

		var index result
		assert.NoError(t, json.Unmarshal(response, &index))
		assert.True(t, len(index.Source) > 0)
		assert.EqualValues(t, expectedDocument, index.Source["title"])
		assert.EqualValues(t, 12, index.Source["num_reviews"])

		a.DeleteIndex("test-index-3")
	})
}

func (a *OpenSearchTestSuite) TestCurlPut() {
	request := platform.CurlCommandRequest{
		Action: "PUT",
		Pretty: true,
	}
	a.T().Run("index a document", func(t *testing.T) {
		ctx := context.Background()
		request.Path = "test-index-2/_doc/1"
		request.QueryParams = ""
		request.Data = `{"message": "insert document","address": "127.0.0.1"}`
		response, err := a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)

		//get document
		request.Action = "GET"
		request.Path = "test-index-2/_doc/1"
		request.QueryParams = ""
		request.Data = ""
		response, err = a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)
		var index result
		assert.NoError(t, json.Unmarshal(response, &index))
		assert.True(t, len(index.Source) > 0)
		assert.EqualValues(t, "insert document", index.Source["message"])
		assert.EqualValues(t, "127.0.0.1", index.Source["address"])

		a.DeleteIndex("test-index-2")
	})
}

func (a *OpenSearchTestSuite) TestCurlDelete() {
	request := platform.CurlCommandRequest{
		Pretty: true,
	}
	a.T().Run("delete index document", func(t *testing.T) {
		ctx := context.Background()
		request.Action = "PUT"
		request.Path = "test-index-delete/_bulk"
		request.QueryParams = "refresh"
		request.Data = `@testdata/sample-index`
		response, err := a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)

		//delete document
		request.Action = "DELETE"
		request.Path = "test-index-delete/_doc/2"
		request.QueryParams = ""
		request.Data = ""
		response, err = a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)
		var result map[string]interface{}
		assert.NoError(t, json.Unmarshal(response, &result))
		assert.EqualValues(t, "deleted", result["result"])

		a.DeleteIndex("test-index-delete")
	})

	a.T().Run("delete index", func(t *testing.T) {
		ctx := context.Background()
		request.Action = "PUT"
		request.Path = "test-index-delete/_bulk"
		request.QueryParams = "refresh"
		request.Data = `@testdata/sample-index`
		response, err := a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)

		//delete document
		request.Action = "DELETE"
		request.Path = "test-index-delete"
		request.QueryParams = ""
		request.Data = ""
		response, err = a.Controller.Curl(ctx, request)
		assert.NoError(t, err, "failed to get response")
		assert.NotNil(t, response)
		var result map[string]interface{}
		assert.NoError(t, json.Unmarshal(response, &result))
		assert.EqualValues(t, true, result["acknowledged"])

		if err != nil {
			a.DeleteIndex("test-index-delete")
		}
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestESGETSuite(t *testing.T) {
	suite.Run(t, new(OpenSearchTestSuite))
}
