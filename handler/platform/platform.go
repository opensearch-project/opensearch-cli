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
