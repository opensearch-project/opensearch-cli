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

package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func TestGetStatistics(t *testing.T) {
	t.Run("test stats command arguments", func(t *testing.T) {
		rootCmd := GetRoot()
		knnCommand := GetKNNCommand()
		knnStatsCmd := GetKNNStatsCommand()
		knnCommand.AddCommand(knnStatsCmd)
		rootCmd.AddCommand(knnCommand)
		_, err := executeCommand(rootCmd, knnCommandName, knnStatsCommandName, "--nodes", "node1,node2", "--stat-names", "stat1")
		assert.NoError(t, err)
		statNames, err := knnStatsCmd.Flags().GetString(knnStatsNamesFlagName)
		assert.NoError(t, err)
		assert.EqualValues(t, "stat1", statNames)
		nodeNames, err := knnStatsCmd.Flags().GetString(knnStatsNodesFlagName)
		assert.NoError(t, err)
		assert.EqualValues(t, "node1,node2", nodeNames)
	})
}

func TestWarmupIndices(t *testing.T) {
	t.Run("test warmup command failed", func(t *testing.T) {
		rootCmd := GetRoot()
		knnCommand := GetKNNCommand()
		knnWarmupCmd := GetKNNWarmupCommand()
		knnCommand.AddCommand(knnWarmupCmd)
		rootCmd.AddCommand(knnCommand)
		_, err := executeCommand(rootCmd, knnCommandName, knnWarmupCommandName)
		assert.Error(t, err)
	})
	t.Run("test warmup command", func(t *testing.T) {
		rootCmd := GetRoot()
		knnCommand := GetKNNCommand()
		knnWarmupCmd := GetKNNWarmupCommand()
		knnCommand.AddCommand(knnWarmupCmd)
		rootCmd.AddCommand(knnCommand)
		result, err := executeCommand(rootCmd, knnCommandName, knnWarmupCommandName, "index1", "index2")
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}
