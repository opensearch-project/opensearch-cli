// +build integration

/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package it

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"odfe-cli/client"
	ctrl "odfe-cli/controller/knn"
	"odfe-cli/entity"
	gateway "odfe-cli/gateway/knn"
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
	ODFECLISuite
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
		Endpoint: os.Getenv("ODFE_ENDPOINT"),
		UserName: os.Getenv("ODFE_USER"),
		Password: os.Getenv("ODFE_PASSWORD"),
	}
	if err = a.ValidateProfile(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.Gateway = gateway.New(a.Client, a.Profile)
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
