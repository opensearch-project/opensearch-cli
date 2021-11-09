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

package config

import (
	"fmt"
	"io/ioutil"
	"opensearch-cli/entity"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/assert"
)

const testFileName = "config.yaml"
const testFolderName = "testdata"

func getSampleConfig() entity.Config {
	return entity.Config{
		Profiles: []entity.Profile{
			{
				Name:     "local",
				Endpoint: "https://localhost:9200",
				UserName: "admin", Password: "admin",
			},
			{
				Name:     "default",
				Endpoint: "https://127.0.0.1:9200",
				UserName: "dadmin", Password: "dadmin",
			},
		}}
}

func TestControllerRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := New(filepath.Join(testFolderName, testFileName))
		cfg, err := ctrl.Read()
		assert.NoError(t, err)
		expected := getSampleConfig()
		assert.EqualValues(t, expected, cfg)
	})
	t.Run("fail", func(t *testing.T) {
		fileName := filepath.Join(testFolderName, "invalid", testFileName)
		ctrl := New(fileName)
		_, err := ctrl.Read()
		assert.EqualError(t, err, fmt.Sprintf("open %s: no such file or directory", fileName))
	})
}
func TestControllerWrite(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		f, err := ioutil.TempFile("", "config")
		assert.NoError(t, err)
		defer func() {
			err = os.Remove(f.Name())
			assert.NoError(t, err)
		}()
		ctrl := New(f.Name())
		err = ctrl.Write(getSampleConfig())
		assert.NoError(t, err)
		contents, err := ioutil.ReadFile(f.Name())
		assert.NoError(t, err)
		var config entity.Config
		err = yaml.Unmarshal(contents, &config)
		assert.NoError(t, err)
		assert.EqualValues(t, getSampleConfig(), config)
	})
}
