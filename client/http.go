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
