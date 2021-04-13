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

package signer

import (
	"net/http"
	"opensearch-cli/entity"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func buildSigner() *v4.Signer {
	return &v4.Signer{
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET", "SESSION"),
	}
}

func TestV4Signer(t *testing.T) {
	t.Run("sign request success", func(t *testing.T) {
		req, _ := retryablehttp.NewRequest(http.MethodGet, "https://localhost:9200", nil)
		region := os.Getenv("AWS_REGION")
		os.Setenv("AWS_REGION", "us-west-2")
		defer func() {
			os.Setenv("AWS_REGION", region)
		}()
		err := SignRequest(req, entity.AWSIAM{
			ProfileName: "test1",
			ServiceName: "es",
		}, func(c *credentials.Credentials) *v4.Signer {
			return buildSigner()
		})
		assert.NoError(t, err)
		q := req.Header
		assert.NotEmpty(t, q.Get("Authorization"))
		assert.NotEmpty(t, q.Get("X-Amz-Date"))
	})
	t.Run("sign request failed due to no region found", func(t *testing.T) {
		req, _ := retryablehttp.NewRequest(http.MethodGet, "https://localhost:9200", nil)
		region := os.Getenv("AWS_REGION")
		os.Setenv("AWS_REGION", "")
		defer func() {
			os.Setenv("AWS_REGION", region)
		}()
		err := SignRequest(req, entity.AWSIAM{
			ProfileName: "test1",
			ServiceName: "es",
		}, func(c *credentials.Credentials) *v4.Signer {
			return buildSigner()
		})
		assert.EqualErrorf(
			t, err, "aws region is not found. Either set 'AWS_REGION' or add this information during aws profile creation step", "unexpected error")
	})
}
