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
	"context"
	"encoding/json"
	entity "opensearch-cli/entity/knn"
	gateway "opensearch-cli/gateway/knn"
	"strings"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen  -destination=mocks/mock_knn.go -package=mocks . Controller

//Controller is an interface for the k-NN plugin controllers
type Controller interface {
	GetStatistics(context.Context, string, string) ([]byte, error)
	WarmupIndices(context.Context, []string) (*entity.Shards, error)
}

type controller struct {
	gateway gateway.Gateway
}

//GetStatistics gets stats data based on nodes and stat names
func (c controller) GetStatistics(ctx context.Context, nodes string, names string) ([]byte, error) {
	return c.gateway.GetStatistics(ctx, nodes, names)
}

//New returns new Controller instance
func New(gateway gateway.Gateway) Controller {
	return &controller{
		gateway,
	}
}

//WarmupIndices will load all the graphs for all of the shards (primaries and replicas)
//of all the indices specified in the request into native memory
func (c controller) WarmupIndices(ctx context.Context, index []string) (*entity.Shards, error) {
	indices := strings.Join(index, ",")
	response, err := c.gateway.WarmupIndices(ctx, indices)
	if err != nil {
		return nil, err
	}
	var warmupAPI entity.WarmupAPIResponse
	err = json.Unmarshal(response, &warmupAPI)
	if err != nil {
		return nil, err
	}
	return &warmupAPI.Shards, nil
}
