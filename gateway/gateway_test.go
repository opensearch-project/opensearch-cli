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

package gateway

import (
	"opensearch-cli/client/mocks"
	"opensearch-cli/entity"
	"opensearch-cli/environment"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetValidEndpoint(t *testing.T) {
	t.Run("valid endpoint", func(t *testing.T) {

		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
			UserName: "foo",
			Password: "bar",
		}
		url, err := GetValidEndpoint(&profile)
		assert.NoError(t, err)
		assert.EqualValues(t, "https://localhost:9200", url.String())
	})
	t.Run("empty endpoint", func(t *testing.T) {
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "",
			UserName: "foo",
			Password: "bar",
		}
		_, err := GetValidEndpoint(&profile)
		assert.EqualErrorf(t, err, "invalid endpoint:  due to parse \"\": empty url", "failed to get expected error")
	})
}

func TestGatewayRetryVal(t *testing.T) {
	t.Run("default retry max value", func(t *testing.T) {
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
		}
		testClient := mocks.NewTestClient(nil)
		NewHTTPGateway(testClient, &profile)
		assert.EqualValues(t, 4, testClient.HTTPClient.RetryMax)
	})
	t.Run("profile retry max value", func(t *testing.T) {
		valAttempt := 2
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
			MaxRetry: &valAttempt,
		}
		testClient := mocks.NewTestClient(nil)
		NewHTTPGateway(testClient, &profile)
		assert.EqualValues(t, valAttempt, testClient.HTTPClient.RetryMax)
	})

	t.Run("override from os variable", func(t *testing.T) {
		val := os.Getenv(environment.OPENSEARCH_MAX_RETRY)
		defer func() {
			assert.NoError(t, os.Setenv(environment.OPENSEARCH_MAX_RETRY, val))
		}()
		os.Setenv(environment.OPENSEARCH_MAX_RETRY, "10")
		valAttempt := 2
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
			MaxRetry: &valAttempt,
		}
		testClient := mocks.NewTestClient(nil)
		NewHTTPGateway(testClient, &profile)
		assert.EqualValues(t, 10, testClient.HTTPClient.RetryMax)
	})
}

func TestGatewayConnectionTimeout(t *testing.T) {
	t.Run("default timeout", func(t *testing.T) {
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
		}
		testClient := mocks.NewTestClient(nil)
		NewHTTPGateway(testClient, &profile)
		assert.EqualValues(t, 10*time.Second, testClient.HTTPClient.HTTPClient.Timeout)
	})
	t.Run("configure profile timeout", func(t *testing.T) {
		timeout := int64(60)
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
			Timeout:  &timeout,
		}
		testClient := mocks.NewTestClient(nil)
		NewHTTPGateway(testClient, &profile)
		assert.EqualValues(t, time.Duration(timeout)*time.Second, testClient.HTTPClient.HTTPClient.Timeout)
	})

	t.Run("override from os variable", func(t *testing.T) {
		val := os.Getenv(environment.OPENSEARCH_TIMEOUT)
		defer func() {
			assert.NoError(t, os.Setenv(environment.OPENSEARCH_TIMEOUT, val))
		}()
		os.Setenv(environment.OPENSEARCH_TIMEOUT, "5")
		timeout := int64(60)
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
			Timeout:  &timeout,
		}
		testClient := mocks.NewTestClient(nil)
		NewHTTPGateway(testClient, &profile)
		assert.EqualValues(t, 5*time.Second, testClient.HTTPClient.HTTPClient.Timeout)
	})
}
