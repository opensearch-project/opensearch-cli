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

package platform

import (
	"context"
	"opensearch-cli/controller/platform"
	entity "opensearch-cli/entity/platform"
)

//Handler is facade for controller
type Handler struct {
	platform.Controller
}

// New returns new Handler instance
func New(controller platform.Controller) *Handler {
	return &Handler{
		controller,
	}
}

//Curl executes REST API as defined by curl command
func Curl(h *Handler, request entity.CurlCommandRequest) ([]byte, error) {
	return h.Curl(request)
}

//Curl executes REST API as defined by curl command
func (h *Handler) Curl(request entity.CurlCommandRequest) ([]byte, error) {
	ctx := context.Background()
	return h.Controller.Curl(ctx, request)
}
