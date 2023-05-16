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
	"io/ioutil"
	"opensearch-cli/entity"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigFilePath(t *testing.T) {

	t.Run("config file path from environment variable", func(t *testing.T) {
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

func createTempConfigFile(testFilePath string) (*os.File, error) {
	content, err := ioutil.ReadFile(testFilePath)
	if err != nil {
		return nil, err
	}
	tmpfile, err := ioutil.TempFile(os.TempDir(), "test-file")
	if err != nil {
		return nil, err
	}
	if _, err := tmpfile.Write(content); err != nil {
		os.Remove(tmpfile.Name()) // clean up
		return nil, err
	}
	if runtime.GOOS == "windows" {
		return tmpfile, nil
	}
	if err := tmpfile.Chmod(0600); err != nil {
		os.Remove(tmpfile.Name()) // clean up
		return nil, err
	}
	return tmpfile, nil
}

func TestGetProfile(t *testing.T) {
	t.Run("get default profile", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		profileFile, err := createTempConfigFile("testdata/config.yaml")
		assert.NoError(t, err)
		filePath, err := filepath.Abs(profileFile.Name())
		assert.NoError(t, err)
		root.SetArgs([]string{"--config", filePath})
		_, err = root.ExecuteC()
		assert.NoError(t, err)
		actual, err := GetProfile()
		assert.NoError(t, err)
		expectedProfile := entity.Profile{Name: "default", Endpoint: "http://localhost:9200", UserName: "default", Password: "admin"}
		assert.EqualValues(t, expectedProfile, *actual)
		os.Remove(profileFile.Name())
	})
	t.Run("test get profile", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		profileFile, err := createTempConfigFile("testdata/config.yaml")
		assert.NoError(t, err)
		filePath, err := filepath.Abs(profileFile.Name())
		assert.NoError(t, err)
		root.SetArgs([]string{"--config", filePath, "--profile", "test"})
		_, err = root.ExecuteC()
		assert.NoError(t, err)
		actual, err := GetProfile()
		assert.NoError(t, err)
		expectedProfile := entity.Profile{Name: "test", Endpoint: "https://localhost:9200", UserName: "admin", Password: "admin"}
		assert.EqualValues(t, expectedProfile, *actual)
	})
	t.Run("Profile mismatch", func(t *testing.T) {
		root := GetRoot()
		assert.NotNil(t, root)
		profileFile, err := createTempConfigFile("testdata/config.yaml")
		assert.NoError(t, err)
		filePath, err := filepath.Abs(profileFile.Name())
		assert.NoError(t, err)
		root.SetArgs([]string{"--config", filePath, "--profile", "test1"})
		_, err = root.ExecuteC()
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
		assert.EqualError(t, err, "failed to get config file info due to: stat testdata/config1.yaml: no such file or directory", "unexpected error")
	})
	t.Run("invalid config file permission", func(t *testing.T) {

		if runtime.GOOS == "windows" {
			t.Skipf("test case does not work on %s", runtime.GOOS)
		}
		root := GetRoot()
		assert.NotNil(t, root)
		profileFile, err := createTempConfigFile("testdata/config.yaml")
		assert.NoError(t, err)
		assert.NoError(t, profileFile.Chmod(0750))
		filePath, err := filepath.Abs(profileFile.Name())
		assert.NoError(t, err)
		root.SetArgs([]string{"--config", filePath, "--profile", "test"})
		_, err = root.ExecuteC()
		assert.NoError(t, err)
		_, err = GetProfile()
		assert.True(t, strings.Contains(err.Error(), "permissions 750"), "unexpected error")
	})
}
