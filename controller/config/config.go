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
	"io/ioutil"
	"opensearch-cli/entity"
	"os"

	"gopkg.in/yaml.v3"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mocks/mock_config.go -package=mocks . Controller
type Controller interface {
	Read() (entity.Config, error)
	Write(config entity.Config) error
}

type controller struct {
	path string
}

//Read deserialize config file into entity.Config
func (c controller) Read() (result entity.Config, err error) {
	contents, err := ioutil.ReadFile(c.path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(contents, &result)
	return
}

//Write serialize entity.Config into file path
func (c controller) Write(config entity.Config) (err error) {
	file, err := os.Create(c.path) //overwrite if file exists
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()
	contents, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	_, err = file.Write(contents)
	if err != nil {
		return err
	}
	return file.Sync()
}

//New returns config controller instance
func New(path string) Controller {
	return controller{
		path: path,
	}
}
