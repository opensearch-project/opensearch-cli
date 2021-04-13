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
	"fmt"
	"opensearch-cli/client"
	ctrl "opensearch-cli/controller/platform"
	entity "opensearch-cli/entity/platform"
	gateway "opensearch-cli/gateway/platform"
	handler "opensearch-cli/handler/platform"

	"github.com/spf13/cobra"
)

const (
	curlCommandName              = "curl"
	curlPrettyFlagName           = "pretty"
	curlPathFlagName             = "path"
	curlQueryParamsFlagName      = "query-params"
	curlDataFlagName             = "data"
	curlHeadersFlagName          = "headers"
	curlOutputFormatFlagName     = "output-format"
	curlOutputFilterPathFlagName = "filter-path"
)

//curlCommand is base command for OpenSearch REST APIs.
var curlCommand = &cobra.Command{
	Use:   curlCommandName,
	Short: "Manage OpenSearch platform features",
	Long:  "Use the curl command to execute any REST API calls against the cluster.",
}

func init() {
	curlCommand.Flags().BoolP("help", "h", false, "Help for curl command")
	curlCommand.PersistentFlags().Bool(curlPrettyFlagName, false, "Response will be formatted")
	curlCommand.PersistentFlags().StringP(curlOutputFormatFlagName, "o", "",
		"Output format if supported by cluster, else, default format by OpenSearch. Example json, yaml")
	curlCommand.PersistentFlags().StringP(curlOutputFilterPathFlagName, "f", "",
		"Filter output fields returned by OpenSearch. Use comma ',' to separate list of filters")
	GetRoot().AddCommand(curlCommand)
}

//GetCurlCommand returns Curl base command, since this will be needed for subcommands
//to add as parent later
func GetCurlCommand() *cobra.Command {
	return curlCommand
}

//getCurlHandler returns handler by wiring the dependency manually
func getCurlHandler() (*handler.Handler, error) {
	c, err := client.New(nil)
	if err != nil {
		return nil, err
	}
	profile, err := GetProfile()
	if err != nil {
		return nil, err
	}
	g := gateway.New(c, profile)
	facade := ctrl.New(g)
	return handler.New(facade), nil
}

//CurlActionExecute executes API based on user request
func CurlActionExecute(input entity.CurlCommandRequest) error {

	commandHandler, err := getCurlHandler()
	if err != nil {
		return err
	}
	response, err := handler.Curl(commandHandler, input)
	if err == nil {
		fmt.Println(string(response))
		return nil
	}
	if requestError, ok := err.(*entity.RequestError); ok {
		fmt.Println(requestError.GetResponse())
		return nil
	}
	return err
}

func FormatOutput() bool {
	isPretty, _ := curlCommand.PersistentFlags().GetBool(curlPrettyFlagName)
	return isPretty
}

func GetUserInputAsStringForFlag(flagName string) string {
	format, _ := curlCommand.PersistentFlags().GetString(flagName)
	return format
}

func Run(cmd cobra.Command, cmdName string) {
	input := entity.CurlCommandRequest{
		Action:           cmdName,
		Pretty:           FormatOutput(),
		OutputFormat:     GetUserInputAsStringForFlag(curlOutputFormatFlagName),
		OutputFilterPath: GetUserInputAsStringForFlag(curlOutputFilterPathFlagName),
	}
	input.Path, _ = cmd.Flags().GetString(curlPathFlagName)
	input.QueryParams, _ = cmd.Flags().GetString(curlQueryParamsFlagName)
	input.Data, _ = cmd.Flags().GetString(curlDataFlagName)
	input.Headers, _ = cmd.Flags().GetString(curlHeadersFlagName)
	err := CurlActionExecute(input)
	DisplayError(err, cmdName)
}
