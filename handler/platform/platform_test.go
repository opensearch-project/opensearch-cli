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
	"errors"
	"opensearch-cli/controller/platform/mocks"
	entity "opensearch-cli/entity/platform"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandlerCurl(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	arg := entity.CurlCommandRequest{}
	t.Run("success", func(t *testing.T) {
		mockedController := mocks.NewMockController(mockCtrl)
		mockedController.EXPECT().Curl(ctx, arg).Return([]byte(`{"result" : "success"}`), nil)
		instance := New(mockedController)
		response, err := Curl(instance, arg)
		assert.NoError(t, err)
		assert.EqualValues(t, "{\"result\" : \"success\"}", string(response))
	})
	t.Run("failed to execute", func(t *testing.T) {
		mockedController := mocks.NewMockController(mockCtrl)
		mockedController.EXPECT().Curl(ctx, arg).Return(nil, errors.New("failed to execute"))
		instance := New(mockedController)
		_, err := instance.Curl(arg)
		assert.EqualError(t, err, "failed to execute")
	})
}
