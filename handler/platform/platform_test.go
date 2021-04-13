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
