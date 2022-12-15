// +build !windows

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
	"fmt"
	"os"
	"path/filepath"
)

const (
	FolderPermission = 0700 // only owner can read, write and execute
	FilePermission   = 0600 // only owner can read and write
)

// createDefaultConfigFolderIfNotExists creates default config file along with folder if
// it doesn't exists
func createDefaultConfigFileIfNotExists() error {
	defaultFilePath := GetDefaultConfigFilePath()
	if isExists(defaultFilePath) {
		return nil
	}
	folderPath := filepath.Dir(defaultFilePath)
	if !isExists(folderPath) {
		err := os.Mkdir(folderPath, FolderPermission)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(defaultFilePath)
	if err != nil {
		return err
	}
	if err = f.Chmod(FilePermission); err != nil {
		return err
	}
	return f.Close()
}

func checkConfigFilePermission(configFilePath string) error {
	//check for config file permission
	info, err := os.Stat(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to get config file info due to: %w", err)
	}
	mode := info.Mode().Perm()

	if mode != FilePermission {
		return fmt.Errorf("config file '%s' permissions %o must be changed to %o, this will limit access to only your user", configFilePath, mode, FilePermission)
	}
	return nil
}
