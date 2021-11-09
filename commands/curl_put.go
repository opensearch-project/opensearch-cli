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

const curlPutCommandName = "put"

var curlPutExample = `
# Create a knn index from mapping setting saved in file "knn-mapping.json"
opensearch-cli curl put --path "my-knn-index"  \
                  --data "@some-location/knn-mapping.json" \
                  --pretty

# Update cluster settings transiently
opensearch-cli curl put --path "_cluster/settings" \
                  --query-params "flat_settings=true"  \
                  --data '
                    {
                      "transient" : {
                        "indices.recovery.max_bytes_per_sec" : "20mb"
                      }
                    }' \
                  --pretty

`

var curlPutCmd = &cobra.Command{
	Use:     curlPutCommandName + " [flags] ",
	Short:   "Put command to execute requests against cluster",
	Long:    "Put command enables you to run any PUT API against cluster",
	Example: curlPutExample,
	Run: func(cmd *cobra.Command, args []string) {
		Run(*cmd, curlPutCommandName)
	},
}

func init() {
	GetCurlCommand().AddCommand(curlPutCmd)
	curlPutCmd.Flags().StringP(curlPathFlagName, "P", "", "URL path for the REST API")
	_ = curlPutCmd.MarkFlagRequired(curlPathFlagName)
	curlPutCmd.Flags().StringP(curlQueryParamsFlagName, "q", "",
		"URL query parameters (key & value) for the REST API. Use ‘&’ to separate multiple parameters. Ex: -q \"v=true&s=order:desc,index_patterns\"")
	curlPutCmd.Flags().StringP(
		curlDataFlagName, "d", "",
		"Data for the REST API. If value starts with '@', the rest should be a file name to read the data from.")
	curlPutCmd.Flags().StringP(
		curlHeadersFlagName, "H", "",
		"Headers for the REST API. Consists of case-insensitive name followed by a colon (`:`), then by its value. Use ';' to separate multiple parameters. Ex: -H \"content-type:json;accept-encoding:gzip\"")
	curlPutCmd.Flags().BoolP("help", "h", false, "Help for curl "+curlPutCommandName)
}
