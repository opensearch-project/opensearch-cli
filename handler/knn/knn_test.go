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
	"errors"
	"odfe-cli/controller/knn/mocks"
	entity "odfe-cli/entity/knn"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetStatistics(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	t.Run("get stats success", func(t *testing.T) {
		mockedController := mocks.NewMockController(mockCtrl)
		mockedController.EXPECT().GetStatistics(ctx, "node1", "stats-name").Return([]byte("{}"), nil)
		instance := New(mockedController)
		response, err := GetStatistics(instance, "node1", "stats-name")
		assert.NoError(t, err)
		assert.EqualValues(t, "{}", string(response))
	})
	t.Run("get stats failure", func(t *testing.T) {
		mockedController := mocks.NewMockController(mockCtrl)
		mockedController.EXPECT().GetStatistics(ctx, "node1", "stats-name").Return(nil, errors.New("failed to fetch data"))
		instance := New(mockedController)
		_, err := instance.GetStatistics("node1", "stats-name")
		assert.EqualError(t, err, "failed to fetch data")
	})
}

func TestHandlerWarmupIndices(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	t.Run("warmup success", func(t *testing.T) {
		mockedController := mocks.NewMockController(mockCtrl)
		result := &entity.Shards{
			Total:      10,
			Successful: 5,
			Failed:     5,
		}
		mockedController.EXPECT().WarmupIndices(ctx, []string{"index1"}).Return(result, nil)
		instance := New(mockedController)
		response, err := WarmupIndices(instance, []string{"index1"})
		assert.NoError(t, err)
		assert.EqualValues(t, *result, *response)
	})
	t.Run("warmup failure", func(t *testing.T) {
		mockedController := mocks.NewMockController(mockCtrl)
		mockedController.EXPECT().WarmupIndices(ctx, []string{"index1"}).Return(nil, errors.New("failed"))
		instance := New(mockedController)
		_, err := instance.WarmupIndices([]string{"index1"})
		assert.EqualError(t, err, "failed")
	})
}
