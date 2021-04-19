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

package client

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const defaultTimeout = 10

//Client is an Abstraction for actual client
type Client struct {
	HTTPClient *retryablehttp.Client
}

//NewDefaultClient return new instance of client
func NewDefaultClient(tripper http.RoundTripper) (*Client, error) {

	client := retryablehttp.NewClient()
	client.HTTPClient.Transport = tripper
	client.HTTPClient.Timeout = defaultTimeout * time.Second
	client.Logger = nil
	return &Client{
		HTTPClient: client,
	}, nil
}

//New takes transport and uses accordingly
func New(tripper http.RoundTripper) (*Client, error) {
	if tripper == nil {
		tripper = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return NewDefaultClient(tripper)
}
