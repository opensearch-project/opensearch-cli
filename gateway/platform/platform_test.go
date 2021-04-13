/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"opensearch-cli/client"
	"opensearch-cli/client/mocks"
	"opensearch-cli/entity"
	"opensearch-cli/entity/platform"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func getTestClient(t *testing.T, responseData string, code int) *client.Client {
	return mocks.NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "http://localhost:9200/test_index/_search")
		resBytes, _ := ioutil.ReadAll(req.Body)
		var body platform.SearchRequest
		err := json.Unmarshal(resBytes, &body)
		assert.NoError(t, err)
		assert.EqualValues(t, body.Size, 0)
		assert.EqualValues(t, body.Agg.Group.Term.Field, "day_of_week")
		assert.EqualValues(t, len(req.Header), 2)
		return &http.Response{
			StatusCode: code,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(responseData)),
			// Must be set to non-nil value or it panics
			Header:  make(http.Header),
			Status:  "SOME OUTPUT",
			Request: req,
		}
	})
}

func getCurlTestClient(t *testing.T, expectedURL string, expectedData []byte, expectedHeader map[string]string, responseData string, code int) *client.Client {
	return mocks.NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, expectedURL, req.URL.String())
		resBytes, _ := ioutil.ReadAll(req.Body)
		assert.EqualValues(t, expectedData, resBytes)

		for k, v := range expectedHeader {
			assert.EqualValues(t, v, req.Header.Get(k))
		}
		return &http.Response{
			StatusCode: code,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(responseData)),
			// Must be set to non-nil value or it panics
			Header:  make(http.Header),
			Status:  "SOME OUTPUT",
			Request: req,
		}
	})
}

func TestGateway_SearchDistinctValues(t *testing.T) {
	responseData, _ := json.Marshal(helperLoadBytes(t, "search_result.json"))
	ctx := context.Background()
	t.Run("search succeeded", func(t *testing.T) {

		testClient := getTestClient(t, string(responseData), 200)
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		actual, err := testGateway.SearchDistinctValues(ctx, "test_index", "day_of_week")
		assert.NoError(t, err)
		assert.EqualValues(t, actual, responseData)
	})
	t.Run("search failed due to 404", func(t *testing.T) {
		testClient := getTestClient(t, "No connection found", 404)
		testGateway := New(testClient, &entity.Profile{
			Endpoint: "http://localhost:9200",
			UserName: "admin",
			Password: "admin",
		})
		_, err := testGateway.SearchDistinctValues(ctx, "test_index", "day_of_week")
		assert.EqualError(t, err, "No connection found")
	})
}

func getErrorResponse() []byte {
	return []byte(`{
  "error" : {
    "root_cause" : [ {
      "type" : "some_exception",
      "reason" : "Failed to execute"
    } ],
    "type" : "some_exception",
    "reason" : "Failed to execute"
  },
  "status" : 400
}`)
}

func TestGatewayCurl(t *testing.T) {
	ctx := context.Background()
	p := &entity.Profile{
		Endpoint: "http://localhost:9200",
		UserName: "admin",
		Password: "admin",
	}
	t.Run("curl succeeded with empty data, headers, params", func(t *testing.T) {
		expectedData := []byte(``)
		expectedHeader := map[string]string{}
		expectedResponse := "OK"
		testClient := getCurlTestClient(t, "http://localhost:9200/_cluster/health", []byte(``), map[string]string{}, expectedResponse, 200)
		testGateway := New(testClient, p)
		actual, err := testGateway.Curl(ctx, platform.CurlRequest{
			Action:      http.MethodGet,
			Path:        "_cluster/health",
			QueryParams: "",
			Headers:     expectedHeader,
			Data:        expectedData,
		})

		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), expectedResponse)
	})
	t.Run("curl succeeded with empty data, headers", func(t *testing.T) {

		expectedData := []byte(``)
		expectedHeader := map[string]string{}
		testClient := getCurlTestClient(t, "http://localhost:9200/_cluster/health?params=true&v=true", expectedData, expectedHeader, "OK", 200)
		testGateway := New(testClient, p)
		actual, err := testGateway.Curl(ctx, platform.CurlRequest{
			Action:      http.MethodGet,
			Path:        "_cluster/health",
			QueryParams: "params=true&v=true",
			Headers:     expectedHeader,
			Data:        expectedData,
		})

		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "OK")
	})

	t.Run("curl succeeded", func(t *testing.T) {

		expectedData := []byte(`{"data": 1}`)
		expectedHeader := map[string]string{
			"one":          "1",
			"two":          "2",
			"content-type": "gzip",
		}
		testClient := getCurlTestClient(t, "http://localhost:9200/_cluster/health?params=true&v=true", expectedData, expectedHeader, "OK", 200)
		testGateway := New(testClient, p)
		actual, err := testGateway.Curl(ctx, platform.CurlRequest{
			Action:      http.MethodGet,
			Path:        "_cluster/health",
			QueryParams: "params=true&v=true",
			Headers:     expectedHeader,
			Data:        expectedData,
		})

		assert.NoError(t, err)
		assert.EqualValues(t, string(actual), "OK")
	})
	t.Run("curl failed due to client error", func(t *testing.T) {
		expectedData := []byte(`{"data": 1}`)
		expectedHeader := map[string]string{
			"one":          "1",
			"two":          "2",
			"content-type": "gzip",
		}
		responseData := getErrorResponse()
		testClient := getCurlTestClient(t, "http://localhost:9200/_cluster/health?params=true&v=true", expectedData, expectedHeader, string(responseData), 400)
		testGateway := New(testClient, p)
		_, err := testGateway.Curl(ctx, platform.CurlRequest{
			Action:      http.MethodGet,
			Path:        "_cluster/health",
			QueryParams: "params=true&v=true",
			Headers:     expectedHeader,
			Data:        expectedData,
		})
		assert.EqualErrorf(t, err, "400 Client Error: SOME OUTPUT for url: http://localhost:9200/_cluster/health?params=true&v=true", "failed to receive expected error")
		assert.IsType(t, &platform.RequestError{}, err, "failed to type cast error")
		requestError, _ := err.(*platform.RequestError)
		assert.True(t, len(requestError.GetResponse()) > 0)
		assert.EqualValues(t, 400, requestError.StatusCode())
	})

	t.Run("curl failed due to server error", func(t *testing.T) {
		expectedData := []byte(`{"data": 1}`)
		expectedHeader := map[string]string{
			"one":          "1",
			"two":          "2",
			"content-type": "gzip",
		}
		responseData := getErrorResponse()
		testClient := getCurlTestClient(t, "http://localhost:9200/_cluster/health?params=true&v=true", expectedData, expectedHeader, string(responseData), 501)
		testGateway := New(testClient, p)
		_, err := testGateway.Curl(ctx, platform.CurlRequest{
			Action:      http.MethodGet,
			Path:        "_cluster/health",
			QueryParams: "params=true&v=true",
			Headers:     expectedHeader,
			Data:        expectedData,
		})
		assert.EqualErrorf(t, err, "501 Server Error: SOME OUTPUT for url: http://localhost:9200/_cluster/health?params=true&v=true", "failed to receive expected error")
		assert.IsType(t, &platform.RequestError{}, err, "failed to type cast error")
		requestError, _ := err.(*platform.RequestError)
		assert.True(t, len(requestError.GetResponse()) > 0)
		assert.EqualValues(t, 501, requestError.StatusCode())
	})
}
