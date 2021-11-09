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

package knn

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"opensearch-cli/client"
	"opensearch-cli/client/mocks"
	"opensearch-cli/entity"
	"opensearch-cli/entity/knn"
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
			Header:  make(http.Header),
			Status:  "SOME OUTPUT",
			Request: req,
		}
	})
}

func TestGatewayGetStatistics(t *testing.T) {
	ctx := context.Background()
	t.Run("full stats succeeded", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_plugins/_knn/stats", 200, []byte("success"))
		testGateway, err := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		assert.NoError(t, err)
		actual, err := testGateway.GetStatistics(ctx, "", "")
		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "success")
	})
	t.Run("filtered node and stats succeeded", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_plugins/_knn/node1,node2/stats/stat1", 200, []byte("success"))
		testGateway, err := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		assert.NoError(t, err)
		actual, err := testGateway.GetStatistics(ctx, "node1,node2", "stat1")
		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "success")
	})
	t.Run("filtered node succeeded", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_plugins/_knn/node1,node2/stats/", 200, []byte("success"))
		testGateway, err := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		assert.NoError(t, err)
		actual, err := testGateway.GetStatistics(ctx, "node1,node2", "")
		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "success")
	})
	t.Run("filtered stats succeeded", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_plugins/_knn//stats/stat1,stat2", 200, []byte("success"))
		testGateway, err := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		assert.NoError(t, err)
		actual, err := testGateway.GetStatistics(ctx, "", "stat1,stat2")
		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "success")
	})
	t.Run("gateway failed due to gateway user config", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_plugins/_knn/stats", 400, []byte("failed"))
		testGateway, err := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		assert.NoError(t, err)
		_, err = testGateway.GetStatistics(ctx, "", "")
		assert.Error(t, err)
	})
	t.Run("failed due to invalid stat names", func(t *testing.T) {
		reason := "request [/_plugins/_knn//stats/graph_count] contains unrecognized stat: [stat1]"
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
		testClient := getTestClient(t, "http://localhost:9200/_plugins/_knn/index1/stats/invalid-stats", 404, response)
		testGateway, err := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		assert.NoError(t, err)
		_, err = testGateway.GetStatistics(ctx, "index1", "invalid-stats")
		assert.EqualErrorf(t, err, reason, "failed to parse error")
	})
}

func TestGatewayWarmupIndices(t *testing.T) {
	ctx := context.Background()
	t.Run("warmup indices", func(t *testing.T) {

		testClient := getTestClient(t, "http://localhost:9200/_plugins/_knn/warmup/index1,index2", 200, []byte("success"))
		testGateway, err := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		assert.NoError(t, err)
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
		testClient := getTestClient(t, "http://localhost:9200/_plugins/_knn/warmup/index1", 404, response)
		testGateway, err := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		assert.NoError(t, err)
		_, err = testGateway.WarmupIndices(ctx, "index1")
		assert.EqualErrorf(t, err, "no such index", "failed to parse error")
	})
}
