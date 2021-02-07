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
	"errors"
	"odfe-cli/entity/knn"
	"odfe-cli/gateway/knn/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestControllerGetStatistics(t *testing.T) {
	t.Run("gateway failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		mockGateway.EXPECT().GetStatistics(ctx, "", "").Return(nil, errors.New("gateway failed"))
		ctrl := New(mockGateway)
		_, err := ctrl.GetStatistics(ctx, "", "")
		assert.Error(t, err)
	})
	t.Run("get stats success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		mockGateway.EXPECT().GetStatistics(ctx, "node1", "stats").Return([]byte(`response succeeded`), nil)
		ctrl := New(mockGateway)
		result, err := ctrl.GetStatistics(ctx, "node1", "stats")
		assert.NoError(t, err)
		assert.EqualValues(t, []byte(`response succeeded`), result)
	})
}

func TestControllerWarmupIndices(t *testing.T) {
	t.Run("gateway failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		mockGateway.EXPECT().WarmupIndices(ctx, "index1").Return(nil, errors.New("gateway failed"))
		ctrl := New(mockGateway)
		_, err := ctrl.WarmupIndices(ctx, []string{"index1"})
		assert.Error(t, err)
	})
	t.Run("warmup success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		expectedResponse := knn.WarmupAPIResponse{
			Shards: knn.Shards{
				Total:      10,
				Successful: 8,
				Failed:     2,
			},
		}
		rawMessage, err := json.Marshal(expectedResponse)
		assert.NoError(t, err)
		mockGateway.EXPECT().WarmupIndices(ctx, "index1").Return(rawMessage, nil)
		ctrl := New(mockGateway)
		result, err := ctrl.WarmupIndices(ctx, []string{"index1"})
		assert.NoError(t, err)
		assert.EqualValues(t, expectedResponse.Shards, *result)
	})
}
