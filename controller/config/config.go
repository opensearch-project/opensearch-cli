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
	"io/ioutil"
	"odfe-cli/entity"
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
