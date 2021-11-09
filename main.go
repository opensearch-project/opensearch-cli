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

// opensearch-cli is an unified command line tool for OpenSearch clusters
package main

import (
	"opensearch-cli/commands"
	"os"
)

func main() {
	if err := commands.Execute(); err != nil {
		// By default every command should handle their error message
		os.Exit(1)
	}
}
