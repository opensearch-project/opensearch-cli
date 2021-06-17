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
	"fmt"
	"opensearch-cli/entity"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	configFileType        = "yaml"
	defaultConfigFileName = "config"
	flagConfig            = "config"
	flagProfileName       = "profile"
	folderPermission      = 0755 // only owner can write, while everyone can read and execute
	ConfigEnvVarName      = "OPENSEARCH_CLI_CONFIG"
	RootCommandName       = "opensearch-cli"
	version               = "1.0.0"
)

func buildVersionString() string {

	return fmt.Sprintf("%s %s/%s", version, runtime.GOOS, runtime.GOARCH)
}

var rootCommand = &cobra.Command{
	Use:     RootCommandName,
	Short:   "opensearch-cli is a unified command line interface for managing OpenSearch clusters",
	Version: buildVersionString(),
}

func GetRoot() *cobra.Command {
	return rootCommand
}

// Execute executes the root command.
func Execute() error {
	err := rootCommand.Execute()
	return err
}

func GetDefaultConfigFilePath() string {
	return filepath.Join(
		getDefaultConfigFolderRootPath(),
		fmt.Sprintf(".%s", RootCommandName),
		fmt.Sprintf("%s.%s", defaultConfigFileName, configFileType),
	)
}

func getDefaultConfigFolderRootPath() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir
	}
	if cwd, err := os.Getwd(); err == nil {
		return cwd
	}
	return ""
}

func init() {
	cobra.OnInitialize()
	configFilePath := GetDefaultConfigFilePath()
	rootCommand.PersistentFlags().StringP(flagConfig, "c", "", fmt.Sprintf("Configuration file for opensearch-cli, default is %s", configFilePath))
	rootCommand.PersistentFlags().StringP(flagProfileName, "p", "", "Use a specific profile from your configuration file")
	rootCommand.Flags().BoolP("version", "v", false, "Version for opensearch-cli")
	rootCommand.Flags().BoolP("help", "h", false, "Help for opensearch-cli")
}

// GetConfigFilePath gets config file path for execution
func GetConfigFilePath(configFlagValue string) (string, error) {

	if configFlagValue != "" {
		return configFlagValue, nil
	}
	if value, ok := os.LookupEnv(ConfigEnvVarName); ok {
		return value, nil
	}
	if err := createDefaultConfigFileIfNotExists(); err != nil {
		return "", err
	}
	return GetDefaultConfigFilePath(), nil
}

// createDefaultConfigFolderIfNotExists creates default config file along with folder if
// it doesn't exists
func createDefaultConfigFileIfNotExists() error {
	defaultFilePath := GetDefaultConfigFilePath()
	if isExists(defaultFilePath) {
		return nil
	}
	folderPath := filepath.Dir(defaultFilePath)
	if !isExists(folderPath) {
		err := os.Mkdir(folderPath, folderPermission)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(defaultFilePath)
	if err != nil {
		return err
	}
	return f.Close()
}

//isExists check if given path exists or not
//if path is just a name, it will check in current directory
func isExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// DisplayError prints command name and error on console and exists as well.
func DisplayError(err error, cmdName string) {
	if err != nil {
		fmt.Println(cmdName, "Command failed.")
		fmt.Println("Reason:", err)
	}
}

// GetProfile gets profile details for current execution
func GetProfile() (*entity.Profile, error) {
	p, err := GetProfileController()
	if err != nil {
		return nil, err
	}
	profileFlagValue, err := rootCommand.PersistentFlags().GetString(flagProfileName)
	if err != nil {
		return nil, err
	}
	profile, ok, err := p.GetProfileForExecution(profileFlagValue)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("no profile found for execution. Try %s %s --help for more information", RootCommandName, ProfileCommandName)
	}
	return &profile, nil
}
