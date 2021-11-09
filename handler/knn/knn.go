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
	"opensearch-cli/controller/knn"
	entity "opensearch-cli/entity/knn"
)

//Handler is facade for controller
type Handler struct {
	knn.Controller
}

// New returns new Handler instance
func New(controller knn.Controller) *Handler {
	return &Handler{
		controller,
	}
}

//GetStatistics gets stats data based on nodes and stat names
func GetStatistics(h *Handler, nodes string, names string) ([]byte, error) {
	return h.GetStatistics(nodes, names)
}

//GetStatistics gets stats data based on nodes and stat names
func (h *Handler) GetStatistics(nodes string, names string) ([]byte, error) {
	ctx := context.Background()
	response, err := h.Controller.GetStatistics(ctx, nodes, names)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}

	if err := json.Unmarshal(response, &data); err != nil {
		return nil, err
	}
	return json.MarshalIndent(data, "", "  ")
}

//WarmupIndices warmups knn index
func WarmupIndices(h *Handler, index []string) (*entity.Shards, error) {
	return h.WarmupIndices(index)
}

//WarmupIndices warmups shard based on knn index and returns status of shards
func (h *Handler) WarmupIndices(index []string) (*entity.Shards, error) {
	ctx := context.Background()
	return h.Controller.WarmupIndices(ctx, index)
}
