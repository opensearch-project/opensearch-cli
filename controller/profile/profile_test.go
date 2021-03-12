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

package profile

import (
	"errors"
	config "odfe-cli/controller/config/mocks"
	"odfe-cli/entity"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func getDefaultConfig() entity.Config {
	return entity.Config{
		Profiles: []entity.Profile{
			{
				Name:     odfeDefaultProfileName,
				Endpoint: "https://localhost:9200",
				UserName: "user", Password: "user123",
			}}}
}

func getSampleConfig() entity.Config {
	return entity.Config{
		Profiles: []entity.Profile{
			{
				Name:     "local",
				Endpoint: "https://localhost:9200",
				UserName: "admin", Password: "admin",
			},
			{
				Name:     "default",
				Endpoint: "https://127.0.0.1:9200",
				UserName: "user", Password: "user123",
			},
		}}
}

func TestControllerGetProfilesMap(t *testing.T) {

	t.Run("get profiles as map", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(getSampleConfig(), nil)
		ctrl := New(mockConfigCtrl)
		actual, err := ctrl.GetProfilesMap()
		assert.NoError(t, err)
		expected := map[string]entity.Profile{}
		for _, p := range getSampleConfig().Profiles {
			expected[p.Name] = p
		}
		assert.NoError(t, err)
		assert.EqualValues(t, expected, actual)
	})
	t.Run("config controller failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(entity.Config{}, errors.New("failed to read"))
		ctrl := New(mockConfigCtrl)
		_, err := ctrl.GetProfilesMap()
		assert.EqualError(t, err, "failed to read")
	})
}

func TestControllerGetProfiles(t *testing.T) {

	profiles := entity.Config{
		Profiles: []entity.Profile{
			{
				Name:     "local",
				Endpoint: "https://localhost:9200",
				UserName: "", Password: "",
			},
			{
				Name:     "default",
				Endpoint: "https://127.0.0.1:9200",
				UserName: "user", Password: "user123",
			},
			{
				Name:     "default1",
				Endpoint: "https://127.0.0.1:9200",
				UserName: "user", Password: "user123",
				AWS: &entity.AWSIAM{ProfileName: "iam"},
			},
			{
				Name:     "default2",
				Endpoint: "https://127.0.0.1:9200",
				UserName: "user", Password: "user123",
				AWS: &entity.AWSIAM{ProfileName: ""},
			},
		}}
	t.Run("get profiles", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(profiles, nil)
		ctrl := New(mockConfigCtrl)
		actual, err := ctrl.GetProfiles()
		assert.NoError(t, err)
		expectedProfiles := []entity.Profile{
			{
				Name:     "local",
				Endpoint: "https://localhost:9200",
				UserName: "", Password: "",
			},
			{
				Name:     "default",
				Endpoint: "https://127.0.0.1:9200",
				UserName: "user", Password: "user123",
			},
			{
				Name:     "default1",
				Endpoint: "https://127.0.0.1:9200",
				UserName: "user", Password: "user123",
				AWS: &entity.AWSIAM{ProfileName: "iam"},
			},
			{
				Name:     "default2",
				Endpoint: "https://127.0.0.1:9200",
				UserName: "user", Password: "user123",
				AWS: &entity.AWSIAM{ProfileName: ""},
			},
		}

		assert.EqualValues(t, expectedProfiles, actual)
	})
}

func TestControllerGetProfileNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(getSampleConfig(), nil)
		ctrl := New(mockConfigCtrl)
		names, err := ctrl.GetProfileNames()
		assert.NoError(t, err)
		assert.EqualValues(t, []string{"local", "default"}, names)
	})
	t.Run("config controller failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(entity.Config{}, errors.New("failed to read"))
		ctrl := New(mockConfigCtrl)
		_, err := ctrl.GetProfileNames()
		assert.EqualError(t, err, "failed to read")
	})
}

func TestControllerGetProfileForExecution(t *testing.T) {
	t.Run("provided profile name: success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(getSampleConfig(), nil)
		ctrl := New(mockConfigCtrl)
		p, ok, err := ctrl.GetProfileForExecution("local")
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.EqualValues(t, getSampleConfig().Profiles[0], p)
	})

	t.Run("select default profile", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		oldValue, ok := os.LookupEnv(odfeProfileEnvVarName)
		if ok {
			assert.NoError(t, os.Unsetenv(odfeProfileEnvVarName))
			defer func() {
				assert.NoError(t, os.Setenv(odfeDefaultProfileName, oldValue))
			}()
		}
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(getDefaultConfig(), nil)
		ctrl := New(mockConfigCtrl)
		p, ok, err := ctrl.GetProfileForExecution("")
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.EqualValues(t, getDefaultConfig().Profiles[0], p)
	})
	t.Run("set environment variable: success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(getSampleConfig(), nil)
		ctrl := New(mockConfigCtrl)
		oldValue, ok := os.LookupEnv(odfeProfileEnvVarName)
		if ok {
			assert.NoError(t, os.Unsetenv(odfeProfileEnvVarName))
			defer func() {
				assert.NoError(t, os.Setenv(odfeDefaultProfileName, oldValue))
			}()
		}
		err := os.Setenv(odfeProfileEnvVarName, "local")
		assert.NoError(t, err)
		defer func() {
			err = os.Unsetenv(odfeProfileEnvVarName)
			assert.NoError(t, err)
		}()
		p, ok, err := ctrl.GetProfileForExecution("")
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.EqualValues(t, getSampleConfig().Profiles[0], p)
	})
	t.Run("config controller failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(entity.Config{}, errors.New("failed to read"))
		ctrl := New(mockConfigCtrl)
		_, _, err := ctrl.GetProfileForExecution("local")
		assert.EqualError(t, err, "failed to read")
	})
	t.Run("no profile found", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(getSampleConfig(), nil)
		ctrl := New(mockConfigCtrl)
		_, ok, err := ctrl.GetProfileForExecution("invalid")
		assert.NoError(t, err)
		assert.False(t, ok)
	})
}
func TestControllerCreateProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(entity.Config{}, nil)
		mockConfigCtrl.EXPECT().Write(getDefaultConfig()).Return(nil)
		ctrl := New(mockConfigCtrl)
		err := ctrl.CreateProfile(getDefaultConfig().Profiles[0])
		assert.NoError(t, err)
	})
	t.Run("config controller read failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(entity.Config{}, errors.New("failed to read"))
		ctrl := New(mockConfigCtrl)
		err := ctrl.CreateProfile(getDefaultConfig().Profiles[0])
		assert.EqualError(t, err, "failed to read")
	})
	t.Run("config controller write failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(entity.Config{}, nil)
		mockConfigCtrl.EXPECT().Write(getDefaultConfig()).Return(errors.New("failed to write"))
		ctrl := New(mockConfigCtrl)
		err := ctrl.CreateProfile(getDefaultConfig().Profiles[0])
		assert.EqualError(t, err, "failed to write")
	})
}

func TestControllerDeleteProfile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(getSampleConfig(), nil).Times(2)
		expectedConfig := getSampleConfig()
		expectedConfig.Profiles = []entity.Profile{expectedConfig.Profiles[1]}
		mockConfigCtrl.EXPECT().Write(expectedConfig).Return(nil)
		ctrl := New(mockConfigCtrl)
		err := ctrl.DeleteProfiles([]string{getSampleConfig().Profiles[0].Name})
		assert.NoError(t, err)
	})
	t.Run("failed to delete only invalid profiles", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(getSampleConfig(), nil).Times(2)
		expectedConfig := getSampleConfig()
		expectedConfig.Profiles = []entity.Profile{expectedConfig.Profiles[1]}
		mockConfigCtrl.EXPECT().Write(expectedConfig).Return(nil)
		ctrl := New(mockConfigCtrl)
		err := ctrl.DeleteProfiles([]string{getSampleConfig().Profiles[0].Name, "invalid-profile1", "invalid-profile2"})
		assert.EqualError(t, err, "no profiles found for: invalid-profile1, invalid-profile2")
	})
	t.Run("config controller read failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(entity.Config{}, errors.New("failed to read"))
		ctrl := New(mockConfigCtrl)
		err := ctrl.DeleteProfiles([]string{getSampleConfig().Profiles[0].Name})
		assert.EqualError(t, err, "failed to read")
	})
	t.Run("config controller write failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockConfigCtrl := config.NewMockController(mockCtrl)
		mockConfigCtrl.EXPECT().Read().Return(getSampleConfig(), nil).Times(2)
		expectedConfig := getSampleConfig()
		expectedConfig.Profiles = []entity.Profile{expectedConfig.Profiles[1]}
		mockConfigCtrl.EXPECT().Write(expectedConfig).Return(errors.New("failed to write"))
		ctrl := New(mockConfigCtrl)
		err := ctrl.DeleteProfiles([]string{getSampleConfig().Profiles[0].Name})
		assert.EqualError(t, err, "failed to write")
	})
}
