//go:build windows
// +build windows

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

import "fmt"

func createDefaultConfigFileIfNotExists() error {
	return fmt.Errorf("creating default config file is not supported for windows. Please create manually")
}

func checkConfigFilePermission(configFilePath string) error {
	// since windows doesn't support create default config file, no validation is required
	return nil
}
