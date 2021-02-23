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
	"github.com/spf13/cobra"
)

const curlGetCommandName = "get"

var curlGetExample = `
# get document count for an index
odfe-cli curl get --path "_cat/count/my-index-01" --query-params "v=true" --pretty

# get health status of a cluster.
odfe-cli curl get --path "_cluster/health" --pretty --filter-path "status"

# get explanation for cluster allocation for a given index and shard number
odfe-cli curl get --path "_cluster/allocation/explain" \
                  --data '{
                    "index": "my-index-01",
                    "shard": 0,
                    "primary": false,
                    "current_node": "nodeA"                         
                  }'
`

var curlGetCmd = &cobra.Command{
	Use:     curlGetCommandName + " [flags] ",
	Short:   "Get command to execute requests against cluster",
	Long:    "Get command enables you to run any GET API against cluster",
	Example: curlGetExample,
	Run: func(cmd *cobra.Command, args []string) {
		Run(*cmd, curlGetCommandName)
	},
}

func init() {
	GetCurlCommand().AddCommand(curlGetCmd)
	curlGetCmd.Flags().StringP(curlPathFlagName, "P", "", "URL path for the REST API")
	_ = curlGetCmd.MarkFlagRequired(curlPathFlagName)
	curlGetCmd.Flags().StringP(curlQueryParamsFlagName, "q", "",
		"URL query parameters (key & value) for the REST API. Use ‘&’ to separate multiple parameters. Ex: -q \"v=true&s=order:desc,index_patterns\"")
	curlGetCmd.Flags().StringP(
		curlDataFlagName, "d", "",
		"Data for the REST API. If value starts with '@', the rest should be a file name to read the data from.")
	curlGetCmd.Flags().StringP(
		curlHeadersFlagName, "H", "",
		"Headers for the REST API. Consists of case-insensitive name followed by a colon (`:`), then by its value. Use ';' to separate multiple parameters. Ex: -H \"content-type:json;accept-encoding:gzip\"")
	curlGetCmd.Flags().BoolP("help", "h", false, "Help for curl "+curlGetCommandName)
}
