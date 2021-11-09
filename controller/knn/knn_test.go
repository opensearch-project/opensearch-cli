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
	"errors"
	entity "opensearch-cli/entity/knn"
	gateway "opensearch-cli/gateway/knn/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestControllerGetStatistics(t *testing.T) {
	t.Run("gateway failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockGateway := gateway.NewMockGateway(mockCtrl)
		ctx := context.Background()
		mockGateway.EXPECT().GetStatistics(ctx, "", "").Return(nil, errors.New("gateway failed"))
		ctrl := New(mockGateway)
		_, err := ctrl.GetStatistics(ctx, "", "")
		assert.Error(t, err)
	})
	t.Run("get stats success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGateway := gateway.NewMockGateway(mockCtrl)
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

		mockGateway := gateway.NewMockGateway(mockCtrl)
		ctx := context.Background()
		mockGateway.EXPECT().WarmupIndices(ctx, "index1").Return(nil, errors.New("gateway failed"))
		ctrl := New(mockGateway)
		_, err := ctrl.WarmupIndices(ctx, []string{"index1"})
		assert.Error(t, err)
	})
	t.Run("warmup success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGateway := gateway.NewMockGateway(mockCtrl)
		ctx := context.Background()
		expectedResponse := entity.WarmupAPIResponse{
			Shards: entity.Shards{
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
