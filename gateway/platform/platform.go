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
)

const search = "_search"

//go:generate go run -mod=mod github.com/golang/mock/mockgen  -destination=mocks/mock_platform.go -package=mocks . Gateway

//Gateway interface to call OpenSearch
type Gateway interface {
	SearchDistinctValues(ctx context.Context, index string, field string) ([]byte, error)
	Curl(ctx context.Context, request platform.CurlRequest) ([]byte, error)
}

type gateway struct {
	gw.HTTPGateway
}

// New returns new Gateway instance
func New(c *client.Client, p *entity.Profile) Gateway {
	return &gateway{
		*gw.NewHTTPGateway(c, p),
	}
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

//SearchDistinctValues gets distinct values on index for given field
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

//Curl executes REST request based on request parameters
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
	curlRequest, err := g.BuildCurlRequest(ctx, request.Action, request.Data, requestURL.String(), headers)
	if err != nil {
		return nil, err
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
