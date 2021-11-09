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

package gateway

import (
	"opensearch-cli/client/mocks"
	"opensearch-cli/entity"
	"opensearch-cli/environment"
	"opensearch-cli/mapper"
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
		_, err := NewHTTPGateway(testClient, &profile)
		assert.NoError(t, err)
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
		_, err := NewHTTPGateway(testClient, &profile)
		assert.NoError(t, err)
		assert.EqualValues(t, valAttempt, testClient.HTTPClient.RetryMax)
	})

	t.Run("override from environment variable", func(t *testing.T) {
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
		_, err := NewHTTPGateway(testClient, &profile)
		assert.NoError(t, err)
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
		_, err := NewHTTPGateway(testClient, &profile)
		assert.NoError(t, err)
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
		_, err := NewHTTPGateway(testClient, &profile)
		assert.NoError(t, err)
		assert.EqualValues(t, time.Duration(timeout)*time.Second, testClient.HTTPClient.HTTPClient.Timeout)
	})

	t.Run("override from environment variable", func(t *testing.T) {
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
		_, err := NewHTTPGateway(testClient, &profile)
		assert.NoError(t, err)
		assert.EqualValues(t, 5*time.Second, testClient.HTTPClient.HTTPClient.Timeout)
	})
}

func TestGatewayTLSConnection(t *testing.T) {

	t.Run("valid certificate path", func(t *testing.T) {
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
			Certificate: &entity.Trust{
				CAFilePath:                mapper.StringToStringPtr("testdata/ca.cert"),
				ClientCertificateFilePath: mapper.StringToStringPtr("testdata/client.cert"),
				ClientKeyFilePath:         mapper.StringToStringPtr("testdata/client.key"),
			},
		}
		testClient := mocks.NewTestClient(nil)
		val, err := NewHTTPGateway(testClient, &profile)
		assert.NoError(t, err)
		assert.NotNil(t, val)
	})
	t.Run("invalid CA certificate path", func(t *testing.T) {
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
			Certificate: &entity.Trust{
				CAFilePath:                mapper.StringToStringPtr("testdata/ca1.cert"),
				ClientCertificateFilePath: mapper.StringToStringPtr("testdata/client.cert"),
				ClientKeyFilePath:         mapper.StringToStringPtr("testdata/client.key"),
			},
		}
		testClient := mocks.NewTestClient(nil)
		_, err := NewHTTPGateway(testClient, &profile)
		assert.EqualError(t, err, "error opening certificate file testdata/ca1.cert, error: open testdata/ca1.cert: no such file or directory")
	})

	t.Run("invalid client certificate path", func(t *testing.T) {
		profile := entity.Profile{
			Name:     "test1",
			Endpoint: "https://localhost:9200",
			Certificate: &entity.Trust{
				CAFilePath:                mapper.StringToStringPtr("testdata/ca.cert"),
				ClientCertificateFilePath: mapper.StringToStringPtr("testdata/client1.cert"),
				ClientKeyFilePath:         mapper.StringToStringPtr("testdata/client.key"),
			},
		}
		testClient := mocks.NewTestClient(nil)
		_, err := NewHTTPGateway(testClient, &profile)
		assert.EqualError(t, err, "error creating x509 keypair from client cert file testdata/client1.cert and client key file testdata/client.key")
	})
}
