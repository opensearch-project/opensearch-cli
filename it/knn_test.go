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

package it

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"opensearch-cli/client"
	ctrl "opensearch-cli/controller/knn"
	"opensearch-cli/entity"
	"opensearch-cli/environment"
	gateway "opensearch-cli/gateway/knn"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	KNNSampleIndexFileName        = "knn-sample-index"
	KnnSampleIndexMappingFileName = "knn-sample-index-mapping"
)

//KNNTestSuite suite specific to k-NN plugin
type KNNTestSuite struct {
	CLISuite
	Gateway    gateway.Gateway
	Controller ctrl.Controller
}

//SetupSuite runs once for every test suite
func (a *KNNTestSuite) SetupSuite() {
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
	a.Plugins = append(a.Plugins, "opensearch-knn")
	a.Gateway, _ = gateway.New(a.Client, a.Profile)
	a.Controller = ctrl.New(a.Gateway)
	a.CreateIndex(KNNSampleIndexFileName, KnnSampleIndexMappingFileName)
}
func (a *KNNTestSuite) TearDownSuite() {
	a.DeleteIndex(KNNSampleIndexFileName)
}

//GetNodesIDUsingRESTAPI helper to get node id using rest api
func (a *KNNTestSuite) GetNodesIDUsingRESTAPI(t *testing.T) string {
	indexURL := fmt.Sprintf("%s/_cat/nodes?full_id=true&h=id", a.Profile.Endpoint)
	response, err := a.callRequest(http.MethodGet, []byte(""), indexURL)
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSuffix(string(response), "\n")
}

func (a *KNNTestSuite) TestGetStatistics() {
	if !a.IsPluginInstalled() {
		a.T().Skipf("plugin %s is not installed", a.Plugins)
	}
	a.T().Run("test get full stats", func(t *testing.T) {
		ctx := context.Background()
		response, err := a.Controller.GetStatistics(ctx, "", "")
		assert.NoError(t, err, "failed to get stats")
		assert.NotNil(t, string(response))
	})
	a.T().Run("test filtered full stats", func(t *testing.T) {
		ctx := context.Background()
		nodeID := a.GetNodesIDUsingRESTAPI(t)
		response, err := a.Controller.GetStatistics(ctx, nodeID, "graph_index_errors,knn_query_requests")
		assert.NoError(t, err, "failed to get stats")
		assert.NotNil(t, string(response))
		var data map[string]interface{}
		if err := json.Unmarshal(response, &data); err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, data)
		assert.NotNil(t, data["nodes"])
		nodes := data["nodes"].(map[string]interface{})
		if _, ok := nodes[nodeID]; !ok {
			t.Fatal("Node id is not found")
		}
		stats := nodes[nodeID].(map[string]interface{})
		if _, ok := stats["graph_index_errors"]; !ok {
			t.Fatal("graph_index_errors is not found")
		}
		if _, ok := stats["knn_query_requests"]; !ok {
			t.Fatal("knn_query_requests is not found")
		}
	})
	a.T().Run("test filtered nodes", func(t *testing.T) {
		ctx := context.Background()
		nodeID := a.GetNodesIDUsingRESTAPI(t)
		response, err := a.Controller.GetStatistics(ctx, nodeID, "")
		assert.NoError(t, err, "failed to get stats")
		assert.NotNil(t, string(response))
		var data map[string]interface{}
		if err := json.Unmarshal(response, &data); err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, data)
		assert.NotNil(t, data["nodes"])
		nodes := data["nodes"].(map[string]interface{})
		if _, ok := nodes[nodeID]; !ok {
			t.Fatal("Node id is not found")
		}
		stats := nodes[nodeID].(map[string]interface{})
		if _, ok := stats["graph_index_errors"]; !ok {
			t.Fatal("graph_index_errors is not found")
		}
		if _, ok := stats["knn_query_requests"]; !ok {
			t.Fatal("knn_query_requests is not found")
		}
	})
	a.T().Run("test filtered only stats", func(t *testing.T) {
		ctx := context.Background()
		response, err := a.Controller.GetStatistics(ctx, "", "graph_index_errors,knn_query_requests")
		assert.NoError(t, err, "failed to get stats")
		assert.NotNil(t, string(response))
		var data map[string]interface{}
		if err := json.Unmarshal(response, &data); err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, data)
		assert.NotNil(t, data["nodes"])
		nodes := data["nodes"].(map[string]interface{})
		nodeID := a.GetNodesIDUsingRESTAPI(t)
		if _, ok := nodes[nodeID]; !ok {
			t.Fatal("Node id is not found")
		}
		stats := nodes[nodeID].(map[string]interface{})
		if _, ok := stats["graph_index_errors"]; !ok {
			t.Fatal("graph_index_errors is not found")
		}
		if _, ok := stats["knn_query_requests"]; !ok {
			t.Fatal("knn_query_requests is not found")
		}
	})
}

func (a *KNNTestSuite) TestWarmupIndices() {
	if !a.IsPluginInstalled() {
		a.T().Skipf("plugin %s is not installed", a.Plugins)
	}
	a.T().Run("test warmup success", func(t *testing.T) {
		ctx := context.Background()
		response, err := a.Controller.WarmupIndices(ctx, []string{KNNSampleIndexFileName})
		assert.NoError(t, err, "failed to load graph into memory")
		assert.NotNil(t, response)
		assert.True(t, response.Total > 0)
		assert.EqualValues(t, response.Total, response.Successful)
	})

	a.T().Run("test warmup failure", func(t *testing.T) {
		ctx := context.Background()
		_, err := a.Controller.WarmupIndices(ctx, []string{"invalid-index-name"})
		assert.Error(t, err, "failed to load graph into memory")
		assert.EqualErrorf(t, err, "no such index [invalid-index-name]", "failed to parse error")
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestKNNSuite(t *testing.T) {
	suite.Run(t, new(KNNTestSuite))
}
