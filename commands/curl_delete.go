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

const curlDeleteCommandName = "delete"

var curlDeleteExample = `
# Delete a document from an index. 
odfe-cli curl delete --path         "my-index/_doc/1" \
                     --query-params "routing=odfe-node1"
`

var curlDeleteCmd = &cobra.Command{
	Use:     curlDeleteCommandName + " [flags] ",
	Short:   "Delete command to execute requests against cluster",
	Long:    "Delete command enables you to run any DELETE API against cluster",
	Example: curlDeleteExample,
	Run: func(cmd *cobra.Command, args []string) {
		Run(*cmd, curlDeleteCommandName)
	},
}

func init() {
	GetCurlCommand().AddCommand(curlDeleteCmd)
	curlDeleteCmd.Flags().StringP(curlPathFlagName, "P", "", "URL path for the REST API")
	_ = curlDeleteCmd.MarkFlagRequired(curlPathFlagName)
	curlDeleteCmd.Flags().StringP(curlQueryParamsFlagName, "q", "",
		"URL query parameters (key & value) for the REST API. Use ‘&’ to separate multiple parameters. Ex: -q \"v=true&s=order:desc,index_patterns\"")
	curlDeleteCmd.Flags().StringP(
		curlHeadersFlagName, "H", "",
		"Headers for the REST API. Consists of case-insensitive name followed by a colon (`:`), then by its value. Use ';' to separate multiple parameters. Ex: -H \"content-type:json;accept-encoding:gzip\"")
	curlDeleteCmd.Flags().BoolP("help", "h", false, "Help for curl "+curlDeleteCommandName)
}
