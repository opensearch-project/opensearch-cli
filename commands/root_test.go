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

package commands

import (
	"opensearch-cli/entity"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigFilePath(t *testing.T) {

	t.Run("config file path from os environment variable", func(t *testing.T) {
		err := os.Setenv(ConfigEnvVarName, "test/config.yml")
		assert.NoError(t, err)
		filePath, err := GetConfigFilePath("")
		assert.NoError(t, err)
		assert.EqualValues(t, "test/config.yml", filePath)
	})
	t.Run("config file path from command line arguments", func(t *testing.T) {
		filePath, err := GetConfigFilePath("test/config.yml")
		assert.NoError(t, err)
		assert.EqualValues(t, "test/config.yml", filePath)
	})
}

func TestGetRoot(t *testing.T) {
	t.Run("test root command", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{"--config", "test/config.yml"})
		cmd, err := root.ExecuteC()
		assert.NoError(t, err)
		actual, err := cmd.Flags().GetString(flagConfig)
		assert.NoError(t, err)
		assert.EqualValues(t, "test/config.yml", actual)
	})
}

func TestVersionString(t *testing.T) {
	t.Run("test version flag", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{"--version"})
		cmd, err := root.ExecuteC()
		assert.NoError(t, err)
		expected := "1.1.0 " + runtime.GOOS + "/" + runtime.GOARCH
		assert.EqualValues(t, expected, cmd.Version)
	})
}

func TestGetProfile(t *testing.T) {
	t.Run("get default profile", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{"--config", "testdata/config.yaml"})
		_, err := root.ExecuteC()
		assert.NoError(t, err)
		actual, err := GetProfile()
		assert.NoError(t, err)
		expectedProfile := entity.Profile{Name: "default", Endpoint: "http://localhost:9200", UserName: "default", Password: "admin"}
		assert.EqualValues(t, expectedProfile, *actual)
	})
	t.Run("test get profile", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{"--config", "testdata/config.yaml", "--profile", "test"})
		_, err := root.ExecuteC()
		assert.NoError(t, err)
		actual, err := GetProfile()
		assert.NoError(t, err)
		expectedProfile := entity.Profile{Name: "test", Endpoint: "https://localhost:9200", UserName: "admin", Password: "admin"}
		assert.EqualValues(t, expectedProfile, *actual)
	})
	t.Run("Profile mismatch", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{"--config", "testdata/config.yaml", "--profile", "test1"})
		_, err := root.ExecuteC()
		assert.NoError(t, err)
		_, err = GetProfile()
		assert.EqualErrorf(t, err, "profile 'test1' does not exist", "unexpected error")
	})
	t.Run("no config file found", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		root.SetArgs([]string{"--config", "testdata/config1.yaml", "--profile", "test1"})
		_, err := root.ExecuteC()
		assert.NoError(t, err)
		_, err = GetProfile()
		assert.EqualError(t, err, "open testdata/config1.yaml: no such file or directory", "unexpected error")
	})
}
