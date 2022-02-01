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

const curlPatchCommandName = "patch"

var patchExample = `
# creates or updates multiple role mappings in a single call.
opensearch-cli curl patch --path "_plugins/_security/api/rolesmapping" \
                 --data '
								 [
								   {
								     "op": "add", "path": "/human_resources", "value": { "users": ["user1"], "backend_roles": ["backendrole2"] }
								   },
								   {
								     "op": "add", "path": "/finance", "value": { "users": ["user2"], "backend_roles": ["backendrole2"] }
								   }
								 ]' \
				--pretty

# add, delete, or modify multiple tenants in a single call.
opensearch-cli curl patch --path "_plugins/_security/api/tenants/" \
                   --data '
									 [
									   {
									     "op": "replace",
									     "path": "/human_resources/description",
									     "value": "An updated description"
									   },
									   {
									     "op": "add",
									     "path": "/another_tenant",
									     "value": {
									       "description": "Another description."
									     }
									   }
									 ]'

`
var curlPatchCmd = &cobra.Command{
	Use:     curlPatchCommandName + " [flags] ",
	Short:   "Patch command to execute requests against cluster",
	Long:    "Patch command enables you to run any PATCH API against cluster",
	Example: patchExample,
	Run: func(cmd *cobra.Command, args []string) {
		Run(*cmd, curlPatchCommandName)
	},
}

func init() {
	GetCurlCommand().AddCommand(curlPatchCmd)
	curlPatchCmd.Flags().StringP(curlPathFlagName, "P", "", "URL path for the REST API")
	_ = curlPatchCmd.MarkFlagRequired(curlPathFlagName)
	curlPatchCmd.Flags().StringP(curlQueryParamsFlagName, "q", "",
		"URL query parameters (key & value) for the REST API. Use ‘&’ to separate multiple parameters. Ex: -q \"v=true&s=order:desc,index_patterns\"")
	curlPatchCmd.Flags().StringP(
		curlDataFlagName, "d", "",
		"Data for the REST API. If value starts with '@', the rest should be a file name to read the data from.")
	curlPatchCmd.Flags().StringP(
		curlHeadersFlagName, "H", "",
		"Headers for the REST API. Consists of case-insensitive name followed by a colon (`:`), then by its value. Use ';' to separate multiple parameters. Ex: -H \"content-type:json;accept-encoding:gzip\"")
	curlPatchCmd.Flags().BoolP("help", "h", false, "Help for curl "+curlPatchCommandName)
}
