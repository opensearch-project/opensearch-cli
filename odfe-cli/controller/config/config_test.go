/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package config

import (
	"es-cli/odfe-cli/entity"
	"fmt"
	"io/ioutil"
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
