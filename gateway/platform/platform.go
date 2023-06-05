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

package platform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"opensearch-cli/client"
	"opensearch-cli/entity"
	"opensearch-cli/entity/platform"
	gw "opensearch-cli/gateway"

	"github.com/hashicorp/go-retryablehttp"
)

const search = "_search"

//go:generate go run -mod=mod github.com/golang/mock/mockgen  -destination=mocks/mock_platform.go -package=mocks . Gateway

// Gateway interface to call OpenSearch
type Gateway interface {
	SearchDistinctValues(ctx context.Context, index string, field string) ([]byte, error)
	Curl(ctx context.Context, request platform.CurlRequest) ([]byte, error)
}

type gateway struct {
	gw.HTTPGateway
}

// New returns new Gateway instance
func New(c *client.Client, p *entity.Profile) (Gateway, error) {
	g, err := gw.NewHTTPGateway(c, p)
	if err != nil {
		return nil, err
	}
	return &gateway{*g}, nil
}
func buildPayload(field string) *platform.SearchRequest {
	return &platform.SearchRequest{
		Size: 0, // This will skip data in the response
		Agg: platform.Aggregate{
			Group: platform.DistinctGroups{
				Term: platform.Terms{
					Field: field,
				},
			},
		},
	}
}

func (g *gateway) buildSearchURL(index string) (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	endpoint.Path = fmt.Sprintf("%s/%s", index, search)
	return endpoint, nil
}

// SearchDistinctValues gets distinct values on index for given field
func (g *gateway) SearchDistinctValues(ctx context.Context, index string, field string) ([]byte, error) {
	searchURL, err := g.buildSearchURL(index)
	if err != nil {
		return nil, err
	}
	searchRequest, err := g.BuildRequest(ctx, http.MethodGet, buildPayload(field), searchURL.String(), gw.GetDefaultHeaders())
	if err != nil {
		return nil, err
	}
	response, err := g.Call(searchRequest, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Curl executes REST request based on request parameters
func (g *gateway) Curl(ctx context.Context, request platform.CurlRequest) ([]byte, error) {
	requestURL, err := g.buildURL(request)
	if err != nil {
		return nil, err
	}
	//append request headers with gateway default headers
	headers := gw.GetDefaultHeaders()
	for k, v := range request.Headers {
		headers[k] = v
	}

	var curlRequest *retryablehttp.Request
	var buildErr error

	// when formDataFile is provided, build multipart/form-data request
	if len(request.FormDataFile) > 0 {
		curlRequest, buildErr = g.BuildCurlMultipartFormRequest(ctx, request.Action, request.FormDataFile, requestURL.String(), request.Headers)
	} else {
		// else build "normal" rest request
		curlRequest, buildErr = g.BuildCurlRequest(ctx, request.Action, request.Data, requestURL.String(), headers)
	}

	if buildErr != nil {
		return nil, buildErr
	}
	response, err := g.Execute(curlRequest)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (g *gateway) buildURL(request platform.CurlRequest) (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	endpoint.Path = request.Path
	endpoint.RawQuery = request.QueryParams
	return endpoint, nil
}
