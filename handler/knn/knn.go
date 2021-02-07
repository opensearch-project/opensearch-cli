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

package knn

import (
	"context"
	"encoding/json"
	"odfe-cli/controller/knn"
	entity "odfe-cli/entity/knn"
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
