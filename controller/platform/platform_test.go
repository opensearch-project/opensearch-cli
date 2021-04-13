/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"io/ioutil"
	"net/http"
	"opensearch-cli/entity/platform"
	"opensearch-cli/gateway/platform/mocks"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func helperConvertToInterface(input []string) []interface{} {
	s := make([]interface{}, len(input))
	for i, v := range input {
		s[i] = v
	}
	return s
}

func TestController_GetDistinctValues(t *testing.T) {
	t.Run("empty index name", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		ctrl := New(mockGateway)
		_, err := ctrl.GetDistinctValues(ctx, "", "f1")
		assert.Error(t, err)
	})
	t.Run("empty field name", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		ctrl := New(mockGateway)
		_, err := ctrl.GetDistinctValues(ctx, "", "")
		assert.Error(t, err)
	})
	t.Run("gateway failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		mockGateway.EXPECT().SearchDistinctValues(ctx, "example", "f1").Return(nil, errors.New("search failed"))
		ctrl := New(mockGateway)
		_, err := ctrl.GetDistinctValues(ctx, "example", "f1")
		assert.Error(t, err)
	})
	t.Run("gateway response failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		mockGateway.EXPECT().SearchDistinctValues(ctx, "example", "f1").Return([]byte("No response"), nil)
		ctrl := New(mockGateway)
		_, err := ctrl.GetDistinctValues(ctx, "example", "f1")
		assert.Error(t, err)
	})
	t.Run("get distinct success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		expectedResult := helperConvertToInterface([]string{"Packaged Foods", "Dairy", "Meat and Seafood"})
		mockGateway.EXPECT().SearchDistinctValues(ctx, "example", "f1").Return(helperLoadBytes(t, "search_result.json"), nil)
		ctrl := New(mockGateway)
		result, err := ctrl.GetDistinctValues(ctx, "example", "f1")
		assert.NoError(t, err)
		assert.EqualValues(t, expectedResult, result)

	})
}

func TestController_Curl(t *testing.T) {
	commandRequest := platform.CurlCommandRequest{
		Action:      "post",
		Path:        "",
		QueryParams: "",
		Headers:     "",
		Data:        "",
		Pretty:      false,
	}

	request := platform.CurlRequest{
		Action:      http.MethodPost,
		Path:        "",
		QueryParams: "",
		Headers:     nil,
		Data:        nil,
	}
	t.Run("gateway success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		mockGateway.EXPECT().Curl(ctx, request).Return([]byte("response"), nil)
		ctrl := New(mockGateway)
		data, err := ctrl.Curl(ctx, commandRequest)
		assert.NoError(t, err, "received error")
		assert.EqualValues(t, []byte("response"), data)
	})
	t.Run("gateway response failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		mockGateway.EXPECT().Curl(ctx, request).Return(nil, errors.New("gateway failed"))
		ctrl := New(mockGateway)
		_, err := ctrl.Curl(ctx, commandRequest)
		assert.Error(t, err)
	})
	t.Run("mapper failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockGateway := mocks.NewMockGateway(mockCtrl)
		ctx := context.Background()
		//mockGateway.EXPECT().Curl(ctx, request).Return(nil, errors.New("gateway failed"))
		ctrl := New(mockGateway)
		_, err := ctrl.Curl(ctx, platform.CurlCommandRequest{})
		assert.EqualErrorf(t, err, "action cannot be empty", "wrong error message")
	})
}
