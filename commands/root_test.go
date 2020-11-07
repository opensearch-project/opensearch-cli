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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigFilePath(t *testing.T) {

	t.Run("config file path from os environment variable", func(t *testing.T) {
		err := os.Setenv(odfeConfigEnvVarName, "test/config.yml")
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
		expected, err := cmd.Flags().GetString(flagConfig)
		assert.NoError(t, err)
		assert.EqualValues(t, expected, "test/config.yml")
	})
}
