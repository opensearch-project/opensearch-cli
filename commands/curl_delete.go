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
	"github.com/spf13/cobra"
)

const curlDeleteCommandName = "delete"

var curlDeleteExample = `
# Delete a document from an index. 
opensearch-cli curl delete --path         "my-index/_doc/1" \
                     --query-params "routing=node1"
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
