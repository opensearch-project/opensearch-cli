/*
 * Copyright 2021 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package knn

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"odfe-cli/client"
	"odfe-cli/client/mocks"
	"odfe-cli/entity"
	"odfe-cli/entity/knn"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestClient(t *testing.T, url string, code int, response []byte) *client.Client {
	return mocks.NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), url)
		assert.EqualValues(t, len(req.Header), 2)
		return &http.Response{
			StatusCode: code,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBuffer(response)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
}

func TestGatewayGetStatistics(t *testing.T) {
	ctx := context.Background()
	t.Run("full stats succeeded", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_opendistro/_knn/stats", 200, []byte("success"))
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		actual, err := testGateway.GetStatistics(ctx, "", "")
		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "success")
	})
	t.Run("filtered node and stats succeeded", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_opendistro/_knn/node1,node2/stats/stat1", 200, []byte("success"))
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		actual, err := testGateway.GetStatistics(ctx, "node1,node2", "stat1")
		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "success")
	})
	t.Run("filtered node succeeded", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_opendistro/_knn/node1,node2/stats/", 200, []byte("success"))
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		actual, err := testGateway.GetStatistics(ctx, "node1,node2", "")
		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "success")
	})
	t.Run("filtered stats succeeded", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_opendistro/_knn//stats/stat1,stat2", 200, []byte("success"))
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		actual, err := testGateway.GetStatistics(ctx, "", "stat1,stat2")
		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "success")
	})
	t.Run("gateway failed due to gateway user config", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_opendistro/_knn/stats", 400, []byte("failed"))
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		_, err := testGateway.GetStatistics(ctx, "", "")
		assert.Error(t, err)
	})
	t.Run("failed due to invalid stat names", func(t *testing.T) {
		reason := "request [/_opendistro/_knn//stats/graph_count] contains unrecognized stat: [stat1]"
		response, _ := json.Marshal(knn.ErrorResponse{
			KNNError: knn.Error{
				RootCause: []knn.RootCause{
					{
						Type:   "stat_not_found_exception",
						Reason: reason,
					},
				},
			},
			Status: 404,
		})
		testClient := getTestClient(t, "http://localhost:9200/_opendistro/_knn/index1/stats/invalid-stats", 404, response)
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		_, err := testGateway.GetStatistics(ctx, "index1", "invalid-stats")
		assert.EqualErrorf(t, err, reason, "failed to parse error")
	})
}

func TestGatewayWarmupIndices(t *testing.T) {
	ctx := context.Background()
	t.Run("warmup indices", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_opendistro/_knn/warmup/index1,index2", 200, []byte("success"))
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		actual, err := testGateway.WarmupIndices(ctx, "index1,index2")
		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "success")
	})
	t.Run("failed due to invalid index", func(t *testing.T) {

		response, _ := json.Marshal(knn.ErrorResponse{
			KNNError: knn.Error{
				RootCause: []knn.RootCause{
					{
						Type:   "index_not_found_exception",
						Reason: "no such index",
					},
				},
			},
			Status: 404,
		})
		testClient := getTestClient(t, "http://localhost:9200/_opendistro/_knn/warmup/index1", 404, response)
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		_, err := testGateway.WarmupIndices(ctx, "index1")
		assert.EqualErrorf(t, err, "no such index", "failed to parse error")
	})
}
