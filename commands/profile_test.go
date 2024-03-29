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

package commands

import (
	"errors"
	"fmt"
	"opensearch-cli/controller/profile/mocks"
	"opensearch-cli/entity"
	"os"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func fakeInputProfile() entity.Profile {
	return entity.Profile{
		Name:     "default",
		Endpoint: "localhost:9200",
		UserName: "admin",
		Password: "admin",
	}
}

func fakeInSecuredInputProfile() entity.Profile {
	return entity.Profile{
		Name:     "default",
		Endpoint: "localhost:9200",
	}
}

func fakeAWSIAMInputProfile() entity.Profile {
	return entity.Profile{
		Name:     "default",
		Endpoint: "localhost:9200",
		AWS: &entity.AWSIAM{
			ProfileName: "iam-test",
		},
	}
}

func TestCreateProfile(t *testing.T) {
	t.Run("create profile successfully", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockProfileCtrl := mocks.NewMockController(mockCtrl)
		mockProfileCtrl.EXPECT().CreateProfile(fakeInputProfile()).Return(nil)
		err := CreateProfile(mockProfileCtrl, fakeInputProfile())
		assert.NoError(t, err)
	})
	t.Run("create security disabled profile successfully", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockProfileCtrl := mocks.NewMockController(mockCtrl)
		mockProfileCtrl.EXPECT().CreateProfile(fakeInSecuredInputProfile()).Return(nil)
		err := CreateProfile(mockProfileCtrl, fakeInSecuredInputProfile())
		assert.NoError(t, err)
	})
	t.Run("create profile using aws iam successfully", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockProfileCtrl := mocks.NewMockController(mockCtrl)
		mockProfileCtrl.EXPECT().CreateProfile(fakeAWSIAMInputProfile()).Return(nil)
		err := CreateProfile(mockProfileCtrl, fakeAWSIAMInputProfile())
		assert.NoError(t, err)
	})
	t.Run("create profile failed", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockProfileCtrl := mocks.NewMockController(mockCtrl)
		mockProfileCtrl.EXPECT().CreateProfile(fakeInputProfile()).Return(errors.New("error"))
		err := CreateProfile(mockProfileCtrl, fakeInputProfile())
		assert.EqualError(t, err, fmt.Sprintf("failed to create profile %v due to: error", fakeInputProfile()))
	})
	t.Run("check mandatory create parameters are provided", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{
			ProfileCommandName, CreateNewProfileCommandName,
		})
		_, err := root.ExecuteC()
		assert.EqualErrorf(t, err, "required flag(s) \"auth-type\", \"endpoint\", \"name\" not set", "unexpected error")
	})
	t.Run("create security disabled profile", func(t *testing.T) {
		f, err := os.CreateTemp("", "profile")
		assert.NoError(t, err)
		defer func() {
			err := os.Remove(f.Name())
			assert.NoError(t, err)
		}()
		root := GetRoot()
		assert.NotNil(t, root)
		testProfileName := "pname"
		testProfileEndpoint := "some-endpoint"
		root.SetArgs([]string{
			ProfileCommandName, CreateNewProfileCommandName,
			"--" + flagConfig, f.Name(),
			"--" + FlagProfileCreateAuthType, "disabled",
			"--" + FlagProfileCreateEndpoint, testProfileEndpoint,
			"--" + FlagProfileCreateName, testProfileName,
			"--" + FlagProfileMaxRetry, "2",
			"--" + FlagProfileTimeout, "10",
		})
		_, err = root.ExecuteC()
		assert.NoError(t, err)
		contents, _ := os.ReadFile(f.Name())
		var actual entity.Config
		assert.NoError(t, yaml.Unmarshal(contents, &actual))
		retryVal := 2
		timeout := int64(10)
		assert.EqualValues(t, []entity.Profile{
			{
				Name:     testProfileName,
				Endpoint: testProfileEndpoint,
				MaxRetry: &retryVal,
				Timeout:  &timeout,
			},
		}, actual.Profiles)

	})
}

func TestDeleteProfileCommand(t *testing.T) {
	t.Run("test delete profile command", func(t *testing.T) {
		f, err := os.CreateTemp("", "profile-delete")
		assert.NoError(t, err)
		defer func() {
			err := os.Remove(f.Name())
			assert.NoError(t, err)
		}()
		config := entity.Config{Profiles: []entity.Profile{fakeInputProfile()}}
		bytes, err := yaml.Marshal(config)
		assert.NoError(t, err)
		assert.NoError(t, os.WriteFile(f.Name(), bytes, 0644))
		assert.NoError(t, f.Sync())
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{ProfileCommandName, DeleteProfilesCommandName, config.Profiles[0].Name, "--config", f.Name()})
		cmd, err := root.ExecuteC()
		assert.NoError(t, err)
		expected, err := cmd.Flags().GetString(flagConfig)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, f.Name())
		var expectedConfig entity.Config
		contents, err := os.ReadFile(f.Name())
		assert.NoError(t, err)
		err = yaml.Unmarshal(contents, &expectedConfig)
		assert.NoError(t, err)
		assert.EqualValues(t, expectedConfig, entity.Config{Profiles: []entity.Profile{}})
	})
}

func TestListsProfileCommand(t *testing.T) {
	t.Run("list profiles", func(t *testing.T) {
		f, err := os.CreateTemp("", "profile-list")
		assert.NoError(t, err)
		defer func() {
			err := os.Remove(f.Name())
			assert.NoError(t, err)
		}()
		config := entity.Config{Profiles: []entity.Profile{fakeInputProfile()}}
		bytes, err := yaml.Marshal(config)
		assert.NoError(t, err)
		assert.NoError(t, os.WriteFile(f.Name(), bytes, 0644))
		assert.NoError(t, f.Sync())
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{ProfileCommandName, ListProfilesCommandName, "--config", f.Name()})
		cmd, err := root.ExecuteC()
		assert.NoError(t, err)
		expected, err := cmd.Flags().GetString(flagConfig)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, f.Name())
	})
	t.Run("list profiles with verbose", func(t *testing.T) {
		f, err := os.CreateTemp("", "profile-list")
		assert.NoError(t, err)
		defer func() {
			err := os.Remove(f.Name())
			assert.NoError(t, err)
		}()
		config := entity.Config{Profiles: []entity.Profile{fakeInputProfile()}}
		bytes, err := yaml.Marshal(config)
		assert.NoError(t, err)
		assert.NoError(t, os.WriteFile(f.Name(), bytes, 0644))
		assert.NoError(t, f.Sync())
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{ProfileCommandName, ListProfilesCommandName, "--" + FlagProfileVerbose, "--" + flagConfig, f.Name()})
		cmd, err := root.ExecuteC()
		assert.NoError(t, err)
		expected, err := cmd.Flags().GetString(flagConfig)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, f.Name())
	})
	t.Run("no profiles found", func(t *testing.T) {
		f, err := os.CreateTemp("", "profile")
		assert.NoError(t, err)
		defer func() {
			err := os.Remove(f.Name())
			assert.NoError(t, err)
		}()
		config := entity.Config{Profiles: []entity.Profile{}}
		bytes, err := yaml.Marshal(config)
		assert.NoError(t, err)
		assert.NoError(t, os.WriteFile(f.Name(), bytes, 0644))
		assert.NoError(t, f.Sync())
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{ProfileCommandName, ListProfilesCommandName, "--" + flagConfig, f.Name()})
		cmd, err := root.ExecuteC()
		assert.NoError(t, err)
		expected, err := cmd.Flags().GetString(flagConfig)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, f.Name())
	})
}
