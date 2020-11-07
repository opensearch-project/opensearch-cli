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

package es

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"odfe-cli/client"
	"odfe-cli/entity"
	"odfe-cli/entity/es"
	gw "odfe-cli/gateway"
)

const search = "_search"

//go:generate go run -mod=mod github.com/golang/mock/mockgen  -destination=mocks/mock_es.go -package=mocks . Gateway

//Gateway interface to call ES
type Gateway interface {
	SearchDistinctValues(ctx context.Context, index string, field string) ([]byte, error)
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
func buildPayload(field string) *es.SearchRequest {
	return &es.SearchRequest{
		Size: 0, // This will skip data in the response
		Agg: es.Aggregate{
			Group: es.DistinctGroups{
				Term: es.Terms{
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
	searchRequest, err := g.BuildRequest(ctx, http.MethodGet, buildPayload(field), searchURL.String(), gw.GetHeaders())
	if err != nil {
		return nil, err
	}
	response, err := g.Call(searchRequest, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return response, nil
}
