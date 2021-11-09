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

const curlPostCommandName = "post"

var postExample = `
# change the allocation of shards in a cluster.
opensearch-cli curl post --path "_cluster/reroute" \
                 --data '
                 {
                    "commands": [
                    {
                        "move": {
                            "index": "opensearch-cli", "shard": 0,
                            "from_node": "node1", "to_node": "node2"
                        }
                    },
                    {
                        "allocate_replica": {
                            "index": "test", "shard": 1,
                            "node": "node3"
                        }
                    }
                ]}' \
				--pretty

# insert a document to an index 
opensearch-cli curl post --path "my-index-01/_doc" \
                   --data '
                    {
                        "message": "insert document",
                        "ip": {
                            "address": "127.0.0.1"
                        }
                    }'

`
var curlPostCmd = &cobra.Command{
	Use:     curlPostCommandName + " [flags] ",
	Short:   "Post command to execute requests against cluster",
	Long:    "Post command enables you to run any POST API against cluster",
	Example: postExample,
	Run: func(cmd *cobra.Command, args []string) {
		Run(*cmd, curlPostCommandName)
	},
}

func init() {
	GetCurlCommand().AddCommand(curlPostCmd)
	curlPostCmd.Flags().StringP(curlPathFlagName, "P", "", "URL path for the REST API")
	_ = curlPostCmd.MarkFlagRequired(curlPathFlagName)
	curlPostCmd.Flags().StringP(curlQueryParamsFlagName, "q", "",
		"URL query parameters (key & value) for the REST API. Use ‘&’ to separate multiple parameters. Ex: -q \"v=true&s=order:desc,index_patterns\"")
	curlPostCmd.Flags().StringP(
		curlDataFlagName, "d", "",
		"Data for the REST API. If value starts with '@', the rest should be a file name to read the data from.")
	curlPostCmd.Flags().StringP(
		curlHeadersFlagName, "H", "",
		"Headers for the REST API. Consists of case-insensitive name followed by a colon (`:`), then by its value. Use ';' to separate multiple parameters. Ex: -H \"content-type:json;accept-encoding:gzip\"")
	curlPostCmd.Flags().BoolP("help", "h", false, "Help for curl "+curlPostCommandName)
}
